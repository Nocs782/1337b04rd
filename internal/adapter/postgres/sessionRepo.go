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
	query := `INSERT INTO sessions (name, avatarurl, created_at) VALUES ($1, $2, $3)`
	_, err := s.db.Exec(query, session.Name, session.AvatarURL, session.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}
