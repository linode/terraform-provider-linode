package user

import (
	"github.com/linode/linodego"
)

func flattenGrantsEntities(entities []linodego.GrantedEntity) []interface{} {
	var result []interface{}

	for _, entity := range entities {
		entity := entity
		// Filter out entities without any permissions set.
		// This is necessary because Linode will automatically
		// create empty entities that will trigger false diffs.
		if entity.Permissions == "" {
			continue
		}

		result = append(result, flattenGrantsEntity(&entity))
	}

	return result
}

func flattenGrantsEntity(entity *linodego.GrantedEntity) map[string]interface{} {
	result := make(map[string]interface{})

	result["id"] = entity.ID
	result["permissions"] = entity.Permissions

	return result
}

func flattenGrantsGlobal(global *linodego.GlobalUserGrants) map[string]interface{} {
	result := make(map[string]interface{})

	result["account_access"] = global.AccountAccess
	result["add_domains"] = global.AddDomains
	result["add_databases"] = global.AddDatabases
	result["add_firewalls"] = global.AddFirewalls
	result["add_images"] = global.AddImages
	result["add_linodes"] = global.AddLinodes
	result["add_longview"] = global.AddLongview
	result["add_nodebalancers"] = global.AddNodeBalancers
	result["add_stackscripts"] = global.AddStackScripts
	result["add_volumes"] = global.AddVolumes
	result["cancel_account"] = global.CancelAccount
	result["longview_subscription"] = global.LongviewSubscription

	return result
}
