package db

import (
	"github.com/golang/protobuf/proto"

	"github.com/ketchuphq/ketchup/proto/ketchup/models"
)

type AddressableProto interface {
	GetUuid() string
	proto.Message
}

type TimestampedProto interface {
	GetTimestamps() *models.Timestamp
	SetTimestamps(*models.Timestamp)
	proto.Message
}

// make sure that the protos implement addressable and timestamp
// interfaces
var (
	_ AddressableProto = &models.Page{}
	_ TimestampedProto = &models.Page{}
)
