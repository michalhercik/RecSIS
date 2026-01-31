import numpy as np
import pandas as pd

from algo.base import Algorithm
from user import User


# NOTE: since there is no table denoting student course registrations per semester, we can imply course registration if there was at least one exam attempt per student for the course,
# but what if the course doesn't have exams only credit? or different type of completion requirements?
class PopularityBasedOnExamAttempts(Algorithm):
    """Recommend courses based on their examination popularity.
    The algorithm ranks courses by the number of distinct student attempts
    (one attempt per student per semester per course). It filters out
    irrelevant courses (e.g., PE, thesis, languages, canceled) and, at
    recommendation time, excludes courses already completed, compulsory
    courses from the student's degree plan, and their prerequisite-related
    requirements. Returns the top-N remaining courses.
    """

    def fit(self):
        # Aggregate to at most one attempt per student per semester per course
        attempts_df = self.data.zkous.groupby(
            ["zident", "zskr", "zsem", "zpovinn"], as_index=False
        ).agg(passed=("zsplcelk", lambda x: 1 if (x == "S").any() else 0))
        attempts_df["attempted"] = 1
        # Calculate popularity of each course based on examination data
        course_stats = (
            attempts_df.groupby("zpovinn")
            .agg(attempts=("attempted", "sum"))
            .rename_axis("course")
            .reset_index()
            .sort_values("attempts", ascending=False)
        )

        # Manual filtering of irrelevant courses (PE, thesis, languages, etc.)
        excluded_course_codes = {
            # Add exact codes here if needed
        }
        code_pattern = r"^(NTVY|NSZ|NSZB|NJAZ)"

        code_mask = course_stats["course"].str.match(code_pattern, na=False)
        manual_mask = course_stats["course"].isin(excluded_course_codes)
        cancelled_courses = course_stats["course"].isin(
            self.data.searchable_povinn["povinn"]
        )

        is_irrelevant = code_mask | manual_mask | cancelled_courses
        course_stats = course_stats[~is_irrelevant]

        course_stats = course_stats.sort_values("attempts", ascending=False)
        self.popularity_ranking = course_stats["course"]

    def recommend(self, user: User, limit: int) -> list[str]:
        # Get list of courses the user has already completed
        completed_courses = user.blueprint_to_df()["course"]

        # Get compulsory courses from the user's degree plan (identified by plan_code and plan_year)
        all_plans_rows = self.data.stud_plan[
            self.data.stud_plan["plan_code"] == user.degree_plan
        ]
        plan_year = all_plans_rows["plan_year"][
            all_plans_rows["plan_year"] <= user.enrollment_year
        ].max()
        plan_rows = all_plans_rows[all_plans_rows["plan_year"] == plan_year]
        compulsory = plan_rows[plan_rows["bloc_type"] == "A"]
        compulsory = pd.concat(
            [compulsory["code"], compulsory["interchangeability"].dropna()],
            ignore_index=True,
        ).drop_duplicates()

        irrelevant_courses = pd.concat(
            [completed_courses, compulsory]
        ).drop_duplicates()
        
        # Filter out courses that are prerequisities (reqtyp = P), non-compatible (reqtyp = N), or interchangeable (reqtyp = Z) with the irrelevant courses
        irrelevant_requisities = self.data.preq[
            (self.data.preq["povinn"].isin(irrelevant_courses))
            & (self.data.preq["reqtyp"].isin(["P", "N", "Z"]))
        ]["reqpovinn"]

        # TODO: filter out courses from other sections/departments 
        # (e.g. degree_plan NISD23N belongs to informatics sections, filter only those or facet search/dropdown if other sections might be relevant for the user. Table missing currently, only meilisearch index available)

        recommendations = [
            course
            for course in self.popularity_ranking
            if course not in irrelevant_courses and course not in irrelevant_requisities
        ]
        return recommendations[:limit]


class PopularityBasedOnCourseReviews(Algorithm):
    """Recommendation algorithm that ranks courses by the popularity derived from course review sentiment analysis.
    """
    pass
