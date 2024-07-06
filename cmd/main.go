package main

import (
	"fmt"
	"log"
	"os"

	"github.com/rinha-golang/config"
	"github.com/rinha-golang/internal/router"
)

func main() {
	port := os.Getenv("PORT")
	config.InitDB()
	r := router.InitRouter()
	if err := r.Listen(fmt.Sprintf(":%s", port)); err != nil {
		log.Panicf("AN ERROR OCURRED ON SERVER CREATION: %v", err)
	}
}
