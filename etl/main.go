package main

import (
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {
	// Load Config
	conf, err := loadConfig()
	if err != nil {
		log.Panicln(err)
	}
	// Create Connections
	sis, err := createSISConn(conf)
	if err != nil {
		log.Panicln(err)
	}
	defer sis.Close()
	recsis, err := createRecSISConn(conf)
	if err != nil {
		log.Panicln(err)
	}
	defer recsis.Close()
	meili, err := createMeilisearchConn(conf)
	if err != nil {
		log.Panicln(err)
	}
	defer meili.Close()

	var report Report
	var start time.Time
	var elapsed time.Duration
	// Extract
	start = time.Now()
	err = extract(sis, recsis)
	elapsed = time.Since(start)
	report = makeReport(err, elapsed)
	log.Println("--------------------------------------------------")
	log.Println(report)
	// Transform
	start = time.Now()
	err = transform(recsis)
	elapsed = time.Since(start)
	report = makeReport(err, elapsed)
	log.Println("--------------------------------------------------")
	log.Println(report)
	// Meilisearch
	start = time.Now()
	err = uploadToMeili(recsis, meili, []meiliUpload{
		{table: "povinn2searchable", index: "courses"},
		{table: "ankecy2searchable", index: "survey"},
	})
	elapsed = time.Since(start)
	report = makeReport(err, elapsed)
	log.Println("--------------------------------------------------")
	log.Println(report)
	// Migration
	err = migrate(recsis)
	if err != nil {
		log.Panicln("Migration failed:", err)
	}
}

func migrate(db *sqlx.DB) error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}
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
			prerequisities,
			corequisities,
			incompatibilities,
			interchangeabilities,
			classes,
			classifications
		FROM povinn2courses;
	`)
	if err != nil {
		return err
	}
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
	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}

func transform(recsis *sqlx.DB) error {
	runner := sequentialRunner{
		parallelRunner{
			fak2JSON,
			ustav2JSON,
			ucit2JSON,
			pamela2JSON,
			rvcem2JSON,
			druh2JSON,
			preq2JSON,
			ptrida2JSON,
			pvyuc2lang,
			jazyk2lang,
			zsem2lang,
			typypov2lang,
			obor2lang,
			klas2lang,
			initFilterTables,
		},
		parallelRunner{
			pklas2JSON,
			povinn2jazykAgg,
			ankecy2JSON,
		},
		parallelRunner{
			povinn2courses,
		},
		parallelRunner{
			povinn2searchable,
		},
		parallelRunner{
			ankecy2searchable,
			createFilterValuesForCredits,
			createFilterValuesForLangs,
			createFilterValuesForFaculties,
			createFilterValuesForTaughtStates,
			createFilterValuesForStartSemesters,
			createFilterValuesForSemesterCounts,
			createFilterValuesForLectureRanges,
			createFilterValuesForSeminarRanges,
			createFilterValuesForExams,
			createFilterValuesForRangeUnits,
			createFilterValuesForDepartments,
			createFilterValuesForSections,
			createFilterValuesForSurveyTeachers,
			createFilterValuesForSurveyStudyFields,
			createFilterValuesForSurveyStudyTypes,
			createFilterValuesForSurveyStudyYears,
			createFilterValuesForSurveyTargetTypes,
		},
	}

	err := runner.run(recsis)
	if err != nil {
		return err
	}
	return nil
}

func extract(sis, recsis *sqlx.DB) error {
	extract := makeExtract(sis, recsis)
	extract.add(&extractPovinn{})
	extract.add(&extractUcitRozvrh{})
	extract.add(&extractUcit{})
	extract.add(&extractAnkecy{})
	extract.add(&extractDruh{})
	extract.add(&extractJazyk{})
	extract.add(&extractKlas{})
	extract.add(&extractPamela{})
	extract.add(&extractPklas{})
	extract.add(&extractPovinn2Jazyk{})
	extract.add(&extractPreq{})
	extract.add(&extractPtrida{})
	extract.add(&extractTrida{})
	extract.add(&extractTypyPov{})
	extract.add(&extractSekce{})
	extract.add(&extractUstav{})
	extract.add(&extractFak{})
	extract.add(&Ciselnik{Table: "rvcem"})
	extract.add(&Ciselnik{Table: "zsem"})
	extract.add(&Ciselnik{Table: "pvyuc"})
	extract.add(&Ciselnik{Table: "typmem"})
	extract.add(&Ciselnik{Table: "obor", KodSize: 12, NazevSize: 250})

	err := extract.run()
	return err
}

// DROP TABLE IF EXISTS
// ankecy2json,
// druh2json,
// fak2json,
// filter_values,
// filter_categories,
// filters,
// jazyk2lang,
// klas2lang,
// obor2lang,
// pamela2json,
// povinn2jazyk_agg,
// preq2json,
// ptrida2json,
// pvyuc2lang,
// rvcem2json,
// typypov2lang,
// ucit2json,
// ustav2json,
// zsem2lang,
// povinn2searchable,
// povinn2courses,
// pklas2json
