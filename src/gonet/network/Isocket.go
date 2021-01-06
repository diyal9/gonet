package network

import (
	"fmt"
	"gonet/base"
	"gonet/base/vector"
	"gonet/rpc"
	"net"
)

const (
	SSF_ACCEPT    		 = iota
	SSF_CONNECT    	 = iota
	SSF_SHUT_DOWN      = iota //已经关闭
)

const (
	CLIENT_CONNECT = iota//对外
	SERVER_CONNECT = iota//对内
)

const (
	MAX_SEND_CHAN = 100
)

type (
	HandleFunc func(uint32,[]byte) bool//回调函数
	Socket struct {
		m_Conn                 net.Conn
		m_nPort                int
		m_sIP                  string
		m_nState			   int
		m_nConnectType		   int
		m_ReceiveBufferSize    int//单次接受缓存
		m_MaxReceiveBufferSize int//最大接受缓存

		m_ClientId uint32
		m_Seq      int64

		m_TotalNum     int
		m_AcceptedNum  int
		m_ConnectedNum int

		m_SendTimes     int
		m_ReceiveTimes  int
		m_bShuttingDown bool
		m_PacketFuncList	*vector.Vector//call back

		m_bHalf		bool
		m_nHalfSize int
		m_MaxReceiveBuffer []byte//max receive buff
	}

	ISocket interface {
		Init(string, int) bool
		Start() bool
		Stop() bool
		Run() bool
		Restart() bool
		Connect() bool
		Disconnect(bool) bool
		OnNetFail(int)
		Clear()
		Close()
		SendMsg(rpc.RpcHead, string, ...interface{})
		Send(rpc.RpcHead, []byte) int
		CallMsg(string, ...interface{})//回调消息处理

		GetState() int
		SetReceiveBufferSize(int)
		GetReceiveBufferSize()int
		SetMaxReceiveBufferSize(int)
		GetMaxReceiveBufferSize()int
		BindPacketFunc(HandleFunc)
		SetConnectType(int)
		SetTcpConn(net.Conn)
		HandlePacket(uint32,	[]byte)
		ReceivePacket(uint32,	[]byte) bool
	}
)

// virtual
func (this *Socket) Init(string, int) bool {
	this.m_PacketFuncList = vector.NewVector()
	this.m_nState = SSF_SHUT_DOWN
	this.m_ReceiveBufferSize =1024
	this.m_MaxReceiveBufferSize = base.MAX_PACKET
	this.m_nConnectType = SERVER_CONNECT
	this.m_bHalf = false
	this.m_nHalfSize = 0
	return true
}

func (this *Socket) Start() bool {
	return true
}

func (this *Socket) Stop() bool {
	this.m_bShuttingDown = true
	return true
}

func (this *Socket) Run()bool {
	return true
}

func (this *Socket) Restart() bool {
	return true
}

func (this *Socket) Connect() bool {
	return true
}

func (this *Socket) Disconnect(bool) bool {
	return true
}

func (this *Socket) OnNetFail(int) {
	this.Stop()
}

func (this *Socket) GetState() int{
	return  this.m_nState
}

func (this *Socket) SendMsg(head rpc.RpcHead, funcName string, params  ...interface{}){
}

func (this *Socket) Send(rpc.RpcHead, []byte) int{
	return  0
}

func (this *Socket) Clear() {
	this.m_nState = SSF_SHUT_DOWN
	//this.m_nConnectType = CLIENT_CONNECT
	this.m_Conn = nil
	this.m_ReceiveBufferSize = 1024
	this.m_MaxReceiveBufferSize = base.MAX_PACKET
	this.m_bShuttingDown = false
	this.m_bHalf = false
	this.m_nHalfSize = 0
}

func (this *Socket) Close() {
	if this.m_Conn != nil{
		this.m_Conn.Close()
	}
	this.Clear()
}

func (this *Socket) GetMaxReceiveBufferSize() int{
	return  this.m_MaxReceiveBufferSize
}

func (this *Socket) SetMaxReceiveBufferSize(maxReceiveSize int){
	this.m_MaxReceiveBufferSize = maxReceiveSize
}

func (this *Socket) GetReceiveBufferSize() int{
	return  this.m_ReceiveBufferSize
}

