package start

import (
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/infraboard/mcube/app"
	"github.com/spf13/cobra"

	// 注册所有服务
	_ "hzh/devcloud/mpaas/apps"
	"hzh/devcloud/mpaas/common/logger"
	"hzh/devcloud/mpaas/conf"
	"hzh/devcloud/mpaas/protocol"
)

// startCmd represents the start command
var Cmd = &cobra.Command{
	Use:   "start",
	Short: "mcenter API服务",
	Long:  "mcenter API服务",
	RunE: func(cmd *cobra.Command, args []string) error {
		conf := conf.C()
		// 启动服务
		ch := make(chan os.Signal, 1)
		defer close(ch)
		signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP, syscall.SIGQUIT)

		// 初始化服务
		svr, err := newService(conf)
		if err != nil {
			return err
		}

		// 等待信号处理
		go svr.waitSign(ch)

		// 启动服务
		if err := svr.start(); err != nil {
			if !strings.Contains(err.Error(), "http: Server closed") {
				return err
			}
		}

		return nil
	},
}

func newService(cnf *conf.Config) (*service, error) {
	http := protocol.NewHTTPService()
	grpc := protocol.NewGRPCService()
	svr := &service{
		http: http,
		grpc: grpc,
	}

	return svr, nil
}

type service struct {
	http *protocol.HTTPService
	grpc *protocol.GRPCService
}

func (s *service) start() error {
	logger.L().Info().Msgf("loaded grpc app: %s", app.LoadedGrpcApp())
	logger.L().Info().Msgf("loaded http app: %s", app.LoadedRESTfulApp())

	logger.L().Info().Msgf("loaded internal app: %s", app.LoadedInternalApp())

	go s.grpc.Start()
	return s.http.Start()
}

func (s *service) waitSign(sign chan os.Signal) {
	for sg := range sign {
		switch v := sg.(type) {
		default:
			logger.L().Info().Msgf("receive signal '%v', start graceful shutdown", v.String())

			if err := s.grpc.Stop(); err != nil {
				logger.L().Info().Msgf("grpc graceful shutdown err: %s, force exit", err)
			} else {
				logger.L().Info().Msgf("grpc service stop complete")
			}

			if err := s.http.Stop(); err != nil {
				logger.L().Info().Msgf("http graceful shutdown err: %s, force exit", err)
			} else {
				logger.L().Info().Msgf("http service stop complete")
			}
			return
		}
	}
}
