package coursedetail

import (
	"github.com/jmoiron/sqlx"
	"github.com/michalhercik/RecSIS/coursedetail/internal/sqlquery"
	"github.com/michalhercik/RecSIS/language"
)

type DBManager struct {
	DB *sqlx.DB
}

func (reader DBManager) Course(sessionID string, code string, lang language.Language) (*Course, error) {
	var course Course
	if err := reader.DB.Get(&course, sqlquery.Course, sessionID, code, lang); err != nil {
		return nil, err
	}
	if err := reader.DB.Select(&course.CategoryRatings, sqlquery.Rating, sessionID, code, lang); err != nil {
		return nil, err
	}
	return &course, nil
}

func (db DBManager) RateCategory(sessionID string, code string, category string, rating int, lang language.Language) ([]CourseCategoryRating, error) {
	var updatedRating []CourseCategoryRating
	_, err := db.DB.Exec(sqlquery.RateCategory, sessionID, code, category, rating)
	if err != nil {
		return updatedRating, err
	}
	if err = db.DB.Select(&updatedRating, sqlquery.Rating, sessionID, code, lang); err != nil {
		return updatedRating, err
	}
	return updatedRating, err
}

func (db DBManager) DeleteCategoryRating(sessionID string, code string, category string, lang language.Language) ([]CourseCategoryRating, error) {
	var updatedRating []CourseCategoryRating
	_, err := db.DB.Exec(sqlquery.DeleteCategoryRating, sessionID, code, category)
	if err != nil {
		return updatedRating, err
	}
	if err = db.DB.Select(&updatedRating, sqlquery.Rating, sessionID, code, lang); err != nil {
		return updatedRating, err
	}
	return updatedRating, err
}

func (db DBManager) Rate(sessionID string, code string, value int) (CourseRating, error) {
	var rating CourseRating
	_, err := db.DB.Exec(sqlquery.Rate, sessionID, code, value)
	if err != nil {
		return rating, err
	}
	if err = db.DB.Get(&rating, sqlquery.CourseOverallRating, sessionID, code); err != nil {
		return rating, err
	}
	return rating, err
}

func (db DBManager) DeleteRating(sessionID string, code string) (CourseRating, error) {
	var rating CourseRating
	_, err := db.DB.Exec(sqlquery.DeleteRating, sessionID, code)
	if err != nil {
		return rating, err
	}
	if err = db.DB.Get(&rating, sqlquery.CourseOverallRating, sessionID, code); err != nil {
		return rating, err
	}
	return rating, err
}
