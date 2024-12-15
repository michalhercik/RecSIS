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

COPY degree_plans(code, plan_year, course, bloc_code, bloc_type, bloc_limit)
FROM '/docker-entrypoint-initdb.d/degree_plans.csv'
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

COPY blueprint_years(student,position)
FROM '/docker-entrypoint-initdb.d/blueprint_years.csv'
DELIMITER ','
CSV HEADER;

-- COPY blueprint_semesters(blueprint_year,course,semester,position)
-- FROM '/docker-entrypoint-initdb.d/blueprint_semesters.csv'
-- DELIMITER ','
-- CSV HEADER;

INSERT INTO blueprint_semesters(blueprint_year,course,semester,position)
VALUES
(1,(SELECT id FROM courses WHERE code='NDMI002'),1,1),
(1,(SELECT id FROM courses WHERE code='NDMI050'),1,2),
(1,(SELECT id FROM courses WHERE code='NJAZ070'),1,3),
(1,(SELECT id FROM courses WHERE code='NMAI057'),1,4),
(1,(SELECT id FROM courses WHERE code='NMAI069'),1,5),
(1,(SELECT id FROM courses WHERE code='NMAT100'),1,6),
(1,(SELECT id FROM courses WHERE code='NPRG030'),1,7),
(1,(SELECT id FROM courses WHERE code='NPRG062'),1,8),
(1,(SELECT id FROM courses WHERE code='NSWI120'),1,9),
(1,(SELECT id FROM courses WHERE code='NSWI141'),1,10),
(1,(SELECT id FROM courses WHERE code='NTVY006'),1,11),
(1,(SELECT id FROM courses WHERE code='NTVY014'),1,12),
(1,(SELECT id FROM courses WHERE code='NJAZ072'),2,1),
(1,(SELECT id FROM courses WHERE code='NMAI054'),2,2),
(1,(SELECT id FROM courses WHERE code='NMAI058'),2,3),
(1,(SELECT id FROM courses WHERE code='NPRG031'),2,4),
(1,(SELECT id FROM courses WHERE code='NSWI170'),2,5),
(1,(SELECT id FROM courses WHERE code='NSWI177'),2,6),
(1,(SELECT id FROM courses WHERE code='NTIN060'),2,7),
(1,(SELECT id FROM courses WHERE code='NTIN107'),2,8),
(1,(SELECT id FROM courses WHERE code='NTVY015'),2,9),
(2,(SELECT id FROM courses WHERE code='NAIL062'),1,1),
(2,(SELECT id FROM courses WHERE code='NDBI025'),1,2),
(2,(SELECT id FROM courses WHERE code='NDMI011'),1,3),
(2,(SELECT id FROM courses WHERE code='NJAZ074'),1,4),
(2,(SELECT id FROM courses WHERE code='NPRG041'),1,5),
(2,(SELECT id FROM courses WHERE code='NSWI142'),1,6),
(2,(SELECT id FROM courses WHERE code='NTIN061'),1,7),
(2,(SELECT id FROM courses WHERE code='NTVY016'),1,8),
(2,(SELECT id FROM courses WHERE code='NJAZ091'),2,1),
(2,(SELECT id FROM courses WHERE code='NJAZ176'),2,2),
(2,(SELECT id FROM courses WHERE code='NMAI059'),2,3),
(2,(SELECT id FROM courses WHERE code='NPRG024'),2,4),
(2,(SELECT id FROM courses WHERE code='NPRG036'),2,5),
(2,(SELECT id FROM courses WHERE code='NPRG045'),2,6),
(2,(SELECT id FROM courses WHERE code='NPRG051'),2,7),
(2,(SELECT id FROM courses WHERE code='NSWI143'),2,8),
(2,(SELECT id FROM courses WHERE code='NSWI153'),2,10),
(2,(SELECT id FROM courses WHERE code='NTIN071'),2,11),
(2,(SELECT id FROM courses WHERE code='NTVY017'),2,12),
(3,(SELECT id FROM courses WHERE code='NPFL129'),1,1),
(3,(SELECT id FROM courses WHERE code='NPGR003'),1,2),
(3,(SELECT id FROM courses WHERE code='NPRG035'),1,3),
(3,(SELECT id FROM courses WHERE code='NPRG073'),1,4),
(3,(SELECT id FROM courses WHERE code='NSWI004'),1,5),
(3,(SELECT id FROM courses WHERE code='NSWI098'),1,6),
(3,(SELECT id FROM courses WHERE code='NSWI154'),1,7),
(3,(SELECT id FROM courses WHERE code='NPRG038'),2,1),
(3,(SELECT id FROM courses WHERE code='NPRG043'),2,2),
(3,(SELECT id FROM courses WHERE code='NPRG074'),2,3),
(3,(SELECT id FROM courses WHERE code='NSWI041'),2,4),
(3,(SELECT id FROM courses WHERE code='NSZZ031'),2,5);