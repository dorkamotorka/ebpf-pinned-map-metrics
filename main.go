package main

import (
	"time"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"

	"github.com/cilium/ebpf"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	ebpfFsPath = "/sys/fs/bpf" 
)

var (
	mapElementCount = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "ebpf_map_element_count",
			Help: "Number of elements in eBPF map",
		},
		[]string{"map_name"},
	)
)

func loadPinnedMaps(dir string) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatalf("Failed to read directory %s: %v", dir, err)
	}

	for _, file := range files {
		if file.IsDir() {
			loadPinnedMaps(filepath.Join(dir, file.Name()))
			continue
		}

		mapPath := filepath.Join(dir, file.Name())
		m, err := ebpf.LoadPinnedMap(mapPath, nil)
		if err != nil {
			log.Printf("Failed to load map %s: %v", mapPath, err)
			continue
		}

		countNonZeroElements(mapPath, m)
	}
}

func countNonZeroElements(mapPath string, m *ebpf.Map) {
	iterator := m.Iterate()
	var key, value []byte
	nonZeroCount := 0

	for iterator.Next(&key, &value) {
		if !isZeroValue(value) {
			nonZeroCount++
		}
	}

	if err := iterator.Err(); err != nil {
		log.Printf("Failed to iterate over map %s: %v", mapPath, err)
	}

	mapElementCount.WithLabelValues(mapPath).Set(float64(nonZeroCount))
}

func isZeroValue(value []byte) bool {
	for _, b := range value {
		if b != 0 {
			return false
		}
	}
	return true
}

func main() {
	reg := prometheus.NewRegistry()
	reg.MustRegister(mapElementCount)
    	handler := promhttp.HandlerFor(reg, promhttp.HandlerOpts{})
	go func() {
		for {
			loadPinnedMaps(ebpfFsPath)
			// Adjust the interval as needed
			time.Sleep(1 * time.Second)
		}
	}()

 	// Start HTTP server for Prometheus metrics
	http.Handle("/metrics", handler)
	log.Fatal(http.ListenAndServe(":2112", nil))
}
