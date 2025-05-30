COPY faculties(sis_id, sis_poid, name_cs, name_en, abbr)
FROM '/docker-entrypoint-initdb.d/faculties.csv'
DELIMITER ','
CSV HEADER;

-- COPY teachers(sis_id,department,faculty,last_name,first_name,title_before,title_after)
-- FROM '/docker-entrypoint-initdb.d/UCIT.csv'
-- DELIMITER ','
-- CSV HEADER;

-- COPY old_courses(code,name_cs,name_en,valid_from,valid_to,faculty,guarantor,taught,start_semester,semester_count,taught_lang,lecture_range1,seminar_range1,lecture_range2,seminar_range2,range_unit,exam_type,credits,teacher1,teacher2,teacher3,min_number,capacity)
-- FROM '/docker-entrypoint-initdb.d/POVINN.csv'
-- DELIMITER ','
-- CSV HEADER;

-- COPY classes(course,class)
-- FROM '/docker-entrypoint-initdb.d/classes.csv'
-- DELIMITER ','
-- CSV HEADER;

-- COPY classifications(course,classification)
-- FROM '/docker-entrypoint-initdb.d/classifications.csv'
-- DELIMITER ','
-- CSV HEADER;

-- COPY requisites(course,requisite_type,requisite,from_year,to_year)
-- FROM '/docker-entrypoint-initdb.d/requisities.csv'
-- DELIMITER ','
-- CSV HEADER;

-- COPY course_texts(course,text_type,lang,title,content,audience)
-- FROM '/docker-entrypoint-initdb.d/course_texts.csv'
-- DELIMITER ','
-- CSV HEADER;

-- COPY course_teachers(course, teacher)
-- FROM '/docker-entrypoint-initdb.d/UCIT_ROZVRH.csv'
-- DELIMITER ','
-- CSV HEADER;

-- COPY degree_plans(code, plan_year, course, bloc_code, bloc_type, bloc_limit)
-- COPY degree_plans(code, plan_year, lang, blocs)
COPY degree_plans(plan_code, plan_year, course_code, interchangeability, bloc_subject_code, bloc_type, bloc_limit, seq, bloc_name, bloc_note, note, lang)
FROM '/docker-entrypoint-initdb.d/degree_plans_transformed.csv'
DELIMITER ','
CSV HEADER;

-- COPY degree_programmes(code, name_cs, name_en, faculty, program_type, program_form, graduate_profile_cs, graduate_profile_en, lang)
-- FROM '/docker-entrypoint-initdb.d/degree_programmes.csv'
-- DELIMITER ','
-- CSV HEADER;

-- COPY studies(sis_id, student, faculty1, faculty2, study_type, study_form, study_specialization, enrollment, study_state, study_state_date, study_year, degree_plan)
-- FROM '/docker-entrypoint-initdb.d/studies.csv'
-- DELIMITER ','
-- CSV HEADER;

COPY courses(title,code,valid_from,valid_to,guarantor,taught,start_semester,semester_count,range_unit,exam_type,credits,min_number,capacity,lang,taught_lang,guarantors,annotation,aim,terms_of_passing,literature,requirements_for_assesment,syllabus,entry_requirements,teachers,faculty,comments,preqrequisities,corequisities,incompatibilities,interchangebilities,classes,classifications,lecture_range1,seminar_range1,lecture_range2,seminar_range2)
FROM '/docker-entrypoint-initdb.d/courses_transformed.csv'
DELIMITER ','
CSV HEADER;

INSERT INTO users (id) VALUES ('81411247'), ('73291111');

COPY blueprint_years(user_id,academic_year)
FROM '/docker-entrypoint-initdb.d/blueprint_years.csv'
DELIMITER ','
CSV HEADER;

INSERT INTO filters (id)
VALUES
    ('courses'),
    ('course-survey');

COPY filter_categories(id, facet_id, title_cs, title_en, description_cs, description_en, displayed_value_limit, position, filter_id)
FROM '/docker-entrypoint-initdb.d/filter_categories.csv'
DELIMITER ','
CSV HEADER;

COPY filter_values(id, title_cs, title_en, category_id, description_cs, description_en, facet_id, position)
FROM '/docker-entrypoint-initdb.d/filter_values.csv'
DELIMITER ','
CSV HEADER;

-- COPY filter_labels(id, lang, label)
-- FROM '/docker-entrypoint-initdb.d/filter_values.csv'
-- DELIMITER ','
-- CSV HEADER;

-- COPY filter_params(param_name, value_id)
-- FROM '/docker-entrypoint-initdb.d/filter_params.csv'
-- DELIMITER ',';

