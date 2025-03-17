package coursedetail

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/michalhercik/RecSIS/coursedetail/internal/sqlquery"
)

type DBManager struct {
	DB *sqlx.DB
}

func (reader DBManager) Course(sessionID string, code string, lang DBLang) (*Course, error) {
	var course Course
	var response []struct {
		CourseInfo
		OverallRating NullInt64      `db:"overall_rating"`
		Category      sql.NullInt64  `db:"category_code"`
		RatingTitle   sql.NullString `db:"rating_title"`
		Rating        sql.NullInt64  `db:"rating"`
	}
	if err := reader.DB.Select(&response, sqlquery.Course, sessionID, code, lang); err != nil {
		return nil, err
	}
	course.CourseInfo = response[0].CourseInfo
	course.OverallRating = response[0].OverallRating
	if response[0].Category.Valid {
		for _, r := range response {
			course.CategoryRatings = append(course.CategoryRatings, CourseCategoryRating{
				Code:   int(r.Category.Int64),
				Title:  r.RatingTitle.String,
				Rating: int(r.Rating.Int64),
			})
		}
	}
	return &course, nil
}

func (db DBManager) RateCategory(sessionID string, code string, category string, rating int) error {
	_, err := db.DB.Exec(sqlquery.RateCategory, sessionID, code, category, rating)
	return err
}

func (db DBManager) DeleteCategoryRating(sessionID string, code string, category string) error {
	_, err := db.DB.Exec(sqlquery.DeleteCategoryRating, sessionID, code, category)
	return err
}

func (db DBManager) Rate(sessionID string, code string, value int) error {
	_, err := db.DB.Exec(sqlquery.Rate, sessionID, code, value)
	return err
}

func (db DBManager) DeleteRating(sessionID string, code string) error {
	_, err := db.DB.Exec(sqlquery.DeleteRating, sessionID, code)
	return err
}
