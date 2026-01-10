package api

type tcpingRes struct {
	Status  string `json:"status"`
	Time    int    `json:"time,omitempty"`
	Message string `json:"message"`
}
