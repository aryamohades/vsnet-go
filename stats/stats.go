package stats

import (
	"encoding/json"
	"log"
	"net/http"
	"runtime"
	"sync/atomic"
)

var connections int64

type Stats struct {
	Connections   int64   `json:"connections"`
	NumGoroutines int     `json:"num_goroutines"`
	TotalAlloc    uint64  `json:"total_alloc"`
	HeapAlloc     uint64  `json:"heap_alloc"`
	HeapSys       uint64  `json:"heap_sys"`
	NumGC         uint32  `json:"num_gc"`
	NextGC        uint64  `json:"next_gc"`
	GCCPUFraction float64 `json:"gccpu_fraction"`
	PauseTotalNS  uint64  `json:"pause_total_ns"`
}

func stats(w http.ResponseWriter, r *http.Request) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	s := Stats{
		Connections:   atomic.LoadInt64(&connections),
		NumGoroutines: runtime.NumGoroutine(),
		HeapAlloc:     m.HeapAlloc,
		TotalAlloc:    m.TotalAlloc,
		HeapSys:       m.HeapSys,
		NumGC:         m.NumGC,
		NextGC:        m.NextGC,
		GCCPUFraction: m.GCCPUFraction,
		PauseTotalNS:  m.PauseTotalNs,
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(s)
}

func Start(addr string) {
	mux := http.NewServeMux()
	mux.HandleFunc("/stats", stats)

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("stats error: %v", err)
	}
}

func AddConnection() {
	atomic.AddInt64(&connections, 1)
}

func RemoveConnection() {
	atomic.AddInt64(&connections, -1)
}
