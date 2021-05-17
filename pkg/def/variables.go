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

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/pkg/errors"
)

// Variables defines project variables which can be defined in multiple places.
type Variables map[string]interface{}

func (v Variables) checkKey(name string) error {
	if strings.Count(name, ":") > 1 {
		return errors.WithStack(fmt.Errorf("variable name should only contain at most one colon (:)"))
	}
	return nil
}

// Get returns given value using a colon as a delimiter for sub values.
func (v Variables) Get(name string) interface{} {
	return (v)[name]
}

// GetString returns value as string.
func (v Variables) GetString(name string) string {
	return InterfaceToString(v.Get(name))
}

// GetSubMap create sub map from given prefix.
func (v Variables) GetSubMap(name string) map[string]interface{} {
	name = strings.TrimRight(name, ":") + ":"
	out := make(map[string]interface{})
	for k, v := range v {
		if strings.HasPrefix(k, name) {
			out[k[len(name):]] = v
		}
	}
	return out
}

// GetStringSubMap return sub map with string values.
func (v Variables) GetStringSubMap(name string) map[string]string {
	omap := v.GetSubMap(name)
	out := make(map[string]string)
	for k, v := range omap {
		out[k] = InterfaceToString(v)
	}
	return out
}

// Merge merges given variables with this one.
func (v *Variables) Merge(m Variables) {
	mergeMaps((*v), m)
}

// Set sets a value.
func (v *Variables) Set(key string, value interface{}) error {
	if err := v.checkKey(key); err != nil {
		return errors.WithStack(err)
	}
	(*v)[key] = value
	return nil
}

// Delete unsets a value.
func (v *Variables) Delete(key string) {
	if v.Get(key) != nil {
		delete(*v, key)
	}
}

// Keys returns list of keys.
func (v Variables) Keys() []string {
	out := make([]string, 0)
	for k := range v {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

func (v *Variables) unmarshalRawMap(data map[string]interface{}) error {
	*v = make(Variables)
	for k, iv := range data {
		switch iv := iv.(type) {
		case map[string]interface{}:
			{
				// handle old multi-level map
				for sk, sv := range iv {
					fullKey := fmt.Sprintf("%s:%s", k, sk)
					if err := v.Set(fullKey, sv); err != nil {
						return errors.WithStack(err)
					}
				}
				break
			}
		default:
			{
				if err := v.Set(k, iv); err != nil {
					return errors.WithStack(err)
				}
				break
			}
		}
	}
	return nil
}

// UnmarshalJSON implements Unmarshaler interface.
func (v *Variables) UnmarshalJSON(data []byte) error {
	umap := make(map[string]interface{})
	if err := json.Unmarshal(data, &umap); err != nil {
		return errors.WithStack(err)
	}
	return errors.WithStack(v.unmarshalRawMap(umap))
}

// UnmarshalYAML implement Unmarshaler interface.
func (v *Variables) UnmarshalYAML(unmarshal func(interface{}) error) error {
	umap := make(map[string]interface{})
	if err := unmarshal(&umap); err != nil {
		return errors.WithStack(err)
	}
	return errors.WithStack(v.unmarshalRawMap(umap))
}
