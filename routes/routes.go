package routes

import (
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginswagger "github.com/swaggo/gin-swagger"

	"github.com/polar-bear-cu/sgt-report-service/controllers"
	"github.com/polar-bear-cu/sgt-report-service/middlewares"
)

func Register(r *gin.Engine, report *controllers.ReportController, jwtSecret string, swaggerEnabled bool) {
	r.GET("/health", controllers.GetHealth)

	if swaggerEnabled {
		r.GET("/swagger/*any", ginswagger.WrapHandler(swaggerfiles.Handler))
	}

	v1 := r.Group("/api/v1")
	v1.Use(middlewares.RequireAuth(jwtSecret))
	v1.GET("/reports/me", report.GetMyReport)
}
