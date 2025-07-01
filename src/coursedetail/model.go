package coursedetail

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"iter"
	"sort"

	"github.com/michalhercik/RecSIS/filters"
	"github.com/michalhercik/RecSIS/language"
)

//================================================================================
// Constants
//================================================================================

const (
	minRating        = 0
	maxRating        = 10
	negativeRating   = 0
	positiveRating   = 1
	numberOfComments = 20
	searchQuery      = "survey-search"
	surveyOffset     = "survey-offset"
)

//================================================================================
// Data Types and Methods
//================================================================================

// all data needed for the course detail page
type courseDetailPage struct {
	course *course
}

// course representation
type course struct {
	code                   string
	title                  string
	faculty                faculty
	guarantorDepartment    department
	state                  string
	semester               teachingSemester
	language               string
	lectureRangeWinter     sql.NullInt64
	seminarRangeWinter     sql.NullInt64
	lectureRangeSummer     sql.NullInt64
	seminarRangeSummer     sql.NullInt64
	rangeUnit              nullRangeUnit
	examType               string
	credits                int
	guarantors             teacherSlice
	teachers               teacherSlice
	capacity               string
	annotation             nullDescription
	syllabus               nullDescription
	passingTerms           nullDescription
	literature             nullDescription
	assessmentRequirements nullDescription
	entryRequirements      nullDescription
	aim                    nullDescription
	prerequisites          []requisite
	corequisites           []requisite
	incompatible           []requisite
	interchange            []requisite
	classes                []string
	classifications        []string
	link                   string // link to course webpage (not SIS)
	blueprintAssignments   assignmentSlice
	blueprintSemesters     []bool
	inDegreePlan           bool
	categoryRatings        []courseCategoryRating
	overallRating          courseRating
}

type rangeUnit struct {
	abbr string
	name string
}

type nullRangeUnit struct {
	rangeUnit
	valid bool
}

type faculty struct {
	abbr string
	name string
}

type department struct {
	id   string
	name string
}

type requisite struct {
	courseCode string
	state      string
}

func (c course) semesterStyleClass() string {
	switch c.semester {
	case teachingBoth:
		return "bg-both"
	case teachingWinterOnly:
		return "bg-winter"
	case teachingSummerOnly:
		return "bg-summer"
	default:
		return ""
	}
}

// semester type - winter, summer, or both
type teachingSemester int

const (
	teachingWinterOnly teachingSemester = iota + 1
	teachingSummerOnly
	teachingBoth
)

func (ts teachingSemester) string(t text) string {
	switch ts {
	case teachingBoth:
		return t.both
	case teachingWinterOnly:
		return t.winter
	case teachingSummerOnly:
		return t.summer
	default:
		return ""
	}
}

// wrapper for teacher slice
type teacherSlice []teacher

// wrapper for Description that allows it to be nullable
type nullDescription struct {
	description
	valid bool
}

// description type - title and content
type description struct {
	title   string
	content string
}

// categorization of course
type class struct {
	code string
	name string
}

// assignment slice for blueprint assignments
type assignmentSlice []assignment

func (a assignmentSlice) sort() assignmentSlice {
	sort.Slice(a, func(i, j int) bool {
		if a[i].year == a[j].year {
			return a[i].semester < a[j].semester
		}
		return a[i].year < a[j].year
	})
	return a
}

// assignment type - year and semester
type assignment struct {
	year     int
	semester semesterAssignment
}

func (a assignment) string(lang language.Language) string {
	t := texts[lang]
	semester := ""
	switch a.semester {
	case assignmentNone:
		semester = "unsupported"
	case assignmentWinter:
		semester = t.winterAssign
	case assignmentSummer:
		semester = t.summerAssign
	default:
		semester = "unsupported"
	}

	result := fmt.Sprintf("%s %s", t.yearStr(a.year), semester)
	if a.year == 0 {
		result = t.unassigned
	}
	return result
}

// semester assignment type - winter or summer (none = unassigned)
type semesterAssignment int

const (
	assignmentNone semesterAssignment = iota
	assignmentWinter
	assignmentSummer
)

