package protocol

import (
	"net"

	"hzh/devcloud/mpaas/conf"

	"google.golang.org/grpc"

	"hzh/devcloud/mpaas/common/logger"

	"github.com/infraboard/mcube/app"
	"github.com/infraboard/mcube/grpc/middleware/recovery"
)

// NewGRPCService todo
func NewGRPCService() *GRPCService {
	rc := recovery.NewInterceptor(recovery.NewZapRecoveryHandler())
	grpcServer := grpc.NewServer(grpc.ChainUnaryInterceptor(
		rc.UnaryServerInterceptor(),
		// 加载gRPC服务端认证中间件
		// auth.NewGrpcAuther().AuthFunc,
	))

	return &GRPCService{
		svr: grpcServer,
		c:   conf.C(),
	}
}

// GRPCService grpc服务
type GRPCService struct {
	svr *grpc.Server
	c   *conf.Config
}

// Start 启动GRPC服务
func (s *GRPCService) Start() {
	// 装载所有GRPC服务
	app.LoadGrpcApp(s.svr)

	// 启动HTTP服务
	lis, err := net.Listen("tcp", s.c.App.GRPC.Addr())
	if err != nil {
		logger.L().Debug().Msgf("listen grpc tcp conn error, %s", err)
		return
	}

	logger.L().Info().Msgf("GRPC 服务监听地址: %s", s.c.App.GRPC.Addr())
	if err := s.svr.Serve(lis); err != nil {
		if err == grpc.ErrServerStopped {
			logger.L().Info().Msg("service is stopped")
		}

		logger.L().Error().Msgf("start grpc service error, %s", err.Error())
		return
	}
}

// Stop GRPC服务
func (s *GRPCService) Stop() error {
	s.svr.GracefulStop()
	return nil
}
