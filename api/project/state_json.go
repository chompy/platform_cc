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

package project

import (
	"encoding/json"

	"gitlab.com/contextualcode/platform_cc/api/def"
)

type serviceState struct {
	Running               bool                   `json:"running"`
	Open                  bool                   `json:"open"`
	Frozen                bool                   `json:"frozen"`
	Image                 string                 `json:"image"`
	Resources             map[string]interface{} `json:"resources"`
	Relationships         map[string]interface{} `json:"relationships"`
	ResolvedRelationships map[string]interface{} `json:"resolved_relationships"`
	Application           *def.App               `json:"application"`
	Configuration         map[string]interface{} `json:"configuration"`
	Endpoints             map[string]interface{} `json:"endpoints"`
	PrimaryIP             string                 `json:"primary_ip"`
	Internal              string                 `json:"internal"`
	Slug                  string                 `json:"slug"`
	ImageInfo             map[string]interface{} `json:"image_info"`
	SlugInfo              map[string]interface{} `json:"slug_info"`
}

func getDefaultServiceState() serviceState {
	return serviceState{
		Running:               true,
		Open:                  false,
		Frozen:                false,
		Image:                 "",
		Resources:             nil,
		Relationships:         map[string]interface{}{},
		ResolvedRelationships: map[string]interface{}{},
		Application:           nil,
		Configuration:         map[string]interface{}{},
		Endpoints:             map[string]interface{}{},
		PrimaryIP:             "127.0.0.1",
		Internal:              "",
		Slug:                  "",
		ImageInfo:             map[string]interface{}{},
		SlugInfo:              map[string]interface{}{},
	}
}

// buildStateJSON builds state change JSON.
func buildStateJSON(instanceID string, current serviceState, desired serviceState) ([]byte, error) {
	out := map[string]interface{}{
		"current_state": map[string]interface{}{
			"instances": map[string]interface{}{
				instanceID: map[string]interface{}{
					"change_id": "",
					"order_key": nil,
					"state":     current,
					"stack":     []string{},
				},
			},
		},
		"desired_state": map[string]interface{}{
			"instances": map[string]interface{}{
				instanceID: map[string]interface{}{
					"change_id": "",
					"order_key": nil,
					"state":     desired,
					"stack":     []string{},
				},
			},
		},
	}
	return json.Marshal(out)
}

/*{
    "current_state" : {
        "instances" : {
            "b17392587b86" : {
                "change_id": "",
                "order_key": null,
                "state": {
                    "running": true,
                    "open": false,
                    "frozen": false,
                    "image": "",
                    "resources": null,
                    "relationships": {},
                    "resolved_relationships": {},
                    "application": null,
                    "configuration": {
                        "authentication": {
                            "enabled": false
                        }
                    },
                    "endpoints": {},
                    "primary_ip": "127.0.0.1",
                    "internal": "",
                    "slug": null,
                    "image_info": {},
                    "slug_info": {}
                },
                "stack": []
            }
        }
    },
    "desired_state": {
        "instances" : {
            "b17392587b86" : {
                "change_id": "",
                "order_key": null,
                "state": {
                    "running": true,
                    "open": false,
                    "frozen": false,
                    "image": "",
                    "resources": null,
                    "relationships": {},
                    "resolved_relationships": {},
                    "application": null,
                    "configuration": {
                        "authentication": {
                            "enabled": true
                        }
                    },
                    "endpoints": {},
                    "primary_ip": "127.0.0.1",
                    "internal": "",
                    "slug": null,
                    "image_info": {},
                    "slug_info": {}
                },
                "stack": []
            }
        }
    }
}*/
