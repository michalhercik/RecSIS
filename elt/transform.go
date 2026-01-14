package main

import (
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/jmoiron/sqlx"
)

type runner interface {
	run(db *sqlx.DB) error
}

type sequentialRunner []runner

func (sr sequentialRunner) run(db *sqlx.DB) error {
	var errs []error
	for _, op := range sr {
		if err := op.run(db); err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) > 0 {
		return listOfErrors(errs)
	}
	return nil
}

type parallelRunner []runner

func (pr parallelRunner) run(db *sqlx.DB) error {
	wg := sync.WaitGroup{}
	errsCh := make(chan error, len(pr))
	for _, op := range pr {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := op.run(db); err != nil {
				errsCh <- err
			}
		}()
	}
	wg.Wait()
	close(errsCh)
	var errs []error
	for err := range errsCh {
		errs = append(errs, err)
	}
	if len(errs) > 0 {
		return listOfErrors(errs)
	}
	return nil
}

type listOfErrors []error

func (le listOfErrors) Error() string {
	var sb strings.Builder
	for i, err := range le {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(err.Error())
	}
	return fmt.Sprintf("Errors: %s", sb.String())
}

type transformation struct {
	name  string
	query string
}

func (t transformation) run(db *sqlx.DB) error {
	_, err := db.Exec(t.query)
	if err != nil {
		log.Printf("❌ %s: Transform: %v", t.name, err)
		return err
	}
	log.Printf("✅ Transformation of %s successfull", t.name)
	return nil
}

/*
Prerequisites:
  - fak
*/
var fak2JSON = transformation{
	name: "fak2json",
	query: `--sql
		DROP TABLE IF EXISTS fak2json;
		CREATE TABLE fak2json (
			id VARCHAR(5),
			lang VARCHAR(2),
			faculty jsonb
		);
		INSERT INTO fak2json(id, lang, faculty)
		SELECT kod id, 'cs' lang, jsonb_object(
			ARRAY['name', 'abbr'],
			ARRAY[nazev, zkratka]
			) faculty
		FROM fak
		UNION
		SELECT kod id, 'en' lang, jsonb_object(
			ARRAY['name', 'abbr'],
			ARRAY[anazev, azkratka]
			) faculty
		FROM fak
	`,
}

/*
Prerequisites:
  - ustav
*/
var ustav2JSON = transformation{
	name: "ustav2json",
	query: `--sql
		DROP TABLE IF EXISTS ustav2json;
		CREATE TABLE ustav2json (
			id VARCHAR(10),
			lang VARCHAR(2),
			department jsonb
		);
		INSERT INTO ustav2json(id, lang, department)
		SELECT kod id, 'cs' lang, jsonb_object(
			ARRAY['id', 'name'],
			ARRAY[kod, nazev]
		) department
		FROM ustav
		UNION
		SELECT kod id, 'en' lang, jsonb_object(
			ARRAY['id', 'name'],
			ARRAY[kod, anazev]
		) department
		FROM ustav
	`,
}

/*
Prerequisites:
  - pvyuc
*/
var pvyuc2lang = transformation{
	name: "pvyuc2lang",
	query: `--sql
		DROP TABLE IF EXISTS pvyuc2lang;
		CREATE TABLE pvyuc2lang (
			id varchar(5),
			lang varchar(2),
			title varchar(120)
		);
		INSERT INTO pvyuc2lang(id, lang, title)
		SELECT kod id, 'cs' lang, nazev title FROM pvyuc
		UNION
		SELECT kod id, 'en' lang, anazev title FROM pvyuc
	`,
}

/*
Prerequisites:
  - jazyk
*/
var jazyk2lang = transformation{
	name: "jazyk2lang",
	query: `--sql
		DROP TABLE IF EXISTS jazyk2lang;
		CREATE TABLE jazyk2lang (
			id varchar(5),
			lang varchar(2),
			title varchar(120)
		);
		INSERT INTO jazyk2lang(id, lang, title)
		SELECT kod id, 'cs' lang, nazev title FROM jazyk
		UNION
		SELECT kod id, 'en' lang, anazev title FROM jazyk
	`,
}

/*
Prerequisites:
  - zsem
*/
var zsem2lang = transformation{
	name: "zsem2lang",
	query: `--sql
		DROP TABLE IF EXISTS zsem2lang;
		CREATE TABLE zsem2lang (
			id varchar(5),
			lang varchar(2),
			title varchar(120)
		);
		INSERT INTO zsem2lang(id, lang, title)
		SELECT kod id, 'cs' lang, nazev title FROM zsem
		UNION
		SELECT kod id, 'en' lang, anazev title FROM zsem
	`,
}

/*
Prerequisites:
  - typypov
*/
var typypov2lang = transformation{
	name: "typypov2lang",
	query: `--sql
		DROP TABLE IF EXISTS typypov2lang;
		CREATE TABLE typypov2lang (
			id VARCHAR(2),
			lang VARCHAR(2),
			title VARCHAR(70),
			exam_winter VARCHAR(15),
			exam_summer VARCHAR(15),
			exam VARCHAR(30)
		);
		INSERT INTO typypov2lang(id, lang, title, exam_winter, exam_summer, exam)
		SELECT
			kod id,
			'cs' lang,
			nazev title,
			exam1 exam_winter,
			exam2 exam_summer,
			(exam1 || COALESCE(' / ' || exam2, '')) exam
		FROM typypov
		UNION
		SELECT
			kod id,
			'en' lang,
			anazev title,
			aexam1 exam_winter,
			aexam2 exam_summer,
			(aexam1 || COALESCE(' / ' || aexam2, '')) exam
		FROM typypov
	`,
}

/*
Prerequisites:
  - ucit
*/
var ucit2JSON = transformation{
	name: "ucit2json",
	query: `--sql
	DROP TABLE IF EXISTS ucit2json;
	CREATE TABLE ucit2json (
		id VARCHAR(10),
		teacher jsonb
	);
	INSERT INTO ucit2json(id, teacher)
	SELECT kod, jsonb_object(
		ARRAY['id', 'first_name', 'last_name', 'title_before', 'title_after'],
		ARRAY[kod, jmeno, prijmeni, titulpred, titulza]
	) teacher
	FROM ucit
	`,
}

/*
Prerequisites:
  - pamela
*/
var pamela2JSON = transformation{
	name: "pamela2json",
	query: `--sql
		DROP TABLE IF EXISTS pamela2json;
		CREATE TABLE pamela2json (
			course_code VARCHAR(10),
			lang VARCHAR(2),
			annotation jsonb,
			syllabus jsonb,
			terms_of_passing jsonb,
			literature jsonb,
			requirements_of_assesment jsonb,
			entry_requirements jsonb,
			aim jsonb
		);
		INSERT INTO pamela2json
		SELECT
			course_code,
			lang,
			CASE
				WHEN annotation IS NULL THEN NULL
				ELSE jsonb_object(ARRAY['title', 'content'], ARRAY[annotation_title, annotation])
			END annotation,
			CASE
				WHEN syllabus IS NULL THEN NULL
				ELSE jsonb_object(ARRAY['title', 'content'], ARRAY[syllabus_title, syllabus])
			END syllabus,
			CASE
				WHEN terms_of_passing IS NULL THEN NULL
				ELSE jsonb_object(ARRAY['title', 'content'], ARRAY[terms_of_passing_title, terms_of_passing])
			END terms_of_passing,
			CASE
				WHEN literature IS NULL THEN NULL
				ELSE jsonb_object(ARRAY['title', 'content'], ARRAY[literature_title, literature])
			END literature,
			CASE
				WHEN requirements_of_assesment IS NULL THEN NULL
				ELSE jsonb_object(ARRAY['title', 'content'], ARRAY[requirements_of_assesment_title, requirements_of_assesment])
			END requirements_of_assesment,
			CASE
				WHEN entry_requirements IS NULL THEN NULL
				ELSE jsonb_object(ARRAY['title', 'content'], ARRAY[entry_requirements_title, entry_requirements])
			END entry_requirements,
			CASE
				WHEN aim IS NULL THEN NULL
				ELSE jsonb_object(ARRAY['title', 'content'], ARRAY[aim_title, aim])
			END aim
		FROM (
			SELECT
				povinn course_code,
				'cs' lang,
				max(nazev) FILTER (WHERE typ='A') annotation_title,
				max(memo) FILTER (WHERE typ='A') annotation,
				max(nazev) FILTER (WHERE typ='S') syllabus_title,
				max(memo) FILTER (WHERE typ='S') syllabus,
				max(nazev) FILTER (WHERE typ='E') terms_of_passing_title,
				max(memo) FILTER (WHERE typ='E') terms_of_passing,
				max(nazev) FILTER (WHERE typ='L') literature_title,
				max(memo) FILTER (WHERE typ='L') literature,
				max(nazev) FILTER (WHERE typ='P') requirements_of_assesment_title,
				max(memo) FILTER (WHERE typ='P') requirements_of_assesment,
				max(nazev) FILTER (WHERE typ='V') entry_requirements_title,
				max(memo) FILTER (WHERE typ='V') entry_requirements,
				max(nazev) FILTER (WHERE typ='C') aim_title,
				max(memo) FILTER (WHERE typ='C') aim
			FROM pamela
			LEFT JOIN typmem ON pamela.typ = typmem.kod
			WHERE jazyk='CZE'
			GROUP BY povinn
			UNION
			SELECT
				povinn course_code,
				'en' lang,
				max(anazev) FILTER (WHERE typ='A') annotation_title,
				max(memo) FILTER (WHERE typ='A') annotation,
				max(anazev) FILTER (WHERE typ='S') syllabus_title,
				max(memo) FILTER (WHERE typ='S') syllabus,
				max(anazev) FILTER (WHERE typ='E') terms_of_passing_title,
				max(memo) FILTER (WHERE typ='E') terms_of_passing,
				max(anazev) FILTER (WHERE typ='L') literature_title,
				max(memo) FILTER (WHERE typ='L') literature,
				max(anazev) FILTER (WHERE typ='P') requirements_of_assesment_title,
				max(memo) FILTER (WHERE typ='P') requirements_of_assesment,
				max(anazev) FILTER (WHERE typ='V') entry_requirements_title,
				max(memo) FILTER (WHERE typ='V') entry_requirements,
				max(anazev) FILTER (WHERE typ='C') aim_title,
				max(memo) FILTER (WHERE typ='C') aim
			FROM pamela
			LEFT JOIN typmem ON pamela.typ = typmem.kod
			WHERE jazyk='ENG'
			GROUP BY povinn
		)
	`,
}

