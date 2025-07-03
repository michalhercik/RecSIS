package main

import (
	"database/sql"
	"fmt"
	"log"
	"reflect"
	"strings"
	"sync"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type Extract struct {
	from       *sqlx.DB
	to         *sqlx.DB
	operations []operation
}

func makeExtract(from, to *sqlx.DB) *Extract {
	e := &Extract{
		from: from,
		to:   to,
	}
	e.operations = []operation{}
	return e
}

func (e *Extract) add(op operation) {
	if op == nil {
		log.Println("Attempted to add nil operation")
		return
	}
	e.operations = append(e.operations, op)
}

func (e *Extract) run() error {
	wg := sync.WaitGroup{}
	errsCh := make(chan ELTResult, len(e.operations))
	for _, op := range e.operations {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if res := e.runOperation(op); res.IsError() {
				errsCh <- res
			}
		}()
	}
	wg.Wait()
	close(errsCh)
	var errs []error
	for res := range errsCh {
		errs = append(errs, res)
	}
	return listOfErrors(errs)
}

func (e *Extract) runOperation(op operation) ELTResult {
	// log.Printf("ℹ️  Extracting %s from SIS", op.name())
	err := op.selectData(e.from, e.to)
	if err != nil {
		log.Printf("❌ %s: Select: %v", op.name(), err)
		return ELTResult{err: err, proc: op}
	}
	// log.Printf("ℹ️  Inserting %s to RecSIS", op.name())
	err = op.insertData(e.to)
	if err != nil {
		log.Printf("❌ %s: Insert: %v", op.name(), err)
		return ELTResult{err: err, proc: op}
	}
	log.Printf("✅ Extraction of %s successfull", op.name())
	return ELTResult{}
}

type ELTResult struct {
	err  error
	proc operation
}

func (e ELTResult) IsError() bool {
	return e.err != nil
}

func (e ELTResult) Error() string {
	return fmt.Sprintf("ExtractError: %s", e.err.Error())
}

type operation interface {
	name() string
	selectData(from *sqlx.DB, to *sqlx.DB) error
	insertData(to *sqlx.DB) error
}

// POVINN
// ===========================================================================================================
type povinnRecord struct {
	POVINN     string         `db:"POVINN"`
	PNAZEV     sql.NullString `db:"PNAZEV"`
	PANAZEV    sql.NullString `db:"PANAZEV"`
	VPLATIOD   sql.NullInt64  `db:"VPLATIOD"`
	VPLATIDO   sql.NullInt64  `db:"VPLATIDO"`
	PFAKULTA   sql.NullString `db:"PFAKULTA"`
	PGARANT    sql.NullString `db:"PGARANT"`
	PVYUCOVAN  sql.NullString `db:"PVYUCOVAN"`
	VSEMZAC    sql.NullString `db:"VSEMZAC"`
	VSEMPOC    sql.NullString `db:"VSEMPOC"`
	PVYJAZYK   sql.NullString `db:"PVYJAZYK"`
	VROZSAHPR1 sql.NullInt64  `db:"VROZSAHPR1"`
	VROZSAHCV1 sql.NullInt64  `db:"VROZSAHCV1"`
	VROZSAHPR2 sql.NullInt64  `db:"VROZSAHPR2"`
	VROZSAHCV2 sql.NullInt64  `db:"VROZSAHCV2"`
	VRVCEM     sql.NullString `db:"VRVCEM"`
	VTYP       sql.NullString `db:"VTYP"`
	VEBODY     sql.NullInt64  `db:"VEBODY"`
	VUCIT1     sql.NullString `db:"VUCIT1"`
	VUCIT2     sql.NullString `db:"VUCIT2"`
	VUCIT3     sql.NullString `db:"VUCIT3"`
	PPOCMIN    sql.NullInt64  `db:"PPOCMIN"`
	PPOCMAX    sql.NullInt64  `db:"PPOCMAX"`
	PURL       sql.NullString `db:"PURL"`
}
type extractPovinn struct {
	data []povinnRecord
}

func (ep *extractPovinn) name() string {
	return "POVINN"
}

func (ep *extractPovinn) selectData(from *sqlx.DB, to *sqlx.DB) error {
	query := `
		SELECT
			POVINN,PNAZEV,PANAZEV,
			VPLATIOD,VPLATIDO,
			PFAKULTA,PGARANT,
			PVYUCOVAN,VSEMZAC,VSEMPOC,PVYJAZYK,
			VROZSAHPR1,VROZSAHCV1,VROZSAHPR2,VROZSAHCV2,
			VRVCEM,VTYP,VEBODY,
			VUCIT1,VUCIT2,VUCIT3,
			PPOCMIN,PPOCMAX,
			PURL
		FROM POVINN
		WHERE TO_CHAR(sysdate, 'YYYY') BETWEEN VPLATIOD AND VPLATIDO
		AND PFAKULTA = '11320'
	`
	err := from.Select(&ep.data, query)
	if err != nil {
		return fmt.Errorf("selectData: %w", err)
	}
	return nil
}

