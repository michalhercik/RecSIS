SET search_path TO webapp;

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
    annotation JSONB,
    syllabus JSONB,
    terms_of_passing JSONB,
    literature JSONB,
    requirements_of_assesment JSONB,
    entry_requirements JSONB,
    aim JSONB,
    prerequisities JSONB,
    corequisities JSONB,
    incompatibilities JSONB,
    interchangeabilities JSONB,
    classes JSONB,
    classifications JSONB
);

CREATE INDEX courses_code_lang_idx ON courses(code, lang);

CREATE TABLE degree_plan_years(
    plan_year INT PRIMARY KEY
);

CREATE TABLE degree_plan_list(
    code VARCHAR(15) PRIMARY KEY
);

CREATE TABLE degree_plans(
    plan_code VARCHAR(15) NOT NULL REFERENCES degree_plan_list(code) DEFERRABLE INITIALLY DEFERRED,
    plan_year INT NOT NULL REFERENCES degree_plan_years(plan_year) DEFERRABLE INITIALLY DEFERRED,
    lang CHAR(2) NOT NULL,
    course_code VARCHAR(10) NOT NULL,
    interchangeability VARCHAR(10),
    bloc_subject_code VARCHAR(20) NOT NULL,
    bloc_type CHAR(1) NOT NULL,
    bloc_limit INT,
    bloc_name VARCHAR(250),
    bloc_note VARCHAR(1000),
    note VARCHAR(250),
    seq VARCHAR(50)
);

CREATE INDEX degree_plan_code_year_lang ON degree_plans(plan_code, plan_year, lang);

CREATE TABLE users (
    id VARCHAR(8) PRIMARY KEY
);

CREATE TABLE blueprint_years (
    id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id VARCHAR(8) NOT NULL,
    academic_year INT NOT NULL,
    UNIQUE (user_id, academic_year),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);


CREATE TABLE blueprint_semesters(
    id int GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    blueprint_year_id INT NOT NULL,
    semester INT NOT NULL,
    folded BOOLEAN NOT NULL DEFAULT false,
    FOREIGN KEY (blueprint_year_id) REFERENCES blueprint_years(id) ON DELETE CASCADE,
    UNIQUE (blueprint_year_id, semester)
);

CREATE TABLE blueprint_courses(
    id int GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    blueprint_semester_id INT NOT NULL,
    course_code VARCHAR(10) NOT NULL,
    course_valid_from INT NOT NULL,
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

CREATE TABLE studies (
    id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id VARCHAR(8) NOT NULL,
    degree_plan_code VARCHAR(15) NOT NULL REFERENCES degree_plan_list(code) DEFERRABLE INITIALLY DEFERRED,
    start_year INT NOT NULL REFERENCES degree_plan_years(plan_year) DEFERRABLE INITIALLY DEFERRED,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);


-- TODO: Clean up expired sessions
CREATE TABLE sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id VARCHAR(8),
    ticket VARCHAR(42) NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE course_rating_categories_domain (
    code INT PRIMARY KEY
);

CREATE TABLE course_rating_categories (
    code INT REFERENCES course_rating_categories_domain(code),
    lang CHAR(2) NOT NULL,
    title VARCHAR(50) NOT NULL
);

CREATE DOMAIN course_overall_rating_domain AS INT CHECK (VALUE = 0 OR VALUE = 1);

CREATE TABLE course_overall_ratings (
    user_id VARCHAR(8) NOT NULL,
    course_code VARCHAR(10) NOT NULL,
    rating course_overall_rating_domain NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (user_id, course_code),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE DOMAIN course_rating_domain AS INT CHECK (VALUE >= 0 AND VALUE <= 10);

CREATE TABLE course_ratings (
    user_id VARCHAR(8) NOT NULL,
    course_code VARCHAR(10) NOT NULL,
    category_code INT NOT NULL REFERENCES course_rating_categories_domain(code),
    rating course_rating_domain NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (user_id, category_code, course_code),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);


CREATE TABLE filters (
    id VARCHAR(50) PRIMARY KEY
);

CREATE TABLE filter_categories (
    id INT PRIMARY KEY,
    filter_id VARCHAR(50) NOT NULL REFERENCES filters(id),
    facet_id VARCHAR(50) NOT NULL,
    title_cs VARCHAR(50) NOT NULL,
    title_en VARCHAR(50) NOT NULL,
    description_cs VARCHAR(200),
    description_en VARCHAR(200),
    displayed_value_limit INT NOT NULL,
    position INT NOT NULL
);

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

SET search_path TO recommender;

CREATE TABLE povinn (
    povinn     VARCHAR(10),
    pnazev     VARCHAR(250),
    panazev    VARCHAR(250),
    vplatiod   INT,
    vplatido   INT,
    pfakulta   VARCHAR(5),
    pgarant    VARCHAR(10),
    pvyucovan  VARCHAR(1),
    vsemzac    VARCHAR(1),
    vsempoc    INT,
    vrozsahpr1 INT,
    vrozsahcv1 INT,
    vrozsahpr2 INT,
    vrozsahcv2 INT,
    vrvcem     VARCHAR(2),
    vtyp       VARCHAR(2),
    vebody     INT,
    vucit1     VARCHAR(10),
    vucit2     VARCHAR(10),
    vucit3     VARCHAR(10)
);

CREATE TABLE studium (
    soident INT,
    sident INT,
    sfak VARCHAR(5),
    sfak2 VARCHAR(5),
    sdruh VARCHAR(2),
    sobor VARCHAR(12),
    srokp VARCHAR(4),
    sstav VARCHAR(6),
    sroc INT,
    splan VARCHAR(15)
);

CREATE TABLE zkous (
    zident INT,
    zskr VARCHAR(4),
    zsem VARCHAR(1),
    zpovinn VARCHAR(10),
    zmarx VARCHAR(5),
    zroc INT,
    zbody INT,
    zsplcelk VARCHAR(1)
);

CREATE TABLE stud_plan (
    code VARCHAR(10),
    interchangeability VARCHAR(10),
    bloc_subject_code VARCHAR(20),
    bloc_type VARCHAR(1),
    bloc_grade VARCHAR(50),
    bloc_limit INT,
    bloc_name_cz VARCHAR(250),
    bloc_name_en VARCHAR(250),
    plan_code VARCHAR(15),
    plan_year INT
)