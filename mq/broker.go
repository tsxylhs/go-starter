package mq

import (
	"crypto/tls"
	"sync"

	mqtt "github.com/eclipse/paho.mqtt.golang"

	common "github.com/tsxylhs/go-starter"
	"github.com/tsxylhs/go-starter/log"
	"go.uber.org/zap"
)

const clientNameDefault = "default"

type Config struct {
	mqtt.ClientOptions `mapstructure:",squash"`
	Host               string
	Default            bool
}
type Broker struct {
	mqtt.Client
}

var SharedBroker *Broker

type IMessage interface {
	GetType() string
	GetPayload() []byte
}

type Message struct {
	Type    string
	Payload []byte
}

func (msg *Message) GetType() string {
	return msg.Type
}

func (msg *Message) GetPayload() []byte {
	return msg.Payload
}

type Handle func(msg IMessage)

type Options map[string]interface{}
type Dispatcher interface {
	//订阅
	Sub(topic string, handle Handle, options ...Options) error
	//发布
	Pub(topic string, payload []byte, options ...Options) error
	UnSub(topic string) error
	//bridging 桥接到消息总线  将实现和分离分开
	Bridging() error
}

// UnSub 取消订阅单个的topic
func (dispatcher *Broker) UnSub(topic string) error {
	return dispatcher.Unsubscribe(topic).Error()
}

//Pub 消息发送
func (dispatcher *Broker) Pub(topic string, payload []byte, options ...Options) error {

	if token := dispatcher.Publish(topic, 2, false, payload); token.Error() != nil {
		log.Logger.Debug("fail to dispatch message", zap.Any("pub", token.Error()))
		return token.Error()
	}
	return nil
}

// Sub 消息订阅
func (dispatcher *Broker) Sub(topic string, handle func(client mqtt.Client, message mqtt.Message), options ...Options) error {
	log.Logger.Debug("开始订阅")
	if token := dispatcher.Subscribe(topic, 2, handle); token.Error() != nil {
		log.Logger.Debug("订阅失败", zap.Any("sub", token.Error()))
		return token.Error()
	}
	return nil
}

func NewBroker(cfg Config) (*Broker, error) {
	client, err := buildClient(cfg)
	return &Broker{
		Client: client,
	}, err
}

func buildClient(cfg Config) (client mqtt.Client, err error) {
	log.Logger.Debug("create mqtt broker", zap.String("host", cfg.Host))
	opts := cfg.AddBroker(cfg.Host)
	opts.SetClientID(cfg.ClientID)
	opts.OnConnect = func(c mqtt.Client) {
		c.Publish("/as/connect", 0, true, []byte("Connected"))
	}
	tlscfg := &tls.Config{
		ClientAuth:         tls.NoClientCert,
		ClientCAs:          nil,
		InsecureSkipVerify: true,
	}
	opts.SetTLSConfig(tlscfg)
	client = mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return client, token.Error()
	} else {
		log.Logger.Debug("Connected to mqtt server", zap.String("host", cfg.Host))
	}
	return client, nil

}

var Clients map[string]mqtt.Client
var m sync.Mutex

func BuildClients(ctx *common.Context) error {
	m.Lock()
	defer m.Unlock()
	if Clients == nil {
		Clients = map[string]mqtt.Client{}
	}

	var confs map[string]Config
	// if app.Service.RawConfig == nil {
	// 	return errors.New("build mqtt clients failure, no config found in core service")
	// }
	// err := app.Service.RawConfig.UnmarshalKey("mqtt", &confs)
	// if err != nil {
	// 	return err
	// }

	if confs == nil {
		return nil
	}
	for key, value := range confs {
		client, err := buildClient(value)
		if err != nil {
			return err
		}
		Clients[key] = client
		if value.Default {
			Clients[clientNameDefault] = client
			SharedBroker = &Broker{
				Client: client,
			}
		}
	}
	if SharedBroker == nil && len(confs) > 0 {
		for _, value := range Clients {
			SharedBroker = &Broker{
				Client: value,
			}
			break
		}
	}
	return nil

}