/*
Prerequisites:
  - povinn
  - povinn2jazyk
  - jazyk2lang
*/
var povinn2jazykAgg = transformation{
	name: "povinn2jazyk_agg",
	query: `--sql
		DROP TABLE IF EXISTS povinn2jazyk_agg;
		CREATE TABLE povinn2jazyk_agg (
			course_code VARCHAR(10),
			lang VARCHAR(2),
			taught_lang jsonb,
			taught_lang_id jsonb,
			taught_lang_str VARCHAR(250)
		);
		WITH taught_lang AS (
			SELECT
				povinn,
				jazyk
			FROM povinn2jazyk
			UNION
			SELECT
				povinn,
				pvyjazyk jazyk
			FROM povinn
			WHERE pvyjazyk IS NOT NULL
		)
		INSERT INTO povinn2jazyk_agg
		SELECT
			povinn course_code,
			jl.lang,
			json_agg(jl.title) taught_lang,
			json_agg(jl.id) taught_lang_id,
			string_agg(jl.title, ', ') taught_lang_str
		FROM taught_lang tl
		RIGHT JOIN jazyk2lang jl ON jl.id = tl.jazyk
		GROUP BY povinn, jl.lang
	`,
}

/*
Prerequisites:
  - rvcem
*/
var rvcem2JSON = transformation{
	name: "rvcem2json",
	query: `--sql
	DROP TABLE IF EXISTS rvcem2json;
	CREATE TABLE rvcem2json (
		id VARCHAR(6),
		lang VARCHAR(2),
		range_unit jsonb
	);
	INSERT INTO rvcem2json
	SELECT kod id, 'cs' lang, jsonb_object(
		ARRAY['abbr', 'name'],
		ARRAY[kod, nazev]
	) range_unit
	FROM rvcem
	UNION
	SELECT kod id, 'en' lang, jsonb_object(
		ARRAY['abbr', 'name'],
		ARRAY[kod, anazev]
	) range_unit
	FROM rvcem
`,
}

/*
Prerequisites:
  - obor
*/
var obor2lang = transformation{
	name: "obor2lang",
	query: `--sql
	DROP TABLE IF EXISTS obor2lang;
	CREATE TABLE obor2lang (
		id VARCHAR(12),
		lang VARCHAR(2),
		title VARCHAR(250)
	);
	INSERT INTO obor2lang(id, lang, title)
	SELECT kod id, 'cs' lang, nazev title
	FROM obor
	UNION
	SELECT kod id, 'en' lang, anazev title
	FROM obor
	`,
}

/*
Prerequisites:
  - druh
*/
var druh2JSON = transformation{
	name: "druh2json",
	query: `--sql
	DROP TABLE IF EXISTS druh2json;
	CREATE TABLE druh2json (
		id VARCHAR(12),
		lang VARCHAR(2),
		study_type jsonb
	);
	INSERT INTO druh2json
	SELECT kod id, 'cs' lang, jsonb_strip_nulls(jsonb_object(
		ARRAY['name', 'abbr'],
		ARRAY[nazev, zkratka]
	)) study_type
	FROM druh
	UNION
	SELECT kod id, 'en' lang, jsonb_strip_nulls(jsonb_object(
		ARRAY['name', 'abbr'],
		ARRAY[anazev, azkratka]
	)) study_type
	FROM druh
`,
}

/*
Prerequisites:
  - zsem2lang
  - obor2lang
  - druh2json
  - ucit2json
*/
var ankecy2JSON = transformation{
	name: "ankecy2json",
	query: `--sql
	DROP TABLE IF EXISTS ankecy2json;
	CREATE TABLE ankecy2json (
		course_code VARCHAR(10),
		lang VARCHAR(2),
		survey jsonb
	);
	with ankecy2json as (
	SELECT povinn course_code, zl.lang, jsonb_object(
		ARRAY['academic_year', 'target', 'content', 'semester', 'field', 'study_type', 'target_teacher'],
		ARRAY[sskr, prdmtyp, memo, zl.title, ol.title, dj.study_type, COALESCE(uj.teacher, jsonb_object(ARRAY['first_name'], ARRAY['global']))]::text[]
	) survey
	FROM ankecy a
	LEFT JOIN zsem2lang zl ON a.sem = zl.id::INT
	LEFT JOIN obor2lang ol ON a.sobor = ol.id AND zl.lang = ol.lang
	LEFT JOIN druh2json dj ON a.sdruh = dj.id AND zl.lang = dj.lang
	LEFT JOIN ucit2json uj ON a.ucit = uj.id
	)
	INSERT INTO ankecy2json
	SELECT course_code, lang, jsonb_agg(survey) survey
	FROM ankecy2json
	GROUP BY course_code, lang
	`,
}

/*
Prerequisites:
  - preq
  - pskup
*/
var preq2requisites = transformation{
	name: "preq2requisites",
	query: `--sql
	DROP TABLE IF EXISTS requisites;
	CREATE TABLE requisites (
		target_course VARCHAR(10),
		parent_course VARCHAR(10),
		child_course  VARCHAR(10),
		req_type      VARCHAR(1),
		group_type    VARCHAR(1)
	);
	WITH RECURSIVE req_tree AS (
		SELECT
			p.POVINN       AS target_course,
			p.POVINN       AS parent_course,
			p.REQPOVINN    AS child_course,
			p.REQTYP       AS req_type,
			p.PSKUPINA     AS group_type
		FROM preq p

		UNION ALL

		SELECT
			rt.target_course,
			rt.child_course AS parent_course,
			s.PSPOVINN      AS child_course,
			rt.req_type,
			s.PSKUPINA      AS group_type
		FROM req_tree rt
		JOIN pskup s ON s.POVINN = rt.child_course
	)
	INSERT INTO requisites
	SELECT DISTINCT
		target_course,
		parent_course,
		child_course,
		req_type,
		group_type
	FROM req_tree;
	`,
}

/*
Prerequisites:
  - klas2lang
  - pklas
*/
var pklas2JSON = transformation{
	name: "pklas2json",
	query: `--sql
	DROP TABLE IF EXISTS pklas2json;
	CREATE TABLE pklas2json (
		course_code VARCHAR(10),
		lang VARCHAR(2),
		classifications jsonb
	);
	INSERT INTO pklas2json
	SELECT
		povinn course_code,
		lang,
		jsonb_agg(title) classifications
	FROM pklas
	LEFT JOIN klas2lang ON pklas.pklas = klas2lang.id
	GROUP BY povinn, lang
	`,
}

