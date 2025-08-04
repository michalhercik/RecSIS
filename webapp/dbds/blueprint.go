package dbds

type BlueprintAssignment struct {
	ID       int `db:"id"`
	Year     int `db:"academic_year"`
	Semester int `db:"semester"`
}
