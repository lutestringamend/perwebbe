package repository

import (
	"github.com/lutestringamend/perwebbe/internal/model"
	"github.com/lutestringamend/perwebbe/pkg/paging"
	"gorm.io/gorm"
)

// BlogRepository defines methods for blog post repository
type BlogRepository interface {
	Create(post *model.BlogPost) error
	GetByID(id uint) (*model.BlogPost, error)
	GetAll(page, pageSize int) (*paging.Paginator, error)
	GetBySlug(slug string) (*model.BlogPost, error)
	Update(post *model.BlogPost) error
	Delete(id uint) error
	GetBaseQuery() *gorm.DB
}

// PortfolioRepository defines methods for portfolio project repository
type PortfolioRepository interface {
	Create(project *model.PortfolioProject) error
	GetByID(id uint) (*model.PortfolioProject, error)
	GetAll(projectType string, page, pageSize int) (*paging.Paginator, error)
	Update(project *model.PortfolioProject) error
	Delete(id uint) error
	GetBaseQuery() *gorm.DB
}

// ContactRepository defines methods for contact submission repository
type ContactRepository interface {
	Create(submission *model.ContactSubmission) error
	GetAll(page, pageSize int) (*paging.Paginator, error)
	MarkAsRead(id uint) error
	Delete(id uint) error
	GetBaseQuery() *gorm.DB
}
