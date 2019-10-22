package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
)

var (
	refreshToken = os.Getenv("RefreshToken")
	basic        = os.Getenv("BASIC")
	ssmSVC       = ssm.New(session.New(), &aws.Config{Region: aws.String(region)})
)

const (
	region                   = "ap-northeast-1"
	urlRefreshToken          = "https://api.fitbit.com/oauth2/token"
	fitbitTokenParameterName = "AccessToken"
)

// AccessToken set value
type AccessToken struct {
	AccessToken string `json:"access_token"`
}

func main() {
	lambda.Start(Handler)
}
// func main() {

// Handler lambda
func Handler() {
	client := &http.Client{}

	values := url.Values{}
	values.Add("grant_type", "refresh_token")
	values.Add("refresh_token", refreshToken)

	req, _ := http.NewRequest("POST", urlRefreshToken, strings.NewReader(values.Encode()))
	req.Header.Set("Authorization", " Basic "+basic)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)
	bodyStr := string(body)
	// fmt.Println(bodyStr)
	jsonBytes := ([]byte)(bodyStr)

	at := &AccessToken{}
	if err = json.Unmarshal(jsonBytes, at); err != nil {
		fmt.Println(err)
	}
	//fmt.Println(at.AccessToken)
	PutToken(at.AccessToken)
}

// PutToken to ssm parameter
func PutToken(accessToken string) {
	_, err := ssmSVC.PutParameter(
		&ssm.PutParameterInput{
			Name:      aws.String(fitbitTokenParameterName),
			Overwrite: aws.Bool(true),
			Type:      aws.String("SecureString"),
			Value:     aws.String(accessToken),
		},
	)
	if err != nil {
		fmt.Println(err)
	}
}
