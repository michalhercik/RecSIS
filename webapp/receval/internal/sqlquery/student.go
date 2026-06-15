package sqlquery

const SaveStudent = `
INSERT INTO receval_saved_students (student_id, title)
VALUES ($1, $2)
`

const DeleteStudent = `
DELETE FROM receval_saved_students WHERE student_id = $1
`

const SelectSavedStudents = `
SELECT student_id, title FROM receval_saved_students
`
