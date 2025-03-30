package model

import (
	"time"

	"gorm.io/gorm"
)

type BlogPost struct {
	gorm.Model
	Title     string    `json:"title" gorm:"not null"`
	Content   string    `json:"content" gorm:"type:text;not null"`
	Summary   string    `json:"summary" gorm:"type:text"`
	ImageURL  string    `json:"image_url"`
	Published bool      `json:"published" gorm:"default:false"`
	PublishAt time.Time `json:"publish_at"`
	Slug      string    `json:"slug" gorm:"uniqueIndex"`
	Tags      []Tag     `json:"tags" gorm:"many2many:blog_tags;"`
}

type Tag struct {
	gorm.Model
	Name      string     `json:"name" gorm:"uniqueIndex;not null"`
	BlogPosts []BlogPost `json:"blog_posts" gorm:"many2many:blog_tags;"`
}

type PortfolioProject struct {
	gorm.Model
	Title        string    `json:"title" gorm:"not null"`
	Description  string    `json:"description" gorm:"type:text;not null"`
	ProjectType  string    `json:"project_type" gorm:"not null"` // "coding" or "music"
	ImageURL     string    `json:"image_url"`
	ProjectURL   string    `json:"project_url"`
	RepoURL      string    `json:"repo_url"`
	Technologies []string  `json:"technologies" gorm:"-"` // Will be handled with JSON serialization
	TechJSON     string    `json:"-" gorm:"column:technologies;type:json"`
	Featured     bool      `json:"featured" gorm:"default:false"`
	StartDate    time.Time `json:"start_date"`
	EndDate      time.Time `json:"end_date"`
}

type ContactSubmission struct {
	gorm.Model
	Name    string `json:"name" gorm:"not null"`
	Email   string `json:"email" gorm:"not null"`
	Subject string `json:"subject"`
	Message string `json:"message" gorm:"type:text;not null"`
	Read    bool   `json:"read" gorm:"default:false"`
}
