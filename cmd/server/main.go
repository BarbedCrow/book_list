package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/BarbedCrow/book_list/internal/auth"
	"github.com/BarbedCrow/book_list/internal/handler"
	"github.com/BarbedCrow/book_list/internal/hasher"
	"github.com/BarbedCrow/book_list/internal/monitor"
	"github.com/BarbedCrow/book_list/internal/postgres"
	authoruc "github.com/BarbedCrow/book_list/internal/usecase/author"
	bookuc "github.com/BarbedCrow/book_list/internal/usecase/book"
	listuc "github.com/BarbedCrow/book_list/internal/usecase/list"
	useruc "github.com/BarbedCrow/book_list/internal/usecase/user"
)

func main() {
	if err := run(); err != nil {
		slog.Error("fatal", "error", err)
		os.Exit(1)
	}
}

func run() error {
	// Config from env
	addr := envOrDefault("SERVER_ADDR", ":8080")
	dbDSN := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		requireEnv("DB_HOST"),
		envOrDefault("DB_PORT", "5432"),
		requireEnv("DB_USER"),
		requireEnv("DB_PASSWORD"),
		requireEnv("DB_NAME"),
		envOrDefault("DB_SSLMODE", "disable"),
	)

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return fmt.Errorf("JWT_SECRET environment variable is required")
	}
	jwtTTL, err := time.ParseDuration(envOrDefault("JWT_TTL", "24h"))
	if err != nil {
		return fmt.Errorf("parse JWT_TTL: %w", err)
	}

	// Database
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	pool, err := pgxpool.New(ctx, dbDSN)
	if err != nil {
		return fmt.Errorf("connect to database: %w", err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		return fmt.Errorf("ping database: %w", err)
	}
	slog.Info("connected to database")

	// Metrics
	reg := prometheus.NewRegistry()
	reg.MustRegister(prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{}))
	reg.MustRegister(prometheus.NewGoCollector())
	reg.MustRegister(monitor.NewPgxPoolCollector(pool))
	handler.RegisterMetrics(reg)

	// Adapters
	bookRepo := postgres.NewBookRepo(pool)
	authorRepo := postgres.NewAuthorRepo(pool)
	userRepo := postgres.NewUserRepo(pool)
	listRepo := postgres.NewListRepo(pool)

	bcryptHasher := hasher.NewBcryptHasher(10)
	jwtProvider := auth.NewJWTProvider(jwtSecret, jwtTTL)

	idGen := func() string { return uuid.New().String() }

	// Use cases
	searchBooks := bookuc.NewSearchBooksByTitle(bookRepo)
	getBookDetails := bookuc.NewGetBookDetails(bookRepo)

	searchAuthors := authoruc.NewSearchAuthorsByName(authorRepo)
	getAuthorDetails := authoruc.NewGetAuthorDetails(authorRepo)
	getBooksByAuthor := authoruc.NewGetBooksByAuthor(authorRepo)

	registerUser := useruc.NewRegisterUser(userRepo, bcryptHasher, idGen)
	authenticateUser := useruc.NewAuthenticateUser(userRepo, bcryptHasher, jwtProvider)

	createList := listuc.NewCreateCustomList(listRepo, idGen)
	getUserLists := listuc.NewGetUserLists(listRepo)
	addBookToList := listuc.NewAddBookToList(listRepo)
	removeBookFromList := listuc.NewRemoveBookFromList(listRepo)

	// Handlers
	bookHandler := handler.NewBookHandler(searchBooks, getBookDetails)
	authorHandler := handler.NewAuthorHandler(searchAuthors, getAuthorDetails, getBooksByAuthor)
	userHandler := handler.NewUserHandler(registerUser, authenticateUser)
	listHandler := handler.NewListHandler(createList, getUserLists, addBookToList, removeBookFromList)

	// Routes
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", handler.HealthHandler(pool))
	mux.Handle("GET /metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))

	mux.HandleFunc("GET /books", bookHandler.Search)
	mux.HandleFunc("GET /books/{id}", bookHandler.GetDetails)

	mux.HandleFunc("GET /authors", authorHandler.Search)
	mux.HandleFunc("GET /authors/{id}", authorHandler.GetDetails)
	mux.HandleFunc("GET /authors/{id}/books", authorHandler.GetBooks)

	mux.HandleFunc("POST /register", userHandler.Register)
	mux.HandleFunc("POST /login", userHandler.Login)

	// Protected routes
	authMW := handler.AuthMiddleware(jwtProvider)
	listMux := http.NewServeMux()
	listMux.HandleFunc("GET /lists", listHandler.GetUserLists)
	listMux.HandleFunc("POST /lists", listHandler.Create)
	listMux.HandleFunc("POST /lists/{id}/books", listHandler.AddBook)
	listMux.HandleFunc("DELETE /lists/{id}/books/{book_id}", listHandler.RemoveBook)

	mux.Handle("/lists", authMW(listMux))
	mux.Handle("/lists/", authMW(listMux))

	// Server
	srv := &http.Server{
		Addr:         addr,
		Handler:      handler.MetricsMiddleware(handler.SecurityHeaders(logRequests(mux))),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		<-ctx.Done()
		slog.Info("shutting down server")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		srv.Shutdown(shutdownCtx)
	}()

	slog.Info("server starting", "addr", addr)
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		return fmt.Errorf("server: %w", err)
	}

	return nil
}

func envOrDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func requireEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		slog.Error("required environment variable is not set", "key", key)
		os.Exit(1)
	}
	return v
}

func logRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.Info("request", "method", r.Method, "path", r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
