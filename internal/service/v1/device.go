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

	devicesCol := s.DB.Collection("devices")

	var device modelsv1.Device
	result := devicesCol.FindOne(ctx, bson.M{"DeviceId": req.DeviceId, "OwnerId": userId})
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

	deviceID, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	devicesCol := s.DB.Collection("devices")
	result, err := devicesCol.InsertOne(ctx, modelsv1.Device{
		Id:      deviceID.String(),
		OwnerId: userId,
		Name:    req.GetDeviceName(),
	})

	if err != nil || result.InsertedID == nil {
		return nil, status.Errorf(codes.Internal, "failed to create device")
	}

	return &authv1.CreateDeviceResponse{
		Device: &authv1.Device{
			DeviceId: deviceID.String(),
			OwnerId:  userId,
			Name:     req.GetDeviceName(),
		},
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

	devicesCol := s.DB.Collection("devices")

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
	userID := jwtauthv1.GetUserIdFromMetadata(md)

	devicesCol := s.DB.Collection("devices")

	cur, err := devicesCol.Find(ctx, bson.M{"OwnerId": userID})
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
