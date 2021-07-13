package packet

import (
	"encoding/binary"
	"encoding/json"
	"io"
	"log"
	"net"
)

type Protocol struct {
	code    uint16 //表示此次的请求code
	Content []byte //携带的内容
}

const (
	HEADER_LEN = 2
	BODY_LEN   = 1024
)
//打包code+head+content  ==2+2+len(content)
func Packet(content interface{},code int16) []byte {
	bytes,_:=json.Marshal(content)
	buffer := make([]byte, HEADER_LEN+HEADER_LEN+len(bytes))
	// 将buffer前面2个字节设置为code
	binary.BigEndian.PutUint16(buffer[0:2], uint16(code))
	binary.BigEndian.PutUint16(buffer[2:4], uint16(len(bytes)))
	copy(buffer[4:], bytes)
	return buffer
}

//解包，先读取2个字节转换成整形，再读包长度字节
func UnPacket(c net.Conn) (*Protocol, error) {
	var (
		p      = &Protocol{}
	)
	header:=make([]byte,2)
	_, err := io.ReadFull(c, header) //读取包长度的2个字节
	if err != nil {
		return p, err
	}
	length := binary.BigEndian.Uint16(header) //转换成10进制的数字
	log.Println("本次解包长度:",length)
	contentByte := make([]byte,length)
	_, e := io.ReadFull(c, contentByte) //继续读取后续内容
	if e != nil {
		return p, e
	}
	p.Content = contentByte
	return p, nil
}

//func (p *Protocol) parseContent() (map[string]interface{}, error) {
//	var object map[string]interface{}
//	unmarshal := json.Unmarshal(p.Content, &object)
//	if unmarshal != nil {
//		return object, unmarshal
//	}
//	return object, nil
//}