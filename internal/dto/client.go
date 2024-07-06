package dto

import (
	"errors"
	"net/http"
	"reflect"
	"time"
)

type MakeTransactionDTO struct {
	Value       uint64 `json:"valor" validate:"required,number,gte=1"`
	Type        string `json:"tipo" validate:"required,oneof=c d"`
	Description string `json:"descricao" validate:"required,lte=10"`
}

func (dto *MakeTransactionDTO) Bind(r *http.Request) error {

	if reflect.TypeOf(dto.Value).Kind() != reflect.Uint64 {
		return errors.New("Apenas números inteiros")
	}

	return nil
}

type MakeTransactionResponse struct {
	Limit   uint64 `json:"limite"`
	Balance int64  `json:"saldo"`
}

type Balance struct {
	Total     int64     `json:"total"`
	CreatedAt time.Time `json:"data_extrato"`
	Limit     string    `json:"limite"`
}

type Historic struct {
	Value       int64     `json:"valor"`
	Type        string    `json:"tipo"`
	Description string    `json:"descricao"`
	CreatedAt   time.Time `json:"realizada_em"`
}

type GetExtractResponse struct {
	Balance   Balance    `json:"saldo"`
	Historico []Historic `json:"ultimas_transacoes"`
}
