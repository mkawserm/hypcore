package views

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

func httpGetHeader(h http.Header, key string) string {
	if h == nil {
		return ""
	}
	v := h[key]
	if len(v) == 0 {
		return ""
	}
	return v[0]
}

// httpJsonError is like the http.Error with WebSocket context exception.
func httpJsonError(w http.ResponseWriter, body string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", strconv.Itoa(len(body)))
	w.WriteHeader(code)

	_, _ = w.Write([]byte(body))
}

func httpBadRequest(w http.ResponseWriter, msg []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	output := fmt.Sprintf("{\"message\":\"%s\", \"error_code\":\"400\"}", msg)
	w.Header().Set("Content-Length", strconv.Itoa(len(output)))

	_, _ = w.Write([]byte(output))
}

func httpNotFound(w http.ResponseWriter, msg []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	output := fmt.Sprintf("{\"message\":\"%s\", \"error_code\":\"404\"}", msg)
	w.Header().Set("Content-Length", strconv.Itoa(len(output)))

	_, _ = w.Write([]byte(output))
}

func httpMessage(w http.ResponseWriter, msg []byte, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	output := fmt.Sprintf("{\"message\":\"%s\", \"code\":\"%d\"}", msg, code)

	w.Header().Set("Content-Length", strconv.Itoa(len(output)))

	_, _ = w.Write([]byte(output))
}

func GraphQLErrorMessage(w http.ResponseWriter, msg []byte, error_code string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	messageFormat := `
{
	"error": {
				"message": "%s",
				"code": "%s"
			 }
}
`

	output := fmt.Sprintf(messageFormat, msg, error_code)

	w.Header().Set("Content-Length", strconv.Itoa(len(output)))

	_, _ = w.Write([]byte(output))
}

func GraphQLSmartErrorMessage(w http.ResponseWriter, msg interface{}, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	messageFormat := `{"data":null,"error":%s}`
	message, err := json.Marshal(msg)

	if err == nil {
		output := fmt.Sprintf(messageFormat, message)
		w.Header().Set("Content-Length", strconv.Itoa(len(output)))
		_, _ = w.Write([]byte(output))
	} else {
		output := `{"data":null,"error":""}`
		w.Header().Set("Content-Length", strconv.Itoa(len(output)))
		_, _ = w.Write([]byte(output))
	}
}
