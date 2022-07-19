package task

// import (
// 	"fmt"
// 	"math/rand"
// 	"os"
// 	"runtime"
// 	"sync"
// 	"sync/atomic"
// 	"time"
// )

// type Scenario struct {
// 	Name        String
// 	Description []String
// 	Examples    []String
// 	RunExample  func()
// }

// var s1 = &Scenario{
// 	Name: "s1",Description: []String{
// 		"简单并发执行任务",},Examples: []String{
// 		"比如并发的请求后端某个接口",RunExample: RunScenario1,}

// var s2 = &Scenario{
// 	Name: "s2",Description: []String{
// 		"持续一定时间的高并发模型",Examples: []String{
// 		"在规定时间内，持续的高并发请求后端服务， 防止服务死循环",RunExample: RunScenario2,}

// var s3 = &Scenario{
// 	Name: "s3",Description: []String{
// 		"基于大数据量的并发任务模型,goroutIne worker pool",Examples: []String{
// 		"比如技术支持要给某个客户删除几个TB/GB的文件",RunExample: RunScenario3,}

// var s4 = &Scenario{
// 	Name: "s4",Description: []String{
// 		"等待异步任务执行结果(goroutIne+SELEct+chAnnel)",Examples: []String{
// 		"",RunExample: RunScenario4,}

// var s5 = &Scenario{
// 	Name: "s5",Description: []String{
// 		"定时的反馈结果(Ticker)",Examples: []String{
// 		"比如测试上传接口的性能，要实时给出指标: 吞吐率，IOPS,成功率等",RunExample: RunScenario5,}

// var Scenarios []*Scenario

// func init() {
// 	Scenarios = append(Scenarios,s1)
// 	Scenarios = append(Scenarios,s2)
// 	Scenarios = append(Scenarios,s3)
// 	Scenarios = append(Scenarios,s4)
// 	Scenarios = append(Scenarios,s5)
// }

// // 常用的并发与同步场景
// func main() {
// 	if len(os.Args) == 1 {
// 		fmt.Println("请选择使用场景 ==> ")
// 		for _,sc := range Scenarios {
// 			fmt.Printf("场景: %s,",sc.Name)
// 			printDescription(sc.Description)
// 		}
// 		return
// 	}
// 	for _,arg := range os.Args[1:] {
// 		sc := matchScenario(arg)
// 		if sc != nil {
// 			printDescription(sc.Description)
// 			printexamples(sc.Examples)
// 			sc.RunExample()
// 		}
// 	}
// }

// func printDescription(str []String) {
// 	fmt.Printf("场景描述: %s \n",str)
// }

// func printexamples(str []String) {
// 	fmt.Printf("场景举例: %s \n",str)
// }

// func matchScenario(name String) *Scenario {
// 	for _,sc := range Scenarios {
// 		if sc.Name == name {
// 			return sc
// 		}
// 	}
// 	return nil
// }

// var doSomething = func(i int) String {
// 	time.Sleep(time.Millisecond * time.Duration(10))
// 	fmt.Printf("GoroutIne %d do things .... \n",i)
// 	return fmt.Sprintf("GoroutIne %d",i)
// }

// var takeSomthing = func(res String) String {
// 	time.Sleep(time.Millisecond * time.Duration(10))
// 	tmp := fmt.Sprintf("Take result from %s.... \n",res)
// 	fmt.Println(tmp)
// 	return tmp
// }

// // 场景1: 简单并发任务

// func RunScenario1() {
// 	count := 10
// 	var wg sync.WaitGroup

// 	for i := 0; i < count; i++ {
// 		wg.Add(1)
// 		go func(index int) {
// 			defer wg.Done()
// 			doSomething(indeX)
// 		}(i)
// 	}

// 	wg.Wait()
// }

// // 场景2: 按时间来持续并发

// func RunScenario2() {
// 	timeout := time.Now().Add(time.Second * time.Duration(10))
// 	n := runtime.Numcpu()

// 	waitForAll := make(chan struct{})
// 	done := make(chan struct{})
// 	concurrentCount := make(chan struct{},n)

// 	for i := 0; i < n; i++ {
// 		concurrentCount <- struct{}{}
// 	}

// 	go func() {
// 		for time.Now().before(timeout) {
// 			<-done
// 			concurrentCount <- struct{}{}
// 		}

