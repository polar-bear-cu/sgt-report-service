package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/polar-bear-cu/sgt-report-service/clients"
	"github.com/polar-bear-cu/sgt-report-service/config"
	"github.com/polar-bear-cu/sgt-report-service/controllers"
	_ "github.com/polar-bear-cu/sgt-report-service/docs"
	"github.com/polar-bear-cu/sgt-report-service/routes"
	"github.com/polar-bear-cu/sgt-report-service/usecases"
)

// @title        Report Service API
// @version      1.0
// @description  Subscription spending report for the Subglutee project
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
	routes.Register(r, reportCtrl, cfg.SwaggerEnabled)
	srv := &http.Server{Addr: ":" + cfg.Port, Handler: r}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		log.Println("listening :" + cfg.Port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
		}
	}()

	<-ctx.Done()
	stop()
	log.Println("shutting down")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("http shutdown: %v", err)
	}
}
