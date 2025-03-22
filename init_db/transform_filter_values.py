import pandas as pd

courses = pd.read_csv('./init_db/courses_transformed.csv', usecols=[
    "POVINN", "LANG", "PGARANT", "PVYUCOVAN", "PPOCMIN", "PPOCMAX", "FACULTY_NAME",
    "VSEMZAC", "VSEMPOC", "VTYP", "VRVCEM", "VEBODY",
    "LECTURE_RANGE_WINTER", "SEMINAR_RANGE_WINTER", "LECTURE_RANGE_SUMMER", "SEMINAR_RANGE_SUMMER",
    # "VROZSAHPR1", "VROZSAHCV1", "VROZSAHPR2", "VROZSAHCV2",
    ], dtype={
        "VYJAZYK": str,
        "VEBODY": str,
        "VSEMZAC": str,
        "VSEMPOC": str,
        "PPOCMIN": str,
        "PPOCMAX": str,
        "LECTURE_RANGE_WINTER": str,
        "SEMINAR_RANGE_WINTER": str,
        "LECTURE_RANGE_SUMMER": str,
        "SEMINAR_RANGE_SUMMER": str,
    }
)
langs = pd.read_csv('./init_db/JAZYK.csv', usecols=["NAZEV", "ANAZEV"]).rename(columns={"NAZEV": "cs", "ANAZEV": "en"}).map(lambda x: x.capitalize())

courses = courses[~courses["PVYUCOVAN"].isin(["Zrušen", "Cancelled"])]

cs = courses[courses["LANG"] == "cs"]
en = courses[courses["LANG"] == "en"]
df_keys = pd.merge(cs, en, on="POVINN", suffixes=("CS", "EN"))
df_keys = [df_keys[[col+"CS", col+"EN"]].rename(columns={col+"CS":"cs", col+"EN":"en"}) for col in courses.drop(columns=["POVINN", "LANG"])]
df_keys.append(langs)
df_keys = pd.concat(df_keys).dropna()
df_keys = df_keys.drop_duplicates()
# df_keys = df_keys.reset_index(drop=True)
# df = pd.concat([df_keys[["cs"]].rename(columns={"cs": "VALUE"}), df_keys[["en"]].rename(columns={"en": "VALUE"})])

df = df_keys.reset_index(drop=True).reset_index(drop=False).melt(id_vars=["index"])
df.to_csv('./init_db/filter_values.csv', index=False)

# df = df_keys.T.apply(list, axis=1)
# df = df.apply(lambda x: "{" + ",".join([f"\"{v}\"" for v in x]) + "}")
# df.to_csv('./init_db/filter_values_array.csv', index=True, header=False, sep=";", quotechar="'")


# series_list = [courses[["LANG", col]].rename(columns={col:"VALUE"}) for col in courses.drop(columns=["POVINN", "LANG"])]
# df = pd.concat(series_list).dropna().drop_duplicates().reset_index(drop=True)
# df.to_csv('./init_db/filter_values.csv', index=True)