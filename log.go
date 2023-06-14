package aliyunlog

import (
	"fmt"
	"time"

	sls "github.com/aliyun/aliyun-log-go-sdk"
	"github.com/aliyun/aliyun-log-go-sdk/producer"
	config "gitlab.qkids.com/group-api-common/go-conf.git"
	"google.golang.org/protobuf/proto"
)

const (
	DEBUG = "DEBUG"
	INFO  = "INFO"
	ERROR = "ERROR"
)

var producerInstance *producer.Producer

func init() {
}

func getInstance() *producer.Producer {
	if producerInstance != nil {
		return producerInstance
	}

	producerConfig := producer.GetDefaultProducerConfig()
	producerConfig.AccessKeyID = config.GetString("aliyun.log.accessKeyId")
	producerConfig.AccessKeySecret = config.GetString("aliyun.log.accessKeySecret")
	producerConfig.Endpoint = config.GetString("aliyun.log.endpoint")
	producerInstance = producer.InitProducer(producerConfig)
	producerInstance.Start()
	return producerInstance
}

func LogMessage(level, message string) {
	if config.GetBool("aliyun.log.report", false) {
		Log(level, []*sls.LogContent{
			{
				Key:   proto.String("message"),
				Value: proto.String(message),
			},
		})
	} else {
		fmt.Printf("[%s] %s", level, message)
	}
}

func Log(level string, contents []*sls.LogContent) {
	go func() {
		contents = append(contents, &sls.LogContent{
			Key:   proto.String("level"),
			Value: proto.String(level),
		})

		log := &sls.Log{
			Time:     proto.Uint32(uint32(time.Now().Unix())),
			Contents: contents,
		}
		getInstance().SendLog(config.GetString("aliyun.log.project"), config.GetString("aliyun.log.logstore"), config.GetString("aliyun.log.topic"), "127.0.0.1", log)
	}()
}

func Debug(message string) {
	LogMessage(DEBUG, message)
}

func Info(message string) {
	LogMessage(INFO, message)
}

func Error(message string) {
	LogMessage(ERROR, message)
}
