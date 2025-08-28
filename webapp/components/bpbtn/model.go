package bpbtn

const (
	yearParam     = "year"
	semesterParam = "semester"
	courseParam   = "course"
)

const endpointPath = "blueprint"

const (
	uniqueViolationCode       = "23505"
	duplicateCoursesViolation = "blueprint_courses_blueprint_semester_id_course_code_key"
)

type ViewModel struct {
	course     string
	semesters  []bool
	hxPostBase string
	hxSwap     string
	hxTarget   string
	hxInclude  string
}

type overYearsIterator struct {
	disableWinter bool
	disableSummer bool
}
