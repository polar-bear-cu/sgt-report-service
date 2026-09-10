package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/polar-bear-cu/sgt-report-service/controllers"
)

func Register(r *gin.Engine, report *controllers.ReportController) {
	r.GET("/health", controllers.GetHealth)

	v1 := r.Group("/api/v1")
	v1.GET("/reports/me", report.GetMyReport)
}
