package model

import "fmt"

type UserTripRole int

const (
	Owner UserTripRole = iota
	Editor
	Reader
)

var userTripRoleStrings = map[UserTripRole]string{
	Owner:  "owner",
	Editor: "editor",
	Reader: "reader",
}

func (r UserTripRole) String() string {
	if role, exists := userTripRoleStrings[r]; exists {
		return role
	}
	return ""
}

func RoleFromString(roleStr string) (UserTripRole, error) {
	switch roleStr {
	case "owner":
		return Owner, nil
	case "editor":
		return Editor, nil
	case "reader":
		return Reader, nil
	default:
		return 0, fmt.Errorf("invalid role: %s", roleStr)
	}
}
