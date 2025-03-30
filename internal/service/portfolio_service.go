package service

import (
	"github.com/lutestringamend/perwebbe/internal/model"
	"github.com/lutestringamend/perwebbe/internal/repository"
	"github.com/lutestringamend/perwebbe/pkg/paging"
	"gorm.io/gorm"
)

type portfolioService struct {
	repo repository.PortfolioRepository
}

func NewPortfolioService(repo repository.PortfolioRepository) PortfolioService {
	return &portfolioService{repo: repo}
}

func (s *portfolioService) CreateProject(project *model.PortfolioProject) error {
	return s.repo.Create(project)
}

func (s *portfolioService) GetProjectByID(id uint) (*model.PortfolioProject, error) {
	return s.repo.GetByID(id)
}

func (s *portfolioService) GetAllProjects(projectType string, page, pageSize int) (*paging.Paginator, error) {
	return s.repo.GetAll(projectType, page, pageSize)
}

func (s *portfolioService) UpdateProject(project *model.PortfolioProject) error {
	return s.repo.Update(project)
}

func (s *portfolioService) DeleteProject(id uint) error {
	return s.repo.Delete(id)
}

func (s *portfolioService) GetBlogBaseQuery() *gorm.DB {
	return s.repo.GetBaseQuery()
}
