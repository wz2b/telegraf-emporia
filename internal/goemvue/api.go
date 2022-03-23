package goemvue

import (
	"errors"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"log"
	"net/http"
	"net/url"
	"time"
)

type EmVueCloudSession struct {
	ClientId             string
	httpClient           *http.Client
	CustomerInfo         *CustomerInfoWithDevices
	authenticationResult *cognitoidentityprovider.AuthenticationResultType
	cognitoClient        *cognitoidentityprovider.CognitoIdentityProvider
	tokenExpiresAt       *time.Time
	DebugLog             *log.Logger
	username             string `sensitive:"true"`
	password             string `sensitive:"true"`
}

func NewEmVueCloud(username string, password string) *EmVueCloudSession {
	return &EmVueCloudSession{
		httpClient: &http.Client{},
		username:   username,
		password:   password,
	}
}

func (t *EmVueCloudSession) apiGet(u *url.URL) (*http.Response, error) {
	err := t.reauthorizeIfRequired()
	if err != nil {
		return nil, err
	}

	urlString := u.String()

	req, err := http.NewRequest("GET", urlString, nil)

	if err != nil {
		return nil, err
	}

	token := t.authenticationResult.IdToken
	if token == nil {
		return nil, errors.New("Not authorized")
	}
	req.Header.Add("authtoken", *token)

	resp, err := t.httpClient.Do(req)

	return resp, err
}

func NewApiUrl(path string) *url.URL {
	return &url.URL{
		Scheme:      "https",
		Opaque:      "",
		User:        nil,
		Host:        "api.emporiaenergy.com",
		Path:        path,
		RawPath:     "",
		ForceQuery:  false,
		RawQuery:    "",
		Fragment:    "",
		RawFragment: "",
	}
}
