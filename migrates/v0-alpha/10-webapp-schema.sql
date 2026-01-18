SET search_path TO webapp;

ALTER TABLE studies ADD CONSTRAINT studies_user_id_unique_constraint UNIQUE (user_id);
ALTER TABLE studies ALTER degree_plan_code DROP NOT NULL;

ALTER TABLE studies DROP COLUMN IF EXISTS start_year;

DROP TABLE IF EXISTS degree_plans;
DROP TABLE IF EXISTS degree_plan_years;

CREATE TABLE IF NOT EXISTS degree_plans(
    plan_code VARCHAR(15) NOT NULL,
    lang CHAR(2) NOT NULL,
    title VARCHAR(250),
    valid_from INT NOT NULL,
    valid_to INT NOT NULL,
    faculty VARCHAR(5),
    section CHAR(2),
    field_code VARCHAR(20),
    study_type VARCHAR(5),
    PRIMARY KEY (plan_code, lang)
);

CREATE TABLE IF NOT EXISTS degree_plan_courses(
    plan_code VARCHAR(15) NOT NULL REFERENCES degree_plan_list(code) DEFERRABLE INITIALLY DEFERRED,
    lang CHAR(2) NOT NULL,
    course_code VARCHAR(10) NOT NULL,
    interchangeability VARCHAR(10),
    recommended_year_from INT,
    recommended_year_to INT,
    recommended_semester INT,
    bloc_name VARCHAR(250),
    bloc_subject_code VARCHAR(20) NOT NULL,
    bloc_type CHAR(1) NOT NULL,
    bloc_limit INT,
    seq VARCHAR(50),
    FOREIGN KEY (plan_code, lang) REFERENCES degree_plans(plan_code, lang) DEFERRABLE INITIALLY DEFERRED
);

CREATE INDEX IF NOT EXISTS degree_plan_code_lang ON degree_plan_courses(plan_code, lang);

DELETE FROM filter_values WHERE TRUE;
DROP TABLE IF EXISTS filter_categories CASCADE;

CREATE TABLE IF NOT EXISTS filter_categories (
    id INT PRIMARY KEY,
    filter_id VARCHAR(50) NOT NULL REFERENCES filters(id),
    facet_id VARCHAR(50) NOT NULL,
    title_cs VARCHAR(50) NOT NULL,
    title_en VARCHAR(50) NOT NULL,
    description_cs VARCHAR(200),
    description_en VARCHAR(200),
	condition VARCHAR(100),
    displayed_value_limit INT NOT NULL,
    position INT NOT NULL
);

CREATE TABLE IF NOT EXISTS requisites(
    target_course VARCHAR(10) NOT NULL,
    parent_course VARCHAR(10) NOT NULL,
    child_course VARCHAR(10) NOT NULL,
    req_type CHAR(1) NOT NULL,
    group_type VARCHAR(1)
);