-- COPY blueprint_semesters(blueprint_year,course,semester,position)
-- FROM '/docker-entrypoint-initdb.d/blueprint_semesters.csv'
-- DELIMITER ','
-- CSV HEADER;

INSERT INTO blueprint_semesters(blueprint_year_id, semester)
VALUES
((SELECT id FROM blueprint_years WHERE user_id='81411247' AND academic_year=0), 0),
((SELECT id FROM blueprint_years WHERE user_id='81411247' AND academic_year=1), 1),
((SELECT id FROM blueprint_years WHERE user_id='81411247' AND academic_year=1), 2),
((SELECT id FROM blueprint_years WHERE user_id='81411247' AND academic_year=2), 1),
((SELECT id FROM blueprint_years WHERE user_id='81411247' AND academic_year=2), 2),
((SELECT id FROM blueprint_years WHERE user_id='81411247' AND academic_year=3), 1),
((SELECT id FROM blueprint_years WHERE user_id='81411247' AND academic_year=3), 2);

WITH year_semester AS (
    SELECT s.id AS semester_id, y.academic_year, s.semester FROM blueprint_years y
    LEFT JOIN blueprint_semesters s ON y.id = s.blueprint_year_id
    WHERE y.user_id = '81411247'
)

INSERT INTO blueprint_courses(blueprint_semester_id,course_code,course_valid_from,position)
VALUES
((SELECT semester_id FROM year_semester WHERE academic_year=1 AND semester=1),'NDMI002',(SELECT MAX(valid_from) FROM courses WHERE code='NDMI002'),1),
((SELECT semester_id FROM year_semester WHERE academic_year=1 AND semester=1),'NDMI050',(SELECT MAX(valid_from) FROM courses WHERE code='NDMI050'),2),
((SELECT semester_id FROM year_semester WHERE academic_year=1 AND semester=1),'NJAZ070',(SELECT MAX(valid_from) FROM courses WHERE code='NJAZ070'),3),
((SELECT semester_id FROM year_semester WHERE academic_year=1 AND semester=1),'NMAI057',(SELECT MAX(valid_from) FROM courses WHERE code='NMAI057'),4),
((SELECT semester_id FROM year_semester WHERE academic_year=1 AND semester=1),'NMAI069',(SELECT MAX(valid_from) FROM courses WHERE code='NMAI069'),5),
((SELECT semester_id FROM year_semester WHERE academic_year=1 AND semester=1),'NMAT100',(SELECT MAX(valid_from) FROM courses WHERE code='NMAT100'),6),
((SELECT semester_id FROM year_semester WHERE academic_year=1 AND semester=1),'NPRG030',(SELECT MAX(valid_from) FROM courses WHERE code='NPRG030'),7),
((SELECT semester_id FROM year_semester WHERE academic_year=1 AND semester=1),'NPRG062',(SELECT MAX(valid_from) FROM courses WHERE code='NPRG062'),8),
((SELECT semester_id FROM year_semester WHERE academic_year=1 AND semester=1),'NSWI120',(SELECT MAX(valid_from) FROM courses WHERE code='NSWI120'),9),
((SELECT semester_id FROM year_semester WHERE academic_year=1 AND semester=1),'NSWI141',(SELECT MAX(valid_from) FROM courses WHERE code='NSWI141'),10),
((SELECT semester_id FROM year_semester WHERE academic_year=1 AND semester=1),'NTVY006',(SELECT MAX(valid_from) FROM courses WHERE code='NTVY006'),11),
((SELECT semester_id FROM year_semester WHERE academic_year=1 AND semester=1),'NTVY014',(SELECT MAX(valid_from) FROM courses WHERE code='NTVY014'),12),
((SELECT semester_id FROM year_semester WHERE academic_year=1 AND semester=2),'NJAZ072',(SELECT MAX(valid_from) FROM courses WHERE code='NJAZ072'),1),
((SELECT semester_id FROM year_semester WHERE academic_year=1 AND semester=2),'NMAI054',(SELECT MAX(valid_from) FROM courses WHERE code='NMAI054'),2),
((SELECT semester_id FROM year_semester WHERE academic_year=1 AND semester=2),'NMAI058',(SELECT MAX(valid_from) FROM courses WHERE code='NMAI058'),3),
((SELECT semester_id FROM year_semester WHERE academic_year=1 AND semester=2),'NPRG031',(SELECT MAX(valid_from) FROM courses WHERE code='NPRG031'),4),
((SELECT semester_id FROM year_semester WHERE academic_year=1 AND semester=2),'NSWI170',(SELECT MAX(valid_from) FROM courses WHERE code='NSWI170'),5),
((SELECT semester_id FROM year_semester WHERE academic_year=1 AND semester=2),'NSWI177',(SELECT MAX(valid_from) FROM courses WHERE code='NSWI177'),6),
((SELECT semester_id FROM year_semester WHERE academic_year=1 AND semester=2),'NTIN060',(SELECT MAX(valid_from) FROM courses WHERE code='NTIN060'),7),
((SELECT semester_id FROM year_semester WHERE academic_year=1 AND semester=2),'NTIN107',(SELECT MAX(valid_from) FROM courses WHERE code='NTIN107'),8),
((SELECT semester_id FROM year_semester WHERE academic_year=1 AND semester=2),'NTVY015',(SELECT MAX(valid_from) FROM courses WHERE code='NTVY015'),9),
((SELECT semester_id FROM year_semester WHERE academic_year=2 AND semester=1),'NAIL062',(SELECT MAX(valid_from) FROM courses WHERE code='NAIL062'),1),
((SELECT semester_id FROM year_semester WHERE academic_year=2 AND semester=1),'NDBI025',(SELECT MAX(valid_from) FROM courses WHERE code='NDBI025'),2),
((SELECT semester_id FROM year_semester WHERE academic_year=2 AND semester=1),'NDMI011',(SELECT MAX(valid_from) FROM courses WHERE code='NDMI011'),3),
((SELECT semester_id FROM year_semester WHERE academic_year=2 AND semester=1),'NJAZ074',(SELECT MAX(valid_from) FROM courses WHERE code='NJAZ074'),4),
((SELECT semester_id FROM year_semester WHERE academic_year=2 AND semester=1),'NPRG041',(SELECT MAX(valid_from) FROM courses WHERE code='NPRG041'),5),
((SELECT semester_id FROM year_semester WHERE academic_year=2 AND semester=1),'NSWI142',(SELECT MAX(valid_from) FROM courses WHERE code='NSWI142'),6),
((SELECT semester_id FROM year_semester WHERE academic_year=2 AND semester=1),'NTIN061',(SELECT MAX(valid_from) FROM courses WHERE code='NTIN061'),7),
((SELECT semester_id FROM year_semester WHERE academic_year=2 AND semester=1),'NTVY016',(SELECT MAX(valid_from) FROM courses WHERE code='NTVY016'),8),
((SELECT semester_id FROM year_semester WHERE academic_year=2 AND semester=2),'NJAZ091',(SELECT MAX(valid_from) FROM courses WHERE code='NJAZ091'),1),
((SELECT semester_id FROM year_semester WHERE academic_year=2 AND semester=2),'NJAZ176',(SELECT MAX(valid_from) FROM courses WHERE code='NJAZ176'),2),
((SELECT semester_id FROM year_semester WHERE academic_year=2 AND semester=2),'NMAI059',(SELECT MAX(valid_from) FROM courses WHERE code='NMAI059'),3),
((SELECT semester_id FROM year_semester WHERE academic_year=2 AND semester=2),'NPRG024',(SELECT MAX(valid_from) FROM courses WHERE code='NPRG024'),4),
((SELECT semester_id FROM year_semester WHERE academic_year=2 AND semester=2),'NPRG036',(SELECT MAX(valid_from) FROM courses WHERE code='NPRG036'),5),
((SELECT semester_id FROM year_semester WHERE academic_year=2 AND semester=2),'NPRG045',(SELECT MAX(valid_from) FROM courses WHERE code='NPRG045'),6),
((SELECT semester_id FROM year_semester WHERE academic_year=2 AND semester=2),'NPRG051',(SELECT MAX(valid_from) FROM courses WHERE code='NPRG051'),7),
((SELECT semester_id FROM year_semester WHERE academic_year=2 AND semester=2),'NSWI143',(SELECT MAX(valid_from) FROM courses WHERE code='NSWI143'),8),
((SELECT semester_id FROM year_semester WHERE academic_year=2 AND semester=2),'NSWI153',(SELECT MAX(valid_from) FROM courses WHERE code='NSWI153'),9),
((SELECT semester_id FROM year_semester WHERE academic_year=2 AND semester=2),'NTIN071',(SELECT MAX(valid_from) FROM courses WHERE code='NTIN071'),10),
((SELECT semester_id FROM year_semester WHERE academic_year=2 AND semester=2),'NTVY017',(SELECT MAX(valid_from) FROM courses WHERE code='NTVY017'),11),
((SELECT semester_id FROM year_semester WHERE academic_year=3 AND semester=1),'NPFL129',(SELECT MAX(valid_from) FROM courses WHERE code='NPFL129'),1),
((SELECT semester_id FROM year_semester WHERE academic_year=3 AND semester=1),'NPGR003',(SELECT MAX(valid_from) FROM courses WHERE code='NPGR003'),2),
((SELECT semester_id FROM year_semester WHERE academic_year=3 AND semester=1),'NPRG035',(SELECT MAX(valid_from) FROM courses WHERE code='NPRG035'),3),
((SELECT semester_id FROM year_semester WHERE academic_year=3 AND semester=1),'NPRG073',(SELECT MAX(valid_from) FROM courses WHERE code='NPRG073'),4),
((SELECT semester_id FROM year_semester WHERE academic_year=3 AND semester=1),'NSWI004',(SELECT MAX(valid_from) FROM courses WHERE code='NSWI004'),5),
((SELECT semester_id FROM year_semester WHERE academic_year=3 AND semester=1),'NSWI098',(SELECT MAX(valid_from) FROM courses WHERE code='NSWI098'),6),
((SELECT semester_id FROM year_semester WHERE academic_year=3 AND semester=1),'NSWI154',(SELECT MAX(valid_from) FROM courses WHERE code='NSWI154'),7),
((SELECT semester_id FROM year_semester WHERE academic_year=3 AND semester=2),'NPRG038',(SELECT MAX(valid_from) FROM courses WHERE code='NPRG038'),1),
((SELECT semester_id FROM year_semester WHERE academic_year=3 AND semester=2),'NPRG043',(SELECT MAX(valid_from) FROM courses WHERE code='NPRG043'),2),
((SELECT semester_id FROM year_semester WHERE academic_year=3 AND semester=2),'NPRG074',(SELECT MAX(valid_from) FROM courses WHERE code='NPRG074'),3),
((SELECT semester_id FROM year_semester WHERE academic_year=3 AND semester=2),'NSWI041',(SELECT MAX(valid_from) FROM courses WHERE code='NSWI041'),4),
((SELECT semester_id FROM year_semester WHERE academic_year=3 AND semester=2),'NSZZ031',(SELECT MAX(valid_from) FROM courses WHERE code='NSZZ031'),5);

