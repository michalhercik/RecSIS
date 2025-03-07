CREATE TABLE faculties(
    id SERIAL PRIMARY KEY,
    sis_id INT UNIQUE NOT NULL,
    sis_poid INT,
    name_cs VARCHAR(150) NOT NULL,
    name_en VARCHAR(150) NOT NULL,
    abbr VARCHAR(10) NOT NULL
);

-- CREATE TABLE departments (
--     id SERIAL PRIMARY KEY,
--     sis_id VARCHAR(10) UNIQUE NOT NULL
-- );

CREATE TABLE teachers(
    id SERIAL PRIMARY KEY,
    sis_id INT UNIQUE NOT NULL,
    department VARCHAR(10), -- INT REFERENCES departments(id),
    faculty INT REFERENCES faculties(sis_id), -- INT REFERENCES faculties(id),
    first_name VARCHAR(50),
    last_name VARCHAR(50) NOT NULL,
    title_before VARCHAR(20),
    title_after VARCHAR(20)
);

CREATE TABLE old_courses(
    id SERIAL PRIMARY KEY,
    code VARCHAR(10) NOT NULL,
    name_cs VARCHAR(250) NOT NULL,
    name_en VARCHAR(250), -- NOT NULL,
    valid_from INT NOT NULL,
    valid_to INT NOT NULL,
    faculty INT REFERENCES faculties(sis_id), -- INT REFERENCES faculties(id),
    guarantor VARCHAR(10), -- INT REFERENCES departments(id),
    taught CHAR(1) NOT NULL,
    start_semester INT, -- NOT NULL,
    semester_count INT, -- NOT NULL,
    taught_lang CHAR(3), -- change to CHAR(2) NOT NULL (cs, en),
    lecture_range1 INT, -- NOT NULL,
    seminar_range1 INT, -- NOT NULL,
    lecture_range2 INT,
    seminar_range2 INT,
    range_unit CHAR(2),
    exam_type VARCHAR(2), -- NOT NULL,
    credits INT, -- NOT NULL,
    teacher1 INT, --REFERENCES teachers(sis_id), --REFERENCES teachers(id),
    teacher2 INT, --REFERENCES teachers(sis_id), --REFERENCES teachers(id),
    teacher3 INT, --REFERENCES teachers(sis_id), --REFERENCES teachers(id),
    min_number INT,
    capacity INT
);

CREATE TABLE start_semester_to_desc(
    id INT,
    lang CHAR(2),
    semester_description VARCHAR(7)
);

CREATE TABLE courses(
    -- id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    code VARCHAR(10) NOT NULL,
    lang CHAR(2) NOT NULL,
    title VARCHAR(250),
    valid_from INT NOT NULL,
    valid_to INT NOT NULL,
    faculty VARCHAR(150), -- INT REFERENCES faculties(id),
    guarantor VARCHAR(10), -- INT REFERENCES departments(id),
    taught VARCHAR(10) NOT NULL,
    start_semester INT, -- NOT NULL,
    semester_count INT, -- NOT NULL,
    taught_lang VARCHAR(10), -- change to CHAR(2) NOT NULL (cs, en),
    lecture_range1 INT, -- NOT NULL,
    seminar_range1 INT, -- NOT NULL,
    lecture_range2 INT,
    seminar_range2 INT,
    range_unit CHAR(2),
    exam_type VARCHAR(4), -- NOT NULL,
    credits INT, -- NOT NULL,
    guarantors JSONB,
    teachers JSONB,
    min_number INT,
    capacity VARCHAR(9),
    annotation JSONB,
    syllabus JSONB,
    aim JSONB,
    requirements JSONB
);

CREATE TABLE classifications(
    course VARCHAR(10),
    classification VARCHAR(6) NOT NULL
);

CREATE TABLE classes(
    course VARCHAR(10),
    class VARCHAR(7) NOT NULL
);

CREATE TABLE requisites(
    course VARCHAR(10) NOT NULL,
    requisite VARCHAR(10) NOT NULL,
    requisite_type CHAR(1) NOT NULL,
    from_year INT NOT NULL,
    to_year INT NOT NULL
);

CREATE TABLE course_texts(
    id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    course VARCHAR(10) NOT NULL,
    text_type CHAR(1) NOT NULL,
    lang CHAR(3) NOT NULL,
    title VARCHAR(120) NOT NULL,
    content TEXT NOT NULL,
    audience VARCHAR(6) NOT NULL,
    UNIQUE (course, text_type, lang)
);

CREATE TABLE course_teachers(
    course VARCHAR(10) NOT NULL,
    teacher INT NOT NULL -- REFERENCES teachers(sis_id) NOT NULL,
);

