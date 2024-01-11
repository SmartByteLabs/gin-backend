package rbac

import (
	"context"
	"errors"

	"github.com/princeparmar/9and9-templeCMS-backend.git/pkg/database"
)

// type ConditionalData[IDTYPE int64 | string] struct {
// 	UserID IDTYPE

// 	Where string
// 	Args  []any

// 	UserLevelData map[string]string
// 	RoleLevelData map[string]string
// }

// func (r *ConditionalData[IDTYPE]) GetWhere(where string, args ...any) (string, []any) {
// 	if r == nil || len(r.Where) == 0 {
// 		return where, args
// 	}

// 	if len(where) == 0 {
// 		return r.Where, r.Args
// 	}

// 	return "(" + where + ") AND (" + r.Where + ")", append(args, r.Args...)
// }

type CrudHelper[T any, MODEL database.TableWithID[IDTYPE], IDTYPE int64 | string] struct {
	helper            *RbacHelper[T, IDTYPE]
	dbHelper          database.CrudHelper[T, MODEL, IDTYPE]
	getUserfunc       func(ctx context.Context) (*User[IDTYPE], error)
	referenceRequired bool
}

func NewCrudHelper[T any, MODEL database.TableWithID[IDTYPE], IDTYPE int64 | string](helper *RbacHelper[T, IDTYPE],
	dbHelper database.CrudHelper[T, MODEL, IDTYPE], getUserfunc func(ctx context.Context) (*User[IDTYPE], error),
) *CrudHelper[T, MODEL, IDTYPE] {
	return &CrudHelper[T, MODEL, IDTYPE]{
		helper:      helper,
		dbHelper:    dbHelper,
		getUserfunc: getUserfunc,
	}
}

func (rbacCrudHelper *CrudHelper[T, MODEL, IDTYPE]) ReferenceRequired() *CrudHelper[T, MODEL, IDTYPE] {
	rbacCrudHelper.referenceRequired = true
	return rbacCrudHelper
}

func (rbacCrudHelper *CrudHelper[T, MODEL, IDTYPE]) GetColumns(project []string, withoutID bool) []string {
	return rbacCrudHelper.dbHelper.GetColumns(project, withoutID)
}

func (rbacCrudHelper *CrudHelper[T, MODEL, IDTYPE]) GetTableName() string {
	return rbacCrudHelper.dbHelper.GetTableName()
}

func (rbacCrudHelper *CrudHelper[T, MODEL, IDTYPE]) checkAccess(ctx context.Context, method string) (*AccessWithReferenceIDMap[IDTYPE], error) {
	user, err := rbacCrudHelper.getUserfunc(ctx)
	if err != nil {
		return nil, err
	}

	return rbacCrudHelper.helper.GetAccessForUser(ctx, user.GetID(), rbacCrudHelper.dbHelper.GetTableName()+"_"+method)
}

func (rbacCrudHelper *CrudHelper[T, MODEL, IDTYPE]) checkAccessWithConditionUpdate(ctx context.Context, method string,
	conditionHelper database.Condition[T]) (*AccessWithReferenceIDMap[IDTYPE], error) {
	access, err := rbacCrudHelper.checkAccess(ctx, method)
	if err != nil {
		return nil, err
	}

	if !rbacCrudHelper.referenceRequired {
		return access, nil
	}

	allReference := access.GetAllReference()
	if len(allReference) == 0 {
		return nil, errors.New("no data found")
	}

	conditionHelperForReference := conditionHelper.New()
	for _, id := range allReference {
		conditionHelperForReference.Or(conditionHelper.New().Set("id", database.ConditionOperationEqual, id))
	}
	conditionHelper.And(conditionHelperForReference)

	return access, nil
}

func (rbacCrudHelper *CrudHelper[T, MODEL, IDTYPE]) Create(ctx context.Context, body *MODEL) (*MODEL, error) {
	_, err := rbacCrudHelper.checkAccess(ctx, "CREATE")
	if err != nil {
		return nil, err
	}

	return rbacCrudHelper.dbHelper.Create(ctx, body)
}

func (rbacCrudHelper *CrudHelper[T, MODEL, IDTYPE]) Get(ctx context.Context, project []string, conditionHelper database.Condition[T]) ([]MODEL, error) {

	access, err := rbacCrudHelper.checkAccessWithConditionUpdate(ctx, "GET", conditionHelper)
	if err != nil {
		return nil, err
	}

	data, err := rbacCrudHelper.dbHelper.Get(ctx, access.GetAllProject(project), conditionHelper)
	if err != nil {
		return nil, err
	}

	if len(data) == 0 {
		return nil, errors.New("no data found")
	}

	return data, nil
}

// ISSUE: when we pass empty project and after doing union with access project is nil then it should not return error but it will result all columns.
func (rbacCrudHelper *CrudHelper[T, MODEL, IDTYPE]) Update(ctx context.Context, m *MODEL, project []string, conditionHelper database.Condition[T]) error {

	access, err := rbacCrudHelper.checkAccessWithConditionUpdate(ctx, "UPDATE", conditionHelper)
	if err != nil {
		return err
	}

	return rbacCrudHelper.dbHelper.Update(ctx, m, access.GetAllProject(project), conditionHelper)
}

func (rbacCrudHelper *CrudHelper[T, MODEL, IDTYPE]) Delete(ctx context.Context, conditionHelper database.Condition[T]) error {

	_, err := rbacCrudHelper.checkAccessWithConditionUpdate(ctx, "DELETE", conditionHelper)
	if err != nil {
		return err
	}

	return rbacCrudHelper.dbHelper.Delete(ctx, conditionHelper)
}
