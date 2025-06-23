package cas

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/michalhercik/RecSIS/errorx"
	"github.com/michalhercik/RecSIS/language"
)

type DBManager struct {
	DB *sqlx.DB
}

func (m DBManager) authenticate(sessionID string, lang language.Language) (string, error) {
	var userID sql.NullString
	err := m.DB.Get(&userID, "SELECT user_id FROM sessions WHERE id = $1", sessionID)
	if err != nil {
		return "", errorx.NewHTTPErr(
			errorx.AddContext(err),
			http.StatusInternalServerError,
			texts[lang].errCannotGetUserIDFromSession,
		)
	}
	if !userID.Valid {
		return "", errorx.NewHTTPErr(
			errorx.AddContext(fmt.Errorf("session %s does not exist in DB - userID is NULL", sessionID)),
			http.StatusUnauthorized,
			texts[lang].errUnauthorized,
		)
	}
	return userID.String, nil
}

func (m DBManager) login(userID, ticket string, lang language.Language) (string, error) {
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
		err = m.createUser(userID, lang)
		if err != nil {
			return "", errorx.AddContext(err)
		}
		sessionID, err = m.login(userID, ticket, lang)
		if err != nil {
			return "", errorx.AddContext(err)
		}
	} else if err != nil {
		//fmt.Println(len(ticket))
		return "", errorx.NewHTTPErr(
			errorx.AddContext(err),
			http.StatusInternalServerError,
			texts[lang].errCannotCreateSession,
		)
	}
	return sessionID, nil
}

func (m DBManager) logoutWithSession(userID, sessionID string, lang language.Language) error {
	query := "DELETE FROM sessions WHERE user_id = $1 AND id = $2"
	_, err := m.DB.Exec(query, userID, sessionID)
	if err != nil {
		return errorx.NewHTTPErr(
			errorx.AddContext(err),
			http.StatusInternalServerError,
			texts[lang].errCannotLogout,
		)
	}
	return nil
}

func (m DBManager) logoutWithTicket(userID, ticket string, lang language.Language) error {
	query := "DELETE FROM sessions WHERE user_id = $1 AND ticket = $2"
	_, err := m.DB.Exec(query, userID, ticket)
	if err != nil {
		return errorx.NewHTTPErr(
			errorx.AddContext(err),
			http.StatusInternalServerError,
			texts[lang].errCannotLogout,
		)
	}
	return nil
}

func (m DBManager) createUser(userID string, lang language.Language) error {
	tx, err := m.DB.Beginx()
	if err != nil {
		return errorx.NewHTTPErr(
			errorx.AddContext(err),
			http.StatusInternalServerError,
			texts[lang].errCannotCreateUser,
		)
	}
	defer tx.Rollback()
	createUserQuery := "INSERT INTO users (id) VALUES ($1)"
	_, err = tx.Exec(createUserQuery, userID)
	if err != nil {
		return errorx.NewHTTPErr(
			errorx.AddContext(err),
			http.StatusInternalServerError,
			texts[lang].errCannotCreateUser,
		)
	}
	createBlueprintYearQuery := "INSERT INTO blueprint_years (user_id, academic_year) VALUES ($1, 0) RETURNING id"
	var unassignedYearID int
	err = tx.Get(&unassignedYearID, createBlueprintYearQuery, userID)
	if err != nil {
		return errorx.NewHTTPErr(
			errorx.AddContext(err),
			http.StatusInternalServerError,
			texts[lang].errCannotCreateUser,
		)
	}
	createBlueprintUnassignedQuery := "INSERT INTO blueprint_semesters (blueprint_year_id, semester) VALUES ($1, 0)"
	_, err = tx.Exec(createBlueprintUnassignedQuery, unassignedYearID)
	if err != nil {
		return errorx.NewHTTPErr(
			errorx.AddContext(err, errorx.P("blueprint_year_id", unassignedYearID)),
			http.StatusInternalServerError,
			texts[lang].errCannotCreateUser,
		)
	}
	// TODO: remove this after SIS integration
	createStudy := "INSERT INTO bla_studies (user_id, degree_plan_code, start_year) VALUES ($1, 'NIPVS19B', 2020)"
	_, err = m.DB.Exec(createStudy, userID)
	if err != nil {
		return errorx.NewHTTPErr(
			errorx.AddContext(err),
			http.StatusInternalServerError,
			texts[lang].errCannotCreateUser,
		)
	}
	err = tx.Commit()
	if err != nil {
		return errorx.NewHTTPErr(
			errorx.AddContext(err),
			http.StatusInternalServerError,
			texts[lang].errCannotCreateUser,
		)
	}
	return nil
}
