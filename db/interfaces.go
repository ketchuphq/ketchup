package db

import (
	"github.com/golang/protobuf/proto"

	"github.com/octavore/ketchup/proto/ketchup/models"
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
