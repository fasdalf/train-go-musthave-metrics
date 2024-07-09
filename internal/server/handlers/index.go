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

		gauges := ""
		for _, key := range ms.ListGauges() {
			gauges += fmt.Sprintf(rowTemplate, key, fmt.Sprint(ms.GetGauge(key)))
		}
		counters := ""
		for _, key := range ms.ListCounters() {
			counters += fmt.Sprintf(rowTemplate, key, fmt.Sprint(ms.GetCounter(key)))
		}

		w.Header().Set("Content-Type", "text/html")
		_, _ = w.Write([]byte(fmt.Sprintf(pageTemplate, gauges, counters)))

		slog.Info("Processed OK")
	}
}
