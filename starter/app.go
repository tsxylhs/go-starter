package starter

import (
	"errors"
	"fmt"

	code "github.com/tsxylhs/go-starter/domain"

	"github.com/tsxylhs/go-starter/config"

	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"xorm.io/xorm"
)

// app interface 实现
type App interface {
	Starter
	Mount(app ...App) App
	SetMaster(master App)
	Master() App
	IsMaster() bool
}

//基类
type BaseApp struct {
	Configurator
	name     string
	master   App
	Mounts   *[]App
	isMaster bool
	Rpc      bool
	modules  []code.IModule
	isDB     bool
	isRedis  bool
	DB       *xorm.Engine
	Redis    *redis.Client
}

// 封装
func (app *BaseApp) SetMaster(master App) {
	app.master = master
}
func (app *BaseApp) Name() string {
	return app.name
}
func (app *BaseApp) Master() App {
	return app.master
}
func (app *BaseApp) IsMaster() bool {
	return app.isMaster
}
func (app *BaseApp) SetRedisConnection(client *redis.Client) {
	app.Redis = client
}
func (app *BaseApp) SetConfigFileName(name string) App {
	app.Configurator.FileName = name
	return app
}

// 启动项挂载
func (app *BaseApp) Mount(apps ...App) App {
	app.isMaster = true
	if app.Mounts == nil {
		app.Mounts = &[]App{}
	}
	for _, ap := range apps {
		*app.Mounts = append(*app.Mounts, ap)
		if ap.Priority() > app.priority {
			ap.SetPriority(app.priority - 1)
		}
		ap.SetMaster(app)
	}
	return app
}

// 启动
func (app *BaseApp) Start(ctx *code.Context) error {
	(&app.Configurator).SetApp(app)
	err := (&app.Configurator).Start(ctx)
	if err != nil {
		return err
	}
	dbn := app.RawConfig.GetString(app.name + ".db")
	for _, module := range app.modules {
		if module.DbEnabled() {
			if module.GetDbName() == "" {
				module.SetDbName(dbn)
			}
			if dbn == "" {
				panic("db enabled for module " + module.GetName() + ", but no name specified for module or app " + app.name)
			}
			if ctx.Get("db."+dbn) == nil {
				//db启动器
				ListenDB(module)
			} else {
				module.SetDB(ctx.Get("db." + dbn).(*xorm.Engine))
			}
		}
	}
	//
	//数据库
	if app.isDB {
		RegisterStarter(&DbStarter{
			BaseStarter: BaseStarter{
				name:     app.Name() + ".DB",
				priority: PriorityMiddle,
			},
			Namespace: app.name,
			//DbHolder:  app,
		})
		fmt.Println("数据库的启动器已注册")
	}
	if app.isRedis {
		//redis
		RegisterStarter(&RedisStarter{
			BaseStarter: BaseStarter{
				name:     app.Name() + ".REDIS",
				priority: PriorityMiddle,
			},
			Namespace:   app.name,
			RedisHolder: app,
		})
		fmt.Println("redis的启动器已注册")
	}
	if app.isMaster && app.Mounts != nil {
		fmt.Println("register mounts")
		for _, mnt := range *app.Mounts {
			RegisterStarter(mnt)
		}
	}

	return nil
}

// 添加注册多个模块模型表
func (app *BaseApp) Register(modules ...code.IModule) {
	app.modules = append(app.modules, modules...)
}

type Configurator struct {
	BaseStarter
	FileName     string
	RawConfig    *viper.Viper
	Subscription []config.Pair
}

// 封装set
func (configurator *Configurator) Subscribe(key string, target interface{}) {
	configurator.Subscription = append(configurator.Subscription, config.Pair{Key: key, Target: target})
}

func (configurator *Configurator) Start(ctx *code.Context) error {
	fileName := configurator.FileName
	if fileName == "" {
		fileName = configurator.app.Name()
	}
	configurator.RawConfig = viper.New()
	err := LoadConfig(fileName, configurator.RawConfig)
	if err != nil {
		if configurator.app.IsMaster() {
			return errors.New("未找到配置文件")
		} else {
			configurator.RawConfig = ctx.Get(configurator.app.Master().Name() + ".config").(*viper.Viper)
		}
	}
	ctx.Set(configurator.app.Name()+".config", configurator.RawConfig)
	if configurator.RawConfig != nil && configurator.Subscription != nil {
		for _, pair := range configurator.Subscription {
			//if configurator.RawConfig == nil {
			//	fmt.Println("unmarshal config: ", pair.Key, nil)
			//} else {
			fmt.Println("unmarshal config: ", pair.Key, configurator.RawConfig.Get(pair.Key))
			//}

			err = configurator.RawConfig.UnmarshalKey(pair.Key, pair.Target)
			if err != nil {
				return err
			}
			//fmt.Println(pair.Target)
		}
	}

	return nil
}
