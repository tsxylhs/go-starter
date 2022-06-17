package app

import "code"

type Service struct {
	BaseApp
}

func NewService(name string) *Service {
	service := &Service{
		BaseApp: BaseApp{
			name: name,
		},
	}
	service.SetPriority(PriorityHigh)
	return service
}

func (app *Service) Start(ctx *code.Context) error {
	//配置文件
	app.Subscribe(app.name, app)

	err := (&app.BaseApp).Start(ctx)
	if err != nil {
		return err
	}

	return nil
}
