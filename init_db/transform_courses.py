import pandas as pd
import json

courses = pd.read_csv('init_db/courses.csv', dtype={"VUCIT1": str, "VUCIT2": str, "VUCIT3": str})
texts = pd.read_csv('init_db/course_texts.csv', usecols=["POVINN", "JAZYK", "TITLE", "MEMO", "TYP"])
teachers = pd.read_csv('init_db/teachers.csv', usecols=["KOD", "PRIJMENI", "JMENO", "TITULPRED", "TITULZA"], dtype={"KOD": str})
course_teachers = pd.read_csv('init_db/course_teachers.csv', dtype={"UCIT": str})

tea = pd.merge(course_teachers, teachers, left_on="UCIT", right_on="KOD")
tea = tea.drop(columns="UCIT").reset_index(drop=True).set_index("POVINN")
tea = tea.apply(lambda x: x.to_dict(), axis=1)
tea = tea.groupby("POVINN").agg(list).rename("teachers")

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
tex = tex.pivot_table(index=["POVINN", "JAZYK"], columns="TYP", values="MEMO", aggfunc="first")
tex = tex.reset_index().rename(columns={"JAZYK": "LANG"}).set_index(["POVINN", "LANG"])

gua_1 = courses[["POVINN", "VUCIT1"]].rename(columns={"VUCIT1": "VUCIT"})
gua_2 = courses[["POVINN", "VUCIT2"]].rename(columns={"VUCIT2": "VUCIT"})
gua_3 = courses[["POVINN", "VUCIT3"]].rename(columns={"VUCIT3": "VUCIT"})
gua = pd.concat([gua_1, gua_2, gua_3])
gua = pd.merge(gua, teachers, how="left", left_on="VUCIT", right_on="KOD")
gua = gua.drop(columns=["VUCIT"])
gua = gua.reset_index(drop=True).set_index("POVINN")
gua = gua.apply(lambda x: x.to_dict(), axis=1)
gua = gua.groupby("POVINN").agg(list).rename("GUARANTORS")

common = ["POVINN", "VPLATIOD", "VPLATIDO", "PFAKULTA",
       "PGARANT", "PVYUCOVAN", "VSEMZAC", "VSEMPOC", "PVYJAZYK", "VROZSAHPR1",
       "VROZSAHCV1", "VROZSAHPR2", "VROZSAHCV2", "VRVCEM", "VTYP", "VEBODY",
       "PPOCMIN", "PPOCMAX"]
courses_cs = courses[["PNAZEV"] + common].rename(columns={"PNAZEV": "NAME"})
courses_cs["LANG"] = "cs"
courses_en = courses[["PANAZEV"] + common].rename(columns={"PANAZEV": "NAME"})
courses_en["LANG"] = "en"
cou = pd.concat([courses_cs, courses_en])
cou = pd.merge(cou, gua, on="POVINN", how="left")
cou = pd.merge(cou, tex, on=["POVINN", "LANG"], how="left")

cou["GUARANTORS"] = cou["GUARANTORS"].apply(lambda x: json.dumps(x))
# cou["TEXTS"] = cou["TEXTS"].apply(lambda x: json.dumps(x))

print(cou)

cou.to_csv('init_db/courses_transformed.csv', index=False)