func (ep *extractPovinn) insertData(to *sqlx.DB) error {
	drop := `--sql
		DROP TABLE IF EXISTS povinn
	`
	create := `
		CREATE TABLE povinn (
			POVINN     VARCHAR(10),
			PNAZEV     VARCHAR(250),
			PANAZEV    VARCHAR(250),
			VPLATIOD   INT,
			VPLATIDO   INT,
			PFAKULTA   VARCHAR(5),
			PGARANT    VARCHAR(10),
			PVYUCOVAN  VARCHAR(1),
			VSEMZAC    VARCHAR(1),
			VSEMPOC    INT,
			PVYJAZYK   VARCHAR(6),
			VROZSAHPR1 INT,
			VROZSAHCV1 INT,
			VROZSAHPR2 INT,
			VROZSAHCV2 INT,
			VRVCEM     VARCHAR(2),
			VTYP       VARCHAR(2),
			VEBODY     INT,
			VUCIT1     VARCHAR(10),
			VUCIT2     VARCHAR(10),
			VUCIT3     VARCHAR(10),
			PPOCMIN    INT,
			PPOCMAX    INT,
			PURL 	   VARCHAR(250)
		)
	`
	insert := `
		INSERT INTO povinn (
			POVINN,PNAZEV,PANAZEV,
			VPLATIOD,VPLATIDO,
			PFAKULTA,PGARANT,
			PVYUCOVAN,VSEMZAC,VSEMPOC,PVYJAZYK,
			VROZSAHPR1,VROZSAHCV1,VROZSAHPR2,VROZSAHCV2,
			VRVCEM,VTYP,VEBODY,
			VUCIT1,VUCIT2,VUCIT3,
			PPOCMIN,PPOCMAX,
			PURL
		)
		(SELECT * FROM unnest(
			$1::text[], $2::text[], $3::text[],
			$4::int[], $5::int[],
			$6::text[], $7::text[], $8::text[], $9::text[],
			$10::int[],
			$11::text[],
			$12::int[], $13::int[], $14::int[], $15::int[],
			$16::text[], $17::text[],
			$18::int[],
			$19::text[], $20::text[], $21::text[],
			$22::int[], $23::int[],
			$24::text[]
			))
	`
	err := insertAsColumns(to, drop, create, insert, toColumns(ep.data))
	if err != nil {
		return fmt.Errorf("insertData: %w", err)
	}
	return nil
}

// UCIT_ROZVRH
// ===========================================================================================================
type extractUcitRozvrh struct {
	data []struct {
		POVINN string `db:"POVINN"`
		UCIT   string `db:"UCIT"`
	}
}

func (ep *extractUcitRozvrh) name() string {
	return "UCIT_ROZVRH"
}
func (ep *extractUcitRozvrh) selectData(from *sqlx.DB, to *sqlx.DB) error {
	var (
		query string
		err   error
	)
	query = `
		SELECT ur.povinn, ur.ucit FROM ucit_rozvrh ur
		LEFT JOIN povinn p ON ur.povinn = p.povinn
		WHERE ur.skr = ( SELECT MAX(skr) FROM ucit_rozvrh )
		AND PFAKULTA = '11320'
	`
	if err = from.Select(&ep.data, query); err != nil {
		return fmt.Errorf("selectData: retrieve ucit_rozvrh: %w", err)
	}
	return nil
}
func (ep *extractUcitRozvrh) insertData(to *sqlx.DB) error {
	drop := `--sql
		DROP TABLE IF EXISTS ucit_rozvrh
	`
	create := `
		CREATE TABLE ucit_rozvrh (
			POVINN VARCHAR(10),
			UCIT VARCHAR(10)
		)
	`
	insert := `
		INSERT INTO ucit_rozvrh (
			POVINN, UCIT
		)
		( SELECT * FROM unnest(
			$1::text[], $2::text[]
		))
	`
	err := insertAsColumns(to, drop, create, insert, toColumns(ep.data))
	if err != nil {
		return err
	}
	return nil
}

// UCIT
// ===========================================================================================================

type extractUcit struct {
	data []struct {
		KOD       string         `db:"KOD"`
		USTAV     sql.NullString `db:"USTAV"`
		FAKULTA   sql.NullString `db:"FAKULTA"`
		PRIJMENI  sql.NullString `db:"PRIJMENI"`
		JMENO     sql.NullString `db:"JMENO"`
		TITULPRED sql.NullString `db:"TITULPRED"`
		TITULZA   sql.NullString `db:"TITULZA"`
	}
}

func (ep *extractUcit) name() string {
	return "UCIT"
}

