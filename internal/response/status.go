package response

type StatusCode int

const (
	StatusOK        = 200
	StatusCreated   = 201
	StatusAccepted  = 202
	StatusNoContent = 204

	StatusMovedPermanently  = 301
	StatusFound             = 302
	StatusNotModified       = 304
	StatusTemporaryRedirect = 307
	StatusPermanentRedirect = 308

	StatusBadRequest            = 400
	StatusUnauthorized          = 401
	StatusForbidden             = 403
	StatusNotFound              = 404
	StatusMethodNotAllowed      = 405
	StatusRequestTimeout        = 408
	StatusConflict              = 409
	StatusRequestEntityTooLarge = 413
	StatusTooManyRequests       = 429

	StatusInternalServerError = 500
	StatusNotImplemented      = 501
	StatusBadGateway          = 502
	StatusServiceUnavailable  = 503
	StatusGatewayTimeout      = 504
)

var statusMessage = map[StatusCode]string{
	StatusOK:        "OK",
	StatusCreated:   "Created",
	StatusAccepted:  "Accepted",
	StatusNoContent: "No Content",

	StatusMovedPermanently:  "Moved Permanently",
	StatusFound:             "Found",
	StatusNotModified:       "Not Modified",
	StatusTemporaryRedirect: "Temporary Redirect",
	StatusPermanentRedirect: "Permanent Redirect",

	StatusBadRequest:            "Bad Request",
	StatusUnauthorized:          "Unauthorized",
	StatusForbidden:             "Forbidden",
	StatusNotFound:              "Not Found",
	StatusMethodNotAllowed:      "Method Not Allowed",
	StatusRequestTimeout:        "Request Timeout",
	StatusConflict:              "Conflict",
	StatusRequestEntityTooLarge: "Payload Too Large",
	StatusTooManyRequests:       "Too Many Requests",

	StatusInternalServerError: "Internal Server Error",
	StatusNotImplemented:      "Not Implemented",
	StatusBadGateway:          "Bad Gateway",
	StatusServiceUnavailable:  "Service Unavailable",
	StatusGatewayTimeout:      "Gateway Timeout",
}
