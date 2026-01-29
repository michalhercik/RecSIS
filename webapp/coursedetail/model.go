package coursedetail

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"iter"
	"net/url"
	"sort"

	"github.com/michalhercik/RecSIS/filters"
	"github.com/michalhercik/RecSIS/language"
)

//================================================================================
// Constants
//================================================================================

const (
	defaultSurveyOffset = 0
	minRating           = 0
	maxRating           = 10
	negativeRating      = 0
	positiveRating      = 1
	resultsPerPage      = 20
	ttDelay             = 200 // tooltip delay in ms
)

const (
	searchQuery  = "survey-search"
	surveyOffset = "survey-offset"
	ratingParam  = "rating"
)

const (
	courseCode     = "code"
	ratingCategory = "category"
)

const (
	meiliCourseCode = "course_code"
	meiliSort       = "academic_year:desc"
)

//================================================================================
// Data Types and Methods
//================================================================================

type courseDetailPage struct {
	course *course
}

func urlHostPath(urlStr string) string {
	u, err := url.Parse(urlStr)
	if err != nil {
		return urlStr
	}
	return u.Host + u.Path
}

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
	url                    sql.NullString
	annotation             nullDescription
	syllabus               nullDescription
	passingTerms           nullDescription
	literature             nullDescription
	assessmentRequirements nullDescription
	entryRequirements      nullDescription
	aim                    nullDescription
	prerequisites          requisiteSlice
	corequisites           requisiteSlice
	incompatibles          requisiteSlice
	interchanges           requisiteSlice
	classes                []string
	classifications        []string
	blueprintAssignments   assignmentSlice
	blueprintSemesters     []bool
	inDegreePlan           bool
	categoryRatings        []courseCategoryRating
	overallRating          courseRating
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

type faculty struct {
	abbr string
	name string
}

type department struct {
	id   string
	name string
}

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

type nullRangeUnit struct {
	rangeUnit
	valid bool
}

type rangeUnit struct {
	abbr string
	name string
}

type teacherSlice []teacher

type nullDescription struct {
	description
	valid bool
}

type description struct {
	title   string
	content string
}

type requisiteSlice []requisite

type requisite struct {
	courseCode string
	children   requisiteSlice
	group      sql.NullString
}

func (r requisite) isNode() bool {
	return !r.group.Valid
}

func (r requisite) isDisjunction() bool {
	return r.group.Valid && r.group.String == "M"
}

func (r requisite) isConjunction() bool {
	return r.group.Valid && r.group.String == "V"
}

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

type surveyViewModel struct {
	lang   language.Language
	code   string
	query  string
	survey []survey
	offset int
	isEnd  bool
	facets iter.Seq[filters.FacetIterator]
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
	Type          string  `json:"target_type"`
	CourseCode    string  `json:"course_code"`
	TargetTeacher teacher `json:"teacher"`
}

type teacher struct {
	SisID       string `json:"id"`
	LastName    string `json:"last_name"`
	FirstName   string `json:"first_name"`
	TitleBefore string `json:"title_before"`
	TitleAfter  string `json:"title_after"`
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
		len(c.interchanges) > 0 ||
		len(c.incompatibles) > 0 ||
		len(c.teachers) > 0
}

func (c *course) hoursString() string {
	result := ""
	winter := c.lectureRangeWinter.Valid && c.seminarRangeWinter.Valid
	summer := c.lectureRangeSummer.Valid && c.seminarRangeSummer.Valid
	if winter {
		result += fmt.Sprintf("%d/%d", c.lectureRangeWinter.Int64, c.seminarRangeWinter.Int64)
	}
	if winter && summer {
		result += ", "
	}
	if summer {
		result += fmt.Sprintf("%d/%d", c.lectureRangeSummer.Int64, c.seminarRangeSummer.Int64)
	}
	return result
}

func courseSISLink(code string, t text) string {
	var (
		csLink = "https://is.cuni.cz/studium/predmety/index.php?do=predmet&kod=" + code
		enLink = "https://is.cuni.cz/studium/eng/predmety/index.php?do=predmet&kod=" + code
	)
	if t.language == language.EN {
		return enLink
	}
	return csLink
}
