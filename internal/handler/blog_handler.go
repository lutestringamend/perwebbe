package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lutestringamend/perwebbe/internal/model"
	"github.com/lutestringamend/perwebbe/internal/service"
)

type BlogHandler struct {
	service service.BlogService
}

func NewBlogHandler(service service.BlogService) *BlogHandler {
	return &BlogHandler{service: service}
}

func (h *BlogHandler) CreateBlog(c *gin.Context) {
	var blog model.BlogPost
	if err := c.ShouldBindJSON(&blog); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.CreateBlog(&blog); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, blog)
}

func (h *BlogHandler) GetBlog(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID format"})
		return
	}

	blog, err := h.service.GetBlogByID(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if blog == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Blog post not found"})
		return
	}
	c.JSON(http.StatusOK, blog)
}

func (h *BlogHandler) GetBlogBySlug(c *gin.Context) {
	slug := c.Param("slug")
	if slug == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid slug"})
		return
	}

	blog, err := h.service.GetBlogBySlug(slug)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if blog == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Blog post not found"})
		return
	}
	c.JSON(http.StatusOK, blog)
}

func (h *BlogHandler) GetAllBlogs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	paginator, err := h.service.GetAllBlogs(page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, paginator)
}

func (h *BlogHandler) UpdateBlog(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID format"})
		return
	}

	blog, err := h.service.GetBlogByID(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if blog == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Blog post not found"})
		return
	}

	var updatedBlog model.BlogPost
	if err := c.ShouldBindJSON(&updatedBlog); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedBlog.ID = blog.ID
	if err := h.service.UpdateBlog(&updatedBlog); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedBlog)
}

func (h *BlogHandler) DeleteBlog(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID format"})
		return
	}

	if err := h.service.DeleteBlog(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "blog post deleted successfully"})
}
