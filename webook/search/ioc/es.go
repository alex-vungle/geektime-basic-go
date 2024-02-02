package ioc

import (
	"fmt"
	"gitee.com/geekbang/basic-go/webook/search/repository/dao"
	"github.com/olivere/elastic/v7"
	"github.com/spf13/viper"
	"time"
)

// InitESClient 读取配置文件，进行初始化ES客户端
func InitESClient() *elastic.Client {
	type Config struct {
		Url   string `yaml:"url"`
		Sniff bool   `yaml:"sniff"`
	}
	var cfg Config
	err := viper.UnmarshalKey("es", &cfg)
	if err != nil {
		panic(fmt.Errorf("读取 ES 配置失败 %w", err))
	}
	const timeout = 100 * time.Second
	opts := []elastic.ClientOptionFunc{
		elastic.SetURL(cfg.Url),
		elastic.SetHealthcheckTimeoutStartup(timeout),
	}
	client, err := elastic.NewClient(opts...)
	if err != nil {
		panic(err)
	}
	err = dao.InitES(client)
	if err != nil {
		panic(err)
	}
	return client
}
