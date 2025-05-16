package dbds

import "github.com/lib/pq"

type DegreePlanRecord struct {
	BlocCode       int           `db:"bloc_subject_code"`
	BlocLimit      int           `db:"bloc_limit"`
	BlocName       string        `db:"bloc_name"`
	BlocNote       string        `db:"bloc_note"`
	Note           string        `db:"note"`
	BlueprintYears pq.Int64Array `db:"academic_years"`
	Course
}
