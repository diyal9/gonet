package netgate

/*import (
	"gonet/actor"
	"gonet/base"
	"strings"
)

type (
	UserServerProcess struct {
		actor.Actor
	}
)

func (this *UserServerProcess) Init(num int) {
	this.Actor.Init(num)
	this.RegisterCall("DISCONNECT", func(ctx context.Context, socketId uint32) {
		SERVER.GetPlayerMgr().SendMsg("DEL_ACCOUNT", socketId)
	})

	this.Actor.Start()
}

func (this *UserServerProcess) PacketFunc(id int, buff []byte) bool{
	var io actor.CallIO
	io.Buff = buff
	io.SocketId = id

	rpcPacket, head := rpc.UnmarshalHead(io.Buff)
	if this.FindCall(rpcPacket.FuncName) != nil{
		this.Send(io)
		return true
	}

	return false
}*/