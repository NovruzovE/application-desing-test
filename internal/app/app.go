package app

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/NovruzovE/application-design-test/internal/config"
	"github.com/NovruzovE/application-design-test/internal/core/usecase/order"
	orderController "github.com/NovruzovE/application-design-test/internal/handler/order"
	"github.com/NovruzovE/application-design-test/internal/repo"
	"github.com/NovruzovE/application-design-test/internal/transaction_manager"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

type App struct {
	orderController *orderController.OrderController
	router          *chi.Mux
	httpServer      *http.Server
	logger          *slog.Logger
	config          *config.Config
}

func NewApp() *App {
	a := &App{}

	var err error
	a.config, err = config.NewConfig()
	if err != nil {
		panic(err)
	}

	a.initLogger()
	a.initOrderController()
	a.initHTTPRouter()
	a.initHTTPServer()

	return a
}

func (a *App) MustRun() {
	a.logger.Info("Starting HTTP server", "address", a.config.HTTP.Address)
	a.logger.Debug("Debug logs enabled")

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	defer stop()

	go func() {
		err := a.httpServer.ListenAndServe()
		if errors.Is(err, http.ErrServerClosed) {
			a.logger.Info("HTTP server stopped")
		} else if err != nil {
			a.logger.Error("HTTP server stopped with error", err)
			os.Exit(1)
		}
	}()

	<-ctx.Done()

	a.ShutdownHTTPServer()

}

func (a *App) initHTTPRouter() {
	a.router = chi.NewRouter()

	a.router.Use(middleware.RequestID)
	a.router.Use(middleware.Logger)
	a.router.Use(middleware.Recoverer)

	a.router.Post("/orders", a.orderController.CreateOrder)
}

func (a *App) initHTTPServer() {
	a.httpServer = &http.Server{
		Addr:         a.config.Address,
		Handler:      a.router,
		ReadTimeout:  a.config.HTTP.Timeout * time.Second,
		WriteTimeout: a.config.HTTP.Timeout * time.Second,
		IdleTimeout:  a.config.HTTP.IdleTimeout * time.Second,
	}
}

func (a *App) initLogger() {

	switch a.config.Env {
	case envLocal:
		a.logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		a.logger = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		a.logger = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}
}

func (a *App) initOrderController() {
	roomAvailRepo := repo.NewRoomAvailabilityInMemRepo(a.logger)
	if a.config.Env == envLocal {
		roomAvailRepo.PrepareRepo()
	}
	orderRepo := repo.NewOrderInMemRepo(a.logger)
	transactionManager := transaction_manager.NewMemTransactionManager(a.logger)

	orderUseCase := order.NewOrderUseCase(roomAvailRepo, orderRepo, transactionManager, a.logger)
	a.orderController = orderController.New(orderUseCase, a.logger)
}

func (a *App) ShutdownHTTPServer() {
	a.logger.Info("Shutting down HTTP server")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := a.httpServer.Shutdown(ctx); err != nil {
		a.logger.Error("failed to stop HTTP server", err)
		return
	}
}
