/*
This file is part of Platform.CC.

Platform.CC is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

Platform.CC is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with Platform.CC.  If not, see <https://www.gnu.org/licenses/>.
*/

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
