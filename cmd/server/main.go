package main

import (
	"github.com/fasdalf/train-go-musthave-metrics/internal/common/metricstorage"
	"github.com/fasdalf/train-go-musthave-metrics/internal/server/handlers"
	"net/http"
)

func main() {
	memStorage := metricstorage.NewMemStorage()
	mux := http.NewServeMux()
	//mux.HandleFunc(`/`, mainPage)
	mux.Handle("/update/", http.StripPrefix("/update/", handlers.UpdateMetricHandler(memStorage)))
	////mux.HandleFunc(`/json`, JSONHandler)
	//fs := http.FileServer(http.Dir("."))
	//mux.Handle("/golang/", http.StripPrefix("/golang/", fs))
	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
}
