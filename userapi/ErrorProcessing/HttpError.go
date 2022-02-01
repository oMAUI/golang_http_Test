package ErrorProcessing

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
)

type CustomError struct {
	Message string `json:"message"`
}

func HttpError(w http.ResponseWriter, err error, msgForLogger string, msgForResponse string, code int) {
	w.Header().Set("Content-Type", "application/json")

	customError := CustomError{
		Message: msgForResponse,
	}

	zap.S().Errorw(msgForLogger, "error", err)
	res, errGetJson := json.Marshal(customError)
	if errGetJson != nil {
		zap.S().Errorw("marshal", "error", errGetJson)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("server error"))
		return
	}

	w.WriteHeader(code)
	w.Write(res)
}
