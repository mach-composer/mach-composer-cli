package utils

import (
	"log"
	"net/http"
	"net/http/httputil"
	"os"
)

type LogTransport struct {
	transport http.RoundTripper
}

var DebugTransport = &LogTransport{
	transport: http.DefaultTransport,
}

func (c *LogTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	if os.Getenv("MCC_DEBUG") != "" {
		logRequest(request)
	}
	response, err := c.transport.RoundTrip(request)
	if os.Getenv("MCC_DEBUG") != "" {
		logResponse(response, err)
	}
	return response, err
}

const logRequestTemplate = `DEBUG:
---[ REQUEST ]--------------------------------------------------------
%s
----------------------------------------------------------------------
`

const logResponseTemplate = `DEBUG:
---[ RESPONSE ]-------------------------------------------------------
%s
----------------------------------------------------------------------
`

func logRequest(r *http.Request) {
	body, err := httputil.DumpRequestOut(r, true)
	if err != nil {
		return
	}
	log.Printf(logRequestTemplate, body)
}

func logResponse(r *http.Response, err error) {
	if err != nil {
		log.Printf(logResponseTemplate, err)
		return
	}
	body, err := httputil.DumpResponse(r, true)
	if err != nil {
		return
	}
	log.Printf(logResponseTemplate, body)
}
