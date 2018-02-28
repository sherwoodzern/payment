package payment

import (
	"github.com/go-kit/kit/endpoint"
	//"github.com/go-kit/kit/tracing/opentracing"
	//stdopentracing "github.com/opentracing/opentracing-go"
	"golang.org/x/net/context"
)

// Endpoints collects the endpoints that comprise the Service.
type Endpoints struct {
	AuthoriseEndpoint endpoint.Endpoint
	HealthEndpoint    endpoint.Endpoint
}

// MakeEndpoints returns an Endpoints structure, where each endpoint is
// backed by the given service.
/*func MakeEndpoints(s Service, tracer stdopentracing.Tracer) Endpoints {
	return Endpoints{
		AuthoriseEndpoint: opentracing.TraceServer(tracer, "POST /paymentAuth")(MakeAuthoriseEndpoint(s)),
		HealthEndpoint:    opentracing.TraceServer(tracer, "GET /health")(MakeHealthEndpoint(s)),
	}
} */


func MakeEndpoints(s Service) Endpoints {
	return Endpoints{
		AuthoriseEndpoint: MakeAuthoriseEndpoint(s),
		HealthEndpoint:    MakeHealthEndpoint(s),
	}
}

// MakeListEndpoint returns an endpoint via the given service.

//		var span stdopentracing.Span
//		span, ctx = stdopentracing.StartSpanFromContext(ctx, "authorize payment")
//		span.SetTag("service", "payment")
//		defer span.Finish()

func MakeAuthoriseEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(AuthoriseRequest)
		authorisation, err := s.Authorise(req.Amount)
		return AuthoriseResponse{Authorisation: authorisation, Err: err}, err
	}
}

// MakeHealthEndpoint returns current health of the given service.

//		var span stdopentracing.Span
//		span, ctx = stdopentracing.StartSpanFromContext(ctx, "health check")
//		span.SetTag("service", "payment")
//		defer span.Finish()

func MakeHealthEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		health := s.Health()
		return healthResponse{Health: health}, err
	}
}

// AuthoriseRequest represents a request for payment authorisation.
// The Amount is the total amount of the transaction
type AuthoriseRequest struct {
	Amount float32 `json:"amount"`
}

// AuthoriseResponse returns a response of type Authorisation and an error, Err.
type AuthoriseResponse struct {
	Authorisation Authorisation
	Err           error
}

type healthRequest struct {
	//
}

type healthResponse struct {
	Health []Health `json:"health"`
}
