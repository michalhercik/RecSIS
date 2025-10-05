import os
import importlib.util
import inspect
import pandas as pd
from fastapi import FastAPI, Request
from fastapi.responses import JSONResponse
from pydantic import BaseModel
from typing import Optional, Any
from data_repository import DataRepository
from algo.base import Algorithm
from user import User

pd.set_option('future.no_silent_downcasting', True)

app = FastAPI()

def discover_algorithm_classes(directory):
    classes = []
    for filename in os.listdir(directory):
        if filename.endswith('.py') and not filename.startswith('__'):
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

algorithm_classes = instantiate_classes(discover_algorithm_classes('algo'))
fit_algorithms(algorithm_classes)

class RecommendRequest(BaseModel):
    algo: Optional[str] = None
    limit: Optional[int] = 10
    user_id: Optional[Any] = None
    degree_plan: Optional[Any] = None
    enrollment_year: Optional[Any] = None
    blueprint: Optional[Any] = None

@app.post("/recommended")
async def recommended(req: RecommendRequest):
    algo_name = req.algo
    algo_class = None
    if algo_name in algorithm_classes:
        algo_class = algorithm_classes[algo_name]
    if algo_class:
        user = User(req.user_id, req.degree_plan, req.enrollment_year, req.blueprint)
        result = algo_class.recommend(user, req.limit)
    else:
        result = None
    return JSONResponse(content={"recommended": result})

@app.get("/algorithms")
async def get_algorithms():
    return JSONResponse(content={"algorithms": list(algorithm_classes.keys())})