/*
Prerequisites:
  - ptrida
*/
var ptrida2JSON = transformation{
	name: "ptrida2json",
	query: `--sql
	DROP TABLE IF EXISTS ptrida2json;
	CREATE TABLE ptrida2json (
		course_code VARCHAR(10),
		classes jsonb
	);
	INSERT INTO ptrida2json
	SELECT
		povinn course_code,
		jsonb_agg(nazev) classes
	FROM ptrida
	LEFT JOIN trida ON ptrida.ptrida = trida.kod
	GROUP BY povinn
	`,
}

/*
Prerequisites:
  - ucit_rozvrh
  - ucit2json
  - povinn
  - pamela2json
  - ustav2json
  - fak2json
  - pvyuc2lang
  - zsem2lang
  - povinn2jazyk_agg
  - rvcem2json
  - typypov2lang
  - ankecy2json
  - pklas2json
  - ptrida2json
*/
var povinn2courses = transformation{
	name: "povinn2courses",
	query: `--sql
	DROP TABLE IF EXISTS povinn2courses;
	CREATE TABLE povinn2courses (
		code VARCHAR(10),
		lang VARCHAR(2),
		title VARCHAR(250),
		valid_from INT,
		valid_to INT,
		capacity VARCHAR(10),
		min_occupancy VARCHAR(10),
		credits INT,
		course_url VARCHAR(250),
		lecture_range_winter INT,
		lecture_range_summer INT,
		seminar_range_winter INT,
		seminar_range_summer INT,
		guarantors jsonb,
		teachers jsonb,
		semester_count INT,
		annotation jsonb,
		syllabus jsonb,
		terms_of_passing jsonb,
		literature jsonb,
		requirements_of_assesment jsonb,
		entry_requirements jsonb,
		aim jsonb,
		department jsonb,
		faculty jsonb,
		taught_state_title VARCHAR(120),
		taught_state VARCHAR(1),
		start_semester VARCHAR(5),
		start_semester_title VARCHAR(120),
		taught_lang VARCHAR(250),
		range_unit jsonb,
		exam VARCHAR(30),
		exam_winter VARCHAR(15),
		exam_summer VARCHAR(15),
		survey jsonb,
		classifications jsonb,
		classes jsonb
	);
	WITH distinct_teachers AS (
		SELECT DISTINCT povinn, ucit
		FROM ucit_rozvrh
	), course_teachers AS (
		SELECT ur.povinn, jsonb_agg(uj.teacher) teachers
		FROM distinct_teachers ur
		LEFT JOIN ucit2json uj ON ur.ucit=uj.id
		GROUP BY ur.povinn
	), course_general AS (
		SELECT
			p.povinn code,
			p.vplatiod valid_from,
			p.vplatido valid_to,
			vsemzac start_semester,
			vsempoc semester_count,
			vebody credits,
			purl course_url,
			pgarant guarantor,
			pfakulta faculty,
			pvyucovan taught,
			vrvcem range_unit,
			vtyp exam_type,
			CASE
				WHEN vsemzac='1' THEN vrozsahpr1
				WHEN vsemzac='2' THEN vrozsahpr2
				WHEN vsemzac='3' THEN vrozsahpr1
			END lecture_range_winter,
			CASE
				WHEN vsemzac='1' THEN vrozsahcv1
				WHEN vsemzac='2' THEN vrozsahcv2
				WHEN vsemzac='3' THEN vrozsahcv1
			END seminar_range_winter,
			CASE
				WHEN vsemzac='1' THEN vrozsahpr2
				WHEN vsemzac='2' THEN vrozsahpr1
				WHEN vsemzac='3' THEN vrozsahpr2
			END lecture_range_summer,
			CASE
				WHEN vsemzac='1' THEN vrozsahcv2
				WHEN vsemzac='2' THEN vrozsahcv1
				WHEN vsemzac='3' THEN vrozsahcv2
			END seminar_range_summer,
			to_jsonb(array_remove(ARRAY[u1.teacher, u2.teacher, u3.teacher], NULL)) guarantors,
			ct.teachers
		FROM povinn p
		LEFT JOIN ucit2json u1 ON u1.id = vucit1
		LEFT JOIN ucit2json u2 ON u2.id = vucit2
		LEFT JOIN ucit2json u3 ON u3.id = vucit3
		LEFT JOIN course_teachers ct ON ct.povinn = p.povinn
	), course_lang AS (
	SELECT
		povinn code,
		'cs' lang,
		pnazev title,
		CASE
			WHEN ppocmin IS NULL THEN 'Neomezená'
			ELSE TO_CHAR(ppocmin, '9')
		END min_occupancy,
		CASE
			WHEN ppocmax IS NULL THEN 'Neomezená'
			ELSE TO_CHAR(ppocmax, '9')
		END capacity
	FROM povinn
	UNION
	SELECT
		povinn code,
		'en' lang,
		panazev title,
		CASE
			WHEN ppocmin IS NULL THEN 'Unlimited'
			ELSE TO_CHAR(ppocmin, '9')
		END min_occupancy,
		CASE
			WHEN ppocmax IS NULL THEN 'Unlimited'
			ELSE TO_CHAR(ppocmax, '9')
		END capacity
	FROM povinn
	)
	INSERT INTO povinn2courses
	SELECT
		cl.code,
		cl.lang,
		cl.title,
		cg.valid_from,
		cg.valid_to,
		cl.capacity,
		cl.min_occupancy,
		cg.credits,
		cg.course_url,
		cg.lecture_range_winter,
		cg.lecture_range_summer,
		cg.seminar_range_winter,
		cg.seminar_range_summer,
		cg.guarantors,
		cg.teachers,
		cg.semester_count,
		pj.annotation,
		pj.syllabus,
		pj.terms_of_passing,
		pj.literature,
		pj.requirements_of_assesment,
		pj.entry_requirements,
		pj.aim,
		uj.department,
		fj.faculty,
		pl.title taught_state_title,
		pl.id taught_state,
		zl.id start_semester,
		zl.title start_semester_title,
		pja.taught_lang_str,
		rj.range_unit,
		tl.exam,
		tl.exam_winter,
		tl.exam_summer,
		aj.survey,
		pklas.classifications,
		ptrida.classes
	FROM course_lang cl
	LEFT JOIN course_general cg ON cl.code = cg.code
	LEFT JOIN pamela2json pj ON cl.code = pj.course_code AND cl.lang = pj.lang
	LEFT JOIN ustav2json uj ON cg.guarantor = uj.id AND cl.lang = uj.lang
	LEFT JOIN fak2json fj ON cg.faculty = fj.id AND cl.lang = fj.lang
	LEFT JOIN pvyuc2lang pl ON cg.taught = pl.id AND cl.lang = pl.lang
	LEFT JOIN zsem2lang zl ON cg.start_semester = zl.id AND cl.lang = zl.lang
	LEFT JOIN povinn2jazyk_agg pja ON cg.code = pja.course_code AND cl.lang = pja.lang
	LEFT JOIN rvcem2json rj ON cg.range_unit = rj.id AND cl.lang = rj.lang
	LEFT JOIN typypov2lang tl ON cg.exam_type = tl.id AND cl.lang = tl.lang
	LEFT JOIN ankecy2json aj ON cg.code = aj.course_code AND cl.lang = aj.lang
	LEFT JOIN pklas2json pklas ON cg.code = pklas.course_code AND cl.lang = pklas.lang
	LEFT JOIN ptrida2json ptrida ON cg.code = ptrida.course_code
	`,
}

