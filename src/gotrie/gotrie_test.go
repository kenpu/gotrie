package gotrie

import (
	"math/rand"
	"testing"
)

func TestIndex(t *testing.T) {
	N := 1000000
	var idx = NewIndex(N)
	rand.Seed(0)
	for i := 0; i < N; i++ {
		p := uint64(rand.Int63())
		idx.Add(uint64(i), p)
	}

	// the split must be monotonically increasing
	for i, node := range idx.nodes {
		if i > 0 {
			for j := 0; j < 2; j++ {
				child_i := node.leftRight[j]
				if child_i > 0 {
					child := idx.nodes[child_i]
					if node.splitPos >= child.splitPos {
						t.Fatalf(
							"Non-monotonic split position detected. %v -> %v\n",
							node, child)
					}
				}
			}
		}
	}
}

func TestSearch(t *testing.T) {
	N := 1000000
	K := 1000
	r := uint8(5)
	var idx = NewIndex(N)
	rand.Seed(0)
	for i := 0; i < N; i++ {
		p := uint64(rand.Int63())
		idx.Add(uint64(i), p)
	}
	searcher := NewSearcher(idx)

	var q uint64
	var p int
	for i := 0; i < K; i++ {
		if i%2 == 0 {
			p = rand.Intn(len(idx.nodes)-1) + 1
			q = idx.nodes[p].value
		} else {
			q = uint64(rand.Int63())
		}
		count := searcher.Search(q, r)
		if i%2 == 0 {
			if count != 1 {
				t.Fatalf("(i=%d) count = %d != 1", i, count)
			}
		} else {
			if count != 0 {
				t.Fatalf("(i=%d) count = %d != 0", i, count)
			}
		}
	}
}
