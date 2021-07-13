package elevator

import (
	"EleManager/public/packet"
	"EleManager/public/session"
	"encoding/json"
	"log"
)

func ReplyTaskEnd(s *session.Session)  {
	q, err := packet.UnPacket(s.C)
	if err != nil {
		log.Println(err)
	}
	ts:=new(string)
	_ = json.Unmarshal(q.Content,ts)
	log.Println(*ts)
}

func ReplyUpdateEle(s *session.Session)  {
	q, err := packet.UnPacket(s.C)
	if err != nil {
		log.Println(err)
	}
	ts:=new(string)
	_ = json.Unmarshal(q.Content,ts)
	log.Println(*ts)
}

func ReplyRegister(s *session.Session)  {
	q, err := packet.UnPacket(s.C)
	if err != nil {
		log.Println(err)
	}
	ts:=new(string)
	_ = json.Unmarshal(q.Content,ts)
	log.Println(*ts)
}
