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

	"golang.org/x/oauth2"

	"gitlab.com/contextualcode/platform_cc/v2/pkg/output"

	"github.com/pkg/errors"
)

const apiURL = "https://api.platform.sh/"

// check performs a check to ensure we're dealing with a valid platform.sh project.
func (p *Project) check() error {
	if p.ID == "" {
		return errors.Wrap(ErrProjectNotFound, "platform.sh project id not found")
	}
	return nil
}

// request performs a platform.sh API request.
func (p *Project) request(endpoint string, post map[string]interface{}, respData interface{}) error {
	// build post data
	method := "GET"
	rawPost := []byte{}
	if post != nil {
		method = "POST"
		var err error
		rawPost, err = json.Marshal(post)
		if err != nil {
			return errors.WithStack(err)
		}
	}
	// create request
	fullURL := p.apiURL + strings.TrimLeft(endpoint, "/")
	req, err := http.NewRequest(
		method,
		fullURL,
		bytes.NewReader(rawPost),
	)
	if err != nil {
		return errors.WithStack(err)
	}
	req.Header.Set("Content-Type", "application/json")
	output.LogDebug("Created Platform.sh API request.", map[string]interface{}{
		"method":    method,
		"endpoint":  endpoint,
		"url":       fullURL,
		"post_data": post,
	})
	// send request
	client, err := GetAPIClient()
	if err != nil {
		return errors.WithStack(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		return errors.WithStack(err)
	}
	defer resp.Body.Close()
	// retrieve updated token
	tok, err := client.Transport.(*oauth2.Transport).Source.Token()
	if err != nil {
		return errors.WithStack(err)
	}
	if err := saveToken(tok); err != nil {
		return errors.WithStack(err)
	}
	// process response
	rawResp, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.WithStack(err)
	}
	output.LogDebug("Recieved Platform.sh API response.", string(rawResp))
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		if respData != nil {
			json.Unmarshal(rawResp, respData)
		}
		return errors.Wrapf(ErrBadAPIResponse, "platform.sh api returned status code %d", resp.StatusCode)
	}
	if respData != nil {
		return errors.WithStack(json.Unmarshal(rawResp, respData))
	}
	return nil
}

// FetchEnvironments populates environments list.
func (p *Project) FetchEnvironments() error {
	if len(p.Environments) > 0 {
		return nil
	}
	if err := p.check(); err != nil {
		return errors.WithStack(err)
	}
	if err := p.request("/projects/"+p.ID+"/environments", nil, &p.Environments); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// Fetch populates project with API data.
func (p *Project) Fetch() error {
	if err := p.check(); err != nil {
		return errors.WithStack(err)
	}
	if err := p.FetchEnvironments(); err != nil {
		return errors.WithStack(err)
	}
	return nil
}