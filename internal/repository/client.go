package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/rinha-golang/internal/domain"
	"github.com/rinha-golang/internal/dto"
)

type ClientRepository struct {
	db *sql.DB
}

func NewEventRepository(db *sql.DB) *ClientRepository {
	return &ClientRepository{
		db: db,
	}
}

func (repo *ClientRepository) GetById(id string) (*domain.Client, error) {

	client := &domain.Client{}
	const sqlQuery = "SELECT id, limite, saldo FROM cliente WHERE id = $1"
	err := repo.db.QueryRow(sqlQuery, id).Scan(&client.ID, &client.Limit, &client.Balance)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("Cliente não encontrado: %w", err)
		}
		return nil, fmt.Errorf("Erro ao buscar cliente: %v", err)
	}
	return client, nil
}

func (repo *ClientRepository) GetExtract(id string) (*dto.GetExtractResponse, error) {
	tx, err := repo.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("erro ao iniciar transação: %v", err)
	}
	// Rollback in case anything fails
	defer tx.Rollback()

	var balance dto.Balance
	var query string
	query = "SELECT limite, saldo, NOW() as created_at FROM cliente WHERE id = $1"
	err = tx.QueryRow(query, id).Scan(&balance.Limit, &balance.Total, &balance.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("Cliente %s não encontrado: %v", id, err)
	}

	query = "SELECT valor, tipo, descricao, created_at FROM historico WHERE client_id = $1 ORDER BY created_at DESC LIMIT 10"
	rows, err := tx.Query(query, id)
	if err != nil {
		return nil, fmt.Errorf("Erro ao buscar historico do cliente %s: %v", id, err)
	}
	defer rows.Close()

	var historic []dto.Historic
	for rows.Next() {
		var transaction dto.Historic
		err := rows.Scan(&transaction.Value, &transaction.Type, &transaction.Description, &transaction.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("Erro ao buscar historico do cliente %s: %v", id, err)
		}
		historic = append(historic, transaction)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("Erro ao buscar historico do cliente %s: %v", id, err)
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("Erro ao confirmar transação: %v", err)
	}

	return &dto.GetExtractResponse{
		Balance:   balance,
		Historico: historic,
	}, nil
}

func (repo *ClientRepository) MakeTx(id string, payload *dto.MakeTransactionDTO) (*domain.Client, error) {
	tx, err := repo.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("erro ao iniciar transação: %v", err)
	}
	// Rollback in case anything fails
	defer tx.Rollback()

	// Check if client has enough limit
	var query string
	var client domain.Client
	if payload.Type == "c" {
		query = "UPDATE cliente SET saldo = saldo + $1 WHERE id = $2 RETURNING limite, saldo"
	} else {
		query = "UPDATE cliente SET saldo = saldo - $1 WHERE id = $2 AND saldo - $1 >= -ABS(limite) RETURNING limite, saldo"
	}
	err = tx.QueryRow(query, payload.Value, id).Scan(&client.Limit, &client.Balance)
	if err != nil {
		return nil, fmt.Errorf("Erro ao realizar transação: %v", err)
	}

	if err != nil {
		return nil, fmt.Errorf("Erro ao realizar transação: %v", err)
	}

	query = "INSERT INTO historico (valor, descricao, tipo, client_id) VALUES ($1, $2, $3, $4)"
	result, err := tx.Exec(query, payload.Value, payload.Description, payload.Type, id)
	if err != nil {
		return nil, fmt.Errorf("Erro ao realizar transação: %v", err)
	}
	rows, err := result.RowsAffected()
	if err != nil || rows != 1 {
		return nil, fmt.Errorf("Erro ao realizar transação: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("Erro ao concluir transação: %v", err)
	}
	return &domain.Client{
		Limit:   client.Limit,
		Balance: client.Balance,
	}, nil
}
