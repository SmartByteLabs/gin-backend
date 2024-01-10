package rbac

import "github.com/princeparmar/9and9-templeCMS-backend.git/pkg/database"

type User[IDTYPE int64 | string] struct {
	database.TableID[IDTYPE]
	Username string `json:"username"`
	Password string `json:"password,omitempty"`
}

type UserRoleMapping[IDTYPE int64 | string] struct {
	database.TableID[IDTYPE]
	UserID IDTYPE `json:"user_id"`
	RoleID IDTYPE `json:"role_id"`
}

type Role[IDTYPE int64 | string] struct {
	database.TableID[IDTYPE]
	Name string `json:"name"`
}

type RoleAccessMapping[IDTYPE int64 | string] struct {
	database.TableID[IDTYPE]
	RoleID   IDTYPE `json:"role_id"`
	AccessID IDTYPE `json:"access_id"`
}

type Access[IDTYPE int64 | string] struct {
	database.TableID[IDTYPE]
	Name          string                       `json:"name"`
	RoleLevelData database.DbMap[string, bool] `json:"role_level_data"`
	UserLevelData database.DbMap[string, bool] `json:"user_level_data"`
}

type RequiredData[IDTYPE int64 | string] struct {
	database.TableID[IDTYPE]
	Level    string `json:"level"`     // role or user
	ParentID IDTYPE `json:"parent_id"` // role_access_mapping_id or user_role_mapping_id
	Key      string `json:"key"`
	Value    string `json:"value"`
}

type AccessMap struct {
	Name         string            `json:"name"`
	RequiredData map[string]string `json:"required_data"`
}
