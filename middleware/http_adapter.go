package middleware

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/google/uuid"
	"github.com/ramoncl001/comet/log"
	"github.com/ramoncl001/comet/rest"
)

var HTTPAdapter = func(next rest.RequestHandler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		bytes, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "error parsing body", 500)
			return
		}

		request := rest.Request{
			Url:           r.URL,
			Method:        r.Method,
			QueryParams:   r.URL.Query(),
			PathParams:    make(map[string]string),
			Headers:       r.Header,
			Body:          bytes,
			UserAgent:     r.UserAgent(),
			RemoteAddress: r.RemoteAddr,
		}

		ctx := context.WithValue(r.Context(), log.TRACE_ID, uuid.New().String())

		response := next(request.WithContext(ctx))

		responseBytes, err := json.Marshal(response.Data)
		if err != nil {
			http.Error(w, "error deserializing response", 500)
			return
		}

		w.WriteHeader(response.Status)
		w.Write(responseBytes)
	})
}
