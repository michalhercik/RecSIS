from typing import Any, Optional

import pandas as pd
from fastapi import APIRouter, FastAPI
from fastapi.responses import JSONResponse
from pydantic import BaseModel
from typing_extensions import List
from user import User

from recommender import EvalRecommender, ProdRecommender

pd.set_option("future.no_silent_downcasting", True)

class RecommendRequest(BaseModel):
    user_id: str
    offset: int = 0
    limit: int = 10
    blueprint: Optional[Any] = None
    groups: bool = False
    categories: bool = False


class EvalRecommendRequest(BaseModel):
    algo: List[str]
    limit: int = 10
    user_id: str
    student: Optional[str] = None
    degree_plan: Optional[str] = None
    enrollment_year: Optional[int] = None
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


class FitRequest(BaseModel):
    algo: list[str]


config = {"env": "dev"}
prod_router = APIRouter()
eval_router = APIRouter()

app = FastAPI()
app.include_router(prod_router)
if config["env"] == "dev":
    app.mount("/eval", eval_router)

recommender = ProdRecommender()
recommender.fit(cache=True)
# eval_recommender = EvalRecommender()



@prod_router.post("/foryou")
async def recommend(req: RecommendRequest):
    user = User(req.user_id, None, None, req.blueprint)
    result = recommender.recommend(user, req.offset, req.limit, req.groups, req.categories)

    return JSONResponse(content=result)


@prod_router.post("/fit")
async def fit():
    recommender.fit()


# @eval_router.post("/recommended")
# async def eval_recommend(req: EvalRecommendRequest):
#     user = User(req.user_id, req.degree_plan, req.enrollment_year, req.blueprint)
#     if req.student is not None and len(req.student) > 0:
#         user.id = req.student
#         user.fetch = True

#     limit = req.limit
#     algo = req.algo
#     result = eval_recommender.recommend(user, algo, limit)
#     return JSONResponse(content=result)


# @eval_router.post("/fit")
# async def eval_fit(req: FitRequest):
#     algos = req.algo
#     eval_recommender.fit(algos)


# @eval_router.get("/algorithms")
# async def algorithms():
#     result = eval_recommender.algorithms()
#     return JSONResponse(content=result)
