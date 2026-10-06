package handler

import (
	pbcommon "github.com/MamangRust/microservice-ecommerce-grpc-pb/common"
	pbrole "github.com/MamangRust/microservice-ecommerce-grpc-pb/role"
	pbuserrole "github.com/MamangRust/microservice-ecommerce-grpc-pb/user_role"
	"math"
	"time"

	"github.com/MamangRust/microservice-ecommerce-grpc-role/repository"
	"github.com/MamangRust/microservice-ecommerce-pkg/database/models"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func normalizePage(page, pageSize int) (int, int) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	return page, pageSize
}

func createPaginationMeta(page, pageSize, totalRecords int) *pbcommon.PaginationMeta {
	totalPages := int(math.Ceil(float64(totalRecords) / float64(pageSize)))
	return &pbcommon.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(totalRecords),
	}
}

func formatTimePtr(t *time.Time) string {
	if t != nil {
		return t.Format("2006-01-02 15:04:05.000")
	}
	return ""
}

// mapToProtoRoleResult maps a RoleResult to RoleResponse
func mapToProtoRoleResult(r *repository.RoleResult) *pbrole.RoleResponse {
	if r == nil {
		return nil
	}
	return &pbrole.RoleResponse{
		Id:        r.RoleID,
		Name:      r.RoleName,
		CreatedAt: formatTimePtr(r.CreatedAt),
		UpdatedAt: formatTimePtr(r.UpdatedAt),
	}
}

// mapToProtoRoleResultDeleteAt maps a RoleResult to RoleResponseDeleteAt
func mapToProtoRoleResultDeleteAt(r *repository.RoleResult) *pbrole.RoleResponseDeleteAt {
	if r == nil {
		return nil
	}
	res := &pbrole.RoleResponseDeleteAt{
		Id:        r.RoleID,
		Name:      r.RoleName,
		CreatedAt: formatTimePtr(r.CreatedAt),
		UpdatedAt: formatTimePtr(r.UpdatedAt),
	}
	if r.DeletedAt != nil {
		res.DeletedAt = &wrapperspb.StringValue{Value: formatTimePtr(r.DeletedAt)}
	}
	return res
}

// mapToProtoRoleResponse maps a *models.Role to RoleResponse
func mapToProtoRoleResponse(m *models.Role) *pbrole.RoleResponse {
	if m == nil {
		return nil
	}
	return &pbrole.RoleResponse{
		Id:        m.RoleID,
		Name:      m.RoleName,
		CreatedAt: formatTimePtr(m.CreatedAt),
		UpdatedAt: formatTimePtr(m.UpdatedAt),
	}
}

// mapToProtoRoleResponseDeleteAt maps a *models.Role to RoleResponseDeleteAt
func mapToProtoRoleResponseDeleteAt(m *models.Role) *pbrole.RoleResponseDeleteAt {
	if m == nil {
		return nil
	}
	res := &pbrole.RoleResponseDeleteAt{
		Id:        m.RoleID,
		Name:      m.RoleName,
		CreatedAt: formatTimePtr(m.CreatedAt),
		UpdatedAt: formatTimePtr(m.UpdatedAt),
	}
	if m.DeletedAt != nil {
		res.DeletedAt = &wrapperspb.StringValue{Value: formatTimePtr(m.DeletedAt)}
	}
	return res
}

func mapToProtoUserRoleResponse(v *models.UserRole) *pbuserrole.UserRoleResponse {
	if v == nil {
		return nil
	}
	return &pbuserrole.UserRoleResponse{
		UserRoleId: v.UserRoleID,
		UserId:     v.UserID,
		RoleId:     v.RoleID,
	}
}