func (ep *extractUcit) selectData(from *sqlx.DB, to *sqlx.DB) error {
	var query string
	var err error
	query = `
		SELECT
			KOD, USTAV, FAKULTA,
			PRIJMENI, JMENO, TITULPRED, TITULZA
		FROM UCIT
	`
	err = from.Select(&ep.data, query)
	return err
}

func (ep *extractUcit) insertData(to *sqlx.DB) error {
	drop := `--sql
		DROP TABLE IF EXISTS ucit
	`
	create := `
		CREATE TABLE ucit (
			KOD VARCHAR(10),
			USTAV VARCHAR(10),
			FAKULTA VARCHAR(5),
			PRIJMENI VARCHAR(50),
			JMENO VARCHAR(50),
			TITULPRED VARCHAR(100),
			TITULZA VARCHAR(100)
		)
	`
	insert := `
		INSERT INTO ucit (
			KOD, USTAV, FAKULTA,
			PRIJMENI, JMENO, TITULPRED, TITULZA
		)
		( SELECT * FROM unnest(
			$1::text[], $2::text[], $3::text[],
			$4::text[], $5::text[], $6::text[], $7::text[]
		))
	`
	err := insertAsColumns(to, drop, create, insert, toColumns(ep.data))
	if err != nil {
		return err
	}
	return nil
}

// ANKECY
// ===========================================================================================================

type extractAnkecy struct {
	data []struct {
		SDRUH   sql.NullString `db:"SDRUH"`
		SROC    sql.NullInt64  `db:"SROC"`
		POVINN  sql.NullString `db:"POVINN"`
		SSKR    sql.NullInt64  `db:"SSKR"`
		SEM     sql.NullInt64  `db:"SEM"`
		SOBOR   sql.NullString `db:"SOBOR"`
		PRDMTYP sql.NullString `db:"PRDMTYP"`
		UCIT    sql.NullString `db:"UCIT"`
		MEMO    sql.NullString `db:"MEMO"`
	}
}

func (ep *extractAnkecy) name() string {
	return "ANKECY"
}

func (ep *extractAnkecy) selectData(from *sqlx.DB, to *sqlx.DB) error {
	var query string
	var err error
	query = `
		SELECT
			SDRUH,
			SROC,
			POVINN,
			SSKR,
			SEM,
			SOBOR,
			PRDMTYP,
			UCIT,
			MEMO
		FROM ANKECY
	`
	err = from.Select(&ep.data, query)
	return err
}

func (ep *extractAnkecy) insertData(to *sqlx.DB) error {
	drop := `--sql
		DROP TABLE IF EXISTS ankecy
	`
	create := `
		CREATE TABLE ankecy (
			SDRUH VARCHAR(2),
			SROC INT,
			POVINN VARCHAR(10),
			SSKR INT,
			SEM INT,
			SOBOR VARCHAR(12),
			PRDMTYP VARCHAR(30),
			UCIT VARCHAR(10),
			MEMO TEXT
		)
	`
	insert := `
		INSERT INTO ankecy (
			SDRUH, SROC, POVINN, SSKR,
			SEM, SOBOR, PRDMTYP, UCIT, MEMO
		)
		( SELECT * FROM unnest(
			$1::text[], $2::int[], $3::text[],
			$4::int[], $5::int[], $6::text[],
			$7::text[], $8::text[], $9::text[]
		))
	`
	err := insertAsColumns(to, drop, create, insert, toColumns(ep.data))
	if err != nil {
		return err
	}
	return nil
}

// PAMELA
// ===========================================================================================================

type extractPamela struct {
	data []struct {
		POVINN string         `db:"POVINN"`
		TYP    string         `db:"TYP"`
		JAZYK  string         `db:"JAZYK"`
		MEMO   sql.NullString `db:"MEMO"`
	}
}

func (ep *extractPamela) name() string {
	return "PAMELA"
}

func (ep *extractPamela) selectData(from *sqlx.DB, to *sqlx.DB) error {
	var query string
	var err error
	query = `
		SELECT
			pamela.POVINN, pamela.TYP, pamela.JAZYK, pamela.MEMO
		FROM pamela
		LEFT JOIN POVINN ON PAMELA.POVINN = POVINN.POVINN
		WHERE TO_CHAR(sysdate, 'YYYY') BETWEEN PAMELA.VPLATIOD AND PAMELA.VPLATIDO
		AND PRAVO='ALL'
		AND TO_CHAR(sysdate, 'YYYY') BETWEEN POVINN.VPLATIOD AND POVINN.VPLATIDO
		AND POVINN.PFAKULTA='11320'
		AND (POVINN.PVYUCOVAN = 'V' OR POVINN.PVYUCOVAN = 'N' OR POVINN.PVYUCOVAN = 'P')
	`
	err = from.Select(&ep.data, query)
	if err != nil {
		return fmt.Errorf("selectData: %w", err)
	}
	return err
}

