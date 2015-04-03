package server

import "time"

const (
	version               = "0.1.0"     // Application and server version.
	DefaultHostname       = "localhost" // The hostname of the server.
	DefaultPort           = 49152       // Port to receive requests: see IANA Port Numbers.
	DefaultProfPort       = 0           // Profiler port to receive requests.*
	DefaultMaxConnections = 0           // Maximum number of connections allowed.*
	DefaultMaxWorkers     = 1000        // Maximum number of running workers allowed.
	DefaultMaxProcs       = 0           // Maximum number of computer processors to utilize.*

	// * zeros = no change or no limitations or not enabled.

	// Listener and connections.
	TCPConnectionTimeout = 3 * time.Minute

	// http: routes.
	httpRouteAliveV1  = "/v1.0/alive"
	httpRouteParseV1  = "/v1.0/parse"
	httpRouteStatusV1 = "/v1.0/status"

	httpGet    = "GET"
	httpPost   = "POST"
	httpPut    = "PUT"
	httpDelete = "DELETE"
	httpHead   = "HEAD"
	httpTrace  = "TRACE"
	httpPatch  = "PATCH"

	// Error messages.
	invalidMediaType     = "Invalid Content-Type or Accept header value."
	invalidBody          = "Invalid body of text in request."
	invalidJSONText      = "Invalid JSON format in text of body in request."
	invalidJSONAttribute = "Invalid - 'text' attribute in JSON not found."
)
