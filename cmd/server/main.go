package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"connectrpc.com/connect"
	pricingv1connect "github.com/blaisealbuquerque/pricing-radar/gen/pricing/v1/pricingv1connect"
	"github.com/blaisealbuquerque/pricing-radar/internal/config"
	"github.com/blaisealbuquerque/pricing-radar/internal/handler"
	"github.com/blaisealbuquerque/pricing-radar/internal/storage"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	cfg, err := config.Load()
	if err != nil {
		logger.Error("config load failed", "error", err)
		os.Exit(1)
	}

	db, err := storage.NewGhostDB(context.Background(), cfg.DatabaseURL)
	if err != nil {
		logger.Error("failed to connect to Ghost DB", "error", err)
		os.Exit(1)
	}
	if err := db.AutoMigrate(context.Background()); err != nil {
		logger.Error("AutoMigrate failed", "error", err)
		os.Exit(1)
	}
	logger.Info("ghost db connected and migrated")

	h := handler.NewPricingHandler(cfg, db)

	mux := http.NewServeMux()
	path, httpHandler := pricingv1connect.NewPricingServiceHandler(
		h,
		connect.WithCompressMinBytes(1024),
	)
	mux.Handle(path, httpHandler)

	// Health check endpoint
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	p := new(http.Protocols)
	p.SetHTTP1(true)
	p.SetUnencryptedHTTP2(true)

	srv := &http.Server{
		Addr:      ":" + cfg.Port,
		Handler:   mux,
		Protocols: p,
	}

	logger.Info("server listening", "port", cfg.Port)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Error("server error", "error", err)
		os.Exit(1)
	}
}
