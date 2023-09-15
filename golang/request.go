package golang

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"os"

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
			if data, err := io.ReadAll(resp.Body); err == nil {
				if err := json.Unmarshal(data, &egeonErr); err == nil {
					return resp, egeonErr
				}
				return resp, errors.New(string(data))
			}
		}
		return resp, nil
	}
}

type RequestEditorFn func(ctx context.Context, req *retry.Request) error

// DoRequest - create request and read answer
// method can be GET, POST, PUT, DELETE (http method)
// user in context is required
// reqBody - can be nil
func DoRequest(ctx context.Context, client *retry.Client, method string, reqURL url.URL, reqBody []byte, reqEditors ...RequestEditorFn) ([]byte, error) {
	req, err := retry.NewRequest(method, reqURL.String(), reqBody)
	if err != nil {
		return nil, err
	}
	user, _ := ctx.Value(UserKey).(User)
	user.UsersGroups = nil
	reqID, _ := ctx.Value(RequestID).(string)
	if len(reqID) == 0 {
		reqID = FormRequestID(&user)
	}
	allowedRole, _ := ctx.Value(AllowedRoleKey).(string)
	userJSON, _ := json.Marshal(&user)
	sign := CreateSignature([]byte(os.Getenv(EgeonSecretKeyEnviron)), userJSON)
	req.Header.Add(SignatureHeaderKey, sign)
	req.Header.Add(UserHeaderKey, string(userJSON))
	req.Header.Add(RequestIDHeaderKey, reqID)
	req.Header.Add(AllowedRoleHeaderKey, allowedRole)
	req.Header.Add("Content-Type", "application/json")
	for i := range reqEditors {
		if err := reqEditors[i](ctx, req); err != nil {
			return nil, err
		}
	}

	if client.ErrorHandler == nil {
		appendDefaultErrorHandler(client)
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, EgeonError{Code: InternalError, Description: "Request failed " + " error " + err.Error()}
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if resp.StatusCode >= http.StatusMultipleChoices {
		return nil, errors.New(string(data))
	}
	return data, err
}
