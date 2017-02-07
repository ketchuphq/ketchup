package db

import (
	"testing"
	"time"

	"github.com/octavore/ketchup/proto/ketchup/models"
)

func TestSort(t *testing.T) {
	now := time.Now()
	nowUnix := now.UnixNano() / 1e6
	nowPlus1 := now.Add(time.Minute).UnixNano() / 1e6
	nowPlus2 := now.Add(time.Minute*2).UnixNano() / 1e6
	pages := []*models.Page{
		{Timestamps: &models.Timestamp{UpdatedAt: &nowUnix}},
		{Timestamps: &models.Timestamp{UpdatedAt: &nowPlus2}},
		{Timestamps: &models.Timestamp{UpdatedAt: &nowPlus1}},
	}

	SortPagesByUpdatedAt(pages, true)
	expected := []int64{nowUnix, nowPlus1, nowPlus2}
	for i, n := range expected {
		ps := pages[i].GetTimestamps().GetUpdatedAt()
		if ps != n {
			t.Fatalf("expected %v but got %v at %v", n, ps, i)
		}
	}

	SortPagesByUpdatedAt(pages, false)
	expected = []int64{nowPlus2, nowPlus1, nowUnix}
	for i, n := range expected {
		ps := pages[i].GetTimestamps().GetUpdatedAt()
		if ps != n {
			t.Fatalf("expected %v but got %v at %v", n, ps, i)
		}
	}
}
