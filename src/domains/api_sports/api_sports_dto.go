package api_sports

type Paging struct {
	Current int64 `json:"current"`
	Total   int64 `json:"total"`
}

type Errors struct {
	Time     string `json:"time"`
	Bug      string `json:"bug"`
	Report   string `json:"report"`
	Endpoint string `json:"endpoint"`
}
type CountriesResponse struct {
	Name string `json:"name"`
	Code string `json:"code"`
	Flag string `json:"flag"`
}
type ErrorResponse struct {
	Message    string `json:"message"`
	StatusCode int64  `json:"status_code"`
}
type GetCountriesOutput struct {
	Get      string              `json:"get"`
	Errors   []Errors            `json:"errors"`
	Results  int64               `json:"results"`
	Paging   Paging              `json:"paging"`
	Response []CountriesResponse `json:"response"`
}
