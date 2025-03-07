import pandas as pd

degree_plans = pd.read_csv("./init_db/final_degree_plan.csv", dtype={
    "BLOC_LIMIT": pd.Int64Dtype()
})
common = ["PLAN_CODE","PLAN_YEAR","CODE","INTERCHANGEABILITY","BLOC_SUBJECT_CODE","BLOC_TYPE","BLOC_LIMIT","SEQ"]
dp_cs = degree_plans[common + ["BLOC_NAME_CZ","BLOC_NOTE_CZ","NOTE_CZ"]].rename(columns={"BLOC_NAME_CZ":"BLOC_NAME","BLOC_NOTE_CZ":"BLOC_NOTE","NOTE_CZ":"NOTE"})
dp_cs["LANG"] = "cs"
dp_en = degree_plans[common + ["BLOC_NAME_EN","BLOC_NOTE_EN","NOTE_EN"]].rename(columns={"BLOC_NAME_EN":"BLOC_NAME","BLOC_NOTE_EN":"BLOC_NOTE","NOTE_EN":"NOTE"})
dp_en["LANG"] = "en"
dp = pd.concat([dp_cs, dp_en])
# print(dp)
dp.to_csv("./init_db/transformed_degree_plan.csv", index=False)


#====================================================================================================
# JSON transformation
#====================================================================================================

# load = pd.read_csv("./init_db/export_degree_plan.csv", dtype={
#     "PLAN_CODE": str,
#     "PLAN_YEAR": pd.Int64Dtype(),
#     "CODE": str,
#     "NAME_CZ": str,
#     "NAME_EN": str,
#     "SUBJECT_STATUS": str,
#     "SEMESTER_PRIMARY": pd.Int64Dtype(),
#     "SEMESTER_COUNT": pd.Int64Dtype(),
#     "SUBJECT_TYPE": str,
#     "WORKLOAD_PRIMARY1": pd.Int64Dtype(),
#     "WORKLOAD_SECONDARY1": pd.Int64Dtype(),
#     "WORKLOAD_PRIMARY2": pd.Int64Dtype(),
#     "WORKLOAD_SECONDARY2": pd.Int64Dtype(),
#     "CREDITS": pd.Int64Dtype(),
#     "INTERCHANGEABILITY": str,
#     "BLOC_SUBJECT_CODE": str,
#     "BLOC_TYPE": str,
#     "BLOC_LIMIT": pd.Int64Dtype(),
#     "BLOC_NAME_CZ": str,
#     "BLOC_NAME_EN": str,
#     "SEQ": str,
#     "NOTE_CZ": str,
#     "NOTE_EN": str,
#     "BLOC_NOTE_CZ": str,
#     "BLOC_NOTE_EN": str
# })
# load = load[load["SUBJECT_STATUS"] == "V"]
# load = load[load["INTERCHANGEABILITY"].isna()]

# blocs_cs = load[["PLAN_CODE", "PLAN_YEAR", "BLOC_SUBJECT_CODE", "BLOC_TYPE", "BLOC_LIMIT", "BLOC_NAME_CZ", "BLOC_NOTE_CZ", "SEQ"]] \
#     .rename(columns={"BLOC_NAME_CZ": "BLOC_NAME", "BLOC_NOTE_CZ": "BLOC_NOTE"})
# blocs_cs["LANG"] = "cs"
# blocs_en = load[["PLAN_CODE", "PLAN_YEAR", "BLOC_SUBJECT_CODE", "BLOC_TYPE", "BLOC_LIMIT", "BLOC_NAME_EN", "BLOC_NOTE_EN", "SEQ"]] \
#     .rename(columns={"BLOC_NAME_EN": "BLOC_NAME", "BLOC_NOTE_EN": "BLOC_NOTE"})
# blocs_en["LANG"] = "en"
# blocs = pd.concat([blocs_cs, blocs_en]) \
#     .drop_duplicates(subset=["PLAN_CODE", "PLAN_YEAR", "BLOC_SUBJECT_CODE"]) \
#     .set_index(["PLAN_CODE", "PLAN_YEAR", "BLOC_SUBJECT_CODE", "LANG"])

# courses_cs = load.drop(columns=["BLOC_TYPE", "BLOC_LIMIT", "BLOC_NAME_CZ", "BLOC_NAME_EN", "BLOC_NOTE_CZ", "BLOC_NOTE_EN", "NAME_EN", "NOTE_EN"]) \
#     .rename(columns={"NAME_CZ": "NAME", "NOTE_CZ": "NOTE"})
# courses_cs["LANG"] = "cs"
# courses_en = load.drop(columns=["BLOC_TYPE", "BLOC_LIMIT", "BLOC_NAME_CZ", "BLOC_NAME_EN", "BLOC_NOTE_CZ", "BLOC_NOTE_EN", "NAME_CZ", "NOTE_CZ"]) \
#     .rename(columns={"NAME_EN": "NAME", "NOTE_EN": "NOTE"})
# courses_en["LANG"] = "en"
# courses = pd.concat([courses_cs, courses_en]) \
#     .set_index(["PLAN_CODE", "PLAN_YEAR", "BLOC_SUBJECT_CODE", "LANG", "SEQ"]) \
#     .apply(axis=1, func=lambda x: x.to_dict()) \
#     .reset_index() \
#     .rename(columns={0: "COURSE"}) \
#     .sort_values("SEQ") \
#     .groupby(["PLAN_CODE", "PLAN_YEAR", "BLOC_SUBJECT_CODE", "LANG"]) \
#     ["COURSE"] \
#     .apply(list) \
#     .rename("COURSES")

# degree_plans = pd.concat([blocs, courses], axis=1) \
#     .reset_index() \
#     .set_index(["PLAN_CODE", "PLAN_YEAR", "LANG", "SEQ"]) \
#     .apply(axis=1, func=lambda x: x.to_dict()) \
#     .reset_index() \
#     .rename(columns={0: "BLOC"}) \
#     .sort_values("SEQ") \
#     .groupby(["PLAN_CODE", "PLAN_YEAR", "LANG"])["BLOC"] \
#     .apply(lambda x: x.to_json(orient="records")) \
#     .rename("BLOCS")

# degree_plans.to_csv("./init_db/transformed_degree_plan.csv")