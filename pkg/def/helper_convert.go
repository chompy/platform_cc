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

import "fmt"

// InterfaceToString returns conversion of given interface to string.
func InterfaceToString(v interface{}) string {
	switch v := v.(type) {
	case string:
		{
			return v
		}
	case int:
		{
			return fmt.Sprintf("%d", v)
		}
	case bool:
		{
			if v {
				return "true"
			}
			return ""
		}
	case float32:
		{
			return fmt.Sprintf("%f", v)
		}
	case float64:
		{
			return fmt.Sprintf("%f", v)
		}
	default:
		{
			return ""
		}
	}
}