/*
Prerequisites:
  - povinn2courses
  - povinn
  - povinn2jazyk_agg
  - typypov
*/
var povinn2searchable = transformation{
	name: "povinn2searchable",
	query: `--sql
		DROP TABLE IF EXISTS povinn2searchable;
		CREATE TABLE povinn2searchable (
			id INT,
			code VARCHAR(10),
			credits INT,
			start_semester VARCHAR(5),
			semester_count INT,
			taught_state VARCHAR(1),
			exam JSONB,
			range_unit VARCHAR(2),
			faculty VARCHAR(5),
			department VARCHAR(10),
			section VARCHAR(10),
			taught_lang JSONB,
			lecture_range JSONB,
			seminar_range JSONB,
			guarantors JSONB,
			teachers JSONB,
			title JSONB,
			annotation JSONB,
			syllabus JSONB,
			terms_of_passing JSONB,
			literature JSONB,
			requirements_of_assesment JSONB,
			aim JSONB
		);
			WITH guarantors as (
			SELECT DISTINCT
				code,
				jsonb_agg(
					jsonb_build_object(
					'first_name', guar->'first_name',
					'last_name', guar->'last_name'
					)
				) AS guarantors
			FROM povinn2courses
			, LATERAL jsonb_array_elements(guarantors) AS guar
			GROUP BY code, lang
		), teachers as (
			SELECT DISTINCT
				code,
				jsonb_agg(
					jsonb_build_object(
					'first_name', teach->'first_name',
					'last_name', teach->'last_name'
					)
				) AS teachers
			FROM povinn2courses
			, LATERAL jsonb_array_elements(teachers) AS teach
			GROUP BY code, lang
		), descriptions as (
			SELECT
				code,
				jsonb_build_object(
					'cs', MAX(title) FILTER (WHERE lang = 'cs'),
					'en', MAX(title) FILTER (WHERE lang = 'en')
				) title,
				jsonb_agg(annotation->'content') annotation,
				jsonb_agg(syllabus->'content') syllabus,
				jsonb_agg(terms_of_passing->'content') terms_of_passing,
				jsonb_agg(literature->'content') literature,
				jsonb_agg(requirements_of_assesment->'content') requirements_of_assesment,
				jsonb_agg(aim->'content') aim
			FROM povinn2courses
			GROUP BY code
		)
		INSERT INTO povinn2searchable
		SELECT
			ROW_NUMBER() OVER () AS id,
			pc.code,
			pc.credits,
			pc.start_semester,
			pc.semester_count,
			pc.taught_state,
			jsonb_build_array(exam_winter.kod, exam_summer.kod) exam,
			pc.range_unit->>'abbr' range_unit,
			pvn.pfakulta,
			pvn.pgarant department,
			u.sekce section,
			pja.taught_lang_id,
			jsonb_build_array(pc.lecture_range_summer, pc.lecture_range_winter) lecture_range,
			jsonb_build_array(pc.seminar_range_summer, pc.seminar_range_winter) seminar_range,
			g.guarantors,
			t.teachers,
			d.title,
			d.annotation,
			d.syllabus,
			d.terms_of_passing,
			d.literature,
			d.requirements_of_assesment,
			d.aim
		FROM povinn2courses pc
		LEFT JOIN povinn p ON p.povinn = pc.code
		LEFT JOIN povinn2jazyk_agg pja ON pc.code = pja.course_code AND pc.lang = pja.lang
		LEFT JOIN typypov exam_winter ON SUBSTRING(p.vtyp, 1, 1) = exam_winter.kod
		LEFT JOIN typypov exam_summer ON SUBSTRING(p.vtyp, 2, 1) = exam_summer.kod
		LEFT JOIN guarantors g on pc.code = g.code
		LEFT JOIN teachers t on pc.code = t.code
		LEFT JOIN povinn pvn on pvn.povinn = pc.code
		LEFT JOIN descriptions d on pc.code = d.code
		LEFT JOIN ustav u ON pvn.pgarant = u.kod
		WHERE pc.lang='cs'
		AND credits > 0
		AND (taught_state = 'V' OR taught_state = 'N')
		AND pvn.pgarant NOT LIKE '%STUD'
	`,
}

/*
Prerequisites:
  - klas
*/
var klas2lang = simple2lang("klas", "klas2lang", 6, 60)

func simple2lang(fromTable, toTable string, IDSize, titleSize int) transformation {
	query := `--sql
		DROP TABLE IF EXISTS %s;
		CREATE TABLE %s (
			id VARCHAR(%d),
			lang VARCHAR(2),
			title VARCHAR(%d)
		);
		INSERT INTO %s(id, lang, title)
		SELECT kod id, 'cs' lang, nazev title
		FROM %s
		UNION
		SELECT kod id, 'en' lang, anazev title
		FROM %s
		`
	result := transformation{
		name:  fmt.Sprintf("%s2lang", toTable),
		query: fmt.Sprintf(query, toTable, toTable, IDSize, titleSize, toTable, fromTable, fromTable),
	}
	return result
}

/*
Prerequisites:
*/
var initFilterTables = transformation{
	name: "init_filter_tables",
	query: `--sql
	 	DROP TABLE IF EXISTS filter_values;
		DROP TABLE IF EXISTS filter_categories;
		DROP TABLE IF EXISTS filters;
		CREATE TABLE filters (
			id VARCHAR(50) PRIMARY KEY
		);
		CREATE TABLE filter_categories (
			id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
			filter_id VARCHAR(50) NOT NULL REFERENCES filters(id),
			facet_id VARCHAR(50) NOT NULL,
			title_cs VARCHAR(50) NOT NULL,
			title_en VARCHAR(50) NOT NULL,
			description_cs VARCHAR(200),
			description_en VARCHAR(200),
			condition VARCHAR(100),
			displayed_value_limit INT NOT NULL,
			position INT NOT NULL
		);
		CREATE TABLE filter_values (
			id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
			category_id INT NOT NULL REFERENCES filter_categories(id),
			facet_id VARCHAR(50) NOT NULL,
			title_cs VARCHAR(250) NOT NULL,
			title_en VARCHAR(250) NOT NULL,
			description_cs VARCHAR(200),
			description_en VARCHAR(200),
			position INT NOT NULL
		);
		INSERT INTO filters(id) VALUES
		('courses'),
		('course-survey'),
		('degree-plans');
		INSERT INTO filter_categories(filter_id, facet_id, title_cs, title_en, description_cs, description_en, condition, displayed_value_limit, position)
		VALUES ` +
		categoriesToSQL("courses", []category{
			makeCategory("faculty", "Fakulta", "Faculty", 3),
			makeCategory("taught_state", "Stav předmětu", "Course status", 3),
			makeCategory("start_semester", "Semestr", "Semester", 3),
			makeCategory("credits", "Kredity", "Credits", 6),
			makeCategory("semester_count", "Počet semestrů", "Number of semesters", 3),
			makeCategory("lecture_range", "Rozsah přednášky", "Lecture range", 3).withDescription(
				"Rozsah přednášky pro letní nebo zimní semestr.",
				"Lecture range for summer or winter semester.",
			),
			makeCategory("seminar_range", "Rozsah cvičení", "Seminar range", 3).withDescription(
				"Rozsah cvičení pro letní nebo zimní semestr.",
				"Seminar range for summer or winter semester.",
			),
			makeCategory("taught_lang", "Jazyk výuky", "Language", 3),
			makeCategory("exam", "Typ examinace", "Exam type", 3),
			makeCategory("range_unit", "Jednotka rozsahu", "Range unit", 3).withDescription(
				"Jednotka rozsahu přednášky a cvičení.",
				"Unit of lecture and seminar range.",
			),
			makeCategory("section", "Sekce", "Section", 3),
			makeCategory("department", "Katedra", "Department", 3),
		}) + `,` +
		categoriesToSQL("course-survey", []category{
			makeCategory("teacher.id", "Učitelé", "Teachers", 5),
			makeCategory("academic_year", "Rok", "Year", 5),
			makeCategory("study_field.id", "Obor", "Field", 5),
			makeCategory("study_type.id", "Forma studia", "Study form", 5),
			makeCategory("study_year", "Ročník", "Year of study", 5),
			makeCategory("target_type", "Přednáška/Cvičení", "Lecture/Seminar", 5),
		}) + `,` +
		categoriesToSQL("degree-plans", []category{
			makeCategory("faculty", "Fakulta", "Faculty", 1),
			makeCategory("section", "Sekce", "Section", 3),
			makeCategory("study_type", "Forma studia", "Study form", 3),
			makeCategory("teaching_lang", "Jazyk výuky", "Teaching language", 2),
			makeCategory("validity", "Platí v roce", "Valid in year", 30).withCondition("validity.from <= {VAL} AND validity.to >= {VAL}"),
			makeCategory("field.code", "Obor", "Field", 5),
		}),
}

type category struct {
	facetID             string
	titleCS             string
	titleEN             string
	descriptionCS       string
	descriptionEN       string
	condition           string
	displayedValueLimit int
}

func makeCategory(facetID, titleCS, titleEN string, displayedValueLimit int) category {
	return category{
		facetID:             facetID,
		titleCS:             titleCS,
		titleEN:             titleEN,
		displayedValueLimit: displayedValueLimit,
	}
}

func (c category) withDescription(descriptionCS, descriptionEN string) category {
	c.descriptionCS = descriptionCS
	c.descriptionEN = descriptionEN
	return c
}

func (c category) withCondition(condition string) category {
	c.condition = condition
	return c
}

