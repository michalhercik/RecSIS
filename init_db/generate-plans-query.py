import pandas as pd

plans = pd.read_csv("./init_search/stud-plans.csv")
now = pd.Timestamp.now().year

query = """
(
	SELECT  PLAN_CODE, PLAN_YEAR, CODE, INTERCHANGEABILITY, BLOC_SUBJECT_CODE, BLOC_TYPE, BLOC_LIMIT, BLOC_NAME_CZ, BLOC_NAME_EN, SEQ, NOTE_CZ, NOTE_EN, BLOC_NOTE_CZ, BLOC_NOTE_EN
	FROM table(study_plan.stud_plan('{code}', {year}))
)
"""
queries = []
for plan_code in plans["SPLAN"]:
    for year in range(now - 9, now + 1):
        queries.append(query.format(
            code=plan_code,
            year=year
        ))


with open("./dp-query.sql", "w") as f:
  f.write("UNION".join(queries) + ";")
