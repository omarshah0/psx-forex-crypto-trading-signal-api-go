package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/omarshah0/rest-api-with-social-auth/internal/config"
	"github.com/omarshah0/rest-api-with-social-auth/internal/database"
	"github.com/omarshah0/rest-api-with-social-auth/internal/handlers"
	"github.com/omarshah0/rest-api-with-social-auth/internal/middleware"
	"github.com/omarshah0/rest-api-with-social-auth/internal/repositories"
	"github.com/omarshah0/rest-api-with-social-auth/internal/services"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	log.Printf("Starting server in %s mode on port %s", cfg.Server.Environment, cfg.Server.Port)

	// Initialize databases
	postgresDB, err := database.NewPostgresDB(cfg.Database.PostgresURL)
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}
	defer postgresDB.Close()

	mongoDB, err := database.NewMongoDB(cfg.Database.MongoDBURL)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer mongoDB.Close()

	redisDB, err := database.NewRedisDB(cfg.Database.RedisURL, cfg.Database.RedisDB)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer redisDB.Close()

	// Initialize repositories
	userRepo := repositories.NewUserRepository(postgresDB.DB)
	adminRepo := repositories.NewAdminRepository(postgresDB.DB)
	oauthProviderRepo := repositories.NewOAuthProviderRepository(postgresDB.DB)
	tradingSignalRepo := repositories.NewTradingSignalRepository(postgresDB.DB)
	logRepo := repositories.NewLogRepository(mongoDB.Database)

	// Initialize services
	jwtService := services.NewJWTService(&cfg.JWT, redisDB)
	passwordService := services.NewPasswordService()
	emailService := services.NewEmailService(cfg.Email.ServiceEnabled, cfg.Email.FrontendURL, cfg.Email.FromAddress, cfg.Email.FromName)
	oauthService := services.NewOAuthService(
		cfg.OAuth.Google.ClientID,
		cfg.OAuth.Google.ClientSecret,
		cfg.OAuth.Google.RedirectURL,
		cfg.OAuth.Google.Enabled,
		cfg.OAuth.Facebook.ClientID,
		cfg.OAuth.Facebook.ClientSecret,
		cfg.OAuth.Facebook.RedirectURL,
		cfg.OAuth.Facebook.Enabled,
	)
	authService := services.NewAuthService(userRepo, oauthProviderRepo, adminRepo, jwtService, oauthService, passwordService, emailService)
	tradingSignalService := services.NewTradingSignalService(tradingSignalRepo)

	// Initialize middleware
	authMiddleware := middleware.NewAuthMiddleware(jwtService)
	adminMiddleware := middleware.NewAdminMiddleware(adminRepo)
	loggingMiddleware := middleware.NewLoggingMiddleware(logRepo, cfg.Logging)
	rateLimitMiddleware := middleware.NewRateLimitMiddleware(redisDB)

	// Initialize handlers
	healthHandler := handlers.NewHealthHandler(postgresDB, mongoDB, redisDB)
	authHandler := handlers.NewAuthHandler(authService, oauthService, cfg.Cookie)
	emailAuthHandler := handlers.NewEmailAuthHandler(authService, cfg.Cookie)
	tradingSignalHandler := handlers.NewTradingSignalHandler(tradingSignalService)
	profileHandler := handlers.NewProfileHandler(authService)

	// Setup router
	router := mux.NewRouter()

	// Apply global middleware
	router.Use(middleware.CORS)
	router.Use(loggingMiddleware.LogRequests)
	router.Use(rateLimitMiddleware.RateLimit(100)) // 100 requests per minute per IP

	// Enable OPTIONS for all routes (CORS preflight)
	router.Methods("OPTIONS").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Health check endpoint (no auth required)
	router.HandleFunc("/health", healthHandler.HealthCheck).Methods("GET", "OPTIONS")

	// Auth routes
	authRouter := router.PathPrefix("/auth").Subrouter()

	// Email/Password Authentication (can be disabled via config)
	if cfg.Auth.EmailPasswordEnabled {
		authRouter.HandleFunc("/register", emailAuthHandler.Register).Methods("POST")
		authRouter.HandleFunc("/login", emailAuthHandler.Login).Methods("POST")
		authRouter.HandleFunc("/verify-email", emailAuthHandler.VerifyEmail).Methods("GET")
		authRouter.HandleFunc("/forgot-password", emailAuthHandler.ForgotPassword).Methods("POST")
		authRouter.HandleFunc("/reset-password", emailAuthHandler.ResetPassword).Methods("POST")
		authRouter.HandleFunc("/resend-verification", emailAuthHandler.ResendVerification).Methods("POST")
		authRouter.Handle("/change-password", authMiddleware.Authenticate(http.HandlerFunc(emailAuthHandler.ChangePassword))).Methods("POST")
	}

	// Token endpoints
	authRouter.HandleFunc("/refresh", authHandler.Refresh).Methods("POST")
	authRouter.Handle("/logout", authMiddleware.Authenticate(http.HandlerFunc(authHandler.Logout))).Methods("POST")

	// Code Exchange endpoints (for React Web & React Native with code flow)
	authRouter.HandleFunc("/google/exchange", authHandler.ExchangeGoogleCode).Methods("POST")
	authRouter.HandleFunc("/facebook/exchange", authHandler.ExchangeFacebookCode).Methods("POST")

	// ID Token verification endpoints (for React Native/Expo with SDK flow)
	authRouter.HandleFunc("/google/verify", authHandler.VerifyGoogleIDToken).Methods("POST")
	authRouter.HandleFunc("/facebook/verify", authHandler.VerifyFacebookAccessToken).Methods("POST")

	// API routes (authenticated)
	apiRouter := router.PathPrefix("/api").Subrouter()
	apiRouter.Use(authMiddleware.Authenticate)

	// Profile route
	apiRouter.HandleFunc("/profile", profileHandler.GetProfile).Methods("GET")

	// Trading signals routes
	signalsRouter := apiRouter.PathPrefix("/trading-signals").Subrouter()

	// Public routes (authenticated users can read)
	signalsRouter.HandleFunc("", tradingSignalHandler.GetAll).Methods("GET")
	signalsRouter.HandleFunc("/{id}", tradingSignalHandler.GetByID).Methods("GET")

	// Admin-only routes
	adminSignalsRouter := signalsRouter.NewRoute().Subrouter()
	adminSignalsRouter.Use(adminMiddleware.RequireAdmin)
	adminSignalsRouter.HandleFunc("", tradingSignalHandler.Create).Methods("POST")
	adminSignalsRouter.HandleFunc("/{id}", tradingSignalHandler.Update).Methods("PUT")
	adminSignalsRouter.HandleFunc("/{id}", tradingSignalHandler.Delete).Methods("DELETE")

	// Create HTTP server
	srv := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine
	go func() {
		log.Printf("Server listening on http://localhost:%s", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited gracefully")
}