func (c category) ToSQL(categoryID string, position int) string {
	var dcs string
	if c.descriptionCS == "" {
		dcs = "NULL"
	} else {
		dcs = fmt.Sprintf("'%s'", c.descriptionCS)
	}
	var den string
	if c.descriptionEN == "" {
		den = "NULL"
	} else {
		den = fmt.Sprintf("'%s'", c.descriptionEN)
	}
	var condition string
	if c.condition == "" {
		condition = "NULL"
	} else {
		condition = fmt.Sprintf("'%s'", c.condition)
	}
	return fmt.Sprintf(
		"('%s', '%s', '%s', '%s', %s, %s, %s, %d, %d)",
		categoryID, c.facetID, c.titleCS, c.titleEN, dcs, den, condition, c.displayedValueLimit, position,
	)
}

func categoriesToSQL(categoryID string, categories []category) string {
	var cstrs = make([]string, len(categories))
	for i, c := range categories {
		cstrs[i] = c.ToSQL(categoryID, i+1)
	}
	return strings.Join(cstrs, ",\n")
}

/*
Prerequisites:
  - povinn2searchable
*/
var createFilterValuesForCredits = transformation{
	name: "create_filter_values_for_credits",
	query: `--sql
		WITH category_id AS (
			SELECT fc.id FROM filters f
			LEFT JOIN filter_categories fc ON f.id = fc.filter_id
			WHERE f.id='courses'
			AND fc.facet_id='credits'
		), distinct_credits AS (
			SELECT DISTINCT credits
			FROM povinn2searchable
			WHERE credits IS NOT NULL
		)
		INSERT INTO filter_values (category_id, facet_id, title_cs, title_en, position)
		SELECT
			cid.id,
			credits facet_id,
			credits title_cs,
			credits title_en,
			ROW_NUMBER() OVER (ORDER BY credits) position
		FROM distinct_credits dc
		LEFT JOIN category_id cid ON true
	`,
}

/*
Prerequisites:
  - povinn2searchable
  - povinn2jazyk
  - povinn
*/
var createFilterValuesForLangs = transformation{
	name: "create_filter_values_for_langs",
	query: `--sql
		WITH langs AS (
			SELECT
				povinn,
				jazyk
			FROM povinn2jazyk
			UNION
			SELECT
				povinn,
				pvyjazyk jazyk
			FROM povinn
		), ORDER_LIST(lang, position) AS (
			VALUES
				('CZE', 1),
				('ENG', 2)
		), distinct_langs AS (
			SELECT DISTINCT ON (jazyk)
				l.jazyk,
				j.nazev,
				j.anazev
			FROM povinn2searchable ps
			LEFT JOIN langs l ON ps.code = l.povinn
			LEFT JOIN jazyk j ON j.kod = l.jazyk
			WHERE l.jazyk IS NOT NULL
		), category_id AS (
			SELECT fc.id FROM filters f
			LEFT JOIN filter_categories fc ON f.id = fc.filter_id
			WHERE f.id='courses'
			AND fc.facet_id='taught_lang'
		)
		INSERT INTO filter_values (category_id, facet_id, title_cs, title_en, position)
		SELECT
			cid.id category_id,
			jazyk facet_id,
			nazev title_cs,
			anazev title_en,
			COALESCE(position, ROW_NUMBER() OVER ()) position
		FROM distinct_langs dl
		LEFT JOIN order_list o ON dl.jazyk = o.lang
		LEFT JOIN category_id cid ON true
		ORDER BY position, jazyk
`,
}

/*
Prerequisites:
  - povinn2searchable
  - fak
*/
var createFilterValuesForFaculties = transformation{
	name: "create_filter_values_for_faculties",
	query: `--sql
		with category_id AS (
			SELECT fc.id FROM filters f
			LEFT JOIN filter_categories fc ON f.id = fc.filter_id
			WHERE f.id='courses'
			AND fc.facet_id='faculty'
		)
		INSERT INTO filter_values (category_id, facet_id, title_cs, title_en, description_cs, description_en, position)
		SELECT DISTINCT ON (faculty)
			cid.id category_id,
			kod facet_id,
			zkratka title_cs,
			azkratka title_en,
			nazev description_cs,
			anazev description_en,
			ROW_NUMBER() OVER () position
		FROM povinn2searchable
		LEFT JOIN fak ON faculty = fak.kod
		LEFT JOIN category_id cid ON true
`,
}

/*
Prerequisites:
  - povinn2searchable
  - pvyuc
*/
var createFilterValuesForTaughtStates = transformation{
	name: "create_filter_values_for_taught_states",
	query: `--sql
		WITH category_id AS (
			SELECT fc.id FROM filters f
			LEFT JOIN filter_categories fc ON f.id = fc.filter_id
			WHERE f.id='courses'
			AND fc.facet_id='taught_state'
		), distinct_states AS (
			SELECT DISTINCT taught_state
			FROM povinn2searchable
		)
		INSERT INTO filter_values (category_id, facet_id, title_cs, title_en, position)
		SELECT
			cid.id category_id,
			taught_state facet_id,
			pv.nazev title_cs,
			pv.anazev title_en,
			ROW_NUMBER() OVER (ORDER BY taught_state DESC) position
		FROM distinct_states ds
		LEFT JOIN pvyuc pv ON ds.taught_state = pv.kod
		LEFT JOIN category_id cid ON true
	`,
}

/*
Prerequisites:
  - zsem
*/
var createFilterValuesForStartSemesters = transformation{
	name: "create_filter_values_for_start_semesters",
	query: `--sql
		WITH category_id AS (
			SELECT fc.id FROM filters f
			LEFT JOIN filter_categories fc ON f.id = fc.filter_id
			WHERE f.id='courses'
			AND fc.facet_id='start_semester'
		)
		INSERT INTO filter_values (category_id, facet_id, title_cs, title_en, position)
		SELECT
			cid.id category_id,
			kod facet_id,
			nazev title_cs,
			anazev title_en,
			ROW_NUMBER() OVER (ORDER BY kod) position
		FROM zsem
		LEFT JOIN category_id cid ON true
	`,
}

/*
Prerequisites:
  - povinn2searchable
*/
var createFilterValuesForSemesterCounts = transformation{
	name: "create_filter_values_for_semester_counts",
	query: `--sql
		WITH category_id AS (
			SELECT fc.id FROM filters f
			LEFT JOIN filter_categories fc ON f.id = fc.filter_id
			WHERE f.id='courses'
			AND fc.facet_id='semester_count'
		), distinct_semester_count AS (
			SELECT DISTINCT semester_count
			FROM povinn2searchable
		)
		INSERT INTO filter_values (category_id, facet_id, title_cs, title_en, position)
		SELECT
			cid.id category_id,
			semester_count facet_id,
			semester_count title_cs,
			semester_count title_en,
			ROW_NUMBER() OVER (ORDER BY semester_count) position
		FROM distinct_semester_count ps
		LEFT JOIN category_id cid ON true
	`,
}

/*
Prerequisites:
  - povinn2searchable
*/
var createFilterValuesForLectureRanges = transformation{
	name: "create_filter_values_for_lecture_ranges",
	query: `--sql
		WITH category_id AS (
			SELECT fc.id FROM filters f
			LEFT JOIN filter_categories fc ON f.id = fc.filter_id
			WHERE f.id='courses'
			AND fc.facet_id='lecture_range'
		), lecture_ranges AS (
			SELECT
				(lecture_range->>0)::INT lecture_range
			FROM povinn2searchable ps
			UNION
			SELECT
				(lecture_range->>1)::INT lecture_range
			FROM povinn2searchable
		)
		INSERT INTO filter_values (category_id, facet_id, title_cs, title_en, position)
		SELECT DISTINCT ON (lecture_range)
			cid.id category_id,
			lecture_range facet_id,
			lecture_range title_cs,
			lecture_range title_en,
			ROW_NUMBER() OVER (ORDER BY lecture_range) position
		FROM lecture_ranges
		LEFT JOIN category_id cid ON true
		WHERE lecture_range IS NOT NULL
	`,
}

/*
Prerequisites:
  - povinn2searchable
*/
var createFilterValuesForSeminarRanges = transformation{
	name: "create_filter_values_for_seminar_ranges",
	query: `--sql
		WITH category_id AS (
			SELECT fc.id FROM filters f
			LEFT JOIN filter_categories fc ON f.id = fc.filter_id
			WHERE f.id='courses'
			AND fc.facet_id='seminar_range'
		), seminar_ranges AS (
			SELECT
				(seminar_range->>0)::INT seminar_range
			FROM povinn2searchable ps
			UNION
			SELECT
				(seminar_range->>1)::INT seminar_range
			FROM povinn2searchable
		)
		INSERT INTO filter_values (category_id, facet_id, title_cs, title_en, position)
		SELECT DISTINCT ON (seminar_range)
			cid.id category_id,
			seminar_range facet_id,
			seminar_range title_cs,
			seminar_range title_en,
			ROW_NUMBER() OVER (ORDER BY seminar_range) position
		FROM seminar_ranges
		LEFT JOIN category_id cid ON true
		WHERE seminar_range IS NOT NULL
	`,
}

