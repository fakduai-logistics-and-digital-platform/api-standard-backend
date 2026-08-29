package dto

type StandardReponseBody struct {
	Ok    bool `json:"ok"`
	Data  any  `json:"data,omitempty"`
	Error any  `json:"error,omitempty"`
}
