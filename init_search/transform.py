import pandas as pd


def converter(data):
    import json
    if len(data) == 0:
        return None
    return json.loads(data)

def com_agg(row):
    data = {
        "course_code": row.iloc[0]["POVINN"],
        "target_type": row.iloc[0]["PRDMTYP"],
        "content": row.iloc[0]["MEMO"],
        "study_year": row.iloc[0]["SROC"],
        "academic_year": row.iloc[0]["SSKR"],
        "semester": row.iloc[0]["SEM"],
        "study_field": row.iloc[0]["SOBOR"],
        "teacher": row.iloc[0]["TEACHER"],
        "teacher_facet": row.iloc[0]["TEACHER"]["JMENO"] + " " + row.iloc[0]["TEACHER"]["PRIJMENI"] if not pd.isna(row.iloc[0]["TEACHER"]) else None,
        "study_type": {
            "code": row.iloc[0]["KOD"],
            "abbr": row.iloc[0]["ZKRATKA"],
            "name_cs": row[row["LANG"] == "cs"]["NAZEV"].iloc[0],
            "name_en": row[row["LANG"] == "en"]["NAZEV"].iloc[0]
        }
    }
    return pd.Series(data=data, index=data.keys())


# comments = pd.read_csv('./init_db/ankecy_transformed.csv', converters={"TEACHER": converter}, dtype={"SROC": pd.Int32Dtype()})
# comments = comments.groupby(["POVINN", "SOBOR", "SSKR", "SROC", "SEM", "KOD"]).apply(com_agg)
# comments = comments.reset_index(drop=True).reset_index().rename(columns={"index": "id"})
# comments = comments.to_json('./init_search/comments.json', orient='records', lines=True)

povinn = pd.read_csv('./init_db/POVINN.csv')
courses = pd.read_csv('./init_db/courses_transformed.csv', usecols=[
    "POVINN", "NAME", "LANG", "PGARANT", "PVYUCOVAN", "PPOCMAX", "PPOCMIN",
    "VYJAZYK", "FACULTY_NAME", "GUARANTORS", "TEACHERS",
    "A", "C", "E", "L", "S", "V", "P_x",
    "VSEMZAC", "VSEMPOC", "VTYP", "VRVCEM", "VEBODY",
    "LECTURE_RANGE_WINTER", "SEMINAR_RANGE_WINTER", "LECTURE_RANGE_SUMMER", "SEMINAR_RANGE_SUMMER",
    # "VROZSAHPR1", "VROZSAHCV1", "VROZSAHPR2", "VROZSAHCV2",
    ], dtype={
        "VEBODY": pd.Int32Dtype(),
        "VSEMZAC": pd.Int32Dtype(),
        "VSEMPOC": pd.Int32Dtype(),
        "PPOCMIN": pd.Int32Dtype(),
        "PPOCMAX": str,
        "LECTURE_RANGE_WINTER": pd.Int32Dtype(),
        "SEMINAR_RANGE_WINTER": pd.Int32Dtype(),
        "LECTURE_RANGE_SUMMER": pd.Int32Dtype(),
        "SEMINAR_RANGE_SUMMER": pd.Int32Dtype(),
    }, converters={
    "A": converter,
    "C": converter,
    "E": converter,
    "L": converter,
    "S": converter,
    "V": converter,
    "P_x": converter,
    "GUARANTORS": converter,
    "TEACHERS": converter
})
# values = pd.read_csv('./init_db/filter_values.csv', usecols=["index", "value"], header=0, dtype={"index": int}).set_index("value")["index"].to_dict()
# values = pd.read_csv('./init_db/filter_values.csv', usecols=[0, 1, 2, 3], dtype={0: int})#, column_names=["id", "title_cs", "title_en", "category"])
# categories = pd.read_csv('./init_db/filter_categories.csv', usecols=[0, 1], dtype={0: int}, index_col=1).iloc[:, 0].to_dict()

