package service

import (
	"github.com/lutestringamend/perwebbe/internal/model"
	"github.com/lutestringamend/perwebbe/internal/repository"
	"github.com/lutestringamend/perwebbe/pkg/paging"
	"gorm.io/gorm"
)

type blogService struct {
	repo repository.BlogRepository
}

func NewBlogService(repo repository.BlogRepository) BlogService {
	return &blogService{repo: repo}
}

func (s *blogService) CreateBlog(post *model.BlogPost) error {
	return s.repo.Create(post)
}

func (s *blogService) GetBlogByID(id uint) (*model.BlogPost, error) {
	return s.repo.GetByID(id)
}

func (s *blogService) GetAllBlogs(page, pageSize int) (*paging.Paginator, error) {
	return s.repo.GetAll(page, pageSize)
}

func (s *blogService) GetBlogBySlug(slug string) (*model.BlogPost, error) {
	return s.repo.GetBySlug(slug)
}

func (s *blogService) UpdateBlog(post *model.BlogPost) error {
	return s.repo.Update(post)
}

func (s *blogService) DeleteBlog(id uint) error {
	return s.repo.Delete(id)
}

func (s *blogService) GetBlogBaseQuery() *gorm.DB {
	return s.repo.GetBaseQuery()
}
