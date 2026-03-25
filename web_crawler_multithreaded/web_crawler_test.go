package webcrawlermultithreaded

import (
	"slices"
	"sync"
	"testing"
	"time"
)

type mockHtmlParser struct {
	graph map[string][]string
	delay time.Duration

	mu        sync.Mutex
	active    int
	maxActive int
}

func (p *mockHtmlParser) GetUrls(url string) []string {
	p.mu.Lock()
	p.active++
	if p.active > p.maxActive {
		p.maxActive = p.active
	}
	p.mu.Unlock()

	if p.delay > 0 {
		time.Sleep(p.delay)
	}

	p.mu.Lock()
	p.active--
	p.mu.Unlock()

	return append([]string(nil), p.graph[url]...)
}

func (p *mockHtmlParser) MaxActive() int {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.maxActive
}

func TestCrawlExample1(t *testing.T) {
	parser := &mockHtmlParser{
		graph: map[string][]string{
			"http://news.yahoo.com": {
				"http://news.yahoo.com/news",
				"http://news.yahoo.com/us",
				"http://news.google.com",
			},
			"http://news.yahoo.com/news": {
				"http://news.yahoo.com",
				"http://news.yahoo.com/news/topics/",
			},
			"http://news.yahoo.com/news/topics/": {
				"http://news.yahoo.com/news",
			},
			"http://news.yahoo.com/us": {
				"http://news.yahoo.com",
			},
			"http://news.google.com": {
				"http://news.google.com/world",
			},
		},
	}

	got := Crawl("http://news.yahoo.com", parser)
	slices.Sort(got)

	want := []string{
		"http://news.yahoo.com",
		"http://news.yahoo.com/news",
		"http://news.yahoo.com/news/topics/",
		"http://news.yahoo.com/us",
	}
	slices.Sort(want)

	if !slices.Equal(got, want) {
		t.Fatalf("crawl() = %v, want %v", got, want)
	}
}

func TestCrawlExample2(t *testing.T) {
	parser := &mockHtmlParser{
		graph: map[string][]string{
			"http://news.google.com": {
				"http://news.yahoo.com",
				"http://news.yahoo.com/news",
			},
			"http://news.yahoo.com": {
				"http://news.yahoo.com/us",
			},
		},
	}

	got := Crawl("http://news.google.com", parser)
	if want := []string{"http://news.google.com"}; !slices.Equal(got, want) {
		t.Fatalf("crawl() = %v, want %v", got, want)
	}
}

func TestCrawlVisitsEachSameHostURLOnceAndRunsConcurrently(t *testing.T) {
	parser := &mockHtmlParser{
		delay: 20 * time.Millisecond,
		graph: map[string][]string{
			"http://site.com": {
				"http://site.com/a",
				"http://site.com/b",
				"http://site.com/c",
				"http://other.com/ignored",
			},
			"http://site.com/a": {
				"http://site.com/shared",
			},
			"http://site.com/b": {
				"http://site.com/shared",
			},
			"http://site.com/c": {
				"http://site.com/shared",
			},
			"http://site.com/shared": {
				"http://site.com",
			},
		},
	}

	got := Crawl("http://site.com", parser)
	slices.Sort(got)

	want := []string{
		"http://site.com",
		"http://site.com/a",
		"http://site.com/b",
		"http://site.com/c",
		"http://site.com/shared",
	}
	slices.Sort(want)

	if !slices.Equal(got, want) {
		t.Fatalf("crawl() = %v, want %v", got, want)
	}
	if parser.MaxActive() < 2 {
		t.Fatalf("expected concurrent crawling, max active GetUrls calls = %d", parser.MaxActive())
	}
}
