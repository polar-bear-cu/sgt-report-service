package controllers

import (
	"log"
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

// GetMonthly godoc
// @Summary  spending per month, the last item is the current month (still in progress)
// @Tags     reports
// @Produce  json
// @Param    query  query     dtos.MonthlySpendingQuery  false  "number of months"
// @Success  200    {object}  dtos.MonthlySpendingResponse
// @Failure  400    {object}  map[string]string
// @Failure  401    {object}  map[string]string
// @Failure  502    {object}  map[string]string
// @Security BearerAuth
// @Router   /api/v1/reports/monthly [get]
func (ctl *ReportController) GetMonthly(c *gin.Context) {
	var q dtos.MonthlySpendingQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := ctl.uc.MonthlySpending(c.Request.Context(), c.GetString("user_id"), q.Months)
	if err != nil {
		respondErr(c, err)
		return
	}

	months := make([]dtos.MonthTotalResponse, 0, len(res.Months))
	for _, m := range res.Months {
		months = append(months, dtos.MonthTotalResponse{Month: m.Month, Total: m.Total})
	}
	c.JSON(http.StatusOK, dtos.MonthlySpendingResponse{
		Months:          months,
		Total:           res.Total,
		AveragePerMonth: res.AveragePerMonth,
	})
}

// Helpers
func respondErr(c *gin.Context, err error) {
	log.Printf("report: subscription service: %v", err)
	c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
}
