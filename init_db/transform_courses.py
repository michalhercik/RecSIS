import pandas as pd
import numpy as np
import json

courses = pd.read_csv('init_db/POVINN.csv', dtype={
    "VUCIT1": str,
    "VUCIT2": str,
    "VUCIT3": str,
    "VSEMZAC": pd.Int32Dtype(),
    "VSEMPOC": pd.Int32Dtype(),
    "VROZSAHPR1": pd.Int32Dtype(),
    "VROZSAHCV1": pd.Int32Dtype(),
    "VROZSAHPR2": pd.Int32Dtype(),
    "VROZSAHCV2": pd.Int32Dtype(),
    "VEBODY": pd.Int32Dtype(),
    "PPOCMIN": pd.Int32Dtype(),
    "PPOCMAX": str,
    "PFAKULTA": str,
    })
courses = courses[courses["VPLATIDO"] == 9999]
texts = pd.read_csv('init_db/PAMELA.csv', usecols=["POVINN", "JAZYK", "MEMO", "TYP"])
teachers = pd.read_csv('init_db/UCIT.csv', usecols=["KOD", "PRIJMENI", "JMENO", "TITULPRED", "TITULZA"], dtype={"KOD": str})
course_teachers = pd.read_csv('init_db/UCIT_ROZVRH.csv', dtype={"UCIT": str})
faculties = pd.read_csv('init_db/faculties.csv', dtype={"FACULTY_ID": str}, usecols=["FACULTY_ID", "FACULTY_NAME_CS", "FACULTY_NAME_EN"])
texts_titles = pd.read_csv('init_db/memo_title.csv')
course_langs = pd.read_csv('init_db/POVINN2JAZYK.csv')
languages = pd.read_csv('init_db/JAZYK.csv')
ankecy = pd.read_csv('init_db/ankecy_transformed.csv', dtype={"SROC": pd.Int32Dtype()}, usecols=["POVINN", "LANG", "NAZEV", "SROC", "SOBOR", "SSKR", "PRDMTYP", "TEACHER", "MEMO"])
preq = pd.read_csv('init_db/PREQ.csv')
ptrida  = pd.read_csv('init_db/PTRIDA.csv')
trida = pd.read_csv('init_db/TRIDA.csv', usecols=["KOD", "NAZEV"])
pklas = pd.read_csv('init_db/PKLAS.csv')
klas = pd.read_csv('init_db/KLAS.csv')

klas = pd.merge(klas, pklas, left_on="KOD", right_on="PKLAS")
klas_cs = klas[["POVINN", "KOD", "NAZEV"]].copy()
klas_cs["LANG"] = "cs"
klas_en = klas[["POVINN", "KOD", "ANAZEV"]].rename(columns={"ANAZEV": "NAZEV"})
klas_en["LANG"] = "en"
klas = pd.concat([klas_cs, klas_en])
klas = klas.set_index(["POVINN", "LANG"])
klas = klas.apply(lambda x: x.replace({pd.NA: None}).to_dict(), axis=1)
klas = klas.groupby(["POVINN", "LANG"]).agg(lambda x: json.dumps(list(x))).rename("CLASSIFICATIONS")

trida = pd.merge(ptrida, trida, left_on="PTRIDA", right_on="KOD")
trida = trida.drop(columns=["PTRIDA"])
trida = trida.set_index("POVINN")
trida = trida.apply(lambda x: x.replace({pd.NA: None}).to_dict(), axis=1)
trida = trida.groupby("POVINN").agg(lambda x: json.dumps(list(x))).rename("CLASSES")

preq = preq.groupby(["POVINN", "REQTYP"]).agg(lambda x: json.dumps(list(x)))
preq = preq.pivot_table(columns=["REQTYP"], values=["REQPOVINN"], index=["POVINN"], aggfunc="first")
preq = preq.droplevel(0, axis=1)
preq.columns.name = None

ankecy["TEACHER"] = ankecy["TEACHER"].apply(lambda x: json.loads(x) if pd.notna(x) else None)
comments = ankecy.set_index(["POVINN", "LANG"])
comments = comments.apply(lambda x: x.replace({pd.NA: None}).to_dict(), axis=1)
comments = comments.groupby(["POVINN", "LANG"]).agg(list).rename("COMMENTS")

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
tex = pd.concat([
    pd.merge(tex[tex["JAZYK"] == "cs"], texts_titles[["KOD", "NAZEV"]], left_on="TYP", right_on="KOD").rename(columns={"NAZEV": "TITLE"}),
    pd.merge(tex[tex["JAZYK"] == "en"], texts_titles[["KOD", "ANAZEV"]], left_on="TYP", right_on="KOD").rename(columns={"ANAZEV": "TITLE"})
], axis=0).drop(columns=["KOD"])
tex = tex.reset_index(drop=True).set_index(["POVINN", "JAZYK"])
# translate = {
#     "A": "ANOTACE",
#     "C": "CIL",
#     "E": "ZAKONCENI",
#     "S": "SYLABUS",
#     "P": "ZKOUSKA",
#     "L": "LITERATURA",
#     "V": "VSTUP",
# }
tex = tex[tex["TYP"].isin(["A", "C", "E", "S", "P", "L", "V"])]
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

