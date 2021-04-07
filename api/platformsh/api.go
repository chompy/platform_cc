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

	"gitlab.com/contextualcode/platform_cc/api/output"

	"github.com/ztrue/tracerr"
)

const platformshAccessTokenUrl = "https://accounts.platform.sh/oauth2/token"
const platformshApiUrl = "https://api.platform.sh/"

// SetAPIToken sets the platform.sh API token.
func (p *Project) SetAPIToken(value string) {
	p.apiToken = value
}

// getAccessToken fetches the access token from the platform.sh.
func (p *Project) getAccessToken() error {
	if p.apiAccessToken != "" {
		return nil
	}
	if p.apiToken == "" {
		return tracerr.Errorf("platform.sh api token not set")
	}
	data := map[string]string{
		"client_id":  "platform-api-user",
		"grant_type": "api_token",
		"api_token":  p.apiToken,
	}
	output.LogDebug("Send request for access token to Platform.sh API.", data)
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
	output.LogDebug("Recieved access token response from Platform.sh API.", respDataRaw)
	respData := map[string]interface{}{}
	if err := json.Unmarshal(respDataRaw, &respData); err != nil {
		return tracerr.Wrap(err)
	}
	if respData["access_token"] != nil {
		p.apiAccessToken = respData["access_token"].(string)
	}
	return nil
}

// check performs a check to ensure we're dealing with a valid platform.sh project.
func (p *Project) check() error {
	if p.ID == "" {
		return tracerr.Errorf("platform.sh project id not found")
	}
	return nil
}

// request performs a platform.sh API request.
func (p *Project) request(endpoint string, post map[string]interface{}, respData interface{}) error {
	if err := p.getAccessToken(); err != nil {
		return tracerr.Wrap(err)
	}
	// build post data
	method := "GET"
	rawPost := []byte{}
	if post != nil {
		method = "POST"
		var err error
		rawPost, err = json.Marshal(post)
		if err != nil {
			return tracerr.Wrap(err)
		}
	}
	// create request
	fullURL := platformshApiUrl + strings.TrimLeft(endpoint, "/")
	req, err := http.NewRequest(
		method,
		fullURL,
		bytes.NewReader(rawPost),
	)
	if err != nil {
		return tracerr.Wrap(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+p.apiAccessToken)
	output.LogDebug("Created Platform.sh API request.", map[string]interface{}{
		"method":    method,
		"endpoint":  endpoint,
		"url":       fullURL,
		"post_data": post,
	})
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
	output.LogDebug("Recieved Platform.sh API response.", rawResp)
	if respData != nil {
		return tracerr.Wrap(json.Unmarshal(rawResp, respData))
	}
	return nil
}

// FetchEnvironments populates environments list.
func (p *Project) FetchEnvironments() error {
	if len(p.Environments) > 0 {
		return nil
	}
	if err := p.check(); err != nil {
		return tracerr.Wrap(err)
	}
	if err := p.request("/projects/"+p.ID+"/environments", nil, &p.Environments); err != nil {
		return tracerr.Wrap(err)
	}
	return nil
}

// Fetch populates project with API data.
func (p *Project) Fetch() error {
	if err := p.check(); err != nil {
		return tracerr.Wrap(err)
	}
	if err := p.FetchEnvironments(); err != nil {
		return tracerr.Wrap(err)
	}
	return nil
}
