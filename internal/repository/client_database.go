package repository

import (
	"database/sql"
	"github.com/pkg/errors"
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
	defer tx.Rollback()

	st, err := tx.Prepare("INSERT INTO clients (guid, username, password) VALUES ($1, $2, $3)")
	if err != nil {
		return nil, errors.Wrap(err, "unable prepare query for client insert")
	}

	_, err = st.Exec(client.GUID, client.Username, client.Password)
	if err != nil {
		return nil, errors.Wrap(err, "unable insert client")
	}

	if err := tx.Commit(); err != nil {
		return nil, errors.Wrap(err, "unable commit on insert client")
	}

	return client, nil
}

// Update used to update client in database
func (repo *ClientDatabase) Update(client *models.Client) (*models.Client, error) {
	tx, err := repo.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	st, err := tx.Prepare("UPDATE clients SET username = $1, password = $2 WHERE guid = $3")
	if err != nil {
		return nil, errors.Wrap(err, "unable prepare query for client update")
	}

	_, err = st.Exec(client.Username, client.Password, client.GUID)
	if err != nil {
		return nil, errors.Wrap(err, "unable update client")
	}

	if err := tx.Commit(); err != nil {
		return nil, errors.Wrap(err, "unable commit update client")
	}

	return client, nil
}

// Delete used to delete client from database
func (repo *ClientDatabase) Delete(client *models.Client) error {
	tx, err := repo.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	st, err := tx.Prepare("DELETE FROM clients WHERE guid = $1")
	if err != nil {
		return errors.Wrap(err, "unable prepare query for client delete")
	}

	_, err = st.Exec(client.GUID)
	if err != nil {
		return errors.Wrap(err, "unable delete client")
	}

	if err := tx.Commit(); err != nil {
		return errors.Wrap(err, "unable commit delete client")
	}

	return nil
}

// FindByName used to find one client by username
func (repo *ClientDatabase) FindByName(username string) (*models.Client, error) {

	rows, err := repo.db.Query("SELECT guid, created, updated, username, password FROM clients WHERE username = $1", username)
	if err != nil {
		return nil, errors.Wrapf(err, "error on querying client by username %s", username)
	}
	defer rows.Close()

	if false == rows.Next() {
		return nil, errors.Errorf("not found client by username %s", username)
	}

	client := &models.Client{}
	if err := rows.Scan(&client.GUID, &client.Created, &client.Updated, &client.Username, &client.Password); err != nil {
		return nil, errors.Wrap(err, "unable scan client")
	}

	return client, nil
}
