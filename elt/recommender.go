package main

import "github.com/jmoiron/sqlx"

func migrateStudium(tx *sqlx.Tx) error {
	var err error
	_, err = tx.Exec(`--sql
		DELETE FROM recommender.studium WHERE TRUE;
		INSERT INTO recommender.studium (
			soident, sident, sfak, sfak2,
			sdruh, sfst, sobor, srokp, 
			sstav, sroc, splan
		) SELECT
			soident, sident, sfak, sfak2,
			sdruh, sfst, sobor, srokp, 
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
