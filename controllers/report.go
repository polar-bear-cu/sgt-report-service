package controllers

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

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

// GetBySubscription godoc
// @Summary  total spent per subscription in a month range (default: last 12 months), highest first
// @Tags     reports
// @Produce  json
// @Param    query  query     dtos.SpendingBySubscriptionQuery  false  "from / to as YYYY-MM, both included"
// @Success  200    {object}  dtos.SpendingBySubscriptionResponse
// @Failure  400    {object}  map[string]string
// @Failure  401    {object}  map[string]string
// @Failure  502    {object}  map[string]string
// @Security BearerAuth
// @Router   /api/v1/reports/subscriptions [get]
func (ctl *ReportController) GetBySubscription(c *gin.Context) {
	var q dtos.SpendingBySubscriptionQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := ctl.uc.SpendingBySubscription(c.Request.Context(), c.GetString("user_id"), q.From, q.To)
	if err != nil {
		respondErr(c, err)
		return
	}

	items := make([]dtos.SubscriptionSpendingResponse, 0, len(res.Items))
	for _, s := range res.Items {
		items = append(items, dtos.SubscriptionSpendingResponse{
			SubscriptionID: s.SubscriptionID,
			Name:           s.Name,
			Category:       s.Category,
			Type:           s.Type,
			Status:         s.Status,
			Count:          s.Count,
			Total:          s.Total,
		})
	}
	c.JSON(http.StatusOK, dtos.SpendingBySubscriptionResponse{
		From:  res.From,
		To:    res.To,
		Total: res.Total,
		Items: items,
	})
}

// GetCharges godoc
// @Summary  charges of one month (default: current month) with the total of each category
// @Tags     reports
// @Produce  json
// @Param    query  query     dtos.MonthChargesQuery  false  "month as YYYY-MM"
// @Success  200    {object}  dtos.MonthChargesResponse
// @Failure  400    {object}  map[string]string
// @Failure  401    {object}  map[string]string
// @Failure  502    {object}  map[string]string
// @Security BearerAuth
// @Router   /api/v1/reports/charges [get]
func (ctl *ReportController) GetCharges(c *gin.Context) {
	var q dtos.MonthChargesQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := ctl.uc.MonthCharges(c.Request.Context(), c.GetString("user_id"), q.Month)
	if err != nil {
		respondErr(c, err)
		return
	}

	categories := make([]dtos.CategoryTotalResponse, 0, len(res.Categories))
	for _, cat := range res.Categories {
		categories = append(categories, dtos.CategoryTotalResponse{Category: cat.Category, Total: cat.Total})
	}
	items := make([]dtos.ChargeResponse, 0, len(res.Items))
	for _, ch := range res.Items {
		items = append(items, dtos.ChargeResponse{
			SubscriptionID: ch.SubscriptionID,
			Name:           ch.Name,
			Category:       ch.Category,
			Type:           ch.Type,
			Amount:         ch.Amount,
			ChargedAt:      ch.ChargedAt,
		})
	}
	c.JSON(http.StatusOK, dtos.MonthChargesResponse{
		Month:      res.Month,
		Total:      res.Total,
		Count:      res.Count,
		Categories: categories,
		Items:      items,
	})
}

// ExportCharges godoc
// @Summary  download every charge in a month range (default: last 12 months) as CSV
// @Tags     reports
// @Produce  text/csv
// @Param    query  query  dtos.ExportChargesQuery  false  "from / to as YYYY-MM, both included"
// @Success  200    {file}    file
// @Failure  400    {object}  map[string]string
// @Failure  401    {object}  map[string]string
// @Failure  502    {object}  map[string]string
// @Security BearerAuth
// @Router   /api/v1/reports/export [get]
func (ctl *ReportController) ExportCharges(c *gin.Context) {
	var q dtos.ExportChargesQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := ctl.uc.ExportCharges(c.Request.Context(), c.GetString("user_id"), q.From, q.To)
	if err != nil {
		respondErr(c, err)
		return
	}

	var buf bytes.Buffer
	buf.WriteString("\ufeff") // BOM so Excel opens UTF-8 (e.g. Thai names) correctly
	w := csv.NewWriter(&buf)
	_ = w.Write([]string{"date", "subscription", "category", "type", "amount"})
	for _, ch := range res.Items {
		_ = w.Write([]string{
			ch.ChargedAt.Format(time.DateOnly),
			ch.Name,
			ch.Category,
			ch.Type,
			strconv.FormatFloat(ch.Amount, 'f', 2, 64),
		})
	}
	w.Flush()
	if err := w.Error(); err != nil {
		log.Printf("report: write csv: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="sgt-charges-%s_%s.csv"`, res.From, res.To))
	c.Data(http.StatusOK, "text/csv; charset=utf-8", buf.Bytes())
}

// Helpers
func respondErr(c *gin.Context, err error) {
	if errors.Is(err, usecases.ErrValidation) {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Printf("report: subscription service: %v", err)
	c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
}
