package postgres

import (
	"workflow-service/internal/domain/entity"

	"gorm.io/gorm"
)

type AuditRepoPG struct {
	db *gorm.DB
}

func NewAuditRepoPG(db *gorm.DB) *AuditRepoPG {
	return &AuditRepoPG{db: db}
}

func (r *AuditRepoPG) Create(audit *entity.AuditLog) error {
	return r.db.Create(audit).Error
}
