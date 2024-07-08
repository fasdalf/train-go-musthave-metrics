package handlers

import (
	"fmt"
	"github.com/fasdalf/train-go-musthave-metrics/internal/common/constants"
	"github.com/fasdalf/train-go-musthave-metrics/internal/common/metricstorage"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
)

func NewUpdateMetricHandler(s metricstorage.Storage) http.HandlerFunc {
	const indexType = 0
	const indexName = 1
	const indexValue = 2
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Processing update", "method", r.Method, "path", r.URL.Path)
		if r.Method != http.MethodPost {
			// job for gin
			http.Error(w, "POST requests only", http.StatusBadRequest)
			return
		}

		// job for gin
		path := strings.Split(r.URL.Path, "/")
		fmt.Println("path is", len(path), path)
		if len(path) > 0 && path[0] == "" {
			path = path[1:]
		}

		// job for gin
		if len(path) < 3 {
			slog.Error("GET requests only", `requested`, r.Method)
			fmt.Println("Invalid path")
			http.Error(w, "Invalid path", http.StatusNotFound)
			return
		}

		// job for gin
		if path[indexName] == "" {
			fmt.Println("Metric not found")
			http.Error(w, "Metric not found", http.StatusNotFound)
			return
		}

		switch path[indexType] {
		case constants.GaugeStr:
			floatValue, err := strconv.ParseFloat(path[indexValue], 64)
			if err != nil {
				slog.Error("Invalid metric value, float64 required", `value`, path[indexValue])
				fmt.Println("Invalid metric value")
				http.Error(w, "", http.StatusBadRequest)
				return
			}
			s.UpdateGauge(path[indexName], floatValue)
			slog.Info("value set", "new", s.GetGauge(path[indexName]))
		case constants.CounterStr:
			intValue, err := strconv.Atoi(path[indexValue])
			if err != nil {
				slog.Error("Invalid metric value, integer required", `value`, path[indexValue])
				http.Error(w, "Invalid metric value, integer required", http.StatusBadRequest)
				return
			}
			s.UpdateCounter(path[indexName], intValue)
			slog.Info("value set", "new", s.GetGauge(path[indexName]))
		default:
			slog.Error("Type is invalid", "type", path[indexValue])
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
