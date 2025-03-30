package service

import (
	"github.com/lutestringamend/perwebbe/internal/model"
	"github.com/lutestringamend/perwebbe/pkg/paging"
	"gorm.io/gorm"
)

// BlogService defines methods for blog service
type BlogService interface {
	CreateBlog(post *model.BlogPost) error
	GetBlogByID(id uint) (*model.BlogPost, error)
	GetAllBlogs(page, pageSize int) (*paging.Paginator, error)
	GetBlogBySlug(slug string) (*model.BlogPost, error)
	UpdateBlog(post *model.BlogPost) error
	DeleteBlog(id uint) error
	GetBlogBaseQuery() *gorm.DB
}

// PortfolioService defines methods for portfolio service
type PortfolioService interface {
	CreateProject(project *model.PortfolioProject) error
	GetProjectByID(id uint) (*model.PortfolioProject, error)
	GetAllProjects(projectType string, page, pageSize int) (*paging.Paginator, error)
	UpdateProject(project *model.PortfolioProject) error
	DeleteProject(id uint) error
	GetBlogBaseQuery() *gorm.DB
}

// ContactService defines methods for contact service
type ContactService interface {
	CreateContact(submission *model.ContactSubmission) error
	GetAllContacts(page, pageSize int) (*paging.Paginator, error)
	MarkContactAsRead(id uint) error
	DeleteContact(id uint) error
	GetBlogBaseQuery() *gorm.DB
}
