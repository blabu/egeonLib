package golang

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"

	retry "github.com/hashicorp/go-retryablehttp"
	"github.com/rs/zerolog/log"
)

//DoRequest - create request and read answer
//method can be GET, POST, PUT, DELETE (http method)
//reqBody - can be nil
func DoRequest(ctx context.Context, client *retry.Client, method string, reqURL url.URL, reqBody []byte) ([]byte, error) {
	reqID, _ := ctx.Value(RequestID).(string)
	user, ok := ctx.Value(UserKey).(User)
	if !ok {
		log.Error().Str(RequestIDHeaderKey, reqID).Str(UserHeaderKey, user.Email).Msg("Undefined user")
		return nil, Errors["undefUser"]
	}
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
				} else {
					return resp, errors.New(string(data))
				}
			}
		}
		return resp, err
	}

	if len(reqID) == 0 {
		reqID = FormRequestID(&user)
	}
	sign, ok := ctx.Value(SignKey).(string)
	if !ok {
		log.Warn().Str(RequestIDHeaderKey, reqID).Str(UserHeaderKey, user.Email).Msg("Undefined user signature. But we continue")
	}
	userJSON, _ := json.Marshal(&user)
	req, err := retry.NewRequest(method, reqURL.String(), reqBody)
	if err != nil {
		log.Err(err).Str(RequestIDHeaderKey, reqID).Msg("When try create request. Can not continue")
		return nil, err
	}
	req.Header.Add(SignatureHeaderKey, sign)
	req.Header.Add(UserHeaderKey, string(userJSON))
	req.Header.Add(RequestIDHeaderKey, reqID)
	resp, err := client.Do(req)
	if err != nil {
		log.Err(err).Str(RequestIDHeaderKey, reqID).Str(UserHeaderKey, user.Email).Msg("When send request " + req.URL.String())
		return nil, err
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if resp.StatusCode >= 300 {
		log.Info().Str(RequestIDHeaderKey, reqID).Str(UserHeaderKey, user.Email).Msg(string(data))
		return nil, errors.New(string(data))
	}
	return data, err
}
