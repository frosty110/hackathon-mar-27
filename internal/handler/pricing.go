package handler

import (
	"context"
	"log/slog"

	pricingv1 "github.com/blaisealbuquerque/pricing-radar/gen/pricing/v1"
	"github.com/blaisealbuquerque/pricing-radar/internal/config"
	"github.com/blaisealbuquerque/pricing-radar/internal/scraper"
	"github.com/blaisealbuquerque/pricing-radar/internal/storage"
)

// PricingHandler implements pricingv1connect.PricingServiceHandler.
type PricingHandler struct {
	cfg     *config.Config
	db      *storage.GhostDB
	targets []scraper.Target
}

func NewPricingHandler(cfg *config.Config, db *storage.GhostDB) *PricingHandler {
	return &PricingHandler{
		cfg:     cfg,
		db:      db,
		targets: scraper.DefaultTargets(),
	}
}

func (h *PricingHandler) RunScan(
	ctx context.Context,
	req *pricingv1.RunScanRequest,
) (*pricingv1.RunScanResponse, error) {
	// 1. Create a new scan run in Ghost DB.
	scanRunID, err := h.db.NewScanRun(ctx)
	if err != nil {
		return nil, err
	}

	// 2. Fetch all 8 target pages concurrently (with fallback).
	rawPages := scraper.FetchAll(ctx, h.targets)

	// 3. Strip boilerplate HTML and build response.
	resp := &pricingv1.RunScanResponse{ScanRunId: scanRunID}
	for _, page := range rawPages {
		stripped, err := scraper.StripBoilerplate(page.HTML)
		if err != nil {
			slog.Warn("StripBoilerplate failed", "competitor", page.Competitor, "error", err)
			stripped = page.HTML // fall back to raw HTML if stripping fails
		}

		// 4. Save snapshot to Ghost DB.
		if saveErr := h.db.SaveSnapshot(ctx, scanRunID, page.Competitor, stripped, page.FromCache); saveErr != nil {
			slog.Warn("SaveSnapshot failed", "competitor", page.Competitor, "error", saveErr)
			// Non-fatal — continue processing other competitors.
		}

		resp.Results = append(resp.Results, &pricingv1.CompetitorResult{
			Competitor:      page.Competitor,
			RawHtmlStripped: stripped,
			FromCache:       page.FromCache,
			ScrapedAt:       page.FetchedAt.Format("2006-01-02T15:04:05Z07:00"),
		})
	}

	// 5. Mark scan run as finished.
	if err := h.db.FinishScanRun(ctx, scanRunID); err != nil {
		slog.Warn("FinishScanRun failed", "scan_run_id", scanRunID, "error", err)
	}

	return resp, nil
}

func (h *PricingHandler) GetComparison(ctx context.Context, req *pricingv1.GetComparisonRequest) (*pricingv1.GetComparisonResponse, error) {
	return &pricingv1.GetComparisonResponse{}, nil
}

func (h *PricingHandler) GetChanges(ctx context.Context, req *pricingv1.GetChangesRequest) (*pricingv1.GetChangesResponse, error) {
	return &pricingv1.GetChangesResponse{}, nil
}

func (h *PricingHandler) GetRecommendation(ctx context.Context, req *pricingv1.GetRecommendationRequest) (*pricingv1.GetRecommendationResponse, error) {
	return &pricingv1.GetRecommendationResponse{}, nil
}

func (h *PricingHandler) GetClusters(ctx context.Context, req *pricingv1.GetClustersRequest) (*pricingv1.GetClustersResponse, error) {
	return &pricingv1.GetClustersResponse{}, nil
}
