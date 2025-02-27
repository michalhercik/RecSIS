import pandas as pd
import numpy as np
import json

courses = pd.read_csv('init_db/courses.csv', dtype={
    "VUCIT1": str,
    "VUCIT2": str,
    "VUCIT3": str,
    "VSEMZAC": str,
    "VSEMPOC": pd.Int32Dtype,
    "VROZSAHPR1": pd.Int32Dtype,
    "VROZSAHCV1": pd.Int32Dtype,
    "VROZSAHPR2": pd.Int32Dtype,
    "VROZSAHCV2": pd.Int32Dtype,
    "VEBODY": pd.Int32Dtype,
    "PPOCMIN": pd.Int32Dtype,
    "PPOCMAX": str,
    "PFAKULTA": str,
    })
texts = pd.read_csv('init_db/course_texts.csv', usecols=["POVINN", "JAZYK", "TITLE", "MEMO", "TYP"])
teachers = pd.read_csv('init_db/teachers.csv', usecols=["KOD", "PRIJMENI", "JMENO", "TITULPRED", "TITULZA"], dtype={"KOD": str})
course_teachers = pd.read_csv('init_db/course_teachers.csv', dtype={"UCIT": str})
faculties = pd.read_csv('init_db/faculties.csv', dtype={"FACULTY_ID": str}, usecols=["FACULTY_ID", "FACULTY_NAME_CS", "FACULTY_NAME_EN"])

fac = faculties.reset_index().set_index("FACULTY_ID")
fac_cs = pd.DataFrame(fac["FACULTY_NAME_CS"].rename("FACULTY_NAME"))
fac_cs["LANG"] = "cs"
fac_en = pd.DataFrame(fac["FACULTY_NAME_EN"].rename("FACULTY_NAME"))
fac_en["LANG"] = "en"
fac = pd.concat([fac_cs, fac_en])

tea = pd.merge(course_teachers, teachers, left_on="UCIT", right_on="KOD")
tea = tea.drop(columns="UCIT").reset_index(drop=True).set_index("POVINN")
tea = tea.apply(lambda x: x.replace({pd.NA: None}).to_dict(), axis=1)
tea = tea.groupby("POVINN").agg(list).rename("TEACHERS")

tex = texts
tex["JAZYK"] = tex["JAZYK"].map({"CZE": "cs", "ENG": "en"})
tex = tex.reset_index(drop=True).set_index(["POVINN", "JAZYK"])
translate = {
    "A": "ANOTACE",
    "C": "CIL",
    "S": "SYLABUS",
    "P": "POZADAVKY",
}
tex["TYP"] = tex["TYP"].map(translate)
tex["DESCRIPTION"] = tex[["TITLE", "MEMO"]].apply(lambda x: x.to_json(), axis=1)
tex = tex.pivot_table(index=["POVINN", "JAZYK"], columns="TYP", values="DESCRIPTION", aggfunc="first")
tex = tex.reset_index().rename(columns={"JAZYK": "LANG"}).set_index(["POVINN", "LANG"])

gua_1 = courses[["POVINN", "VUCIT1"]].rename(columns={"VUCIT1": "VUCIT"})
gua_2 = courses[["POVINN", "VUCIT2"]].rename(columns={"VUCIT2": "VUCIT"})
gua_3 = courses[["POVINN", "VUCIT3"]].rename(columns={"VUCIT3": "VUCIT"})
gua = pd.concat([gua_1, gua_2, gua_3])
gua = pd.merge(gua, teachers, how="left", left_on="VUCIT", right_on="KOD")
gua = gua.drop(columns=["VUCIT"])
gua = gua.reset_index(drop=True).set_index("POVINN")
gua = gua.dropna()
gua = gua.apply(lambda x: x.replace({pd.NA: None}).to_dict(), axis=1)
gua = gua.groupby("POVINN").agg(list).rename("GUARANTORS")

common = ["POVINN", "VPLATIOD", "VPLATIDO", "PFAKULTA",
       "PGARANT", "PVYUCOVAN", "VSEMZAC", "VSEMPOC", "PVYJAZYK", "VROZSAHPR1",
       "VROZSAHCV1", "VROZSAHPR2", "VROZSAHCV2", "VRVCEM", "VTYP", "VEBODY",
       "PPOCMIN", "PPOCMAX"]
courses_cs = courses[["PNAZEV"] + common].rename(columns={"PNAZEV": "NAME"})
courses_cs["PPOCMAX"] = courses_cs["PPOCMAX"].replace({np.nan: "Neomezená"})
courses_cs["VTYP"] = courses_cs["VTYP"].replace({"Z": "Z", "F": "KZ", "K": "Zk", "*": "Z+Zk"})
courses_cs["LANG"] = "cs"
courses_cs["VSEMZAC"] = courses_cs["VSEMZAC"].replace({"1": "Zimní", "2": "Letní", "3": "Oba"})
courses_cs["PVYJAZYK"] = courses_cs["PVYJAZYK"].replace({"CZE": "Čeština", "ENG": "Angličtina", np.nan: "Neznámý"})
courses_cs["PVYUCOVAN"] = courses_cs["PVYUCOVAN"].replace({"V": "Vyučován", "N": "Nevyučován", "Z": "Zrušen"})
courses_en = courses[["PANAZEV"] + common].rename(columns={"PANAZEV": "NAME"})
courses_en["PPOCMAX"] = courses_en["PPOCMAX"].replace({np.nan: "Unlimited"})
courses_en["VTYP"] = courses_en["VTYP"].replace({"Z": "C", "F": "MC", "K": "Ex", "*": "C+Ex"})
courses_en["LANG"] = "en"
courses_en["VSEMZAC"] = courses_en["VSEMZAC"].replace({"1": "Winter", "2": "Summer", "3": "Both"})
courses_en["PVYJAZYK"] = courses_en["PVYJAZYK"].replace({"CZE": "Czech", "ENG": "English", np.nan: "Unknown"})
courses_en["PVYUCOVAN"] = courses_en["PVYUCOVAN"].replace({"V": "Taught", "N": "Not taught", "Z": "Cancelled"})
cou = pd.concat([courses_cs, courses_en])
cou = pd.merge(cou, gua, on="POVINN", how="left")
cou = pd.merge(cou, tex, on=["POVINN", "LANG"], how="left")
cou = pd.merge(cou, tea, on="POVINN", how="left")
cou = pd.merge(cou, fac, left_on=["PFAKULTA", "LANG"], right_on=["FACULTY_ID", "LANG"], how="left")
cou = cou.drop(columns="PFAKULTA")


cou["GUARANTORS"] = cou["GUARANTORS"].replace({pd.NA: None})
cou["GUARANTORS"] = cou["GUARANTORS"].apply(lambda x: json.dumps(x))
cou["TEACHERS"] = cou["TEACHERS"].replace({pd.NA: None})
cou["TEACHERS"] = cou["TEACHERS"].apply(lambda x: json.dumps(x))

cou.to_csv('init_db/courses_transformed.csv', index=False)