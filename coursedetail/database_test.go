package coursedetail

import (
	"database/sql"
	"fmt"
	"testing"

	_ "github.com/lib/pq"
)

func equalFaculty(f1, f2 *Faculty) error {
	if f1.Id != f2.Id {
		return fmt.Errorf("Id: %d != %d", f1.Id, f2.Id)
	}
	if f1.SisId != f2.SisId {
		return fmt.Errorf("SisId: %d != %d", f1.SisId, f2.SisId)
	}
	if f1.NameCs != f2.NameCs {
		return fmt.Errorf("nameCs: %s != %s", f1.NameCs, f2.NameCs)
	}
	if f1.NameEn != f2.NameEn {
		return fmt.Errorf("NameEn: %s != %s", f1.NameEn, f2.NameEn)
	}
	if f1.Abbr != f2.Abbr {
		return fmt.Errorf("Abbr: %s != %s", f1.Abbr, f2.Abbr)
	}
	return nil
}

func equalTeacher(t1, t2 *Teacher) error {
	if t1.Id != t2.Id {
		return fmt.Errorf("Id: %d != %d", t1.Id, t2.Id)
	}
	if t1.SisId != t2.SisId {
		return fmt.Errorf("SisId: %d != %d", t1.SisId, t2.SisId)
	}
	if t1.Department != t2.Department {
		return fmt.Errorf("Department: %s != %s", t1.Department, t2.Department)
	}
	err := equalFaculty(&t1.Faculty, &t2.Faculty)
	if err != nil {
		return fmt.Errorf("Faculty: %v", err)
	}
	if t1.FirstName != t2.FirstName {
		return fmt.Errorf("FirstName: %s != %s", t1.FirstName, t2.FirstName)
	}
	if t1.LastName != t2.LastName {
		return fmt.Errorf("LastName: %s != %s", t1.LastName, t2.LastName)
	}
	if t1.TitleBefore != t2.TitleBefore {
		return fmt.Errorf("TitleBefore: %s != %s", t1.TitleBefore, t2.TitleBefore)
	}
	if t1.TitleAfter != t2.TitleAfter {
		return fmt.Errorf("TitleAfter: %s != %s", t1.TitleAfter, t2.TitleAfter)
	}
	return nil
}

func equalCourse(c1, c2 *Course) error {
	if c1.Id != c2.Id {
		return fmt.Errorf("Id: %d != %d", c1.Id, c2.Id)
	}
	if c1.Code != c2.Code {
		return fmt.Errorf("Code: %s != %s", c1.Code, c2.Code)
	}
	if c1.NameCs != c2.NameCs {
		return fmt.Errorf("NameCs: %s != %s", c1.NameCs, c2.NameCs)
	}
	if c1.NameEn != c2.NameEn {
		return fmt.Errorf("NameEn: %s != %s", c1.NameEn, c2.NameEn)
	}
	if c1.ValidFrom != c2.ValidFrom {
		return fmt.Errorf("ValidFrom: %d != %d", c1.ValidFrom, c2.ValidFrom)
	}
	if c1.ValidTo != c2.ValidTo {
		return fmt.Errorf("ValidTo: %d != %d", c1.ValidTo, c2.ValidTo)
	}
	err := equalFaculty(&c1.Faculty, &c2.Faculty)
	if err != nil {
		return fmt.Errorf("Faculty: %v", err)
	}
	if c1.Guarantor != c2.Guarantor {
		return fmt.Errorf("Guarantor: %s != %s", c1.Guarantor, c2.Guarantor)
	}
	if c1.State != c2.State {
		return fmt.Errorf("State: %s != %s", c1.State, c2.State)
	}
	if c1.Start != c2.Start {
		return fmt.Errorf("Start: %d != %d", c1.Start, c2.Start)
	}
	if c1.SemesterCount != c2.SemesterCount {
		return fmt.Errorf("SemesterCount: %d != %d", c1.SemesterCount, c2.SemesterCount)
	}
	if c1.Language != c2.Language {
		return fmt.Errorf("Language: %s != %s", c1.Language, c2.Language)
	}
	if c1.LectureRange1 != c2.LectureRange1 {
		return fmt.Errorf("LectureRange1: %d != %d", c1.LectureRange1, c2.LectureRange1)
	}
	if c1.SeminarRange1 != c2.SeminarRange1 {
		return fmt.Errorf("SeminarRange1: %d != %d", c1.SeminarRange1, c2.SeminarRange1)
	}
	if c1.LectureRange2 != c2.LectureRange2 {
		return fmt.Errorf("LectureRange2: %v != %v", c1.LectureRange2, c2.LectureRange2)
	}
	if c1.SeminarRange2 != c2.SeminarRange2 {
		return fmt.Errorf("SeminarRange2: %v != %v", c1.SeminarRange2, c2.SeminarRange2)
	}
	if c1.ExamType != c2.ExamType {
		return fmt.Errorf("ExamType: %s != %s", c1.ExamType, c2.ExamType)
	}
	if c1.Credits != c2.Credits {
		return fmt.Errorf("Credits: %d != %d", c1.Credits, c2.Credits)
	}
	for i := range c1.Teachers {
		err = equalTeacher(&c1.Teachers[i], &c2.Teachers[i])
		if err != nil {
			return fmt.Errorf("Teacher1: %v", err)
		}
	}
	if c1.MinEnrollment != c2.MinEnrollment {
		return fmt.Errorf("MinEnrollment: %d != %d", c1.MinEnrollment, c2.MinEnrollment)
	}
	if c1.Capacity != c2.Capacity {
		return fmt.Errorf("Capacity: %d != %d", c1.Capacity, c2.Capacity)
	}
	if c1.AnnotationCs != c2.AnnotationCs {
		return fmt.Errorf("AnnotationCs: \n%s!=\n%s", c1.AnnotationCs, c2.AnnotationCs)
	}
	if c1.AnnotationEn != c2.AnnotationEn {
		return fmt.Errorf("AnnotationEn: \n%s!=\n%s", c1.AnnotationEn, c2.AnnotationEn)
	}
	if c1.SylabusCs != c2.SylabusCs {
		return fmt.Errorf("SylabusCs: \n%s!=\n%s", c1.SylabusCs, c2.SylabusCs)
	}
	if c1.SylabusEn != c2.SylabusEn {
		return fmt.Errorf("SylabusEn: \n%s!=\n%s", c1.SylabusEn, c2.SylabusEn)
	}
	if c1.Link != c2.Link {
		return fmt.Errorf("Link: \n%s!=\n%s", c1.Link, c2.Link)
	}
	return nil
}