def memo(data):
    if data is not None:
        return data["MEMO"]
    return None

def select(row, lang):
    # TODO: exam_type, faculty, taught, taught_lang, lecture/seminar range, range_unit, capacity, min_number,
    return pd.concat([row[row["LANG"]==lang][["A", "C", "E", "S", "V", "P_x"]].map(memo), row[row["LANG"]==lang][["NAME"]]], axis=1)

def select_teachers(row):
    return [{"JMENO": y["JMENO"], "PRIJMENI": y["PRIJMENI"]} for y in row] if row is not None else None

# def id(value):
#     if pd.isna(value):
#         return value
#     elif isinstance(value, str):
#         return values[value]
#     return values[value.astype(str)]

def id(row, col):
    # value = row.iloc[0][col]
    # category = col2category[col]
    # print(row[["NAME", "LANG"]])
    v = row.iloc[1][col]
    return v

def lang_id(row, col):
    value = row.iloc[1][col]
    if pd.isna(value):
        return value
    return [format_str(x) for x in value.split(", ")]

format_str = lambda x: x.lower().replace(" ", "_") if not pd.isna(x) else x

def aggregate(row):
    data={
        "code": row.iloc[0]["POVINN"],
        "start_semester": id(row, "VSEMZAC"),
        "semester_count": id(row, "VSEMPOC"),
        "lecture_range": [id(row, "LECTURE_RANGE_WINTER"), id(row, "LECTURE_RANGE_SUMMER")],
        "seminar_range": [id(row, "SEMINAR_RANGE_WINTER"), id(row, "SEMINAR_RANGE_SUMMER")],
        # "lecture_range_winter": id(row, "LECTURE_RANGE_WINTER"),
        # "seminar_range_winter": id(row, "SEMINAR_RANGE_WINTER"),
        # "lecture_range_summer": id(row, "LECTURE_RANGE_SUMMER"),
        # "seminar_range_summer": id(row, "SEMINAR_RANGE_SUMMER"),
        "credits": id(row, "VEBODY"),
        "department": format_str(id(row, "PGARANT")),
        "exam_type": format_str(id(row, "VTYP")),
        "range_unit": format_str(id(row, "VRVCEM")),
        "taught": format_str(id(row, "PVYUCOVAN")),
        "language": lang_id(row, "VYJAZYK"),
        "faculty": format_str(id(row, "FACULTY_NAME")),
        "capacity": id(row, "PPOCMAX"),
        "min_occupancy": id(row, "PPOCMIN"),
        "cs": select(row, "cs").to_dict(orient="records")[0],
        "en": select(row, "en").to_dict(orient="records")[0],
        "guarantors": select_teachers(row.iloc[0]["GUARANTORS"]),
        "teachers": select_teachers(row.iloc[0]["TEACHERS"])
    }
    return pd.Series(data=data, index=data.keys())

courses = courses[~courses["PVYUCOVAN"].isin(["Zrušen", "Cancelled"])]
courses_data = courses[pd.notna(courses["VEBODY"])].copy()
courses_data["POVINNGROUP"] = courses_data["POVINN"]
courses_data = courses_data.groupby(["POVINNGROUP"]).apply(aggregate, include_groups=False)
courses_data = courses_data.reset_index(drop=True).reset_index(drop=False).rename(columns={"index": "id"})
courses_data.to_json('./init_search/courses.json', orient='records', lines=True)

# df = courses_data[[
#     "faculty", "department", "exam_type", "taught", "language", "start_semester", "semester_count",
#     "lecture_range_winter", "seminar_range_winter", "lecture_range_summer", "seminar_range_summer",
#     "range_unit", "credits", "capacity", "min_occupancy"]]
# df = df.melt().explode(["value"]).dropna().drop_duplicates()
# df["value"] = df["value"].astype(pd.Int32Dtype())
# df.to_csv("./init_db/filter_params.csv", index=False, header=False)
