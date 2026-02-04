package organizationuser

import "fmt"

type UserRole string

const (
	UserRoleOwner UserRole = "OWNER"
	UserRoleAdmin UserRole = "ADMIN"
	UserRoleUser  UserRole = "USER"
)

func (r UserRole) IsValid() bool {
	switch r {
	case UserRoleOwner, UserRoleAdmin, UserRoleUser:
		return true
	default:
		return false
	}
}

func (r *UserRole) Scan(src any) error {
	switch s := src.(type) {
	case string:
		*r = UserRole(s)
	case []byte:
		*r = UserRole(s)
	default:
		return fmt.Errorf("cannot scan %T into UserRole", src)
	}
	return nil
}

func (r UserRole) String() string {
	return string(r)
}
