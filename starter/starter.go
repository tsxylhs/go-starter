package starter

//启动类
import (
	"fmt"
	"reflect"
	"strconv"
	"sync"

	code "github.com/tsxylhs/go-starter/domain"
)

const (
	PriorityHighest = 1000
	PriorityHigh    = 900
	PriorityMiddle  = 600
	PriorityLow     = 300
	PriorityLowest  = 0
)

type Starter interface {
	Name() string
	Priority() int
	SetPriority(int) Starter
	SetAppName(appName string) Starter
	SetApp(app App) Starter
	AppName() string
	Start(ctx *code.Context) error
	Started() bool
	SetStarted(bool) Starter
}

type BaseStarter struct {
	name     string
	priority int
	started  bool
	appName  string
	app      App
	action   func(ctx *code.Context) error
}

func NewBaseStarter(name string, priority int) *BaseStarter {
	return &BaseStarter{
		name:     name,
		priority: priority,
	}
}

// basseStarter set封装
func (base *BaseStarter) Name() string {
	return base.name
}

func (base *BaseStarter) Priority() int {
	return base.priority
}

func (base *BaseStarter) SetPriority(priority int) Starter {
	base.priority = priority
	return base
}

func (base *BaseStarter) SetAppName(appName string) Starter {
	base.appName = appName
	return base
}

func (base *BaseStarter) AppName() string {
	return base.appName
}

func (base *BaseStarter) Started() bool {
	return base.started
}

func (base *BaseStarter) SetStarted(started bool) Starter {
	base.started = started
	return base
}

func (base *BaseStarter) SetApp(app App) Starter {
	base.app = app
	return base
}

func (base *BaseStarter) Action(action func(ctx *code.Context) error) Starter {
	base.action = action
	return base
}

func (base *BaseStarter) Start(ctx *code.Context) error {
	if base.action != nil {
		return base.action(ctx)
	}

	return nil
}

// 启动器监听类
type StartListener func(ctx code.Context) error

func OnStarted(starterName string, listener StartListener) {
	if controller.listenersMap == nil {
		controller.listenersMap = map[string][]StartListener{}
	}
	controller.listenersMap[starterName] = append(controller.listenersMap[starterName], listener)
}

// 启动控制类
type StartController struct {
	ctx           code.Context
	mu            sync.RWMutex
	startersMap   map[string]Starter
	startersArray []Starter
	listenersMap  map[string][]StartListener
}

var (
	controller = &StartController{}
)

func RegisterStarter(code Starter) {
	fmt.Println("Register code >> " + code.Name() + " [" + reflect.TypeOf(code).String() + "]")
	controller.register(code)
}
func Start() error {
	controller.ctx = code.Context{}
	err := controller.startNext()
	controller = nil
	return err
}

// 注册启动器
func (controller *StartController) register(code Starter) {
	controller.mu.Lock()
	if controller.startersMap == nil {
		controller.startersMap = make(map[string]Starter)
	}
	if controller.startersMap[code.Name()] == nil {
		controller.startersMap[code.Name()] = code
		var arr []Starter
		var added = false
		for _, str := range controller.startersArray {
			if code.Priority() > str.Priority() && !added {
				arr = append(arr, code)
				added = true
			}
			arr = append(arr, str)
		}
		if !added {
			arr = append(arr, code)
		}
		controller.startersArray = arr
	} else {
		panic("重复添加启动器：" + code.Name())
	}
	controller.mu.Unlock()

}

// 打印装装载好的启动器
func printStarters(prefix string, starts []Starter) {
	str := ""
	for index, st := range controller.startersArray {
		if index > 0 {
			str += ","
		}
		str += st.Name() + ":" + strconv.Itoa(st.Priority())
	}
	fmt.Println(prefix + "-----" + str)
}

// 启动启动项
func (controller *StartController) startStarter(code Starter) error {
	if code.Started() {
		panic("该启动器已启动" + code.Name())
	}

	// 传入全局ctx
	err := code.Start(&controller.ctx)
	if err != nil {
		return err
	} else {
		code.SetStarted(true)
		listeners := controller.listenersMap[code.Name()]
		if listeners != nil {
			for _, listener := range listeners {
				if err = listener(controller.ctx); err != nil {
					return err
				}
			}
		}
		fmt.Println("该启动项已启动" + code.Name())
		return nil
	}

}

// 递归启动
func (controller *StartController) startNext() error {
	var err error
	var starter Starter
	controller.mu.Lock()

	if len(controller.startersArray) == 0 {
		return nil
	}
	starter = controller.startersArray[0]

	if len(controller.startersArray) > 0 {
		controller.startersArray = controller.startersArray[1:]
		printStarters(starter.Name(), controller.startersArray)
	} else {
		controller.startersArray = []Starter{}

	}
	controller.mu.Unlock()
	err = controller.startStarter(starter)
	if err != nil {
		fmt.Println("启动器启动失败", starter.Name())
		fmt.Println(err.Error())
		return err
	}
	return controller.startNext()

}
