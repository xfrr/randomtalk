// semantic package provides constants for semantic conventions used in
// observability and telemetry data.
package semantic

import semconv "go.opentelemetry.io/otel/semconv/v1.26.0"

// client
const (
	ClientAddressKey = string(semconv.ClientAddressKey)
	ClientPortKey    = string(semconv.ClientPortKey)
	EnduserIDKey     = string(semconv.EnduserIDKey)
)

// http
const (
	HTTPMethodKey = string(semconv.HTTPRequestMethodKey)
	HTTPRouteKey  = string(semconv.HTTPRouteKey)

	HTTPRequestBodySizeKey    = string(semconv.HTTPRequestBodySizeKey)
	HTTPRequestSizeKey        = string(semconv.HTTPRequestSizeKey)
	HTTPRequestContentTypeKey = "http.request.content_type"

	HTTPResponseBodySizeKey    = string(semconv.HTTPResponseBodySizeKey)
	HTTPResponseSizeKey        = string(semconv.HTTPResponseSizeKey)
	HTTPResponseStatusCodeKey  = string(semconv.HTTPResponseStatusCodeKey)
	HTTPResponseContentTypeKey = "http.response.content_type"
)

// messaging
const (
	MessageIDKey   = string(semconv.MessagingMessageIDKey)
	MessageTypeKey = string(semconv.MessagingOperationTypeKey)
)
