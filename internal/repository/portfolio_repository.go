package repository

import (
	"encoding/json"
	"errors"

	"github.com/lutestringamend/perwebbe/internal/model"
	"github.com/lutestringamend/perwebbe/pkg/paging"
	"gorm.io/gorm"
)

type portfolioRepository struct {
	db *gorm.DB
}

func NewPortfolioRepository(db *gorm.DB) PortfolioRepository {
	return &portfolioRepository{db: db}
}

func (r *portfolioRepository) GetBaseQuery() *gorm.DB {
	return r.db.Model(&model.PortfolioProject{})
}

func (r *portfolioRepository) Create(project *model.PortfolioProject) error {
	techJSON, err := json.Marshal(project.Technologies)
	if err != nil {
		return err
	}
	project.TechJSON = string(techJSON)
	return r.db.Create(project).Error
}

func (r *portfolioRepository) GetByID(id uint) (*model.PortfolioProject, error) {
	var project model.PortfolioProject
	err := r.db.First(&project, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	if project.TechJSON != "" {
		if err := json.Unmarshal([]byte(project.TechJSON), &project.Technologies); err != nil {
			return nil, err
		}
	}
	return &project, nil
}

func (r *portfolioRepository) GetAll(projectType string, page, pageSize int) (*paging.Paginator, error) {
	var projects []model.PortfolioProject
	query := r.db.Model(model.PortfolioProject{})

	if projectType != "" {
		query = query.Where("project_type = ?", projectType)
	}
	pagingParam := &paging.Param{
		DB:      query,
		Page:    page,
		Limit:   pageSize,
		OrderBy: []string{"created_at DESC"},
	}

	paginator := paging.Paging(pagingParam, &projects)

	for i := range projects {
		if projects[i].TechJSON != "" {
			if err := json.Unmarshal([]byte(projects[i].TechJSON), &projects[i].Technologies); err != nil {
				return nil, err
			}
		}
	}

	return paginator, nil
}

func (r *portfolioRepository) Update(project *model.PortfolioProject) error {
	techJSON, err := json.Marshal(project.Technologies)
	if err != nil {
		return err
	}
	project.TechJSON = string(techJSON)
	return r.db.Session(&gorm.Session{FullSaveAssociations: true}).Save(project).Error
}

func (r *portfolioRepository) Delete(id uint) error {
	return r.db.Delete(&model.PortfolioProject{}, id).Error
}
