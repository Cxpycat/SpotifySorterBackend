package main

import (
	"SpotifySorter/internal/config"
	userHandlers "SpotifySorter/internal/http-server/handlers/user"
	jwtMiddleware "SpotifySorter/internal/http-server/middleware/jwt"
	"SpotifySorter/internal/lib/logger/handlers/slogpretty"
	"SpotifySorter/internal/lib/logger/slog"
	"SpotifySorter/internal/storage/mysql"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	cfg := config.MustLoad()

	logger := setupLogger(cfg.Env)

	logger.Info("Starting application")
	//logger.Debug("Environment:", cfg)

	dbConfig := mysql.Config{
		Host:     cfg.Database.Host,
		Port:     cfg.Database.Port,
		User:     cfg.Database.User,
		Password: cfg.Database.Password,
		Database: cfg.Database.Database,
	}

	storage, err := mysql.Init(dbConfig)
	if err != nil {
		logger.Error("failed to init storage", sl.Err(err))
		os.Exit(1)
	}

	logger.Info("Starting application")
	router := chi.NewRouter()
	logger.Info("Router created")

	router.Use(middleware.RequestID)
	//router.Use(mwLogger.New(logger))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	corsMiddleware(router)

	router.Post("/auth/code", userHandlers.AuthUser(logger, storage))

	router.Group(func(r chi.Router) {
		r.Use(jwtMiddleware.JWTMiddleware(os.Getenv("JWT_SECRET"), storage))
		r.Get("/user/playlist", userHandlers.GetAllPlaylists(logger, storage))
		r.Get("/user/playlist/{id}", userHandlers.GetPlaylistById(logger, storage))
	})

	logger.Info("Starting server")

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	initServer(cfg, router, logger)
	logger.Info("server started")

	<-done
	logger.Info("stopping server")

	logger.Info("server stopped")
}

func setupLogger(env string) *slog.Logger {
	var logger *slog.Logger

	switch env {
	case envLocal:
		logger = setupPrettySlog()
	case envDev:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	default:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return logger
}

func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}

func corsMiddleware(router *chi.Mux) {
	corsOptions := cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	}
	router.Use(cors.New(corsOptions).Handler)
}

func initServer(cfg *config.Config, router *chi.Mux, logger *slog.Logger) {
	{
		srv := &http.Server{
			Addr:         cfg.Address,
			Handler:      router,
			ReadTimeout:  cfg.HTTPServer.Timeout,
			WriteTimeout: cfg.HTTPServer.Timeout,
			IdleTimeout:  cfg.HTTPServer.IdleTimeout,
		}
		go func() {
			if err := srv.ListenAndServe(); err != nil {
				logger.Error("failed to start server")
			}
		}()
	}
}
