package main

import "github.com/jmoiron/sqlx"

func migratePovinn(tx *sqlx.Tx) error {
	var err error
	_, err = tx.Exec(`--sql
		DELETE FROM recommender.povinn WHERE TRUE;
		INSERT INTO recommender.povinn (
			povinn,pnazev,panazev,
			vplatiod,vplatido,
			pfakulta,pgarant,
			pvyucovan,vsemzac,vsempoc,
			vrozsahpr1,vrozsahcv1,vrozsahpr2,vrozsahcv2,
			vrvcem,vtyp,vebody,
			vucit1,vucit2,vucit3
		) SELECT
			povinn,pnazev,panazev,
			vplatiod,vplatido,
			pfakulta,pgarant,
			pvyucovan,vsemzac,vsempoc,
			vrozsahpr1,vrozsahcv1,vrozsahpr2,vrozsahcv2,
			vrvcem,vtyp,vebody,
			vucit1,vucit2,vucit3
		FROM povinn
	`)
	if err != nil {
		return err
	}
	return nil
}

func migrateStudium(tx *sqlx.Tx) error {
	var err error
	_, err = tx.Exec(`--sql
		DELETE FROM recommender.studium WHERE TRUE;
		INSERT INTO recommender.studium (
			soident, sident, sfak, sfak2,
			sdruh, sobor, srokp, 
			sstav, sroc, splan
		) SELECT
			soident, sident, sfak, sfak2,
			sdruh, sobor, srokp, 
			sstav, sroc, splan 
		FROM studium
	`)
	if err != nil {
		return err
	}
	return nil
}

func migrateZkous(tx *sqlx.Tx) error {
	var err error
	_, err = tx.Exec(`--sql
		DELETE FROM recommender.zkous WHERE TRUE;
		INSERT INTO recommender.zkous (
			zident, zskr, zsem, zpovinn,
			zmarx, zroc, zbody, zsplcelk	
		) SELECT
			zident, zskr, zsem, zpovinn,
			zmarx, zroc, zbody, zsplcelk	
		FROM zkous
	`)
	if err != nil {
		return err
	}
	return nil
}

func migrateStudPlan(tx *sqlx.Tx) error {
	var err error
	_, err = tx.Exec(`--sql
		DELETE FROM recommender.stud_plan WHERE TRUE;
		INSERT INTO recommender.stud_plan (
			code,
			interchangeability,
			bloc_subject_code,
			bloc_type,
			bloc_grade,
			bloc_limit,
			bloc_name_cz,
			bloc_name_en,
			plan_code,
			plan_year
		) SELECT
			code,
			interchangeability,
			bloc_subject_code,
			bloc_type,
			bloc_grade,
			bloc_limit,
			bloc_name_cz,
			bloc_name_en,
			plan_code,
			plan_year
		FROM stud_plan
	`)
	if err != nil {
		return err
	}
	return nil
}
