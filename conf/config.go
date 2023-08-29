package conf

import (
	"context"
	"fmt"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func newConfig() *Config {
	return &Config{
		App:   newDefaultAPP(),
		Mongo: newDefaultMongoDB(),
	}
}

type Config struct {
	App   *app     `toml:"app"`
	Mongo *mongodb `toml:"mongodb"`
}

type app struct {
	Name string `toml:"name" env:"APP_NAME"`
	HTTP *http  `toml:"http"`
	GRPC *grpc  `toml:"grpc"`
}

func newDefaultAPP() *app {
	return &app{
		Name: "mcenter",
		HTTP: newDefaultHTTP(),
		GRPC: newDefaultGRPC(),
	}
}

type http struct {
	Host string `toml:"host" env:"HTTP_HOST"`
	Port string `toml:"port" env:"HTTP_PORT"`
}

func (a *http) Addr() string {
	return a.Host + ":" + a.Port
}

func newDefaultHTTP() *http {
	return &http{
		Host: "127.0.0.1",
		Port: "8020",
	}
}

type grpc struct {
	Host string `toml:"host" env:"GRPC_HOST"`
	Port string `toml:"port" env:"GRPC_PORT"`
}

func (a *grpc) Addr() string {
	return a.Host + ":" + a.Port
}

func newDefaultGRPC() *grpc {
	return &grpc{
		Host: "127.0.0.1",
		Port: "18020",
	}
}
func newDefaultMongoDB() *mongodb {
	return &mongodb{
		Database:   "mpaas",
		AuthSource: "mpaas",
		Endpoints:  []string{"127.0.0.1:27017"},
	}
}

type mongodb struct {
	Endpoints  []string `toml:"endpoints" env:"MONGO_ENDPOINTS" envSeparator:","`
	UserName   string   `toml:"username" env:"MONGO_USERNAME"`
	Password   string   `toml:"password" env:"MONGO_PASSWORD"`
	Database   string   `toml:"database" env:"MONGO_DATABASE"`
	AuthSource string   `toml:"auth_source" env:"MONGO_AUTH_SOURCE"`

	lock   sync.Mutex
	client *mongo.Client
}

// Client 获取一个全局的mongodb客户端连接
func (m *mongodb) Client() (*mongo.Client, error) {
	// 加载全局数据量单例
	m.lock.Lock()
	defer m.lock.Unlock()

	if m.client == nil {
		conn, err := m.getClient()
		if err != nil {
			return nil, err
		}
		m.client = conn
	}

	return m.client, nil
}

func (m *mongodb) GetDB() (*mongo.Database, error) {
	conn, err := m.Client()
	if err != nil {
		return nil, err
	}
	return conn.Database(m.Database), nil
}

// AuthSource: 用于认证的数据, 使用的Database
func (m *mongodb) getClient() (*mongo.Client, error) {
	opts := options.Client()

	cred := options.Credential{
		AuthSource: m.AuthSource,
	}

	if m.UserName != "" && m.Password != "" {
		cred.Username = m.UserName
		cred.Password = m.Password
		cred.PasswordSet = true
		opts.SetAuth(cred)
	}
	opts.SetHosts(m.Endpoints)
	opts.SetConnectTimeout(5 * time.Second)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*5))
	defer cancel()

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("new mongodb client error, %s", err)
	}

	// Ping下MongoDB
	if err = client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("ping mongodb server(%s) error, %s", m.Endpoints, err)
	}

	return client, nil
}
