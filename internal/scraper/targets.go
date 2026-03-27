package scraper

// Target defines a competitor pricing page to scrape.
type Target struct {
	Name         string
	URL          string
	FallbackPath string
}

// DefaultTargets returns the 8 pre-selected AI code review pricing pages.
// All URLs are verified as httpx-fetchable (no JS rendering needed).
func DefaultTargets() []Target {
	return []Target{
		{Name: "CodeRabbit", URL: "https://www.coderabbit.ai/pricing", FallbackPath: "demo-data/cached/coderabbit.html"},
		{Name: "Codacy", URL: "https://www.codacy.com/pricing", FallbackPath: "demo-data/cached/codacy.html"},
		{Name: "DeepSource", URL: "https://deepsource.com/pricing", FallbackPath: "demo-data/cached/deepsource.html"},
		{Name: "Sourcery AI", URL: "https://www.sourcery.ai/pricing", FallbackPath: "demo-data/cached/sourcery.html"},
		{Name: "Qodo", URL: "https://www.qodo.ai/pricing/", FallbackPath: "demo-data/cached/qodo.html"},
		{Name: "Snyk", URL: "https://snyk.io/plans/", FallbackPath: "demo-data/cached/snyk.html"},
		{Name: "Greptile", URL: "https://www.greptile.com/pricing", FallbackPath: "demo-data/cached/greptile.html"},
		{Name: "CodeAnt AI", URL: "https://www.codeant.ai/pricing", FallbackPath: "demo-data/cached/codeant.html"},
	}
}
