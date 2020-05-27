package sqlite

import (
	"billing_api/pkg/models"
	"database/sql"
)

// CountersModel to hold the DB connections
type CountersModel struct {
	DB *sql.DB
}

// CreateTable will create table if not exists in the db
func (r *CountersModel) CreateTable() (bool, error) {
	sqlStmt := `CREATE TABLE IF NOT EXISTS requests (authenticatorID TEXT, scanID TEXT PRIMARY KEY, pageCount INT);`
	_, err := r.DB.Exec(sqlStmt)
	if err != nil {
		return false, err
	}
	return true, nil
}

// Insert used to insert the data in table
func (r *CountersModel) Insert(authenticatorID string, scanID string, pageCount int) (int, error) {
	stmt := `INSERT INTO requests (authenticatorID, scanID, pageCount) 
	VALUES(?, ?, ?);`
	result, err := r.DB.Exec(stmt, authenticatorID, scanID, pageCount)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

// GetCounts calculates the request count & page count.
func (r *CountersModel) GetCounts(authenticatorID string) (*models.Counters, error) {
	stmt := `SELECT count(scanID), sum(pageCount) from requests WHERE authenticatorID=?;`
	row := r.DB.QueryRow(stmt, authenticatorID)
	s := &models.Counters{}
	s.AuthenticatorID = authenticatorID
	err := row.Scan(&s.RequestCount, &s.PageCount)
	if err == sql.ErrNoRows {
		return nil, models.ErrNoRecord
	} else if err != nil {
		return nil, err
	}
	return s, nil
}
