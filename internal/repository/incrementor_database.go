package repository

import (
	"database/sql"
	"fmt"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
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
		return nil, err
	}

	st, err := tx.Prepare("INSERT INTO increments (guid, username, number, step, maxvalue) VALUES ($1, $2, $3, $4, $5)")
	if err != nil {

		if rErr := tx.Rollback(); rErr != nil {
			logrus.Errorf("tx.Rollback on Prepare: %s", rErr.Error())
		}

		return nil, errors.Wrap(err, fmt.Sprintf("Insert.Prepare"))
	}

	_, err = st.Exec(inc.GUID, inc.Username, inc.Number, inc.Step, inc.MaxValue)
	if err != nil {
		if rErr := tx.Rollback(); rErr != nil {
			logrus.Errorf("tx.Rollback on Exec: %s", rErr.Error())
		}

		return nil, errors.Wrap(err, fmt.Sprintf("Insert.Exec: %#v", inc))
	}

	if err := tx.Commit(); err != nil {
		logrus.Errorf("tx.Commit: %s", err.Error())
	}

	return inc, nil
}

// Update used to update client in database
func (repo *IncrementorDatabase) Update(inc *models.Increment) (*models.Increment, error) {
	tx, err := repo.db.Begin()
	if err != nil {
		return nil, err
	}

	st, err := tx.Prepare("UPDATE increments SET username = $1, number = $2, step = $3, maxvalue = $4 WHERE guid = $5")
	if err != nil {

		if rErr := tx.Rollback(); rErr != nil {
			logrus.Errorf("tx.Rollback on Prepare: %s", rErr.Error())
		}

		return nil, errors.Wrap(err, fmt.Sprintf("UPDATE.Prepare"))
	}

	_, err = st.Exec(inc.Username, inc.Number, inc.Step, inc.MaxValue, inc.GUID)
	if err != nil {
		if rErr := tx.Rollback(); rErr != nil {
			logrus.Errorf("tx.Rollback on Exec: %s", rErr.Error())
		}

		return nil, errors.Wrap(err, fmt.Sprintf("Update.Exec: %#v", inc))
	}

	if err := tx.Commit(); err != nil {
		logrus.Errorf("tx.Commit: %s", err.Error())
	}

	return inc, nil
}

// Delete used to delete client from database
func (repo *IncrementorDatabase) Delete(inc *models.Increment) error {
	tx, err := repo.db.Begin()
	if err != nil {
		return err
	}

	st, err := tx.Prepare("DELETE FROM increments WHERE guid = $1")
	if err != nil {

		if rErr := tx.Rollback(); rErr != nil {
			logrus.Errorf("tx.Rollback on Prepare: %s", rErr.Error())
		}

		return errors.Wrap(err, fmt.Sprintf("Delete.Prepare"))
	}

	_, err = st.Exec(inc.GUID)
	if err != nil {

		if rErr := tx.Rollback(); rErr != nil {
			logrus.Errorf("tx.Rollback on Exec: %s", rErr.Error())
		}

		return errors.Wrap(err, fmt.Sprintf("DELETE.Exec: %#v", inc))
	}

	if err := tx.Commit(); err != nil {
		logrus.Errorf("tx.Commit: %s", err.Error())
	}

	return nil
}

// FindByUsername used to find one token by username
func (repo *IncrementorDatabase) FindByUsername(username string) (*models.Increment, error) {

	rows, err := repo.db.Query("SELECT guid, created, updated, username, number, step, maxvalue FROM increments WHERE username = $1", username)

	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error on querieng increment by username: %s", username))
	}

	defer func() {
		if err := rows.Close(); err != nil {
			logrus.Errorf("Close row error: %s", err.Error())
		}
	}()

	if false == rows.Next() {
		return nil, errors.Wrap(err, fmt.Sprintf("Not found increment by username: %s", username))
	}

	inc := &models.Increment{}
	if err := rows.Scan(&inc.GUID, &inc.Created, &inc.Updated, &inc.Username, &inc.Number, &inc.Step, &inc.MaxValue); err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Could not scan incremnt"))
	}

	return inc, nil
}
