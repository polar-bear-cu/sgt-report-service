package controllers

import (
	"github.com/polar-bear-cu/sgt-report-service/usecases"
)

type ReportController struct {
	uc *usecases.ReportUsecase
}

func NewReport(uc *usecases.ReportUsecase) *ReportController {
	return &ReportController{uc: uc}
}
