package cloudmonitor

import "github.com/volcengine/volc-sdk-golang/base"

// SendDataRequest .
type SendDataRequest struct {
	DataLines string
}

// Response .
type Response struct {
	ResponseMetadata base.ResponseMetadata
	Result           *Result `json:"Result,omitempty"`
}

// Result .
type Result struct {
	Code string `json:"Code,omitempty"`
}
