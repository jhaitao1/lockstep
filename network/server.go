package network

import (
	"fmt"

	"github.com/xtaci/kcp-go/v5"
)

//Run 开始运行server
func Run(address string) {
	fmt.Println("star run")
	lis, err := kcp.Listen(":8008")
	if err != nil {
		panic("监听端口失败")
	}
	for {
		conn, err := lis.Accept()
		if err != nil {
			continue
		}
		//创建conn单独开一个goroutine去完成
		go func() {
			NewConn(conn).Work()
		}()
	}
}
