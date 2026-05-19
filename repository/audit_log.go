package repository

import (
	"encoding/json"
	"progas-wms-be/global"
	"progas-wms-be/model"

	"gorm.io/gorm"
)

type AuditLogRepository interface {
	Create(log *model.AuditLog) global.ErrorResponse
	Log(userId, action, objectType, objectId string, details any) global.ErrorResponse
}

type auditLogRepository struct {
	db *gorm.DB
}

func NewAuditLogRepository(db *gorm.DB) AuditLogRepository {
	return &auditLogRepository{db: db}
}

func (r *auditLogRepository) Create(log *model.AuditLog) global.ErrorResponse {
	if err := r.db.Create(log).Error; err != nil {
		return global.InternalServerError(err)
	}
	return nil
}

func (r *auditLogRepository) Log(userId, action, objectType, objectId string, details any) global.ErrorResponse {
	logEntry := &model.AuditLog{
		UserId:     userId,
		Action:     action,
		ObjectType: objectType,
		ObjectId:   objectId,
	}
	if details != nil {
		if b, err := json.Marshal(details); err == nil {
			logEntry.Details = string(b)
		}
	}
	return r.Create(logEntry)
}
