package handler

type response struct {
	Status string `json:"status"`
}

type key string

func (k key) String() string {
	return string(k)
}
