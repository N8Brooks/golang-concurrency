package solutions_test

import (
	"testing"

	webcrawler "github.com/N8Brooks/golang-concurrency/web_crawler"
	"github.com/N8Brooks/golang-concurrency/web_crawler/solutions"
	"github.com/N8Brooks/golang-concurrency/web_crawler/testsuite"
)

func TestWaitGroup(t *testing.T) {
	testsuite.Run(t, func() webcrawler.Crawler {
		return solutions.NewWaitGroup()
	})
}

func BenchmarkWaitGroup(b *testing.B) {
	testsuite.Benchmark(b, func() webcrawler.Crawler {
		return solutions.NewWaitGroup()
	})
}
