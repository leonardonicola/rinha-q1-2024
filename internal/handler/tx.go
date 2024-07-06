package handler

import (
	"github.com/kataras/iris/v12"
	"github.com/rinha-golang/internal/dto"
	"github.com/rinha-golang/internal/repository"
)

type TxHandler struct {
	clientRepo *repository.ClientRepository
}

func NewTxHandler(clientRepo *repository.ClientRepository) *TxHandler {
	return &TxHandler{
		clientRepo: clientRepo,
	}
}

func (tx *TxHandler) MakeTx(r iris.Context) {
	var payload dto.MakeTransactionDTO
	if err := r.ReadJSON(&payload); err != nil {
		r.StopWithStatus(iris.StatusUnprocessableEntity)
		return
	}
	id := r.Params().Get("id")
	client, err := tx.clientRepo.GetById(id)
	if err != nil {
		r.StopWithStatus(iris.StatusNotFound)
		return
	}

	if payload.Type == "d" && client.Balance-int64(payload.Value) <= -(int64(client.Limit)) {
		r.StopWithStatus(iris.StatusUnprocessableEntity)
		return
	}

	newClient, err := tx.clientRepo.MakeTx(id, &payload)
	if err != nil {
		r.StopWithStatus(iris.StatusUnprocessableEntity)
		return
	}
	res := dto.MakeTransactionResponse{
		Limit:   newClient.Limit,
		Balance: newClient.Balance,
	}
	r.StatusCode(iris.StatusOK)
	r.JSON(res)
}
