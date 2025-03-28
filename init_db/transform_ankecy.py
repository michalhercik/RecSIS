import pandas as pd
import numpy as np
import json

ankecy = pd.read_csv('init_db/ANKECY.csv', dtype={"SROC": pd.Int32Dtype(), "UCIT": str})
teachers = pd.read_csv('init_db/UCIT.csv', usecols=["KOD", "PRIJMENI", "JMENO", "TITULPRED", "TITULZA"], dtype={"KOD": str})
druh = pd.read_csv('init_db/DRUH.csv')

druh_cs = druh[["KOD", "NAZEV", "ZKRATKA"]].copy()
druh_cs["LANG"] = "cs"
druh_en = druh[["KOD", "ANAZEV", "AZKRATKA"]].rename(columns={"ANAZEV": "NAZEV", "AZKRATKA": "ZKRATKA"})
druh_en["LANG"] = "en"

dr = pd.concat([druh_cs, druh_en], axis=0)

dict_tea = teachers.copy()
dict_tea["id"] = dict_tea["KOD"]
dict_tea = dict_tea.set_index("id")
dict_tea = dict_tea.apply(lambda x: x.replace({pd.NA: None}).to_dict(), axis=1).rename("TEACHER")

comments = pd.merge(ankecy, dict_tea, left_on="UCIT", right_on="id", how="left")
comments = pd.merge(comments, dr, left_on="SDRUH", right_on="KOD", how="left")
comments = comments.drop(columns=["UCIT", "SDRUH"])

comments["TEACHER"] = comments["TEACHER"].apply(lambda x: json.dumps(x))

comments.to_csv("init_db/ankecy_transformed.csv", index=False)