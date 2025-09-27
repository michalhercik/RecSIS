
# FastAPI version
import os
import importlib.util
import inspect
from fastapi import FastAPI, Request
from fastapi.responses import JSONResponse
from pydantic import BaseModel
from typing import Optional, Any

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
                classes.append(obj)
    return classes

algorithm_classes = discover_algorithm_classes('algo')

class RecommendRequest(BaseModel):
    algo: Optional[str] = None
    user_id: Optional[Any] = None
    user_study_info: Optional[Any] = None
    blueprint: Optional[Any] = None

@app.post("/recommended")
async def recommended(req: RecommendRequest):
    algo_name = req.algo
    algo_class = None
    if algo_name:
        for cls in algorithm_classes:
            if cls.__name__ == algo_name:
                algo_class = cls
                break
    if not algo_class and algorithm_classes:
        algo_class = algorithm_classes[0]
    if algo_class:
        algo_instance = algo_class()
        if hasattr(algo_instance, "recommend"):
            result = algo_instance.recommend(req.user_id, req.user_study_info, req.blueprint)
        else:
            result = None
    else:
        result = None
    return JSONResponse(content={"recommended": result})

@app.get("/algorithms")
async def get_algorithms():
    algo_names = [cls.__name__ for cls in algorithm_classes]
    return JSONResponse(content={"algorithms": algo_names})

