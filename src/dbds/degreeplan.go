package dbds

type DegreePlanRecord struct {
	BlocCode         int    `db:"bloc_subject_code"`
	BlocLimit        int    `db:"bloc_limit"`
	BlocName         string `db:"bloc_name"`
	BlocNote         string `db:"bloc_note"`
	IsBlocCompulsory bool   `db:"is_compulsory"`
	Note             string `db:"note"`
	Course
}
