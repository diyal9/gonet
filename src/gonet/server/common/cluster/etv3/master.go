package etv3

import (
	"encoding/json"
	"go.etcd.io/etcd/clientv3"
	"gonet/actor"
	"gonet/rpc"
	"gonet/server/common"
	"log"

	"golang.org/x/net/context"
)

//监控服务器
type Master struct {
	m_ServiceMap map[uint32]*common.ClusterInfo
	m_Client *clientv3.Client
	m_Actor actor.IActor
	*common.ClusterInfo
}

//监控服务器
func (this *Master) Init(info *common.ClusterInfo, Endpoints []string, pActor actor.IActor) {
	cfg := clientv3.Config{
		Endpoints:               Endpoints,
	}

	etcdClient, err := clientv3.New(cfg)
	if err != nil {
		log.Fatal("Error: cannot connec to etcd:", err)
	}

	this.m_ServiceMap = make(map[uint32]*common.ClusterInfo)
	this.m_Client = etcdClient
	this.BindActor(pActor)
	this.Start()
	this.ClusterInfo = info
}

func (this *Master) Start() {
	go this.Run()
}

func (this *Master) BindActor(pActor actor.IActor) {
	this.m_Actor = pActor
}

func (this *Master) AddService(info *common.ClusterInfo) {
	_, bEx := this.m_ServiceMap[info.Id()]
	if !bEx{
		this.m_Actor.SendMsg(rpc.RpcHead{},"Cluster_Add", info)
	}
	this.m_ServiceMap[info.Id()] = info
}

func (this *Master) DelService(info *common.ClusterInfo) {
	delete(this.m_ServiceMap, info.Id())
	this.m_Actor.SendMsg(rpc.RpcHead{},"Cluster_Del", info)
}

func (this *Master) InitService(info *common.ClusterInfo) {
}

func NodeToService(val []byte) *common.ClusterInfo {
	info := &common.ClusterInfo{}
	err := json.Unmarshal([]byte(val), info)
	if err != nil {
		log.Print(err)
	}
	return info
}

func (this *Master) Run() {
	wch := this.m_Client.Watch(context.Background(), ETCD_DIR + this.String(), clientv3.WithPrefix(), clientv3.WithPrevKV())
	for v := range wch{
		for _, v1 := range v.Events{
			if v1.Type.String() == "PUT"{
				info := NodeToService(v1.Kv.Value)
				this.AddService(info)
			}else {
				info := NodeToService(v1.PrevKv.Value)
				this.DelService(info)
			}
		}
	}
}