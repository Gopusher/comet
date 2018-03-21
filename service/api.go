package service

import (
	"net/rpc"
	"gopusher/comet/contracts"
	"net"
	"net/rpc/jsonrpc"
	"reflect"
	"encoding/json"
	"errors"
	"github.com/fatih/color"
)

type Server struct {
	server contracts.Server
}

func InitRpcServer(server contracts.Server) {
	rpc.Register(&Server{
		server: server,
	})
	listener, err := net.Listen("tcp", server.GetRpcAddr())
	if err != nil {
		panic("rpc服务初始化失败, " + err.Error())
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		//新协程来处理--json
		go jsonrpc.ServeConn(conn)
	}
}

func (s *Server) SendToConnections(body string, reply *string) error {
	//const messageMaxLen = 200
	//if strings.Count(body, "") - 1 > messageMaxLen {
	//	return errors.New(fmt.Sprintf("消息体过长，最大允许长度: %d", messageMaxLen))
	//}

	type Message struct {
		To   	[]string	`json:"to"`	//消息接受者
		Msg 	string		`json:"msg"` //为一个json，里边包含 type 消息类型
	}

	var message Message
	if err := json.Unmarshal([]byte(body), &message); err != nil {
		color.Red("消息体异常, 不能解析 %v %v", body, reflect.TypeOf(body))
		return errors.New("消息体异常, 不能解析")
	}

	if err := s.server.SendToConnections(message.To, message.Msg); err != nil {
		return errors.New("消息发送失败" + err.Error())
	}

	*reply = "消息发送成功"
	return nil
}
