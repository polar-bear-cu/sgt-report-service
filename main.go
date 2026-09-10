package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/polar-bear-cu/sgt-report-service/clients"
	"github.com/polar-bear-cu/sgt-report-service/controllers"
	"github.com/polar-bear-cu/sgt-report-service/routes"
	"github.com/polar-bear-cu/sgt-report-service/usecases"
)

func main() {
	src := clients.StubSubscriptionSource{}
	uc := usecases.NewReport(src)
	reportCtrl := controllers.NewReport(uc)

	r := gin.Default()
	routes.Register(r, reportCtrl)

	log.Println("listening :8086")
	if err := r.Run(":8086"); err != nil {
		log.Fatal(err)
	}
}
