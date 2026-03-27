package scraper

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"
)

// RawPage holds the raw (un-stripped) HTML for one competitor.
type RawPage struct {
	Competitor string
	URL        string
	HTML       string
	FromCache  bool
	FetchedAt  time.Time
}

// FetchAll fetches all targets concurrently with a 10-second per-URL timeout.
// On any fetch error, falls back to the cached HTML file and marks FromCache=true.
// Returns all results regardless of individual failures (best-effort fan-out).
func FetchAll(ctx context.Context, targets []Target) []RawPage {
	results := make([]RawPage, len(targets))
	var mu sync.Mutex
	var g errgroup.Group

	for i, t := range targets {
		i, t := i, t
		g.Go(func() error {
			html, err := fetchWithTimeout(ctx, t.URL, 10*time.Second)
			if err != nil {
				slog.Warn("fetch failed, using cache", "competitor", t.Name, "error", err)
				cached, readErr := os.ReadFile(t.FallbackPath)
				if readErr != nil {
					slog.Error("cache read failed", "competitor", t.Name, "path", t.FallbackPath, "error", readErr)
					// Still return nil — don't cancel other goroutines.
					// Leave results[i] as zero value (empty HTML).
					return nil
				}
				mu.Lock()
				results[i] = RawPage{
					Competitor: t.Name,
					URL:        t.URL,
					HTML:       string(cached),
					FromCache:  true,
					FetchedAt:  time.Now(),
				}
				mu.Unlock()
				return nil
			}
			mu.Lock()
			results[i] = RawPage{
				Competitor: t.Name,
				URL:        t.URL,
				HTML:       html,
				FromCache:  false,
				FetchedAt:  time.Now(),
			}
			mu.Unlock()
			return nil
		})
	}

	// errgroup.Wait() cannot return a non-nil error here (all goroutines return nil).
	_ = g.Wait()
	return results
}

func fetchWithTimeout(ctx context.Context, url string, timeout time.Duration) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", fmt.Errorf("build request: %w", err)
	}
	// Set User-Agent to avoid blocks on Go's default "Go-http-client/1.1"
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; PricingRadar/1.0)")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("http.Do: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP %d from %s", resp.StatusCode, url)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read body: %w", err)
	}
	return string(body), nil
}
