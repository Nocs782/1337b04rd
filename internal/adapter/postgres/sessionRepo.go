package postgres

import (
	"1337b04rd/internal/domain"
	"database/sql"
)

type SessionRepo struct {
	db *sql.DB
}

func NewSessionRepo(db *sql.DB) *SessionRepo {
	return &SessionRepo{
		db: db,
	}
}

func (s *SessionRepo) CreateSession(session domain.Session) error {
	query := `INSERT INTO sessions (id, name, avatar_url, created_at, expires_at) VALUES ($1, $2, $3, $4, $5)`
	_, err := s.db.Exec(query, session.ID, session.Name, session.AvatarURL, session.CreatedAt, session.ExpiresAt)
	if err != nil {
		return err
	}
	return nil
}

func (s *SessionRepo) GetSessionByID(id string) (*domain.Session, error) {
	query := `SELECT id, name, avatar_url, created_at, expires_at FROM sessions WHERE id = $1`
	row := s.db.QueryRow(query, id)

	var session domain.Session
	err := row.Scan(&session.ID, &session.Name, &session.AvatarURL, &session.CreatedAt, &session.ExpiresAt)
	if err != nil {
		return nil, err
	}

	return &session, nil
}
