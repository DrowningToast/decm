package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	customerror "apps/backend/common/customerror"
	"apps/backend/common/log"
	"apps/backend/common/pgclient"
	"apps/backend/core-api/config"
	auth_handler "apps/backend/core-api/internal/handler/auth"
	"apps/backend/core-api/internal/handler/event"
	"apps/backend/core-api/internal/handler/event_registration_invitation"
	eventconfig_handler "apps/backend/core-api/internal/handler/eventconfig"
	"apps/backend/core-api/internal/handler/issuer"
	"apps/backend/core-api/internal/handler/onboard"
	"apps/backend/core-api/internal/handler/profile"
	authenticationguard "apps/backend/core-api/internal/middleware/authentication_guard"
	verifyjwt "apps/backend/core-api/internal/middleware/verify_jwt"
	"apps/backend/core-api/internal/repositories/postgres"
	event_usecase "apps/backend/core-api/internal/usecase/event"
	event_registration_invitation_usecase "apps/backend/core-api/internal/usecase/event_registration_invitation"
	eventconfig_usecase "apps/backend/core-api/internal/usecase/eventconfig"
	issuer_usecase "apps/backend/core-api/internal/usecase/issuer"
	oauth_usecase "apps/backend/core-api/internal/usecase/oauth"
	onboard_usecase "apps/backend/core-api/internal/usecase/onboard"
	profile_usecase "apps/backend/core-api/internal/usecase/profile"
	"apps/backend/services/auth"
	"apps/backend/services/oauth"
	"apps/backend/services/s3"

	json "github.com/goccy/go-json"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/fiber/v2/middleware/timeout"
	"github.com/gofiber/swagger"

	// fiber-swagger middleware
	_ "apps/backend/core-api/docs"
)

