package metricstorage

import (
	"context"
	"strconv"
	"testing"

	"github.com/fasdalf/train-go-musthave-metrics/internal/common/apimodels"
	"github.com/fasdalf/train-go-musthave-metrics/internal/common/constants"
)

func BenchmarkModelStorageSingleUpdates(b *testing.B) {
	count := 20

	ds := NewDirtyStorage(NewMemStorage())
	ms := NewSavableModelStorage(ds)
	m := &apimodels.Metrics{}
	for i := 0; i < count; i++ {
		if i%2 == 0 {
			v := int64(i)
			m.MType = constants.CounterStr
			m.ID = "counter" + strconv.Itoa(i)
			m.Delta = &v
		} else {
			v := float64(i)
			m.MType = constants.GaugeStr
			m.ID = "gauge" + strconv.Itoa(i)
			m.Value = &v
		}
		_ = ms.SaveCommonModel(m)
	}
	b.ResetTimer()
	b.Run("single updates", func(b *testing.B) {
		ds.Clear()
		var n int
		var f float64
		for i := 0; i < b.N; i++ {
			if i%2 == 0 {
				v := float64(i)
				m.MType = constants.CounterStr
				m.ID = "counter" + strconv.Itoa(i%count)
				m.Value = &v
				_ = ms.SaveCommonModel(m)
				n, _ = ms.GetCounter("counter" + strconv.Itoa(i%count))
			} else {
				v := float64(i)
				m.MType = constants.GaugeStr
				m.ID = "gauge" + strconv.Itoa(i)
				m.Value = &v
				_ = ms.SaveCommonModel(m)
				f, _ = ms.GetGauge("gauge" + strconv.Itoa(i%count))
			}
		}
		b.ReportMetric(float64(n)/f, "useless/op")
	})
}

func TestModelStorageSaveCommonModelsSuccess(t *testing.T) {
	ms := NewSavableModelStorage(NewMemStorage())

	const (
		counterID    = "YTGBDRTHGUJ"
		counterDelta = 456
		gaugeID      = "JYFVNHYUTRE"
		gaugeValue   = 456.789
	)

	var i64 int64 = counterDelta
	var f64 float64 = gaugeValue

	metrics := []apimodels.Metrics{
		{
			ID:    counterID,
			MType: constants.CounterStr,
			Delta: &i64,
			Value: nil,
		},
		{
			ID:    gaugeID,
			MType: constants.GaugeStr,
			Delta: nil,
			Value: &f64,
		},
	}
	t.Run("SaveCommonModels", func(t *testing.T) {
		err := ms.SaveCommonModels(context.Background(), metrics)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		counter, err := ms.GetCounter(counterID)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if counter != counterDelta {
			t.Errorf("Expected counter value %d, got %d", counterDelta, counter)
		}
		gauge, err := ms.GetGauge(gaugeID)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if gauge != gaugeValue {
			t.Errorf("Expected gauge value %f, got %f", gaugeValue, gauge)
		}

		names, err := ms.ListCounters()
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if len(names) != 1 || names[0] != counterID {
			t.Errorf("counter not listed")
		}
		names, err = ms.ListGauges()
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if len(names) != 1 || names[0] != gaugeID {
			t.Errorf("gauge not listed")
		}

	})
}

func TestModelStorageSaveCommonModelErrors(t *testing.T) {
	ms := NewSavableModelStorage(NewMemStorage())

	var i64 int64
	var f64 float64
	i64 = 456
	f64 = 456.789

	tests := []struct {
		name    string
		metrics apimodels.Metrics
		wantErr bool
	}{
		{
			name: "wrong type",
			metrics: apimodels.Metrics{
				ID:    "a",
				MType: "wrong",
				Delta: nil,
				Value: nil,
			},
			wantErr: true,
		},
		{
			name: "empty delta",
			metrics: apimodels.Metrics{
				ID:    "a",
				MType: constants.CounterStr,
				Delta: nil,
				Value: nil,
			},
			wantErr: true,
		},
		{
			name: "enmpty value",
			metrics: apimodels.Metrics{
				ID:    "a",
				MType: constants.GaugeStr,
				Delta: nil,
				Value: nil,
			},
			wantErr: true,
		},
		{
			name: "skip unused value",
			metrics: apimodels.Metrics{
				ID:    "a",
				MType: constants.CounterStr,
				Delta: &i64,
				Value: &f64,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ms.SaveCommonModel(&tt.metrics)
			if (err != nil) != tt.wantErr {
				t.Errorf("SaveCommonModel() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
