import pandas as pd

courses = pd.read_csv('./init_db/courses_transformed.csv', usecols=[
    "POVINN", "LANG", "PGARANT", "PVYUCOVAN", "PPOCMIN", "PPOCMAX", "FACULTY_NAME",
    "VSEMZAC", "VSEMPOC", "VTYP", "VRVCEM", "VEBODY", "VYJAZYK",
    "LECTURE_RANGE_WINTER", "SEMINAR_RANGE_WINTER", "LECTURE_RANGE_SUMMER", "SEMINAR_RANGE_SUMMER",
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
# langs = pd.read_csv('./init_db/JAZYK.csv', usecols=["NAZEV", "ANAZEV"]).rename(columns={"NAZEV": "title_cs", "ANAZEV": "title_en"})
courses = courses[~courses["PVYUCOVAN"].isin(["Zrušen", "Cancelled"])]
courses["PPOCMAX"] = courses["PPOCMAX"].replace(["Neomezená", "Unlimited"], -1)
courses["PPOCMAX"] = courses["PPOCMAX"].astype(int)

cs = courses[courses["LANG"] == "cs"]
en = courses[courses["LANG"] == "en"]
df = pd.merge(cs, en, on="POVINN", suffixes=("CS", "EN"))

facet2column = {
    # "lecture_range_winter": "LECTURE_RANGE_WINTER",
    # "lecture_range_summer": "LECTURE_RANGE_SUMMER",
    # "seminar_range_winter": "SEMINAR_RANGE_WINTER",
    # "seminar_range_summer": "SEMINAR_RANGE_SUMMER",
    "lecture_range": "LECTURE_RANGE",
    "seminar_range": "SEMINAR_RANGE",
    "exam_type": "VTYP",
    "credits": "VEBODY",
    "faculty": "FACULTY_NAME",
    "department": "PGARANT",
    "taught": "PVYUCOVAN",
    "capacity": "PPOCMAX",
    "min_occupancy": "PPOCMIN",
    "semester_count": "VSEMPOC",
    "start_semester": "VSEMZAC",
    "range_unit": "VRVCEM",
}

asint = lambda x: x.astype(pd.Int64Dtype())
exam_order = {
    "Z+Zk": 0,
    "Zk": 1,
    "KZ": 2,
    "Z": 3,
}
facet2sort = {
    # "lecture_range_winter": asint,
    # "lecture_range_summer": asint,
    # "seminar_range_winter": asint,
    # "seminar_range_summer": asint,
    "lecture_range": asint,
    "seminar_range": asint,
    "credits": asint,
    "min_occupancy": asint,
    "capacity": asint,
    "semester_count": asint,
    "exam_type": lambda x: x.apply(lambda e: exam_order.get(e, len(exam_order)))
}

format_str = lambda x: x.astype(str).str.lower().str.replace(" ", "_")
format_facet_id = {
    "faculty": format_str,
    "language": format_str,
    "taught": format_str,
    "exam_type": format_str,
    "department": format_str,
    "range_unit": format_str,
}

categories = pd.DataFrame(
    data=[
        ["faculty", "Fakulta", "Faculty", "", "", 3],
        ["taught", "Stav předmětu", "Course status", "", "", 3],
        ["start_semester", "Semestr", "Semester", "", "", 3],
        ["credits", "Kredity", "Credits", "", "", 6],
        ["semester_count", "Počet semestrů", "Number of semesters", "", "", 3],
        ["lecture_range", "Rozsah přednášky", "Lecture range", "", "", 3],
        ["seminar_range", "Rozsah cvičení", "Seminar range", "", "", 3],
        # ["lecture_range_winter", "Rozsah přednášky zima", "Winter lecture range", "", ""],
        # ["lecture_range_summer", "Rozsah přednášky léto", "Summer Lecture range", "", ""],
        # ["seminar_range_winter", "Rozsah semináře zima", "Winter seminar range", "", ""],
        # ["seminar_range_summer", "Rozsah semináře léto", "Summer seminar range", "", ""],
        ["language", "Jazyk výuky", "Language", "", "", 2],
        ["exam_type", "Typ examinace", "Exam type", "", "", 4],
        ["range_unit", "Jednotka rozsahu", "Range unit", "", "", 4],
        ["department", "Katedra", "Department", "", "", 6],
        # ["capacity", "Kapacita", "Capacity", "", ""],
        # ["min_occupancy", "Minimální Obsazenost", "Minimum Occupancy", "", ""],
    ],
    columns=["facet_id", "title_cs", "title_en", "desc_cs", "desc_en", "displayed_value_limit"]
)
categories["position"] = categories.index

language_order = {
    "čeština": 0,
    "slovenština": 1,
    "angličtina": 2,
    "němčina": 3,
    "španělština": 4,
    "francouština": 5,
}
# langs["description_cs"] = ""
# langs["description_en"] = ""
# langs["category"] = categories[categories["facet_id"] == "language"].index[0]
# langs["facet_id"] = format_str(langs["title_en"]) #.astype(str).str.lower()
# langs = langs.sort_values(by=["title_cs", "title_en"], key=lambda x: x.apply(lambda e: language_order.get(e, len(language_order))))
# langs["position"] = langs.index

def reverse(x):
    x.reverse()
    return x
langs = df[["VYJAZYKCS", "VYJAZYKEN"]].rename({"VYJAZYKCS": "title_cs", "VYJAZYKEN": "title_en"}, axis=1)
langs = langs.drop_duplicates().dropna()
langs = langs.apply(lambda x: x.str.split(","))
langs["title_cs"] = langs["title_cs"].apply(reverse)
langs = langs.explode(["title_cs", "title_en"])
langs = langs.map(lambda x: x.strip())
langs = langs.drop_duplicates(ignore_index=True)
langs["category"] = categories[categories["facet_id"] == "language"].index[0]
langs["description_cs"] = ""
langs["description_en"] = ""
langs["facet_id"] = format_str(langs["title_en"])
langs = langs.sort_values(by=["title_cs", "title_en"], key=lambda x: x.apply(lambda e: language_order.get(e, len(language_order))))
langs = langs.reset_index(drop=True)
langs["position"] = langs.index

lrange = pd.concat([df["LECTURE_RANGE_WINTERCS"], df["LECTURE_RANGE_SUMMERCS"]])
lrange = lrange.rename("facet_id")
lrange = lrange.drop_duplicates().dropna()
lrange = lrange.sort_values(key=facet2sort["lecture_range"])
lrange = lrange.reset_index(drop=True)
lrange = lrange.to_frame()
lrange["category"] = categories[categories["facet_id"] == "lecture_range"].index[0]
lrange["title_cs"] = lrange["facet_id"]
lrange["title_en"] = lrange["facet_id"]
lrange["description_cs"] = ""
lrange["description_en"] = ""
lrange["position"] = lrange.index

srange = pd.concat([df["SEMINAR_RANGE_WINTERCS"], df["SEMINAR_RANGE_SUMMERCS"]])
srange = srange.rename("facet_id")
srange = srange.drop_duplicates().dropna()
srange = srange.sort_values(key=facet2sort["seminar_range"])
srange = srange.reset_index(drop=True)
srange = srange.to_frame()
srange["category"] = categories[categories["facet_id"] == "seminar_range"].index[0]
srange["title_cs"] = srange["facet_id"]
srange["title_en"] = srange["facet_id"]
srange["description_cs"] = ""
srange["description_en"] = ""
srange["position"] = srange.index

category_values = [langs, lrange, srange]
for row in categories[~categories["facet_id"].isin(["language", "lecture_range", "seminar_range"])].iterrows():
    col = facet2column[row[1]["facet_id"]]
    cv = df[[col+"CS", col+"EN"]].drop_duplicates().rename(columns={col+"CS": "title_cs", col+"EN": "title_en"})
    formater = format_facet_id.get(row[1]["facet_id"], lambda x: x)
    cv["facet_id"] = formater(cv["title_en"]) #.astype(str).str.lower()
    cv["category"] = row[0]
    cv["description_cs"] = ""
    cv["description_en"] = ""
    sortkey = facet2sort.get(row[1]["facet_id"], lambda x: x)
    cv = cv.sort_values(by=["title_cs", "title_en"], key=sortkey)
    cv = cv.reset_index(drop=True)
    cv["position"] = cv.index
    category_values.append(cv)

category_values = pd.concat(category_values).reset_index(drop=True).dropna()

categories.to_csv('./init_db/filter_categories.csv', index=True)
category_values.to_csv('./init_db/filter_values.csv', index=True)
# langs = pd.read_csv('./init_db/JAZYK.csv', usecols=["NAZEV", "ANAZEV"]).rename(columns={"NAZEV": "cs", "ANAZEV": "en"}).map(lambda x: x.capitalize())

# courses = courses[~courses["PVYUCOVAN"].isin(["Zrušen", "Cancelled"])]

# cs = courses[courses["LANG"] == "cs"]
# en = courses[courses["LANG"] == "en"]
# df_keys = pd.merge(cs, en, on="POVINN", suffixes=("CS", "EN"))
# df_keys = [df_keys[[col+"CS", col+"EN"]].rename(columns={col+"CS":"cs", col+"EN":"en"}) for col in courses.drop(columns=["POVINN", "LANG"])]
# df_keys.append(langs)
# df_keys = pd.concat(df_keys).dropna()
# df_keys = df_keys.drop_duplicates()
# # df_keys = df_keys.reset_index(drop=True)
# # df = pd.concat([df_keys[["cs"]].rename(columns={"cs": "VALUE"}), df_keys[["en"]].rename(columns={"en": "VALUE"})])

# df = df_keys.reset_index(drop=True).reset_index(drop=False).melt(id_vars=["index"])
# df.to_csv('./init_db/filter_values.csv', index=False)

# # df = df_keys.T.apply(list, axis=1)
# # df = df.apply(lambda x: "{" + ",".join([f"\"{v}\"" for v in x]) + "}")
# # df.to_csv('./init_db/filter_values_array.csv', index=True, header=False, sep=";", quotechar="'")


# # series_list = [courses[["LANG", col]].rename(columns={col:"VALUE"}) for col in courses.drop(columns=["POVINN", "LANG"])]
# # df = pd.concat(series_list).dropna().drop_duplicates().reset_index(drop=True)
# # df.to_csv('./init_db/filter_values.csv', index=True)