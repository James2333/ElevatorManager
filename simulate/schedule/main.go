package main

import (
	"EleManager/public/code"
	"EleManager/public/packet"
	"EleManager/public/session"
	"EleManager/public/task"
	"EleManager/schedule"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
)

func main() {
	NewClientSocket()
}
func NewClientSocket() {
	conn, err := net.Dial("tcp", "127.0.0.1:20000")
	if err != nil {
		fmt.Println("客户端建立连接失败")
		return
	}
	log.Println("客户端建立连接成功...")
	cConnHandler(conn)
}

func cConnHandler(c net.Conn) {
	//这个client模拟的调度请求电梯操作
	if c == nil {
		log.Println("conn无效")
		return
	}
	in := make(chan []byte, 16)
	sess := session.NewSession(c, in)
	//连接时注册
	schedule.ReqRegister(sess)
	defer func() {
		log.Println("disconnect", c.RemoteAddr().String())
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

	//模拟发起一次任务
	task1 := &task.Task{
		TaskID: "111",
		Start:  3,
		End:    5,
	}
	NewTestTask(sess, task1)
	for {
		//此处应该先 解包识别byte[0:2]的code 然后去传入 不同的方法。
		head := make([]byte, packet.HEADER_LEN)
		_, err := io.ReadFull(c, head) //读取头部的2个字节
		if err != nil {
			log.Println(err)
		}
		cod := binary.BigEndian.Uint16(head)
		schedule.ParseCodeScheduling(cod, sess)
	}
}

func NewTestTask(s *session.Session, task *task.Task) {
	b := packet.Packet(task, code.CHOOSE_ELE)
	s.Ch <- b
}
