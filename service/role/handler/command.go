package handler

import (
	"context"
	pbrole "github.com/MamangRust/microservice-ecommerce-grpc-pb/role"

	"github.com/MamangRust/microservice-ecommerce-grpc-role/service"
	"github.com/MamangRust/microservice-ecommerce-pkg/logger"
	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
	"github.com/MamangRust/microservice-ecommerce-shared/errors"
	"github.com/MamangRust/microservice-ecommerce-shared/errors/role_errors"
	"google.golang.org/protobuf/types/known/emptypb"
)

type roleCommandHandler struct {
	pbrole.UnimplementedRoleCommandServiceServer
	roleCommand service.RoleCommandService
	logger      logger.LoggerInterface
}

func NewRoleCommandHandler(roleCommand service.RoleCommandService, logger logger.LoggerInterface) pbrole.RoleCommandServiceServer {
	return &roleCommandHandler{
		roleCommand: roleCommand,
		logger:      logger,
	}
}

func (s *roleCommandHandler) CreateRole(ctx context.Context, request *pbrole.CreateRoleRequest) (*pbrole.ApiResponseRole, error) {
	req := &requests.CreateRoleRequest{
		Name: request.GetName(),
	}

	if err := req.Validate(); err != nil {
		return nil, role_errors.ErrGrpcValidateCreateRole
	}

	role, err := s.roleCommand.Create(ctx, req)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pbrole.ApiResponseRole{
		Status:  "success",
		Message: "Successfully created role",
		Data:    mapToProtoRoleResponse(role),
	}, nil
}

func (s *roleCommandHandler) UpdateRole(ctx context.Context, request *pbrole.UpdateRoleRequest) (*pbrole.ApiResponseRole, error) {
	id := int(request.GetId())
	req := &requests.UpdateRoleRequest{
		ID:   &id,
		Name: request.GetName(),
	}

	if err := req.Validate(); err != nil {
		return nil, role_errors.ErrGrpcValidateUpdateRole
	}

	role, err := s.roleCommand.Update(ctx, req)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pbrole.ApiResponseRole{
		Status:  "success",
		Message: "Successfully updated role",
		Data:    mapToProtoRoleResponse(role),
	}, nil
}

func (s *roleCommandHandler) TrashedRole(ctx context.Context, request *pbrole.FindByIdRoleRequest) (*pbrole.ApiResponseRole, error) {
	id := int(request.GetRoleId())
	if id == 0 {
		return nil, role_errors.ErrGrpcRoleInvalidId
	}

	role, err := s.roleCommand.Trash(ctx, id)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pbrole.ApiResponseRole{
		Status:  "success",
		Message: "Successfully trashed role",
		Data:    mapToProtoRoleResponse(role),
	}, nil
}

func (s *roleCommandHandler) RestoreRole(ctx context.Context, request *pbrole.FindByIdRoleRequest) (*pbrole.ApiResponseRole, error) {
	id := int(request.GetRoleId())
	if id == 0 {
		return nil, role_errors.ErrGrpcRoleInvalidId
	}

	role, err := s.roleCommand.Restore(ctx, id)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pbrole.ApiResponseRole{
		Status:  "success",
		Message: "Successfully restored role",
		Data:    mapToProtoRoleResponse(role),
	}, nil
}

func (s *roleCommandHandler) DeleteRolePermanent(ctx context.Context, request *pbrole.FindByIdRoleRequest) (*pbrole.ApiResponseRoleDelete, error) {
	id := int(request.GetRoleId())
	if id == 0 {
		return nil, role_errors.ErrGrpcRoleInvalidId
	}

	_, err := s.roleCommand.DeletePermanent(ctx, id)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pbrole.ApiResponseRoleDelete{
		Status:  "success",
		Message: "Successfully deleted role permanently",
	}, nil
}

func (s *roleCommandHandler) RestoreAllRole(ctx context.Context, _ *emptypb.Empty) (*pbrole.ApiResponseRoleAll, error) {
	_, err := s.roleCommand.RestoreAll(ctx)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pbrole.ApiResponseRoleAll{
		Status:  "success",
		Message: "Successfully restored all roles",
	}, nil
}

func (s *roleCommandHandler) DeleteAllRolePermanent(ctx context.Context, _ *emptypb.Empty) (*pbrole.ApiResponseRoleAll, error) {
	_, err := s.roleCommand.DeleteAll(ctx)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pbrole.ApiResponseRoleAll{
		Status:  "success",
		Message: "Successfully deleted all roles permanently",
	}, nil
}