/*
Prerequisites:
  - povinn2searchable
  - povinn
  - typypov
*/
var createFilterValuesForExams = transformation{
	name: "create_filter_values_for_exams",
	query: `--sql
		WITH category_id AS (
			SELECT
				fc.id
			FROM filters f
			LEFT JOIN filter_categories fc ON f.id = fc.filter_id
			WHERE f.id='courses'
			AND fc.facet_id='exam'
		), ORDER_LIST(exam, position) AS (
			VALUES
				('*', 1),
				('K', 2),
				('Z', 3),
				('F', 4)
		), exam_codes AS (
			SELECT DISTINCT
				t.kod
			FROM povinn2searchable ps
			LEFT JOIN povinn p ON ps.code = p.povinn
			LEFT join typypov t ON p.vtyp = t.kod
		), exams AS (
			SELECT
				SUBSTRING(kod, 1, 1) exam
			FROM exam_codes
			UNION
			SELECT
				SUBSTRING(kod, 2, 1) exam
			FROM exam_codes
		), distinct_exams AS (
		SELECT DISTINCT
			exam
		FROM exams
		)
		INSERT INTO filter_values (category_id, facet_id, title_cs, title_en, description_cs, description_en, position)
		SELECT
			cid.id,
			t.kod facet_id,
			t.exam1 title_cs,
			t.aexam1 title_en,
			t.nazev description_cs,
			t.anazev description_en,
			ROW_NUMBER() OVER (ORDER BY ol.position) position
		FROM distinct_exams e
		LEFT JOIN typypov t ON t.kod = e.exam
		LEFT JOIN ORDER_LIST ol ON ol.exam = e.exam
		LEFT JOIN category_id cid ON true
		WHERE e.exam IS NOT NULL
		AND LENGTH(e.exam) > 0
	`,
}

/*
Prerequisites:
  - povinn2searchable
  - rvcem
*/
var createFilterValuesForRangeUnits = transformation{
	name: "create_filter_values_for_range_units",
	query: `--sql
	WITH category_id AS (
		SELECT
			fc.id
		FROM filters f
		LEFT JOIN filter_categories fc ON f.id = fc.filter_id
		WHERE f.id='courses'
		AND fc.facet_id='range_unit'
	), distinct_range_units AS (
		SELECT DISTINCT
			range_unit
		FROM povinn2searchable
		WHERE range_unit IS NOT NULL
	)
	INSERT INTO filter_values (category_id, facet_id, title_cs, title_en, description_cs, description_en, position)
	SELECT
		cid.id,
		r.kod facet_id,
		r.kod title_cs,
		r.kod title_en,
		r.nazev description_cs,
		r.anazev description_en,
		ROW_NUMBER() OVER () position
	FROM distinct_range_units dru
	LEFT JOIN rvcem r ON dru.range_unit = r.kod
	LEFT JOIN category_id cid ON true
	`,
}

/*
Prerequisites:
  - povinn2searchable
  - ustav
*/
var createFilterValuesForDepartments = transformation{
	name: "create_filter_values_for_departments",
	query: `--sql
		WITH category_id AS (
			SELECT
				fc.id
			FROM filters f
			LEFT JOIN filter_categories fc ON f.id = fc.filter_id
			WHERE f.id='courses'
			AND fc.facet_id='department'
		), distinct_departments AS (
			SELECT DISTINCT
				department
			FROM povinn2searchable
		)
		INSERT INTO filter_values (category_id, facet_id, title_cs, title_en, description_cs, description_en, position)
		SELECT
			cid.id category_id,
			u.kod facet_id,
			u.kod title_cs,
			u.kod title_en,
			u.nazev description_cs,
			u.anazev description_en,
			ROW_NUMBER() OVER (ORDER BY u.sekce, u.kod) position
		FROM distinct_departments d
		LEFT JOIN ustav u ON d.department=u.kod
		LEFT JOIN category_id cid ON true
	`,
}

/*
Prerequisites:
  - povinn2searchable
  - sekce
  - ustav
*/
var createFilterValuesForSections = transformation{
	name: "create_filter_values_for_sections",
	query: `--sql
		WITH category_id AS (
			SELECT
				fc.id
			FROM filters f
			LEFT JOIN filter_categories fc ON f.id = fc.filter_id
			WHERE f.id='courses'
			AND fc.facet_id='section'
		), order_list(section, position) AS (
			VALUES
				('NI', 1),
				('NM', 2),
				('NF', 3)
		), distinct_sections AS (
			SELECT DISTINCT
				s.kod,
				s.nazev,
				s.anazev
			FROM povinn2searchable ps
			LEFT JOIN ustav u ON ps.department = u.kod
			LEFT JOIN sekce s ON u.sekce = s.kod
		)
		INSERT INTO filter_values (category_id, facet_id, title_cs, title_en, position)
		SELECT
			cid.id category_id,
			kod facet_id,
			nazev title_cs,
			anazev title_en,
			ROW_NUMBER() OVER (ORDER BY ol.position) position
		FROM distinct_sections s
		LEFT JOIN order_list ol on ol.section = s.kod
		LEFT JOIN category_id cid ON true
	`,
}

/*
Prerequisites:
  - povinn2searchable
  - ankecy
  - ucit
*/
var createFilterValuesForSurveyTeachers = transformation{
	name: "create_filter_values_for_survey_teachers",
	query: `--sql
			WITH category_id AS (
			SELECT
				fc.id
			FROM filters f
			LEFT JOIN filter_categories fc ON f.id = fc.filter_id
			WHERE f.id='course-survey'
			AND fc.facet_id='teacher.id'
		), distinct_teachers AS (
			SELECT DISTINCT ON (u.kod)
			    u.kod,
				u.jmeno,
				u.prijmeni
			FROM povinn2searchable ps
			LEFT JOIN ankecy a ON ps.code = a.povinn
			LEFT JOIN ucit u ON u.kod = a.ucit
			WHERE u.kod IS NOT NULL
		)
		INSERT INTO filter_values (category_id, facet_id, title_cs, title_en, position)
		SELECT
			cid.id category_id,
			kod facet_id,
			CONCAT(jmeno, ' ', prijmeni) title_cs,
			CONCAT(jmeno, ' ', prijmeni) title_en,
			ROW_NUMBER() OVER (ORDER BY CONCAT(jmeno, prijmeni)) position
		FROM distinct_teachers
		LEFT JOIN category_id cid ON true
	`,
}

/*
Prerequisites:
  - povinn2searchable
  - ucit2json
  - ankecy
  - obor
  - druh
*/
var ankecy2searchable = transformation{
	name: "ankecy2searchable",
	query: `--sql
		DROP TABLE IF EXISTS ankecy2searchable;
		CREATE TABLE ankecy2searchable (
			id INT,
			course_code VARCHAR(10) NOT NULL,
			study_year INT,
			academic_year INT,
			study_field JSONB,
			study_type JSONB,
			teacher JSONB,
			target_type VARCHAR(30),
			content TEXT
		);
		INSERT INTO ankecy2searchable
		SELECT
			ROW_NUMBER() OVER () id,
			a.povinn course_code,
			a.sroc study_year,
			a.sskr academic_year,
			JSONB_BUILD_OBJECT(
				'id', o.kod,
				'name', JSONB_BUILD_OBJECT(
					'cs', o.nazev,
					'en', o.anazev
				)
			) study_field,
			JSONB_BUILD_OBJECT(
				'id', d.kod,
				'abbr', JSONB_BUILD_OBJECT(
					'cs', COALESCE(d.zkratka, d.nazev),
					'en', COALESCE(d.azkratka, d.anazev)
				),
				'name', JSONB_BUILD_OBJECT(
					'cs', d.nazev,
					'en', d.anazev
				)
			) study_type,
			uj.teacher,
			a.prdmtyp target_type,
			a.memo content
		FROM povinn2searchable ps
		LEFT JOIN ankecy a ON a.povinn = ps.code
		LEFT JOIN obor o ON a.sobor = o.kod
		LEFT JOIN druh d ON a.sdruh = d.kod
		LEFT JOIN ucit2json uj ON a.ucit = uj.id
		WHERE a.povinn IS NOT NULL
	`,
}