func (ep *extractPamela) insertData(to *sqlx.DB) error {
	drop := `--sql
		DROP TABLE IF EXISTS pamela
	`
	create := `
		CREATE TABLE pamela (
			POVINN VARCHAR(10),
			TYP VARCHAR(1),
			JAZYK VARCHAR(6),
			MEMO TEXT
		)
	`
	insert := `
		INSERT INTO pamela (
			POVINN, TYP, JAZYK, MEMO
		)
		( SELECT * FROM unnest(
			$1::text[], $2::text[], $3::text[], $4::text[]
		))
	`
	// PostgreSQL does not support null characters in text fields, so we need to remove them
	ep.removeNullCharsFromMemo()
	err := insertAsColumns(to, drop, create, insert, toColumns(ep.data))
	if err != nil {
		return err
	}
	return nil
}

func (ep *extractPamela) removeNullCharsFromMemo() {
	for i := range ep.data {
		if ep.data[i].MEMO.Valid {
			ep.data[i].MEMO.String = strings.Replace(ep.data[i].MEMO.String, "\x00", "", -1)
		}
	}
}

// JAZYK
// ============================================================================================================

type extractJazyk struct {
	data []struct {
		KOD    string         `db:"KOD"`
		NAZEV  sql.NullString `db:"NAZEV"`
		ANAZEV sql.NullString `db:"ANAZEV"`
	}
}

func (ep *extractJazyk) name() string {
	return "JAZYK"
}

func (ep *extractJazyk) selectData(from *sqlx.DB, to *sqlx.DB) error {
	var query string
	var err error
	query = `
		SELECT KOD, NAZEV, ANAZEV FROM JAZYK
		WHERE (DDO > sysdate OR DDO IS NULL)
		AND (DOD < sysdate OR DOD IS NULL)
	`
	err = from.Select(&ep.data, query)
	if err != nil {
		return fmt.Errorf("selectData: %w", err)
	}
	return err
}

func (ep *extractJazyk) insertData(to *sqlx.DB) error {
	drop := `--sql
		DROP TABLE IF EXISTS jazyk
	`
	create := `
		CREATE TABLE jazyk (
			KOD VARCHAR(6),
			NAZEV VARCHAR(120),
			ANAZEV VARCHAR(120)
		)
	`
	insert := `
		INSERT INTO jazyk (
			KOD,
			NAZEV, ANAZEV
		) VALUES(
			:KOD,
			:NAZEV, :ANAZEV
		)
	`
	err := simpleInsert(to, drop, create, insert, ep.data)
	if err != nil {
		return fmt.Errorf("insertData: %w", err)
	}
	return nil
}

// DRUH
// ============================================================================================================

type extractDruh struct {
	data []struct {
		KOD      string         `db:"KOD"`
		NAZEV    sql.NullString `db:"NAZEV"`
		ANAZEV   sql.NullString `db:"ANAZEV"`
		ZKRATKA  sql.NullString `db:"ZKRATKA"`
		AZKRATKA sql.NullString `db:"AZKRATKA"`
	}
}

func (ep *extractDruh) name() string {
	return "DRUH"
}

func (ep *extractDruh) selectData(from *sqlx.DB, to *sqlx.DB) error {
	var query string
	var err error
	query = `
		SELECT
			KOD,
			NAZEV, ANAZEV,
			ZKRATKA, AZKRATKA
		FROM DRUH
	`
	err = from.Select(&ep.data, query)
	if err != nil {
		return fmt.Errorf("selectData: %w", err)
	}
	return err
}

func (ep *extractDruh) insertData(to *sqlx.DB) error {
	drop := `--sql
		DROP TABLE IF EXISTS druh
	`
	create := `
		CREATE TABLE druh (
			KOD VARCHAR(2),
			NAZEV VARCHAR(60),
			ANAZEV VARCHAR(60),
			ZKRATKA VARCHAR(10),
			AZKRATKA VARCHAR(10)
		)
	`
	insert := `
		INSERT INTO druh (
			KOD,
			NAZEV, ANAZEV,
			ZKRATKA, AZKRATKA
		) VALUES(
			:KOD,
			:NAZEV, :ANAZEV,
			:ZKRATKA, :AZKRATKA
		)
	`
	err := simpleInsert(to, drop, create, insert, ep.data)
	if err != nil {
		return fmt.Errorf("insertData: %w", err)
	}
	return nil
}

// KLAS
// ============================================================================================================

type extractKlas struct {
	data []struct {
		KOD    string         `db:"KOD"`
		NAZEV  sql.NullString `db:"NAZEV"`
		ANAZEV sql.NullString `db:"ANAZEV"`
	}
}

func (ep *extractKlas) name() string {
	return "KLAS"
}

func (ep *extractKlas) selectData(from *sqlx.DB, to *sqlx.DB) error {
	var query string
	var err error
	query = `
		SELECT
			KOD,
			NAZEV, ANAZEV
		FROM KLAS
	`
	err = from.Select(&ep.data, query)
	if err != nil {
		return fmt.Errorf("selectData: %w", err)
	}
	return err
}

