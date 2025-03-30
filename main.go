package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/lutestringamend/perwebbe/internal/config"
	"github.com/lutestringamend/perwebbe/internal/handler"
	"github.com/lutestringamend/perwebbe/internal/repository"
	"github.com/lutestringamend/perwebbe/internal/service"
)

func main() {
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	db, err := config.SetupDatabase(cfg)
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}

	withMigrations := flag.Bool("with-migrations", false, "with db migrations")
	flag.Parse()

	if *withMigrations {
		if err := config.RunMigrations(db); err != nil {
			log.Fatalf("failed to run db migrations : %v", err)
		}
	}

	blogRepo := repository.NewBlogRepository(db)
	portfolioRepo := repository.NewPortfolioRepository(db)
	contactRepo := repository.NewContactRepository(db)

	blogService := service.NewBlogService(blogRepo)
	portfolioService := service.NewPortfolioService(portfolioRepo)
	contactService := service.NewContactService(contactRepo)

	blogHandler := handler.NewBlogHandler(blogService)
	portfolioHandler := handler.NewPortfolioHandler(portfolioService)
	contactHandler := handler.NewContactHandler(contactService)

	router := gin.Default()

	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	api := router.Group("/api")
	{
		blogRoutes := api.Group("/blogs")
		{
			blogRoutes.GET("/", blogHandler.GetAllBlogs)
			blogRoutes.GET("/:slug", blogHandler.GetBlogBySlug)
			blogRoutes.POST("/", blogHandler.CreateBlog)
			blogRoutes.PUT("/:id", blogHandler.UpdateBlog)
			blogRoutes.DELETE("/:id", blogHandler.DeleteBlog)
		}

		portfolioRoutes := api.Group("/portfolio")
		{
			portfolioRoutes.GET("/", portfolioHandler.GetAllProjects)
			portfolioRoutes.GET("/:id", portfolioHandler.GetProject)
			portfolioRoutes.POST("/", portfolioHandler.CreateProject)
			portfolioRoutes.PUT("/:id", portfolioHandler.UpdateProject)
			portfolioRoutes.DELETE("/:id", portfolioHandler.DeleteProject)
		}

		contactRoutes := api.Group("/contacts")
		{
			contactRoutes.GET("/", contactHandler.GetAllContacts)
			contactRoutes.POST("/", contactHandler.CreateContact)
			contactRoutes.PUT("/:id/read", contactHandler.MarkContactAsRead)
			contactRoutes.DELETE("/:id", contactHandler.DeleteContact)
		}
	}

	serverAddr := fmt.Sprintf(":%s", cfg.ServerPort)
	log.Printf("server starting on %s", serverAddr)
	if err := router.Run(serverAddr); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}

}
