package security

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"ghwebhooks/config"
	"io/ioutil"
	"net/http"
)

func retrieveSignature(req *http.Request) (string, error) {
	var signature string
	var err error
	var hubSignature = req.Header.Get("X-Hub-Signature")

	if len(hubSignature) < 5 {
		err = errors.New("signature invalid")
	} else {
		signature = hubSignature[5:]
	}

	return signature, err
}

func valid(signature string, req *http.Request, config *config.Config) bool {
	body, err := ioutil.ReadAll(req.Body)
	req.Body = ioutil.NopCloser(bytes.NewReader(body))

	if err != nil {
		return false
	}

	mac := hmac.New(sha1.New, []byte(config.Secret))
	mac.Write(body)
	expectedMAC := hex.EncodeToString(mac.Sum(nil))

	return hmac.Equal([]byte(signature), []byte(expectedMAC))
}

func Secured(handler http.HandlerFunc, config *config.Config) http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		if config.Dev {
			handler(resp, req)
			return
		}

		signature, err := retrieveSignature(req)

		if err != nil {
			resp.WriteHeader(403)
		} else if valid(signature, req, config) {
			handler(resp, req)
		} else {
			resp.WriteHeader(403)
		}
	}
}
