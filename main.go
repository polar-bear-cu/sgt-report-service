package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/polar-bear-cu/sgt-report-service/clients"
	"github.com/polar-bear-cu/sgt-report-service/config"
	"github.com/polar-bear-cu/sgt-report-service/controllers"
	"github.com/polar-bear-cu/sgt-report-service/routes"
	"github.com/polar-bear-cu/sgt-report-service/usecases"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	subConn, err := grpc.NewClient(cfg.SubscriptionAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = subConn.Close() }()

	src := clients.NewSubscriptionClient(subConn)
	uc := usecases.NewReport(src)
	reportCtrl := controllers.NewReport(uc)

	r := gin.Default()
	routes.Register(r, reportCtrl)

	log.Println("listening :" + cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}
