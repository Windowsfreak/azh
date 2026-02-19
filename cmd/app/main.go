package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"azh/internal/config"
	"azh/internal/handler"
	"azh/internal/model"
	"azh/internal/repository"
	"azh/internal/service"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize database connection
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto-migrate models
	err = db.AutoMigrate(&model.Course{}, &model.Member{}, &model.MemberCourse{}, &model.Participation{})
	if err != nil {
		log.Fatalf("Failed to auto-migrate database: %v", err)
	}

	// Initialize repositories
	courseRepo := repository.NewCourseRepository(db)
	memberCourseRepo := repository.NewMemberCourseRepository(db)
	participationRepo := repository.NewParticipationRepository(db)
	memberRepo := repository.NewMemberRepository(db)

	// Initialize services
	courseService := service.NewCourseService(courseRepo)
	participationService := service.NewParticipationService(courseRepo, memberCourseRepo, participationRepo, memberRepo)
	importService := service.NewImportService(db, courseRepo, memberRepo, memberCourseRepo, participationRepo)

	// Initialize handlers
	courseHandler := handler.NewCourseHandler(courseService)
	participationHandler := handler.NewParticipationHandler(participationService)
	importHandler := handler.NewImportHandler(importService)

	// Set up router
	router := httprouter.New()

	// Course endpoints
	router.GET("/api/courses", courseHandler.GetCourses)
	router.GET("/api/courses/:id/occurrences", courseHandler.GetOccurrences)

	// Participation endpoints
	router.GET("/api/courses/:id/dates/:date/participants", participationHandler.GetParticipants)
	router.POST("/api/courses/:id/dates/:date/participants/:participantId/attendance", participationHandler.SetAttendance)

	// Export endpoint
	router.GET("/api/export", participationHandler.ExportData)

	// Import endpoint
	router.POST("/api/import", importHandler.ImportCSV)

	router.GET("/", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		w.Header().Set("Content-Type", "text/html")
		http.ServeFile(w, r, "index.html")
	})

	// Start server
	port := cfg.Port
	if port == "" {
		port = "8572"
	}
	log.Printf("Server starting on port %s", port)
	err = http.ListenAndServe(":"+port, router)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
