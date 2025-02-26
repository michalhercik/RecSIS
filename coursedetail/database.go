package coursedetail

import (
	"github.com/jmoiron/sqlx"
	"github.com/michalhercik/RecSIS/coursedetail/internal/sqlquery"
)

type DBManager struct {
	DB *sqlx.DB
}

func (reader DBManager) Course(code string, lang DBLang) (*Course, error) {
	var course Course
	if err := reader.DB.Get(&course, sqlquery.Course, code, lang); err != nil {
		return nil, err
	}
	// TODO: this is mock - change to real data
	course.Comments = []Comment{
		{ID: 1, UserID: 1, Content: "This is a comment"},
		{ID: 2, UserID: 2, Content: "This is another comment"},
		{ID: 3, UserID: 3, Content: "This is yet another comment"},
	}
	course.Ratings = []Rating{
		{ID: 1, UserID: 1, Rating: 1},
		{ID: 2, UserID: 2, Rating: 1},
		{ID: 3, UserID: 3, Rating: -1},
	}

	return &course, nil
}

// TODO: MOCK - implement
func (reader DBManager) AddComment(code, commentContent string) error {
	//_, err := reader.DB.Exec(sqlquery.AddComment, code, commentContent)
	//return err
	return nil
}

// TODO: MOCK - implement
func (reader DBManager) GetComments(code string) ([]Comment, error) {
	comments := []Comment{
		{ID: 1, UserID: 1, Content: "This is a comment"},
		{ID: 2, UserID: 2, Content: "This is another comment"},
		{ID: 3, UserID: 3, Content: "This is yet another comment"},
		{ID: 4, UserID: 4, Content: "I think that Michal is great name"},
	}
	return comments, nil
}

func (m DBManager) AddCourseToBlueprint(user int, code string) ([]Assignment, error) {
	// TODO: Implement this method
	// year=0, semester=course.semester, position=-1
	// this must be done in a transaction, must return the all assignments
	return nil, nil
}

func (m DBManager) RemoveCourseFromBlueprint(user int, code string) error {
	// TODO: Implement this method
	// it is expected that there is only one course with the given code
	// if not return an error
	// by that we do not have to return the assignments - there should be none after the removal
	return nil
}
