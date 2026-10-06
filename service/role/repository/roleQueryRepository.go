package repository

import (
	"context"
	"time"

	"github.com/MamangRust/microservice-ecommerce-pkg/database/models"
	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
	sharedErrors "github.com/MamangRust/microservice-ecommerce-shared/errors"
	"github.com/MamangRust/microservice-ecommerce-shared/errors/role_errors"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type roleQueryRepository struct {
	db *mongo.Database
}

func NewRoleQueryRepository(db *mongo.Database) *roleQueryRepository {
	return &roleQueryRepository{db: db}
}

func (r *roleQueryRepository) FindAll(ctx context.Context, req *requests.FindAllRole) ([]*RoleResult, error) {
	offset := (req.Page - 1) * req.PageSize
	coll := r.db.Collection("roles")

	filter := bson.M{"deleted_at": bson.M{"$eq": nil}}
	if req.Search != "" {
		filter["role_name"] = bson.M{"$regex": req.Search, "$options": "i"}
	}

	pipeline := []bson.M{
		{"$match": filter},
		{"$sort": bson.M{"created_at": 1}},
		{"$facet": bson.M{
			"metadata": []bson.M{{"$count": "total_count"}},
			"data": []bson.M{
				{"$skip": offset},
				{"$limit": req.PageSize},
			},
		}},
	}

	cursor, err := coll.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, role_errors.ErrFindAllRoles.WithInternal(err)
	}
	defer cursor.Close(ctx)

	var facetResults []struct {
		Metadata []struct{ TotalCount int64 `bson:"total_count"` } `bson:"metadata"`
		Data     []bson.M `bson:"data"`
	}
	if err := cursor.All(ctx, &facetResults); err != nil {
		return nil, role_errors.ErrFindAllRoles.WithInternal(err)
	}
	if len(facetResults) == 0 {
		return []*RoleResult{}, nil
	}

	var results []*RoleResult
	for _, d := range facetResults[0].Data {
		item := &RoleResult{
			RoleID:    getInt32(d, "_id"),
			RoleName:  getString(d, "role_name"),
			CreatedAt: bsonMToTimePtr(d, "created_at"),
			UpdatedAt: bsonMToTimePtr(d, "updated_at"),
			DeletedAt: bsonMToTimePtr(d, "deleted_at"),
		}
		if len(facetResults[0].Metadata) > 0 {
			item.TotalCount = facetResults[0].Metadata[0].TotalCount
		}
		results = append(results, item)
	}
	return results, nil
}

func (r *roleQueryRepository) FindByID(ctx context.Context, id int) (*models.Role, error) {
	coll := r.db.Collection("roles")
	var role models.Role
	err := coll.FindOne(ctx, bson.M{"_id": id, "deleted_at": bson.M{"$eq": nil}}).Decode(&role)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, role_errors.ErrRoleNotFound.WithInternal(err)
		}
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}
	return &role, nil
}

func (r *roleQueryRepository) FindByName(ctx context.Context, name string) (*models.Role, error) {
	coll := r.db.Collection("roles")
	var role models.Role
	err := coll.FindOne(ctx, bson.M{"role_name": name, "deleted_at": bson.M{"$eq": nil}}).Decode(&role)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, role_errors.ErrRoleNotFound.WithInternal(err)
		}
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}
	return &role, nil
}

func (r *roleQueryRepository) FindByUserId(ctx context.Context, userID int) ([]*models.Role, error) {
	urColl := r.db.Collection("user_roles")

	// Find role_ids assigned to this user
	urCursor, err := urColl.Find(ctx, bson.M{"user_id": userID, "deleted_at": bson.M{"$eq": nil}})
	if err != nil {
		return nil, role_errors.ErrRoleNotFound.WithInternal(err)
	}
	defer urCursor.Close(ctx)

	var userRoles []bson.M
	if err := urCursor.All(ctx, &userRoles); err != nil {
		return nil, role_errors.ErrRoleNotFound.WithInternal(err)
	}

	var roleIDs []int32
	for _, ur := range userRoles {
		roleIDs = append(roleIDs, getInt32(ur, "role_id"))
	}
	if len(roleIDs) == 0 {
		return []*models.Role{}, nil
	}

	rolesColl := r.db.Collection("roles")
	rolesCursor, err := rolesColl.Find(ctx, bson.M{
		"_id":        bson.M{"$in": roleIDs},
		"deleted_at": bson.M{"$eq": nil},
	})
	if err != nil {
		return nil, role_errors.ErrRoleNotFound.WithInternal(err)
	}
	defer rolesCursor.Close(ctx)

	var roles []*models.Role
	if err := rolesCursor.All(ctx, &roles); err != nil {
		return nil, role_errors.ErrRoleNotFound.WithInternal(err)
	}
	return roles, nil
}

