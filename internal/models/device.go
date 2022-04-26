package models

import (
	auth "go.buf.build/grpc/go/corux/gps-auth/v1"
)

const (
	DEVICE_COLLECTION = "devices"
)

type Device struct {
	Id      string `json:"device_id" bson:"device_id"`
	OwnerId string `json:"owner_id" bson:"owner_id"`
	Name    string `json:"name" bson:"name"`
}

func (d *Device) AsProtoBuf() *auth.Device {
	return &auth.Device{
		DeviceId: d.Id,
		OwnerId:  d.OwnerId,
		Name:     d.Name,
	}
}
