package data

import (
	"context"
	"database/sql"
)

type Storage struct {
	db  *sql.DB
	cfg *Config
	// queries map[string]*sql.Stmt
}

func (s *Storage) Open() (err error) {
	s.cfg = &Config{}
	s.cfg.Read()

	err = s.connectToDatabase()

	if err != nil {
		return err
	}

	return s.prepareQueries()
}

func (s *Storage) Close(ctx context.Context) (err error) {
	closed := make(chan struct{}, 1)

	go func() {
		// for _, query := range s.queries {
		// 	err = errors.Join(err, query.Close())
		// }

		// err = errors.Join(err, s.db.Close())

		closed <- struct{}{}
	}()

	select {
	case <-closed:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (s *Storage) connectToDatabase() (err error) {
	// dbAddress := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", s.cfg.user, s.cfg.password, s.cfg.host, s.cfg.port, s.cfg.database)

	// s.db, err = sql.Open("pgx", dbAddress)

	// if err == nil {
	// 	err = s.db.Ping()
	// }

	// if err == nil {
	// 	log.Info().Msg(fmt.Sprintf("Successfully connected to DB at %s", dbAddress))
	// }

	return err
}

func (s *Storage) prepareQueries() (err error) {
	// s.queries = make(map[string]*sql.Stmt)

	// err = errors.Join(err, s.prepare(queryNewUserCreate, queryNewUserCreateName))
	// err = errors.Join(err, s.prepare(queryPasswordIsValid, queryPasswordIsValidName))

	return err
}

func (s *Storage) prepare(query, name string) (err error) {
	// s.queries[name], err = s.db.Prepare(query)
	return err
}
