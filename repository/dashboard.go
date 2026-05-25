package repository

import (
	"progas-wms-be/enum"
	"progas-wms-be/global"
	"progas-wms-be/model"
	"time"

	"gorm.io/gorm"
)

type DashboardRepository interface {
	CountCylindersByStatus() (map[string]int, global.ErrorResponse)
	CountOutstandingCylinders() (int64, global.ErrorResponse)
	FindCustomersOverQuota() ([]model.Customer, global.ErrorResponse)
	CountHydrotestExpired() (int64, global.ErrorResponse)
	CountHydrotestDueSoon(withinDays int) (int64, global.ErrorResponse)
}

type dashboardRepository struct {
	db *gorm.DB
}

func NewDashboardRepository(db *gorm.DB) DashboardRepository {
	return &dashboardRepository{db: db}
}

func (r *dashboardRepository) CountCylindersByStatus() (map[string]int, global.ErrorResponse) {
	type row struct {
		Status string
		Count  int
	}
	var rows []row
	if err := r.db.Model(&model.Cylinder{}).
		Select("status, COUNT(*) as count").
		Group("status").Scan(&rows).Error; err != nil {
		return nil, global.InternalServerError(err)
	}
	result := make(map[string]int)
	for _, row := range rows {
		result[row.Status] = row.Count
	}
	return result, nil
}

func (r *dashboardRepository) CountOutstandingCylinders() (int64, global.ErrorResponse) {
	var count int64
	if err := r.db.Model(&model.Cylinder{}).
		Where("status = ?", enum.CylinderStatusOutstanding).
		Count(&count).Error; err != nil {
		return 0, global.InternalServerError(err)
	}
	return count, nil
}

func (r *dashboardRepository) FindCustomersOverQuota() ([]model.Customer, global.ErrorResponse) {
	var customers []model.Customer
	err := r.db.Where("cylinder_quota_limit > 0 AND outstanding_count > cylinder_quota_limit AND is_active = ?", true).
		Order("outstanding_count desc").Find(&customers).Error
	if err != nil {
		return nil, global.InternalServerError(err)
	}
	return customers, nil
}

func (r *dashboardRepository) CountHydrotestExpired() (int64, global.ErrorResponse) {
	expiryBefore := time.Now().AddDate(-5, 0, 0)
	var count int64
	if err := r.db.Model(&model.Cylinder{}).
		Where("last_hydrotest_date < ?", expiryBefore).
		Count(&count).Error; err != nil {
		return 0, global.InternalServerError(err)
	}
	return count, nil
}

func (r *dashboardRepository) CountHydrotestDueSoon(withinDays int) (int64, global.ErrorResponse) {
	now := time.Now()
	expiredBefore := now.AddDate(-5, 0, 0)
	dueSoonBy := now.AddDate(0, 0, withinDays)
	var count int64
	if err := r.db.Model(&model.Cylinder{}).
		Where("last_hydrotest_date < ? OR DATE_ADD(last_hydrotest_date, INTERVAL 5 YEAR) <= ?",
			expiredBefore, dueSoonBy).
		Count(&count).Error; err != nil {
		return 0, global.InternalServerError(err)
	}
	return count, nil
}
