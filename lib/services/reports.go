package services

import (
	"context"

	"github.com/officer47p/addressport/lib/db"
	"github.com/officer47p/addressport/lib/types"
)

func NewReportsService(reportsStore db.ReportsStore) ReportsService {
	return ReportsService{reportsStore: reportsStore}
}

type ReportsService struct {
	reportsStore db.ReportsStore
}

func (r *ReportsService) CreateReport(c context.Context, params types.CreateReportParams) (*types.Report, error) {
	report, err := types.NewReportFromParams(params)
	if err != nil {
		return nil, err
	}

	insertedReport, err := r.reportsStore.InsertReport(c, report)
	if err != nil {
		return nil, err
	}

	return insertedReport, nil
}

func (r *ReportsService) GetAllReports(c context.Context) ([]*types.Report, error) {
	reports, err := r.reportsStore.GetReports(c)
	if err != nil {
		return nil, err
	}

	return reports, nil
}

func (r *ReportsService) GetReportsForAddress(c context.Context, address string) ([]*types.Report, error) {
	reports, err := r.reportsStore.GetReportsByAddress(c, address)
	if err != nil {
		return nil, err
	}
	return reports, nil
}

func (r *ReportsService) DeleteReportById(c context.Context, id string) error {
	err := r.reportsStore.DeleteReport(c, id)
	if err != nil {
		return err
	}
	return nil
}

func (r *ReportsService) UpdateReportById(c context.Context, id string, updateParams types.UpdateReportParams) error {
	err := r.reportsStore.DeleteReport(c, id)
	if err != nil {
		return err
	}
	return nil
}
