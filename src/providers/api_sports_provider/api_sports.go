package api_sports_provider

import (
	"encoding/json"
	"fmt"
	"github.com/development-raul/footy-predictor/src/clients/restclient"
	"github.com/development-raul/footy-predictor/src/domains/api_sports"
	"github.com/development-raul/footy-predictor/src/zlog"
	"io/ioutil"
	"net/http"
	"os"
)

func makeRequest(url string, action string) ([]byte, *api_sports.ErrorResponse) {
	// Make API Sports request
	res, err := restclient.Get(url, setHeaders())
	if err != nil {
		zlog.Logger.Error(fmt.Sprintf("APISportsProvider %s requestData: ", action), err)
		return nil, &api_sports.ErrorResponse{
			Message:    "Error making API request",
			StatusCode: http.StatusInternalServerError,
		}
	}
	// Read the response
	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		zlog.Logger.Error(fmt.Sprintf("APISportsProvider %s ReadAll: ", action), err)
		return nil, &api_sports.ErrorResponse{
			Message:    "Error reading API response",
			StatusCode: http.StatusInternalServerError,
		}
	}
	defer res.Body.Close()

	// Handle errors from API Sports
	if res.StatusCode != http.StatusOK {
		zlog.Logger.Warn("API Sports non 200 response: ", string(bytes))
		// Attempt to unmarshall the response into the API Sports ErrorResponse struct
		var errResponse api_sports.ErrorResponse
		if err := json.Unmarshal(bytes, &errResponse); err != nil {
			zlog.Logger.Error(fmt.Sprintf("APISportsProvider %s Unmarshal: ", action), err)
			return nil, &api_sports.ErrorResponse{
				Message:    "Error decoding API response",
				StatusCode: http.StatusInternalServerError,
			}
		}
		// Add the status code of the response
		errResponse.StatusCode = int64(res.StatusCode)
		return nil, &errResponse
	}

	zlog.Logger.Info("API Sports 200 response: ", string(bytes))

	return bytes, nil
}

func GetCountries() ([]api_sports.CountriesResponse, *api_sports.ErrorResponse) {
	url := fmt.Sprintf("%s/countries", os.Getenv("AS_BASE_URL"))
	// Make the request
	bytes, err := makeRequest(url, "GetCountries")
	if err != nil {
		return nil, err
	}

	var result api_sports.GetCountriesOutput
	// Handle success response from API Sports
	if err := json.Unmarshal(bytes, &result); err != nil {
		zlog.Logger.Error("APISportsProvider GetCountries Unmarshal: ", err)
		return nil, &api_sports.ErrorResponse{
			Message:    "Error decoding API response",
			StatusCode: http.StatusInternalServerError,
		}
	}

	return result.Response, nil
}

func GetSeasons() ([]int64, *api_sports.ErrorResponse) {
	url := fmt.Sprintf("%s/leagues/seasons", os.Getenv("AS_BASE_URL"))
	// Make the request
	bytes, err := makeRequest(url, "GetSeasons")
	if err != nil {
		return nil, err
	}
	// Handle success response from API Sports
	var result api_sports.GetSeasonsOutput
	if err := json.Unmarshal(bytes, &result); err != nil {
		zlog.Logger.Error("APISportsProvider GetSeasons Unmarshal: ", err)
		return nil, &api_sports.ErrorResponse{
			Message:    "Error decoding API response",
			StatusCode: http.StatusInternalServerError,
		}
	}
	return result.Response, nil
}

func setHeaders() http.Header {
	headers := http.Header{}
	headers.Set("Content-type", "application/json")
	headers.Set("Accept", "application/json")
	headers.Set("x-rapidapi-key", os.Getenv("AS_KEY"))
	headers.Set("x-rapidapi-host", os.Getenv("AS_HOST"))

	return headers
}
