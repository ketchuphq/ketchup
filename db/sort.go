package db

import (
	"sort"

	"github.com/ketchuphq/ketchup/proto/ketchup/models"
)

func SortPagesByUpdatedAt(data []*models.Page, ascending bool) {
	sort.Slice(data, func(i, j int) bool {
		a := data[i].GetTimestamps().GetUpdatedAt()
		b := data[j].GetTimestamps().GetUpdatedAt()
		if ascending {
			return a < b
		}
		return a > b
	})
}

func SortRoutesByPath(data []*models.Route) []*models.Route {
	sort.Slice(data, func(i, j int) bool {
		return data[i].GetPath() < data[j].GetPath()
	})
	return data
}
