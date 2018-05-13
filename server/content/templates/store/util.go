package store

import (
	"github.com/golang/protobuf/proto"

	"github.com/ketchuphq/ketchup/proto/ketchup/models"
)

func StripData(in *models.Theme) *models.Theme {
	out := proto.Clone(in).(*models.Theme)
	for _, template := range out.Templates {
		template.Data = nil
	}
	for _, asset := range out.Assets {
		asset.Data = nil
	}
	return out
}
