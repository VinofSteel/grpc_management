package handlers

import "github.com/vinofsteel/grpc-management/internal/handlers/proto_user"

type Handlers struct {
	proto_user.UnimplementedUserServiceServer
}
