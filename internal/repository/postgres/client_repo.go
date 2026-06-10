package postgres

import (
	"context"
	"shop_api/internal/entity"
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

func (r *clientRepoPostgres) Save(client *entity.Client) error {
	if client == nil {
		return repository.ErrNilEntity
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
		return repository.ErrQueryExec
	}

	return nil
}

func (r *clientRepoPostgres) RemoveByID(id uuid.UUID) error {
	query := `
		DELETE FROM clients
		WHERE id = $1
	`

	cmdTag, err := r.pg.Pool.Exec(context.Background(), query, id)

	if isFKViolation(err) {
		return repository.ErrDependentEntity
	} else if err != nil {
		return repository.ErrQueryExec
	}

	if cmdTag.RowsAffected() == 0 {
		return repository.ErrNoRows
	}

	return nil
}

func (r *clientRepoPostgres) GetWithParams(name, surname *string, limit, offset int) ([]*entity.Client, error) {
	query := `
		SELECT
			id, client_name, client_surname, birthdate,
			gender, registration_date, address_id
		FROM clients
		WHERE
			($1::text IS NULL OR client_name = $1)
			AND ($2::text IS NULL OR client_surname = $2)
		LIMIT $3
		OFFSET $4
	`

	rows, err := r.pg.Pool.Query(context.Background(), query, name, surname, limit, offset)
	if err != nil {
		return nil, repository.ErrQueryExec
	}
	defer rows.Close()

	var clients []*entity.Client
	for rows.Next() {
		var client entity.Client
		err := rows.Scan(
			&client.ID, &client.Name, &client.Surname, &client.Birthdate,
			&client.Gender, &client.RegistrationDate, &client.AddressID,
		)
		if err != nil {
			return nil, repository.ErrRowScan
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

	cmdTag, err := r.pg.Pool.Exec(context.Background(), query, clientID, addressID)

	if err != nil {
		return repository.ErrQueryExec
	}

	if cmdTag.RowsAffected() == 0 {
		return repository.ErrNoRows
	}

	return nil
}