func (ep *extractKlas) insertData(to *sqlx.DB) error {
	drop := `--sql
		DROP TABLE IF EXISTS klas
	`
	create := `
		CREATE TABLE klas (
			KOD VARCHAR(6),
			NAZEV VARCHAR(60),
			ANAZEV VARCHAR(60)
		)
	`
	insert := `
		INSERT INTO klas (
			KOD,
			NAZEV, ANAZEV
		) VALUES(
			:KOD,
			:NAZEV, :ANAZEV
		)
	`
	err := simpleInsert(to, drop, create, insert, ep.data)
	if err != nil {
		return fmt.Errorf("insertData: %w", err)
	}
	return nil
}

// PKLAS
// ===========================================================================================================

type extractPklas struct {
	data []struct {
		POVINN string `db:"POVINN"`
		PKLAS  string `db:"PKLAS"`
	}
}

func (ep *extractPklas) name() string {
	return "PKLAS"
}

func (ep *extractPklas) selectData(from *sqlx.DB, to *sqlx.DB) error {
	var query string
	var err error
	query = `
		SELECT
			POVINN, PKLAS
		FROM PKLAS
	`
	err = from.Select(&ep.data, query)
	if err != nil {
		return err
	}
	return err
}

func (ep *extractPklas) insertData(to *sqlx.DB) error {
	drop := `--sql
		DROP TABLE IF EXISTS pklas
	`
	create := `
		CREATE TABLE pklas (
			POVINN VARCHAR(10),
			PKLAS VARCHAR(6)
		)
	`
	insert := `
		INSERT INTO pklas (
			POVINN, PKLAS
		)
		( SELECT * FROM unnest(
			$1::text[], $2::text[]
		))
	`
	err := insertAsColumns(to, drop, create, insert, toColumns(ep.data))
	if err != nil {
		return err
	}
	return nil
}

// POVINN2JAZYK
// ============================================================================================================

type extractPovinn2Jazyk struct {
	data []struct {
		POVINN string `db:"POVINN"`
		JAZYK  string `db:"JAZYK"`
	}
}

func (ep *extractPovinn2Jazyk) name() string {
	return "POVINN2JAZYK"
}

func (ep *extractPovinn2Jazyk) selectData(from *sqlx.DB, to *sqlx.DB) error {
	var query string
	var err error
	query = `
		SELECT
			POVINN, JAZYK
		FROM POVINN2JAZYK
		WHERE TO_CHAR(sysdate, 'YYYY') BETWEEN PLATIOD AND PLATIDO
	`
	err = from.Select(&ep.data, query)
	if err != nil {
		return fmt.Errorf("selectData: %w", err)
	}
	return err
}

func (ep *extractPovinn2Jazyk) insertData(to *sqlx.DB) error {
	drop := `--sql
		DROP TABLE IF EXISTS povinn2jazyk
	`
	create := `
		CREATE TABLE povinn2jazyk (
			POVINN VARCHAR(10),
			JAZYK VARCHAR(6)
		)
	`
	insert := `
		INSERT INTO povinn2jazyk (
			POVINN, JAZYK
		) VALUES(
			:POVINN, :JAZYK
		)
	`
	err := simpleInsert(to, drop, create, insert, ep.data)
	if err != nil {
		return err
	}
	return nil
}

// POVINN2JAZYK
// ============================================================================================================

type extractPreq struct {
	data []struct {
		POVINN    string `db:"POVINN"`
		REQTYP    string `db:"REQTYP"`
		REQPOVINN string `db:"REQPOVINN"`
	}
}

func (ep *extractPreq) name() string {
	return "PREQ"
}

func (ep *extractPreq) selectData(from *sqlx.DB, to *sqlx.DB) error {
	var query string
	var err error
	query = `
		SELECT
			PREQ.POVINN, PREQ.REQTYP, REQPOVINN
		FROM PREQ
		LEFT JOIN POVINN ON PREQ.POVINN = POVINN.POVINN
		WHERE to_char(sysdate, 'YYYY') BETWEEN PREQ.REQOD AND PREQ.REQDO
		AND TO_CHAR(sysdate, 'YYYY') BETWEEN POVINN.VPLATIOD AND POVINN.VPLATIDO
		AND POVINN.PFAKULTA='11320'
		AND (POVINN.PVYUCOVAN = 'V' OR POVINN.PVYUCOVAN = 'N' OR POVINN.PVYUCOVAN = 'P')
	`
	err = from.Select(&ep.data, query)
	if err != nil {
		return err
	}
	return err
}

