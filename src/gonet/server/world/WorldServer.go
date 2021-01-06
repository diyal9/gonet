package world

import (
	"database/sql"
	"gonet/base"
	"gonet/db"
	"gonet/network"
	"gonet/rd"
	"gonet/server/common/cluster"
	"log"
)

type(
	ServerMgr struct{
		m_pService	*network.ServerSocket
		m_pClusterMgr *ClusterManager
		m_pActorDB *sql.DB
		m_Inited bool
		m_config base.Config
		m_Log	base.CLog
		m_SnowFlake *cluster.Snowflake
	}

	IServerMgr interface{
		Init() bool
		InitDB() bool
		GetDB() *sql.DB
		GetLog() *base.CLog
		GetServer() *network.ServerSocket
		GetClusterMgr() *ClusterManager
	}
)

var(
	UserNetIP string
	UserNetPort string
	DB_Server string
	DB_Name string
	DB_UserId string
	DB_Password string
	//Web_Url string
	SERVER ServerMgr
	RdID int
	OpenRedis bool
	EtcdEndpoints []string
)

func (this *ServerMgr)Init() bool{
	if(this.m_Inited){
		return true
	}

	//test reload file
	/*file := &common.FileMonitor{}
	file.Init(1000)
	file.AddFile("GONET_SERVER.CFG", func() {this.m_config.Read("GONET_SERVER.CFG")})
	file.AddFile(data.SKILL_DATA_NAME, func() {
		data.SKILLDATA.Read()
	})*/

	//初始化log文件
	this.m_Log.Init("world")
	//初始ini配置文件
	this.m_config.Read("GONET_SERVER.CFG")
	EtcdEndpoints = this.m_config.Get5("Etcd_Cluster", ",")
	UserNetIP, UserNetPort 	= this.m_config.Get2("World_LANAddress", ":")
	DB_Server 	= this.m_config.Get3("WorldDB", "DB_LANIP")
	DB_Name		= this.m_config.Get3("WorldDB","DB_Name");
	DB_UserId	= this.m_config.Get3("WorldDB", "DB_UserId");
	DB_Password	= this.m_config.Get3("WorldDB", "DB_Password")
	RdID 		= 0//this.m_config.Int("WorkID") / 10
	OpenRedis	= this.m_config.Bool("Redis_Open")
	//Web_Url		= this.m_config.Get("World_Url")

	ShowMessage := func(){
		this.m_Log.Println("**********************************************************")
		this.m_Log.Printf("\tWorldServer Version:\t%s",base.BUILD_NO)
		this.m_Log.Printf("\tWorldServerIP(LAN):\t%s:%s", UserNetIP, UserNetPort)
		this.m_Log.Printf("\tActorDBServer(LAN):\t%s", DB_Server)
		this.m_Log.Printf("\tActorDBName:\t\t%s", DB_Name)
		this.m_Log.Println("**********************************************************");
	}
	ShowMessage()

	this.m_Log.Println("正在初始化数据库连接...")
	if (this.InitDB()){
		this.m_Log.Printf("[%s]数据库连接是失败...", DB_Name)
		log.Fatalf("[%s]数据库连接是失败...", DB_Name)
		return false
	}
	this.m_Log.Printf("[%s]数据库初始化成功!", DB_Name)

	if OpenRedis{
		rd.OpenRedisPool(this.m_config.Get("Redis_Host"), this.m_config.Get("Redis_Pwd"))
	}

	//初始化socket
	this.m_pService = new(network.ServerSocket)
	port := base.Int(UserNetPort)
	this.m_pService.Init(UserNetIP, port)
	this.m_pService.Start()

	//snowflake
	this.m_SnowFlake = cluster.NewSnowflake(this.m_config.Get5("Etcd_SnowFlake_Cluster", ","))

	//本身world集群管理
	this.m_pClusterMgr = new(ClusterManager)
	this.m_pClusterMgr.Init(1000)
	this.m_pClusterMgr.BindServer(this.m_pService)

	var packet EventProcess
	packet.Init(1000)
	this.m_pService.BindPacketFunc(packet.PacketFunc)
	this.m_pService.BindPacketFunc(this.m_pClusterMgr.PacketFunc)

	return  false
}

func (this *ServerMgr)InitDB() bool{
	this.m_pActorDB = db.OpenDB(DB_Server, DB_UserId, DB_Password, DB_Name)
	err := this.m_pActorDB.Ping()
	this.m_pActorDB.SetMaxOpenConns(base.Int(this.m_config.Get3("WorldDB", "DB_MaxOpenConns")))
	this.m_pActorDB.SetMaxIdleConns(base.Int(this.m_config.Get3("WorldDB", "DB_MaxIdleConns")))
	return  err != nil
}

func (this *ServerMgr) GetDB() *sql.DB{
	return this.m_pActorDB
}

func (this *ServerMgr) GetLog() *base.CLog{
	return &this.m_Log
}

func (this *ServerMgr) GetServer() *network.ServerSocket{
 	return this.m_pService
}

func (this *ServerMgr) GetClusterMgr() *ClusterManager{
	return this.m_pClusterMgr
}