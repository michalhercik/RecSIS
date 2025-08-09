\c recsis

GRANT USAGE ON SCHEMA webapp TO webapp;
GRANT SELECT, INSERT, DELETE, UPDATE ON ALL TABLES IN SCHEMA webapp TO webapp;

GRANT USAGE ON SCHEMA webapp TO elt;
GRANT
    DELETE,
    INSERT
ON
    webapp.courses,
    webapp.degree_plan_list,
    webapp.degree_plan_years,
    webapp.degree_plans,
    webapp.filters,
    webapp.filter_categories,
    webapp.filter_values
TO elt;

GRANT USAGE ON SCHEMA webapp TO recommender;
GRANT SELECT ON ALL TABLES IN SCHEMA webapp TO recommender;