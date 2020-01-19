package main

import (
	"fmt"
	"github.com/viphxin/xingo/fserver"
	"github.com/viphxin/xingo/iface"
	"github.com/viphxin/xingo/logger"
	"github.com/viphxin/xingo/utils"
	"xingo_demo/api"
	"xingo_demo/core"

	_ "net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	_ "runtime/pprof"
	_ "time"
	"time"
)

func DoConnectionMade(fconn iface.Iconnection) {
	logger.Debug("111111111111111111111111")
	p, _ := core.WorldMgrObj.AddPlayer(fconn)
	fconn.SetProperty("pid", p.Pid)
}

func DoConnectionLost(fconn iface.Iconnection) {
	logger.Debug("222222222222222222222222")
	pid, _ := fconn.GetProperty("pid")
	p, _ := core.WorldMgrObj.GetPlayer(pid.(int32))
	//移除玩家
	core.WorldMgrObj.RemovePlayer(pid.(int32))
	//消失在地图
	p.LostConnection()
}

func testTimer(args ...interface {}){
	logger.Info(fmt.Sprintf("%s-%d-%f", args[0], args[1], args[2]))
}

func main() {
	s := fserver.NewServer()

	//add api ---------------start
	TestRouterObj := &api.TestRouter{}
	s.AddRouter(TestRouterObj)
	//add api ---------------end
	//regest callback
	utils.GlobalObject.OnConnectioned = DoConnectionMade
	utils.GlobalObject.OnClosed = DoConnectionLost

	// go func() {
	// 	fmt.Println(http.ListenAndServe("localhost:6061", nil))
	// 	// for {
	// 	// 	time.Sleep(time.Second * 10)
	// 	// 	fm, err := os.OpenFile("./memory.log", os.O_RDWR|os.O_CREATE, 0644)
	// 	// 	if err != nil {
	// 	// 		fmt.Println(err)
	// 	// 	}
	// 	// 	pprof.WriteHeapProfile(fm)
	// 	// 	fm.Close()
	// 	// }
	// }()
	s.CallLater(5*time.Second, testTimer, "viphxin", 10009, 10.999)
	s.CallWhen("2016-12-15 18:35:10", testTimer, "viphxin", 10009, 10.999)
	s.CallLoop(5*time.Second, testTimer, "loop--viphxin", 10009, 10.999)
	s.Start()
	// close
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	sig := <-c
	fmt.Println("=======", sig)
	s.Stop()
}
