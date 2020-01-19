package cmd

import (
	"fmt"
	"xingo_demo/core"
)

type OnlineCommand struct {
}

func NewOnlineCommand() *OnlineCommand{
	return &OnlineCommand{}
}
func (this *OnlineCommand)Name()string{
	return "online"
}

func (this *OnlineCommand)Help()string{
	return fmt.Sprintf("online:\r\n" +
		"----------- login: 登陆玩家数\r\n" +
		"----------- nologin:  未登陆玩家数\r\n" +
		"----------- kick [userId]:  踢出玩家")
}

func (this *OnlineCommand)Run(args []string) string{
	if len(args) == 0{
		return this.Help()
	}else{
		switch args[0] {
		case "login":
			core.WorldMgrObj.RLock()
			ss := len(core.WorldMgrObj.Players)
			core.WorldMgrObj.RUnlock()
			return fmt.Sprintf("login player: %d", ss)
		default:
			return "未实现"
		}
	}
	return "OK"
}
