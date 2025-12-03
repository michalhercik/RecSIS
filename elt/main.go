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
	// start = time.Now()
	// err = extract(sis, recsis)
	// elapsed = time.Since(start)
	// report = makeReport(err, elapsed)
	// log.Println("--------------------------------------------------")
	// log.Println(report)
	// // Transform
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
		{table: "studplanlist2searchable", index: "degree-plans"},
	})
	elapsed = time.Since(start)
	report = makeReport(err, elapsed)
	log.Println("--------------------------------------------------")
	log.Println(report)
	// // Migration
	// err = migrate(recsis)
	// if err != nil {
	// 	log.Panicln("Migration failed:", err)
	// }
}

func migrate(db *sqlx.DB) error {
	var err error
	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	err = migrateCourses(tx)
	if err != nil {
		return err
	}
	err = migrateFilters(tx)
	if err != nil {
		return err
	}
	err = migrateStudPlanList(tx)
	if err != nil {
		return err
	}
	err = migrateDegreePlanYears(tx)
	if err != nil {
		return err
	}
	err = migrateDegreePlans(tx)
	if err != nil {
		return err
	}
	err = migrateStudium(tx)
	if err != nil {
		return err
	}
	err = migrateZkous(tx)
	if err != nil {
		return err
	}
	err = migratePovinn(tx)
	if err != nil {
		return err
	}
	err = migrateSearchablePovinn(tx)
	if err != nil {
		return err
	}
	err = migratePreq(tx)
	if err != nil {
		return err
	}
	err = migratePamela(tx)
	if err != nil {
		return err
	}
	err = migrateStudPlan(tx)
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
			studplan2lang,
			studplanlist2searchable,
			degreePlanYears,
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
			createFilterValuesForSurveyAcademicYears,
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
	extract.add(&extractStudPlan{})
	extract.add(&extractStudium{})
	extract.add(&extractZkous{})

	err := extract.run()
	return err
}
