package cas

import (
	"database/sql"
	"fmt"
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
		return "", fmt.Errorf("Authenticate: %w", err)
	}
	if !userID.Valid {
		return "", fmt.Errorf("Authenticate: %w", sql.ErrNoRows)
	}
	return userID.String, nil
}

func (m DBManager) Login(userID, ticket string) (string, error) {
	var sessionID string
	expiresAt := time.Now().Add(24 * time.Hour)
	query := `
		INSERT INTO sessions (user_id, ticket, expires_at)
		SELECT $1::VARCHAR(8), $2, $3
		WHERE EXISTS ( SELECT id FROM users WHERE id = $1::VARCHAR(8) )
		RETURNING id;
	`
	err := m.DB.Get(&sessionID, query, userID, ticket, expiresAt)
	if err == sql.ErrNoRows {
		err = m.createUser(userID)
		if err != nil {
			return "", fmt.Errorf("Login: %w", err)
		}
		return m.Login(userID, ticket)
	}
	if err != nil {
		fmt.Println(len(ticket))
		return "", fmt.Errorf("Login: %w", err)
	}
	return sessionID, nil
}

func (m DBManager) LogoutWithSession(userID, sessionID string) error {
	query := "DELETE FROM sessions WHERE user_id = $1 AND id = $2"
	_, err := m.DB.Exec(query, userID, sessionID)
	if err != nil {
		return fmt.Errorf("Logout: %w", err)
	}
	return nil
}

func (m DBManager) LogoutWithTicket(userID, ticket string) error {
	query := "DELETE FROM sessions WHERE user_id = $1 AND ticket = $2"
	_, err := m.DB.Exec(query, userID, ticket)
	if err != nil {
		return fmt.Errorf("Logout: %w", err)
	}
	return nil
}

func (m DBManager) createUser(userID string) error {
	tx, err := m.DB.Beginx()
	if err != nil {
		return fmt.Errorf("createUser: %w", err)
	}
	defer tx.Rollback()
	createUserQuery := "INSERT INTO users (id) VALUES ($1)"
	_, err = tx.Exec(createUserQuery, userID)
	if err != nil {
		return fmt.Errorf("createUser: %w", err)
	}
	createBlueprintYearQuery := "INSERT INTO blueprint_years (user_id, academic_year) VALUES ($1, 0) RETURNING id"
	var unassignedYearID int
	err = tx.Get(&unassignedYearID, createBlueprintYearQuery, userID)
	if err != nil {
		return fmt.Errorf("createUser: %w", err)
	}
	createBlueprintUnassignedQuery := "INSERT INTO blueprint_semesters (blueprint_year_id, semester) VALUES ($1, 0)"
	_, err = tx.Exec(createBlueprintUnassignedQuery, unassignedYearID)
	if err != nil {
		return fmt.Errorf("createUser: %w", err)
	}
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("createUser: %w", err)
	}
	return nil
}
