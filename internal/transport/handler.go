package transport

import (
	"context"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/skandyla/s3-uploader/internal/models"
)

type Health interface {
	//Liveness(ctx context.Context) error
	Ping(ctx context.Context) error
	Info(ctx context.Context) (models.InfoDependencyItem, error)
}

type User interface {
	SignUp(ctx context.Context, inp models.SignUpInput) error
	SignIn(ctx context.Context, inp models.SignInInput) (string, string, error)
	ParseToken(ctx context.Context, accessToken string) (int64, error)
	RefreshTokens(ctx context.Context, refreshToken string) (string, string, error)
}

type Handler struct {
	healthService Health
	usersService  User
}

func NewHandler(health Health, users User) *Handler {
	return &Handler{
		healthService: health,
		usersService:  users,
	}
}

func (h *Handler) InitRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	//r.Use(middleware.Logger)
	r.Use(loggingMiddleware) //test our own middleware implementation

	// health
	r.Route("/liveness", func(r chi.Router) {
		r.Get("/", h.liveness)
	})
	r.Route("/__service", func(r chi.Router) {
		r.Get("/info", h.info)
	})

	// Users
	r.Route("/auth", func(r chi.Router) {
		r.Post("/sign-up", h.signUp)
		r.Get("/sign-in", h.signIn)
		r.Get("/refresh", h.refresh)
	})

	return r
}
