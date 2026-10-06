package handler

import (
	pbbanner "github.com/MamangRust/microservice-ecommerce-grpc-pb/banner"
	pbcommon "github.com/MamangRust/microservice-ecommerce-grpc-pb/common"
	"math"
	"time"

	"github.com/MamangRust/microservice-ecommerce-grpc-banner/repository"
	"github.com/MamangRust/microservice-ecommerce-pkg/database/models"
	"github.com/MamangRust/microservice-ecommerce-shared/convert"
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

func fmtTime(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format("2006-01-02 15:04:05.000")
}

func fmtDate(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format("2006-01-02")
}

func fmtTimeOnlyStr(t *string) string {
	if t == nil {
		return ""
	}
	return *t
}

func boolVal(b *bool) bool {
	if b == nil {
		return false
	}
	return *b
}

func mapToProtoBannerResponse(m interface{}) *pbbanner.BannerResponse {
	switch v := m.(type) {
	case *models.Banner:
		return &pbbanner.BannerResponse{
			BannerId:  v.BannerID,
			Name:      v.Name,
			StartDate: fmtDate(v.StartDate),
			EndDate:   fmtDate(v.EndDate),
			StartTime: fmtTimeOnlyStr(v.StartTime),
			EndTime:   fmtTimeOnlyStr(v.EndTime),
			IsActive:  boolVal(v.IsActive),
			CreatedAt: fmtTime(v.CreatedAt),
			UpdatedAt: fmtTime(v.UpdatedAt),
		}
	case *repository.BannerResult:
		return &pbbanner.BannerResponse{
			BannerId:  v.BannerID,
			Name:      v.Name,
			StartDate: convert.StrVal(v.StartDate),
			EndDate:   convert.StrVal(v.EndDate),
			StartTime: convert.StrVal(v.StartTime),
			EndTime:   convert.StrVal(v.EndTime),
			IsActive:  boolVal(v.IsActive),
			CreatedAt: convert.StrVal(v.CreatedAt),
			UpdatedAt: convert.StrVal(v.UpdatedAt),
		}
	default:
		return nil
	}
}

func mapToProtoBannerResponseDeleteAt(m interface{}) *pbbanner.BannerResponseDeleteAt {
	switch v := m.(type) {
	case *models.Banner:
		res := &pbbanner.BannerResponseDeleteAt{
			BannerId:  v.BannerID,
			Name:      v.Name,
			StartDate: fmtDate(v.StartDate),
			EndDate:   fmtDate(v.EndDate),
			StartTime: fmtTimeOnlyStr(v.StartTime),
			EndTime:   fmtTimeOnlyStr(v.EndTime),
			IsActive:  boolVal(v.IsActive),
			CreatedAt: fmtTime(v.CreatedAt),
			UpdatedAt: fmtTime(v.UpdatedAt),
		}
		if v.DeletedAt != nil {
			res.DeletedAt = &wrapperspb.StringValue{Value: fmtTime(v.DeletedAt)}
		}
		return res
	case *repository.BannerResult:
		res := &pbbanner.BannerResponseDeleteAt{
			BannerId:  v.BannerID,
			Name:      v.Name,
			StartDate: convert.StrVal(v.StartDate),
			EndDate:   convert.StrVal(v.EndDate),
			StartTime: convert.StrVal(v.StartTime),
			EndTime:   convert.StrVal(v.EndTime),
			IsActive:  boolVal(v.IsActive),
			CreatedAt: convert.StrVal(v.CreatedAt),
			UpdatedAt: convert.StrVal(v.UpdatedAt),
		}
		if v.DeletedAt != nil && *v.DeletedAt != "" {
			res.DeletedAt = &wrapperspb.StringValue{Value: *v.DeletedAt}
		}
		return res
	default:
		return nil
	}
}
