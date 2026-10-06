package handler

import (
	"context"

	pbrole "github.com/MamangRust/microservice-ecommerce-grpc-pb/role"
	pbuserrole "github.com/MamangRust/microservice-ecommerce-grpc-pb/user_role"

	"github.com/MamangRust/microservice-ecommerce-grpc-role/service"
	"github.com/MamangRust/microservice-ecommerce-pkg/logger"
	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
	"github.com/MamangRust/microservice-ecommerce-shared/errors"
	"github.com/MamangRust/microservice-ecommerce-shared/errors/role_errors"
	"google.golang.org/protobuf/types/known/emptypb"
)

type userRoleHandler struct {
	pbuserrole.UnimplementedUserRoleServiceServer
	roleQuery   service.RoleQueryService
	roleCommand service.RoleCommandService
	logger      logger.LoggerInterface
}

func NewUserRoleHandler(roleQuery service.RoleQueryService, roleCommand service.RoleCommandService, logger logger.LoggerInterface) pbuserrole.UserRoleServiceServer {
	return &userRoleHandler{
		roleQuery:   roleQuery,
		roleCommand: roleCommand,
		logger:      logger,
	}
}

func (s *userRoleHandler) FindByUserId(ctx context.Context, req *pbuserrole.FindByIdUserRoleRequest) (*pbrole.ApiResponsesRole, error) {
	userID := int(req.GetUserId())
	if userID == 0 {
		return nil, role_errors.ErrGrpcRoleInvalidId
	}

	roles, err := s.roleQuery.FindByUserId(ctx, userID)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoRoles := make([]*pbrole.RoleResponse, len(roles))
	for i, role := range roles {
		protoRoles[i] = mapToProtoRoleResponse(role)
	}

	return &pbrole.ApiResponsesRole{
		Status:  "success",
		Message: "Successfully fetched role by user id",
		Data:    protoRoles,
	}, nil
}

func (s *userRoleHandler) AssignRoleToUser(ctx context.Context, request *pbuserrole.AssignRoleToUserRequest) (*pbuserrole.ApiResponseUserRole, error) {
	req := &requests.CreateUserRoleRequest{
		UserId: int(request.GetUserId()),
		RoleId: int(request.GetRoleId()),
	}

	userRole, err := s.roleCommand.AssignRoleToUser(ctx, req)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pbuserrole.ApiResponseUserRole{
		Status:  "success",
		Message: "Successfully assigned role to user",
		Data:    mapToProtoUserRoleResponse(userRole),
	}, nil
}

func (s *userRoleHandler) RemoveRoleFromUser(ctx context.Context, request *pbuserrole.RemoveRoleFromUserRequest) (*emptypb.Empty, error) {
	req := &requests.RemoveUserRoleRequest{
		UserId: int(request.GetUserId()),
		RoleId: int(request.GetRoleId()),
	}

	if err := s.roleCommand.RemoveRoleFromUser(ctx, req); err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &emptypb.Empty{}, nil
}
