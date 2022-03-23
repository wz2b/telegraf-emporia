package goemvue

import (
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	cognito "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"time"
)

var CLIENT_ID = "4qte47jbstod8apnfic0bunmrq"

func (t *EmVueCloudSession) authorize() error {
	params := map[string]*string{
		"USERNAME": aws.String(t.username),
		"PASSWORD": aws.String(t.password),
	}

	authTry := &cognito.InitiateAuthInput{
		AuthFlow:       aws.String("USER_PASSWORD_AUTH"),
		AuthParameters: params,
		ClientId:       aws.String(CLIENT_ID),
	}

	conf := &aws.Config{Region: aws.String("us-east-2")}
	sess, err := session.NewSession(conf)
	if err != nil {
		panic(err)
	}

	t.cognitoClient = cognito.New(sess)

	authResult, err := t.cognitoClient.InitiateAuth(authTry)

	if err != nil {
		return err
	}

	if authResult.AuthenticationResult != nil {
		t.authenticationResult = authResult.AuthenticationResult

		/*
		 * Calculate the expiration time
		 */
		var expSeconds int64 = *(authResult.AuthenticationResult.ExpiresIn)
		exp := time.Now().Add(time.Duration(expSeconds) * time.Second)
		t.tokenExpiresAt = &exp
	}

	if t.DebugLog != nil {
		t.DebugLog.Println(authResult.AuthenticationResult)
	}

	return nil
}

func (t *EmVueCloudSession) refreshToken() error {

	if t.authenticationResult == nil || t.authenticationResult.RefreshToken == nil {
		return errors.New("There is no refresh token, cannot refresh")
	}
	params := map[string]*string{
		"REFRESH_TOKEN": aws.String(*t.authenticationResult.RefreshToken),
	}

	authTry := &cognito.InitiateAuthInput{
		AuthFlow:       aws.String("REFRESH_TOKEN_AUTH"),
		AuthParameters: params,
		ClientId:       aws.String(CLIENT_ID),
	}

	conf := &aws.Config{Region: aws.String("us-east-2")}
	sess, err := session.NewSession(conf)
	if err != nil {
		panic(err)
	}

	t.cognitoClient = cognito.New(sess)

	authResult, err := t.cognitoClient.InitiateAuth(authTry)

	if err != nil {
		return err
	}

	if authResult.AuthenticationResult != nil {
		oldRefreshToken := t.authenticationResult.RefreshToken
		t.authenticationResult = authResult.AuthenticationResult
		t.authenticationResult.RefreshToken = oldRefreshToken

		/*
		 * Calculate the new expiration time
		 */
		var expSeconds int64 = *(authResult.AuthenticationResult.ExpiresIn)
		exp := time.Now().Add(time.Duration(expSeconds) * time.Second)
		t.tokenExpiresAt = &exp
	}


	return nil
}

func (t *EmVueCloudSession) timeLeftOnToken() (time.Duration, error) {

	if t.authenticationResult == nil {
		return 0, errors.New("There is no token to check")
	}

	left := t.tokenExpiresAt.Sub(time.Now())
	return left, nil
}

func (t *EmVueCloudSession) percentLeftOnToken() (float32, error) {
	left, err := t.timeLeftOnToken()

	if err != nil {
		return 0.0, err
	}

	life := time.Duration(*t.authenticationResult.ExpiresIn) * time.Second


	t1 := left.Milliseconds()
	t2 := life.Milliseconds()
	pctLeft := float32(t1) / float32(t2)

	if pctLeft <= 0 {
		return 0.0, nil
	} else {
		return pctLeft, nil
	}
}

func (t *EmVueCloudSession) reauthorizeIfRequired() error {
	/*
	 * See how much time is left on the token
	 */
	pctTimeLeft, err  := t.percentLeftOnToken()

	if t.DebugLog != nil {
		t.DebugLog.Printf("%3.1f%% left on auth token\n", pctTimeLeft * 100.0)
	}

	if err != nil {
		return t.authorize()
	} else if pctTimeLeft < 0.3333 {
		/*
		 * The token has expired or is about to expire, try to use the
		 * refresh token to get a new one
		 */
		err = t.refreshToken()
		if err == nil {
			if t.DebugLog != nil {
				t.DebugLog.Println("Sucessfully refreshed token")
			}
		} else {

			/*
			 * Refresh failed, get a new token
			 */
			if t.DebugLog != nil {
				t.DebugLog.Printf("Token refresh failed (will retry as a new authorization): %s\n", err)
			}


			err := t.authorize()
			if err != nil {
				return err
			}

		}
	}

	return nil
}
