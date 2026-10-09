package api

type tcpingRes struct {
	// Status is "true" / "false" for a probe verdict and "error" when no verdict was reached
	Status string `json:"status"`
	// Time is the connect latency in milliseconds, only present when the connection succeeded
	Time    *int64 `json:"time,omitempty"`
	Message string `json:"message"`
}
