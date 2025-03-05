COPY faculties(sis_id, sis_poid, name_cs, name_en, abbr)
FROM '/docker-entrypoint-initdb.d/faculties.csv'
DELIMITER ','
CSV HEADER;

COPY teachers(sis_id,department,faculty,last_name,first_name,title_before,title_after)
FROM '/docker-entrypoint-initdb.d/teachers.csv'
DELIMITER ','
CSV HEADER;

COPY courses(code,name_cs,name_en,valid_from,valid_to,faculty,guarantor,taught,start_semester,semester_count,taught_lang,lecture_range1,seminar_range1,lecture_range2,seminar_range2,range_unit,exam_type,credits,teacher1,teacher2,teacher3,min_number,capacity)
FROM '/docker-entrypoint-initdb.d/courses.csv'
DELIMITER ','
CSV HEADER;

COPY classes(course,class)
FROM '/docker-entrypoint-initdb.d/classes.csv'
DELIMITER ','
CSV HEADER;

COPY classifications(course,classification)
FROM '/docker-entrypoint-initdb.d/classifications.csv'
DELIMITER ','
CSV HEADER;

COPY requisites(course,requisite_type,requisite,from_year,to_year)
FROM '/docker-entrypoint-initdb.d/requisities.csv'
DELIMITER ','
CSV HEADER;

COPY course_texts(course,text_type,lang,title,content,audience)
FROM '/docker-entrypoint-initdb.d/course_texts.csv'
DELIMITER ','
CSV HEADER;

COPY course_teachers(course, teacher)
FROM '/docker-entrypoint-initdb.d/course_teachers.csv'
DELIMITER ','
CSV HEADER;

-- COPY degree_plans(code, plan_year, course, bloc_code, bloc_type, bloc_limit)
COPY degree_plans(code, plan_year, lang, blocs)
FROM '/docker-entrypoint-initdb.d/transformed_degree_plan.csv'
DELIMITER ','
CSV HEADER;

COPY degree_programmes(code, name_cs, name_en, faculty, program_type, program_form, graduate_profile_cs, graduate_profile_en, lang)
FROM '/docker-entrypoint-initdb.d/degree_programmes.csv'
DELIMITER ','
CSV HEADER;

COPY studies(sis_id, student, faculty1, faculty2, study_type, study_form, study_specialization, enrollment, study_state, study_state_date, study_year, degree_plan)
FROM '/docker-entrypoint-initdb.d/studies.csv'
DELIMITER ','
CSV HEADER;

INSERT INTO users (id) VALUES (81411247);

COPY blueprint_years(user_id,academic_year)
FROM '/docker-entrypoint-initdb.d/blueprint_years.csv'
DELIMITER ','
CSV HEADER;

-- COPY blueprint_semesters(blueprint_year,course,semester,position)
-- FROM '/docker-entrypoint-initdb.d/blueprint_semesters.csv'
-- DELIMITER ','
-- CSV HEADER;

INSERT INTO blueprint_semesters(blueprint_year_id, semester)
VALUES
((SELECT id FROM blueprint_years WHERE user_id=81411247 AND academic_year=0), 0),
((SELECT id FROM blueprint_years WHERE user_id=81411247 AND academic_year=1), 1),
((SELECT id FROM blueprint_years WHERE user_id=81411247 AND academic_year=1), 2),
((SELECT id FROM blueprint_years WHERE user_id=81411247 AND academic_year=2), 1),
((SELECT id FROM blueprint_years WHERE user_id=81411247 AND academic_year=2), 2),
((SELECT id FROM blueprint_years WHERE user_id=81411247 AND academic_year=3), 1),
((SELECT id FROM blueprint_years WHERE user_id=81411247 AND academic_year=3), 2);

WITH year_semester AS (
    SELECT s.id AS semester_id, y.academic_year, s.semester FROM blueprint_years y
    LEFT JOIN blueprint_semesters s ON y.id = s.blueprint_year_id
    WHERE y.user_id = 81411247
)

