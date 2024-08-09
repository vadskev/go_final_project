package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/vadskev/go_final_project/internal/handlers/done"
	"github.com/vadskev/go_final_project/internal/handlers/nextdate"
	"github.com/vadskev/go_final_project/internal/handlers/signin"
	"github.com/vadskev/go_final_project/internal/handlers/task"
	"github.com/vadskev/go_final_project/internal/handlers/tasks"
	"github.com/vadskev/go_final_project/internal/logger"
	"github.com/vadskev/go_final_project/internal/middleware/auth"
	mwLogger "github.com/vadskev/go_final_project/internal/middleware/logger"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/vadskev/go_final_project/internal/config"
	"go.uber.org/zap"
)

const (
	taskPath     = "/api/task"
	tasksPath    = "/api/tasks"
	nextDatePath = "/api/nextdate"
	taskDonePath = "/api/task/done"
	singPath     = "/api/signin"
)
const (
	ReadTimeout        = 4 * time.Second
	WriteTimeout       = 4 * time.Second
	IdleTimeout        = 3 * time.Second
	shutDownCtxTimeout = 1 * time.Second
)

type App struct {
	serviceProvider *serviceProvider
	httpServer      *http.Server
}

func NewApp(ctx context.Context) (*App, error) {
	app := &App{}

	err := app.loadDeps(ctx)
	if err != nil {
		return nil, err
	}

	return app, nil
}

func (a *App) loadDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.loadConfig,
		a.loadServiceProvider,
		a.loadLogger,
	}
	for _, f := range inits {
		err := f(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *App) loadConfig(_ context.Context) error {
	err := config.Load()
	if err != nil {
		return err
	}
	return nil
}

func (a *App) loadServiceProvider(_ context.Context) error {
	a.serviceProvider = newServiceProvider()
	return nil
}

func (a *App) loadLogger(_ context.Context) error {
	err := logger.Init(a.serviceProvider.LogConfig().Level())
	if err != nil {
		return err
	}
	return nil
}

func (a *App) RunServer(ctx context.Context) error {
	router := chi.NewRouter()

	router.Use(mwLogger.New())

	router.Route(taskPath, func(r chi.Router) {
		r.Use(auth.New(a.serviceProvider.PassConfig()))
		r.Get("/", task.New(ctx, a.serviceProvider.DBRepository()).HandleGet)
		r.Post("/", task.New(ctx, a.serviceProvider.DBRepository()).HandlePost)
		r.Put("/", task.New(ctx, a.serviceProvider.DBRepository()).HandlePut)
		r.Delete("/", task.New(ctx, a.serviceProvider.DBRepository()).HandleDelete)
	})

	router.Route(tasksPath, func(r chi.Router) {
		r.Use(auth.New(a.serviceProvider.PassConfig()))
		r.Get("/", tasks.New(ctx, a.serviceProvider.DBRepository()).Handle)
	})

	router.Route(nextDatePath, func(r chi.Router) {
		r.Get("/", nextdate.New(ctx, a.serviceProvider.DBRepository()).HandleGet)
	})

	router.Route(taskDonePath, func(r chi.Router) {
		r.Use(auth.New(a.serviceProvider.PassConfig()))
		r.Post("/", done.New(ctx, a.serviceProvider.DBRepository()).HandlePost)
	})

	router.Route(singPath, func(r chi.Router) {
		r.Post("/", signin.New(ctx, a.serviceProvider.DBRepository(), a.serviceProvider.PassConfig()).HandlePost)
	})

	a.httpServer = &http.Server{
		Addr:         a.serviceProvider.HTTPConfig().Address(),
		Handler:      router,
		ReadTimeout:  ReadTimeout,
		WriteTimeout: WriteTimeout,
		IdleTimeout:  IdleTimeout,
	}

	a.FileServer(router)

	go func() {
		if err := a.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Info("failed to start server")
			os.Exit(1)
		}
	}()

	logger.Info("HTTP server is running on ", zap.String("address", a.httpServer.Addr))

	// wait for gracefully shutdown
	<-ctx.Done()

	logger.Info("Shutting down server gracefully")

	shutDownCtx, cancel := context.WithTimeout(context.Background(), shutDownCtxTimeout)
	defer cancel()

	if err := a.httpServer.Shutdown(shutDownCtx); err != nil {
		return fmt.Errorf("shutdown: %w", err)
	}
	<-shutDownCtx.Done()

	return nil
}

func (a *App) FileServer(router *chi.Mux) {
	root := "./web"
	fs := http.FileServer(http.Dir(root))

	router.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		if _, err := os.Stat(root + r.RequestURI); os.IsNotExist(err) {
			http.StripPrefix(r.RequestURI, fs).ServeHTTP(w, r)
		} else {
			fs.ServeHTTP(w, r)
		}
	})
}
