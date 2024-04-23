package startup

import (
	"gitee.com/geekbang/basic-go/webook/search/events"
	"gitee.com/geekbang/basic-go/webook/search/ioc"
	"gitee.com/geekbang/basic-go/webook/search/repository"
	"gitee.com/geekbang/basic-go/webook/search/repository/dao"
	"gitee.com/geekbang/basic-go/webook/search/service"
	"github.com/IBM/sarama"
	"github.com/olivere/elastic/v7"
	"github.com/spf13/viper"
)

func InitTestSvc() (service.SearchService, sarama.SyncProducer, dao.AnyDAO, dao.ArticleDAO, dao.LikeDAO, dao.CollectDAO, *elastic.Client) {
	// 直接指定文件路径
	viper.SetConfigFile("startup/dev.yaml")
	viper.WatchConfig()
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	client := ioc.InitESClient()
	anyDAO := dao.NewAnyESDAO(client)
	anyRepository := repository.NewAnyRepository(anyDAO)
	userDAO := dao.NewUserElasticDAO(client)
	userRepository := repository.NewUserRepository(userDAO)
	articleDAO := dao.NewArticleElasticDAO(client)
	tagDAO := dao.NewTagESDAO(client)
	likeDAO := dao.NewLikeDAO(client)
	collectDAO := dao.NewCollectDAO(client)
	articleRepository := repository.NewArticleRepository(articleDAO, tagDAO, collectDAO, likeDAO)
	syncService := service.NewSyncService(anyRepository, userRepository, articleRepository)
	searchService := service.NewSearchService(userRepository, articleRepository)
	loggerV1 := ioc.InitLogger()
	saramaClient := ioc.InitKafka()
	createTopic(saramaClient, events.InteractiveTopic)
	interactiveConsumer := events.NewInteractiveConsumer(saramaClient, loggerV1, syncService)
	err = interactiveConsumer.Start()
	if err != nil {
		panic(err)
	}
	p, err := sarama.NewSyncProducerFromClient(saramaClient)
	if err != nil {
		panic(err)
	}
	return searchService, p, anyDAO, articleDAO, likeDAO, collectDAO, client
}

func createTopic(client sarama.Client, topic string) {

	partitions := int32(1)
	replicationFactor := int16(1)
	detail := &sarama.TopicDetail{
		NumPartitions:     partitions,
		ReplicationFactor: replicationFactor,
	}

	admin, err := sarama.NewClusterAdminFromClient(client)
	if err != nil {
		panic(err)
	}

	// 检查话题是否已存在
	existingTopics, err := admin.ListTopics()
	if err != nil {
		panic(err)
	}

	if _, exists := existingTopics[topic]; !exists {
		// 话题不存在，创建它
		err = admin.CreateTopic(topic, detail, false)
		if err != nil {
			panic(err)
		}
	}
}
