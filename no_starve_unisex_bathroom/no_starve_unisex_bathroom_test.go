//go:build challenge

package nostarveunisexbathroom_test

import (
	"testing"

	nostarveunisexbathroom "github.com/N8Brooks/golang-concurrency/no_starve_unisex_bathroom"
	"github.com/N8Brooks/golang-concurrency/no_starve_unisex_bathroom/testsuite"
)

func TestNoStarveUnisexBathroom(t *testing.T) {
	testsuite.Run(t, func() testsuite.NoStarveUnisexBathroom {
		return nostarveunisexbathroom.NewNoStarveUnisexBathroom()
	})
}

func BenchmarkNoStarveUnisexBathroom(b *testing.B) {
	testsuite.Benchmark(b, func() testsuite.NoStarveUnisexBathroom {
		return nostarveunisexbathroom.NewNoStarveUnisexBathroom()
	})
}
