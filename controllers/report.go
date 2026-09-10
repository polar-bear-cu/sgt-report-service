package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/polar-bear-cu/sgt-report-service/dtos"
	"github.com/polar-bear-cu/sgt-report-service/usecases"
)

type ReportController struct {
	uc *usecases.ReportUsecase
}

func NewReport(uc *usecases.ReportUsecase) *ReportController {
	return &ReportController{uc: uc}
}

func (ctl *ReportController) GetMyReport(c *gin.Context) {
	report, err := ctl.uc.GenerateReport(c.Request.Context(), c.GetString("user_id"))
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, dtos.ReportResponse{
		UserID:            report.UserID,
		SubscriptionCount: report.SubscriptionCount,
	})
}
