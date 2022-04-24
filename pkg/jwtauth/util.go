package jwtauth

import (
	"context"

	"google.golang.org/grpc/metadata"
)

func GetValueFromMetadata(md metadata.MD, key string) string {
	values := md.Get(key)
	if len(values) == 1 {
		return values[0]
	}
	return ""
}

func GetUserIdFromMetadata(md metadata.MD) string {
	return GetValueFromMetadata(md, JWT_METADATA_SUB_KEY)
}

func GetUserIdFromContext(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ""
	}
	return GetValueFromMetadata(md, JWT_METADATA_SUB_KEY)
}

func GetSessionIdFromMetadata(md metadata.MD) string {
	return GetValueFromMetadata(md, JWT_METADATA_ID_KEY)
}

func GetSessionIdFromContext(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ""
	}
	return GetSessionIdFromMetadata(md)
}
