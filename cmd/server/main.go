package main

import (
	"fmt"
	"github.com/fasdalf/train-go-musthave-metrics/internal/constants"
	"github.com/fasdalf/train-go-musthave-metrics/internal/server"
	"net/http"
	"strconv"
	"strings"
)

func updateMetric(w http.ResponseWriter, r *http.Request, s server.Storage) {
	fmt.Println("Processing", r.Method, r.URL.Path)
	if r.Method != http.MethodPost {
		fmt.Println("POST requests only")
		http.Error(w, "POST requests only", http.StatusBadRequest)
		return
	}

	path := strings.Split(r.URL.Path, "/")
	fmt.Println("path is", len(path), path)
	if len(path) > 0 && path[0] == "" {
		path = path[1:]
	}

	if len(path) < 4 {
		fmt.Println("Invalid path")
		http.Error(w, "Invalid path", http.StatusNotFound)
		return
	}

	if path[2] == "" {
		fmt.Println("Metric not found")
		http.Error(w, "Metric not found", http.StatusNotFound)
		return
	}

	switch path[1] {
	case constants.GaugeStr:
		floatValue, err := strconv.ParseFloat(path[3], 64)
		if err != nil {
			fmt.Println("Invalid metric value", path[3])
			http.Error(w, "Invalid metric value, float64 required", http.StatusBadRequest)
			return
		}
		s.UpdateGauge(path[2], floatValue)
		fmt.Println("New value", s.GetGauge(path[2]))
	case constants.CounterStr:
		intValue, err := strconv.Atoi(path[3])
		if err != nil {
			fmt.Println("Invalid metric value", path[3])
			http.Error(w, "Invalid metric value, integer required", http.StatusBadRequest)
			return
		}
		s.UpdateCounter(path[2], intValue)
		fmt.Println("New value", s.GetCounter(path[2]))
	default:
		fmt.Println("Invalid type", path[1])
		http.Error(w, fmt.Sprintf(
			"Invalid type, only %s and %s supported",
			constants.GaugeStr,
			constants.CounterStr,
		), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Println("Processed OK")
}

func main() {
	storage := server.NewMemStorage()
	mux := http.NewServeMux()
	//mux.HandleFunc(`/`, mainPage)
	mux.HandleFunc("/update/", func(w http.ResponseWriter, r *http.Request) {
		updateMetric(w, r, storage)
	})
	////mux.HandleFunc(`/json`, JSONHandler)
	//fs := http.FileServer(http.Dir("."))
	//mux.Handle("/golang/", http.StripPrefix("/golang/", fs))
	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
}