INSERT INTO blueprint_courses(blueprint_semester_id,course,position)
VALUES
((SELECT semester_id FROM year_semester WHERE academic_year=1 AND semester=1),(SELECT id FROM courses WHERE code='NDMI002'),1),
((SELECT semester_id FROM year_semester WHERE academic_year=1 AND semester=1),(SELECT id FROM courses WHERE code='NDMI050'),2),
((SELECT semester_id FROM year_semester WHERE academic_year=1 AND semester=1),(SELECT id FROM courses WHERE code='NJAZ070'),3),
((SELECT semester_id FROM year_semester WHERE academic_year=1 AND semester=1),(SELECT id FROM courses WHERE code='NMAI057'),4),
((SELECT semester_id FROM year_semester WHERE academic_year=1 AND semester=1),(SELECT id FROM courses WHERE code='NMAI069'),5),
((SELECT semester_id FROM year_semester WHERE academic_year=1 AND semester=1),(SELECT id FROM courses WHERE code='NMAT100'),6),
((SELECT semester_id FROM year_semester WHERE academic_year=1 AND semester=1),(SELECT id FROM courses WHERE code='NPRG030'),7),
((SELECT semester_id FROM year_semester WHERE academic_year=1 AND semester=1),(SELECT id FROM courses WHERE code='NPRG062'),8),
((SELECT semester_id FROM year_semester WHERE academic_year=1 AND semester=1),(SELECT id FROM courses WHERE code='NSWI120'),9),
((SELECT semester_id FROM year_semester WHERE academic_year=1 AND semester=1),(SELECT id FROM courses WHERE code='NSWI141'),10),
((SELECT semester_id FROM year_semester WHERE academic_year=1 AND semester=1),(SELECT id FROM courses WHERE code='NTVY006'),11),
((SELECT semester_id FROM year_semester WHERE academic_year=1 AND semester=1),(SELECT id FROM courses WHERE code='NTVY014'),12),
((SELECT semester_id FROM year_semester WHERE academic_year=1 AND semester=2),(SELECT id FROM courses WHERE code='NJAZ072'),1),
((SELECT semester_id FROM year_semester WHERE academic_year=1 AND semester=2),(SELECT id FROM courses WHERE code='NMAI054'),2),
((SELECT semester_id FROM year_semester WHERE academic_year=1 AND semester=2),(SELECT id FROM courses WHERE code='NMAI058'),3),
((SELECT semester_id FROM year_semester WHERE academic_year=1 AND semester=2),(SELECT id FROM courses WHERE code='NPRG031'),4),
((SELECT semester_id FROM year_semester WHERE academic_year=1 AND semester=2),(SELECT id FROM courses WHERE code='NSWI170'),5),
((SELECT semester_id FROM year_semester WHERE academic_year=1 AND semester=2),(SELECT id FROM courses WHERE code='NSWI177'),6),
((SELECT semester_id FROM year_semester WHERE academic_year=1 AND semester=2),(SELECT id FROM courses WHERE code='NTIN060'),7),
((SELECT semester_id FROM year_semester WHERE academic_year=1 AND semester=2),(SELECT id FROM courses WHERE code='NTIN107'),8),
((SELECT semester_id FROM year_semester WHERE academic_year=1 AND semester=2),(SELECT id FROM courses WHERE code='NTVY015'),9),
((SELECT semester_id FROM year_semester WHERE academic_year=2 AND semester=1),(SELECT id FROM courses WHERE code='NAIL062'),1),
((SELECT semester_id FROM year_semester WHERE academic_year=2 AND semester=1),(SELECT id FROM courses WHERE code='NDBI025'),2),
((SELECT semester_id FROM year_semester WHERE academic_year=2 AND semester=1),(SELECT id FROM courses WHERE code='NDMI011'),3),
((SELECT semester_id FROM year_semester WHERE academic_year=2 AND semester=1),(SELECT id FROM courses WHERE code='NJAZ074'),4),
((SELECT semester_id FROM year_semester WHERE academic_year=2 AND semester=1),(SELECT id FROM courses WHERE code='NPRG041'),5),
((SELECT semester_id FROM year_semester WHERE academic_year=2 AND semester=1),(SELECT id FROM courses WHERE code='NSWI142'),6),
((SELECT semester_id FROM year_semester WHERE academic_year=2 AND semester=1),(SELECT id FROM courses WHERE code='NTIN061'),7),
((SELECT semester_id FROM year_semester WHERE academic_year=2 AND semester=1),(SELECT id FROM courses WHERE code='NTVY016'),8),
((SELECT semester_id FROM year_semester WHERE academic_year=2 AND semester=2),(SELECT id FROM courses WHERE code='NJAZ091'),1),
((SELECT semester_id FROM year_semester WHERE academic_year=2 AND semester=2),(SELECT id FROM courses WHERE code='NJAZ176'),2),
((SELECT semester_id FROM year_semester WHERE academic_year=2 AND semester=2),(SELECT id FROM courses WHERE code='NMAI059'),3),
((SELECT semester_id FROM year_semester WHERE academic_year=2 AND semester=2),(SELECT id FROM courses WHERE code='NPRG024'),4),
((SELECT semester_id FROM year_semester WHERE academic_year=2 AND semester=2),(SELECT id FROM courses WHERE code='NPRG036'),5),
((SELECT semester_id FROM year_semester WHERE academic_year=2 AND semester=2),(SELECT id FROM courses WHERE code='NPRG045'),6),
((SELECT semester_id FROM year_semester WHERE academic_year=2 AND semester=2),(SELECT id FROM courses WHERE code='NPRG051'),7),
((SELECT semester_id FROM year_semester WHERE academic_year=2 AND semester=2),(SELECT id FROM courses WHERE code='NSWI143'),8),
((SELECT semester_id FROM year_semester WHERE academic_year=2 AND semester=2),(SELECT id FROM courses WHERE code='NSWI153'),9),
((SELECT semester_id FROM year_semester WHERE academic_year=2 AND semester=2),(SELECT id FROM courses WHERE code='NTIN071'),10),
((SELECT semester_id FROM year_semester WHERE academic_year=2 AND semester=2),(SELECT id FROM courses WHERE code='NTVY017'),11),
((SELECT semester_id FROM year_semester WHERE academic_year=3 AND semester=1),(SELECT id FROM courses WHERE code='NPFL129'),1),
((SELECT semester_id FROM year_semester WHERE academic_year=3 AND semester=1),(SELECT id FROM courses WHERE code='NPGR003'),2),
((SELECT semester_id FROM year_semester WHERE academic_year=3 AND semester=1),(SELECT id FROM courses WHERE code='NPRG035'),3),
((SELECT semester_id FROM year_semester WHERE academic_year=3 AND semester=1),(SELECT id FROM courses WHERE code='NPRG073'),4),
((SELECT semester_id FROM year_semester WHERE academic_year=3 AND semester=1),(SELECT id FROM courses WHERE code='NSWI004'),5),
((SELECT semester_id FROM year_semester WHERE academic_year=3 AND semester=1),(SELECT id FROM courses WHERE code='NSWI098'),6),
((SELECT semester_id FROM year_semester WHERE academic_year=3 AND semester=1),(SELECT id FROM courses WHERE code='NSWI154'),7),
((SELECT semester_id FROM year_semester WHERE academic_year=3 AND semester=2),(SELECT id FROM courses WHERE code='NPRG038'),1),
((SELECT semester_id FROM year_semester WHERE academic_year=3 AND semester=2),(SELECT id FROM courses WHERE code='NPRG043'),2),
((SELECT semester_id FROM year_semester WHERE academic_year=3 AND semester=2),(SELECT id FROM courses WHERE code='NPRG074'),3),
((SELECT semester_id FROM year_semester WHERE academic_year=3 AND semester=2),(SELECT id FROM courses WHERE code='NSWI041'),4),
((SELECT semester_id FROM year_semester WHERE academic_year=3 AND semester=2),(SELECT id FROM courses WHERE code='NSZZ031'),5);

