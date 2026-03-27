package handler

import (
	"context"

	pricingv1 "github.com/blaisealbuquerque/pricing-radar/gen/pricing/v1"
	"github.com/blaisealbuquerque/pricing-radar/internal/config"
	"github.com/blaisealbuquerque/pricing-radar/internal/storage"
)

// PricingHandler implements pricingv1connect.PricingServiceHandler.
// Phase 1: stub — returns empty responses. Scraper is wired in plan 01-04.
type PricingHandler struct {
	cfg *config.Config
	db  *storage.GhostDB
}

func NewPricingHandler(cfg *config.Config, db *storage.GhostDB) *PricingHandler {
	return &PricingHandler{cfg: cfg, db: db}
}

func (h *PricingHandler) RunScan(
	ctx context.Context,
	req *pricingv1.RunScanRequest,
) (*pricingv1.RunScanResponse, error) {
	return &pricingv1.RunScanResponse{}, nil
}

func (h *PricingHandler) GetComparison(
	ctx context.Context,
	req *pricingv1.GetComparisonRequest,
) (*pricingv1.GetComparisonResponse, error) {
	return &pricingv1.GetComparisonResponse{}, nil
}

func (h *PricingHandler) GetChanges(
	ctx context.Context,
	req *pricingv1.GetChangesRequest,
) (*pricingv1.GetChangesResponse, error) {
	return &pricingv1.GetChangesResponse{}, nil
}

func (h *PricingHandler) GetRecommendation(
	ctx context.Context,
	req *pricingv1.GetRecommendationRequest,
) (*pricingv1.GetRecommendationResponse, error) {
	return &pricingv1.GetRecommendationResponse{}, nil
}

func (h *PricingHandler) GetClusters(
	ctx context.Context,
	req *pricingv1.GetClustersRequest,
) (*pricingv1.GetClustersResponse, error) {
	return &pricingv1.GetClustersResponse{}, nil
}
