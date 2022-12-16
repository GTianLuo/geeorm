package session

import (
	"GeeORM/log"
)

func (s *Session) Begin() (err error) {
	if s.tx, err = s.db.Begin(); err != nil {
		log.Error(err)
		return
	}
	log.Info("Start Transaction")
	return
}

func (s *Session) Commit() (err error) {
	if err = s.tx.Commit(); err != nil {
		log.Error(err)
		return
	}
	log.Info("Commit Transaction")
	return
}

func (s *Session) RollBack() (err error) {
	if err = s.tx.Rollback(); err != nil {
		log.Error(err)
		return
	}
	log.Info("RollBack")
	return
}
