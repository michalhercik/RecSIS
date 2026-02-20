import pandas as pd

class User:
    def __init__(self, id: str, degree_plan: str, enrollment_year: int, blueprint):
        self.id = id
        self.degree_plan = degree_plan # Degree plan code e.g. "NISD23N"
        self.enrollment_year = enrollment_year
        self.blueprint = sorted(blueprint, key=lambda x: x["year"]) if blueprint else []
        self.fetch = False

        # """ Blueprint format example:
        # [
        #     {"year": 0, "unassigned": ["NDBI021", ...]},
        #     {"year": 1, "winter": ["NPRG021", ...], "summer": [...]},
        #     ...
        # ]
        # """

    # """
    # Converts the user's blueprint into a DataFrame with columns: year, semester, course
    # """
    def blueprint_to_df(self) -> pd.DataFrame:
        if len(self.blueprint) == 0:
            return pd.DataFrame(columns=["year", "semester", "course"])
        records = []
        for yearRecord in self.blueprint:
            for semester in ["winter", "summer", "unassigned"]:
                if semester in yearRecord:
                    for course in yearRecord[semester]:
                        records.append({"year": yearRecord["year"], "semester": semester, "course": course})
        return pd.DataFrame(records, columns=["year", "semester", "course"])

    def finished(self) -> list[str]:
        courses = []
        for yearRecord in self.blueprint:
            for semester in ["winter", "summer", "unassigned"]:
                if semester in yearRecord:
                    courses.extend(yearRecord[semester])
        return courses
