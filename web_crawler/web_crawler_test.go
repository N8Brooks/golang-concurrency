//go:build challenge

package webcrawler_test

import (
	"testing"

	webcrawler "github.com/N8Brooks/golang-concurrency/web_crawler"
	"github.com/N8Brooks/golang-concurrency/web_crawler/testsuite"
)

func TestWebCrawler(t *testing.T) {
	testsuite.Run(t, func() webcrawler.Crawler {
		return webcrawler.NewChallenge()
	})
}

func BenchmarkWebCrawler(b *testing.B) {
	testsuite.Benchmark(b, func() webcrawler.Crawler {
		return webcrawler.NewChallenge()
	})
}
