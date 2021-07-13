package tcpService

import (
	"EleManager/eleManager"
	"EleManager/public/elevator"
	"EleManager/public/packet"
	"EleManager/public/session"
	"encoding/binary"
	"fmt"
	"github.com/golang/glog"
	"io"
	"log"
	"net"
	"time"
)

func NewTcpService() {
	//1.建立监听端口
	listen, err := net.Listen("tcp", "0.0.0.0:20000")
	if err != nil {
		fmt.Println("listen failed, err:", err)
		return
	}
	fmt.Println("listen Start...:")
	//els := NewTestEls()
	fmt.Println("初始化电梯数据...")
	go func() {
		for  {
			//log.Println("打印电梯信息:")
			for _,j:=range elevator.Els {
				log.Printf("打印电梯信息:ElevatorId:%s,Floor:%d,State:%s,CurrentState:%s,IsInFloor:%t\n",j.ElevatorId,j.Floor,j.State,j.CurrentState,j.IsInFloor)
			}
			time.Sleep(time.Second*10)
		}
	}()
	for {
		//2.接收客户端的链接
		conn, err := listen.Accept()
		if err != nil {
			fmt.Printf("accept failed, err:%v\n", err)
			continue
		}
		//3.开启一个Goroutine，处理链接
		go connSt(conn)
	}
}

func connSt(c net.Conn)  {
	in := make(chan []byte, 1024)
	sess := session.NewSession(c,in)
	defer func() {
		glog.Info("disconnect:" + c.RemoteAddr().String())
		c.Close()
	}()
	go func() {
		for {
			select {
			case msg := <-in:
				c.Write(msg)
			}
		}
	}()
	for {
		//此处应该先 解包识别byte[0:2]的code 然后去传入 不同的方法。
		head := make([]byte, packet.HEADER_LEN)
		_, err := io.ReadFull(c, head) //读取头部的2个字节
		if err != nil {
			log.Println(err)
		}
		code := binary.BigEndian.Uint16(head)
		eleManager.ParseCode(code,sess)
	}
}
