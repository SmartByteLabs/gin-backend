package rbac

import (
	"context"

	"github.com/princeparmar/gin-backend.git/pkg/database"
)

type DatabaseHelper[IDTYPE int | string] interface {
	GetUserHelper() database.CRUDDatabaseHelper[User[IDTYPE], IDTYPE]
	GetUserRoleMappingHelper() database.CRUDDatabaseHelper[UserRoleMapping[IDTYPE], IDTYPE]
	GetRoleHelper() database.CRUDDatabaseHelper[Role[IDTYPE], IDTYPE]
	GetRoleAccessMappingHelper() database.CRUDDatabaseHelper[RoleAccessMapping[IDTYPE], IDTYPE]
	GetAccessHelper() database.CRUDDatabaseHelper[Access[IDTYPE], IDTYPE]
	GetRequiredDataHelper() database.CRUDDatabaseHelper[RequiredData[IDTYPE], IDTYPE]

	GetAccessMap(ctx context.Context, UserID, accessID IDTYPE) ([]AccessMap, error)
	GetAccessMapWithCondition(ctx context.Context, UserID, accessID IDTYPE, userLevelCondition, roleLevelCondition map[string]string) ([]AccessMap, error)
}

type RBACHelper interface {
}
