package api

import (
	"xingo_demo/pb"
	"xingo_demo/core"
	"github.com/golang/protobuf/proto"
	"github.com/viphxin/xingo/fnet"
	"github.com/viphxin/xingo/logger"
	_ "time"
	"fmt"
	"github.com/viphxin/xingo/utils"
	"github.com/viphxin/xingo/iface"
)

type Api0Router struct {
	fnet.BaseRouter
}

/*
ping test
*/
func (this *Api0Router) Handle(request iface.IRequest) {
	logger.Debug("call Api_0")
	// request.Fconn.SendBuff(0, nil)
	packdata, err := utils.GlobalObject.Protoc.GetDataPack().Pack(0, nil)
	if err == nil{
		request.GetConnection().Send(packdata)
	}else{
		logger.Error("pack data error")
	}
}


type Api2Router struct {
	fnet.BaseRouter
}
/*
世界聊天
 */
func (this *Api2Router) Handle(request iface.IRequest) {
	msg := &pb.Talk{}
	err := proto.Unmarshal(request.GetData(), msg)
	if err == nil {
		logger.Debug(fmt.Sprintf("user talk: content: %s.", msg.Content))
		pid, err1 := request.GetConnection().GetProperty("pid")
		if err1 == nil{
			p, _ := core.WorldMgrObj.GetPlayer(pid.(int32))
			p.Talk(msg.Content)
		}else{
			logger.Error(err1)
			request.GetConnection().LostConnection()
		}

	} else {
		logger.Error(err)
		request.GetConnection().LostConnection()
	}
}


type Api3Router struct {
	fnet.BaseRouter
}
/*
移动
 */
func (this *Api3Router) Handle(request iface.IRequest) {
	msg := &pb.Position{}
	err := proto.Unmarshal(request.GetData(), msg)
	if err == nil {
		logger.Debug(fmt.Sprintf("user move: (%f, %f, %f, %f)", msg.X, msg.Y, msg.Z, msg.V))
		pid, err1 := request.GetConnection().GetProperty("pid")
		if err1 == nil{
			p, _ := core.WorldMgrObj.GetPlayer(pid.(int32))
			p.UpdatePos(msg.X, msg.Y, msg.Z, msg.V)
		}else{
			logger.Error(err1)
			request.GetConnection().LostConnection()
		}

	} else {
		logger.Error(err)
		request.GetConnection().LostConnection()
	}
}