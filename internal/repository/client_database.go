package repository

import (
	"database/sql"
	"fmt"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"incrementor/internal/interfaces"
	"incrementor/internal/models"
)

// ClientDatabase struct for manipulating Client on database
type ClientDatabase struct {
	db *sql.DB
}

// NewClientRepositry creates interfaces.ClientRepository
func NewClientRepositry(db *sql.DB) interfaces.ClientRepository {
	return &ClientDatabase{
		db: db,
	}
}

// Create used to insert client in database
func (repo *ClientDatabase) Create(client *models.Client) (*models.Client, error) {
	tx, err := repo.db.Begin()
	if err != nil {
		return nil, err
	}

	st, err := tx.Prepare("INSERT INTO clients (guid, username, password) VALUES ($1, $2, $3)")
	if err != nil {

		if rErr := tx.Rollback(); rErr != nil {
			logrus.Errorf("tx.Rollback on Prepare: %s", rErr.Error())
		}

		return nil, errors.Wrap(err, fmt.Sprintf("Insert.Prepare"))
	}

	_, err = st.Exec(client.GUID, client.Username, client.Password)
	if err != nil {
		if rErr := tx.Rollback(); rErr != nil {
			logrus.Errorf("tx.Rollback on Exec: %s", rErr.Error())
		}

		return nil, errors.Wrap(err, fmt.Sprintf("Insert.Exec: %#v", client))
	}

	if err := tx.Commit(); err != nil {
		logrus.Errorf("tx.Commit: %s", err.Error())
	}

	return client, nil
}

// Update used to update client in database
func (repo *ClientDatabase) Update(client *models.Client) (*models.Client, error) {
	tx, err := repo.db.Begin()
	if err != nil {
		return nil, err
	}

	st, err := tx.Prepare("UPDATE clients SET username = $1, password = $2 WHERE guid = $3")
	if err != nil {

		if rErr := tx.Rollback(); rErr != nil {
			logrus.Errorf("tx.Rollback on Prepare: %s", rErr.Error())
		}

		return nil, errors.Wrap(err, fmt.Sprintf("UPDATE.Prepare"))
	}

	_, err = st.Exec(client.Username, client.Password, client.GUID)
	if err != nil {
		if rErr := tx.Rollback(); rErr != nil {
			logrus.Errorf("tx.Rollback on Exec: %s", rErr.Error())
		}

		return nil, errors.Wrap(err, fmt.Sprintf("Update.Exec: %#v", client))
	}

	if err := tx.Commit(); err != nil {
		logrus.Errorf("tx.Commit: %s", err.Error())
	}

	return client, nil
}

// Delete used to delete client from database
func (repo *ClientDatabase) Delete(client *models.Client) error {
	tx, err := repo.db.Begin()
	if err != nil {
		return err
	}

	st, err := tx.Prepare("DELETE FROM clients WHERE guid = $1")
	if err != nil {

		if rErr := tx.Rollback(); rErr != nil {
			logrus.Errorf("tx.Rollback on Prepare: %s", rErr.Error())
		}

		return errors.Wrap(err, fmt.Sprintf("Delete.Prepare"))
	}

	_, err = st.Exec(client.GUID)
	if err != nil {

		if rErr := tx.Rollback(); rErr != nil {
			logrus.Errorf("tx.Rollback on Exec: %s", rErr.Error())
		}

		return errors.Wrap(err, fmt.Sprintf("DELETE.Exec: %#v", client))
	}

	if err := tx.Commit(); err != nil {
		logrus.Errorf("tx.Commit: %s", err.Error())
	}

	return nil
}

// FindByName used to find one client by username
func (repo *ClientDatabase) FindByName(username string) (*models.Client, error) {

	rows, err := repo.db.Query("SELECT guid, created, updated, username, password FROM clients WHERE username = $1", username)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error on querieng client by username: %s", username))
	}

	defer func() {
		if err := rows.Close(); err != nil {
			logrus.Errorf("Close row error: %s", err.Error())
		}
	}()

	if false == rows.Next() {
		return nil, errors.New(fmt.Sprintf("Not found client by username: %s", username))
	}

	client := &models.Client{}
	if err := rows.Scan(&client.GUID, &client.Created, &client.Updated, &client.Username, &client.Password); err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Could not scan client"))
	}

	return client, nil
}
