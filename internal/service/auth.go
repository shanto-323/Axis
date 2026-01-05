package service

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/shanto-323/axis/config"
	"github.com/shanto-323/axis/internal/database"
	"github.com/shanto-323/axis/internal/errs"
	"github.com/shanto-323/axis/internal/model/dto"
	"github.com/shanto-323/axis/internal/server/middleware"
	"github.com/shanto-323/axis/pkg"
	"go.opentelemetry.io/otel/trace"
)

type AuthService interface {
	Login(c echo.Context, payload *dto.LoginRequest) (*dto.AuthResponse, error)
	Register(c echo.Context, payload *dto.RegisterRequest) (*dto.AuthResponse, error)
}

type authService struct {
	cfg    *config.Config
	db     database.Database
	tracer trace.Tracer
}

func NewAuthService(cfg *config.Config, db database.Database, tracer trace.Tracer) AuthService {
	return &authService{
		cfg:    cfg,
		db:     db,
		tracer: tracer,
	}
}

func (a *authService) Login(c echo.Context, payload *dto.LoginRequest) (*dto.AuthResponse, error) {
	ctx := c.Request().Context()

	ctx, span := a.tracer.Start(ctx, "service")
	defer span.End()

	c.SetRequest(c.Request().WithContext(ctx))

	user, err := a.db.GetUserByEmail(context.Background(), payload.Email)
	if err != nil {
		return nil, err
	}

	if pkg.CompareWithHash(user.PasswordHash, payload.Password) != nil {
		return nil, errs.NewForbiddenError(
			"Invalid credentials",
			true,
		)
	}

	token, err := pkg.CreateAccessToken(a.cfg, user.ID)
	if err != nil {
		return nil, errs.NewInternalServerError()
	}

	return &dto.AuthResponse{
		AccessToken: token,
	}, nil
}

func (a *authService) Register(c echo.Context, payload *dto.RegisterRequest) (*dto.AuthResponse, error) {
	ctx := c.Request().Context()

	ctx, span := a.tracer.Start(ctx, "service")
	defer span.End()

	c.SetRequest(c.Request().WithContext(ctx))

	logger := middleware.GetLogger(c)

	user, err := a.db.GetUserByEmail(context.Background(), payload.Email)
	_, ok := err.(*errs.HTTPError)
	if !ok {
		code := "USER_ALREADY_EXISTS"
		return nil, errs.NewBadRequestError("email already exists", false, &code, nil, nil)
	}

	if user != nil {
		return nil, err
	}

	hashPassword, err := pkg.CreateHash(payload.Password)
	if err != nil {
		return nil, errs.NewInternalServerError()
	}

	// Swap with hash one
	payload.Password = hashPassword

	user, err = a.db.CreateUser(context.Background(), payload)
	if err != nil {
		return nil, err
	}

	token, err := pkg.CreateAccessToken(a.cfg, user.ID)
	if err != nil {
		return nil, errs.NewInternalServerError()
	}

	logger.Info().
		Str("event", "login").
		Any("user", user.ID).
		Msg("registered successfully")

	return &dto.AuthResponse{
		AccessToken: token,
	}, nil
}
