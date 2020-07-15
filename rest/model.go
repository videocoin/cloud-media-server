package rest

type requestData struct {
	URL string `json:"url"`
}

type progressResponse struct {
	Progress int64 `json:"progress"`
}
