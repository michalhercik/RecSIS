COPY faculties(sis_id, sis_poid, name_cs, name_en, abbr)
FROM '/docker-entrypoint-initdb.d/faculties.csv'
DELIMITER ','
CSV HEADER;

COPY teachers(sis_id,department,faculty,last_name,first_name,title_before,title_after)
FROM '/docker-entrypoint-initdb.d/teachers.csv'
DELIMITER ','
CSV HEADER;

COPY courses(code,name_cs,name_en,valid_from,valid_to,faculty,guarantor,taught,start_semester,semester_count,taught_lang,lecture_range1,seminar_range1,lecture_range2,seminar_range2,range_unit,exam_type,credits,teacher1,teacher2,min_number,capacity,annotation_cs,annotation_en,sylabus_cs,sylabus_en)
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