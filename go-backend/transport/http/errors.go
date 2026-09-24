package transporthttp

import (
	"errors"
	"log/slog"
	"net/http"
)

type Error struct {
	Status int
	Msg    string
	Cause  error
}

func (e *Error) Error() string { return e.Msg }
func (e *Error) Unwrap() error { return e.Cause }

func writeError(w http.ResponseWriter, r *http.Request, err error) {
	var herr *Error
	if errors.As(err, &herr) {
		if herr.Cause != nil {
			slog.Error(r.Pattern+": request failed", "error", herr.Cause)
		}
		http.Error(w, herr.Msg, herr.Status)
		return
	}

	slog.Error(r.Pattern+": unhandled error", "error", err)
	http.Error(w, "internal error", http.StatusInternalServerError)
}
