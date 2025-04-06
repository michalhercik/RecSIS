package coursedetail

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/michalhercik/RecSIS/coursedetail/internal/sqlquery"
	"github.com/michalhercik/RecSIS/language"
)

type DBManager struct {
	DB *sqlx.DB
}

func (reader DBManager) Course(sessionID string, code string, lang language.Language) (*Course, error) {
	var course Course
	var response []struct {
		CourseInfo
		UserOverallRating NullInt64       `db:"overall_rating"`
		Category          sql.NullInt64   `db:"category_code"`
		RatingTitle       sql.NullString  `db:"rating_title"`
		UserRating        sql.NullInt64   `db:"rating"`
		AvgOverallRating  NullFloat64     `db:"avg_overall_rating"`
		AvgRating         sql.NullFloat64 `db:"avg_rating"`
	}
	if err := reader.DB.Select(&response, sqlquery.Course, sessionID, code, lang); err != nil {
		return nil, err
	}
	course.CourseInfo = response[0].CourseInfo
	course.UserOverallRating = response[0].UserOverallRating
	course.AvgOverallRating = response[0].AvgOverallRating
	if response[0].Category.Valid {
		for _, r := range response {
			course.CategoryRatings = append(course.CategoryRatings, CourseCategoryRating{
				Code:       int(r.Category.Int64),
				Title:      r.RatingTitle.String,
				UserRating: int(r.UserRating.Int64),
				AvgRating:  float64(r.AvgRating.Float64),
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
