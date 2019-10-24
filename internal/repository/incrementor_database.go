package repository

import (
	"database/sql"
	"fmt"
	"github.com/pkg/errors"
	"incrementor/internal/interfaces"
	"incrementor/internal/models"
)

// IncrementorDatabase struct for manipulating Client on database
type IncrementorDatabase struct {
	db *sql.DB
}

// NewIncrementorRepositry creates interfaces.IncrementRepository
func NewIncrementorRepositry(db *sql.DB) interfaces.IncrementRepository {
	return &IncrementorDatabase{
		db: db,
	}
}

// Create used to insert token in database
func (repo *IncrementorDatabase) Create(inc *models.Increment) (*models.Increment, error) {
	tx, err := repo.db.Begin()
	if err != nil {
		return nil, errors.Wrap(err, "unable start transaction for increment create")
	}
	tx.Rollback()

	st, err := tx.Prepare("INSERT INTO increments (guid, client_guid, number, step, maxvalue) VALUES ($1, $2, $3, $4, $5)")
	if err != nil {
		return nil, errors.Wrap(err, "unable prepare query for increment insert")
	}

	_, err = st.Exec(inc.GUID, inc.ClientGUID, inc.Number, inc.Step, inc.MaxValue)
	if err != nil {
		return nil, errors.Wrap(err, "unable execute query increment insert")
	}

	if err := tx.Commit(); err != nil {
		return nil, errors.Wrap(err, "unable commit increment insert")
	}

	return inc, nil
}

// Update used to update client in database
func (repo *IncrementorDatabase) Update(inc *models.Increment) (*models.Increment, error) {
	tx, err := repo.db.Begin()
	if err != nil {
		return nil, errors.Wrap(err, "unable start transaction for increment update")
	}
	tx.Rollback()

	st, err := tx.Prepare("UPDATE increments SET client_guid = $1, number = $2, step = $3, maxvalue = $4 WHERE guid = $5")
	if err != nil {
		return nil, errors.Wrap(err, "unable prepare query increment update")
	}

	_, err = st.Exec(inc.ClientGUID, inc.Number, inc.Step, inc.MaxValue, inc.GUID)
	if err != nil {
		return nil, errors.Wrap(err, "unable execute query increment update")
	}

	if err := tx.Commit(); err != nil {
		return nil, errors.Wrap(err, "unable commit increment update")
	}

	return inc, nil
}

// Delete used to delete client from database
func (repo *IncrementorDatabase) Delete(inc *models.Increment) error {
	tx, err := repo.db.Begin()
	if err != nil {
		return errors.Wrap(err, "unable start transaction for increment delete")
	}
	tx.Rollback()

	st, err := tx.Prepare("DELETE FROM increments WHERE guid = $1")
	if err != nil {
		return errors.Wrap(err, "unable prepare query for increment delete")
	}

	_, err = st.Exec(inc.GUID)
	if err != nil {
		return errors.Wrap(err, "unable execute query increment delete")
	}

	if err := tx.Commit(); err != nil {
		return errors.Wrap(err, "unable commit increment delete")
	}

	return nil
}

// FindByUsername used to find one token by username
func (repo *IncrementorDatabase) FindByUsername(username string) (*models.Increment, error) {

	rows, err := repo.db.Query("SELECT guid, created, updated, client_guid, number, step, maxvalue FROM increments WHERE username = $1", username)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error on querieng increment by username: %s", username))
	}
	defer rows.Close()

	if false == rows.Next() {
		return nil, errors.Wrapf(err, "not found increment by username %s", username)
	}

	inc := &models.Increment{}
	if err := rows.Scan(&inc.GUID, &inc.Created, &inc.Updated, &inc.ClientGUID, &inc.Number, &inc.Step, &inc.MaxValue); err != nil {
		return nil, errors.Wrap(err, "unable scan increment")
	}

	return inc, nil
}
