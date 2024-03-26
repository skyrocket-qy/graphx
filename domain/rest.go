package domain

type Response struct {
	Message string `json:"message"`
}

type ErrResponse struct {
	Error string `json:"error"`
}

type EdgesResponse struct {
	Edges []Edge `json:"data"`
}

type StringsResponse struct {
	Data []string `json:"data"`
}
