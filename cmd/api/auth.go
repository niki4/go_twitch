package api

import (
	"encoding/json"
	"fmt"
	"github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
	"html/template"
	"net/http"
)

type URLData struct {
	Base         string
	ClientID     string
	RedirectURI  string
	GrantType    string
	ResponseType string
	Scope        string
}

type SecretData struct {
	ClientSecret  string
	UserOAuthCode []byte
}

type LoginData struct {
	Title  string
	URL    URLData
	Secret SecretData
}

type UAT struct {
	UserAccessToken string   `json:"access_token"`
	RefreshToken    string   `json:"refresh_token"`
	ExpiresIn       int      `json:"expires_in"`
	Scope           []string `json:"scope"`
	TokenType       string   `json:"token_type"`
}

// ShowLoginPage dispatch request for Login Page
func (r *Router) ShowLoginPage(ctx *routing.Context) error {

	r.logger.Info(
		"Inside ShowLoginPage() view function",
		zap.ByteString("Method", ctx.Request.Header.Method()),
	)

	userAuthCookie := ctx.Request.Header.Cookie("Authorization")
	if userAuthCookie != nil {
		r.logger.Info("User provided Auth token", zap.ByteString("Authorization", userAuthCookie))
		ctx.Redirect("/streams", fasthttp.StatusOK)
	}

	// OAuth Authorization Code Flow. 1st stage
	loginData := LoginData{
		Title: "Login with Twitch account",
		URL: URLData{
			Base:         "https://id.twitch.tv/oauth2/authorize",
			ClientID:     r.ClientID,
			RedirectURI:  fmt.Sprintf("http://%s:%s", r.Host, r.Port),
			ResponseType: "code",
			Scope:        "viewing_activity_read",
		},
	}

	args := ctx.QueryArgs()
	OAuthCode := args.Peek("code")
	if OAuthCode != nil {
		// moving to 2nd stage
		loginData.URL.Base = "https://id.twitch.tv/oauth2/token"
		loginData.Secret.ClientSecret = r.ClientSecret
		loginData.Secret.UserOAuthCode = OAuthCode
		loginData.URL.GrantType = "authorization_code"

		UATdata, err := requestUserAccessToken(loginData)
		if err != nil {
			r.logger.Error("Failure on requestUserAccessToken:", zap.Error(err))
		}
		r.logger.Info("UserAccessToken received", zap.Any("UAT", UATdata))

		userCookie := fasthttp.Cookie{}
		userCookie.SetKey("Authorization")
		userCookie.SetValue("Bearer " + UATdata.UserAccessToken)
		ctx.Response.Header.SetCookie(&userCookie)
		ctx.Redirect("/streams", fasthttp.StatusOK)
	}

	ctx.SetContentType("text/html")

	tmpl := template.Must(template.ParseFiles("templates/login.html"))
	err := tmpl.Execute(ctx, loginData)
	if err != nil {
		fmt.Fprintf(ctx, "%s", err)
	}

	return nil
}

func requestUserAccessToken(ld LoginData) (UAT, error) {
	u := fmt.Sprintf("%s?client_id=%s&client_secret=%s&code=%s&grant_type=%s&redirect_uri=%s",
		ld.URL.Base,
		ld.URL.ClientID,
		ld.Secret.ClientSecret,
		ld.Secret.UserOAuthCode,
		ld.URL.GrantType,
		ld.URL.RedirectURI,
	)

	resp, err := http.Post(u, "application/json", nil)
	if err != nil {
		return UAT{}, err
	}
	defer resp.Body.Close()

	uat := UAT{}
	err = json.NewDecoder(resp.Body).Decode(&uat)
	if err != nil {
		return UAT{}, err
	}

	return uat, nil
}
