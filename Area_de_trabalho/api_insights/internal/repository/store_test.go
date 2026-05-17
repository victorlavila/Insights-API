package repository

import (
	"context"
	"sync"
	"testing"
)

func TestMockStoreReturnsDefensiveCopies(t *testing.T) {
	store := NewMockStore()

	sales, err := store.ListSales(context.Background())
	if err != nil {
		t.Fatalf("ListSales retornou erro: %v", err)
	}
	sales[0].Product = "alterado fora do repositorio"

	nextSales, err := store.ListSales(context.Background())
	if err != nil {
		t.Fatalf("ListSales retornou erro: %v", err)
	}
	if nextSales[0].Product == "alterado fora do repositorio" {
		t.Fatal("repositorio vazou slice interno")
	}
}

func TestMockStoreConcurrentReads(t *testing.T) {
	store := NewMockStore()
	var wg sync.WaitGroup

	for range 100 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if _, err := store.ListSales(context.Background()); err != nil {
				t.Errorf("ListSales retornou erro: %v", err)
			}
			if _, err := store.ListSuggestions(context.Background()); err != nil {
				t.Errorf("ListSuggestions retornou erro: %v", err)
			}
			if _, err := store.ListComplaints(context.Background()); err != nil {
				t.Errorf("ListComplaints retornou erro: %v", err)
			}
			if _, err := store.ListUserCaptures(context.Background()); err != nil {
				t.Errorf("ListUserCaptures retornou erro: %v", err)
			}
		}()
	}

	wg.Wait()
}
