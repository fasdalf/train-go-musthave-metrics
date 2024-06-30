package handlers

import (
	"fmt"
	"github.com/fasdalf/train-go-musthave-metrics/internal/common/metricstorage"
	"github.com/fasdalf/train-go-musthave-metrics/internal/constants"
	"net/http"
	"strconv"
	"strings"
)

func UpdateMetricHandler(s metricstorage.Storage) http.HandlerFunc {
	const indexType = 0
	const indexName = 1
	const indexValue = 2
	return func(w http.ResponseWriter, r *http.Request) {
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

		if len(path) < 3 {
			fmt.Println("Invalid path")
			http.Error(w, "Invalid path", http.StatusNotFound)
			return
		}

		if path[indexName] == "" {
			fmt.Println("Metric not found")
			http.Error(w, "Metric not found", http.StatusNotFound)
			return
		}

		switch path[indexType] {
		case constants.GaugeStr:
			floatValue, err := strconv.ParseFloat(path[indexValue], 64)
			if err != nil {
				fmt.Println("Invalid metric value", path[indexValue])
				http.Error(w, "Invalid metric value, float64 required", http.StatusBadRequest)
				return
			}
			s.UpdateGauge(path[indexName], floatValue)
			fmt.Println("New value", s.GetGauge(path[indexName]))
		case constants.CounterStr:
			intValue, err := strconv.Atoi(path[indexValue])
			if err != nil {
				fmt.Println("Invalid metric value", path[indexValue])
				http.Error(w, "Invalid metric value, integer required", http.StatusBadRequest)
				return
			}
			s.UpdateCounter(path[indexName], intValue)
			fmt.Println("New value", s.GetCounter(path[indexName]))
		default:
			fmt.Println("Invalid type", path[indexType])
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
}
