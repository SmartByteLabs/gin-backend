package rbac

import (
	"github.com/princeparmar/9and9-templeCMS-backend.git/pkg/database"
	"github.com/princeparmar/9and9-templeCMS-backend.git/pkg/utils"
)

type User[IDTYPE int64 | string] struct {
	database.TableID[IDTYPE]
	Username string `json:"username"`
	Password string `json:"password,omitempty"`
}

type UserRoleMapping[IDTYPE int64 | string] struct {
	database.TableID[IDTYPE]
	UserID IDTYPE `json:"user_id"`
	RoleID IDTYPE `json:"role_id"`

	// ReferenceID is used to store the id of the reference table
	// this can be nil if that model is not supporting row level access control.
	ReferenceID *IDTYPE `json:"reference_id"`
}

type Role[IDTYPE int64 | string] struct {
	database.TableID[IDTYPE]
	Name string `json:"name"`
}

type RoleAccessMapping[IDTYPE int64 | string] struct {
	database.TableID[IDTYPE]
	RoleID   IDTYPE `json:"role_id"`
	AccessID IDTYPE `json:"access_id"`

	Project database.DbSlice[string] `json:"project"`
}

type Access[IDTYPE int64 | string] struct {
	database.TableID[IDTYPE]
	Name string `json:"name"`
}

type AccessWithReferenceID[IDTYPE int64 | string] struct {
	AccessName  string
	Project     database.DbSlice[string]
	ReferenceID *IDTYPE
}

type AccessWithReferenceIDMap[IDTYPE int64 | string] struct {
	Name          string                        `json:"name"`
	Reference     map[IDTYPE]*utils.Set[string] `json:"reference"`
	GlobalProject *utils.Set[string]            `json:"global_project"`
}

func (a *AccessWithReferenceIDMap[IDTYPE]) GetAllReference() []IDTYPE {
	var result []IDTYPE
	for k := range a.Reference {
		result = append(result, k)
	}
	return result
}

func (a *AccessWithReferenceIDMap[IDTYPE]) GetAllProject(result []string) []string {
	if len(result) == 0 || result[0] == "" || result[0] == "*" {
		result = nil
	}

	if a.GlobalProject != nil {
		if result == nil {
			return a.GlobalProject.ToSlice()
		}

		return a.GlobalProject.GetCommonElements(result)
	}

	for _, v := range a.Reference {
		if result == nil {
			result = v.ToSlice()
			continue
		}

		result = v.GetCommonElements(result)
	}
	return result
}
