package golang

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
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
	errorWriter, ok := client.Logger.(retry.Logger)
	if !ok {
		errorWriter = LoggerWrapper(func(format string, args ...interface{}) {
			fmt.Printf(format, args...)
		})
	}
	reqID, _ := ctx.Value(RequestID).(string)
	user, ok := ctx.Value(UserKey).(User)
	if !ok {
		errorWriter.Printf("User undefined\n")
		return nil, Errors["undefUser"]
	}
	if client.ErrorHandler == nil {
		appendDefaultErrorHandler(client)
	}
	if len(reqID) == 0 {
		reqID = FormRequestID(&user)
	}
	sign, ok := ctx.Value(SignKey).(string)
	if !ok {
		errorWriter.Printf("Undefined user signature. But we continue\n")
	}
	userJSON, _ := json.Marshal(&user)
	req, err := retry.NewRequest(method, reqURL.String(), reqBody)
	if err != nil {
		return nil, err
	}
	req.Header.Add(SignatureHeaderKey, sign)
	req.Header.Add(UserHeaderKey, string(userJSON))
	req.Header.Add(RequestIDHeaderKey, reqID)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if resp.StatusCode >= 300 {
		return nil, errors.New(string(data))
	}
	return data, err
}
