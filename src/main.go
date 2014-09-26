package main

import (
	"flag"
	"fmt"
	. "gotrie"
	"math/rand"
	"runtime"
	"time"
)

type Stat struct {
	IndexTime   time.Duration
	QueryTime   time.Duration
	MemoryUsage uint64
	start       time.Time
}

func (this *Stat) Start() {
	this.start = time.Now()
}
func (this *Stat) MeasureIndexTime() {
	var now = time.Now()
	this.IndexTime = now.Sub(this.start)
}
func (this *Stat) MeasureQueryTime() {
	var now = time.Now()
	this.QueryTime = now.Sub(this.start)
}
func (this *Stat) MeasureMemory(idx *Index) {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)
	this.MemoryUsage = mem.Alloc
}
func (this *Stat) String() string {
	return fmt.Sprintf("%.4f,%.4f,%d",
		float64(this.IndexTime)/float64(time.Millisecond),
		float64(this.QueryTime)/float64(time.Millisecond),
		this.MemoryUsage)
}

func main() {
	var N int = 1000
	var K int = 1000
	var R int = 5
	var stat Stat
	var graphviz bool
	var optlevel int = 0

	flag.IntVar(&N, "points", 1000, "Number of data points")
	flag.IntVar(&K, "queries", 1000, "Number of queries")
	flag.IntVar(&R, "radius", 5, "Hamming distance radius")
	flag.IntVar(&optlevel, "O", 0, "Optimization level")
	flag.BoolVar(&graphviz, "graphviz", false, "Graphviz")
	flag.Parse()

	var idx = NewIndex(N)
	var points = make([]uint64, N)
	stat.Start()
	for i := range points {
		points[i] = uint64(rand.Int63())
	}
	stat.MeasureIndexTime()

	stat.Start()
	for i := uint64(0); i < uint64(N); i++ {
		idx.Add(i, points[i])
	}
	stat.MeasureQueryTime()
	stat.MeasureMemory(idx)

	if graphviz {
		idx.Graphviz()
		return
	}

	searcher := NewSearcher(idx)
	searcher.O = optlevel

	stat.Start()
	for i := 0; i < K; i++ {
		q := uint64(rand.Int63())
		searcher.Search(q, uint8(R))
	}
	stat.MeasureQueryTime()

	fmt.Println(stat.String())
}
