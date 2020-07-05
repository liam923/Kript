package data

import "github.com/liam923/Kript/server/pkg/proto/kript/api"

// A rawPermission gives permission for a specific operation on a datum. This is distinct from an api.Permission,
// which is a group of rawPermissions.
type rawPermission int

const (
	rawPermissionRead   = 1
	rawPermissionWrite  = 2
	rawPermissionDelete = 3
	rawPermissionShare  = 4
)

var allRawPermissions = newRawPermissionSet(
	rawPermissionRead,
	rawPermissionWrite,
	rawPermissionDelete,
	rawPermissionShare)

type rawPermissionSet struct {
	m map[rawPermission]struct{}
}

func newRawPermissionSet(permissions ...rawPermission) rawPermissionSet {
	m := map[rawPermission]struct{}{}
	for _, p := range permissions {
		m[p] = struct{}{}
	}
	return rawPermissionSet{m: m}
}

var apiPermissionToRaw = map[api.Permission]rawPermissionSet{
	api.Permission_UNKNOWN: newRawPermissionSet(),
	api.Permission_READ:    newRawPermissionSet(rawPermissionRead),
	api.Permission_WRITE:   newRawPermissionSet(rawPermissionRead, rawPermissionWrite),
	api.Permission_DELETE:  newRawPermissionSet(rawPermissionRead, rawPermissionDelete),
	api.Permission_SHARE:   newRawPermissionSet(rawPermissionRead, rawPermissionShare),
	api.Permission_ADMIN:   allRawPermissions,
}

// Confirm the user with the given id has one of the given permissions for the given datum. If the user is the owner,
// they are considered to have all permissions.
func confirmPermission(userId string, datum *datum, requiredPermissions rawPermissionSet) bool {
	if access, ok := datum.Accessors[userId]; ok {
		allRaw := newRawPermissionSet() // set of all rawPermissions the user has
		for _, apiPermission := range access.Permissions {
			for rawPermission := range apiPermissionToRaw[apiPermission].m {
				allRaw.m[rawPermission] = struct{}{}
			}
		}
		for requiredPermission := range requiredPermissions.m {
			_, contains := allRaw.m[requiredPermission]
			if !contains {
				return false
			}
		}
		return true
	} else {
		return false
	}
}

func contains(permissions []api.Permission, permission api.Permission) bool {
	for _, e := range permissions {
		if e == permission {
			return true
		}
	}
	return false
}
