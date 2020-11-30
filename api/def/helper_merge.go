package def

func mergeMaps(a map[string]interface{}, b map[string]interface{}) {
	for kb, vb := range b {
		hasKey := false
		for ka, va := range a {
			if ka == kb {
				hasKey = true
				switch va.(type) {
				case map[string]interface{}:
					{
						mergeMaps(a[ka].(map[string]interface{}), b[kb].(map[string]interface{}))
						break
					}
				case []interface{}:
					{
						a[ka] = append(a[ka].([]interface{}), vb.([]interface{})...)
						break
					}
				default:
					{
						a[ka] = vb
						break
					}
				}
				break
			}
		}
		if !hasKey {
			a[kb] = vb
		}
	}
}