func (sa semesterAssignment) stringID() string {
	switch sa {
	case assignmentNone:
		return "none"
	case assignmentWinter:
		return "winter"
	case assignmentSummer:
		return "summer"
	default:
		return "unsupported"
	}
}

// rating structures
type courseCategoryRating struct {
	code  int
	title string
	courseRating
}

type courseRating struct {
	userRating  sql.NullInt64
	avgRating   sql.NullFloat64
	ratingCount sql.NullInt64
}

// surveys structs and methods
type surveyViewModel struct {
	lang   language.Language
	code   string
	query  string
	survey []survey
	offset int
	isEnd  bool
	facets iter.Seq[filters.FacetIterator] // TODO
}

type survey struct {
	student
	surveyTarget
	AcademicYear int    `json:"academic_year"`
	Content      string `json:"content"`
}

type student struct {
	Year  int        `json:"study_year"`
	Field StudyField `json:"study_field"`
	Study studyType  `json:"study_type"`
}

type StudyField struct {
	ID   string
	Name string
}

func (s *StudyField) UnmarshalJSON(data []byte) error {
	var temp struct {
		ID   string `json:"id"`
		Name struct {
			CS string `json:"cs"`
			EN string `json:"en"`
		} `json:"name"`
	}
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}
	s.ID = temp.ID
	s.Name = temp.Name.CS
	if s.Name == "" {
		s.Name = temp.Name.EN
	}
	return nil
}

type studyType struct {
	Abbr string
	Name string
}

func (s *studyType) UnmarshalJSON(data []byte) error {
	var temp struct {
		Abbr struct {
			CS string `json:"cs"`
			EN string `json:"en"`
		} `json:"abbr"`
		Name struct {
			CS string `json:"cs"`
			EN string `json:"en"`
		} `json:"name"`
	}
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}
	s.Abbr = temp.Abbr.CS
	s.Name = temp.Name.CS
	if s.Abbr == "" && s.Name == "" {
		s.Abbr = temp.Abbr.EN
		s.Name = temp.Name.EN
	}
	return nil
}

type surveyTarget struct {
	Type          string  `json:"target_type"` // Lecture or Seminar
	CourseCode    string  `json:"course_code"`
	TargetTeacher teacher `json:"teacher"`
}

type teacher struct {
	SisID       string `json:"KOD"`
	LastName    string `json:"PRIJMENI"`
	FirstName   string `json:"JMENO"`
	TitleBefore string `json:"TITULPRED"`
	TitleAfter  string `json:"TITULZA"`
}

// ================================================================================
// Helper Functions
// ================================================================================

func (c course) hasDescriptions() bool {
	return c.annotation.valid ||
		c.syllabus.valid ||
		c.passingTerms.valid ||
		c.literature.valid ||
		c.assessmentRequirements.valid ||
		c.entryRequirements.valid ||
		c.aim.valid
}

func (c course) hasDetailInfo() bool {
	return len(c.classes) > 0 ||
		len(c.classifications) > 0 ||
		len(c.corequisites) > 0 ||
		len(c.interchange) > 0 ||
		len(c.incompatible) > 0 ||
		len(c.teachers) > 0
}

func hoursString(course *course) string {
	result := ""
	winter := course.lectureRangeWinter.Valid && course.seminarRangeWinter.Valid
	summer := course.lectureRangeSummer.Valid && course.seminarRangeSummer.Valid
	if winter {
		result += fmt.Sprintf("%d/%d", course.lectureRangeWinter.Int64, course.seminarRangeWinter.Int64)
	}
	if winter && summer {
		result += ", "
	}
	if summer {
		result += fmt.Sprintf("%d/%d", course.lectureRangeSummer.Int64, course.seminarRangeSummer.Int64)
	}
	return result
}

func courseSISLink(code string, t text) string {
	if t.language == language.CS {
		return "https://is.cuni.cz/studium/predmety/index.php?do=predmet&kod=" + code
	} else if t.language == language.EN {
		return "https://is.cuni.cz/studium/eng/predmety/index.php?do=predmet&kod=" + code
	}
	// default to Czech if language is not recognized
	return "https://is.cuni.cz/studium/predmety/index.php?do=predmet&kod=" + code
}
