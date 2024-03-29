package starter

import (
	"errors"
	"html/template"

	code "github.com/tsxylhs/go-starter/domain"

	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"xorm.io/xorm"
)

type DbHolder interface {
	SetDb(*xorm.Engine)
}

type DbStarter struct {
	BaseStarter
	Namespace string
	DbHolder  DbHolder
	listeners map[string][]code.DBListener
}

var dbListeners map[string][]code.DBListener

func ListenDB(listeners ...code.DBListener) {
	if dbListeners == nil {
		dbListeners = map[string][]code.DBListener{}
	}

	for _, listener := range listeners {
		if listener.DbEnabled() {
			dbListeners[listener.GetDbName()] = append(dbListeners[listener.GetDbName()], listener)
		}
	}
}

func (code *DbStarter) Start(ctx *code.Context) error {
	cfg := ctx.MustGet(code.Namespace + ".config").(*viper.Viper)

	dbns := cfg.GetStringMap("db")
	if len(dbns) == 0 {
		//log.Slog.Warn("no db config found for db code ", code.name)
		return nil
	}

	for dbn, _ := range dbns {
		var conn *xorm.Engine
		var err error

		if ctx.Get("db."+dbn) == nil {
			conn, err = BuildDBConnection(cfg.Sub("db." + dbn))
			if err != nil {
				return err
			}
			ctx.Set("db."+dbn, conn)

			if len(dbListeners) == 0 || len(dbListeners[dbn]) == 0 {
				continue
			}
			for _, listener := range dbListeners[dbn] {
				listener.SetDB(conn)
			}
		}
	}

	return nil
}

type RedisHolder interface {
	SetRedisConnection(*redis.Client)
}

type RedisStarter struct {
	BaseStarter
	Namespace   string
	RedisHolder RedisHolder
}

func (code *RedisStarter) Start(ctx *code.Context) error {
	cfg := ctx.MustGet(code.Namespace + ".config").(*viper.Viper)

	dbn := cfg.GetString(code.Namespace + ".redis")
	if dbn == "" {
		return nil
	}

	var conn *redis.Client
	var err error
	if ctx.Get("redis."+dbn) == nil {
		conn, err = BuildRedisConnection(cfg.Sub("redis." + dbn))
		if err != nil {
			return err
		}
		ctx.Set("redis."+dbn, conn)
	} else {
		conn = ctx.Get("redis." + dbn).(*redis.Client)
	}

	code.RedisHolder.SetRedisConnection(conn)

	return nil
}

type dbConfig struct {
	Clustered bool
	Name      string
	Ref       string
	Type      string
	Uri       string
	MaxIdle   int
	MaxOpen   int
	ShowSQL   bool
}

func BuildDBConnection(config *viper.Viper) (*xorm.Engine, error) {
	conf := dbConfig{}
	err := config.Unmarshal(&conf)
	if err != nil {
		return nil, err
	}
	if conf.Clustered {
		//TODO build by cluster config

		return nil, nil
	}

	engine, err := xorm.NewEngine(conf.Type, conf.Uri)
	if err != nil {
		return engine, err
	}

	if conf.MaxIdle > 0 {
		engine.SetMaxIdleConns(conf.MaxIdle)
	}

	if conf.MaxOpen > 0 {
		engine.SetMaxOpenConns(conf.MaxOpen)
	}

	engine.ShowSQL(conf.ShowSQL)

	return engine, err
}

// type redisConfig struct {
// 	MaxIdle     int
// 	IdleTimeout int
// 	Server      string
// 	Auth        bool
// 	Password    string
// }

type redisConfig struct {
	redis.Options `mapstructure:",squash"`
}

func BuildRedisConnection(config *viper.Viper) (*redis.Client, error) {
	if config == nil {
		return nil, errors.New("nil config when build redis connection")
	}
	conf := redisConfig{}
	err := config.Unmarshal(&conf)
	if err != nil {
		return nil, err
	}
	return redis.NewClient(&conf.Options), nil
}

type HtmlTemplateStarter struct {
	*BaseStarter
	RootDir             string
	HtmlTemplateHolder  **template.Template
	HtmlTemplateFuncMap template.FuncMap
}

func (code *HtmlTemplateStarter) Start() (err error) {
	if code.RootDir == "" {
		return errors.New("no template root")
	}

	*code.HtmlTemplateHolder = template.Must(template.New("").Funcs(code.HtmlTemplateFuncMap).ParseGlob(code.RootDir + "/*.html"))
	return nil
}
