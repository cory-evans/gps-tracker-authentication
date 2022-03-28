package modelsv1

import authv1 "github.com/cory-evans/gps-tracker-authentication/pkg/auth/v1"

type Device struct {
	Id      string `json:"device_id" bson:"device_id"`
	OwnerId string `json:"owner_id" bson:"owner_id"`
	Name    string `json:"name" bson:"name"`
}

func (d *Device) AsProtoBuf() *authv1.Device {
	return &authv1.Device{
		DeviceId: d.Id,
		OwnerId:  d.OwnerId,
	}
}