lang = course_langs[course_langs["PLATIDO"] == 9999].set_index("POVINN")[["JAZYK"]]
a = courses.set_index("POVINN")[["PVYJAZYK"]].rename(columns={"PVYJAZYK": "JAZYK"})
lang = pd.concat([a, lang]).dropna()
lang = pd.merge(lang.reset_index(), languages, left_on="JAZYK", right_on="KOD")
lang_cs = lang.copy()
lang_cs["LANG"] = "cs"
lang_cs = lang_cs[["POVINN", "NAZEV", "LANG"]]
lang_en = lang.copy()
lang_en["LANG"] = "en"
lang_en = lang_en[["POVINN", "ANAZEV", "LANG"]]
lang_en = lang_en.rename(columns={"ANAZEV": "NAZEV"})
lang = pd.concat([lang_cs, lang_en])
lang["NAZEV"] = lang["NAZEV"].apply(lambda x: x.capitalize())
lang = lang.groupby(["POVINN", "LANG"]).agg(list)
lang = lang.apply(lambda x: ", ".join(sorted(x["NAZEV"])), axis=1).rename("VYJAZYK")

common = ["POVINN", "VPLATIOD", "VPLATIDO", "PFAKULTA",
       "PGARANT", "PVYUCOVAN", "VSEMZAC", "VSEMPOC", "VROZSAHPR1",
       "VROZSAHCV1", "VROZSAHPR2", "VROZSAHCV2", "VRVCEM", "VTYP", "VEBODY",
       "PPOCMIN", "PPOCMAX"]
courses_cs = courses[["PNAZEV"] + common].rename(columns={"PNAZEV": "NAME"})
courses_cs["PPOCMAX"] = courses_cs["PPOCMAX"].replace({np.nan: "Neomezená"})
# TODO: not complete
courses_cs["VTYP"] = courses_cs["VTYP"].replace({"Z": "Z", "F": "KZ", "K": "Zk", "*": "Z+Zk"})
courses_cs["LANG"] = "cs"
courses_cs["PVYUCOVAN"] = courses_cs["PVYUCOVAN"].replace({"V": "Vyučován", "N": "Nevyučován", "Z": "Zrušen"})
courses_en = courses[["PANAZEV"] + common].rename(columns={"PANAZEV": "NAME"})
courses_en["PPOCMAX"] = courses_en["PPOCMAX"].replace({np.nan: "Unlimited"})
courses_en["VTYP"] = courses_en["VTYP"].replace({"Z": "C", "F": "MC", "K": "Ex", "*": "C+Ex"})
courses_en["LANG"] = "en"
courses_en["PVYUCOVAN"] = courses_en["PVYUCOVAN"].replace({"V": "Taught", "N": "Not taught", "Z": "Cancelled"})
cou = pd.concat([courses_cs, courses_en])
cou = pd.merge(cou, lang, on=["POVINN", "LANG"], how="left")
cou = pd.merge(cou, gua, on="POVINN", how="left")
cou = pd.merge(cou, tex, on=["POVINN", "LANG"], how="left")
cou = pd.merge(cou, tea, on="POVINN", how="left")
cou = pd.merge(cou, fac, left_on=["PFAKULTA", "LANG"], right_on=["FACULTY_ID", "LANG"], how="left")
cou = pd.merge(cou, comments, on=["POVINN", "LANG"], how="left")
cou = pd.merge(cou, preq, how="left", on=["POVINN"])
cou = pd.merge(cou, trida, how="left", on=["POVINN"])
cou = pd.merge(cou, klas, how="left", on=["POVINN", "LANG"])
cou = cou.drop(columns="PFAKULTA")


cou["GUARANTORS"] = cou["GUARANTORS"].replace({pd.NA: None})
cou["GUARANTORS"] = cou["GUARANTORS"].apply(lambda x: json.dumps(x))
cou["TEACHERS"] = cou["TEACHERS"].replace({pd.NA: None})
cou["TEACHERS"] = cou["TEACHERS"].apply(lambda x: json.dumps(x))
cou["COMMENTS"] = cou["COMMENTS"].replace({pd.NA: None})
cou["COMMENTS"] = cou["COMMENTS"].apply(lambda x: json.dumps(x))

def condition(row, first, second):
    if pd.notna(row["VSEMZAC"]) and row["VSEMZAC"] == 2:
        return row[second]
    return row[first]

cou["LECTURE_RANGE_WINTER"] = cou[["VROZSAHPR1", "VROZSAHPR2", "VSEMZAC"]].apply(lambda x: condition(x, "VROZSAHPR1", "VROZSAHPR2"), axis=1)
cou["SEMINAR_RANGE_WINTER"] = cou[["VROZSAHCV1", "VROZSAHCV2", "VSEMZAC"]].apply(lambda x: condition(x, "VROZSAHCV1", "VROZSAHCV2"), axis=1)
cou["LECTURE_RANGE_SUMMER"] = cou[["VROZSAHPR1", "VROZSAHPR2", "VSEMZAC"]].apply(lambda x: condition(x, "VROZSAHPR2", "VROZSAHPR1"), axis=1)
cou["SEMINAR_RANGE_SUMMER"] = cou[["VROZSAHCV1", "VROZSAHCV2", "VSEMZAC"]].apply(lambda x: condition(x, "VROZSAHCV2", "VROZSAHCV1"), axis=1)

cou = cou.drop(columns=["VROZSAHPR1", "VROZSAHPR2", "VROZSAHCV1", "VROZSAHCV2"])

cou.to_csv('init_db/courses_transformed.csv', index=False)