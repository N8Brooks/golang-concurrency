package webcrawlermultithreaded

import (
	"net/url"
	"sync"
)

type HtmlParser interface {
	GetUrls(url string) []string
}

func Crawl(startURL string, htmlParser HtmlParser) []string {
	startHost := hostName(startURL)

	seen := map[string]struct{}{
		startURL: {},
	}
	urls := []string{startURL}

	var mu sync.Mutex
	var wg sync.WaitGroup

	var visit func(string)
	visit = func(current string) {
		defer wg.Done()

		for _, next := range htmlParser.GetUrls(current) {
			if hostName(next) != startHost {
				continue
			}

			mu.Lock()
			if _, ok := seen[next]; ok {
				mu.Unlock()
				continue
			}
			seen[next] = struct{}{}
			urls = append(urls, next)
			mu.Unlock()
			wg.Add(1)

			go visit(next)
		}
	}

	wg.Add(1)
	go visit(startURL)
	wg.Wait()

	return urls
}

func hostName(rawURL string) string {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return ""
	}
	return parsed.Host
}
