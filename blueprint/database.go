package blueprint

import (
	"database/sql"
)

type DBManager struct {
	Db *sql.DB
}

func (m DBManager) BluePrint(user int) (*Blueprint, error) {
	// rows, err := m.Db.Query(sqlquery.Course, user)
	// if err != nil {
	// 	return nil, err
	// }
	// defer rows.Close()
	// for rows.Next() {
	// 	if err := rows.Scan(); err != nil {
	// 		return nil, err
	// 	}

	// }
	// if err := rows.Err(); err != nil {
	// 	return nil, err
	// }
	return &Blueprint{years: []AcademicYear{}}, nil
}

func (m DBManager) AddCourse(user int, course string, year int, semester int, position int) {

}

func (db DBManager) RemoveCourse(user int, course string, year int, semester int) {

}

func (db DBManager) AddYear(user int) {

}

func (db DBManager) RemoveYear(user int, year int) {

}
