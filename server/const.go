package server

const (
	version               = "0.1.0" // Application and server version.
	DefaultPort           = 49152   // Port to receive requests: see IANA Port Numbers.
	DefaultMaxConnections = 1000    // Maximum number of concurrent connections allowed.
	DefaultMaxProcs       = 0       // Maximum number of computer processors to utilize.*

	// * zeros = no change

	// http: routes

	httpRouteAliveV1  = "/v1.0/alive"
	httpRouteParseV1  = "/v1.0/parse"
	httpRouteStatusV1 = "/v1.0/status"

	httpGet = "GET"

	// error messages

	invalidMediaType     = "Invalid Content-Type or Accept header value."
	invalidBody          = "Invalid body of text in request."
	invalidJSONText      = "Invalid JSON format in text of body in request."
	invalidJSONAttribute = "Invalid - 'text' attribute in JSON not found."
)
