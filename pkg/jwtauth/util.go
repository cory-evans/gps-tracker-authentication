package jwtauth

import (
	"context"

	"google.golang.org/grpc/metadata"
)

func GetUserIdFromMetadata(md metadata.MD) string {
	values := md.Get(JWT_METADATA_SUB_KEY)
	if len(values) == 1 {
		return values[0]
	}
	return ""
}

func GetUserIdFromContext(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ""
	}
	return GetUserIdFromMetadata(md)
}
