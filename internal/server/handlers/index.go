package handlers

import (
	"fmt"
	"log/slog"
	"net/http"
)

// NewIndexHandler list all stored metrics
func NewIndexHandler(ms Storage) http.HandlerFunc {
	const (
		pageTemplate = `<html><body>
<table>
<tr><td colspan=2>Gauges</td></tr>
%s<tr><td colspan=2>counters</td></tr>
%s</table>
</body></html>
`
		rowTemplate = `<tr><td>%s</td><td>%s</td></tr>
`
	)

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			slog.Error("GET requests only", `requested`, r.Method)
			http.Error(w, "GET requests only", http.StatusMethodNotAllowed)
			return
		}
		slog.Info("showing index")

		l, err := ms.ListGauges()
		if err != nil {
			slog.Error("can't list gauges", "error", err)
			http.Error(w, `unexpected error`, http.StatusInternalServerError)
			return
		}
		gauges := ""
		for _, key := range l {
			gauges += fmt.Sprintf(rowTemplate, key, fmt.Sprint(ms.GetGauge(key)))
		}

		l, err = ms.ListCounters()
		if err != nil {
			slog.Error("can't list counters", "error", err)
			http.Error(w, `unexpected error`, http.StatusInternalServerError)
			return
		}
		counters := ""
		for _, key := range l {
			counters += fmt.Sprintf(rowTemplate, key, fmt.Sprint(ms.GetCounter(key)))
		}

		w.Header().Set("Content-Type", "text/html")
		_, _ = w.Write([]byte(fmt.Sprintf(pageTemplate, gauges, counters)))

		slog.Info("Processed OK")
	}
}
