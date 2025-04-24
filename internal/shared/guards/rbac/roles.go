package rbac

type UserRole string

const (
	Admin     UserRole = "ADMIN"
	Regular   UserRole = "REGULAR"
	Moderator UserRole = "MODERATOR"
	Guest     UserRole = "GUEST"
)

func (r UserRole) String() string {
	return string(r)
}
