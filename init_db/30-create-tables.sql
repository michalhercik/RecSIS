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

CREATE TABLE courses(
    id SERIAL PRIMARY KEY,
    code VARCHAR(10) NOT NULL,
    name_cs VARCHAR(100) NOT NULL,
    name_en VARCHAR(100) NOT NULL,
    valid_from INT NOT NULL,
    valid_to INT NOT NULL,
    faculty INT REFERENCES faculties(sis_id), -- INT REFERENCES faculties(id),
    guarantor VARCHAR(10), -- INT REFERENCES departments(id),
    taught CHAR(1) NOT NULL,
    start_semester INT NOT NULL,
    semester_count INT NOT NULL,
    taught_lang CHAR(3),
    lecture_range1 INT NOT NULL,
    seminar_range1 INT NOT NULL,
    lecture_range2 INT,
    seminar_range2 INT,
    range_unit CHAR(2),
    exam_type VARCHAR(2) NOT NULL,
    credits INT NOT NULL,
    teacher1 INT, --REFERENCES teachers(id),
    teacher2 INT, --REFERENCES teachers(id),
    min_number INT,
    capacity INT,
    annotation_cs TEXT,
    annotation_en TEXT,
    sylabus_cs TEXT,
    sylabus_en TEXT
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

CREATE TABLE blueprint_years(
    id SERIAL PRIMARY KEY,
    student INT NOT NULL,
    position INT NOT NULL
);

CREATE TABLE blueprint_semesters(
    blueprint_year INT REFERENCES blueprint_years(id),
    course INT REFERENCES courses(id) NOT NULL,
    semester INT NOT NULL,
    position INT NOT NULL
);