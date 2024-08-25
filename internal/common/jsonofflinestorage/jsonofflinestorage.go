package jsonofflinestorage

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/fasdalf/train-go-musthave-metrics/internal/common/apimodels"
	"github.com/fasdalf/train-go-musthave-metrics/internal/common/constants"
	"io"
	"io/fs"
	"log/slog"
	"os"
	"sync"
	"time"
)

type Storage interface {
	GetCounter(key string) (int, error)
	GetGauge(key string) (float64, error)
	ListGauges() ([]string, error)
	ListCounters() ([]string, error)
	SaveCommonModel(metric *apimodels.Metrics) error
}

type JSONFileStorage struct {
	storage       Storage
	fileName      string
	restore       bool
	storeInterval time.Duration
	mu            *sync.Mutex
}

func NewJSONFileStorage(storage Storage, fileName string, restore bool, storeInterval int) *JSONFileStorage {
	return &JSONFileStorage{
		storage:       storage,
		fileName:      fileName,
		restore:       restore,
		storeInterval: time.Duration(storeInterval) * time.Second,
		mu:            new(sync.Mutex),
	}
}

func (l *JSONFileStorage) Save() error {
	l.mu.Lock()
	defer l.mu.Unlock()
	file, err := os.OpenFile(l.fileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return fmt.Errorf("can't open file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)

	list, err := l.storage.ListGauges()
	if err != nil {
		return fmt.Errorf("can't ListGauges: %w", err)
	}
	for _, key := range list {
		g, err := l.storage.GetGauge(key)
		if err != nil {
			return fmt.Errorf("can't GetGauge: %w", err)
		}
		err = encoder.Encode(apimodels.Metrics{
			ID:    key,
			MType: constants.GaugeStr,
			Delta: nil,
			Value: &g,
		})
		if err != nil {
			return fmt.Errorf("can't write JSON: %w", err)
		}
	}

	list, err = l.storage.ListCounters()
	if err != nil {
		return fmt.Errorf("can't ListCounters: %w", err)
	}
	for _, key := range list {
		c, err := l.storage.GetCounter(key)
		if err != nil {
			return fmt.Errorf("can't GetCounter: %w", err)
		}
		c64 := int64(c)
		err = encoder.Encode(apimodels.Metrics{
			ID:    key,
			MType: constants.CounterStr,
			Delta: &c64,
			Value: nil,
		})
		if err != nil {
			return fmt.Errorf("can't write JSON: %w", err)
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
		return fmt.Errorf("can't open file: %w", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)

	for {
		v := &apimodels.Metrics{}
		if err = decoder.Decode(v); errors.Is(err, io.ErrUnexpectedEOF) || errors.Is(err, io.EOF) {
			// just done.
			break
		} else if err != nil {
			return fmt.Errorf("can't decode json: %w", err)
		}

		if err = l.storage.SaveCommonModel(v); err != nil {
			return fmt.Errorf("can't save value: %w", err)
		}
	}

	slog.Info("Restored from file", "file", l.fileName)
	return nil
}

type SavedChan = chan struct{}

func (l *JSONFileStorage) SaveMetrics(ctx context.Context, saved SavedChan, wg *sync.WaitGroup) error {
	defer wg.Done()
	do := true
	if l.storeInterval > 0 {
		t := time.NewTimer(l.storeInterval)
		defer t.Stop()
		for do {
			select {
			case <-ctx.Done():
				do = false
			case <-t.C:
			}

			if err := l.Save(); err != nil {
				slog.Error("SaveMetrics error", "error", err)
			}

			if do {
				t.Reset(l.storeInterval)
			}
		}
		return nil
	}

	for do {
		select {
		case <-ctx.Done():
			do = false
		case <-saved:
		}

		if err := l.Save(); err != nil {
			slog.Error("SaveMetrics error", "error", err)
		}
	}
	return nil
}
