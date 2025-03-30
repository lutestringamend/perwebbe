package repository

import (
	"github.com/lutestringamend/perwebbe/internal/model"
	"github.com/lutestringamend/perwebbe/pkg/paging"
	"gorm.io/gorm"
)

type contactRepository struct {
	db *gorm.DB
}

func NewContactRepository(db *gorm.DB) ContactRepository {
	return &contactRepository{db: db}
}

func (r *contactRepository) GetBaseQuery() *gorm.DB {
	return r.db.Model(&model.ContactSubmission{})
}

func (r *contactRepository) Create(submission *model.ContactSubmission) error {
	return r.db.Create(submission).Error
}

func (r *contactRepository) GetAll(page, pageSize int) (*paging.Paginator, error) {
	var submissions []model.ContactSubmission

	pagingParam := &paging.Param{
		DB:      r.db,
		Page:    page,
		Limit:   pageSize,
		OrderBy: []string{"created_at DESC"},
	}

	paginator := paging.Paging(pagingParam, &submissions)
	return paginator, nil
}

func (r *contactRepository) MarkAsRead(id uint) error {
	return r.db.Model(&model.ContactSubmission{}).Where("id = ?", id).Update("read", true).Error
}

func (r *contactRepository) Delete(id uint) error {
	return r.db.Delete(&model.ContactSubmission{}, id).Error
}
