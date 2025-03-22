import pandas as pd

def converter(data):
    import json
    if len(data) == 0:
        return None
    return json.loads(data)


povinn = pd.read_csv('./init_db/POVINN.csv')
courses = pd.read_csv('./init_db/courses_transformed.csv', usecols=[
    "POVINN", "NAME", "LANG", "PGARANT", "PVYUCOVAN", "PPOCMAX", "PPOCMIN",
    "VYJAZYK", "FACULTY_NAME", "GUARANTORS", "TEACHERS",
    "A", "C", "E", "L", "S", "V", "P",
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
    "P": converter,
    "GUARANTORS": converter,
    "TEACHERS": converter
})
values = pd.read_csv('./init_db/filter_values.csv', usecols=["index", "value"], header=0, dtype={"index": int}).set_index("value")["index"].to_dict()

def memo(data):
    if data is not None:
        return data["MEMO"]
    return None

def select(row, lang):
    # TODO: exam_type, faculty, taught, taught_lang, lecture/seminar range, range_unit, capacity, min_number,
    return pd.concat([row[row["LANG"]==lang][["A", "C", "E", "S", "V", "P"]].map(memo), row[row["LANG"]==lang][["NAME"]]], axis=1)

def select_teachers(row):
    return [{"JMENO": y["JMENO"], "PRIJMENI": y["PRIJMENI"]} for y in row] if row is not None else None

def id(value):
    if pd.isna(value):
        return value
    elif isinstance(value, str):
        return values[value]
    return values[value.astype(str)]

def lang_id(value):
    if pd.isna(value):
        return value
    return [values[x] for x in value.split(", ")]

def aggregate(row):
    data={
        "code": row.iloc[0]["POVINN"],
        "start_semester": id(row.iloc[0]["VSEMZAC"]),
        "semester_count": id(row.iloc[0]["VSEMPOC"]),
        "lecture_range_winter": id(row.iloc[0]["LECTURE_RANGE_WINTER"]),
        "seminar_range_winter": id(row.iloc[0]["SEMINAR_RANGE_WINTER"]),
        "lecture_range_summer": id(row.iloc[0]["LECTURE_RANGE_SUMMER"]),
        "seminar_range_summer": id(row.iloc[0]["SEMINAR_RANGE_SUMMER"]),
        "credits": id(row.iloc[0]["VEBODY"]),
        "faculty_guarantor": id(row.iloc[0]["PGARANT"]),
        "exam_type": id(row.iloc[0]["VTYP"]),
        "range_unit": id(row.iloc[0]["VRVCEM"]),
        "taught": id(row.iloc[0]["PVYUCOVAN"]),
        "taught_lang": lang_id(row.iloc[0]["VYJAZYK"]),
        "faculty": id(row.iloc[0]["FACULTY_NAME"]),
        "capacity": id(row.iloc[0]["PPOCMAX"]),
        "min_number": id(row.iloc[0]["PPOCMIN"]),
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

df = courses_data[[
    "faculty", "faculty_guarantor", "exam_type", "taught", "taught_lang",
    "lecture_range_winter", "seminar_range_winter", "lecture_range_summer", "seminar_range_summer",
    "range_unit", "credits", "capacity", "min_number"]]
df = df.melt().explode(["value"]).dropna().drop_duplicates()
df.to_csv("./init_db/filter_params.csv", index=False, header=False)