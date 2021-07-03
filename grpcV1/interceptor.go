package grpcV1

import (
	"context"
	"os"
	"time"

	"github.com/onrik/logrus/filename"
	//pb "github.com/seoyhaein/caleb-middle/protos"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	//"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

const (
	maxNoLeaderCnt          = 3
	warnUnaryRequestLatency = 300 * time.Millisecond
	snapshotMethod          = "/etcdserverpb.Maintenance/Snapshot"
)

var (
	log *logrus.Logger = logrus.New()
)

func init() {

	// 특정위치에 저장하는 코드 추가해줘야 함. 일단은
	_, err := os.OpenFile("lo.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		//log.Out = file
		log.Out = os.Stdout
	} else {
		log.Info("Failed to log to file, using default stderr")
	}

	filenameHook := filename.NewHook()
	//filenameHook.Field = "custom_source_field" // Customize source field name
	log.AddHook(filenameHook)

	//log.Println("something")

}

func newLogUnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		startTime := time.Now()
		// 전처리
		resp, err := handler(ctx, req)

		//logUnaryRequestStats(ctx, info, startTime, req, resp)
		//후처리
		defer logUnaryRequestStats(ctx, info, startTime, req, resp)
		log.Info("newLogUnaryInterceptor")

		return resp, err
	}
}

func logUnaryRequestStats(ctx context.Context, info *grpc.UnaryServerInfo, startTime time.Time, req interface{}, resp interface{}) {
	// 시간이 너무 많이 지체하면 지체한 api 에 대한 로그를 기록한다.
	// 시간이 너무 많이 지체하지 않으면 로그를 기록하지 않는다.
	// 테스트를 위해 일단 주석처리
	duration := time.Since(startTime)
	//var enabledDebugLevel, expensiveRequest bool

	/*if duration > warnUnaryRequestLatency {
		expensiveRequest = true
	}*/
	/*log.Info("logUnaryRequestStats")
	if !enabledDebugLevel && !expensiveRequest {
		log.Info("logUnaryRequestStatsdfdfdfdfdf")
		return
	}*/

	remote := "No remote client info."
	// 이 함수 알아봐야함.
	peerInfo, ok := peer.FromContext(ctx)
	if ok {
		remote = peerInfo.Addr.String()
	}

	responseType := info.FullMethod
	var reqCount, respCount int64
	var reqSize, respSize int
	var reqContent string
	log.Info("1111111")
	/*switch _resp := resp.(type) {
	case *pb.Product:
		_req, ok := req.(*pb.Product)
		log.Info("222222222")
		if ok {
			//reqCount = 0
			log.Info("33333333")
			reqSize = _req.Size()
			//reqContent = _req.String()
		}
		log.Info("4444444444")
		if _resp != nil {
			log.Info("55555555")
			//respCount = _resp.GetCount()
			respSize = _resp.Size()
		}
	case *pb.ProductID:
		log.Info("666666666")
		_req, ok := req.(*pb.ProductID)
		if ok {
			log.Info("777777777")
			//reqCount = 1
			reqSize = _req.Size()
			//reqContent = pb.NewLoggablePutRequest(_req).String()
			// redact value field from request content, see PR #9821
		}
		log.Info("8888888888888")
		if _resp != nil {
			//respCount = 0
			log.Info("99999999999")
			respSize = _resp.Size()
		}
	default:
		log.Info("AAAAAAAAAAA")
		//reqCount = -1
		reqSize = -1
		//respCount = -1
		respSize = -1
	}*/
	log.Info("YY")
	logGenericRequestStats(startTime, duration, remote, responseType, reqCount, reqSize, respCount, respSize, reqContent)

	/*if enabledDebugLevel {
		logGenericRequestStats(lg, startTime, duration, remote, responseType, reqCount, reqSize, respCount, respSize, reqContent)
	} else if expensiveRequest {
		logExpensiveRequestStats(lg, startTime, duration, remote, responseType, reqCount, reqSize, respCount, respSize, reqContent)
	}*/
}

func logGenericRequestStats(startTime time.Time, duration time.Duration, remote string, responseType string,
	reqCount int64, reqSize int, respCount int64, respSize int, reqContent string) {

	log.Println("XX")

	//fmt.Println("일단 이렇게")

	log.WithFields(logrus.Fields{
		"start time":      startTime,
		"time spent":      duration,
		"remote":          remote,
		"response type":   responseType,
		"request count":   reqCount,
		"request size":    reqSize,
		"response count":  respCount,
		"response size":   respSize,
		"request content": reqContent,
	}).Info("request stats")

	// 로그 기록을 남겨야 한다.
	/*lg.Debug("request stats",
		zap.Time("start time", startTime),
		zap.Duration("time spent", duration),
		zap.String("remote", remote),
		zap.String("response type", responseType),
		zap.Int64("request count", reqCount),
		zap.Int("request size", reqSize),
		zap.Int64("response count", respCount),
		zap.Int("response size", respSize),
		zap.String("request content", reqContent),
	)*/
}

/*func newStreamInterceptor(s *etcdserver.EtcdServer) grpc.StreamServerInterceptor {
	smap := monitorLeader(s)

	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		if !api.IsCapabilityEnabled(api.V3rpcCapability) {
			return rpctypes.ErrGRPCNotCapable
		}

		if s.IsMemberExist(s.ID()) && s.IsLearner() && info.FullMethod != snapshotMethod { // learner does not support stream RPC except Snapshot
			return rpctypes.ErrGPRCNotSupportedForLearner
		}

		md, ok := metadata.FromIncomingContext(ss.Context())
		if ok {
			ver, vs := "unknown", md.Get(rpctypes.MetadataClientAPIVersionKey)
			if len(vs) > 0 {
				ver = vs[0]
			}
			clientRequests.WithLabelValues("stream", ver).Inc()

			if ks := md[rpctypes.MetadataRequireLeaderKey]; len(ks) > 0 && ks[0] == rpctypes.MetadataHasLeader {
				if s.Leader() == types.ID(raft.None) {
					return rpctypes.ErrGRPCNoLeader
				}

				ctx := newCancellableContext(ss.Context())
				ss = serverStreamWithCtx{ctx: ctx, ServerStream: ss}

				smap.mu.Lock()
				smap.streams[ss] = struct{}{}
				smap.mu.Unlock()

				defer func() {
					smap.mu.Lock()
					delete(smap.streams, ss)
					smap.mu.Unlock()
					// TODO: investigate whether the reason for cancellation here is useful to know
					ctx.Cancel(nil)
				}()
			}
		}

		return handler(srv, ss)
	}
}*/
