package aliyunsls

import (
	"fmt"
	sls "github.com/aliyun/aliyun-log-go-sdk"
	"github.com/sirupsen/logrus"
	"github.com/thalesfu/golangutils/logging"
	"google.golang.org/protobuf/proto"
)

// SlsHook 实现 logrus.Hook 接口
type SlsHook struct {
	Client   sls.ClientInterface
	Project  string
	Logstore string
	Topic    string
	Source   string
}

func (h *SlsHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (h *SlsHook) Fire(entry *logrus.Entry) error {
	timestamp := uint32(entry.Time.Unix())
	contents := []*sls.LogContent{
		{Key: proto.String("level"), Value: proto.String(entry.Level.String())},
		{Key: proto.String("severityText"), Value: proto.String(entry.Level.String())},
		{Key: proto.String("content"), Value: proto.String(entry.Message)},
	}

	logDataGroup := logging.GetLogContext(entry.Context)

	for _, logData := range logDataGroup {
		for k, v := range logData {
			contents = append(contents, &sls.LogContent{
				Key:   proto.String(k),
				Value: proto.String(fmt.Sprintf("%v", v)),
			})
		}
	}

	for k, v := range entry.Data {
		contents = append(contents, &sls.LogContent{
			Key:   proto.String(k),
			Value: proto.String(fmt.Sprintf("%v", v)),
		})
	}

	log := &sls.Log{
		Time:     &timestamp,
		Contents: contents,
	}
	group := &sls.LogGroup{
		Topic:  proto.String(h.Topic),
		Source: proto.String(h.Source),
		Logs:   []*sls.Log{log},
	}
	return h.Client.PutLogs(h.Project, h.Logstore, group)
}

func GetLogrusHood(region string, project string, logstore string, topic string, accessKeyId string, accessKeySecret string, endpoint string) *SlsHook {
	provider := sls.NewStaticCredentialsProvider(accessKeyId, accessKeySecret, "")
	client := sls.CreateNormalInterfaceV2(endpoint, provider)
	// 设置使用 v4 签名
	client.SetAuthVersion(sls.AuthV4)
	// 设置地域
	client.SetRegion(region)

	return &SlsHook{
		Client:   client,
		Project:  project,
		Logstore: logstore,
		Topic:    topic,
	}
}
