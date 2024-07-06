package router

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/recover"

	"github.com/go-playground/validator/v10"
	"github.com/rinha-golang/config"
	"github.com/rinha-golang/internal/handler"
	"github.com/rinha-golang/internal/repository"
)

func InitRouter() *iris.Application {
	db := config.GetDB()
	r := iris.New()
	v := validator.New()
	r.Validator = v
	r.UseRouter(recover.New())
	clientRepo := repository.NewEventRepository(db)
	txHandler := handler.NewTxHandler(clientRepo)
	exHandler := handler.NewExtractHandler(clientRepo)
	client := r.Party("/clientes/{id}")
	{
		client.Post("/transacoes", txHandler.MakeTx)
		client.Get("/extrato", exHandler.GetExtract)
	}
	return r
}
