package repository

import (
	"context"
	"time"

	"github.com/MamangRust/microservice-ecommerce-pkg/database"
	"github.com/MamangRust/microservice-ecommerce-pkg/database/models"
	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
	"github.com/MamangRust/microservice-ecommerce-shared/errors/role_errors"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type roleCommandRepository struct {
	db *mongo.Database
}

func NewRoleCommandRepository(db *mongo.Database) *roleCommandRepository {
	return &roleCommandRepository{db: db}
}

func roleGetNextID(db *mongo.Database, collName string, ctx context.Context) (int32, error) {
	return database.GetNextID(ctx, db, collName)
}

func (r *roleCommandRepository) Create(ctx context.Context, req *requests.CreateRoleRequest) (*models.Role, error) {
	id, err := roleGetNextID(r.db, "roles", ctx)
	if err != nil {
		return nil, role_errors.ErrCreateRole.WithInternal(err)
	}
	now := time.Now()
	role := &models.Role{RoleID: id, RoleName: req.Name, CreatedAt: &now, UpdatedAt: &now}
	_, err = r.db.Collection("roles").InsertOne(ctx, role)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return nil, role_errors.ErrRoleConflict
		}
		return nil, role_errors.ErrCreateRole.WithInternal(err)
	}
	return role, nil
}

func (r *roleCommandRepository) Update(ctx context.Context, req *requests.UpdateRoleRequest) (*models.Role, error) {
	coll := r.db.Collection("roles")
	var role models.Role
	if err := coll.FindOne(ctx, bson.M{"_id": *req.ID}).Decode(&role); err != nil {
		return nil, role_errors.ErrUpdateRole.WithInternal(err)
	}
	role.RoleName = req.Name
	now := time.Now()
	role.UpdatedAt = &now
	if _, err := coll.ReplaceOne(ctx, bson.M{"_id": *req.ID}, role); err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return nil, role_errors.ErrRoleConflict
		}
		return nil, role_errors.ErrUpdateRole.WithInternal(err)
	}
	return &role, nil
}

func (r *roleCommandRepository) Trash(ctx context.Context, id int) (*models.Role, error) {
	coll := r.db.Collection("roles")
	now := time.Now()
	res, err := coll.UpdateOne(ctx, bson.M{"_id": id, "deleted_at": bson.M{"$eq": nil}}, bson.M{"$set": bson.M{"deleted_at": now}})
	if err != nil {
		return nil, role_errors.ErrTrashedRole.WithInternal(err)
	}
	if res.MatchedCount == 0 {
		return nil, role_errors.ErrRoleNotFound
	}
	var role models.Role
	if err := coll.FindOne(ctx, bson.M{"_id": id}).Decode(&role); err != nil {
		return nil, role_errors.ErrTrashedRole.WithInternal(err)
	}
	return &role, nil
}

func (r *roleCommandRepository) Restore(ctx context.Context, id int) (*models.Role, error) {
	coll := r.db.Collection("roles")
	res, err := coll.UpdateOne(ctx, bson.M{"_id": id, "deleted_at": bson.M{"$ne": nil}}, bson.M{"$set": bson.M{"deleted_at": nil}})
	if err != nil {
		return nil, role_errors.ErrRestoreRole.WithInternal(err)
	}
	if res.MatchedCount == 0 {
		return nil, role_errors.ErrRoleNotFound
	}
	var role models.Role
	if err := coll.FindOne(ctx, bson.M{"_id": id}).Decode(&role); err != nil {
		return nil, role_errors.ErrRestoreRole.WithInternal(err)
	}
	return &role, nil
}

func (r *roleCommandRepository) DeletePermanent(ctx context.Context, roleID int) (bool, error) {
	res, err := r.db.Collection("roles").DeleteOne(ctx, bson.M{"_id": roleID})
	if err != nil {
		return false, role_errors.ErrDeleteRolePermanent.WithInternal(err)
	}
	if res.DeletedCount == 0 {
		return false, role_errors.ErrRoleNotFound
	}
	return true, nil
}

func (r *roleCommandRepository) RestoreAll(ctx context.Context) (bool, error) {
	_, err := r.db.Collection("roles").UpdateMany(ctx, bson.M{"deleted_at": bson.M{"$ne": nil}}, bson.M{"$set": bson.M{"deleted_at": nil}})
	if err != nil {
		return false, role_errors.ErrRestoreAllRoles.WithInternal(err)
	}
	return true, nil
}

func (r *roleCommandRepository) DeleteAll(ctx context.Context) (bool, error) {
	_, err := r.db.Collection("roles").DeleteMany(ctx, bson.M{"deleted_at": bson.M{"$ne": nil}})
	if err != nil {
		return false, role_errors.ErrDeleteAllRoles.WithInternal(err)
	}
	return true, nil
}
