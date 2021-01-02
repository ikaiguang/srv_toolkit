package tkcurl

import (
	"context"
	"crypto/tls"
	"github.com/go-kratos/kratos/pkg/net/http/blademaster/binding"
	tkapp "github.com/ikaiguang/srv_toolkit/app"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

// content type
const (
	ContentTypeJSON     = binding.MIMEJSON
	ContentTypeJSONUtf8 = tkapp.ContentTypeJSON
	ContentTypePB       = tkapp.ContentTypePB
)

const (
	// timeout
	DefaultTimeout = time.Second * 3

	// token
	DefaultTokenKey = "Authorization"
)

// RequestParam .
type RequestParam struct {
	URL string

	Timeout     time.Duration
	ContentType string

	NeedToken  bool
	TokenKey   string
	TokenValue string

	Body io.Reader
}

// JSONRequestParam .
func JSONRequestParam() *RequestParam {
	return &RequestParam{
		Timeout:     DefaultTimeout,
		ContentType: ContentTypeJSON,

		NeedToken: false,
		TokenKey:  DefaultTokenKey,
	}
}

// ProtobufRequestParam .
func ProtobufRequestParam() *RequestParam {
	return &RequestParam{
		Timeout:     DefaultTimeout,
		ContentType: ContentTypePB,

		NeedToken: false,
		TokenKey:  DefaultTokenKey,
	}
}

// Get get
func Get(requestParam *RequestParam) (code int, bodyBytes []byte, err error) {
	// http
	httpReq, err := http.NewRequest(http.MethodGet, requestParam.URL, nil)
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	httpReq.Header.Set("Content-Type", requestParam.ContentType)
	if requestParam.NeedToken {
		httpReq.Header.Set(requestParam.TokenKey, requestParam.TokenValue)
	}

	// 超时
	ctx, cancel := context.WithTimeout(context.Background(), requestParam.Timeout)
	defer cancel()
	httpReq.WithContext(ctx)

	return requestResp(httpReq, requestParam)
}

// Post post
func Post(requestParam *RequestParam) (httpCode int, bodyBytes []byte, err error) {
	// http
	httpReq, err := http.NewRequest(http.MethodPost, requestParam.URL, requestParam.Body)
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	httpReq.Header.Set("Content-Type", requestParam.ContentType)
	if requestParam.NeedToken {
		httpReq.Header.Set(requestParam.TokenKey, requestParam.TokenValue)
	}

	// 超时
	ctx, cancel := context.WithTimeout(context.Background(), requestParam.Timeout)
	defer cancel()
	httpReq.WithContext(ctx)

	return requestResp(httpReq, requestParam)
}

// requestResp request
func requestResp(httpReq *http.Request, requestParam *RequestParam) (httpCode int, bodyBytes []byte, err error) {
	// client
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	httpClient := &http.Client{Transport: tr}
	defer httpClient.CloseIdleConnections()

	// 请求
	httpClient.Timeout = requestParam.Timeout
	httpResp, err := httpClient.Do(httpReq)
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	defer func() { _ = httpResp.Body.Close() }()

	// resp
	httpCode = httpResp.StatusCode
	bodyBytes, err = ioutil.ReadAll(httpResp.Body)
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}
