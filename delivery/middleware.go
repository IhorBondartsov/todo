package delivery

import (
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	log "github.com/sirupsen/logrus"
)

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func NewLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{w, http.StatusOK}
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func logMiddleware(h httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		t1 := time.Now()

		lrw := NewLoggingResponseWriter(w)
		h(lrw, r, ps)

		timeProcessing := time.Since(t1)

		logger := log.WithFields(
			log.Fields{
				"duration": timeProcessing,
				"method":   r.Method,
				"path":     r.URL.Path,
			})

		switch lrw.statusCode {
		case http.StatusOK:
			logger.Info("Success")
		case http.StatusBadRequest:
			logger.Info("bad request")
		case http.StatusNotFound:
			logger.Info("Not Found")
		case http.StatusInternalServerError:
			logger.Error("Return an error")
		default:
			logger.Warning("unhandle http status code")
		}
	}
}
