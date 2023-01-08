package sqlreq

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/scootpl/sqlreq/parser"
)

func doReq(r *result) (*http.Response, error) {
	var payload any
	body := bytes.NewBuffer([]byte{})

	if r.payload != nil {
		payload = r.payload
	} else {
		payload = r.jsonParam
	}
	if err := json.NewEncoder(body).Encode(payload); err != nil {
		return nil, err
	}

	client := http.Client{
		Timeout: time.Duration(r.timeout) * time.Second,
	}

	req, err := http.NewRequest(r.method, r.url, body)
	if err != nil {
		return nil, err
	}

	for k, v := range r.headerParam {
		req.Header.Add(k, v)
	}

	if r.basicLogin != "" && r.basicPass != "" {
		req.SetBasicAuth(r.basicLogin, r.basicPass)
	}

	return client.Do(req)
}

// decompositionParams looks for a payload and string objects
func decompositionParams(params ...any) (payload []any, str []any) {
	for _, param := range params {
		if _, ok := param.(string); ok {
			str = append(str, param)
		} else {
			payload = append(payload, param)
		}
	}
	return
}

func doSelect(mode int, query string, result *result, params ...any) error {
	p := parser.New(mode)

	payloadObjects, stringObjects := decompositionParams(params...)
	q := fmt.Sprintf(query, stringObjects...)

	if err := p.Parse(p.Tokenizer(q)); err != nil {
		return err
	}

	result.jsonParam = make(map[string]any)
	result.headerParam = make(map[string]string)

	for i, k := range p.HeaderKey {
		result.headerParam[k] = p.HeaderValue[i]
	}

	for i, k := range p.JsonKey {
		result.jsonParam[k] = p.JsonValue[i]
	}

	if p.Payload && len(payloadObjects) == 0 {
		return errors.New("no payload")
	}

	if p.Payload {
		result.payload = payloadObjects[0]
	}

	if p.Timeout != "" {
		t, err := strconv.Atoi(p.Timeout)
		if err != nil {
			return fmt.Errorf("wrong timeout format error: %w", err)
		}
		result.timeout = t
	} else {
		result.timeout = 30
	}

	result.url = p.Url
	result.method = p.Method
	result.outputHeaders = p.OutputHeaders
	result.basicLogin = p.BasicLogin
	result.basicPass = p.BasicPass
	return nil
}
