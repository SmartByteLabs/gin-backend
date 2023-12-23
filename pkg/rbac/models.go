package rbac

import "github.com/princeparmar/gin-backend.git/pkg/database"

type User[IDTYPE int | string] struct {
	database.TableID[IDTYPE]
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserRoleMapping[IDTYPE int | string] struct {
	database.TableID[IDTYPE]
	UserID IDTYPE `json:"user_id"`
	RoleID IDTYPE `json:"role_id"`

	RequiredData map[string]string `json:"required_data"`
}

type Role[IDTYPE int | string] struct {
	database.TableID[IDTYPE]
	Name string `json:"name"`
}

type RoleAccessMapping[IDTYPE int | string] struct {
	database.TableID[IDTYPE]
	RoleID   IDTYPE `json:"role_id"`
	AccessID IDTYPE `json:"access_id"`

	RequiredData map[string]string `json:"required_data"`
}

type Access[IDTYPE int | string] struct {
	database.TableID[IDTYPE]
	Name                  string   `json:"name"`
	RoleLevelRequiredData []string `json:"role_level_required_data"`
	UserLevelRequiredData []string `json:"user_level_required_data"`
}

type RequiredData[IDTYPE int | string] struct {
	database.TableID[IDTYPE]
	Level    string `json:"level"`     // role or user
	ParentID IDTYPE `json:"parent_id"` // role_access_mapping_id or user_role_mapping_id
	Key      string `json:"key"`
	Value    string `json:"value"`
}
