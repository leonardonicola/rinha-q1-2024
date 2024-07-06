package handler

import (
	"github.com/kataras/iris/v12"
	"github.com/rinha-golang/internal/repository"
)

type ExtractHandler struct {
	clientRepo *repository.ClientRepository
}

func NewExtractHandler(repo *repository.ClientRepository) *ExtractHandler {
	return &ExtractHandler{
		clientRepo: repo,
	}
}

func (ex *ExtractHandler) GetExtract(r iris.Context) {
	id := r.Params().Get("id")
	_, err := ex.clientRepo.GetById(id)
	if err != nil {
		r.StopWithStatus(iris.StatusNotFound)
		return
	}

	res, err := ex.clientRepo.GetExtract(id)
	if err != nil {
		r.StopWithStatus(iris.StatusUnprocessableEntity)
		return
	}
	r.StatusCode(iris.StatusOK)
	r.JSON(res)
}