func (r *roleQueryRepository) FindActive(ctx context.Context, req *requests.FindAllRole) ([]*RoleResult, error) {
	offset := (req.Page - 1) * req.PageSize
	coll := r.db.Collection("roles")

	filter := bson.M{"deleted_at": bson.M{"$eq": nil}}
	if req.Search != "" {
		filter["role_name"] = bson.M{"$regex": req.Search, "$options": "i"}
	}

	pipeline := []bson.M{
		{"$match": filter},
		{"$sort": bson.M{"created_at": 1}},
		{"$facet": bson.M{
			"metadata": []bson.M{{"$count": "total_count"}},
			"data": []bson.M{
				{"$skip": offset},
				{"$limit": req.PageSize},
			},
		}},
	}

	cursor, err := coll.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, role_errors.ErrFindActiveRoles.WithInternal(err)
	}
	defer cursor.Close(ctx)

	var facetResults []struct {
		Metadata []struct{ TotalCount int64 `bson:"total_count"` } `bson:"metadata"`
		Data     []bson.M `bson:"data"`
	}
	if err := cursor.All(ctx, &facetResults); err != nil {
		return nil, role_errors.ErrFindActiveRoles.WithInternal(err)
	}
	if len(facetResults) == 0 {
		return []*RoleResult{}, nil
	}

	var results []*RoleResult
	for _, d := range facetResults[0].Data {
		item := &RoleResult{
			RoleID:    getInt32(d, "_id"),
			RoleName:  getString(d, "role_name"),
			CreatedAt: bsonMToTimePtr(d, "created_at"),
			UpdatedAt: bsonMToTimePtr(d, "updated_at"),
			DeletedAt: bsonMToTimePtr(d, "deleted_at"),
		}
		if len(facetResults[0].Metadata) > 0 {
			item.TotalCount = facetResults[0].Metadata[0].TotalCount
		}
		results = append(results, item)
	}
	return results, nil
}

func (r *roleQueryRepository) FindTrashed(ctx context.Context, req *requests.FindAllRole) ([]*RoleResult, error) {
	offset := (req.Page - 1) * req.PageSize
	coll := r.db.Collection("roles")

	filter := bson.M{"deleted_at": bson.M{"$ne": nil}}
	if req.Search != "" {
		filter["role_name"] = bson.M{"$regex": req.Search, "$options": "i"}
	}

	pipeline := []bson.M{
		{"$match": filter},
		{"$sort": bson.M{"created_at": 1}},
		{"$facet": bson.M{
			"metadata": []bson.M{{"$count": "total_count"}},
			"data": []bson.M{
				{"$skip": offset},
				{"$limit": req.PageSize},
			},
		}},
	}

	cursor, err := coll.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, role_errors.ErrFindTrashedRoles.WithInternal(err)
	}
	defer cursor.Close(ctx)

	var facetResults []struct {
		Metadata []struct{ TotalCount int64 `bson:"total_count"` } `bson:"metadata"`
		Data     []bson.M `bson:"data"`
	}
	if err := cursor.All(ctx, &facetResults); err != nil {
		return nil, role_errors.ErrFindTrashedRoles.WithInternal(err)
	}
	if len(facetResults) == 0 {
		return []*RoleResult{}, nil
	}

	var results []*RoleResult
	for _, d := range facetResults[0].Data {
		item := &RoleResult{
			RoleID:    getInt32(d, "_id"),
			RoleName:  getString(d, "role_name"),
			CreatedAt: bsonMToTimePtr(d, "created_at"),
			UpdatedAt: bsonMToTimePtr(d, "updated_at"),
			DeletedAt: bsonMToTimePtr(d, "deleted_at"),
		}
		if len(facetResults[0].Metadata) > 0 {
			item.TotalCount = facetResults[0].Metadata[0].TotalCount
		}
		results = append(results, item)
	}
	return results, nil
}

// --- helpers ---

func bsonMToTimePtr(m bson.M, key string) *time.Time {
	v, ok := m[key]
	if !ok || v == nil {
		return nil
	}
	switch t := v.(type) {
	case time.Time:
		return &t
	case bson.DateTime:
		ts := t.Time()
		return &ts
	default:
		return nil
	}
}

func getInt32(m bson.M, key string) int32 {
	v, ok := m[key]
	if !ok {
		return 0
	}
	switch t := v.(type) {
	case int32:
		return t
	case int64:
		return int32(t)
	case float64:
		return int32(t)
	case int:
		return int32(t)
	default:
		return 0
	}
}

func getString(m bson.M, key string) string {
	v, ok := m[key]
	if !ok {
		return ""
	}
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}
