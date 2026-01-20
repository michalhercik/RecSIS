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
    webapp.filters,
    webapp.filter_categories,
    webapp.filter_values,
    webapp.requisites
TO elt;
