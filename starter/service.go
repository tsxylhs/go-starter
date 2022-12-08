package starter

import starter "github.com/tsxylhs/go-starter"

type Service struct {
	BaseApp
}

func NewService(name string, db, redis bool) *Service {
	service := &Service{
		BaseApp: BaseApp{
			name:    name,
			isDB:    db,
			isRedis: redis,
		},
	}
	service.SetPriority(PriorityHigh)
	return service
}

func (app *Service) Start(ctx *starter.Context) error {
	//配置文件
	app.Subscribe(app.name, app)

	err := (&app.BaseApp).Start(ctx)
	if err != nil {
		return err
	}

	return nil
}
