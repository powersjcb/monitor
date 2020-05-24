package gateway

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/powersjcb/monitor/go/src/lib/crypto"
	"github.com/powersjcb/monitor/go/src/lib/httpclient"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"gopkg.in/go-playground/validator.v9"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

type googleUserInfo struct {
	ID    		  string `json:"id"`
	Email 		  string `json:"email" validate:"email"`
	VerifiedEmail bool	 `json:"verified_email"`
	Picture  	  string `json:"picture"`
}

func getConfig() *oauth2.Config {
	return &oauth2.Config{
		RedirectURL:  "http://localhost:8080/auth/google/callback",
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}
}

type JWTConfig struct {
	PublicKey ecdsa.PublicKey
	PrivateKey ecdsa.PrivateKey
}

const csrfCookieName = "monitor_csrf"

func (s HTTPServer) GoogleLoginHandler(rw http.ResponseWriter, r *http.Request) {
	stateString := crypto.GetToken(32)
	url := getConfig().AuthCodeURL(stateString)
	var cookie = http.Cookie{
		Name:       csrfCookieName,
		Value:      stateString,
		Path:       "/",
		Domain:     r.URL.Host,
		Expires:    time.Now().Add(5 * time.Minute),
		MaxAge:     0, // session only cookie
		Secure:     gaeHTTPS(r),
		HttpOnly:   true,
		SameSite:   http.SameSiteStrictMode,
	}
	fmt.Println(url)
	http.SetCookie(rw, &cookie) // csrf cookie
	http.Redirect(rw, r, url, http.StatusTemporaryRedirect)
}

func (s HTTPServer) GoogleCallbackHandler(rw http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(csrfCookieName)
	if err != nil || cookie == nil || cookie.Value == "" {
		rw.Write([]byte("invalid csrf token"))
		rw.WriteHeader(400)
		return
	}
	content, err := getUserInfo(r.Context(), cookie.Value, r.FormValue("state"), r.FormValue("code"))
	if err != nil {
		fmt.Println(err.Error())
		http.Redirect(rw, r, "/", http.StatusTemporaryRedirect)
		return
	}

	var user googleUserInfo
	err = json.Unmarshal(content, &user)
	if err != nil {
		fmt.Println(err.Error())
		rw.WriteHeader(500)
		return
	}
	err = s.setLoginCookie(rw, r, user, "google")
	if err != nil {
		fmt.Println(err.Error())
		rw.WriteHeader(500)
		return
	}
	fmt.Println("logged in!")
	http.Redirect(rw, r, "/", http.StatusTemporaryRedirect)
}

const jwtCookieName = "jtw-cookie"

func (s HTTPServer) setLoginCookie(rw http.ResponseWriter, r *http.Request, userInfo googleUserInfo, provider string) error {
	v := validator.New()
	err := v.Struct(userInfo)
	if err != nil {
		return err
	}
	if !userInfo.VerifiedEmail {
		return errors.New("email has not been verified with provider")
	}
	// https://tools.ietf.org/html/rfc7519#section-4.1
	token := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.StandardClaims{
		Issuer:    r.Host,
		Subject:   userInfo.Email,
		Audience:  "https://monitor.jacobpowers.me,https://jacobpowers.me",
		ExpiresAt: time.Now().Add(time.Minute * 60).Unix(),
		IssuedAt:  time.Now().Unix(),
	})
	ss, err := token.SignedString(s.jwtConfig.PrivateKey)
	if err != nil {
		return err
	}
	c := http.Cookie{
		Name:       jwtCookieName,
		Value:      ss,
		Path:       "/",
		Domain:     r.Host,
		Expires:    time.Now().Add(7 * 24 * time.Hour),
		Secure:     gaeHTTPS(r),
		HttpOnly:   true,
		SameSite:   http.SameSiteStrictMode,
	}
	http.SetCookie(rw, &c)
	return nil
}

func gaeHTTPS(r *http.Request) bool {
	return r.Header.Get("X-AppEngine-Https") == "on"
}

func getUserInfo(ctx context.Context, csrfState, state, code string) ([]byte, error) {
	if state != csrfState {
		return nil, fmt.Errorf("invalid oauth state")
	}
	token, err := getConfig().Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("code exchange failed: %s", err.Error())
	}
	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed getting user info: %s", err.Error())
	}
	defer httpclient.CloseBody(response)
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed reading response body: %s", err.Error())
	}
	return contents, nil
}