func (ep *extractPreq) insertData(to *sqlx.DB) error {
	drop := `--sql
		DROP TABLE IF EXISTS preq
	`
	create := `
		CREATE TABLE preq (
			POVINN VARCHAR(10),
			REQTYP VARCHAR(1),
			REQPOVINN VARCHAR(10)
		)
	`
	insert := `
		INSERT INTO preq (
			POVINN, REQTYP, REQPOVINN
		) VALUES(
			:POVINN, :REQTYP, :REQPOVINN
		)
	`
	err := simpleInsert(to, drop, create, insert, ep.data)
	if err != nil {
		return err
	}
	return nil
}

// PTRIDA
// ============================================================================================================

type extractPtrida struct {
	data []struct {
		POVINN string `db:"POVINN"`
		PTRIDA string `db:"PTRIDA"`
	}
}

func (ep *extractPtrida) name() string {
	return "PTRIDA"
}

func (ep *extractPtrida) selectData(from *sqlx.DB, to *sqlx.DB) error {
	var query string
	var err error
	query = `
		SELECT
			POVINN, PTRIDA
		FROM PTRIDA
	`
	err = from.Select(&ep.data, query)
	if err != nil {
		return err
	}
	return err
}

func (ep *extractPtrida) insertData(to *sqlx.DB) error {
	drop := `--sql
		DROP TABLE IF EXISTS ptrida
	`
	create := `
		CREATE TABLE ptrida (
			POVINN VARCHAR(10),
			PTRIDA VARCHAR(7)
		)
	`
	insert := `
		INSERT INTO ptrida (
			POVINN, PTRIDA
		) VALUES(
			:POVINN, :PTRIDA
		)
	`
	err := simpleInsert(to, drop, create, insert, ep.data)
	if err != nil {
		return err
	}
	return nil
}

// TRIDA
// ============================================================================================================

type extractTrida struct {
	data []struct {
		KOD     string         `db:"KOD"`
		FAKULTA string         `db:"FAKULTA"`
		NAZEV   sql.NullString `db:"NAZEV"`
	}
}

func (ep *extractTrida) name() string {
	return "TRIDA"
}

func (ep *extractTrida) selectData(from *sqlx.DB, to *sqlx.DB) error {
	var query string
	var err error
	query = `
		SELECT
			KOD, FAKULTA, NAZEV
		FROM TRIDA
		WHERE (DDO >= to_char(sysdate, 'YYYY') OR DDO IS NULL)
		AND (DOD <= to_char(sysdate, 'YYYY') OR DOD IS NULL)
	`
	err = from.Select(&ep.data, query)
	if err != nil {
		return err
	}
	return err
}

func (ep *extractTrida) insertData(to *sqlx.DB) error {
	drop := `--sql
		DROP TABLE IF EXISTS trida
	`
	create := `
		CREATE TABLE trida (
			KOD VARCHAR(7),
			FAKULTA VARCHAR(5),
			NAZEV VARCHAR(50)
		)
	`
	insert := `
		INSERT INTO trida (
			KOD, FAKULTA, NAZEV
		) VALUES(
			:KOD, :FAKULTA, :NAZEV
		)
	`
	err := simpleInsert(to, drop, create, insert, ep.data)
	if err != nil {
		return err
	}
	return nil
}

// TYPYPOV
// ============================================================================================================

type extractTypyPov struct {
	data []struct {
		KOD    string         `db:"KOD"`
		NAZEV  sql.NullString `db:"NAZEV"`
		EXAM1  sql.NullString `db:"EXAM1"`
		EXAM2  sql.NullString `db:"EXAM2"`
		ANAZEV sql.NullString `db:"ANAZEV"`
		AEXAM1 sql.NullString `db:"AEXAM1"`
		AEXAM2 sql.NullString `db:"AEXAM2"`
	}
}

func (ep *extractTypyPov) name() string {
	return "TYPYPOV"
}

func (ep *extractTypyPov) selectData(from *sqlx.DB, to *sqlx.DB) error {
	var query string
	var err error
	query = `
		SELECT
			KOD,
			NAZEV, EXAM1, EXAM2,
			ANAZEV, AEXAM1, AEXAM2
		FROM TYPYPOV
		WHERE (DDO >= to_char(sysdate, 'YYYY') OR DDO IS NULL)
		AND (DOD <= to_char(sysdate, 'YYYY') OR DOD IS NULL)
	`
	err = from.Select(&ep.data, query)
	if err != nil {
		return err
	}
	return err
}

func (ep *extractTypyPov) insertData(to *sqlx.DB) error {
	drop := `--sql
		DROP TABLE IF EXISTS typypov
	`
	create := `
		CREATE TABLE typypov (
			KOD VARCHAR(2),
			NAZEV VARCHAR(70),
			EXAM1 VARCHAR(15),
			EXAM2 VARCHAR(15),
			ANAZEV VARCHAR(70),
			AEXAM1 VARCHAR(15),
			AEXAM2 VARCHAR(15)
		)
	`
	insert := `
		INSERT INTO typypov (
			KOD,
			NAZEV, EXAM1, EXAM2,
			ANAZEV, AEXAM1, AEXAM2
		) VALUES(
			:KOD,
			:NAZEV, :EXAM1, :EXAM2,
			:ANAZEV, :AEXAM1, :AEXAM2
		)
	`
	err := simpleInsert(to, drop, create, insert, ep.data)
	if err != nil {
		return err
	}
	return nil
}