INSERT INTO bla_studies(user_id, degree_plan_code, start_year)
VALUES
    ('81411247', 'NIPVS19B', 2020),
    ('81411247', 'NISD23N', 2023),
    ('73291111', 'NIPVS19B', 2020);

-- INSERT INTO sessions(id, user_id, expires_at)
-- VALUES
--     ('977e69df-0b48-4790-a409-b86656ff86bc', '81411247', '2200-01-01 00:00:00-00'::timestamptz),
--     ('17924537-c555-4fc7-b83c-f1e839cfee61', '73291111', '2200-01-01 00:00:00-00'::timestamptz);

INSERT INTO start_semester_to_desc(id, lang, semester_description) VALUES
    (1, 'cs', 'Zimní'),
    (2, 'cs', 'Letní'),
    (3, 'cs', 'Oba'),
    (1, 'en', 'Winter'),
    (2, 'en', 'Summer'),
    (3, 'en', 'Both');

INSERT INTO course_overall_ratings(user_id, course_code, rating)
VALUES
    ('81411247', 'NDMI002', 1),
    ('81411247', 'NDMI050', 0),
    ('73291111', 'NDMI002', 1),
    ('73291111', 'NDMI050', 0);

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

INSERT INTO course_ratings(user_id, course_code, category_code, rating)
VALUES
    ('81411247', 'NDMI002', 1, 1),
    ('81411247', 'NDMI002', 2, 4),
    ('81411247', 'NDMI002', 3, 3),
    ('81411247', 'NDMI002', 4, 2),
    ('81411247', 'NDMI002', 5, 1),
    ('81411247', 'NDMI050', 1, 0),
    ('81411247', 'NDMI050', 2, 1),
    ('81411247', 'NDMI050', 3, 5),
    ('81411247', 'NDMI050', 4, 4),
    ('81411247', 'NDMI050', 5, 3),
    ('73291111', 'NDMI002', 1, 1),
    ('73291111', 'NDMI002', 2, 4),
    ('73291111', 'NDMI002', 3, 3),
    ('73291111', 'NDMI002', 4, 2),
    ('73291111', 'NDMI002', 5, 1),
    ('73291111', 'NDMI050', 1, 0),
    ('73291111', 'NDMI050', 2, 1),
    ('73291111', 'NDMI050', 3, 5),
    ('73291111', 'NDMI050', 4, 4),
    ('73291111', 'NDMI050', 5, 3);