package metricstorage

import "sync"

type DirtyStorage struct {
	basicStorage
	SavedChan chan struct{}
	isDirty   bool
	dirtyMx   *sync.Mutex
}

// NewDirtyStorage creates new storage
func NewDirtyStorage(storage basicStorage) *DirtyStorage {
	return &DirtyStorage{
		basicStorage: storage,
		isDirty:      false,
		SavedChan:    make(chan struct{}),
		dirtyMx:      &sync.Mutex{},
	}
}

func (ds *DirtyStorage) UpdateCounter(key string, value int) error {
	err := ds.basicStorage.UpdateCounter(key, value)
	if err == nil {
		ds.markDirty()
	}
	return err
}

func (ds *DirtyStorage) UpdateGauge(key string, value float64) error {
	err := ds.basicStorage.UpdateGauge(key, value)
	if err == nil {
		ds.markDirty()
	}
	return err
}

func (ds *DirtyStorage) markDirty() {
	ds.dirtyMx.Lock()
	defer ds.dirtyMx.Unlock()
	ds.isDirty = true
	select {
	case ds.SavedChan <- struct{}{}:
	default:
	}
}

// Clear removes dirty flag and returns its current value.
func (ds *DirtyStorage) Clear() bool {
	ds.dirtyMx.Lock()
	defer ds.dirtyMx.Unlock()
	result := ds.isDirty
	ds.isDirty = false
	return result
}
