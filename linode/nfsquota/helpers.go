package nfsquota

import "reflect"

// userGroupConfigEqual reports whether two user_group_config models are
// semantically equal. It is used to determine whether an update is required.
func userGroupConfigEqual(a, b *UserGroupConfigModel) bool {
	return reflect.DeepEqual(a, b)
}
