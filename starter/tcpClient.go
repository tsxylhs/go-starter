package starter

import (
	"errors"
	"net"
	"sync"
	"time"

	starter "github.com/tsxylhs/go-starter"
	"github.com/tsxylhs/go-starter/log"
)

type TcpClient struct {
	BaseApp
	SharedBroker *TcpBroker
}

//var clientConn sync.Map

func NewTcpClient(name string) *TcpClient {
	tcpClient := &TcpClient{
		BaseApp: BaseApp{
			name: name,
		},
	}
	tcpClient.SetPriority(PriorityHigh)
	return tcpClient

}
func (app *TcpClient) Start(ctx *starter.Context) error {
	app.Subscribe(app.name, app)
	err := (&app.BaseApp).Start(ctx)
	if err != nil {
		return err
	}
	app.BuildClient(ctx)
	return nil
}

type TcpBroker struct {
	Dailer       net.Dialer
	IsConnection bool
	Conn         net.Conn
	ConnTcp      *net.TCPConn
	mutex        sync.Mutex
}

func (tcpClient *TcpClient) BuildClient(ctx *starter.Context) error {
	var tcpConfig TcpConfig
	if tcpClient.RawConfig == nil {
		return errors.New("build tcp clients failure, no config found in core service")
	}
	err := tcpClient.RawConfig.UnmarshalKey("tcp", &tcpConfig)
	if err != nil {
		return err
	}
	if tcpConfig.Host == "" {
		return errors.New("host is null")
	}
	if tcpConfig.IsStart {
		tcpAddr, _ := net.ResolveTCPAddr("tcp4", tcpConfig.Host+":"+tcpConfig.Port)
		conn, err := net.DialTCP(tcpConfig.Protocol, nil, tcpAddr)
		if err != nil {
			log.Logger.Logger.Error("tcp 链接失败")
			tcpClient.SharedBroker.IsConnection = false
			return err
		} else {
			tcpClient.SharedBroker = &TcpBroker{
				IsConnection: true,
				ConnTcp:      conn,
			}
		}
		var heatbeatCount = 0
		for {
			if !tcpClient.SharedBroker.IsConnection {
				conn, err := net.DialTCP("tcp", nil, tcpAddr)
				tcpClient.SharedBroker.mutex.Lock()
				if err != nil {
					log.Logger.Logger.Error("tcp 链接失败")
					tcpClient.SharedBroker.IsConnection = false
				} else {

					conn.SetKeepAlive(true)
					tcpClient.SharedBroker = &TcpBroker{
						IsConnection: true,
						ConnTcp:      conn,
					}
				}

			} else {
				//心跳检验
				heatbeatCount++
				if heatbeatCount == tcpConfig.Heartbeat {
					heatbeatCount = 0
					tcpClient.SharedBroker.Conn.Write([]byte("is connected"))
					time.Sleep(100 * time.Millisecond)

				}

			}
			time.Sleep(1 * time.Second)
		}
	}
	return nil

}

type TcpConfig struct {
	Protocol  string
	Host      string
	Port      string
	Heartbeat int
	IsStart   bool
}
