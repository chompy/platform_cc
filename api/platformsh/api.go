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
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"gitlab.com/contextualcode/platform_cc/api/output"

	"github.com/ztrue/tracerr"
)

const platformshAccessTokenUrl = "https://accounts.platform.sh/oauth2/token"
const platformshApiUrl = "https://api.platform.sh/"
const apiTokenPath = "~/.pcc/psh_api_token.txt"
const apiOauthPort = 31601
const apiOauthClientId = "platform-cli"
const apiOauth = "https://auth.api.platform.sh/oauth2/authorize?redirect_uri=http%%3A%%2F%%2F127.0.0.1%%3A%d&amp;state=%s&amp;client_id=%s&amp;prompt=consent%%20select_account&amp;response_type=code&amp;code_challenge=%s&amp;code_challenge_method=S256&amp;scope=offline_access"

// GetAccessToken returns the currently stored API token.
func GetAccessToken() string {
	data, err := ioutil.ReadFile(expandPath(apiTokenPath))
	if err != nil {
		return ""
	}
	return string(data)
}

// SetAccessToken sets the Platform.sh API token.
func SetAccessToken(value string) error {
	return tracerr.Wrap(ioutil.WriteFile(
		expandPath(apiTokenPath), []byte(value), 0600,
	))
}

// APILogin creates a temporary web server for handling Platform.sh OAuth login.
func APILogin() error {
	b64Url := func(data []byte) string {
		strVal := base64.StdEncoding.EncodeToString(data)
		strVal = strings.ReplaceAll(strVal, "+", "-")
		strVal = strings.ReplaceAll(strVal, "/", "_")
		strVal = strings.TrimRight(strVal, "=")
		return strVal
	}
	// generate params
	stateRand := make([]byte, 32)
	if _, err := rand.Read(stateRand); err != nil {
		return tracerr.Wrap(err)
	}
	state := b64Url(stateRand)
	codeRand := make([]byte, 32)
	if _, err := rand.Read(codeRand); err != nil {
		return tracerr.Wrap(err)
	}
	codeVerify := b64Url(codeRand)
	h := sha256.New()
	h.Write([]byte(codeVerify))
	codeChallenge := b64Url(h.Sum(nil))

	// generate oauth redirect
	oauthRedirect := fmt.Sprintf(
		apiOauth,
		apiOauthPort,
		state,
		apiOauthClientId,
		codeChallenge,
	)

	// output instructions
	output.Info(fmt.Sprintf("Open http://127.0.0.1:%d in your browser to continue.", apiOauthPort))

	// start server
	httpMux := http.NewServeMux()
	httpServer := http.Server{Addr: fmt.Sprintf("127.0.0.1:%d", apiOauthPort), Handler: httpMux}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	httpMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		renderText := func(text string, status int) {
			w.Header().Set("Content-Type", "text/plain;charset=utf-8")
			w.WriteHeader(status)
			fmt.Fprint(w, text)
		}

		returnCode := r.URL.Query().Get("code")
		//returnState := r.URL.Query().Get("state")
		returnError := r.URL.Query().Get("error")
		returnErrorDesc := r.URL.Query().Get("error_description")
		if returnError != "" || returnErrorDesc != "" || returnCode != "" {
			defer cancel()
			if returnError != "" || returnErrorDesc != "" {
				// error returned
				renderText("An error has occured, check your terminal for more information.", http.StatusInternalServerError)
				output.Warn(fmt.Sprintf("OAuth Error, %s, %s", returnError, returnErrorDesc))
				return
			} else if returnCode != "" {
				// fetch access token
				accessToken, err := fetchAccessToken(
					returnCode, codeVerify,
				)
				if err != nil {
					renderText("An error has occured, check your terminal for more information.", http.StatusInternalServerError)
					output.Error(err)
					return
				}
				if err := SetAccessToken(accessToken); err != nil {
					renderText("An error has occured, check your terminal for more information.", http.StatusInternalServerError)
					output.Error(err)
					return
				}
				renderText("Success! You may close this window.", http.StatusOK)
				output.Info("Success!")
				return
			}
		}
		// perform oauth redirect
		http.Redirect(w, r, oauthRedirect, http.StatusTemporaryRedirect)
	})
	go func() {
		if err := httpServer.ListenAndServe(); err != nil {
			if err != nil && err != http.ErrServerClosed {
				output.Error(tracerr.Wrap(err))
			}
		}
	}()
	<-ctx.Done()
	httpServer.Shutdown(ctx)
	return nil
}

// fetchAccessToken fetches the access token from the Platform.sh API.
func fetchAccessToken(returnCode string, codeVerify string) (string, error) {
	data := map[string]string{
		"grant_type":    "authorization_code",
		"code":          returnCode,
		"client_id":     apiOauthClientId,
		"redirect_uri":  fmt.Sprintf("http://127.0.0.1:%d", apiOauthPort),
		"code_verifier": codeVerify,
	}
	output.LogDebug("Send request for access token to Platform.sh API.", data)
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", tracerr.Wrap(err)
	}
	resp, err := http.Post(
		platformshAccessTokenUrl,
		"application/json",
		bytes.NewReader(jsonData),
	)
	if err != nil {
		return "", tracerr.Wrap(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		r, _ := ioutil.ReadAll(resp.Body)
		output.LogDebug("Invalid response when fetching Platform.sh access token.", string(r))
		return "", tracerr.Errorf("status %s when trying to fetch access token", resp.Status)
	}
	respDataRaw, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", tracerr.Wrap(err)
	}
	output.LogDebug("Recieved access token response from Platform.sh API.", string(respDataRaw))
	respData := map[string]interface{}{}
	if err := json.Unmarshal(respDataRaw, &respData); err != nil {
		return "", tracerr.Wrap(err)
	}
	if respData["access_token"] != nil {
		return respData["access_token"].(string), nil
	}
	return "", tracerr.Errorf("invalid api response")
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
	accessToken := GetAccessToken()
	if accessToken == "" {
		return tracerr.Errorf("access token not set")
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
	req.Header.Set("Authorization", "Bearer "+accessToken)
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
	output.LogDebug("Recieved Platform.sh API response.", string(rawResp))
	if resp.StatusCode != 200 {
		return tracerr.Errorf("platform.sh api returned status code %d", resp.StatusCode)
	}
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
