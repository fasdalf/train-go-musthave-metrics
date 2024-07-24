package jsonofflinestorage

import (
	"encoding/json"
	"errors"
	"github.com/fasdalf/train-go-musthave-metrics/internal/common/apimodels"
	"github.com/fasdalf/train-go-musthave-metrics/internal/common/constants"
	"io"
	"log/slog"
	"os"
	"time"
)

type Storage interface {
	GetCounter(key string) int
	GetGauge(key string) float64
	ListGauges() []string
	ListCounters() []string
	SaveCommonModel(metric *apimodels.Metrics) error
}

type JSONFileStorage struct {
	storage       Storage
	fileName      string
	restore       bool
	storeInterval time.Duration
	lastStore     time.Time
}

func NewJSONFileStorage(storage Storage, fileName string, restore bool, storeInterval int) *JSONFileStorage {
	return &JSONFileStorage{
		storage:       storage,
		fileName:      fileName,
		restore:       restore,
		storeInterval: time.Duration(storeInterval) * time.Second,
	}
}

func (l *JSONFileStorage) SaveWith2Buffers() error {
	now := time.Now()
	if l.lastStore.Add(l.storeInterval).After(now) {
		return nil
	}

	dump := []apimodels.Metrics{}

	for _, key := range l.storage.ListGauges() {
		g := l.storage.GetGauge(key)
		dump = append(dump, apimodels.Metrics{
			ID:    key,
			MType: constants.GaugeStr,
			Delta: nil,
			Value: &g,
		})
	}
	for _, key := range l.storage.ListCounters() {
		c := int64(l.storage.GetCounter(key))
		dump = append(dump, apimodels.Metrics{
			ID:    key,
			MType: constants.CounterStr,
			Delta: &c,
			Value: nil,
		})
	}

	body, err := json.MarshalIndent(dump, "", "  ")
	if err != nil {
		slog.Error("Can't generate JSON", "error", err)
		return err
	}
	err = os.WriteFile(l.fileName, body, 0660)
	if err != nil {
		slog.Error("Can't write to file", "error", err)
		return err
	}

	l.lastStore = now
	return nil
}

func (l *JSONFileStorage) SaveWithInterval() error {
	if l.storeInterval > 0 && l.lastStore.Add(l.storeInterval).After(time.Now()) {
		return nil
	}

	return l.Save()
}

func (l *JSONFileStorage) Save() error {
	file, err := os.OpenFile(l.fileName, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		slog.Error("Can't open file", "error", err)
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)

	for _, key := range l.storage.ListGauges() {
		g := l.storage.GetGauge(key)
		err = encoder.Encode(apimodels.Metrics{
			ID:    key,
			MType: constants.GaugeStr,
			Delta: nil,
			Value: &g,
		})
		if err != nil {
			slog.Error("Can't write JSON", "error", err)
			return err
		}
	}
	for _, key := range l.storage.ListCounters() {
		c := int64(l.storage.GetCounter(key))
		err = encoder.Encode(apimodels.Metrics{
			ID:    key,
			MType: constants.CounterStr,
			Delta: &c,
			Value: nil,
		})
		if err != nil {
			slog.Error("Can't write JSON", "error", err)
			return err
		}
	}

	slog.Info("Saved to file", "file", l.fileName, "error", err)
	l.lastStore = time.Now()
	return nil
}

func (l *JSONFileStorage) Restore() error {
	if !l.restore {
		return nil
	}
	file, err := os.OpenFile(l.fileName, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		slog.Error("Can't open file", "error", err)
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)

	for {
		v := &apimodels.Metrics{}
		if err = decoder.Decode(v); errors.Is(err, io.ErrUnexpectedEOF) || errors.Is(err, io.EOF) {
			// just done.
			break
		} else if err != nil {
			slog.Error("Can't decode json", "error", err)
			return err
		}

		if err = l.storage.SaveCommonModel(v); err != nil {
			slog.Error("Can't save value", "error", err)
			return err
		}
	}

	slog.Info("Restored from file", "file", l.fileName)
	return nil
}
