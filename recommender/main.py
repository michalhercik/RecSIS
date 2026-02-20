import importlib.util
import inspect
import os
from typing import Any, Optional

import pandas as pd
from fastapi import FastAPI, Request
from fastapi.responses import JSONResponse
from pydantic import BaseModel

from algo.base import Algorithm, Recommendation
from data_repository import DataRepository
from user import User

pd.set_option("future.no_silent_downcasting", True)

app = FastAPI()


def discover_algorithm_classes(directory):
    classes = []
    for filename in os.listdir(directory):
        if filename.endswith(".py") and not filename.startswith("__"):
            module_name = filename[:-3]
            file_path = os.path.join(directory, filename)
            spec = importlib.util.spec_from_file_location(module_name, file_path)
            module = importlib.util.module_from_spec(spec)
            spec.loader.exec_module(module)
            for _, obj in inspect.getmembers(module, inspect.isclass):
                if issubclass(obj, Algorithm) and obj is not Algorithm:
                    classes.append(obj)
    return classes


def instantiate_classes(classes):
    instances = {}
    data = DataRepository()
    for cls in classes:
        try:
            instance = cls(data)
            instances[cls.__name__] = instance
        except Exception as e:
            print(f"Could not instantiate {cls.__name__}: {e}")
    return instances


def fit_algorithms(instances):
    for cls in instances.values():
        try:
            cls.fit()
        except NotImplementedError:
            print(f"WARNING: {cls.__class__.__name__}.fit not implemented")
        except Exception as e:
            print(f"Could not fit {cls.__class__.__name__}: {e}")


algorithm_classes = instantiate_classes(discover_algorithm_classes("algo"))
algorithm_fit = {}
# fit_algorithms(algorithm_classes)


class RecommendRequest(BaseModel):
    algo: Optional[str] = None
    limit: Optional[int] = 10
    user_id: Optional[Any] = None
    student: Optional[Any] = None
    degree_plan: Optional[Any] = None
    enrollment_year: Optional[Any] = None
    blueprint: Optional[Any] = None
    # model_config = {
    #     "json_schema_extra": {
    #         "examples": [
    #             {
    #                 "algo": "knn",
    #                 "limit": 10,
    #                 "user_id": "1234",
    #                 "degree_plan": "NIPVS19B",
    #                 "enrollment_year": 2020,
    #                 "blueprint": [
    #                     {"year": 0, "unassigned": []},
    #                     {"year": 1,
    #                         "summer": ["NSWI170", "NMAI054", "NMAI058", "NPRG031", "NSWI177", "NTIN060", "NTVY015"],
    #                         "winter": ["NPRG062", "NSWI120", "NMAI069", "NTVY014", "NSWI141", "NDMI050", "NDMI002", "NMAI057", "NPRG030"]},
    #                     {"year": 2, "summer": [], "winter": []},
    #                     {"year": 3, "summer": [], "winter": []}
    #                 ]
    #             }
    #         ]
    #     }
    # }


@app.post("/recommended")
async def recommended(req: RecommendRequest):
    algo_name = req.algo
    algo_class = None
    if algo_name in algorithm_classes:
        algo_class = algorithm_classes[algo_name]
    if algo_class:
        user = User(req.user_id, req.degree_plan, req.enrollment_year, req.blueprint)
        if req.student is not None:
            user.id = req.student
            user.fetch = True
        result = algo_class.recommend(user, req.limit)
    else:
        result = Recommendation(None, None)
    return JSONResponse(content={"recommended": result.rec, "target": result.target, "finished": result.finished, "expected": result.expected})


@app.post("/fit")
async def fit(req: RecommendRequest):
    algo_name = req.algo
    algo_class = None
    if algo_name in algorithm_classes:
        algo_class = algorithm_classes[algo_name]
    if algo_class:
        try:
            algo_class.fit()
            algorithm_fit[algo_name] = True
        except NotImplementedError:
            print(f"WARNING: {algo_class.__class__.__name__}.fit not implemented")
    return JSONResponse(content={})


@app.get("/algorithms")
async def get_algorithms():
    algos = list(algorithm_classes.keys())
    fitted = [True if algo in algorithm_fit else False for algo in algos ]
    return JSONResponse(content={"algorithms": algos, "fit": fitted})
