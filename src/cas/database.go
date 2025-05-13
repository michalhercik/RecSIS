package cas

import (
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
)

type DBManager struct {
	DB *sqlx.DB
}

func (m DBManager) Authenticate(sessionID string) (string, error) {
	var userID sql.NullString
	err := m.DB.Get(&userID, "SELECT user_id FROM sessions WHERE id = $1", sessionID)
	if err != nil {
		return "", err
	}
	if !userID.Valid {
		return "", sql.ErrNoRows
	}
	return userID.String, nil
}

func (m DBManager) Login(userID string) (string, error) {
	var sessionID string
	expiresAt := time.Now().Add(24 * time.Hour)
	query := `
		INSERT INTO sessions (user_id, expires_at)
		SELECT $1, $2
		WHERE EXISTS ( SELECT id FROM users WHERE id = $1 )
		RETURNING id;
	`
	err := m.DB.Get(&sessionID, query, userID, expiresAt)
	if err == sql.ErrNoRows {
		err = m.createUser(userID)
		if err != nil {
			return "", err
		}
		return m.Login(userID)
	}
	if err != nil {
		return "", err
	}
	return sessionID, nil
}

func (m DBManager) Logout(userID, sessionID string) error {
	query := "DELETE FROM sessions WHERE user_id = $1 AND id = $2"
	_, err := m.DB.Exec(query, userID, sessionID)
	if err != nil {
		return err
	}
	return nil
}

func (m DBManager) createUser(userID string) error {
	tx, err := m.DB.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	createUserQuery := "INSERT INTO users (id) VALUES ($1)"
	_, err = tx.Exec(createUserQuery, userID)
	if err != nil {
		return err
	}
	createBlueprintYearQuery := "INSERT INTO blueprint_years (user_id, academic_year) VALUES ($1, 0) RETURNING id"
	var unassignedYearID int
	err = tx.Get(&unassignedYearID, createBlueprintYearQuery, userID)
	if err != nil {
		return err
	}
	createBlueprintUnassignedQuery := "INSERT INTO blueprint_semesters (blueprint_year_id, semester) VALUES ($1, 0)"
	_, err = tx.Exec(createBlueprintUnassignedQuery, unassignedYearID)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}
