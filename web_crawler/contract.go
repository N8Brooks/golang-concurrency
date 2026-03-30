// Package webcrawler contains the web crawler challenge and shared contract.
//
// Starting from a URL, the crawler should visit only pages on the same host,
// deduplicate URLs it has already seen, and return every reachable same-host
// URL. Implementations may crawl concurrently.
package webcrawler

// HTMLParser returns the outgoing links for a URL.
type HTMLParser interface {
	GetURLs(rawURL string) []string
}

// Crawler is the shared contract used by the challenge implementation,
// reusable tests, and hint solutions.
type Crawler interface {
	Crawl(startURL string, htmlParser HTMLParser) []string
}