INSERT INTO bla_studies(user_id, degree_plan_code, start_year)
VALUES
    (81411247, 'NIPVS19B', 2020),
    (81411247, 'NISD23N', 2023);

INSERT INTO sessions(id, user_id, expires_at)
VALUES ('977e69df-0b48-4790-a409-b86656ff86bc', 81411247, '2200-01-01 00:00:00-00'::timestamptz);

INSERT INTO start_semester_to_desc(id, lang, semester_description) VALUES
    (1, 'cs', 'Zimní'),
    (2, 'cs', 'Letní'),
    (3, 'cs', 'Oba'),
    (1, 'en', 'Winter'),
    (2, 'en', 'Summer'),
    (3, 'en', 'Both');


COPY bla_courses(title,code,valid_from,valid_to,guarantor,taught,start_semester,semester_count,taught_lang,lecture_range1,seminar_range1,lecture_range2,seminar_range2,range_unit,exam_type,credits,min_number,capacity,lang,guarantors,annotation,aim,requirements,syllabus,teachers,faculty)
FROM '/docker-entrypoint-initdb.d/courses_transformed.csv'
DELIMITER ','
CSV HEADER;

INSERT INTO bla_blueprints(user_id, lang, blueprint) VALUES
    (81411247, 'cs' , '{"unassigned": [], "assigned": [{"year": 1, "winter": [], "summer": []},{"year": 2, "winter": [], "summer": []},{"year": 3, "winter": [], "summer": []}]}'),
    (81411247, 'en' , '{"unassigned": [], "assigned": [{"year": 1, "winter": [], "summer": []},{"year": 2, "winter": [], "summer": []},{"year": 3, "winter": [], "summer": []}]}');