var createFilterValuesForSurveyAcademicYears = transformation{
	name: "create_filter_values_for_survey_academic_years",
	query: `--sql
		WITH category_id AS (
			SELECT
				fc.id
			FROM filters f
			LEFT JOIN filter_categories fc ON f.id = fc.filter_id
			WHERE f.id='course-survey'
			AND fc.facet_id='academic_year'
		), distinct_academic_year AS (
			SELECT DISTINCT ON (a.sskr)
				a.sskr
			FROM povinn2searchable ps
			LEFT JOIN ankecy a ON ps.code = a.povinn
			WHERE a.sskr IS NOT NULL
		)
		INSERT INTO filter_values (category_id, facet_id, title_cs, title_en, position)
		SELECT
			cid.id category_id,
			sskr facet_id,
			sskr title_cs,
			sskr title_en,
			ROW_NUMBER() OVER (ORDER BY sskr) position
		FROM distinct_academic_year
		LEFT JOIN category_id cid ON true
		`,
}

/*
Prerequisites:
  - povinn2searchable
  - ankecy
  - obor
*/
var createFilterValuesForSurveyStudyFields = transformation{
	name: "create_filter_values_for_survey_study_fields",
	query: `--sql
		WITH category_id AS (
			SELECT
				fc.id
			FROM filters f
			LEFT JOIN filter_categories fc ON f.id = fc.filter_id
			WHERE f.id='course-survey'
			AND fc.facet_id='study_field.id'
		), distinct_study_field AS (
			SELECT DISTINCT ON (a.sobor)
				a.sobor
			FROM povinn2searchable ps
			LEFT JOIN ankecy a ON ps.code = a.povinn
			WHERE a.sobor IS NOT NULL
		)
		INSERT INTO filter_values (category_id, facet_id, title_cs, title_en, position)
		SELECT
			cid.id category_id,
			o.kod facet_id,
			o.nazev title_cs,
			COALESCE(o.anazev, o.nazev) title_en,
			ROW_NUMBER() OVER (ORDER BY o.nazev) position
		FROM distinct_study_field d
		LEFT JOIN obor o ON d.sobor = o.kod
		LEFT JOIN category_id cid ON true
	`,
}

/*
Prerequisites:
  - povinn2searchable
  - ankecy
  - druh
*/
var createFilterValuesForSurveyStudyTypes = transformation{
	name: "create_filter_values_for_survey_study_types",
	query: `--sql
		WITH category_id AS (
			SELECT
				fc.id
			FROM filters f
			LEFT JOIN filter_categories fc ON f.id = fc.filter_id
			WHERE f.id='course-survey'
			AND fc.facet_id='study_type.id'
		), distinct_study_type AS (
			SELECT DISTINCT ON (a.sdruh)
				a.sdruh
			FROM povinn2searchable ps
			LEFT JOIN ankecy a ON ps.code = a.povinn
			WHERE a.sdruh IS NOT NULL
		)
		INSERT INTO filter_values (category_id, facet_id, title_cs, title_en, position)
		SELECT
			cid.id category_id,
			d.kod facet_id,
			COALESCE(d.zkratka, d.nazev) title_cs,
			COALESCE(d.azkratka, d.anazev) title_en,
			ROW_NUMBER() OVER (ORDER BY d.nazev) position
		FROM distinct_study_type dst
		LEFT JOIN druh d ON dst.sdruh = d.kod
		LEFT JOIN category_id cid ON true
		`,
}

/*
Prerequisites:
  - povinn2searchable
  - ankecy
*/
var createFilterValuesForSurveyStudyYears = transformation{
	name: "create_filter_values_for_survey_study_years",
	query: `--sql
		WITH category_id AS (
			SELECT
				fc.id
			FROM filters f
			LEFT JOIN filter_categories fc ON f.id = fc.filter_id
			WHERE f.id='course-survey'
			AND fc.facet_id='study_year'
		), distinct_study_year AS (
			SELECT DISTINCT ON (a.sroc)
				a.sroc
			FROM povinn2searchable ps
			LEFT JOIN ankecy a ON ps.code = a.povinn
			WHERE a.sroc IS NOT NULL
		)
		INSERT INTO filter_values (category_id, facet_id, title_cs, title_en, position)
		SELECT
			cid.id category_id,
			d.sroc facet_id,
			d.sroc title_cs,
			d.sroc title_en,
			ROW_NUMBER() OVER (ORDER BY d.sroc) position
		FROM distinct_study_year d
		LEFT JOIN category_id cid ON true
	`,
}

/*
Prerequisites:
  - povinn2searchable
  - ankecy
*/
var createFilterValuesForSurveyTargetTypes = transformation{
	name: "create_filter_values_for_survey_target_types",
	query: `--sql
		WITH category_id AS (
			SELECT
				fc.id
			FROM filters f
			LEFT JOIN filter_categories fc ON f.id = fc.filter_id
			WHERE f.id='course-survey'
			AND fc.facet_id='target_type'
		), distinct_target_type AS (
			SELECT DISTINCT
				a.prdmtyp
			FROM povinn2searchable ps
			LEFT JOIN ankecy a ON ps.code = a.povinn
			WHERE a.prdmtyp IS NOT NULL
		)
		INSERT INTO filter_values (category_id, facet_id, title_cs, title_en, position)
		SELECT
			cid.id category_id,
			d.prdmtyp facet_id,
			d.prdmtyp title_cs,
			d.prdmtyp title_en,
			ROW_NUMBER() OVER (ORDER BY d.prdmtyp) position
		FROM distinct_target_type d
		LEFT JOIN category_id cid ON true
		`,
}

/*
Prerequisites:
  - stud_plan
*/
var studplan2lang = transformation{
	name: "studplan2lang",
	query: `--sql
		DROP TABLE IF EXISTS studplan2lang;
		CREATE TABLE studplan2lang (
			plan_code VARCHAR(15) NOT NULL,
			lang VARCHAR(2) NOT NULL,
			course_code VARCHAR(10) NOT NULL,
			interchangeability VARCHAR(10),
			recommended_year_from INT,
			recommended_year_to INT,
			recommended_semester INT,
			bloc_name VARCHAR(250),
			bloc_subject_code VARCHAR(20),
			bloc_type VARCHAR(1),
			bloc_limit INT,
			seq VARCHAR(50)
		);
		INSERT INTO studplan2lang
		SELECT
			plan_code,
			'cs' lang,
			code course_code,
			interchangeability,
			recommended_year_from,
			recommended_year_to,
			recommended_semester,
			bloc_name_cz bloc_name,
			bloc_subject_code,
			bloc_type,
			bloc_limit,
			seq
		FROM stud_plan
		UNION
		SELECT
			plan_code,
			'en' lang,
			code course_code,
			interchangeability,
			recommended_year_from,
			recommended_year_to,
			recommended_semester,
			bloc_name_en bloc_name,
			bloc_subject_code,
			bloc_type,
			bloc_limit,
			seq
		FROM stud_plan
	`,
}

/*
Prerequisites:
  - stud_plan_metadata
*/
var studmetadata2lang = transformation{
	name: "studmetadata2lang",
	query: `--sql
		DROP TABLE IF EXISTS studmetadata2lang;
		CREATE TABLE studmetadata2lang (
			plan_code VARCHAR(15) NOT NULL,
			lang VARCHAR(2) NOT NULL,
			title VARCHAR(250),
			valid_from INT,
			valid_to INT,
			faculty VARCHAR(5),
			section VARCHAR(2),
			field_code VARCHAR(20),
			study_type VARCHAR(5)
		);
		INSERT INTO studmetadata2lang
		SELECT
			code as plan_code,
			'cs' lang,
			name_cz as title,
			valid_from,
			valid_to,
			faculty,
			section,
			field_code,
			druh.zkratka AS study_type
		FROM stud_plan_metadata
		JOIN druh ON stud_plan_metadata.study_type = druh.kod
		UNION
		SELECT
			code as plan_code,
			'en' lang,
			name_en as title,
			valid_from,
			valid_to,
			faculty,
			section,
			field_code,
			druh.zkratka AS study_type
		FROM stud_plan_metadata
		JOIN druh ON stud_plan_metadata.study_type = druh.kod
	`,
}

