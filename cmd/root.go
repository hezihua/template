package cmd

import (
	"errors"
	"fmt"

	"hzh/devcloud/mpaas/cmd/start"
	"hzh/devcloud/mpaas/conf"
	"hzh/devcloud/mpaas/version"

	"github.com/infraboard/mcube/app"
	"github.com/spf13/cobra"
)

var (
	// pusher service config option
	confType string
	confFile string
	confETCD string
)

var vers bool

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "mcenter",
	Short: "用户中心",
	Long:  "用户中心",
	RunE: func(cmd *cobra.Command, args []string) error {
		if vers {
			fmt.Println(version.FullVersion())
			return nil
		}
		return cmd.Help()
	},
}

func initail() {
	// 初始化全局变量
	err := loadGlobalConfig(confType)
	cobra.CheckErr(err)

	// 初始化全局app
	err = app.InitAllApp()
	cobra.CheckErr(err)
}

// config 为全局变量, 只需要load 即可全局可用户
func loadGlobalConfig(configType string) error {
	// 配置加载
	switch configType {
	case "file":
		err := conf.LoadConfigFromToml(confFile)
		if err != nil {
			return err
		}
	case "env":
		err := conf.LoadConfigFromEnv()
		if err != nil {
			return err
		}
	default:
		return errors.New("unknown config type")
	}

	return nil
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	// 初始化设置
	cobra.OnInitialize(initail)
	RootCmd.AddCommand(start.Cmd)
	// RootCmd.AddCommand(initial.Cmd)
	err := RootCmd.Execute()
	cobra.CheckErr(err)
}

func init() {
	RootCmd.PersistentFlags().StringVarP(&confType, "config-type", "t", "file", "the service config type [file/env/etcd]")
	RootCmd.PersistentFlags().StringVarP(&confFile, "config-file", "f", "etc/config.toml", "the service config from file")
	RootCmd.PersistentFlags().StringVarP(&confETCD, "config-etcd", "e", "127.0.0.1:2379", "the service config from etcd")
	RootCmd.PersistentFlags().BoolVarP(&vers, "version", "v", false, "the mcenter version")
}
