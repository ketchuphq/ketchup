package db

import (
	"sort"

	"github.com/ketchuphq/ketchup/proto/ketchup/models"
)

var _ sort.Interface = &sortPages{}

func SortPagesByUpdatedAt(data []*models.Page, ascending bool) {
	sort.Sort(&sortPages{ascending, data})
}

type sortPages struct {
	asc  bool
	data []*models.Page
}

func (s *sortPages) Swap(i, j int) {
	s.data[i], s.data[j] = s.data[j], s.data[i]
}

func (s *sortPages) Len() int {
	return len(s.data)
}

func (s *sortPages) Less(i, j int) bool {
	a := s.data[i].GetTimestamps().GetUpdatedAt()
	b := s.data[j].GetTimestamps().GetUpdatedAt()
	if s.asc {
		return a < b
	}
	return a > b
}
