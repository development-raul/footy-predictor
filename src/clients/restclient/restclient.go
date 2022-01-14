package restclient

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/development-raul/footy-predictor/src/zlog"
	"net/http"
	"net/url"
	"strings"
)

var (
	enabledMocks = false
	mocks        = make(map[string]*Mock)
)

type Mock struct {
	Url        string
	HttpMethod string
	Response   *http.Response
	Err        error
}

func getMockId(httpMethod string, url string) string {
	return fmt.Sprintf("%s_%s", httpMethod, url)
}

func StartMockups() {
	enabledMocks = true
}

func FlushMockups() {
	mocks = make(map[string]*Mock)
}

func StopMockups() {
	enabledMocks = false
}

func AddMockup(mock Mock) {
	mocks[getMockId(mock.HttpMethod, mock.Url)] = &mock
}

func Post(url string, body interface{}, headers http.Header) (*http.Response, error) {
	if enabledMocks {
		mock := mocks[getMockId(http.MethodPost, url)]
		if mock == nil {
			return nil, errors.New("no mockup found for give request")
		}
		return mock.Response, mock.Err
	}

	jsonBytes, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(jsonBytes))
	if err != nil {
		zlog.Logger.Errorw("RestClient Post NewRequest", err)
		return nil, err
	}
	request.Header = headers

	client := http.Client{}
	return client.Do(request)
}

func Get(url string, headers http.Header) (*http.Response, error) {
	if enabledMocks {
		mock := mocks[getMockId(http.MethodGet, url)]
		if mock == nil {
			return nil, errors.New("no mockup found for give request")
		}
		return mock.Response, mock.Err
	}

	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		zlog.Logger.Errorw("RestClient Get NewRequest", err)
		return nil, err
	}
	request.Header = headers

	client := http.Client{}
	return client.Do(request)
}

func PostForm(url string, data url.Values, headers http.Header) (*http.Response, error) {
	if enabledMocks {
		mock := mocks[getMockId(http.MethodPost, url)]
		if mock == nil {
			return nil, errors.New("no mock found for given request")
		}
		return mock.Response, mock.Err
	}

	request, err := http.NewRequest(http.MethodPost, url, strings.NewReader(data.Encode()))
	if err != nil {
		zlog.Logger.Errorw("RestClient PostForm NewRequest", err)
		return nil, err
	}
	request.Header = headers

	client := http.Client{}
	return client.Do(request)
}