SET search_path TO webapp;

DROP TABLE IF EXISTS courses CASCADE;

CREATE TABLE courses(
    code VARCHAR(10) NOT NULL,
    lang CHAR(2) NOT NULL,
    title VARCHAR(250),
    valid_from INT NOT NULL,
    valid_to INT NOT NULL,
    course_url VARCHAR(250),
    faculty JSONB,
    department JSONB,
    taught_state CHAR(1),
    taught_state_title VARCHAR(120),
    start_semester VARCHAR(5),
    start_semester_title VARCHAR(120),
    taught_lang VARCHAR(250),
    lecture_range_winter INT,
    seminar_range_winter INT,
    lecture_range_summer INT,
    seminar_range_summer INT,
    range_unit JSONB,
    exam VARCHAR(30),
    credits INT,
    guarantors JSONB,
    teachers JSONB,
    min_occupancy VARCHAR(10),
    capacity VARCHAR(10),
    prerequisites JSONB,
    corequisites JSONB,
    incompatibilities JSONB,
    interchangeabilities JSONB,
    annotation JSONB,
    syllabus JSONB,
    terms_of_passing JSONB,
    literature JSONB,
    requirements_of_assessment JSONB,
    entry_requirements JSONB,
    aim JSONB,
    classes JSONB,
    classifications JSONB,
    PRIMARY KEY (code, lang)
);

DROP TABLE IF EXISTS degree_plan_years CASCADE;

DROP TABLE IF EXISTS degree_plan_list CASCADE;

CREATE TABLE degree_plan_list(
    code VARCHAR(15) PRIMARY KEY
);

DROP TABLE IF EXISTS degree_plans CASCADE;

CREATE TABLE degree_plans(
    plan_code VARCHAR(15) NOT NULL,
    lang CHAR(2) NOT NULL,
    title VARCHAR(250) NOT NULL,
    valid_from INT NOT NULL,
    valid_to INT NOT NULL,
    field_code VARCHAR(20) NOT NULL,
    field_title VARCHAR(250) NOT NULL,
    study_type VARCHAR(5) NOT NULL,
    required_credits INT,
    required_elective_credits INT,
    total_credits INT NOT NULL,
    studying JSONB,
    graduates JSONB,
    next_plans JSONB,
    previous_plans JSONB,
    requisite_graph_data TEXT,
    PRIMARY KEY (plan_code, lang)
);

DROP TABLE IF EXISTS degree_plan_courses CASCADE;

CREATE TABLE degree_plan_courses(
    plan_code VARCHAR(15) NOT NULL REFERENCES degree_plan_list(code) DEFERRABLE INITIALLY DEFERRED,
    lang CHAR(2) NOT NULL,
    course_code VARCHAR(10) NOT NULL,
    interchangeability VARCHAR(10),
    recommended_year_from INT,
    recommended_year_to INT,
    recommended_semester INT,
    bloc_name VARCHAR(250),
    bloc_subject_code VARCHAR(20) NOT NULL,
    is_required BOOLEAN NOT NULL,
    is_elective BOOLEAN NOT NULL,
    bloc_limit INT NOT NULL,
    seq VARCHAR(50),
    FOREIGN KEY (plan_code, lang) REFERENCES degree_plans(plan_code, lang) DEFERRABLE INITIALLY DEFERRED
);

CREATE INDEX degree_plan_code_lang ON degree_plan_courses(plan_code, lang);

DROP TABLE IF EXISTS filters CASCADE;

CREATE TABLE filters (
    id VARCHAR(50) PRIMARY KEY
);

DROP TABLE IF EXISTS filter_categories CASCADE;

CREATE TABLE filter_categories (
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

DROP TABLE IF EXISTS filter_values CASCADE;

CREATE TABLE filter_values (
    id INT PRIMARY KEY,
    category_id INT NOT NULL REFERENCES filter_categories(id),
    facet_id VARCHAR(50) NOT NULL,
    title_cs VARCHAR(250) NOT NULL,
    title_en VARCHAR(250) NOT NULL,
    description_cs VARCHAR(200),
    description_en VARCHAR(200),
    position INT NOT NULL
);