// SEKCE
// ============================================================================================================

type extractSekce struct {
	data []struct {
		KOD     string         `db:"KOD"`
		NAZEV   sql.NullString `db:"NAZEV"`
		ANAZEV  sql.NullString `db:"ANAZEV"`
		FAKULTA sql.NullString `db:"FAKULTA"`
	}
}

func (ep *extractSekce) name() string {
	return "SEKCE"
}

func (ep *extractSekce) selectData(from *sqlx.DB, to *sqlx.DB) error {
	var query string
	var err error
	query = `
		SELECT
			KOD, NAZEV, ANAZEV, FAKULTA
		FROM SEKCE
	`
	err = from.Select(&ep.data, query)
	if err != nil {
		return fmt.Errorf("selectData: %w", err)
	}
	return err
}

func (ep *extractSekce) insertData(to *sqlx.DB) error {
	drop := `--sql
		DROP TABLE IF EXISTS sekce
	`
	create := `
		CREATE TABLE sekce (
			KOD VARCHAR(10),
			NAZEV VARCHAR(50),
			ANAZEV VARCHAR(50),
			FAKULTA VARCHAR(5)
		)
	`
	insert := `
		INSERT INTO sekce (
			KOD,
			NAZEV, ANAZEV,
			FAKULTA
		) VALUES(
			:KOD,
			:NAZEV, :ANAZEV,
			:FAKULTA
		)
	`
	err := simpleInsert(to, drop, create, insert, ep.data)
	if err != nil {
		return fmt.Errorf("insertData: %w", err)
	}
	return nil
}

// USTAV
// ============================================================================================================

type extractUstav struct {
	data []struct {
		KOD     string         `db:"KOD"`
		NAZEV   sql.NullString `db:"NAZEV"`
		ANAZEV  sql.NullString `db:"ANAZEV"`
		SEKCE   sql.NullString `db:"SEKCE"`
		FAKULTA sql.NullString `db:"FAKULTA"`
	}
}

func (ep *extractUstav) name() string {
	return "USTAV"
}

func (ep *extractUstav) selectData(from *sqlx.DB, to *sqlx.DB) error {
	var query string
	var err error
	query = `
		SELECT
			KOD,
			NAZEV, ANAZEV,
			SEKCE, FAKULTA
		FROM USTAV
		WHERE FAKULTA='11320'
	`
	err = from.Select(&ep.data, query)
	if err != nil {
		return fmt.Errorf("selectData: %w", err)
	}
	return err
}

func (ep *extractUstav) insertData(to *sqlx.DB) error {
	drop := `--sql
		DROP TABLE IF EXISTS ustav
	`
	create := `
		CREATE TABLE ustav (
			KOD VARCHAR(10),
			NAZEV VARCHAR(100),
			ANAZEV VARCHAR(200),
			SEKCE VARCHAR(10),
			FAKULTA VARCHAR(5)
		)
	`
	insert := `
		INSERT INTO ustav (
			KOD,
			NAZEV, ANAZEV,
			SEKCE, FAKULTA
		) VALUES(
			:KOD,
			:NAZEV, :ANAZEV,
			:SEKCE, :FAKULTA
		)
	`
	err := simpleInsert(to, drop, create, insert, ep.data)
	if err != nil {
		return fmt.Errorf("insertData: %w", err)
	}
	return nil
}

// FAK
// ============================================================================================================

type extractFak struct {
	data []struct {
		KOD      string         `db:"KOD"`
		NAZEV    sql.NullString `db:"NAZEV"`
		ANAZEV   sql.NullString `db:"ANAZEV"`
		ZKRATKA  sql.NullString `db:"ZKRATKA"`
		AZKRATKA sql.NullString `db:"AZKRATKA"`
	}
}

func (ep *extractFak) name() string {
	return "FAK"
}

func (ep *extractFak) selectData(from *sqlx.DB, to *sqlx.DB) error {
	var query string
	var err error
	query = `
		SELECT
			KOD,
			TEXT NAZEV, ANAZEV,
			NAZEV ZKRATKA, AZKRATKA
		FROM FAK
	`
	err = from.Select(&ep.data, query)
	if err != nil {
		return fmt.Errorf("selectData: %w", err)
	}
	return err
}

