SET search_path TO webapp;

DROP TABLE IF EXISTS degree_plans;

CREATE TABLE IF NOT EXISTS degree_plans(
    plan_code VARCHAR(15) NOT NULL,
    lang CHAR(2) NOT NULL,
    title VARCHAR(250),
    valid_from INT NOT NULL,
    valid_to INT NOT NULL,
    field_code VARCHAR(20),
    field_title VARCHAR(250),
    study_type VARCHAR(5),
    PRIMARY KEY (plan_code, lang)
);

DROP TABLE IF EXISTS requisites;