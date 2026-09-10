package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/polar-bear-cu/sgt-report-service/routes"
)

func main() {
	r := gin.Default()
	routes.Register(r)

	log.Println("listening :8086")
	if err := r.Run(":8086"); err != nil {
		log.Fatal(err)
	}
}