CREATE TABLE degree_plans(
    plan_code VARCHAR(15) NOT NULL,
    plan_year INT NOT NULL,
    lang CHAR(2) NOT NULL,
    course_code VARCHAR(10) NOT NULL,
    interchangeability VARCHAR(10),
    bloc_subject_code INT NOT NULL,
    bloc_type CHAR(1) NOT NULL,
    bloc_limit INT,
    bloc_name VARCHAR(250),
    bloc_note VARCHAR(1000),
    note VARCHAR(250),
    seq VARCHAR(50)
);

CREATE TABLE degree_programmes(
    id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    code CHAR(12) NOT NULL,
    name_cs VARCHAR(350) NOT NULL,
    name_en VARCHAR(350) NOT NULL,
    faculty INT REFERENCES faculties(sis_id) NOT NULL,
    program_type CHAR(1) NOT NULL,
    program_form CHAR(1) NOT NULL,
    graduate_profile_cs TEXT, --NOT NULL,
    graduate_profile_en TEXT, --NOT NULL,
    lang CHAR(2) NOT NULL
);

CREATE TABLE studies(
    id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    sis_id INT NOT NULL,
    student INT NOT NULL,
    faculty1 INT REFERENCES faculties(sis_id) NOT NULL,
    faculty2 INT REFERENCES faculties(sis_id),
    study_type CHAR(1) NOT NULL,
    study_form CHAR(1) NOT NULL,
    study_specialization VARCHAR(12) NOT NULL,
    enrollment DATE NOT NULL,
    study_state CHAR(1) NOT NULL,
    study_state_date DATE NOT NULL,
    study_year INT NOT NULL,
    degree_plan VARCHAR(15) --INT REFERENCES degree_plans(id) NOT NULL
);

CREATE TABLE users (
    id INT PRIMARY KEY
);

CREATE TABLE blueprint_years(
    id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id),
    academic_year INT NOT NULL,
    UNIQUE (user_id, academic_year)
);

CREATE TABLE blueprint_semesters(
    id int GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    blueprint_year_id INT NOT NULL,
    semester INT NOT NULL,
    FOREIGN KEY (blueprint_year_id) REFERENCES blueprint_years(id) ON DELETE CASCADE,
    UNIQUE (blueprint_year_id, semester)
);

CREATE TABLE blueprint_courses(
    id int GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    blueprint_semester_id INT NOT NULL,
    course_code VARCHAR(10) NOT NULL,
    course_valid_from INT NOT NULL,
    -- course INT REFERENCES courses(id) NOT NULL,
    position INT NOT NULL,
    secondary_position INT NOT NULL DEFAULT 2,
    FOREIGN KEY (blueprint_semester_id) REFERENCES blueprint_semesters(id) ON DELETE CASCADE,
    UNIQUE (blueprint_semester_id, course_code),
    UNIQUE (blueprint_semester_id, position) DEFERRABLE INITIALLY DEFERRED
);

CREATE OR REPLACE FUNCTION blueprint_course_reordering()
   RETURNS TRIGGER
AS
$$
BEGIN
    UPDATE blueprint_courses
    SET position = new_position, secondary_position = 2
    FROM (
        SELECT
            id as sub_id,
            position,
            ROW_NUMBER() OVER (
                PARTITION BY blueprint_semester_id
                ORDER BY position, secondary_position
            ) AS new_position
        FROM blueprint_courses
        WHERE (
            CASE
                WHEN NEW IS NULL THEN false
                ELSE (blueprint_semester_id = NEW.blueprint_semester_id)
            END
            OR (blueprint_semester_id = OLD.blueprint_semester_id)
        )
    )
    WHERE id=sub_id;
    RETURN CASE
        WHEN NEW IS NULL THEN OLD
        ELSE NEW
    END;
END;
$$ LANGUAGE PLPGSQL;

CREATE TRIGGER blueprint_course_reordering_trigger
AFTER UPDATE OR DELETE ON blueprint_courses
FOR EACH ROW
WHEN (pg_trigger_depth() = 0)
EXECUTE FUNCTION blueprint_course_reordering();

CREATE TABLE bla_blueprints (
    user_id INT REFERENCES users(id),
    lang CHAR(2) NOT NULL,
    blueprint JSONB NOT NULL
);

-- TODO: rename
CREATE TABLE bla_studies (
    id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id INT REFERENCES users(id) NOT NULL,
    degree_plan_code VARCHAR(15) NOT NULL,
    start_year INT NOT NULL
);

-- TODO: Clean up expired sessions
CREATE TABLE sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id INT REFERENCES users(id),
    expires_at TIMESTAMPTZ NOT NULL
);