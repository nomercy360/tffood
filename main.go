package main

import (
	"context"
	"errors"
	"github.com/caarlos0/env/v11"
	"github.com/go-playground/validator/v10"
	telegram "github.com/go-telegram/bot"
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"rednit/config"
	"rednit/db"
	"rednit/handler"
	storage "rednit/s3"
	"rednit/terrors"
	"time"
)

func getLoggerMiddleware(logger *slog.Logger) middleware.RequestLoggerConfig {
	return middleware.RequestLoggerConfig{
		LogStatus:   true,
		LogURI:      true,
		LogError:    true,
		HandleError: true,
		LogValuesFunc: func(_ echo.Context, v middleware.RequestLoggerValues) error {
			if v.Error == nil {
				logger.LogAttrs(context.Background(), slog.LevelInfo, "REQUEST",
					slog.String("uri", v.URI),
					slog.Int("status", v.Status),
				)
			} else {
				logger.LogAttrs(context.Background(), slog.LevelError, "REQUEST_ERROR",
					slog.String("uri", v.URI),
					slog.Int("status", v.Status),
					slog.String("err", v.Error.Error()),
				)
			}
			return nil
		},
	}
}

func getServerErrorHandler(e *echo.Echo) func(err error, context2 echo.Context) {
	return func(err error, c echo.Context) {
		var (
			code = http.StatusInternalServerError
			msg  interface{}
		)

		var he *echo.HTTPError
		var terror *terrors.Error
		switch {
		case errors.As(err, &he):
			code = he.Code
			msg = he.Message
		case errors.As(err, &terror):
			code = terror.Code
			msg = terror.Message
		default:
			msg = err.Error()
		}

		if _, ok := msg.(string); ok {
			msg = map[string]interface{}{"error": msg}
		}

		if !c.Response().Committed {
			if c.Request().Method == http.MethodHead {
				err = c.NoContent(code)
			} else {
				err = c.JSON(code, msg)
			}

			if err != nil {
				e.Logger.Error(err)
			}
		}
	}
}

type customValidator struct {
	validator *validator.Validate
}

func (cv *customValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}

func getAuthConfig(secret string) echojwt.Config {
	return echojwt.Config{
		NewClaimsFunc: func(_ echo.Context) jwt.Claims {
			return new(handler.JWTClaims)
		},
		SigningKey:             []byte(secret),
		ContinueOnIgnoredError: true,
		ErrorHandler: func(c echo.Context, err error) error {
			var extErr *echojwt.TokenExtractionError
			if !errors.As(err, &extErr) {
				return echo.NewHTTPError(http.StatusUnauthorized, "auth is invalid")
			}

			claims := &handler.JWTClaims{
				RegisteredClaims: jwt.RegisteredClaims{
					ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour * 30)),
				},
			}

			token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

			c.Set("user", token)

			return nil
		},
	}
}

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	cfg := config.Default{}
	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("Failed to parse config: %v\n", err)
	}

	sql, err := db.ConnectDB(cfg.DBPath)

	if err != nil {
		e.Logger.Fatalf("failed to connect to db: %v", err)
	}

	if err := sql.Migrate(); err != nil {
		e.Logger.Fatalf("failed to migrate db: %v", err)
	}

	s3Client, err := storage.NewS3Client(
		cfg.AWS.AccessKey, cfg.AWS.SecretKey, cfg.AWS.Endpoint, cfg.AWS.Bucket)

	if err != nil {
		log.Fatalf("Failed to initialize AWS S3 client: %v\n", err)
	}

	tg, err := telegram.New(cfg.BotToken)
	if err != nil {
		log.Fatalf("Failed to create telegram bot: %v", err)
	}

	webhook := &telegram.SetWebhookParams{URL: cfg.ExternalURL + "/wh/telegram", MaxConnections: 10}

	if _, err := tg.SetWebhook(context.Background(), webhook); err != nil {
		log.Fatalf("Failed to set webhook: %v", err)
	}

	h := handler.New(sql, cfg, s3Client, tg)

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	}))

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	e.Use(middleware.RequestLoggerWithConfig(getLoggerMiddleware(logger)))

	e.HTTPErrorHandler = getServerErrorHandler(e)

	e.Validator = &customValidator{validator: validator.New()}

	e.POST("/auth/telegram", h.TelegramAuth)
	e.POST("/wh/telegram", h.HandleWebhook)

	// Routes
	g := e.Group("/api")

	authCfg := getAuthConfig(cfg.JWTSecret)

	g.Use(echojwt.WithConfig(authCfg))

	g.GET("/posts", h.GetPosts)
	g.POST("/posts", h.CreatePost)
	g.PUT("/posts/:id", h.UpdatePost)
	g.GET("/posts/:id/ai", h.CreatePostAISuggestions)
	g.GET("/posts/:id", h.GetPost)
	g.GET("/tags", h.GetTags)
	g.POST("/presigned-url", h.GetPresignedURL)
	g.POST("/posts/:id/react/:reaction", h.ReactToPost)
	g.DELETE("/posts/:id/react", h.DropPostReaction)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	// Start server
	go func() {
		if err := e.Start(":8080"); err != nil && !errors.Is(err, http.ErrServerClosed) {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
