package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"google.golang.org/api/iterator"
	"github.com/task-manager/task-service/config"
	"github.com/task-manager/task-service/internal/auth"
	"github.com/task-manager/task-service/internal/projects"
	"github.com/task-manager/task-service/internal/tasks"
	"github.com/task-manager/task-service/internal/users"
	firestoreClient "github.com/task-manager/task-service/pkg/firestore"
	"github.com/task-manager/task-service/pkg/middleware"
	"github.com/task-manager/task-service/pkg/models"
)

func main() {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Load configuration
	cfg := config.Load()

	// Initialize Firestore client
	fsClient := firestoreClient.NewClient(cfg.GCPProjectID)
	defer fsClient.Close()

	// Initialize repositories
	authRepo := auth.NewRepository(fsClient)
	taskRepo := tasks.NewRepository(fsClient)
	projectRepo := projects.NewRepository(fsClient)
	userRepo := users.NewRepository(fsClient)

	// Initialize services
	authService := auth.NewService(authRepo, cfg.JWTSecret, cfg.JWTExpiration)
	taskService := tasks.NewService(taskRepo, projectRepo)
	projectService := projects.NewService(projectRepo)
	userService := users.NewService(userRepo)

	// Initialize handlers
	authHandler := auth.NewHandler(authService)
	taskHandler := tasks.NewHandler(taskService)
	projectHandler := projects.NewHandler(projectService)
	userHandler := users.NewHandler(userService)

	// Setup router
	mux := http.NewServeMux()

	// Auth routes (public)
	mux.HandleFunc("POST /auth/register", authHandler.Register)
	mux.HandleFunc("POST /auth/login", authHandler.Login)

	// Protected routes
	authMiddleware := middleware.Auth(cfg.JWTSecret)
	adminOnly := middleware.RequireRole(models.RoleAdmin)

	// Task routes
	mux.Handle("GET /projects/{projectId}/tasks", authMiddleware(http.HandlerFunc(taskHandler.GetByProject)))
	mux.Handle("GET /projects/{projectId}/tasks/{taskId}", authMiddleware(http.HandlerFunc(taskHandler.GetByID)))
	mux.Handle("POST /projects/{projectId}/tasks", authMiddleware(middleware.RequireRole(models.RoleAdmin, models.RoleMember)(http.HandlerFunc(taskHandler.Create))))
	mux.Handle("PUT /projects/{projectId}/tasks/{taskId}", authMiddleware(middleware.RequireRole(models.RoleAdmin, models.RoleMember)(http.HandlerFunc(taskHandler.Update))))
	mux.Handle("DELETE /projects/{projectId}/tasks/{taskId}", authMiddleware(middleware.RequireRole(models.RoleAdmin, models.RoleMember)(http.HandlerFunc(taskHandler.Delete))))

	// Project routes
	mux.Handle("GET /projects", authMiddleware(http.HandlerFunc(projectHandler.GetByUser)))
	mux.Handle("GET /projects/{projectId}", authMiddleware(http.HandlerFunc(projectHandler.GetByID)))
	mux.Handle("POST /projects", authMiddleware(middleware.RequireRole(models.RoleAdmin, models.RoleMember)(http.HandlerFunc(projectHandler.Create))))
	mux.Handle("PUT /projects/{projectId}", authMiddleware(middleware.RequireRole(models.RoleAdmin, models.RoleMember)(http.HandlerFunc(projectHandler.Update))))
	mux.Handle("DELETE /projects/{projectId}", authMiddleware(middleware.RequireRole(models.RoleAdmin)(http.HandlerFunc(projectHandler.Delete))))
	mux.Handle("POST /projects/{projectId}/members", authMiddleware(middleware.RequireRole(models.RoleAdmin, models.RoleMember)(http.HandlerFunc(projectHandler.AddMember))))
	mux.Handle("DELETE /projects/{projectId}/members/{userId}", authMiddleware(middleware.RequireRole(models.RoleAdmin, models.RoleMember)(http.HandlerFunc(projectHandler.RemoveMember))))

	// User routes (admin only)
	mux.Handle("GET /users", authMiddleware(adminOnly(http.HandlerFunc(userHandler.GetAll))))
	mux.Handle("PUT /users/{userId}/role", authMiddleware(adminOnly(http.HandlerFunc(userHandler.UpdateRole))))

	// Health check with Firestore connectivity
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// Verify Firestore connection
		iter := fsClient.Collections(r.Context())
		_, err := iter.Next()
		// err == iterator.Done is fine, means DB is reachable but empty
		if err != nil && err != iterator.Done {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte(`{"status":"unhealthy","firestore":"disconnected"}`))
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok","firestore":"connected"}`))
	})

	// Apply global middleware
	corsConfig := middleware.DefaultCORSConfig()
	handler := middleware.Recovery(middleware.Logger(middleware.CORS(corsConfig)(mux)))

	// Determine port
	port := cfg.Port
	if envPort := os.Getenv("PORT"); envPort != "" {
		port = envPort
	}

	// Configure server
	server := &http.Server{
		Addr:         ":" + port,
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Graceful shutdown
	go func() {
		log.Printf("Task Service running on port %s", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server stopped gracefully")
}
