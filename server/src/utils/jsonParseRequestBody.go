package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/labstack/echo/v5"
	"github.com/sirupsen/logrus"
)

func JsonParseRequestBody(c echo.Context, result any) error {
	r := c.Request()

	disallowUnknownFields := true

	rawBody, err := io.ReadAll(r.Body)
	if err != nil {
		err = echo.NewHTTPError(http.StatusBadRequest, "Unable to read request body")
		logrus.Warnf("JsonParseRequestBody: %s. Raw body: %s", err.Error(), string(rawBody))
		return err
	}

	r.Body = io.NopCloser(bytes.NewBuffer(rawBody)) // restore io.ReadCloser to original state

	dec := json.NewDecoder(bytes.NewBuffer(rawBody))
	if disallowUnknownFields {
		dec.DisallowUnknownFields()
	}

	if err = dec.Decode(&result); err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError

		switch {
		// Catch any syntax errors in the JSON and send an error message
		// which interpolates the location of the problem to make it
		// easier for the client to fix.
		case errors.As(err, &syntaxError):
			msg := fmt.Sprintf("request body contains badly-formed JSON (at position %d)", syntaxError.Offset)
			logrus.Warnf("JsonParseRequestBody: %s. Error: %s. Raw body: %s", msg, err.Error(), string(rawBody))
			return echo.NewHTTPError(http.StatusBadRequest, msg)

		// In some circumstances Decode() may also return an
		// io.ErrUnexpectedEOF error for syntax errors in the JSON. There
		// is an open issue regarding this at
		// https://github.com/golang/go/issues/25956.
		case errors.Is(err, io.ErrUnexpectedEOF):
			msg := "request body contains badly-formed JSON"
			logrus.Warnf("JsonParseRequestBody: %s: Error: %s. Raw body: %s", msg, err.Error(), string(rawBody))
			return echo.NewHTTPError(http.StatusBadRequest, msg)

		// Catch any type errors, like trying to assign a string in the
		// JSON request body to a int field in our Person struct. We can
		// interpolate the relevant field name and position into the error
		// message to make it easier for the client to fix.
		case errors.As(err, &unmarshalTypeError):
			msg := fmt.Sprintf("request body contains an invalid value for the %q field (at position %d).",
				unmarshalTypeError.Field, unmarshalTypeError.Offset)
			logrus.Warnf("JsonParseRequestBody: %s: Error: %s. Raw body: %s", msg, err.Error(), string(rawBody))
			return echo.NewHTTPError(http.StatusBadRequest, msg)

		// Catch the error caused by extra unexpected fields in the request
		// body. We extract the field name from the error message and
		// interpolate it in our custom error message. There is an open
		// issue at https://github.com/golang/go/issues/29035 regarding
		// turning this into a sentinel error.
		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			msg := fmt.Sprintf("request body contains unknown field %s", fieldName)
			logrus.Warnf("JsonParseRequestBody: %s: Error: %s. Raw body: %s", msg, err.Error(), string(rawBody))
			return echo.NewHTTPError(http.StatusBadRequest, msg)

		// An io.EOF error is returned by Decode() if the request body is
		// empty.
		case errors.Is(err, io.EOF):
			msg := "request body must not be empty"
			logrus.Warnf("JsonParseRequestBody: %s: Error: %s. Raw body: %s", msg, err.Error(), string(rawBody))
			return echo.NewHTTPError(http.StatusBadRequest, msg)

		// Catch the error caused by the request body being too large. Again
		// there is an open issue regarding turning this into a sentinel
		// error at https://github.com/golang/go/issues/30715.
		case err.Error() == "http: request body too large":
			msg := "Request body must not be larger than 1MB"
			logrus.Warnf("JsonParseRequestBody: %s: Error: %s. Raw body: %s", msg, err.Error(), string(rawBody))
			return echo.NewHTTPError(http.StatusRequestEntityTooLarge, msg)

		// Otherwise default to logging the error and sending a 500 Internal
		// Server Error response.
		default:
			logrus.Warnf("JsonParseRequestBody: Unknown error: %s. Raw body: %s", err, string(rawBody))
			return echo.NewHTTPError(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		}
	}

	// Call decode again, using a pointer to an empty anonymous struct as
	// the destination. If the request body only contained a single JSON
	// object this will return an io.EOF error. So if we get anything else,
	// we know that there is additional data in the request body.
	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		msg := "Request body must only contain a single JSON object"
		return echo.NewHTTPError(http.StatusBadRequest, msg)
	}

	c.Set("result", result)
	return nil
}
