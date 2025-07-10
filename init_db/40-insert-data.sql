\c recsis
SET search_path TO webapp;

COPY degree_plans(plan_code, plan_year, course_code, interchangeability, bloc_subject_code, bloc_type, bloc_limit, seq, bloc_name, bloc_note, note, lang)
FROM '/docker-entrypoint-initdb.d/degree_plans_transformed.csv'
DELIMITER ','
CSV HEADER;

INSERT INTO course_rating_categories(code, lang, title)
VALUES
    (1, 'cs', 'Náročnost'),
    (2, 'cs', 'Přínosnost'),
    (3, 'cs', 'Zajímavost'),
    (4, 'cs', 'Zábava'),
    (5, 'cs', 'Zátěž'),
    (1, 'en', 'Difficulty'),
    (2, 'en', 'Usefulness'),
    (3, 'en', 'Interest'),
    (4, 'en', 'Fun'),
    (5, 'en', 'Workload');