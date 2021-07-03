package grpcV1

import (
	"crypto/tls"
	"math"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/encoding"
	//"google.golang.org/grpc/health"
	//healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

const (
	grpcOverheadBytes = 512 * 1024
	maxStreams        = math.MaxUint32
	maxSendBytes      = math.MaxInt32
)

// https://sourcegraph.com/github.com/asim/go-micro/-/blob/plugins/server/grpc/README.md
// 차후 살펴보자
func init() {
	encoding.RegisterCodec(wrapCodec{protoCodec{}})
}

// 먼저 interceptor 들 기타등등을 주석처리하고 동작 위주로 구현하고 채워나가자.
func Server(tls *tls.Config, gopts ...grpc.ServerOption) *grpc.Server {
	var opts []grpc.ServerOption

	// 여기서는 tls 가 일단 nil 로 들어감.
	// nil이 아닌경우 openssl 로 테스트 진행해보자
	if tls != nil {
		bundle := NewBundle(Config{TLSConfig: tls})
		// credentials.NewServerTLSFromCert(&cert)
		opts = append(opts, grpc.Creds(bundle.TransportCredentials()))
	}
	// interceptors
	chainUnaryInterceptors := []grpc.UnaryServerInterceptor{
		newLogUnaryInterceptor(),
		//newUnaryInterceptor(s),
		//grpc_prometheus.UnaryServerInterceptor,
	}
	/*
		chainStreamInterceptors := []grpcV1.StreamServerInterceptor{
			newStreamInterceptor(s),
			//grpc_prometheus.StreamServerInterceptor,
		}*/
	// otelgrpc 살펴봐야 함.
	/*if s.Cfg.ExperimentalEnableDistributedTracing {
		chainUnaryInterceptors = append(chainUnaryInterceptors, otelgrpc.UnaryServerInterceptor(s.Cfg.ExperimentalTracerOptions...))
		chainStreamInterceptors = append(chainStreamInterceptors, otelgrpc.StreamServerInterceptor(s.Cfg.ExperimentalTracerOptions...))

	}*/

	// grpc middleware 살펴보자.
	opts = append(opts, grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(chainUnaryInterceptors...)))
	//opts = append(opts, grpcV1.StreamInterceptor(grpc_middleware.ChainStreamServer(chainStreamInterceptors...)))

	//opts = append(opts, grpcV1.MaxRecvMsgSize(int(s.Cfg.MaxRequestBytes+grpcOverheadBytes)))
	opts = append(opts, grpc.MaxSendMsgSize(maxSendBytes))
	opts = append(opts, grpc.MaxConcurrentStreams(maxStreams))

	grpcServer := grpc.NewServer(append(opts, gopts...)...)

	//ecommerce.RegisterProductInfoServer(grpcServer, apiV1.NewProductServer())
	//pb.RegisterOrderManagementServer(grpcServer,nil)
	/*pb.RegisterKVServer(grpcServer, NewQuotaKVServer(s))
	pb.RegisterWatchServer(grpcServer, NewWatchServer(s))
	pb.RegisterLeaseServer(grpcServer, NewQuotaLeaseServer(s))
	pb.RegisterClusterServer(grpcServer, NewClusterServer(s))
	pb.RegisterAuthServer(grpcServer, NewAuthServer(s))
	pb.RegisterMaintenanceServer(grpcServer, NewMaintenanceServer(s))*/

	// server should register all the services manually
	// use empty service name for all etcd services' health status,
	// see https://github.com/grpc/grpc/blob/master/doc/health-checking.md for more
	/*hsrv := health.NewServer()
	hsrv.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)
	healthpb.RegisterHealthServer(grpcServer, hsrv)*/

	// set zero values for metrics registered for this grpcV1 server
	//grpc_prometheus.Register(grpcServer)

	return grpcServer
}
