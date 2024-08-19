package models

import "time"

type AdminRole int

const (
	ROLE_STAFF AdminRole = iota
	ROLE_SUPER_ADMIN
	EDITOR
)

var AdminRoleMappings = map[AdminRole]string{
	ROLE_STAFF:       "admin",
	ROLE_SUPER_ADMIN: "super_admin",
	EDITOR:           "editor",
	// ... add more roles here
}

type Admin struct {
	Id           string
	Email        string
	Name         string
	PasswordHash []byte
	CreatedAt    time.Time
	UpdatedAt    time.Time
	RoleMask     int
}

func (a Admin) GetRoles() []string {
	adminRoles := make([]string, 0)

	for role, roleName := range AdminRoleMappings {
		if a.HasRole(role) {
			adminRoles = append(adminRoles, roleName)
			break
		}
	}

	return adminRoles
}

func (a Admin) HasRole(role AdminRole) bool {
	// return int(role)&a.RoleMask == int(role)
	return a.RoleMask == int(role)
}