func (ep *extractFak) insertData(to *sqlx.DB) error {
	drop := `--sql
		DROP TABLE IF EXISTS fak
	`
	create := `
		CREATE TABLE fak (
			KOD VARCHAR(5),
			NAZEV VARCHAR(150),
			ANAZEV VARCHAR(150),
			ZKRATKA VARCHAR(10),
			AZKRATKA VARCHAR(10)
		)
	`
	insert := `
		INSERT INTO fak (
			KOD,
			NAZEV, ANAZEV,
			ZKRATKA, AZKRATKA
		) VALUES(
			:KOD,
			:NAZEV, :ANAZEV,
			:ZKRATKA, :AZKRATKA
		)
	`
	err := simpleInsert(to, drop, create, insert, ep.data)
	if err != nil {
		return fmt.Errorf("insertData: %w", err)
	}
	return nil
}

// CISELNIK
// ============================================================================================================

type Ciselnik struct {
	Table     string
	KodSize   int
	NazevSize int
	data      []struct {
		KOD    string         `db:"KOD"`
		NAZEV  sql.NullString `db:"NAZEV"`
		ANAZEV sql.NullString `db:"ANAZEV"`
	}
}

func (ep *Ciselnik) name() string {
	return strings.ToUpper(ep.Table)
}

func (ep *Ciselnik) selectData(from *sqlx.DB, to *sqlx.DB) error {
	var query string
	var err error
	query = `
		SELECT
			KOD,
			NAZEV, ANAZEV
		FROM %s
	`
	query = fmt.Sprintf(query, ep.Table)
	err = from.Select(&ep.data, query)
	if err != nil {
		return fmt.Errorf("selectData: %w", err)
	}
	return err
}

func (ep *Ciselnik) insertData(to *sqlx.DB) error {
	drop := `--sql
		DROP TABLE IF EXISTS %s
	`
	drop = fmt.Sprintf(drop, ep.Table)
	create := `
		CREATE TABLE %s (
			KOD VARCHAR(%d),
			NAZEV VARCHAR(%d),
			ANAZEV VARCHAR(%d)
		)
	`
	if ep.KodSize < 1 {
		ep.KodSize = 6
	}
	if ep.NazevSize < 1 {
		ep.NazevSize = 120
	}
	create = fmt.Sprintf(create, ep.Table, ep.KodSize, ep.NazevSize, ep.NazevSize)
	insert := `
		INSERT INTO %s (
			KOD,
			NAZEV, ANAZEV
		) VALUES(
			:KOD,
			:NAZEV, :ANAZEV
		)
	`
	insert = fmt.Sprintf(insert, ep.Table)
	err := simpleInsert(to, drop, create, insert, ep.data)
	if err != nil {
		return fmt.Errorf("insertData: %w", err)
	}
	return nil
}

// Helpers
// ============================================================================================================

func simpleInsert[T any](db *sqlx.DB, drop, create, insert string, data []T) error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	err = dropCreateTable(tx, drop, create)
	if err != nil {
		return err
	}
	_, err = tx.NamedExec(insert, data)
	if err != nil {
		return err
	}
	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}

func insertAsColumns(db *sqlx.DB, drop, create, insert string, data []any) error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	err = dropCreateTable(tx, drop, create)
	if err != nil {
		return err
	}
	_, err = tx.Exec(insert, data...)
	if err != nil {
		return err
	}
	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}

func dropCreateTable(tx *sqlx.Tx, drop, create string) error {
	_, err := tx.Exec(drop)
	if err != nil {
		return err
	}
	_, err = tx.Exec(create)
	if err != nil {
		return err
	}
	return nil
}

// func Map[T any, V any](data []T, fn func(T) V) []V {
// 	result := make([]V, len(data))
// 	for i, item := range data {
// 		result[i] = fn(item)
// 	}
// 	return result
// }

func toColumns[T any](data []T) []any {
	numFields := reflect.TypeOf((*T)(nil)).Elem().NumField()
	columns := make([][]interface{}, numFields)
	for i := range columns {
		columns[i] = make([]interface{}, 0, len(data))
	}
	for _, row := range data {
		v := reflect.ValueOf(row)
		for i := 0; i < numFields; i++ {
			columns[i] = append(columns[i], v.Field(i).Interface())
		}
	}
	result := make([]any, numFields)
	for i, col := range columns {
		switch reflect.TypeOf(col[0]).String() {
		case "string":
			tmp := make([]string, len(col))
			for j, v := range col {
				tmp[j] = v.(string)
			}
			result[i] = pq.Array(tmp)
		case "sql.NullString":
			tmp := make([]sql.NullString, len(col))
			for j, v := range col {
				tmp[j] = v.(sql.NullString)
			}
			result[i] = pq.Array(tmp)
		case "sql.NullInt64":
			tmp := make([]sql.NullInt64, len(col))
			for j, v := range col {
				tmp[j] = v.(sql.NullInt64)
			}
			result[i] = pq.Array(tmp)
		default:
			result[i] = pq.Array(col)
		}
	}
	return result
}