WITH blueprint_record AS (
    SELECT code, title, valid_from, start_semester, semester_count, lecture_range1, seminar_range1, lecture_range2, seminar_range2, exam_type, credits, guarantors
    FROM bla_courses
    WHERE lang='cs'
)
UPDATE bla_blueprints
SET blueprint = jsonb_build_object('unassigned', '[]', 'assigned', jsonb_build_array(
    (SELECT * from jsonb_build_object('year', 1,
        'winter', (SELECT * FROM jsonb_build_array(
            (SELECT row_to_json(r) FROM (SELECT * FROM blueprint_record WHERE code = 'NDMI002') r),
            (SELECT row_to_json(r) FROM (SELECT * FROM blueprint_record WHERE code = 'NDMI050') r),
            (SELECT row_to_json(r) FROM (SELECT * FROM blueprint_record WHERE code = 'NJAZ070') r),
            (SELECT row_to_json(r) FROM (SELECT * FROM blueprint_record WHERE code = 'NMAI057') r),
            (SELECT row_to_json(r) FROM (SELECT * FROM blueprint_record WHERE code = 'NMAI069') r),
            (SELECT row_to_json(r) FROM (SELECT * FROM blueprint_record WHERE code = 'NMAT100') r),
            (SELECT row_to_json(r) FROM (SELECT * FROM blueprint_record WHERE code = 'NPRG030') r),
            (SELECT row_to_json(r) FROM (SELECT * FROM blueprint_record WHERE code = 'NPRG062') r),
            (SELECT row_to_json(r) FROM (SELECT * FROM blueprint_record WHERE code = 'NSWI120') r),
            (SELECT row_to_json(r) FROM (SELECT * FROM blueprint_record WHERE code = 'NSWI141') r),
            (SELECT row_to_json(r) FROM (SELECT * FROM blueprint_record WHERE code = 'NTVY006') r),
            (SELECT row_to_json(r) FROM (SELECT * FROM blueprint_record WHERE code = 'NTVY014') r)
        )),
        'summer', (SELECT * FROM jsonb_build_array(
            (SELECT row_to_json(r) FROM (SELECT * FROM blueprint_record WHERE code = 'NJAZ072') r),
            (SELECT row_to_json(r) FROM (SELECT * FROM blueprint_record WHERE code = 'NMAI054') r),
            (SELECT row_to_json(r) FROM (SELECT * FROM blueprint_record WHERE code = 'NMAI058') r),
            (SELECT row_to_json(r) FROM (SELECT * FROM blueprint_record WHERE code = 'NPRG031') r),
            (SELECT row_to_json(r) FROM (SELECT * FROM blueprint_record WHERE code = 'NSWI170') r),
            (SELECT row_to_json(r) FROM (SELECT * FROM blueprint_record WHERE code = 'NSWI177') r),
            (SELECT row_to_json(r) FROM (SELECT * FROM blueprint_record WHERE code = 'NTIN060') r),
            (SELECT row_to_json(r) FROM (SELECT * FROM blueprint_record WHERE code = 'NTIN107') r),
            (SELECT row_to_json(r) FROM (SELECT * FROM blueprint_record WHERE code = 'NTVY015') r)
        ))
    )),
    (SELECT * from jsonb_build_object('year', 2,
        'winter', (SELECT * FROM jsonb_build_array(
            (SELECT row_to_json(r) FROM (SELECT * FROM blueprint_record WHERE code = 'NAIL062') r),
            (SELECT row_to_json(r) FROM (SELECT * FROM blueprint_record WHERE code = 'NDBI025') r),
            (SELECT row_to_json(r) FROM (SELECT * FROM blueprint_record WHERE code = 'NDMI011') r),
            (SELECT row_to_json(r) FROM (SELECT * FROM blueprint_record WHERE code = 'NJAZ074') r),
            (SELECT row_to_json(r) FROM (SELECT * FROM blueprint_record WHERE code = 'NPRG041') r),
            (SELECT row_to_json(r) FROM (SELECT * FROM blueprint_record WHERE code = 'NSWI142') r),
            (SELECT row_to_json(r) FROM (SELECT * FROM blueprint_record WHERE code = 'NTIN061') r),
            (SELECT row_to_json(r) FROM (SELECT * FROM blueprint_record WHERE code = 'NTVY016') r)
        )),
        'summer', (SELECT * FROM jsonb_build_array(
            (SELECT row_to_json(r) FROM (SELECT * FROM blueprint_record WHERE code = 'NJAZ091') r),
            (SELECT row_to_json(r) FROM (SELECT * FROM blueprint_record WHERE code = 'NJAZ176') r),
            (SELECT row_to_json(r) FROM (SELECT * FROM blueprint_record WHERE code = 'NMAI059') r),
            (SELECT row_to_json(r) FROM (SELECT * FROM blueprint_record WHERE code = 'NPRG024') r),
            (SELECT row_to_json(r) FROM (SELECT * FROM blueprint_record WHERE code = 'NPRG036') r),
            (SELECT row_to_json(r) FROM (SELECT * FROM blueprint_record WHERE code = 'NPRG045') r),
            (SELECT row_to_json(r) FROM (SELECT * FROM blueprint_record WHERE code = 'NPRG051') r),
            (SELECT row_to_json(r) FROM (SELECT * FROM blueprint_record WHERE code = 'NSWI143') r),
            (SELECT row_to_json(r) FROM (SELECT * FROM blueprint_record WHERE code = 'NSWI153') r),
            (SELECT row_to_json(r) FROM (SELECT * FROM blueprint_record WHERE code = 'NTIN071') r),
            (SELECT row_to_json(r) FROM (SELECT * FROM blueprint_record WHERE code = 'NTVY017') r)
        ))
    )),
    (SELECT * from jsonb_build_object('year', 3,
        'winter', (SELECT * FROM jsonb_build_array(
            (SELECT row_to_json(r) FROM (SELECT * FROM blueprint_record WHERE code = 'NPFL129') r),
            (SELECT row_to_json(r) FROM (SELECT * FROM blueprint_record WHERE code = 'NPGR003') r),
            (SELECT row_to_json(r) FROM (SELECT * FROM blueprint_record WHERE code = 'NPRG035') r),
            (SELECT row_to_json(r) FROM (SELECT * FROM blueprint_record WHERE code = 'NPRG073') r),
            (SELECT row_to_json(r) FROM (SELECT * FROM blueprint_record WHERE code = 'NSWI004') r),
            (SELECT row_to_json(r) FROM (SELECT * FROM blueprint_record WHERE code = 'NSWI098') r),
            (SELECT row_to_json(r) FROM (SELECT * FROM blueprint_record WHERE code = 'NSWI154') r),
            (SELECT row_to_json(r) FROM (SELECT * FROM blueprint_record WHERE code = 'NPRG038') r)
        )),
        'summer', (SELECT * FROM jsonb_build_array(
            (SELECT row_to_json(r) FROM (SELECT * FROM blueprint_record WHERE code = 'NPRG043') r),
            (SELECT row_to_json(r) FROM (SELECT * FROM blueprint_record WHERE code = 'NPRG074') r),
            (SELECT row_to_json(r) FROM (SELECT * FROM blueprint_record WHERE code = 'NSWI041') r),
            (SELECT row_to_json(r) FROM (SELECT * FROM blueprint_record WHERE code = 'NSZZ031') r)
        ))
    ))
))
WHERE user_id = 81411247
AND lang = 'cs';

