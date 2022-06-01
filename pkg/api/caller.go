package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"go.uber.org/zap"
)

type CallerData struct {
	AppId      string `json:"appId"`
	Method     string `json:"method"`
	HTTPMethod string `json:"httpMethod"`
	Payload    string `json:"payload"`
}

func (s *Server) callMethod(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		s.logger.Infow("HTTP method not allowed on callMethod",
			zap.String("error", r.Method))
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		s.logger.Infow("Could not call method: error reading request body",
			zap.Error(err))
		return
	}

	// unmarshal the request body, expecting a JSON object with AppId, Method and Payload
	var callerData CallerData
	err = json.Unmarshal(body, &callerData)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		s.logger.Infow("Could not call method: invalid request body",
			zap.Error(err))
		return
	}

	// return error if AppId, Method or PayLoad is empty
	if callerData.AppId == "" || callerData.Method == "" || callerData.Payload == "" {
		w.WriteHeader(http.StatusBadRequest)
		s.logger.Infow("Could not call method: missing AppId, Method or Payload",
			zap.String("AppId", callerData.AppId),
			zap.String("Method", callerData.Method),
			zap.String("Payload", callerData.Payload))
		return
	}

	// call the method
	var jsonMap map[string]interface{}
	json.Unmarshal([]byte(callerData.Payload), &jsonMap)
	ctx := r.Context()
	// call Dapr method with content
	_, err = s.daprClient.InvokeMethodWithCustomContent(ctx, callerData.AppId, callerData.Method, callerData.HTTPMethod, "application/json", jsonMap)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		s.logger.Infow("Could not call method: invoke with custom content failed",
			zap.Error(err))
		return
	}

	w.WriteHeader(http.StatusOK)
	s.logger.Infow("Called method with custom content", "method", callerData.Method, "appId", callerData.AppId)

}
