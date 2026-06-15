SET search_path TO webapp;

DROP TABLE IF EXISTS receval_saved_students CASCADE;

CREATE TABLE receval_saved_students (
    student_id VARCHAR(8),
    title VARCHAR(255)
);
