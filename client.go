package xin

import (
	"net"
	"fmt"
	"github.com/golang/protobuf/proto"
	"io"
	"encoding/binary"
	"bytes"
	"os"
	"os/signal"
	"xingo_demo/pb"
	"math/rand"
	"time"
)

type PkgData struct {
	Len   uint32
	MsgId uint32
	Data  []byte
}

type TcpClient struct{
	conn *net.TCPConn
	addr *net.TCPAddr
	X float32
	Y float32
	Z float32
	V float32
	Pid int32
}

type INPUT int32

const (
	_       = iota
	UP INPUT = 1 << iota
	DOWN
	LEFT
	RIGHT
	ROTATE_LEFT
	ROTATE_RIGHT
)

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

func NewTcpClient(ip string, port int) *TcpClient{
	addr := &net.TCPAddr{
		IP: net.ParseIP(ip),
		Port: port,
		Zone: "",
	}
	conn, err := net.DialTCP("tcp", nil, addr)
	if err == nil{
		client := &TcpClient{
		conn: conn,
		addr: addr,
		}
		client.ConnectionMade()
		return client
	}else{
		panic(err)
	}

}

func (this *TcpClient)ConnectionMade(){
	fmt.Println("链接建立")
}

func (this *TcpClient)ConnectionLost(){
	fmt.Println("链接断开")
}

func (this *TcpClient) Unpack(headdata []byte) (head *PkgData, err error) {
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

func (this *TcpClient) Pack(msgId uint32, data proto.Message) (out []byte, err error) {
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

func (this *TcpClient)DoMsg(pdata *PkgData){
	//处理消息
	fmt.Println(fmt.Sprintf("msg id :%d, data len: %d", pdata.MsgId, pdata.Len))
	if pdata.MsgId == 1{
		syncpid := &pb.SyncPid{}
		proto.Unmarshal(pdata.Data, syncpid)
		this.Pid = syncpid.Pid
	}else if pdata.MsgId == 200{
		bdata := &pb.BroadCast{}
		proto.Unmarshal(pdata.Data, bdata)
		if bdata.Tp == 2{
			this.X = bdata.GetP().X
			this.Y = bdata.GetP().Y
			this.Z = bdata.GetP().Z
			this.V = bdata.GetP().V
			fmt.Println(fmt.Sprintf("player ID: %d" , bdata.Pid))
		}else{
			fmt.Println(fmt.Sprintf("世界聊天,玩家%d: %s", bdata.Pid, bdata.GetContent()))
		}
		//聊天或者移动
		time.Sleep(3*time.Second)
		tp := rand.Intn(2)
		if tp == 0{
			//聊天
			msg := &pb.Talk{
				Content: "你猜猜我是谁？",
			}
			this.Send(2, msg)
		}else{
			//移动
			x := this.X
			y := this.Y
			if x >=410{
				x -= 1
			}else if x < 85{
				x += 1
			}

			if y >=400{
				y -= 1
			}else if y < 75{
				y += 1
			}
			msg := &pb.MovePackege{
				P : &pb.Position{
				X: this.X + 1,
				Z: this.Z + 1,
				V: this.V,
				},
				ActionData: GenActionData(),
			}
			this.Send(3, msg)
		}
	}
}

func (this *TcpClient)Send(msgID uint32, data proto.Message){
	//fmt.Println("Send")
	dd, err := this.Pack(msgID, data)
	if err == nil{
		this.conn.Write(dd)
	}else{
		fmt.Println(err)
	}

}

func (this *TcpClient)Start(){
	go func() {
		for {
		//read per head data
		headdata := make([]byte, 8)

		if _, err := io.ReadFull(this.conn, headdata); err != nil {
			fmt.Println(err)
			this.ConnectionLost()
			return
		}
		pkgHead, err := this.Unpack(headdata)
		if err != nil {
			this.ConnectionLost()
			return
		}
		//data
		if pkgHead.Len > 0 {
			pkgHead.Data = make([]byte, pkgHead.Len)
			if _, err := io.ReadFull(this.conn, pkgHead.Data); err != nil {
				this.ConnectionLost()
				return
			}
		}
		this.DoMsg(pkgHead)
	}
	}()
}

func main() {
	for i := 0; i< 2; i ++{
		client := NewTcpClient("0.0.0.0", 8909)
		client.Start()
		time.Sleep(1*time.Second)
	}

	// close
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	sig := <-c
	fmt.Println("=======", sig)
}
