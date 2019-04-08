package utils

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"

	"github.com/fankys2012/gostudy/chatroom/common/message"
)

type Transfer struct {
	conn net.Conn
	Buf  [8096]byte //缓冲
}

func (this *Transfer) ReadPkg(conn net.Conn) (mes message.Message, err error) {
	buf := make([]byte, 8096)
	//conn.Read 只有在conn 没有被关闭的情况才会一直阻塞
	//如果客户端关闭了，则
	_, err = conn.Read(buf[:4])
	if err != nil {
		fmt.Printf("数据读取失败err =%s", err)
		// err = errors.New("read pkg header error")
		return
	}

	//根据buf[:4] 转成uint32
	pkgLen := binary.BigEndian.Uint32(buf[0:4])

	//根据pkgLen 读取消息内容  从conn 中读取内容存入buf中
	n, err := conn.Read(buf[:pkgLen])
	if uint32(n) != pkgLen || err != nil {
		fmt.Println("读取消息失败,err ", err)
		// err = errors.New("read message body failed")
		return
	}

	//反序列化message  mes必须传地址; mes 在返回参数中已声明，不用重复声明
	err = json.Unmarshal(buf[:pkgLen], &mes)
	if err != nil {
		fmt.Println("消息解析失败 ", err)
	}
	return
}

func (this *Transfer) WritePkg(conn net.Conn, data []byte) (err error) {
	//发送消息内容长度
	// 先获取 data 的长度-> 转换成一个表示长度的byte切片
	var pkgLen uint32
	pkgLen = uint32(len(data))
	var buf [4]byte
	binary.BigEndian.PutUint32(buf[0:4], pkgLen) //len -> byte
	//发送长度
	n, err := conn.Write(buf[:4])
	if n != 4 || err != nil {
		fmt.Println("字符长度发送失败 ", err)
		return
	}

	//发送消息体内容
	//发送内容
	n, err = conn.Write(data)
	if uint32(n) != pkgLen || err != nil {
		fmt.Println("消息内容发送失败 ", err)
		return
	}
	return

}