func createExpected() *Course {
	return &Course{
		Id:        81,
		Code:      "NPFL104",
		NameCs:    "Metody strojového učení",
		NameEn:    "Machine Learning Methods",
		ValidFrom: 2020,
		ValidTo:   9999,
		Faculty: Faculty{
			Id:     3,
			SisId:  11320,
			NameCs: "Matematicko-fyzikální fakulta",
			NameEn: "Faculty of Mathematics and Physics",
			Abbr:   "MFF",
		},
		Guarantor:     "32-UFAL",
		State:         "N",
		Start:         2,
		SemesterCount: 1,
		Language:      "Czech",
		LectureRange1: 1,
		SeminarRange1: 2,
		LectureRange2: -1,
		SeminarRange2: -1,
		ExamType:      "*",
		Credits:       4,
		Teachers: []Teacher{
			Teacher{
				Id:         16,
				SisId:      11275,
				Department: "32-UFAL",
				Faculty: Faculty{
					Id:     3,
					SisId:  11320,
					NameCs: "Matematicko-fyzikální fakulta",
					NameEn: "Faculty of Mathematics and Physics",
					Abbr:   "MFF",
				},
				FirstName:   "Zdeněk",
				LastName:    "Žabokrtský",
				TitleBefore: "doc. Ing.",
				TitleAfter:  "Ph.D.",
			},
			Teacher{
				Id:         5,
				SisId:      10876,
				Department: "32-UFAL",
				Faculty: Faculty{
					Id:     3,
					SisId:  11320,
					NameCs: "Matematicko-fyzikální fakulta",
					NameEn: "Faculty of Mathematics and Physics",
					Abbr:   "MFF",
				},
				FirstName:   "Ondřej",
				LastName:    "Bojar",
				TitleBefore: "doc. RNDr.",
				TitleAfter:  "Ph.D.",
			},
		},
		MinEnrollment:   10876,
		Capacity:        -1,
		Classifications: nil,
		Classes:         nil,
		AnnotationCs: `Kurs je zaměřen na získání praktických zkušeností s aplikací technik strojového učení na reálná data. U studentů je 
očekávána znalost základních pojmů z oblasti strojového učení. V přednášce jsou stručně zopakovány vybrané 
metody klasifikace, regrese a shlukové analýzy a dále probrány některé přístupy ke zvyšování jejich úspěšnosti, 
například regularizace, transformace množin rysů, diagnostika. Cvičení jsou zaměřena jak na vlastní 
implementace několika metod strojového učení, tak na seznámení se s existujícími implementacemi v jazyce 
Python. 
`,
		AnnotationEn: `The course is focused on practical exercises with applying
machine learning techniques to real data. Students are expected
to be familiar with basic machine learning concepts.
`,
		SylabusCs: `- vlastní implementace základních metod pro klasifikaci a regresi
- seznámení s vybranými knihovnami pro ML
- experimentální srovnávání charakteristik různých klasifikačních metod 
- výběr rysů
- kombinace modelů
- implementace základních technik neřízeného učení
`,
		SylabusEn: `- implementation of basic ML methods for classification and regression
- learning to use selected ML libraries
- experimental comparison of performance characteristics of different classification
methods
- feature engineering
- ensemble techniques
- implementation of basic techniques of unsupervised ML
`,
	}
}

func newDb() (*sql.DB, error) {
	const (
		host     = "localhost"
		port     = 5432
		user     = "recsis"
		password = "recsis"
		dbname   = "recsis"
	)
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	return db, err
}

func TestDbCourse(t *testing.T) {
	expected := createExpected()
	db, err := newDb()
	if err != nil {
		t.Fatalf("error opening database: %v", err)
	}
	defer db.Close()

	courseReader := DbCourseReader{Db: db}
	course, err := courseReader.Course(expected.Code)
	if err != nil {
		t.Fatalf("error getting course: %v", err)
	}
	err = equalCourse(course, expected)
	if err != nil {
		t.Fatalf("%v", err)
	}
}
