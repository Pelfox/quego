package quego

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/Pelfox/quego/internal"
	"github.com/Pelfox/quego/internal/dto"
	"github.com/Pelfox/quego/internal/repositories"
	"github.com/Pelfox/quego/internal/services"
	"github.com/Pelfox/quego/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/mattn/go-sqlite3"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

// ServerConfig holds configuration options for the Server.
type ServerConfig struct {
	// RedisAddr is the address of the Redis server.
	RedisOptions *redis.Options
	// WorkersCount is the number of concurrent workers to process function executions.
	WorkersCount int
	// CORSOrigins is the string of allowed origins for CORS.
	CORSOrigins []string
}

// Server represents the HTTP API server. It wires together the Echo instance
// with services and repositories that provide business and persistence logic.
type Server struct {
	app *echo.Echo

	executionService *services.ExecutionService
	triggerService   *services.TriggerService
}

// NewServer initializes and returns a new Server instance. It creates a SQLite
// database connection, configures repositories, and wires up the execution and
// trigger services.
func NewServer(config ServerConfig) (*Server, error) {
	db, err := sqlx.Connect("sqlite3", fmt.Sprintf("file:%s", internal.DatabaseFile))
	if err != nil {
		return nil, err
	}

	redis := redis.NewClient(config.RedisOptions)
	app := echo.New()
	app.HideBanner = true
	app.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: config.CORSOrigins,
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodOptions},
	}))

	return &Server{
		app: app,
		executionService: services.NewExecutionService(
			config.WorkersCount,
			redis,
			repositories.NewExecutionRepository(db),
		),
		triggerService: services.NewTriggerService(
			repositories.NewTriggerRepository(db),
		),
	}, nil
}

// RegisterFunction registers a function with the ExecutionService. Registered
// functions can later be invoked via triggers.
func (s *Server) RegisterFunction(name string, f models.ExecFunction) {
	s.executionService.RegisterFunction(name, f)
}

// triggerRoute handles `POST /trigger` requests.
// Flow:
//  1. Client submits a trigger in the request body.
//  2. The trigger is parsed and saved in the database.
//  3. The corresponding function is executed.
func (s *Server) triggerRoute(ctx echo.Context) error {
	var triggerPayload dto.CreateTriggerDTO
	if err := ctx.Bind(&triggerPayload); err != nil {
		return internal.RespondError(
			ctx,
			http.StatusBadRequest,
			internal.ErrorCodeInvalidBody,
			"Failed to parse request body",
		)
	}

	trigger := models.Trigger{
		TriggerType:  models.TriggerTypeEvent,
		FunctionName: triggerPayload.FunctionName,
		Payload:      triggerPayload.Payload,
	}
	if err := s.triggerService.Create(&trigger); err != nil {
		log.Error().Err(err).Msg("failed to create trigger")
		return internal.RespondError(
			ctx,
			http.StatusInternalServerError,
			internal.ErrorCodeDatabase,
			"Failed to create trigger",
		)
	}

	execution, err := s.executionService.Process(&trigger)
	if err != nil {
		if errors.Is(err, services.ErrFunctionNotFound) {
			return internal.RespondError(
				ctx,
				http.StatusBadRequest,
				internal.ErrorCodeFunctionNotFound,
				"The requested function is not registered",
			)
		}
		log.Error().Err(err).Str("id", trigger.ID.String()).Msg("failed to enqueue trigger")
		return internal.RespondError(
			ctx,
			http.StatusInternalServerError,
			internal.ErrorCodeDatabase,
			"Failed to process trigger",
		)
	}

	return ctx.JSON(http.StatusOK, execution)
}

// getExecution handles `GET /executions/:id` requests. It retrieves a single
// execution by its UUID and returns it as JSON.
func (s *Server) getExecution(ctx echo.Context) error {
	id := ctx.Param("id")
	if uuid.Validate(id) != nil {
		return internal.RespondError(
			ctx,
			http.StatusBadRequest,
			internal.ErrorCodeInvalidBody,
			"Invalid execution ID",
		)
	}

	executionID, _ := uuid.Parse(id)
	execution, err := s.executionService.GetByID(executionID)
	if err != nil {
		log.Error().Err(err).Str("id", id).Msg("failed to get execution")
		return internal.RespondError(
			ctx,
			http.StatusInternalServerError,
			internal.ErrorCodeDatabase,
			"Failed to retrieve execution",
		)
	}
	if execution == nil {
		return internal.RespondError(
			ctx,
			http.StatusNotFound,
			internal.ErrorCodeInvalidBody,
			"Execution not found",
		)
	}

	return ctx.JSON(http.StatusOK, execution)
}

// ListExecutions handles `GET /executions` requests. It retrieves all
// executions and returns them as JSON.
func (s *Server) ListExecutions(ctx echo.Context) error {
	executions, err := s.executionService.ListAllTriggers()
	if err != nil {
		log.Error().Err(err).Msg("failed to retrieve executions")
		return internal.RespondError(
			ctx,
			http.StatusInternalServerError,
			internal.ErrorCodeDatabase,
			"Failed to retrieve executions",
		)
	}
	return ctx.JSON(http.StatusOK, executions)
}

// Start runs the HTTP server at the given address. Before starting,
// it ensures the database schema is migrated.
func (s *Server) Start(addr string) error {
	if err := internal.MigrateDatabase(); err != nil {
		return err
	}
	if err := s.executionService.RequeueStaled(); err != nil {
		return err
	}

	// logging middleware
	s.app.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			err := next(c)
			log.Info().Str("path", c.Request().RequestURI).
				Str("method", c.Request().Method).
				Str("ip", c.RealIP()).
				Int("status", c.Response().Status).
				Msg("processing request")
			return err
		}
	})

	s.executionService.StartWorkers(context.Background())
	s.app.POST("/trigger", s.triggerRoute)
	s.app.GET("/executions", s.ListExecutions)
	s.app.GET("/executions/:id", s.getExecution)
	return s.app.Start(addr)
}
