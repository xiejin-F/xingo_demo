package main

import (
	"github.com/viphxin/xingo/iface"
	"github.com/viphxin/xingo/logger"
	"github.com/viphxin/xingo/utils"
	"xingo_demo/api"
	"xingo_demo/core"

	_ "net/http"
	_ "net/http/pprof"
	_ "runtime/pprof"
	_ "time"
	"github.com/viphxin/xingo"
	"xingo_demo/cmd"
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

func main() {
	s := xingo.NewXingoTcpServer()

	//add gm command
	if utils.GlobalObject.CmdInterpreter != nil{
		utils.GlobalObject.CmdInterpreter.AddCommand(cmd.NewOnlineCommand())
	}

	//add api ---------------start
	s.AddRouter("0", &api.Api0Router{})
	s.AddRouter("2", &api.Api2Router{})
	s.AddRouter("3", &api.Api3Router{})
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

	//s.Start()
	//// close
	//c := make(chan os.Signal, 1)
	//signal.Notify(c, os.Interrupt, os.Kill)
	//sig := <-c
	//fmt.Println("=======", sig)
	//s.Stop()
	s.Serve()
}
