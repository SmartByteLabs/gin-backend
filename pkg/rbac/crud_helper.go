package rbac

import (
	"context"
	"errors"

	"github.com/princeparmar/9and9-templeCMS-backend.git/pkg/database"
	"github.com/princeparmar/9and9-templeCMS-backend.git/pkg/utils"
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
	helper   *RbacHelper[T, IDTYPE]
	dbHelper database.CrudHelper[T, MODEL, IDTYPE]
}

func NewCrudHelper[T any, MODEL database.TableWithID[IDTYPE], IDTYPE int64 | string](helper *RbacHelper[T, IDTYPE], dbHelper database.CrudHelper[T, MODEL, IDTYPE]) database.CrudHelper[RbacCondition[T, IDTYPE], MODEL, IDTYPE] {
	return &CrudHelper[T, MODEL, IDTYPE]{
		helper:   helper,
		dbHelper: dbHelper,
	}
}

func (rbacCrudHelper *CrudHelper[T, MODEL, IDTYPE]) GetColumns(project []string, withoutID bool) []string {
	return rbacCrudHelper.dbHelper.GetColumns(project, withoutID)
}

func (rbacCrudHelper *CrudHelper[T, MODEL, IDTYPE]) GetColumnsWithAccess(project []string, access []AccessMap) []string {
	return utils.NewSetFromSlice(getLevelFromAccessMap(access)).Intersection(utils.NewSetFromSlice(rbacCrudHelper.GetColumns(project, false))).ToSlice()
}

func (rbacCrudHelper *CrudHelper[T, MODEL, IDTYPE]) GetTableName() string {
	return rbacCrudHelper.dbHelper.GetTableName()
}

func (rbacCrudHelper *CrudHelper[T, MODEL, IDTYPE]) GetAccessForMethod(ctx context.Context, conditionHelper database.Condition[RbacCondition[T, IDTYPE]], method string) ([]AccessMap, error) {

	access, err := rbacCrudHelper.helper.GetAccessMapWithCondition(ctx, rbacCrudHelper.dbHelper.GetTableName()+"_"+method, conditionHelper)
	if err != nil {
		return nil, err
	}

	if len(access) == 0 {
		return nil, errors.New("Access not found")
	}

	return access, nil
}

func (rbacCrudHelper *CrudHelper[T, MODEL, IDTYPE]) Create(ctx context.Context, body *MODEL, conditionHelper database.Condition[RbacCondition[T, IDTYPE]]) (*MODEL, error) {
	_, err := rbacCrudHelper.GetAccessForMethod(ctx, conditionHelper, "CREATE")
	if err != nil {
		return nil, err
	}

	return rbacCrudHelper.dbHelper.Create(ctx, body, conditionHelper.Final().Condition)
}

func (rbacCrudHelper *CrudHelper[T, MODEL, IDTYPE]) Get(ctx context.Context, project []string, conditionHelper database.Condition[RbacCondition[T, IDTYPE]]) ([]MODEL, error) {

	access, err := rbacCrudHelper.GetAccessForMethod(ctx, conditionHelper, "GET")
	if err != nil {
		return nil, err
	}

	data, err := rbacCrudHelper.dbHelper.Get(ctx, rbacCrudHelper.GetColumnsWithAccess(project, access), conditionHelper.Final().Condition)
	if err != nil {
		return nil, err
	}

	if len(data) == 0 {
		return nil, errors.New("No data found")
	}

	return data, nil
}

func (rbacCrudHelper *CrudHelper[T, MODEL, IDTYPE]) Update(ctx context.Context, m *MODEL, project []string, condition database.Condition[RbacCondition[T, IDTYPE]]) error {

	access, err := rbacCrudHelper.GetAccessForMethod(ctx, condition, "UPDATE")
	if err != nil {
		return err
	}

	return rbacCrudHelper.dbHelper.Update(ctx, m, rbacCrudHelper.GetColumnsWithAccess(project, access), condition.Final().Condition)
}

func (rbacCrudHelper *CrudHelper[T, MODEL, IDTYPE]) Delete(ctx context.Context, condition database.Condition[RbacCondition[T, IDTYPE]]) error {

	_, err := rbacCrudHelper.GetAccessForMethod(ctx, condition, "DELETE")
	if err != nil {
		return err
	}

	return rbacCrudHelper.dbHelper.Delete(ctx, condition.Final().Condition)
}
