package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/quantum0hound/gochat/internal/models"
)

type SessionProviderPostgres struct {
	db *sqlx.DB
}

func NewSessionProviderPostgres(db *sqlx.DB) *SessionProviderPostgres {
	return &SessionProviderPostgres{db: db}
}

func (r *SessionProviderPostgres) Create(session *models.Session) error {
	query := fmt.Sprintf(`INSERT INTO %s (refresh_token,expires_in,user_id,fingerprint) 
								 VALUES ($1, $2, $3, $4)`, sessionsTable)
	_, err := r.db.Exec(query, session.RefreshToken, session.ExpiresIn, session.UserId, session.Fingerprint)
	return err
}

func (r *SessionProviderPostgres) Delete(id string) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE refresh_token = $1`, sessionsTable)
	_, err := r.db.Exec(query, id)
	return err
}

func (r *SessionProviderPostgres) Get(id string) (*models.Session, error) {
	var session models.Session
	query := fmt.Sprintf("SELECT * FROM %s WHERE (refresh_token=$1)", sessionsTable)
	err := r.db.Get(&session, query, id)
	if err != nil {
		return nil, err
	}
	return &session, nil
}
