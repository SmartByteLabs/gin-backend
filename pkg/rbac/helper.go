package rbac

import (
	"context"
	"errors"

	"github.com/princeparmar/9and9-templeCMS-backend.git/pkg/database"
)

// type RbacCondition[T any, IDTYPE int64 | string] struct {
// 	database.Condition[T]
// }

// type rbacConditionHelper[T any, IDTYPE int64 | string] struct {
// 	r *RbacCondition[T, IDTYPE]
// }

// func NewRbacCondition[T any, IDTYPE int64 | string](userID IDTYPE, conditionHelper database.Condition[T]) *rbacConditionHelper[T, IDTYPE] {
// 	return &rbacConditionHelper[T, IDTYPE]{
// 		r: &RbacCondition[T, IDTYPE]{
// 			Condition: conditionHelper,
// 		},
// 	}
// }

// func (r *rbacConditionHelper[T, IDTYPE]) New() database.Condition[RbacCondition[T, IDTYPE]] {
// 	newR := *r
// 	newR.r.Condition = r.r.Condition.New()

// 	return &newR
// }

// func (r *rbacConditionHelper[T, IDTYPE]) And(condition ...database.Condition[RbacCondition[T, IDTYPE]]) database.Condition[RbacCondition[T, IDTYPE]] {
// 	for _, c := range condition {
// 		r.r.Condition = r.r.Condition.And(c.Final().Condition)
// 	}
// 	return r
// }

// func (r *rbacConditionHelper[T, IDTYPE]) Or(condition ...database.Condition[RbacCondition[T, IDTYPE]]) database.Condition[RbacCondition[T, IDTYPE]] {
// 	for _, c := range condition {
// 		r.r.Condition = r.r.Condition.Or(c.Final().Condition)
// 	}
// 	return r
// }

// func (r *rbacConditionHelper[T, IDTYPE]) Set(key string, opration database.ConditionOperation, value any) database.Condition[RbacCondition[T, IDTYPE]] {
// 	r.r.Condition = r.r.Condition.Set(key, opration, value)
// 	return r
// }

// func (r *rbacConditionHelper[T, IDTYPE]) Final() *RbacCondition[T, IDTYPE] {
// 	return r.r
// }

type DatabaseHelper[T any, IDTYPE int64 | string] interface {
	GetUserHelper() database.CrudHelper[T, User[IDTYPE], IDTYPE]
	GetUserRoleMappingHelper() database.CrudHelper[T, UserRoleMapping[IDTYPE], IDTYPE]
	GetRoleHelper() database.CrudHelper[T, Role[IDTYPE], IDTYPE]
	GetRoleAccessMappingHelper() database.CrudHelper[T, RoleAccessMapping[IDTYPE], IDTYPE]
	GetAccessHelper() database.CrudHelper[T, Access[IDTYPE], IDTYPE]

	GetAccessForUser(ctx context.Context, userID IDTYPE, access string) ([]AccessWithReferenceID[IDTYPE], error)
	GetAccessForUserWithReference(ctx context.Context, userID IDTYPE, access string, referenceID IDTYPE) ([]AccessWithReferenceID[IDTYPE], error)
}

type RbacHelper[T any, IDTYPE int64 | string] struct {
	db DatabaseHelper[T, IDTYPE]
}

func NewRbacHelper[T any, IDTYPE int64 | string](db DatabaseHelper[T, IDTYPE]) *RbacHelper[T, IDTYPE] {
	return &RbacHelper[T, IDTYPE]{
		db: db,
	}
}

func (r *RbacHelper[T, IDTYPE]) GetAccessForUser(ctx context.Context, userID IDTYPE, access string) (*AccessWithReferenceIDMap[IDTYPE], error) {
	resp, err := r.db.GetAccessForUser(ctx, userID, access)
	if err != nil {
		return nil, err
	}

	if len(resp) == 0 {
		return nil, errors.New("no access found")
	}

	return convertToAccessWithReferenceIDMap(resp), nil
}

func (r *RbacHelper[T, IDTYPE]) GetAccessForUserWithReference(ctx context.Context, userID IDTYPE, access string, referenceID IDTYPE) (*AccessWithReferenceIDMap[IDTYPE], error) {
	resp, err := r.db.GetAccessForUserWithReference(ctx, userID, access, referenceID)
	if err != nil {
		return nil, err
	}

	return convertToAccessWithReferenceIDMap(resp), nil
}
