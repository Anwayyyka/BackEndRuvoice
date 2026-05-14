package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Anwayyyka/ruvoice-backend/internal/config"
	"github.com/Anwayyyka/ruvoice-backend/internal/delivery/http/handlers"
	authmiddleware "github.com/Anwayyyka/ruvoice-backend/internal/delivery/http/middleware"
	"github.com/Anwayyyka/ruvoice-backend/internal/repository"
	"github.com/Anwayyyka/ruvoice-backend/internal/service"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func main() {
	cfg := config.Load()

	db, err := repository.NewPostgresDB(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Репозитории
	userRepo := repository.NewUserRepository(db)
	artistRequestRepo := repository.NewArtistRequestRepository(db)
	trackRepo := repository.NewTrackRepository(db)
	likeRepo := repository.NewLikeRepository(db)
	favRepo := repository.NewFavoriteRepository(db)

	// Сервисы
	userService := service.NewUserService(userRepo, artistRequestRepo, cfg.JWTSecret)
	trackService := service.NewTrackService(trackRepo, userRepo, likeRepo, favRepo)

	// Хендлеры
	jamendoHandler := handlers.NewJamendoHandler(trackService)
	authHandler := handlers.NewAuthHandler(userService)
	userHandler := handlers.NewUserHandler(userService)
	trackHandler := handlers.NewTrackHandler(trackService, userService)
	moderationHandler := handlers.NewModerationHandler(trackService)

	r := chi.NewRouter()

	// Middleware
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173", "http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Статика для загруженных файлов
	r.Handle("/uploads/*", http.StripPrefix("/uploads/", http.FileServer(http.Dir("./uploads"))))

	// Публичные маршруты
	r.Post("/api/auth/register", authHandler.Register)
	r.Post("/api/auth/login", authHandler.Login)
	r.Get("/api/tracks", trackHandler.ListApproved)
	r.Get("/api/tracks/{id}", trackHandler.GetByID)
	r.Get("/api/users/by-email/{email}", userHandler.GetUserByEmail)
	r.Get("/api/artists/{artistId}/tracks", trackHandler.GetByArtist)
	// Jamendo прокси (добавить эту строку)
	r.Get("/api/jamendo/tracks", jamendoHandler.GetTracks)

	// Защищённые маршруты
	r.Group(func(r chi.Router) {
		r.Use(authmiddleware.AuthMiddleware(cfg.JWTSecret))

		r.Get("/api/profile", authHandler.GetProfile)
		r.Put("/api/profile", userHandler.UpdateProfile)
		r.Post("/api/profile/artist-request", userHandler.RequestArtist)

		// Треки
		r.With(authmiddleware.ArtistMiddleware).Post("/api/tracks", trackHandler.Create)

		// Лайки и избранное
		r.Post("/api/tracks/{id}/play", trackHandler.Play)
		r.Post("/api/tracks/{id}/like", trackHandler.Like)
		r.Delete("/api/tracks/{id}/like", trackHandler.Unlike)
		r.Post("/api/tracks/{id}/favorite", trackHandler.AddFavorite)
		r.Delete("/api/tracks/{id}/favorite", trackHandler.RemoveFavorite)
		r.Get("/api/favorites", trackHandler.GetUserFavorites)
	})

	// Админские маршруты
	r.Group(func(r chi.Router) {
		r.Use(authmiddleware.AuthMiddleware(cfg.JWTSecret))
		r.Use(authmiddleware.AdminMiddleware)

		r.Get("/api/moderation/pending", moderationHandler.ListPending)
		r.Post("/api/moderation/approve", moderationHandler.Approve)
		r.Post("/api/moderation/reject", moderationHandler.Reject)
	})

	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	srv := &http.Server{
		Addr:    ":" + cfg.AppPort,
		Handler: r,
	}

	go func() {
		log.Printf("Server starting on port %s", cfg.AppPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}
	log.Println("Server exited")
}
