package servicev1

import (
	"context"

	modelsv1 "github.com/cory-evans/gps-tracker-authentication/internal/models/v1"
	authv1 "github.com/cory-evans/gps-tracker-authentication/pkg/auth/v1"
	jwtauthv1 "github.com/cory-evans/gps-tracker-authentication/pkg/jwtauth/v1"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func (s *AuthService) GetDevice(ctx context.Context, req *authv1.GetDeviceRequest) (*authv1.GetDeviceResponse, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Internal, "missing metadata")
	}

	userId := jwtauthv1.GetUserIdFromMetadata(md)

	devicesCol := s.DB.Collection(modelsv1.DEVICE_COLLECTION)

	var device modelsv1.Device
	result := devicesCol.FindOne(ctx, bson.M{"device_id": req.DeviceId, "owner_id": userId})
	err := result.Decode(&device)

	if err != nil {
		return nil, status.Errorf(codes.NotFound, "device not found")
	}

	return &authv1.GetDeviceResponse{
		Device: device.AsProtoBuf(),
	}, nil
}

func (s *AuthService) CreateDevice(ctx context.Context, req *authv1.CreateDeviceRequest) (*authv1.CreateDeviceResponse, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Internal, "missing metadata")
	}

	userId := jwtauthv1.GetUserIdFromMetadata(md)

	if userId == "" {
		return nil, status.Errorf(codes.Unauthenticated, "user not authenticated")
	}

	deviceName := req.GetDeviceName()
	if deviceName == "" {
		return nil, status.Error(codes.InvalidArgument, "Device name can't be none")
	}

	deviceID, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	devicesCol := s.DB.Collection(modelsv1.DEVICE_COLLECTION)
	dev := modelsv1.Device{
		Id:      deviceID.String(),
		OwnerId: userId,
		Name:    req.GetDeviceName(),
	}
	result, err := devicesCol.InsertOne(ctx, dev)

	if err != nil || result.InsertedID == nil {
		return nil, status.Errorf(codes.Internal, "failed to create device")
	}

	return &authv1.CreateDeviceResponse{
		Token:  "TODO",
		Device: dev.AsProtoBuf(),
	}, nil
}

func (s *AuthService) EditDevice(ctx context.Context, req *authv1.EditDeviceRequest) (*authv1.EditDeviceResponse, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Internal, "missing metadata")
	}

	userId := jwtauthv1.GetUserIdFromMetadata(md)

	if userId == "" {
		return nil, status.Errorf(codes.Unauthenticated, "user not authenticated")
	}

	devicesCol := s.DB.Collection(modelsv1.DEVICE_COLLECTION)

	result := devicesCol.FindOneAndUpdate(ctx, bson.M{"DeviceId": req.GetDeviceId(), "OwnerId": userId}, bson.M{"$set": bson.M{"Name": req.GetDeviceName()}})

	if result.Err() != nil {
		return nil, status.Errorf(codes.NotFound, "device not found")
	}

	return &authv1.EditDeviceResponse{}, nil
}

func (s *AuthService) GetOwnedDevices(ctx context.Context, req *authv1.GetOwnedDevicesRequest) (*authv1.GetOwnedDevicesResponse, error) {
	// check is authenticated
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "Not Authenticated")
	}
	userId := jwtauthv1.GetUserIdFromMetadata(md)

	if userId == "" {
		return nil, status.Errorf(codes.Unauthenticated, "user not authenticated")
	}

	devicesCol := s.DB.Collection(modelsv1.DEVICE_COLLECTION)

	cur, err := devicesCol.Find(ctx, bson.M{"owner_id": userId})
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "%v", err)
	}

	var devices = make([]*authv1.Device, 0)

	for cur.Next(ctx) {
		var d modelsv1.Device
		err = cur.Decode(&d)
		if err != nil {
			// TODO: Log error here
			continue
		}

		devices = append(devices, d.AsProtoBuf())
	}

	return &authv1.GetOwnedDevicesResponse{
		Devices: devices,
	}, nil
}
