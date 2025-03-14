import pandas as pd

courses = pd.read_csv('./init_db/POVINN.csv', usecols=[
    "POVINN", "PNAZEV", "PANAZEV", "VPLATIDO",
    "VSEMZAC", "VSEMPOC",
    "VROZSAHPR1", "VROZSAHCV1", "VROZSAHPR2", "VROZSAHCV2",
    "VTYP", "VEBODY",
    "VUCIT1", "VUCIT2", "VUCIT3",])
teachers = pd.read_csv('./init_db/UCIT.csv', usecols=["KOD", "JMENO", "PRIJMENI"])
texts = pd.read_csv('./init_db/PAMELA.csv', usecols=["POVINN", "TYP", "JAZYK", "MEMO"])

df = courses[courses["VPLATIDO"] == 9999].drop(columns=["VPLATIDO"])
df = pd.merge(df, teachers, how="left", left_on="VUCIT1", right_on="KOD").rename(columns={"KOD": "VUCIT1_KOD", "JMENO": "VUCIT1_JMENO", "PRIJMENI": "VUCIT1_PRIJMENI"})
df = pd.merge(df, teachers, how="left", left_on="VUCIT2", right_on="KOD").rename(columns={"KOD": "VUCIT2_KOD", "JMENO": "VUCIT2_JMENO", "PRIJMENI": "VUCIT2_PRIJMENI"})
df = pd.merge(df, teachers, how="left", left_on="VUCIT3", right_on="KOD").rename(columns={"KOD": "VUCIT3_KOD", "JMENO": "VUCIT3_JMENO", "PRIJMENI": "VUCIT3_PRIJMENI"})
df = df.drop(columns=["VUCIT1", "VUCIT2", "VUCIT3"])
df = df[df["VROZSAHPR1"].notna()]

texts = pd.pivot_table(texts, index=["POVINN", "JAZYK"], columns=["TYP"], values=["MEMO"], aggfunc=lambda x: x)
texts = texts.reset_index().droplevel(0, axis=1)
texts.columns.values[0] = "POVINN"
texts.columns.values[1] = "JAZYK"
texts = texts.rename(columns={
    "A": "ANOTACE",
    "C": "CIL",
    "S": "SYLABUS",
    "P": "POZADAVKY"
    })

df = pd.merge(df, texts[texts["JAZYK"] == "CZE"], how="left", on="POVINN").drop(columns=["JAZYK"])
df = pd.merge(df, texts[texts["JAZYK"] == "ENG"], how="left", on="POVINN", suffixes=("_CZE", "_ENG")).drop(columns=["JAZYK"])

df = df.rename(columns={
    "POVINN": "code",
    "PNAZEV": "nameCs",
    "PANAZEV": "nameEn",
    "VSEMZAC": "start",
    "VSEMPOC": "semesterCount",
    "VROZSAHPR1": "lectureRange1",
    "VROZSAHCV1": "seminarRange1",
    "VROZSAHPR2": "lectureRange2",
    "VROZSAHCV2": "seminarRange2",
    "VTYP": "examType",
    "VEBODY": "credits",
    "VUCIT1_KOD": "teacher1Id",
    "VUCIT1_JMENO": "teacher1Name",
    "VUCIT1_PRIJMENI": "teacher1Lastname",
    "VUCIT2_KOD": "teacher2Id",
    "VUCIT2_JMENO": "teacher2Name",
    "VUCIT2_PRIJMENI": "teacher2Lastname",
    "VUCIT3_KOD": "teacher3Id",
    "VUCIT3_JMENO": "teacher3Name",
    "VUCIT3_PRIJMENI": "teacher3Lastname",
    "ANOTACE_CZE": "annotationCs",
    "CIL_CZE": "aimsCs",
    "SYLABUS_CZE": "syllabusCs",
    "POZADAVKY_CZE": "requirementsCs",
    "ANOTACE_ENG": "annotationEn",
    "CIL_ENG": "aimsEn",
    "SYLABUS_ENG": "syllabusEn",
    "POZADAVKY_ENG": "requirementsEn"
})

df = df.reset_index(drop=False).rename(columns={"index": "id"})
df.to_json('./init_search/courses.json', orient='records', lines=True)