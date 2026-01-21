package auth

import (
	"context"
	"io"
	// "encoding/json"
	"log"
	"net/http"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

type Config struct {
	ClientId     string
	ClientSecret string
	State        string
	RedirectURI  string
	Issuer			 string
}

var oauthConfig oauth2.Config
var ctx = context.Background()
var provider *oidc.Provider

func (auth *Config) SetupAuth() {
	provider, err := oidc.NewProvider(ctx, auth.Issuer)

	if err != nil {
		log.Fatal(err)
	}

	oauthConfig = oauth2.Config{
		ClientID:     auth.ClientId,
		ClientSecret: auth.ClientSecret,
		Endpoint:     provider.Endpoint(),
		RedirectURL:  auth.RedirectURI,
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
	}
}

func (auth *Config) LoginRequest(w http.ResponseWriter, r *http.Request) {
	// requestUri := fmt.Sprintf("https://sso.csh.rit.edu/auth/realms/csh/protocol/openid-connect/auth?client_id=%s&redirect_uri=%s&response_type=code&scope=openid&state=%s", auth.ClientId, auth.RedirectURI, auth.State)
	http.Redirect(w, r, oauthConfig.AuthCodeURL(auth.State), http.StatusFound)
}

func (auth *Config) LoginCallback(w http.ResponseWriter, r *http.Request) {
	// session, err := r.Cookie("AuthSession")
	//
	// if err != nil {
	// 	http.Error(w, "No auth session", http.StatusBadRequest)
	// 	return
	// }

	state := r.URL.Query().Get("state")

	if state != auth.State {
		http.Error(w, "Bad state", http.StatusBadRequest)
	}

	oauthToken, err := oauthConfig.Exchange(ctx, r.URL.Query().Get("code"))

	if err != nil {
		http.Error(w, "Failed to exchange token: "+err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("%s", oauthToken)
	
	req, err := http.NewRequest("GET", auth.Issuer+"/protocol/openid-connect/userinfo", nil)
	oauthToken.SetAuthHeader(req)

	res, err := http.DefaultClient.Do(req)
	jsonRaw, err := io.ReadAll(res.Body)

	log.Printf("%s", &jsonRaw)
	
}
