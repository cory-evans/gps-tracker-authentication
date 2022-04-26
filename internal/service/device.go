package service

import (
	"context"

	"github.com/cory-evans/gps-tracker-authentication/internal/models"
	"github.com/cory-evans/gps-tracker-authentication/pkg/jwtauth"
	"github.com/google/uuid"
	auth "go.buf.build/grpc/go/corux/gps-auth/v1"
	"go.mongodb.org/mongo-driver/bson"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func (s *AuthService) GetDevice(ctx context.Context, req *auth.GetDeviceRequest) (*auth.GetDeviceResponse, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Internal, "missing metadata")
	}

	userId := jwtauth.GetUserIdFromMetadata(md)

	if userId == "" {
		return nil, status.Errorf(codes.Unauthenticated, "user not authenticated")
	}

	devicesCol := s.DB.Collection(models.DEVICE_COLLECTION)

	var device models.Device
	result := devicesCol.FindOne(ctx, bson.M{"device_id": req.DeviceId})
	err := result.Decode(&device)

	if err != nil {
		return nil, status.Errorf(codes.NotFound, "device not found")
	}

	return &auth.GetDeviceResponse{
		Device: device.AsProtoBuf(),
	}, nil
}

func (s *AuthService) CreateDevice(ctx context.Context, req *auth.CreateDeviceRequest) (*auth.CreateDeviceResponse, error) {
	userId := jwtauth.GetUserIdFromContext(ctx)

	deviceName := req.GetDeviceName()
	if deviceName == "" {
		return nil, status.Error(codes.InvalidArgument, "Device name can't be none")
	}

	deviceID, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	devicesCol := s.DB.Collection(models.DEVICE_COLLECTION)
	dev := models.Device{
		Id:      deviceID.String(),
		OwnerId: userId,
		Name:    req.GetDeviceName(),
	}
	result, err := devicesCol.InsertOne(ctx, dev)

	if err != nil || result.InsertedID == nil {
		return nil, status.Errorf(codes.Internal, "failed to create device")
	}

	// create session for the device
	token, sess, err := s.createNewDeviceSession(ctx, deviceID.String())

	return &auth.CreateDeviceResponse{
		Token:        token,
		Device:       dev.AsProtoBuf(),
		RefreshToken: sess.RefreshToken,
		ExpiresAtUtc: sess.ExpiresAtUtc,
	}, err
}

func (s *AuthService) EditDevice(ctx context.Context, req *auth.EditDeviceRequest) (*auth.EditDeviceResponse, error) {
	userId := jwtauth.GetUserIdFromContext(ctx)

	devicesCol := s.DB.Collection(models.DEVICE_COLLECTION)

	result := devicesCol.FindOneAndUpdate(ctx, bson.M{"DeviceId": req.GetDeviceId(), "OwnerId": userId}, bson.M{"$set": bson.M{"Name": req.GetDeviceName()}})

	if result.Err() != nil {
		return nil, status.Errorf(codes.NotFound, "device not found")
	}

	return &auth.EditDeviceResponse{}, nil
}

func (s *AuthService) GetOwnedDevices(ctx context.Context, req *auth.GetOwnedDevicesRequest) (*auth.GetOwnedDevicesResponse, error) {
	userId := jwtauth.GetUserIdFromContext(ctx)

	devicesCol := s.DB.Collection(models.DEVICE_COLLECTION)

	cur, err := devicesCol.Find(ctx, bson.M{"owner_id": userId})
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "%v", err)
	}

	var devices = make([]*auth.Device, 0)

	for cur.Next(ctx) {
		var d models.Device
		err = cur.Decode(&d)
		if err != nil {
			// TODO: Log error here
			continue
		}

		devices = append(devices, d.AsProtoBuf())
	}

	return &auth.GetOwnedDevicesResponse{
		Devices: devices,
	}, nil
}
