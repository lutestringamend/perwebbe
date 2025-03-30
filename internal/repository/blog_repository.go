package repository

import (
	"errors"

	"github.com/lutestringamend/perwebbe/internal/model"
	"github.com/lutestringamend/perwebbe/pkg/paging"
	"gorm.io/gorm"
)

type blogRepository struct {
	db *gorm.DB
}

func NewBlogRepository(db *gorm.DB) BlogRepository {
	return &blogRepository{db: db}
}

func (r *blogRepository) GetBaseQuery() *gorm.DB {
	return r.db.Model(&model.BlogPost{})
}

func (r *blogRepository) Create(post *model.BlogPost) error {
	return r.db.Create(post).Error
}

func (r *blogRepository) GetByID(id uint) (*model.BlogPost, error) {
	var post model.BlogPost
	err := r.db.Preload("Tags").First(&post, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &post, nil
}

func (r *blogRepository) GetAll(page, pageSize int) (*paging.Paginator, error) {
	var posts []model.BlogPost
	pagingParam := &paging.Param{
		DB:      r.db.Preload("Tags"),
		Page:    page,
		Limit:   pageSize,
		OrderBy: []string{"created_at DESC"},
	}

	paginator := paging.Paging(pagingParam, &posts)
	return paginator, nil
}

func (r *blogRepository) GetBySlug(slug string) (*model.BlogPost, error) {
	var post model.BlogPost
	err := r.db.Preload("Tags").Where("slug = ?", slug).First(&post).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &post, nil
}

func (r *blogRepository) Update(post *model.BlogPost) error {
	if err := r.db.Model(post).Association("Tags").Clear(); err != nil {
		return err
	}

	return r.db.Session(&gorm.Session{FullSaveAssociations: true}).Save(post).Error
}

func (r *blogRepository) Delete(id uint) error {
	return r.db.Delete(&model.BlogPost{}, id).Error
}
