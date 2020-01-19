package main

import (
	"net"
	"math/rand"
	"fmt"
	"github.com/golang/protobuf/proto"
	"encoding/binary"
	"bytes"
	"time"
	"io"
	"xingo_demo/pb"
	"os"
	"os/signal"
	"github.com/viphxin/xingo/fnet"
)

type PkgData struct {
	Len   uint32
	MsgId uint32
	Data  []byte
}

/*
简单AI规则
 */
func GenActionData() int32{
	var action int32 = 0
	//gen op count
	opCount := rand.Intn(6) + 1

	for i := 0; i < opCount; i++{
		v := rand.Intn(2)
		if v == 1 {
			action += 1 << uint(i)
		}
	}
	return action
}

type MyPtotoc struct{
	exitChan chan bool
	X float32
	Y float32
	Z float32
	V float32
	Pid int32
}

func (this *MyPtotoc)OnConnectionMade(fconn *net.TCPConn){
	fmt.Println("链接建立")
	this.exitChan = make(chan bool, 1)
}

func (this *MyPtotoc)OnConnectionLost(fconn *net.TCPConn){
	fmt.Println("链接丢失")
	this.exitChan <- true
}

func (this *MyPtotoc) Unpack(headdata []byte) (head *PkgData, err error) {
	headbuf := bytes.NewReader(headdata)

	head = &PkgData{}

	// 读取Len
	if err = binary.Read(headbuf, binary.LittleEndian, &head.Len); err != nil {
		return nil, err
	}

	// 读取MsgId
	if err = binary.Read(headbuf, binary.LittleEndian, &head.MsgId); err != nil {
		return nil, err
	}

	// 封包太大
	//if head.Len > MaxPacketSize {
	//	return nil, packageTooBig
	//}

	return head, nil
}

func (this *MyPtotoc) Pack(msgId uint32, data proto.Message) (out []byte, err error) {
	outbuff := bytes.NewBuffer([]byte{})
	// 进行编码
	dataBytes := []byte{}
	if data != nil {
		dataBytes, err = proto.Marshal(data)
	}

	if err != nil {
		fmt.Println(fmt.Sprintf("marshaling error:  %s", err))
	}
	// 写Len
	if err = binary.Write(outbuff, binary.LittleEndian, uint32(len(dataBytes))); err != nil {
		return
	}
	// 写MsgId
	if err = binary.Write(outbuff, binary.LittleEndian, msgId); err != nil {
		return
	}

	//all pkg data
	if err = binary.Write(outbuff, binary.LittleEndian, dataBytes); err != nil {
		return
	}

	out = outbuff.Bytes()
	return

}

func (this *MyPtotoc)DoMsg(fconn *net.TCPConn, pdata *PkgData){
	//处理消息
	//fmt.Println(fmt.Sprintf("msg id :%d, data len: %d", pdata.MsgId, pdata.Len))
	if pdata.MsgId == 1{
		syncpid := &pb.SyncPid{}
		proto.Unmarshal(pdata.Data, syncpid)
		this.Pid = syncpid.Pid
	}else if pdata.MsgId == 200{
		bdata := &pb.BroadCast{}
		proto.Unmarshal(pdata.Data, bdata)
		if bdata.Tp == 2 && bdata.Pid == this.Pid{
			//本人
			this.X = bdata.GetP().X
			this.Y = bdata.GetP().Y
			this.Z = bdata.GetP().Z
			this.V = bdata.GetP().V
			fmt.Println(fmt.Sprintf("player ID: %d" , bdata.Pid))
			go func() {
				for{
					if len(this.exitChan) >= 1{
						close(this.exitChan)
						return
					}
					this.WalkOrTalk(fconn)
				}
			}()
		}else if bdata.Tp == 1{
			//fmt.Println(fmt.Sprintf("世界聊天,玩家%d: %s", bdata.Pid, bdata.GetContent()))
		}
	}
}

func (this *MyPtotoc)WalkOrTalk(fconn *net.TCPConn){
	//聊天或者移动
	time.Sleep(3*time.Second)
	tp := rand.Intn(2)
	if tp == 0{
		//聊天
		//msg := &pb.Talk{
		//	Content: "你猜猜我是谁？",
		//}
		//this.Send(2, msg)
	}else{
		//移动
		x := this.X
		z := this.Z

		rrr := rand.Intn(2)
		if rrr == 0{
			x -= float32(rand.Intn(10))
			z -= float32(rand.Intn(10))
		}else{
			x += float32(rand.Intn(10))
			z += float32(rand.Intn(10))
		}

		//纠正坐标位置
		if x > 410{
			x = 410
		}else if x < 85{
			x = 85
		}

		if z > 400{
			z = 400
		}else if z < 75{
			z = 75
		}

		rv := rand.Intn(2)
		v := this.V
		if rv == 0{
			v = 25
		}else{
			v = 335
		}
		msg := &pb.Position{
			X: x,
			Y: this.Y,
			Z: z,
			V: v,
			}

		fmt.Println(fmt.Sprintf("player ID: %d. Walking..." , this.Pid))
		this.Send(fconn, 3, msg)
	}
}

func (this *MyPtotoc)Send(fconn *net.TCPConn, msgID uint32, data proto.Message){
	dd, err := this.Pack(msgID, data)
	if err == nil{
		fconn.Write(dd)
	}else{
		fmt.Println(err)
	}

}

func (this *MyPtotoc)StartReadThread(fconn *net.TCPConn){
	go func() {
		for {
		//read per head data
		headdata := make([]byte, 8)

		if _, err := io.ReadFull(fconn, headdata); err != nil {
			fmt.Println(err)
			this.OnConnectionLost(fconn)
			return
		}
		pkgHead, err := this.Unpack(headdata)
		if err != nil {
			this.OnConnectionLost(fconn)
			return
		}
		//data
		if pkgHead.Len > 0 {
			pkgHead.Data = make([]byte, pkgHead.Len)
			if _, err := io.ReadFull(fconn, pkgHead.Data); err != nil {
				this.OnConnectionLost(fconn)
				return
			}
		}
		this.DoMsg(fconn, pkgHead)
	}
	}()
}

func main() {
	for i := 0; i< 10; i ++{
		client := fnet.NewTcpClient("0.0.0.0", 8999, &MyPtotoc{})
		client.Start()
		time.Sleep(1*time.Second)
	}

	// close
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	sig := <-c
	fmt.Println("=======", sig)
}