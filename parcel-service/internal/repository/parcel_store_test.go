package repository

import (
	"fmt"
	"sync"
	"testing"

	"github.com/sachinggsingh/PDTS/parcel-service/internal/models"
)

func TestParcelStoreConcurrency(t *testing.T) {
	store := NewParcelStore()
	var wg sync.WaitGroup
	numGoroutines := 100
	numOperations := 100

	// Concurrent Add
	for i := range numGoroutines {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := range numOperations {
				parcelID := fmt.Sprintf("parcel-%d-%d", id, j)
				store.Add(models.ParcelMeta{
					ParcelID: parcelID,
					Status:   models.ParcelCreated,
				})
			}
		}(i)
	}
	wg.Wait()

	// Concurrent Get and Update
	for i := range numGoroutines {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := range numOperations {
				parcelID := fmt.Sprintf("parcel-%d-%d", id, j)
				meta, ok := store.Get(parcelID)
				if !ok {
					t.Errorf("expected to find parcel %s", parcelID)
					return
				}
				if meta.ParcelID != parcelID {
					t.Errorf("expected parcel ID %s, got %s", parcelID, meta.ParcelID)
				}

				err := store.Update(parcelID, models.ParcelPickedUp)
				if err != nil {
					t.Errorf("unexpected error updating parcel %s: %v", parcelID, err)
				}
			}
		}(i)
	}
	wg.Wait()

	// Concurrent Delete
	for i := range numGoroutines {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := range numOperations {
				parcelID := fmt.Sprintf("parcel-%d-%d", id, j)
				store.Delete(parcelID)
			}
		}(i)
	}
	wg.Wait()

	// Verify all deleted
	for i := range numGoroutines {
		for j := range numOperations {
			parcelID := fmt.Sprintf("parcel-%d-%d", i, j)
			_, ok := store.Get(parcelID)
			if ok {
				t.Errorf("expected parcel %s to be deleted", parcelID)
			}
		}
	}
}

func TestParcelStatusHistory(t *testing.T) {
	store := NewParcelStore()
	parcelID := "test-parcel-1"

	// 1. Add parcel
	store.Add(models.ParcelMeta{
		ParcelID: parcelID,
		Status:   models.ParcelCreated,
	})

	// 2. Update status multiple times
	statuses := []models.ParcelStatus{
		models.ParcelPickedUp,
		models.ParcelInTransit,
		models.ParcelDelivered,
	}

	for _, status := range statuses {
		err := store.Update(parcelID, status)
		if err != nil {
			t.Fatalf("failed to update status: %v", err)
		}
	}

	// 3. Get history
	history, err := store.GetHistory(parcelID)
	if err != nil {
		t.Fatalf("failed to get history: %v", err)
	}

	// 4. Verify history (Created + PickedUp + InTransit + Delivered = 4 nodes)
	expectedCount := 4
	if len(history) != expectedCount {
		t.Errorf("expected history length %d, got %d", expectedCount, len(history))
	}

	expectedStatuses := append([]models.ParcelStatus{models.ParcelCreated}, statuses...)
	for i, node := range history {
		if node.Status != expectedStatuses[i] {
			t.Errorf("at index %d: expected status %s, got %s", i, expectedStatuses[i], node.Status)
		}
	}
}
