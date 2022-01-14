package api_sports_provider

import (
	"errors"
	"github.com/development-raul/footy-predictor/src/clients/restclient"
	"github.com/development-raul/footy-predictor/src/domains/api_sports"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"
)

type FakeReaderCloser struct{}

func (i *FakeReaderCloser) Read([]byte) (int, error) {
	return 0, errors.New("test error")
}

func (i *FakeReaderCloser) Close() error {
	return nil
}

func TestAPISportsProvider_GetCountries(t *testing.T) {
	testCases := []struct {
		title       string
		apiMock     restclient.Mock
		withMock    bool
		baseURL     string
		expectedRes []api_sports.CountriesResponse
		expectedErr *api_sports.ErrorResponse
	}{
		{
			title:       "error restclient.Get",
			baseURL:     "invalid-url",
			expectedRes: nil,
			expectedErr: &api_sports.ErrorResponse{
				Message:    "Error making API request",
				StatusCode: http.StatusInternalServerError,
			},
		},
		{
			title: "error ioutil.ReadAll",
			apiMock: restclient.Mock{
				Url:        "https://test.com/api_sports",
				HttpMethod: http.MethodGet,
				Response: &http.Response{
					StatusCode: 200,
					Body:       &FakeReaderCloser{},
				},
			},
			withMock:    true,
			baseURL:     "https://test.com",
			expectedRes: nil,
			expectedErr: &api_sports.ErrorResponse{
				Message:    "Error reading API response",
				StatusCode: http.StatusInternalServerError,
			},
		},
		{
			title: "error non 200 json.Unmarshal",
			apiMock: restclient.Mock{
				Url:        "https://test.com/api_sports",
				HttpMethod: http.MethodGet,
				Response: &http.Response{
					StatusCode: 499,
					Body:       io.NopCloser(strings.NewReader(`{"response does not match ErrorResponse struct"}`)),
				},
			},
			withMock:    true,
			baseURL:     "https://test.com",
			expectedRes: nil,
			expectedErr: &api_sports.ErrorResponse{
				Message:    "Error decoding API response",
				StatusCode: http.StatusInternalServerError,
			},
		},
		{
			title: "error non 200 response",
			apiMock: restclient.Mock{
				Url:        "https://test.com/api_sports",
				HttpMethod: http.MethodGet,
				Response: &http.Response{
					StatusCode: 499,
					Body:       io.NopCloser(strings.NewReader(`{"message": "Something went wrong while fetching details. Try again later."}`)),
				},
			},
			withMock:    true,
			baseURL:     "https://test.com",
			expectedRes: nil,
			expectedErr: &api_sports.ErrorResponse{
				Message:    "Something went wrong while fetching details. Try again later.",
				StatusCode: 499,
			},
		},
		{
			title: "error 200 json.Unmarshal",
			apiMock: restclient.Mock{
				Url:        "https://test.com/api_sports",
				HttpMethod: http.MethodGet,
				Response: &http.Response{
					StatusCode: 200,
					Body:       io.NopCloser(strings.NewReader(`{"response does not match ErrorResponse struct"}`)),
				},
			},
			withMock:    true,
			baseURL:     "https://test.com",
			expectedRes: nil,
			expectedErr: &api_sports.ErrorResponse{
				Message:    "Error decoding API response",
				StatusCode: http.StatusInternalServerError,
			},
		},
		{
			title: "success",
			apiMock: restclient.Mock{
				Url:        "https://test.com/api_sports",
				HttpMethod: http.MethodGet,
				Response: &http.Response{
					StatusCode: 200,
					Body:       io.NopCloser(strings.NewReader(`{"get":"api_sports","parameters":[],"errors":[],"results":1,"paging":{"current":1,"total":1},"response":[{"name":"Germany","code":"DE","flag":"https://test.com/flags/de.svg"}]}`)),
				},
			},
			withMock: true,
			baseURL:  "https://test.com",
			expectedRes: []api_sports.CountriesResponse{
				{
					Name: "Germany",
					Code: "DE",
					Flag: "https://test.com/flags/de.svg",
				},
			},
			expectedErr: nil,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.title, func(t *testing.T) {
			if testCase.withMock {
				restclient.StartMockups()
				restclient.AddMockup(testCase.apiMock)
			}
			os.Setenv("AS_BASE_URL", testCase.baseURL)

			res, err := GetCountries()
			assert.Equal(t, testCase.expectedRes, res)
			assert.Equal(t, testCase.expectedErr, err)

			restclient.FlushMockups()
		})
	}
}
