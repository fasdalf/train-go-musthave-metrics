package jsonofflinestorage

import (
	"encoding/json"
	"errors"
	"github.com/fasdalf/train-go-musthave-metrics/internal/common/apimodels"
	"github.com/fasdalf/train-go-musthave-metrics/internal/common/constants"
	"io"
	"io/fs"
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
}

func NewJSONFileStorage(storage Storage, fileName string, restore bool, storeInterval int) *JSONFileStorage {
	return &JSONFileStorage{
		storage:       storage,
		fileName:      fileName,
		restore:       restore,
		storeInterval: time.Duration(storeInterval) * time.Second,
	}
}

func (l *JSONFileStorage) Save() error {
	file, err := os.OpenFile(l.fileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
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
	return nil
}

func (l *JSONFileStorage) Restore() error {
	if !l.restore {
		return nil
	}
	file, err := os.OpenFile(l.fileName, os.O_RDONLY, 0666)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			slog.Warn("file does not exist, skipping JSON load", "error", err, "filename", l.fileName)
			return nil
		}
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

type SavedChan = chan struct{}

func (l *JSONFileStorage) SaverRoutine(saved SavedChan) error {
	t := time.NewTimer(l.storeInterval)
	t.Stop()
	if l.storeInterval > 0 {
		t.Reset(l.storeInterval)
	}

	for {
		select {
		case <-saved:
		case <-t.C:
		}

		if err := l.Save(); err != nil {
			slog.Error("SaverRoutine error", "error", err)
		}

		if l.storeInterval > 0 {
			t.Stop()
			t.Reset(l.storeInterval)
		}
	}
}