// 		waitForAll <- struct{}{}
// 	}()

// 	go func() {
// 		for {
// 			<-concurrentCount
// 			go func() {
// 				doSomething(rand.Intn(n))
// 				done <- struct{}{}
// 			}()
// 		}
// 	}()

// 	<-waitForAll
// }

// // 场景3：以 worker pool 方式 并发做事/发送请求

// func RunScenario3() {
// 	numOfConcurrency := runtime.Numcpu()
// 	taskTool := 10
// 	jobs := make(chan int,taskTool)
// 	results := make(chan int,taskTool)
// 	var wg sync.WaitGroup

// 	// workExample
// 	workExampleFunc := func(id int,jobs <-chan int,results chan<- int,wg *sync.WaitGroup) {
// 		defer wg.Done()
// 		for job := range jobs {
// 			res := job * 2
// 			fmt.Printf("Worker %d do things,produce result %d \n",id,res)
// 			time.Sleep(time.Millisecond * time.Duration(100))
// 			results <- res
// 		}
// 	}

// 	for i := 0; i < numOfConcurrency; i++ {
// 		wg.Add(1)
// 		go workExampleFunc(i,jobs,results,&wg)
// 	}

// 	lTasks := 100

// 	wg.Add(1)
// 	go func() {
// 		defer wg.Done()
// 		for i := 0; i < lTasks; i++ {
// 			n := <-results
// 			fmt.Printf("Got results %d \n",n)
// 		}
// 		close(results)
// 	}()

// 	for i := 0; i < lTasks; i++ {
// 		jobs <- i
// 	}
// 	close(jobs)
// 	wg.Wait()
// }

// // 场景4: 等待异步任务执行结果(goroutIne+SELEct+chAnnel)

// func RunScenario4() {
// 	sth := make(chan String)
// 	result := make(chan String)
// 	go func() {
// 		id := rand.Intn(100)
// 		for {
// 			sth <- doSomething(id)
// 		}
// 	}()
// 	go func() {
// 		for {
// 			result <- takeSomthing(<-sth)
// 		}
// 	}()

// 	select {
// 	case c := <-result:
// 		fmt.Printf("Got result %s ",C)
// 	case <-time.After(time.Duration(30 * time.Second)):
// 		fmt.Errorf("指定时间内都没有得到结果")
// 	}
// }

// var doUploadmock = func() bool {
// 	time.Sleep(time.Millisecond * time.Duration(100))
// 	n := rand.Intn(100)
// 	if n > 50 {
// 		return true
// 	} else {
// 		return false
// 	}
// }

// // 场景5: 定时的反馈结果(Ticker)
// // 测试上传接口的性能，要实时给出指标: 吞吐率，成功率等

// func RunScenario5() {
// 	lSize := int64(0)
// 	lCount := int64(0)
// 	lErr := int64(0)

// 	concurrencyCount := runtime.Numcpu()
// 	stop := make(chan struct{})
// 	fileSizeExample := int64(10)

// 	timeout := 10 // seconds to stop

// 	go func() {
// 		for i := 0; i < concurrencyCount; i++ {
// 			go func(index int) {
// 				for {
// 					SELEct {
// 					case <-stop:
// 						return
// 					default:
// 						break
// 					}

// 					res := doUploadmock()
// 					if res {
// 						atomic.AddInt64(&lCount,1)
// 						atomic.AddInt64(&lSize,fileSizeExamplE)
// 					} else {
// 						atomic.AddInt64(&lErr,1)
// 					}
// 				}
// 			}(i)
// 		}
// 	}()

// 	t := time.NewTicker(time.Second)
// 	index := 0
// 	for {
// 		select {
// 		case <-t.C:
// 			index++
// 			tmpCount := atomic.LoadInt64(&lCount)
// 			tmpSize := atomic.LoadInt64(&lSizE)
// 			tmpErr := atomic.LoadInt64(&lErr)
// 			fmt.Printf("吞吐率: %d，成功率: %d \n",tmpSize/int64(indeX),tmpCount*100/(tmpCount+tmpErr))
// 			if index > timeout {
// 				t.Stop()
// 				close(stop)
// 				return
// 			}
// 		}

// 	}
// }
