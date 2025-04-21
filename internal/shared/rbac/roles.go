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

type Permission string

const (
	BanUser       Permission = "ban:user"
	UnbanUser     Permission = "unban:user"
	CreatePost    Permission = "create:post"
	DeletePost    Permission = "delete:post"
	UpdatePost    Permission = "update:post"
	DeleteUser    Permission = "delete:user"
	AwardBadge    Permission = "award:badge"
	RevokeBadge   Permission = "revoke:badge"
	MakeModerator Permission = "make:moderator"
	MakeRegular   Permission = "make:regular"
	CreateAccount Permission = "create:account"
	ViewUser      Permission = "view:user"
	ListUsers     Permission = "list:users"
)
