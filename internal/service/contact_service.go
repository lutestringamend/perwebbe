package service

import (
	"github.com/lutestringamend/perwebbe/internal/model"
	"github.com/lutestringamend/perwebbe/internal/repository"
	"github.com/lutestringamend/perwebbe/pkg/paging"
	"gorm.io/gorm"
)

type contactService struct {
	repo repository.ContactRepository
}

func NewContactService(repo repository.ContactRepository) ContactService {
	return &contactService{repo: repo}
}

func (s *contactService) CreateContact(submission *model.ContactSubmission) error {
	return s.repo.Create(submission)
}

func (s *contactService) GetAllContacts(page, pageSize int) (*paging.Paginator, error) {
	return s.repo.GetAll(page, pageSize)
}

func (s *contactService) MarkContactAsRead(id uint) error {
	return s.repo.MarkAsRead(id)
}

func (s *contactService) DeleteContact(id uint) error {
	return s.repo.Delete(id)
}

func (s *contactService) GetBlogBaseQuery() *gorm.DB {
	return s.repo.GetBaseQuery()
}
