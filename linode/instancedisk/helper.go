package instancedisk

func expandStackScriptData(data any) map[string]string {
	dataMap := data.(map[string]any)
	result := make(map[string]string, len(dataMap))

	for k, v := range dataMap {
		result[k] = v.(string)
	}

	return result
}
