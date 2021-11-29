package golang

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"

	retry "github.com/hashicorp/go-retryablehttp"
)

func appendDefaultErrorHandler(client *retry.Client) {
	client.ErrorHandler = func(resp *http.Response, err error, numTries int) (*http.Response, error) {
		if err != nil {
			return resp, err
		}
		if numTries < client.RetryMax {
			return resp, nil //retry again
		}
		if resp.StatusCode > 399 {
			var egeonErr EgeonError
			defer resp.Body.Close()
			if data, err := ioutil.ReadAll(resp.Body); err == nil {
				if err := json.Unmarshal(data, &egeonErr); err == nil {
					return resp, egeonErr
				}
				return resp, errors.New(string(data))
			}
		}
		return resp, nil
	}
}

//DoRequest - create request and read answer
//method can be GET, POST, PUT, DELETE (http method)
// user in context is required
//reqBody - can be nil
func DoRequest(ctx context.Context, client *retry.Client, method string, reqURL url.URL, reqBody []byte) ([]byte, error) {
	if client.ErrorHandler == nil {
		appendDefaultErrorHandler(client)
	}
	req, err := retry.NewRequest(method, reqURL.String(), reqBody)
	if err != nil {
		return nil, err
	}

	user, _ := ctx.Value(UserKey).(User)
	reqID, _ := ctx.Value(RequestID).(string)
	if len(reqID) == 0 {
		reqID = FormRequestID(&user)
	}
	sign, _ := ctx.Value(SignKey).(string)
	allowedRole, _ := ctx.Value(AllowedRoleKey).(string)
	userJSON, _ := json.Marshal(&user)
	req.Header.Add(SignatureHeaderKey, sign)
	req.Header.Add(UserHeaderKey, string(userJSON))
	req.Header.Add(RequestIDHeaderKey, reqID)
	req.Header.Add(AllowedRoleHeaderKey, allowedRole)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("content-type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if resp.StatusCode >= http.StatusMultipleChoices {
		return nil, errors.New(string(data))
	}
	return data, err
}
