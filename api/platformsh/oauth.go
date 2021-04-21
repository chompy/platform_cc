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
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"gitlab.com/contextualcode/platform_cc/api/config"

	"github.com/ztrue/tracerr"
	"gitlab.com/contextualcode/platform_cc/api/output"
	"golang.org/x/oauth2"
)

const oauthTokenStore = "psh_api_token"
const oauthPort = 31698
const oauthClientId = "platform-cli"
const oauthAuthURL = "https://auth.api.platform.sh/oauth2/authorize"
const oauthTokenURL = "https://accounts.platform.sh/oauth2/token"

func getOauthConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID: oauthClientId,
		Scopes:   []string{"offline_access"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  oauthAuthURL,
			TokenURL: oauthTokenURL,
		},
		RedirectURL: fmt.Sprintf("http://127.0.0.1:%d", oauthPort),
	}
}

func base64URL(data []byte) string {
	strVal := base64.StdEncoding.EncodeToString(data)
	strVal = strings.ReplaceAll(strVal, "+", "-")
	strVal = strings.ReplaceAll(strVal, "/", "_")
	strVal = strings.TrimRight(strVal, "=")
	return strVal
}

func generateRandom() string {
	randBytes := make([]byte, 32)
	if _, err := rand.Read(randBytes); err != nil {
		return ""
	}
	return base64URL(randBytes)
}

func saveToken(t *oauth2.Token) error {
	out, err := json.Marshal(t)
	if err != nil {
		return tracerr.Wrap(err)
	}
	return tracerr.Wrap(ioutil.WriteFile(
		filepath.Join(config.Path(), oauthTokenStore),
		out, 0644,
	))
}

func loadToken() (*oauth2.Token, error) {
	t := &oauth2.Token{}
	rawData, err := ioutil.ReadFile(filepath.Join(config.Path(), oauthTokenStore))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, tracerr.Errorf("platform.sh api token not found, please use the platformsh:login command to generate it")
		}
		return nil, tracerr.Wrap(err)
	}
	if err := json.Unmarshal(rawData, t); err != nil {
		return nil, tracerr.Wrap(err)
	}
	return t, nil
}

func loginServer(authURL string) (string, error) {
	var err error
	code := ""
	m := http.NewServeMux()
	s := http.Server{Addr: fmt.Sprintf("127.0.0.1:%d", oauthPort), Handler: m}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	m.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		errType := r.URL.Query().Get("error")
		errDesc := r.URL.Query().Get("error_description")
		code = r.URL.Query().Get("code")
		if errType != "" {
			// error detected, display to user and close server
			w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
			w.WriteHeader(http.StatusInternalServerError)
			io.WriteString(w, "An error has occured. See terminal for more details.")
			err = fmt.Errorf("oauth2 response error, %s", errDesc)
			cancel()
			return
		} else if code == "" {
			// no code, no error, redirect to auth url
			http.Redirect(w, r, authURL, http.StatusFound)
			return
		}
		// success
		w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, "Success! You may close this window.")
		cancel()
	})
	go func() {
		err = s.ListenAndServe()
		if err == http.ErrServerClosed {
			err = nil
		}
	}()
	<-ctx.Done()
	s.Shutdown(ctx)
	return code, tracerr.Wrap(err)
}

// Login to Platform.sh API via oauth.
func Login() error {
	output.LogInfo("Begin Platform.sh oauth login.")
	// generate codes
	state := generateRandom()
	codeVerify := generateRandom()
	h := sha256.New()
	h.Write([]byte(codeVerify))
	codeChallenge := base64URL(h.Sum(nil))
	// generate psh oauth url
	ctx := context.Background()
	conf := getOauthConfig()
	url := conf.AuthCodeURL(
		state,
		oauth2.AccessTypeOffline,
		oauth2.ApprovalForce,
		oauth2.SetAuthURLParam("code_challenge", codeChallenge),
		oauth2.SetAuthURLParam("code_challenge_method", "S256"),
	)
	output.LogDebug("OAUTH2 redirect URL.", url)
	// open local web server to accept response from psh oauth
	output.Info(fmt.Sprintf("Open %s in your browser to continue.", conf.RedirectURL))
	code, err := loginServer(url)
	if err != nil {
		return tracerr.Wrap(err)
	}
	// exchange code for token
	tok, err := conf.Exchange(
		ctx,
		code,
		oauth2.SetAuthURLParam("code_verifier", codeVerify),
	)
	if err != nil {
		return tracerr.Wrap(err)
	}
	// save token
	if err := saveToken(tok); err != nil {
		return tracerr.Wrap(err)
	}
	output.Info("Login successful!")
	return nil
}

// GetAPIClient returns the HTTP client used for API transactions.
func GetAPIClient() (*http.Client, error) {
	conf := getOauthConfig()
	tok, err := loadToken()
	if err != nil {
		return nil, tracerr.Wrap(err)
	}
	return conf.Client(context.Background(), tok), nil
}