/*
Prerequisites:
  - stud_plan_obor
*/
var studobor2lang = transformation{
	name: "studobor2lang",
	query: `--sql
		DROP TABLE IF EXISTS studobor2lang;
		CREATE TABLE studobor2lang (
			code VARCHAR(15),
			lang VARCHAR(2),
			title VARCHAR(250),
			teaching_lang CHAR(3),
			sims_code VARCHAR(20),
			sims_title VARCHAR(250)
		);
		INSERT INTO studobor2lang
		SELECT
			code,
			'cs' lang,
			name_cz as title,
			teaching_lang,
			sims_code,
			sims_name_cz sims_title
		FROM stud_plan_obor
		UNION
		SELECT
			code,
			'en' lang,
			name_en as title,
			teaching_lang,
			sims_code,
			sims_name_en sims_title
		FROM stud_plan_obor
	`,
}

/*
Prerequisites:
  - stud_plan_metadata
  - stud_plan_obor
*/
var studplan2searchable = transformation{
	name: "studplan2searchable",
	query: `--sql
		DROP TABLE IF EXISTS studplan2searchable;
		CREATE TABLE studplan2searchable (
			id INT PRIMARY KEY,
			code VARCHAR(15),
			title jsonb,
			faculty VARCHAR(5),
			section VARCHAR(2),
			field jsonb,
			study_type VARCHAR(1),
			validity jsonb,
			teaching_lang CHAR(3)
		);
		INSERT INTO studplan2searchable
		SELECT
			ROW_NUMBER() OVER () id,
			smd.code,
			JSONB_BUILD_OBJECT(
				'cs', smd.name_cz,
				'en', smd.name_en
			) as title,
			smd.faculty,
			smd.section,
			JSONB_BUILD_OBJECT(
				'code', sob.code,
				'title', JSONB_BUILD_OBJECT(
					'cs', sob.name_cz,
					'en', sob.name_en
				),
				'sims_code', sob.sims_code,
				'sims_title', JSONB_BUILD_OBJECT(
					'cs', sob.sims_name_cz,
					'en', sob.sims_name_en
				)
			) as field,
			smd.study_type,
			JSONB_BUILD_OBJECT(
				'from', smd.valid_from,
				'to', smd.valid_to
			) as validity,
			sob.teaching_lang
		FROM stud_plan_metadata smd
		LEFT JOIN stud_plan_obor sob ON smd.field_code = sob.code
	`,
}

/*
Prerequisites:
  - studplan2searchable
  - fak
*/
var createFilterValuesForDegreePlanFaculties = transformation{
	name: "create_filter_values_for_degree_plan_faculties",
	query: `--sql
		with category_id AS (
			SELECT fc.id FROM filters f
			LEFT JOIN filter_categories fc ON f.id = fc.filter_id
			WHERE f.id='degree-plans'
			AND fc.facet_id='faculty'
		)
		INSERT INTO filter_values (category_id, facet_id, title_cs, title_en, description_cs, description_en, position)
		SELECT DISTINCT ON (faculty)
			cid.id category_id,
			kod facet_id,
			zkratka title_cs,
			azkratka title_en,
			nazev description_cs,
			anazev description_en,
			ROW_NUMBER() OVER () position
		FROM studplan2searchable
		LEFT JOIN fak ON faculty = fak.kod
		LEFT JOIN category_id cid ON true
	`,
}

/*
Prerequisites:
  - studplan2searchable
  - sekce
*/
var createFilterValuesForDegreePlanSections = transformation{
	name: "create_filter_values_for_degree_plan_sections",
	query: `--sql
		WITH category_id AS (
			SELECT fc.id FROM filters f
			LEFT JOIN filter_categories fc ON f.id = fc.filter_id
			WHERE f.id='degree-plans'
			AND fc.facet_id='section'
		), order_list(section, position) AS (
			VALUES
				('NI', 1),
				('NM', 2),
				('NF', 3)
		), distinct_sections AS (
			SELECT DISTINCT
				section
			FROM studplan2searchable
		)
		INSERT INTO filter_values (category_id, facet_id, title_cs, title_en, position)
		SELECT
			cid.id category_id,
			kod facet_id,
			nazev title_cs,
			anazev title_en,
			ROW_NUMBER() OVER (ORDER BY ol.position) position
		FROM distinct_sections s
		LEFT JOIN sekce sec ON s.section = sec.kod
		LEFT JOIN order_list ol on ol.section = s.section
		LEFT JOIN category_id cid ON true
	`,
}

/*
Prerequisites:
  - stud_plan_obor
*/
var createFilterValuesForDegreePlanFields = transformation{
	name: "create_filter_values_for_degree_plan_fields",
	query: `--sql
		WITH category_id AS (
			SELECT fc.id FROM filters f
			LEFT JOIN filter_categories fc ON f.id = fc.filter_id
			WHERE f.id='degree-plans'
			AND fc.facet_id='field.code'
		)
		INSERT INTO filter_values (category_id, facet_id, title_cs, title_en, position)
		SELECT
			cid.id category_id,
			code facet_id,
			CONCAT(name_cz, ' (', code, ')') title_cs,
			CONCAT(name_en, ' (', code, ')') title_en,
			ROW_NUMBER() OVER (ORDER BY name_cz) position
		FROM stud_plan_obor
		LEFT JOIN category_id cid ON true
	`,
}

/*
Prerequisites:
  - stud_plan_obor
*/
var createFilterValuesForDegreePlanLanguages = transformation{
	name: "create_filter_values_for_degree_plan_languages",
	query: `--sql
		WITH category_id AS (
			SELECT fc.id FROM filters f
			LEFT JOIN filter_categories fc ON f.id = fc.filter_id
			WHERE f.id='degree-plans'
			AND fc.facet_id='teaching_lang'
		),
		distint_teachin_lang AS (
			SELECT DISTINCT
				teaching_lang
			FROM stud_plan_obor
		)
		INSERT INTO filter_values (category_id, facet_id, title_cs, title_en, position)
		SELECT
			cid.id category_id,
			l.teaching_lang facet_id,
			j.nazev title_cs,
			COALESCE(j.anazev, j.nazev) title_en,
			ROW_NUMBER() OVER (ORDER BY teaching_lang) position
		FROM distint_teachin_lang l
		LEFT JOIN jazyk j ON l.teaching_lang = j.kod
		LEFT JOIN category_id cid ON true
	`,
}

/*
Prerequisites:
  - studplan2searchable
*/
var createFilterValuesForDegreePlanValid = transformation{
	name: "create_filter_values_for_degree_plan_valid",
	query: `--sql
		WITH category_id AS (
			SELECT fc.id FROM filters f
			LEFT JOIN filter_categories fc ON f.id = fc.filter_id
			WHERE f.id='degree-plans'
			AND fc.facet_id='validity'
		), year_range AS (
			SELECT generate_series(
				(SELECT MIN((validity->>'from')::INT) FROM studplan2searchable),
				EXTRACT(YEAR FROM CURRENT_DATE)::INT
			) AS year
		)
		INSERT INTO filter_values (category_id, facet_id, title_cs, title_en, position)
		SELECT
			cid.id category_id,
			yr.year facet_id,
			yr.year title_cs,
			yr.year title_en,
			ROW_NUMBER() OVER (ORDER BY yr.year DESC) position
		FROM year_range yr
		LEFT JOIN category_id cid ON true
	`,
}

/*
Prerequisites:
  - studplan2searchable
*/
var createFilterValuesForDegreePlanStudyTypes = transformation{
	name: "create_filter_values_for_degree_plan_study_types",
	query: `--sql
		WITH category_id AS (
			SELECT fc.id FROM filters f
			LEFT JOIN filter_categories fc ON f.id = fc.filter_id
			WHERE f.id='degree-plans'
			AND fc.facet_id='study_type'
		), distinct_study_type AS (
			SELECT DISTINCT
				study_type
			FROM studplan2searchable
		)
		INSERT INTO filter_values (category_id, facet_id, title_cs, title_en, position)
		SELECT
			cid.id category_id,
			d.study_type facet_id,
			CASE
				WHEN d.study_type = 'B' THEN 'Bakalářské'
				WHEN d.study_type = 'N' THEN 'Navazující magisterské'
				WHEN d.study_type = 'M' THEN 'Magisterské'
				ELSE d.study_type
			END title_cs,
			CASE
				WHEN d.study_type = 'B' THEN 'Bachelor''s'
				WHEN d.study_type = 'N' THEN 'Master''s (post-Bachelor)'
				WHEN d.study_type = 'M' THEN 'Master''s'
				ELSE d.study_type
			END title_en,
			ROW_NUMBER() OVER (ORDER BY d.study_type) position
		FROM distinct_study_type d
		LEFT JOIN category_id cid ON true
	`,
}
