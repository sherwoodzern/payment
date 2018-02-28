package main


import (
	"flag"
	"fmt"
	//"net"
	"net/http"
	"os"
	"os/signal"
	//"strings"
	"syscall"

	//"github.com/go-kit/kit/log"
	"github.com/microservices-demo/payment"
	//stdopentracing "github.com/opentracing/opentracing-go"
	//zipkin "github.com/openzipkin/zipkin-go-opentracing"
	"golang.org/x/net/context"

)

const (
	ServiceName = "payment"
)

func main() {
	var (
		port          = flag.String("port", "8080", "Port to bind HTTP listener")
		//zip           = flag.String("zipkin", os.Getenv("ZIPKIN"), "Zipkin address")
		declineAmount = flag.Float64("decline", 105, "Decline payments over certain amount")
	)
	flag.Parse()

	fmt.Println("The Port Value: " + *port)
	// The istio service mesh provides openTracing at the time the service is deployed
	/* var tracer stdopentracing.Tracer
	{
		// Log domain.
		var logger log.Logger
		{
			logger = log.NewLogfmtLogger(os.Stderr)
			logger = log.With(logger, "ts", log.DefaultTimestampUTC)
			logger = log.With(logger, "caller", log.DefaultCaller)
		}
		// Find service local IP.
		conn, err := net.Dial("udp", "8.8.8.8:80")
		if err != nil {
			logger.Log("err", err)
			os.Exit(1)
		}
		localAddr := conn.LocalAddr().(*net.UDPAddr)
		host := strings.Split(localAddr.String(), ":")[0]
		defer conn.Close()
		if *zip == "" {
			tracer = stdopentracing.NoopTracer{}
		} else {
			logger := log.With(	logger,"tracer", "Zipkin")

			logger.Log("addr", zip)
			collector, err := zipkin.NewHTTPCollector(
				*zip,
				zipkin.HTTPLogger(logger),
			)
			if err != nil {
				logger.Log("err", err)
				os.Exit(1)
			}
			tracer, err = zipkin.NewTracer(
				zipkin.NewRecorder(collector, false, fmt.Sprintf("%v:%v", host, port), ServiceName),
			)
			if err != nil {
				logger.Log("err", err)
				os.Exit(1)
			}
		}
		stdopentracing.InitGlobalTracer(tracer)

	} */
	// Mechanical stuff.
	errc := make(chan error)
	ctx := context.Background()  //return a non-nil, empty Context

	// ? Do we need the tracer ?
	//handler, logger := payment.WireUp(ctx, float32(*declineAmount), tracer, ServiceName)

	//The tracer was removed
	handler, logger := payment.WireUp(ctx,float32(*declineAmount),ServiceName)

	// Create and launch the HTTP server.
	go func() {
		logger.Log("transport", "HTTP", "port", "8080")
		logger.Log("transport", "HTTP", "port-pointer", *port)
		//errc <- http.ListenAndServe(":"+*port, handler)
		errc <- http.ListenAndServe(":"+*port, handler)
	}()

	// Capture interrupts.
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errc <- fmt.Errorf("%s", <-c)
	}()

	logger.Log("exit", <-errc)
}
