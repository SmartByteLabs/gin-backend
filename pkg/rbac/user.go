package rbac

type AccessMap struct {
	Name         string            `json:"name"`
	RequiredData map[string]string `json:"required_data"`
}
