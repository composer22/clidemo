package server

const (
	version               = "0.1.0" // Application and server version.
	DefaultPort           = 49152   // Port to receive requests from -- see IANA Private/Dynamic Port Numbers.
	DefaultMaxConnections = 1000    // Maximum number of concurrent connections allowed.
	DefaultMaxProcs       = 0       // Maximum number of computer processors to utilize in the application and/or server. 0 = no change.

	// http: routes

	httpRouteAliveV1  = "/v1.0/alive/"
	httpRouteParseV1  = "/v1.0/parse/"
	httpRouteStatusV1 = "/v1.0/status/"
	httpGet           = "GET"
)
