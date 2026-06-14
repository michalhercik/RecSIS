SET search_path TO public;

GRANT SELECT, INSERT, DELETE, UPDATE ON ALL TABLES IN SCHEMA webapp TO webapp;

GRANT USAGE ON SCHEMA webapp TO elt;
GRANT
    DELETE,
    INSERT
ON
    webapp.courses,
    webapp.degree_plan_courses,
    webapp.degree_plan_list,
    webapp.degree_plans,
    webapp.filter_categories,
    webapp.filter_values,
    webapp.filters
TO elt;

GRANT USAGE ON SCHEMA recommender TO elt;
GRANT
    DELETE,
    INSERT
ON
    recommender.povinn,
    recommender.studium,
    recommender.zkous,
    recommender.stud_plan,
    recommender.searchable_povinn,
    recommender.preq,
    recommender.pamela
TO elt;
