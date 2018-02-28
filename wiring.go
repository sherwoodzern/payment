package payment

import (
	"net/http"
	"os"
	goLogger "log"

	"github.com/go-kit/kit/log"
	"golang.org/x/net/context"

	//stdopentracing "github.com/opentracing/opentracing-go"
	"github.com/prometheus/client_golang/prometheus"
	//"github.com/weaveworks/common/middleware"

)

var (
	HTTPLatency = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "request_duration_seconds",
		Help:    "Time (in seconds) spent serving HTTP requests.",
		Buckets: prometheus.DefBuckets,
	}, []string{"method", "route", "status_code", "isWS"})
)


func init() {
	prometheus.MustRegister(HTTPLatency)
}

//func WireUp(ctx context.Context, declineAmount float32, tracer stdopentracing.Tracer, serviceName string) (http.Handler, log.Logger) {
func WireUp(ctx context.Context, declineAmount float32, serviceName string) (http.Handler, log.Logger) {
	// Log domain.



	var golog = goLogger.New(os.Stdout,"payment.wiring: ", goLogger.Lshortfile)

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger,"ts", log.DefaultTimestampUTC)
		logger = log.With(logger,"caller", log.DefaultCaller)
	}

	// Service domain.
	var service Service
	{
		service = NewAuthorisationService(declineAmount)
		service = LoggingMiddleware(logger)(service)
	}

	// Endpoint domain.
	//endpoints := MakeEndpoints(service, tracer)
	endpoints := MakeEndpoints(service)

	//router := MakeHTTPHandler(ctx, endpoints, logger, tracer)
	router := MakeHTTPHandler(ctx, endpoints, logger)

	/*httpMiddleware := []middleware.Interface{
		middleware.Instrument{
			Duration:     HTTPLatency,
			RouteMatcher: router,
		},
	} */

	// Handler
	//handler := middleware.Merge(httpMiddleware...).Wrap(router)


	golog.Print("Printing the handler")
	//fmt.Println(handler)


	return router, logger
	//return logger
}
