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

package platformsh

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/ztrue/tracerr"
	"gitlab.com/contextualcode/platform_cc/api/def"
)

const platformshAccessTokenUrl = "https://accounts.platform.sh/oauth2/token"
const platformshApiUrl = "https://api.platform.sh/"

// API handles connection to PlatformSH API.
type API struct {
	APIToken    string
	AccessToken string
}

// NewAPI returns new API.
func NewAPI(globalConfig *def.GlobalConfig) API {
	apiToken := ""
	if globalConfig != nil {
		apiToken = globalConfig.PlatformSH.APIToken
	}
	return API{
		APIToken: apiToken,
	}
}

func (p *API) getAccessToken() error {
	if p.AccessToken != "" {
		return nil
	}
	data := map[string]string{
		"client_id":  "platform-api-user",
		"grant_type": "api_token",
		"api_token":  p.APIToken,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return tracerr.Wrap(err)
	}
	resp, err := http.Post(
		platformshAccessTokenUrl,
		"application/json",
		bytes.NewReader(jsonData),
	)
	if err != nil {
		return tracerr.Wrap(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return tracerr.Errorf("status %s when trying to fetch access token", resp.Status)
	}
	respDataRaw, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return tracerr.Wrap(err)
	}
	respData := map[string]interface{}{}
	if err := json.Unmarshal(respDataRaw, &respData); err != nil {
		return tracerr.Wrap(err)
	}
	if respData["access_token"] != nil {
		p.AccessToken = respData["access_token"].(string)
	}
	return nil
}

func (p *API) request(endpoint string, post map[string]interface{}, respData interface{}) error {

	if err := p.getAccessToken(); err != nil {
		return tracerr.Wrap(err)
	}

	// build post data
	rawPost := []byte{}
	if post != nil {
		var err error
		rawPost, err = json.Marshal(post)
		if err != nil {
			return tracerr.Wrap(err)
		}
	}

	// create request
	req, err := http.NewRequest(
		"GET",
		platformshApiUrl+strings.TrimLeft(endpoint, "/"),
		bytes.NewReader(rawPost),
	)
	if err != nil {
		return tracerr.Wrap(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+p.AccessToken)

	// send request
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return tracerr.Wrap(err)
	}

	// process response
	defer resp.Body.Close()
	rawResp, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return tracerr.Wrap(err)
	}
	err = json.Unmarshal(rawResp, respData)
	return tracerr.Wrap(err)
}

// PopulateProjectEnvironments populates environment list for given Platform.sh project.
func (p *API) PopulateProjectEnvironments(project *Project) error {
	if err := p.request("/projects/"+project.ID+"/environments", nil, &project.Environments); err != nil {
		return tracerr.Wrap(err)
	}
	return nil
}

// PopulateProject populates given Platform.sh project with data from API.
func (p *API) PopulateProject(project *Project) error {

	if err := p.PopulateProjectEnvironments(project); err != nil {
		return tracerr.Wrap(err)
	}

	return nil
}
