package sqlreq

import (
	"io"
	"net/http"

	"github.com/scootpl/sqlreq/parser"
)

type result struct {
	jsonParam             map[string]any
	headerParam           map[string]string
	url                   string
	method                string
	outputHeaders         []string
	basicLogin, basicPass string
	timeout               int
	payload               any
}

// SelectRESPONSE returns *http.Response, http.StatusCode and error
func SelectRESPONSE(query string, params ...any) (*http.Response, int, error) {
	r := new(result)
	if err := doSelect(parser.SelectBody, query, r, params...); err != nil {
		return nil, 0, err
	}

	resp, err := doReq(r)
	if err != nil {
		return nil, 0, err
	}

	return resp, resp.StatusCode, nil
}

// SelectHEADERS returns all headers, http.StatusCode and error
func SelectHEADERS(query string, params ...any) (http.Header, int, error) {
	r := new(result)
	if err := doSelect(parser.SelectBody, query, r, params...); err != nil {
		return nil, 0, err
	}

	resp, err := doReq(r)
	if err != nil {
		return nil, 0, err
	}
	resp.Body.Close()

	return resp.Header, resp.StatusCode, nil
}

// SelectHEADER returns the indicated header or headers, http.StatusCode and error
func SelectHEADER(query string, params ...any) (map[string]string, int, error) {
	r := new(result)
	if err := doSelect(parser.SelectHeader, query, r, params...); err != nil {
		return nil, 0, err
	}

	resp, err := doReq(r)
	if err != nil {
		return nil, 0, err
	}

	defer resp.Body.Close()

	headers := make(map[string]string)
	for _, k := range r.outputHeaders {
		headers[k] = resp.Header.Get(k)
	}

	return headers, resp.StatusCode, nil
}

// SelectBODY returns body io.ReadCloser, http.StatusCode and error
func SelectBODY(query string, params ...any) (io.ReadCloser, int, error) {
	r := new(result)
	if err := doSelect(parser.SelectBody, query, r, params...); err != nil {
		return nil, 0, err
	}

	resp, err := doReq(r)
	if err != nil {
		return nil, 0, err
	}

	return resp.Body, resp.StatusCode, nil
}