// @title DECM Core
// @version 1.0
// @description DECM (Decentralized Event Management) platform API for NFT ticketing, digital credentials, and academic identity verification.
// @termsOfService http://swagger.io/terms/
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8080
// @BasePath /
func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfg := config.LoadConfig()

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		logger := log.LoadLogger()
		logger.ErrorContext(ctx, "Configuration validation failed", slog.String("error", err.Error()))
		os.Exit(1)
	}

	logger := log.LoadLogger()

	pgConn, err := pgclient.NewPool(ctx, &cfg.Postgres)
	if err != nil {
		logger.ErrorContext(ctx, "failed to create pgxpool", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer func() {
		pgConn.Close()
		logger.InfoContext(ctx, "Gracefully closed pgxpool connection")
	}()
	logger.Info("Sucessfully connected to pg pool")

	// services
	expiration, err := time.ParseDuration(cfg.Jwt.Expiration)
	if err != nil {
		logger.ErrorContext(ctx, "failed to parse jwt expiration", slog.String("error", err.Error()))
		os.Exit(1)
	}
	authService := auth.NewAuthService(cfg.Jwt.Issuer, cfg.Jwt.SecretKey, expiration)
	googleOAuthService := oauth.NewGoogleOAuthService()
	s3Service, err := s3.NewS3Service()
	if err != nil {
		logger.Error("Failed to initialize S3 service", "error", err)
		panic(fmt.Sprintf("S3 service initialization failed: %v", err))
	}

	// repo
	pgRepo := postgres.NewRepository(pgConn, cfg.PIIEncryptionKey)

	onboardUc := onboard_usecase.NewOnboardUsecase(pgRepo, pgRepo, authService, googleOAuthService)
	oauthUc := oauth_usecase.NewOAuthUsecase(googleOAuthService, pgRepo)
	profileUc := profile_usecase.NewProfileUsecase(pgRepo, pgRepo)
	eventUc := event_usecase.NewEventUsecase(pgRepo, pgRepo, pgRepo, pgRepo, pgRepo, pgRepo, pgRepo, s3Service, logger, authService)
	eventConfigUc := eventconfig_usecase.NewEventConfigUsecase(pgRepo, pgRepo, pgRepo, pgRepo, *s3Service, logger)
	issuerUc := issuer_usecase.NewIssuerUsecase(pgRepo)
	eventRegistrationInvitationUc := event_registration_invitation_usecase.NewEventRegistrationInvitationUsecase(pgRepo, pgRepo, pgRepo)

	// Setup HTTP server
	app := fiber.New(fiber.Config{
		AppName:            cfg.Name,
		JSONEncoder:        json.Marshal, // Optimize JSON encoding with go-json
		JSONDecoder:        json.Unmarshal,
		ReadBufferSize:     cfg.Api.MaxReadBufferSize,
		EnableIPValidation: true,
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			return customerror.GetErrFiberHandler(logger)(ctx, err)
		},
	})

	app.Use(favicon.New()).
		Use(cors.New(cors.Config{
			AllowOrigins:     cfg.CorsAllowedOrigins,
			AllowCredentials: true,
		})).
		Use(requestid.New()).
		Use(recover.New(recover.Config{
			EnableStackTrace: true,
			StackTraceHandler: func(c *fiber.Ctx, e interface{}) {
				buf := make([]byte, 1024) // bufLen = 1024
				buf = buf[:runtime.Stack(buf, false)]
				logger.ErrorContext(c.UserContext(), "Something went wrong, panic in http handler", slog.Any("panic", e), slog.String("stacktrace", string(buf)))
			},
		})).
		Use(healthcheck.New(healthcheck.Config{
			LivenessEndpoint: "/",
			LivenessProbe: func(c *fiber.Ctx) bool {
				return true
			},
			ReadinessEndpoint: "/ready",
			ReadinessProbe: func(c *fiber.Ctx) bool {
				if err := pgConn.Ping(ctx); err != nil {
					return false
				}

				return true
			},
		})).
		Use(compress.New(compress.Config{
			Level: compress.LevelDefault,
		})).
		Use(timeout.NewWithContext(func(c *fiber.Ctx) error { return c.Next() }, cfg.Api.Timeout))

	// Swagger
	if config.IsDevelopment() {
		// programmatically set swagger info
		app.Get("/swagger/*", swagger.HandlerDefault) // default
	}

	// API v1
	apiV1 := app.Group("/api/v1")

	authenticationGuardMiddleware := authenticationguard.New(authService)
	verifyJwtMiddleware := verifyjwt.New(authService)

	// Onboard handler
	onboardHandler, err := onboard.NewHandler(onboardUc, profileUc, authService, googleOAuthService, verifyJwtMiddleware)
	if err != nil {
		logger.ErrorContext(ctx, "failed to create onboard handler", slog.String("error", err.Error()))
		os.Exit(1)
	}
	onboardHandler.Mount(apiV1)

	authHandler := auth_handler.NewHandler(oauthUc, googleOAuthService, authService)
	authHandler.Mount(apiV1)

	profileHandler := profile.NewHandler(profileUc, authService, authenticationGuardMiddleware)
	profileHandler.Mount(apiV1)

	eventHandler := event.NewHandler(eventUc, eventConfigUc, profileUc, eventRegistrationInvitationUc, authService, authenticationGuardMiddleware, logger)
	eventHandler.Mount(apiV1)

	eventConfigHandler := eventconfig_handler.NewHandler(eventConfigUc, eventUc, authService, authenticationGuardMiddleware)
	eventConfigHandler.Mount(apiV1)

	issuerHandler := issuer.NewHandler(issuerUc, authService, authenticationGuardMiddleware)
	issuerHandler.Mount(apiV1)

	eventRegistrationInvitationHandler := event_registration_invitation.NewHandler(*authService, eventRegistrationInvitationUc)
	eventRegistrationInvitationHandler.RegisterRoutes(apiV1)

	// Start HTTP Server
	go func() {
		if err := app.Listen(fmt.Sprintf(":%d", cfg.Port)); err != nil {
			logger.ErrorContext(ctx, "error while server listening", slog.String("error", err.Error()))
			stop() // stop app if HTTP server is stopped
		}
	}()

	logger.InfoContext(ctx, "Starting application...",
		slog.String("name", cfg.Name),
		slog.String("env", cfg.Env),
		slog.Int("port", cfg.Port),
		slog.Duration("timeout", cfg.Api.Timeout),
	)

	// Handle slow gracefully shutdown
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		defer logger.InfoContext(ctx, "Gracefully shutting down server")

		go func() {
			<-ctx.Done()
			if ctx.Err() == context.DeadlineExceeded {
				logger.ErrorContext(ctx, "Graceful shutdown timeout, force shutdown", slog.String("error", "graceful shutdown timeout"))
				os.Exit(1)
			}
		}()

		if err := app.ShutdownWithContext(ctx); err != nil {
			logger.ErrorContext(ctx, "error in shutdown http server", slog.String("error", err.Error()))
		} else {
			logger.InfoContext(ctx, "Gracefully stopped HTTP server")
		}
	}()

	// Listen for signals
	<-ctx.Done()
	logger.InfoContext(ctx, "Received signal to terminate application")
}
