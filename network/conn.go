package network

import (
	"fmt"
	"net"
	"time"
)

//Packet 接口
type Packet interface {
	Serialize() []byte
}

//Conn 对net.Conn的封装
type Conn struct {
	conn      net.Conn    //实际的链接
	closeChan chan byte   //关闭通道
	readChan  chan Packet //读通道
	writeChan chan Packet //读通道
	isClosed  bool        //连接是否已经关闭了
}

//NewConn 创建一个Conn
func NewConn(conn net.Conn) *Conn {
	return &Conn{
		conn:      conn,
		closeChan: make(chan byte),
		readChan:  make(chan Packet, 10), //缓冲大小后期可以考虑以参数方式传进
		writeChan: make(chan Packet, 10),
	}
}

//IsClosed 是否连接已经关了
func (c *Conn) IsClosed() bool {
	return c.isClosed
}

func (c *Conn) readPacket() {
}

//readLoop 读字节后打包发给readChan
func (c *Conn) readLoop() {
	defer func() {
		err := recover()
		if err != nil {
			fmt.Println(err)
		}
		c.conn.Close()
	}()

	for {
		select {
		case <-c.closeChan:
			return
		default:
		}
		c.conn.SetReadDeadline(time.Now().Add(5 * time.Second)) //设置超时时间
		var bytes []byte
		_, err := c.conn.Read(bytes)
		if err != nil {
			fmt.Println("c.conn.Read Error:", err)
			return
		}
		//这部分是把数据打成Packet然后放到readChan上由handleLoop去处理
	}
}

//writeLoop 写
func (c *Conn) writeLoop() {
	defer func() {
		err := recover()
		if err != nil {
			fmt.Println(err)
		}
		c.conn.Close()
	}()

	for {
		select {
		case <-c.closeChan:
			return
		case p, ok := <-c.writeChan:
			if ok {
				if !c.isClosed {
					c.conn.SetWriteDeadline(time.Now().Add(5 * time.Second)) //设置超时时间
					_, err := c.conn.Write(p.Serialize())
					if err != nil {
						fmt.Println("c.conn.Write Error:", err)
						return
					}
				} else {
					fmt.Println("连接已经关闭了")
					return
				}
			}
		}
	}
}

//handleLoop 处理数据
func (c *Conn) handleLoop() {
	defer func() {
		err := recover()
		if err != nil {
			fmt.Println(err)
		}
		c.conn.Close() //这三个loop直接这么关闭会有隐患这里到时会封装一个函数sync.once确保只关闭一次
	}()

	for {
		select {
		case <-c.closeChan:
			return
		case p, ok := <-c.readChan:
			if ok {
				//有数据之后就就行消息的分发反序列化处理
			}
		}
	}
}

func asyncWork(fn func()) {
	go func() {
		fn()
	}()
}

//Work 主要工作是异步读写, 数据都通过channel传递
func (c *Conn) Work() {
	asyncWork(c.readLoop)
	asyncWork(c.writeLoop)
	asyncWork(c.handleLoop)
}
