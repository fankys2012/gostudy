package utils

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"net"

	"github.com/fankys2012/gostudy/chatroom/common/message"
)

type Transfer struct {
	Conn net.Conn
	Buf  [8096]byte //缓冲
}

func (this *Transfer) ReadPkg() (mes message.Message, err error) {
	//conn.Read 只有在conn 没有被关闭的情况才会一直阻塞
	//如果客户端关闭了，则
	_, err = this.Conn.Read(this.Buf[:4])
	//&& err != io.EOF
	if err != nil{
		// err = errors.New("read pkg header error")
		if err != io.EOF {
			fmt.Println("读取消息长度失败. err :",err)
		}
		return
	}

	//根据buf[:4] 转成uint32
	pkgLen := binary.BigEndian.Uint32(this.Buf[0:4])

	//根据pkgLen 读取消息内容  从conn 中读取内容存入buf中
	n, err := this.Conn.Read(this.Buf[:pkgLen])
	if uint32(n) != pkgLen || err != nil {
		fmt.Println("读取消息失败，err：", err)
		// err = errors.New("read message body failed")
		return
	}

	//反序列化message  mes必须传地址; mes 在返回参数中已声明，不用重复声明
	err = json.Unmarshal(this.Buf[:pkgLen], &mes)
	if err != nil {
		fmt.Println("解析消息失败,err：", err)
	}
	return
}

func (this *Transfer) WritePkg(data []byte) (err error) {
	//发送消息内容长度
	// 先获取 data 的长度-> 转换成一个表示长度的byte切片
	var pkgLen uint32
	pkgLen = uint32(len(data))
	// var buf [4]byte
	binary.BigEndian.PutUint32(this.Buf[0:4], pkgLen) //len -> byte
	//发送长度
	n, err := this.Conn.Write(this.Buf[:4])
	if n != 4 || err != nil {
		fmt.Println("字符长度发送失败 ", err)
		return
	}

	//发送消息体内容
	//发送内容
	n, err = this.Conn.Write(data)
	if uint32(n) != pkgLen || err != nil {
		fmt.Println("消息内容发送失败 ", err)
		return
	}
	return

}
