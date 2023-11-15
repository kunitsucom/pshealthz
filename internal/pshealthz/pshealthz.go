package pshealthz

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"regexp"

	errorz "github.com/kunitsucom/util.go/errors"
)

//nolint:gochecknoglobals
var (
	procPath     = os.Getenv("PSHEALTHZ_PROC_EMULATOR_PATH_PREFIX") + "/proc"
	regexNumbers = regexp.MustCompile(`^\d+$`)
)

type Request struct {
	Regex string `json:"regex"`
}

type Response struct {
	OK        bool        `json:"ok"`
	Message   *string     `json:"message,omitempty"`
	Processes *[]*Process `json:"processes,omitempty"`
}

type Process struct {
	PID     string `json:"pid"`
	Cmdline string `json:"cmdline"`
}

func NewErrorResponse(message string) string {
	res := &Response{
		OK:      false,
		Message: &message,
	}

	resBytes, err := json.Marshal(res)
	if err != nil {
		err = errorz.Errorf("json.Marshal: %w", err)
		panic(err)
	}

	return string(resBytes)
}

//nolint:cyclop
func PSHealthz(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	req := &Request{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		http.Error(w, NewErrorResponse("invalid json request"), http.StatusBadRequest)
		return
	}

	psRegex, err := regexp.Compile(req.Regex)
	if err != nil {
		http.Error(w, NewErrorResponse("invalid regex"), http.StatusBadRequest)
		return
	}

	entries, err := os.ReadDir(procPath)
	if err != nil {
		http.Error(w, NewErrorResponse("failed to read /proc"), http.StatusInternalServerError)
		return
	}

	res := &Response{
		OK:        true,
		Processes: &[]*Process{},
	}
	for _, entry := range entries {
		if entry.IsDir() && regexNumbers.MatchString(entry.Name()) {
			cmdline, err := os.ReadFile(filepath.Join(procPath, entry.Name(), "cmdline"))
			if err != nil {
				continue
			}

			sanitized := bytes.TrimSuffix(bytes.ReplaceAll(cmdline, []byte{0}, []byte(" ")), []byte(" "))

			if !psRegex.Match(sanitized) {
				continue
			}

			*res.Processes = append(*res.Processes, &Process{
				PID:     entry.Name(),
				Cmdline: string(sanitized),
			})
		}
	}

	if len(*res.Processes) != 0 {
		resBytes, err := json.Marshal(res)
		if err != nil {
			http.Error(w, NewErrorResponse("failed to marshal response"), http.StatusInternalServerError)
			return
		}

		if _, err := w.Write(resBytes); err != nil {
			http.Error(w, NewErrorResponse("failed to write response"), http.StatusInternalServerError)
			return
		}
		return
	}

	http.Error(w, NewErrorResponse("no processes found"), http.StatusNotFound)
	return //nolint:gosimple
}