func (this *Socket) SetReceiveBufferSize(maxSendSize int){
	this.m_ReceiveBufferSize = maxSendSize
}

func (this *Socket) SetConnectType(nType int){
	this.m_nConnectType = nType
}

func (this *Socket) SetTcpConn(conn net.Conn){
	this.m_Conn = conn
}

func (this *Socket) BindPacketFunc(callfunc HandleFunc){
	this.m_PacketFuncList.PushBack(callfunc)
}

func (this *Socket) CallMsg(funcName string, params ...interface{}){
	buff := rpc.Marshal(rpc.RpcHead{}, funcName, params...)
	this.HandlePacket(this.m_ClientId, buff)
}

func (this *Socket) HandlePacket(Id uint32, buff []byte){
	for _,v := range this.m_PacketFuncList.Values() {
		if (v.(HandleFunc)(Id, buff)){
			break
		}
	}
}

func (this *Socket) ReceivePacket(Id uint32, dat []byte) bool{
	//找包结束
	seekToTcpEnd := func(buff []byte) (bool, int){
		nLen := len(buff)
		if nLen < base.TCP_HEAD_SIZE{
			return false, 0
		}

		nSize := base.BytesToInt(buff[0:4])
		if nSize + base.TCP_HEAD_SIZE <= nLen{
			return true, nSize+base.TCP_HEAD_SIZE
		}
		return false, 0
	}

	buff := append(this.m_MaxReceiveBuffer, dat...)
	this.m_MaxReceiveBuffer = []byte{}
	nCurSize := 0
	//fmt.Println(this.m_MaxReceiveBuffer)
ParsePacekt:
	nPacketSize := 0
	nBufferSize := len(buff[nCurSize:])
	bFindFlag := false
	bFindFlag, nPacketSize = seekToTcpEnd(buff[nCurSize:])
	//fmt.Println(bFindFlag, nPacketSize, nBufferSize)
	if bFindFlag{
		if nBufferSize == nPacketSize{		//完整包
			this.HandlePacket(Id, buff[nCurSize+base.TCP_HEAD_SIZE:nCurSize+nPacketSize])
			nCurSize += nPacketSize
		}else if ( nBufferSize > nPacketSize){
			this.HandlePacket(Id, buff[nCurSize+base.TCP_HEAD_SIZE:nCurSize+nPacketSize])
			nCurSize += nPacketSize
			goto ParsePacekt
		}
	}else if nBufferSize < this.m_MaxReceiveBufferSize{
		this.m_MaxReceiveBuffer = buff[nCurSize:]
	}else{
		fmt.Println("超出最大包限制，丢弃该包")
		return false
	}
	return true
}

//tcp粘包特殊结束标志
/*func (this *Socket) ReceivePacket(Id int, dat []byte) bool{
	//找包结束
	seekToTcpEnd := func(buff []byte) (bool, int) {
		nLen := bytes.Index(buff, []byte(base.TCP_END))
		if nLen != -1{
			return true, nLen+base.TCP_END_LENGTH
		}
		return false, 0
	}

	buff := append(this.m_MaxReceiveBuffer, dat...)
	this.m_MaxReceiveBuffer = []byte{}
	nCurSize := 0
	//fmt.Println(this.m_MaxReceiveBuffer)
ParsePacekt:
	nPacketSize := 0
	nBufferSize := len(buff[nCurSize:])
	bFindFlag := false
	bFindFlag, nPacketSize = seekToTcpEnd(buff[nCurSize:])
	//fmt.Println(bFindFlag, nPacketSize, nBufferSize)
	if bFindFlag {
		if nBufferSize == nPacketSize { //完整包
			this.HandlePacket(Id, buff[nCurSize:nCurSize+nPacketSize-base.TCP_END_LENGTH])
			nCurSize += nPacketSize
		} else if (nBufferSize > nPacketSize) {
			this.HandlePacket(Id, buff[nCurSize:nCurSize+nPacketSize-base.TCP_END_LENGTH])
			nCurSize += nPacketSize
			goto ParsePacekt
		}
	}else if nBufferSize < this.m_MaxReceiveBufferSize{
		this.m_MaxReceiveBuffer = buff[nCurSize:]
	}else{
		fmt.Println("超出最大包限制，丢弃该包")
		return false
	}
	return true
}*/