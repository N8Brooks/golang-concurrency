// Package testsuite contains reusable behavioral tests for web crawler implementations.
package testsuite

import (
	"slices"
	"sync"
	"testing"
	"time"

	webcrawler "github.com/N8Brooks/golang-concurrency/web_crawler"
)

type mockHTMLParser struct {
	graph map[string][]string
	delay time.Duration

	mu        sync.Mutex
	active    int
	maxActive int
}

func (p *mockHTMLParser) GetURLs(rawURL string) []string {
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

	return append([]string(nil), p.graph[rawURL]...)
}

func (p *mockHTMLParser) MaxActive() int {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.maxActive
}

func Run(t *testing.T, newImpl func() webcrawler.Crawler) {
	t.Helper()

	t.Run("Example1", func(t *testing.T) {
		parser := &mockHTMLParser{
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

		got := newImpl().Crawl("http://news.yahoo.com", parser)
		want := []string{
			"http://news.yahoo.com",
			"http://news.yahoo.com/news",
			"http://news.yahoo.com/news/topics/",
			"http://news.yahoo.com/us",
		}

		slices.Sort(got)
		slices.Sort(want)
		if !slices.Equal(got, want) {
			t.Fatalf("crawl() = %v, want %v", got, want)
		}
	})

	t.Run("Example2", func(t *testing.T) {
		parser := &mockHTMLParser{
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

		got := newImpl().Crawl("http://news.google.com", parser)
		if want := []string{"http://news.google.com"}; !slices.Equal(got, want) {
			t.Fatalf("crawl() = %v, want %v", got, want)
		}
	})

	t.Run("Reusable", func(t *testing.T) {
		crawler := newImpl()

		parser1 := &mockHTMLParser{
			graph: map[string][]string{
				"http://a.com": {
					"http://a.com/1",
					"http://b.com/ignored",
				},
				"http://a.com/1": {
					"http://a.com/2",
				},
				"http://a.com/2": nil,
			},
		}
		got1 := crawler.Crawl("http://a.com", parser1)
		want1 := []string{
			"http://a.com",
			"http://a.com/1",
			"http://a.com/2",
		}
		slices.Sort(got1)
		slices.Sort(want1)
		if !slices.Equal(got1, want1) {
			t.Fatalf("first crawl() = %v, want %v", got1, want1)
		}

		parser2 := &mockHTMLParser{
			graph: map[string][]string{
				"http://site.com": {
					"http://site.com/a",
				},
				"http://site.com/a": nil,
			},
		}
		got2 := crawler.Crawl("http://site.com", parser2)
		want2 := []string{
			"http://site.com",
			"http://site.com/a",
		}
		slices.Sort(got2)
		slices.Sort(want2)
		if !slices.Equal(got2, want2) {
			t.Fatalf("second crawl() = %v, want %v", got2, want2)
		}
	})

	t.Run("VisitsEachSameHostURLOnceAndRunsConcurrently", func(t *testing.T) {
		parser := &mockHTMLParser{
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

		got := newImpl().Crawl("http://site.com", parser)
		want := []string{
			"http://site.com",
			"http://site.com/a",
			"http://site.com/b",
			"http://site.com/c",
			"http://site.com/shared",
		}

		slices.Sort(got)
		slices.Sort(want)
		if !slices.Equal(got, want) {
			t.Fatalf("crawl() = %v, want %v", got, want)
		}
		if parser.MaxActive() < 2 {
			t.Fatalf("expected concurrent crawling, max active GetURLs calls = %d", parser.MaxActive())
		}
	})
}

func Benchmark(b *testing.B, newImpl func() webcrawler.Crawler) {
	b.Helper()
	b.ReportAllocs()

	parser := &mockHTMLParser{
		graph: map[string][]string{
			"http://site.com": {
				"http://site.com/a",
				"http://site.com/b",
				"http://site.com/c",
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
			"http://site.com/shared": nil,
		},
	}

	b.Run("SmallGraph", func(b *testing.B) {
		for b.Loop() {
			got := newImpl().Crawl("http://site.com", parser)
			if len(got) != 5 {
				b.Fatalf("len(crawl()) = %d, want 5", len(got))
			}
		}
	})
}
