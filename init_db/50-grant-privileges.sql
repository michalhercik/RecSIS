\c recsis

GRANT USAGE ON SCHEMA webapp TO webapp;
GRANT SELECT, INSERT, DELETE ON ALL TABLES IN SCHEMA webapp TO webapp;

GRANT USAGE ON SCHEMA webapp TO elt;
GRANT
    DELETE,
    INSERT
ON
    webapp.courses,
    webapp.degree_plans,
    webapp.filters,
    webapp.filter_categories,
    webapp.filter_values
TO elt;