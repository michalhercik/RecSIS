package main

import "github.com/jmoiron/sqlx"

func migrateCourses(tx *sqlx.Tx) error {
	var err error
	_, err = tx.Exec(`--sql
		DELETE FROM webapp.courses WHERE TRUE;
		INSERT INTO webapp.courses
		SELECT
			code,
			lang,
			title,
			valid_from,
			valid_to,
			course_url,
			faculty,
			department,
			taught_state,
			taught_state_title,
			start_semester,
			start_semester_title,
			taught_lang,
			lecture_range_winter,
			seminar_range_winter,
			lecture_range_summer,
			seminar_range_summer,
			range_unit,
			exam,
			credits,
			guarantors,
			teachers,
			min_occupancy,
			capacity,
			annotation,
			syllabus,
			terms_of_passing,
			literature,
			requirements_of_assesment,
			entry_requirements,
			aim,
			classes,
			classifications
		FROM povinn2courses;
	`)
	if err != nil {
		return err
	}
	return nil
}

func migrateRequisites(tx *sqlx.Tx) error {
	var err error
	_, err = tx.Exec(`--sql
		DELETE FROM webapp.requisites WHERE TRUE;
		INSERT INTO webapp.requisites
		SELECT
			target_course,
			parent_course,
			child_course,
			req_type,
			group_type
		FROM requisites;
	`)
	if err != nil {
		return err
	}
	return nil
}

func migrateFilters(tx *sqlx.Tx) error {
	var err error
	_, err = tx.Exec(`--sql
		DELETE FROM webapp.filter_values WHERE TRUE;
		DELETE FROM webapp.filter_categories WHERE TRUE;
		DELETE FROM webapp.filters WHERE TRUE;
		INSERT INTO webapp.filters
		SELECT * from filters;
		INSERT INTO webapp.filter_categories
		SELECT * from filter_categories;
		INSERT INTO webapp.filter_values
		SELECT * from filter_values
	`)
	if err != nil {
		return err
	}
	return nil
}

func migrateStudPlanList(tx *sqlx.Tx) error {
	var err error
	_, err = tx.Exec(`--sql
		DELETE FROM webapp.degree_plan_list WHERE TRUE;
		INSERT INTO webapp.degree_plan_list (code)
		SELECT DISTINCT plan_code
		FROM studmetadata2lang;
	`)
	if err != nil {
		return err
	}
	return nil
}

func migrateStudPlanMetadata(tx *sqlx.Tx) error {
	var err error
	_, err = tx.Exec(`--sql
		DELETE FROM webapp.degree_plan_metadata WHERE TRUE;
		INSERT INTO webapp.degree_plan_metadata (
			plan_code,
			lang,
			title,
			valid_from,
			valid_to,
			faculty,
			section,
			field_code,
			study_type
		) SELECT DISTINCT
			plan_code,
			lang,
			title,
			valid_from,
			valid_to,
			faculty,
			section,
			field_code,
			study_type
		FROM studmetadata2lang;
	`)
	if err != nil {
		return err
	}
	return nil
}

func migrateStudPlans(tx *sqlx.Tx) error {
	var err error
	_, err = tx.Exec(`--sql
		DELETE FROM webapp.degree_plans WHERE TRUE;
		INSERT INTO webapp.degree_plans (
			plan_code,
			lang,
			course_code,
			interchangeability,
			recommended_year_from,
			recommended_year_to,
			recommended_semester,
			bloc_name,
			bloc_subject_code,
			bloc_type,
			bloc_limit,
			seq
		) SELECT
			plan_code,
			lang,
			course_code,
			interchangeability,
			recommended_year_from,
			recommended_year_to,
			recommended_semester,
			bloc_name,
			bloc_subject_code,
			bloc_type,
			bloc_limit,
			seq
		FROM studplan2lang;
	`)
	if err != nil {
		return err
	}
	return nil
}
