package rbac

import (
	"context"
	"errors"

	"github.com/princeparmar/9and9-templeCMS-backend.git/pkg/database"
)

type RbacCondition[T any, IDTYPE int64 | string] struct {
	UserID IDTYPE

	UserLevelData map[string]string
	RoleLevelData map[string]string

	database.Condition[T]
}

type rbacConditionHelper[T any, IDTYPE int64 | string] struct {
	r *RbacCondition[T, IDTYPE]
}

func NewRbacCondition[T any, IDTYPE int64 | string](userID IDTYPE, conditionHelper database.Condition[T]) *rbacConditionHelper[T, IDTYPE] {
	return &rbacConditionHelper[T, IDTYPE]{
		r: &RbacCondition[T, IDTYPE]{
			UserID:        userID,
			UserLevelData: make(map[string]string),
			RoleLevelData: make(map[string]string),
			Condition:     conditionHelper,
		},
	}
}

func (r *rbacConditionHelper[T, IDTYPE]) SetUserLevelData(key string, value string) database.Condition[RbacCondition[T, IDTYPE]] {
	r.r.UserLevelData[key] = value
	return r
}

func (r *rbacConditionHelper[T, IDTYPE]) SetRoleLevelData(key string, value string) database.Condition[RbacCondition[T, IDTYPE]] {
	r.r.RoleLevelData[key] = value
	return r
}

func (r *rbacConditionHelper[T, IDTYPE]) New() database.Condition[RbacCondition[T, IDTYPE]] {
	newR := *r
	newR.r.Condition = r.r.Condition.New()

	return &newR
}

func (r *rbacConditionHelper[T, IDTYPE]) And(condition ...database.Condition[RbacCondition[T, IDTYPE]]) database.Condition[RbacCondition[T, IDTYPE]] {
	for _, c := range condition {
		r.r.Condition = r.r.Condition.And(c.Final().Condition)
	}
	return r
}

func (r *rbacConditionHelper[T, IDTYPE]) Or(condition ...database.Condition[RbacCondition[T, IDTYPE]]) database.Condition[RbacCondition[T, IDTYPE]] {
	for _, c := range condition {
		r.r.Condition = r.r.Condition.Or(c.Final().Condition)
	}
	return r
}

func (r *rbacConditionHelper[T, IDTYPE]) Set(key string, opration database.ConditionOperation, value any) database.Condition[RbacCondition[T, IDTYPE]] {
	r.r.Condition = r.r.Condition.Set(key, opration, value)
	return r
}

func (r *rbacConditionHelper[T, IDTYPE]) Final() *RbacCondition[T, IDTYPE] {
	return r.r
}

type DatabaseHelper[T any, IDTYPE int64 | string] interface {
	GetUserHelper() database.CrudHelper[T, User[IDTYPE], IDTYPE]
	GetUserRoleMappingHelper() database.CrudHelper[T, UserRoleMapping[IDTYPE], IDTYPE]
	GetRoleHelper() database.CrudHelper[T, Role[IDTYPE], IDTYPE]
	GetRoleAccessMappingHelper() database.CrudHelper[T, RoleAccessMapping[IDTYPE], IDTYPE]
	GetAccessHelper() database.CrudHelper[T, Access[IDTYPE], IDTYPE]
	GetRequiredDataHelper() database.CrudHelper[T, RequiredData[IDTYPE], IDTYPE]

	GetAccessMap(ctx context.Context, userID IDTYPE, access string) ([]AccessMap, error)
	GetAccessMapWithCondition(ctx context.Context, userID IDTYPE, access string, userLevelCondition, roleLevelCondition map[string]string) ([]AccessMap, error)
}

type RbacHelper[T any, IDTYPE int64 | string] struct {
	db DatabaseHelper[T, IDTYPE]
}

func NewRbacHelper[T any, IDTYPE int64 | string](db DatabaseHelper[T, IDTYPE]) *RbacHelper[T, IDTYPE] {
	return &RbacHelper[T, IDTYPE]{
		db: db,
	}
}

func (r *RbacHelper[T, IDTYPE]) GetAccessMapWithCondition(ctx context.Context, access string, condition database.Condition[RbacCondition[T, IDTYPE]]) ([]AccessMap, error) {
	rData := condition.Final()

	accessAr, err := r.db.GetAccessHelper().Get(ctx, database.AllFields, rData.Condition.New().Set("name", database.ConditionOperationEqual, access))
	if err != nil {
		return nil, err
	}

	if len(accessAr) == 0 {
		return nil, errors.New("Access not found")
	}

	for field, required := range accessAr[0].RoleLevelData.Map() {
		if required {
			if _, ok := rData.RoleLevelData[field]; !ok {
				return nil, errors.New("Role level required data not found")
			}
		}
	}

	for field, required := range accessAr[0].UserLevelData.Map() {
		if required {
			if _, ok := rData.UserLevelData[field]; !ok {
				return nil, errors.New("User level required data not found")
			}
		}
	}

	return r.db.GetAccessMapWithCondition(ctx, rData.UserID, access, rData.UserLevelData, rData.RoleLevelData)
}
