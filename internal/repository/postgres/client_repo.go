package postgres

import (
	"context"
	"shop_api/internal/dao"
	"shop_api/internal/repository"
	"shop_api/pkg/postgres"

	"github.com/google/uuid"
)

type clientRepoPostgres struct {
	pg *postgres.Postgres
}

func NewClientRepoPostgres(pg *postgres.Postgres) repository.ClientRepository {
	return &clientRepoPostgres{
		pg: pg,
	}
}

func (r *clientRepoPostgres) Save(client *dao.Client) error {
	if client == nil {
		return ErrNilEntity
	}

	query := `
		INSERT INTO clients (
			client_name, client_surname, birthdate, gender, address_id
		) VALUES (
			$1, $2, $3, $4, $5
		)
	`

	_, err := r.pg.Pool.Exec(context.Background(), query,
		client.Name, client.Surname,
		client.Birthdate, client.Gender, client.AddressID,
	)

	if err != nil {
		return ErrQueryExec
	}

	return nil
}

func (r *clientRepoPostgres) RemoveByID(id uuid.UUID) error {
	query := `
		DELETE FROM clients
		WHERE id = $1
	`

	_, err := r.pg.Pool.Exec(context.Background(), query, id)

	if err != nil {
		return ErrQueryExec
	}

	return nil
}

func (r *clientRepoPostgres) GetByID(id uuid.UUID) (*dao.Client, error) {
	var client dao.Client
	query := `
		SELECT
			id, client_name, client_surname, birthdate,
			gender, registration_date, address_id
		FROM clients
		WHERE id = $1
	`

	err := r.pg.Pool.QueryRow(context.Background(), query, id).Scan(
		&client.ID, &client.Name, &client.Surname, &client.Birthdate,
		&client.Gender, &client.RegistrationDate, &client.AddressID,
	)

	if err != nil {
		return nil, ErrQueryExec
	}

	return &client, nil
}

func (r *clientRepoPostgres) GetByName(name, surname string) (*dao.Client, error) {
	var client dao.Client
	query := `
		SELECT
			id, client_name, client_surname, birthdate,
			gender, registration_date, address_id
		FROM clients
		WHERE client_name = $1 AND client_surname = $2
	`

	err := r.pg.Pool.QueryRow(context.Background(), query, name, surname).Scan(
		&client.ID, &client.Name, &client.Surname, &client.Birthdate,
		&client.Gender, &client.RegistrationDate, &client.AddressID,
	)

	if err != nil {
		return nil, ErrQueryExec
	}

	return &client, nil
}

func (r *clientRepoPostgres) GetAll(limit, offset int) ([]*dao.Client, error) {
	query := `
		SELECT
			id, client_name, client_surname, birthdate,
			gender, registration_date, address_id
		FROM clients
		LIMIT $1
		OFFSET $2
	`

	rows, err := r.pg.Pool.Query(context.Background(), query, limit, offset)
	if err != nil {
		return nil, ErrQueryExec
	}
	defer rows.Close()

	var clients []*dao.Client
	for rows.Next() {
		var client dao.Client
		err := rows.Scan(
			&client.ID, &client.Name, &client.Surname, &client.Birthdate,
			&client.Gender, &client.RegistrationDate, &client.AddressID,
		)
		if err != nil {
			return nil, ErrRowScan
		}
		clients = append(clients, &client)
	}

	return clients, nil
}

func (r *clientRepoPostgres) UpdateAddress(clientID, addressID uuid.UUID) error {
	query := `
		UPDATE clients
		SET address_id = $2
		WHERE id = $1
	`

	_, err := r.pg.Pool.Exec(context.Background(), query, clientID, addressID)

	if err != nil {
		return ErrQueryExec
	}

	return nil
}
