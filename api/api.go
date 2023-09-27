// Package api gather generic functions to deal with the API.
package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"

	"log/slog"

	"github.com/perebaj/contractus"
)

// Error represents an error returned by the API.
type Error struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
}

func (e Error) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Msg)
}

func sendErr(w http.ResponseWriter, statusCode int, err error) {
	var httpErr Error
	if !errors.As(err, &httpErr) {
		httpErr = Error{
			Code: "unknown_error",
			Msg:  "An unexpected error happened",
		}
	}
	if statusCode >= 500 {
		slog.Error("Unable to process request", "error", err.Error(), "status_code", statusCode)
	}

	send(w, statusCode, httpErr)
}

func send(w http.ResponseWriter, statusCode int, body interface{}) {
	const jsonContentType = "application/json; charset=utf-8"

	w.Header().Set("Content-Type", jsonContentType)
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(body); err != nil {
		slog.Error("Unable to encode body as JSON", "error", err)
	}
}

func parseFile(r *http.Request) (content string, err error) {
	err = r.ParseMultipartForm(32 << 20) // 32MB
	if err != nil {
		return "", fmt.Errorf("failed to parse multipart form: %v", err)
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		return "", fmt.Errorf("failed to parse file: %v", err)
	}

	defer func() {
		_ = file.Close()
	}()

	contentByte, err := io.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %v", err)
	}
	return string(contentByte), nil
}

// convert is an internal function responsible for converting the file content into transactions.
func convert(content string) (t []contractus.Transaction, err error) {
	content = strings.Replace(content, "\t", "", -1)
	var re = regexp.MustCompile(`(?m).*$\n`)
	for _, match := range re.FindAllString(content, -1) {
		rawTransaction := Transaction{
			Type:               match[0:1],
			Date:               match[1:26],
			ProductDescription: match[26:56],
			ProductPriceCents:  match[56:66],
			SellerName:         match[66:],
		}
		transac, err := rawTransaction.Convert()
		if err != nil {
			return nil, fmt.Errorf("failed to convert transaction: %v", err)
		}
		t = append(t, *transac)

	}
	return t, nil
}
