package game_test

import (
	"testing"

	"github.com/RaniSputnik/ok/game"
)

/*
goos: darwin
goarch: amd64
pkg: github.com/RaniSputnik/ko2/game
BenchmarkPlay5-8                   30000             58928 ns/op           21296 B/op        408 allocs/op
BenchmarkPlayTiny-8                10000            199783 ns/op           78000 B/op       1378 allocs/op
BenchmarkPlaySmall-8                3000            422476 ns/op          178384 B/op       2923 allocs/op
BenchmarkPlayNormal-8               2000            971839 ns/op          458864 B/op       6320 allocs/op
PASS
ok      github.com/RaniSputnik/ko2/game 7.791s
*/

func BenchmarkPlay5(b *testing.B) {
	benchmarkPlayEverywhere(b, 5)
}

func BenchmarkPlayTiny(b *testing.B) {
	benchmarkPlayEverywhere(b, game.BoardSizeTiny)
}

func BenchmarkPlaySmall(b *testing.B) {
	benchmarkPlayEverywhere(b, game.BoardSizeSmall)
}

func BenchmarkPlayNormal(b *testing.B) {
	benchmarkPlayEverywhere(b, game.BoardSizeNormal)
}

func benchmarkPlayEverywhere(b *testing.B, size int) {
	b.ReportAllocs()

	for n := 0; n < b.N; n++ {
		m := game.New(size)

		for y := 0; y < size; y++ {
			for x := 0; x < size; x++ {
				m.Play(game.Stone{Colour: m.Next(), Position: game.Position{X: x, Y: y}})
			}
		}
	}
}
