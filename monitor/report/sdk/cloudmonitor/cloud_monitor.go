package cloudmonitor

import (
	"encoding/json"

	"github.com/pkg/errors"
)

// Send data to server
func (m *Metric) Send(req *SendDataRequest) (*Response, int, error) {
	resp := new(Response)
	statusCode, err := m.SendHandler("SendMetricData", req, resp)

	if err != nil {
		return nil, statusCode, err
	}
	return resp, statusCode, nil
}

// SendHandler deal with send req
func (m *Metric) SendHandler(api string, req *SendDataRequest, resp interface{}) (int, error) {
	respBody, statusCode, err := m.Client.PostWithContentType(api, nil, req.DataLines, "text/plain")
	if err != nil {
		return statusCode, errors.Wrap(err, "send http request error")
	}
	if statusCode >= 500 {
		respBody, statusCode, err = m.Client.PostWithContentType(api, nil, req.DataLines, "text/plain")
		if err != nil {
			return statusCode, err
		}
	}

	if err := json.Unmarshal(respBody, resp); err != nil {
		return statusCode, err
	}
	return statusCode, nil
}