WITH blueprint_record AS (
    SELECT code, title, valid_from, start_semester, lecture_range1, seminar_range1, lecture_range2, seminar_range2, exam_type, credits, guarantors
    FROM bla_courses
    WHERE lang='en'
)
UPDATE bla_blueprints
SET blueprint = jsonb_build_object('unassigned', '[]', 'assigned', jsonb_build_array(
    (SELECT * from jsonb_build_object('year', 1,
        'winter', (SELECT * FROM jsonb_build_array(
            (SELECT row_to_json(r) FROM (SELECT * FROM blueprint_record WHERE code = 'NDMI002') r),
            (SELECT row_to_json(r) FROM (SELECT * FROM blueprint_record WHERE code = 'NDMI050') r),
            (SELECT row_to_json(r) FROM (SELECT * FROM blueprint_record WHERE code = 'NJAZ070') r),
            (SELECT row_to_json(r) FROM (SELECT * FROM blueprint_record WHERE code = 'NMAI057') r),
            (SELECT row_to_json(r) FROM (SELECT * FROM blueprint_record WHERE code = 'NMAI069') r),
            (SELECT row_to_json(r) FROM (SELECT * FROM blueprint_record WHERE code = 'NMAT100') r),
            (SELECT row_to_json(r) FROM (SELECT * FROM blueprint_record WHERE code = 'NPRG030') r),
            (SELECT row_to_json(r) FROM (SELECT * FROM blueprint_record WHERE code = 'NPRG062') r),
            (SELECT row_to_json(r) FROM (SELECT * FROM blueprint_record WHERE code = 'NSWI120') r),
            (SELECT row_to_json(r) FROM (SELECT * FROM blueprint_record WHERE code = 'NSWI141') r),
            (SELECT row_to_json(r) FROM (SELECT * FROM blueprint_record WHERE code = 'NTVY006') r),
            (SELECT row_to_json(r) FROM (SELECT * FROM blueprint_record WHERE code = 'NTVY014') r)
        )),
        'summer', (SELECT * FROM jsonb_build_array(
            (SELECT row_to_json(r) FROM (SELECT * FROM blueprint_record WHERE code = 'NJAZ072') r),
            (SELECT row_to_json(r) FROM (SELECT * FROM blueprint_record WHERE code = 'NMAI054') r),
            (SELECT row_to_json(r) FROM (SELECT * FROM blueprint_record WHERE code = 'NMAI058') r),
            (SELECT row_to_json(r) FROM (SELECT * FROM blueprint_record WHERE code = 'NPRG031') r),
            (SELECT row_to_json(r) FROM (SELECT * FROM blueprint_record WHERE code = 'NSWI170') r),
            (SELECT row_to_json(r) FROM (SELECT * FROM blueprint_record WHERE code = 'NSWI177') r),
            (SELECT row_to_json(r) FROM (SELECT * FROM blueprint_record WHERE code = 'NTIN060') r),
            (SELECT row_to_json(r) FROM (SELECT * FROM blueprint_record WHERE code = 'NTIN107') r),
            (SELECT row_to_json(r) FROM (SELECT * FROM blueprint_record WHERE code = 'NTVY015') r)
        ))
    )),
    (SELECT * from jsonb_build_object('year', 2,
        'winter', (SELECT * FROM jsonb_build_array(
            (SELECT row_to_json(r) FROM (SELECT * FROM blueprint_record WHERE code = 'NAIL062') r),
            (SELECT row_to_json(r) FROM (SELECT * FROM blueprint_record WHERE code = 'NDBI025') r),
            (SELECT row_to_json(r) FROM (SELECT * FROM blueprint_record WHERE code = 'NDMI011') r),
            (SELECT row_to_json(r) FROM (SELECT * FROM blueprint_record WHERE code = 'NJAZ074') r),
            (SELECT row_to_json(r) FROM (SELECT * FROM blueprint_record WHERE code = 'NPRG041') r),
            (SELECT row_to_json(r) FROM (SELECT * FROM blueprint_record WHERE code = 'NSWI142') r),
            (SELECT row_to_json(r) FROM (SELECT * FROM blueprint_record WHERE code = 'NTIN061') r),
            (SELECT row_to_json(r) FROM (SELECT * FROM blueprint_record WHERE code = 'NTVY016') r)
        )),
        'summer', (SELECT * FROM jsonb_build_array(
            (SELECT row_to_json(r) FROM (SELECT * FROM blueprint_record WHERE code = 'NJAZ091') r),
            (SELECT row_to_json(r) FROM (SELECT * FROM blueprint_record WHERE code = 'NJAZ176') r),
            (SELECT row_to_json(r) FROM (SELECT * FROM blueprint_record WHERE code = 'NMAI059') r),
            (SELECT row_to_json(r) FROM (SELECT * FROM blueprint_record WHERE code = 'NPRG024') r),
            (SELECT row_to_json(r) FROM (SELECT * FROM blueprint_record WHERE code = 'NPRG036') r),
            (SELECT row_to_json(r) FROM (SELECT * FROM blueprint_record WHERE code = 'NPRG045') r),
            (SELECT row_to_json(r) FROM (SELECT * FROM blueprint_record WHERE code = 'NPRG051') r),
            (SELECT row_to_json(r) FROM (SELECT * FROM blueprint_record WHERE code = 'NSWI143') r),
            (SELECT row_to_json(r) FROM (SELECT * FROM blueprint_record WHERE code = 'NSWI153') r),
            (SELECT row_to_json(r) FROM (SELECT * FROM blueprint_record WHERE code = 'NTIN071') r),
            (SELECT row_to_json(r) FROM (SELECT * FROM blueprint_record WHERE code = 'NTVY017') r)
        ))
    )),
    (SELECT * from jsonb_build_object('year', 3,
        'winter', (SELECT * FROM jsonb_build_array(
            (SELECT row_to_json(r) FROM (SELECT * FROM blueprint_record WHERE code = 'NPFL129') r),
            (SELECT row_to_json(r) FROM (SELECT * FROM blueprint_record WHERE code = 'NPGR003') r),
            (SELECT row_to_json(r) FROM (SELECT * FROM blueprint_record WHERE code = 'NPRG035') r),
            (SELECT row_to_json(r) FROM (SELECT * FROM blueprint_record WHERE code = 'NPRG073') r),
            (SELECT row_to_json(r) FROM (SELECT * FROM blueprint_record WHERE code = 'NSWI004') r),
            (SELECT row_to_json(r) FROM (SELECT * FROM blueprint_record WHERE code = 'NSWI098') r),
            (SELECT row_to_json(r) FROM (SELECT * FROM blueprint_record WHERE code = 'NSWI154') r),
            (SELECT row_to_json(r) FROM (SELECT * FROM blueprint_record WHERE code = 'NPRG038') r)
        )),
        'summer', (SELECT * FROM jsonb_build_array(
            (SELECT row_to_json(r) FROM (SELECT * FROM blueprint_record WHERE code = 'NPRG043') r),
            (SELECT row_to_json(r) FROM (SELECT * FROM blueprint_record WHERE code = 'NPRG074') r),
            (SELECT row_to_json(r) FROM (SELECT * FROM blueprint_record WHERE code = 'NSWI041') r),
            (SELECT row_to_json(r) FROM (SELECT * FROM blueprint_record WHERE code = 'NSZZ031') r)
        ))
    ))
))
WHERE user_id = 81411247
AND lang = 'en';