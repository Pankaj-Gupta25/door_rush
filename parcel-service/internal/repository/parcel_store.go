package repository

import (
	"fmt"
	"sync"

	"github.com/sachinggsingh/PDTS/parcel-service/internal/models"
)

type ParcelStore interface {
	Add(parcelMeta models.ParcelMeta)
	Get(parcelID string) (models.ParcelMeta, bool)
	Delete(parcelID string)
	Update(parcelID string, status models.ParcelStatus) error
	GetHistory(parcelID string) ([]models.StatusNode, error)
}

type parcelStore struct {
	mu      sync.RWMutex
	parcels map[string]models.ParcelMeta
}

func NewParcelStore() ParcelStore {
	return &parcelStore{
		parcels: make(map[string]models.ParcelMeta),
	}
}

func (s *parcelStore) Add(parcelMeta models.ParcelMeta) {
	s.mu.Lock()
	defer s.mu.Unlock()

	fmt.Printf("[ParcelStore] Adding parcel: %s\n", parcelMeta.ParcelID)

	// If the parcel already exists, preserve its history
	if existing, ok := s.parcels[parcelMeta.ParcelID]; ok && parcelMeta.History == nil {
		fmt.Printf("[ParcelStore] Preserving history for parcel: %s\n", parcelMeta.ParcelID)
		parcelMeta.History = existing.History
	}

	// If the history is still nil, create a new history and append the status
	if parcelMeta.History == nil {
		fmt.Printf("[ParcelStore] Initializing new history for parcel: %s with status: %s\n", parcelMeta.ParcelID, parcelMeta.Status)
		parcelMeta.History = &models.StatusHistory{}
		parcelMeta.History.Append(parcelMeta.Status)
	}

	// Add the parcel to the store
	s.parcels[parcelMeta.ParcelID] = parcelMeta
}

func (s *parcelStore) Get(parcelID string) (models.ParcelMeta, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	// Get the parcel from the store
	meta, ok := s.parcels[parcelID]
	return meta, ok
}

func (s *parcelStore) Delete(parcelID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	// Delete the parcel from the store
	delete(s.parcels, parcelID)
}

func (s *parcelStore) Update(parcelID string, status models.ParcelStatus) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Get the parcel from the store
	meta, ok := s.parcels[parcelID]
	if !ok {
		return fmt.Errorf("parcel with ID %s not found in store", parcelID)
	}

	// Update the status and history
	if meta.Status != status {
		meta.Status = status
		if meta.History == nil {
			meta.History = &models.StatusHistory{}
		}
		meta.History.Append(status)
	}

	s.parcels[parcelID] = meta
	return nil
}

func (s *parcelStore) GetHistory(parcelID string) ([]models.StatusNode, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	fmt.Printf("[ParcelStore] Getting history for parcel: %s\n", parcelID)

	// Get the parcel from the store
	meta, ok := s.parcels[parcelID]
	if !ok {
		fmt.Printf("[ParcelStore] Parcel %s NOT FOUND in store\n", parcelID)
		return nil, fmt.Errorf("parcel with ID %s not found in store", parcelID)
	}

	// Get the history
	if meta.History == nil {
		fmt.Printf("[ParcelStore] History for parcel %s is NIL\n", parcelID)
		return []models.StatusNode{}, nil
	}

	nodes := meta.History.GetAll()
	fmt.Printf("[ParcelStore] Found %d history nodes for parcel: %s\n", len(nodes), parcelID)
	return nodes, nil
}
