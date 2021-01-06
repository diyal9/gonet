 package account

 import (
	 "database/sql"
	 "gonet/base"
	 "gonet/db"
	 "gonet/network"
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
		m_AccountMgr *AccountMgr
		m_SnowFlake *cluster.Snowflake
	}

	IServerMgr interface{
		Init() bool
		InitDB() bool
		GetDB() *sql.DB
		GetLog() *base.CLog
		GetServer() *network.ServerSocket
		GetClusterMgr() *ClusterManager
		GetAccountMgr() *AccountMgr
	}
)

var(
	UserNetIP string
	UserNetPort string
	WorkID	int
	DB_Server string
	DB_Name string
	DB_UserId string
	DB_Password string
	EtcdEndpoints []string
	SERVER ServerMgr
)

func (this *ServerMgr)Init() bool{
	if(this.m_Inited){
		return true
	}

	//初始化log文件
	this.m_Log.Init("account")
	//初始ini配置文件
	this.m_config.Read("GONET_SERVER.CFG")
	EtcdEndpoints = this.m_config.Get5("Etcd_Cluster", ",")
	UserNetIP, UserNetPort = this.m_config.Get2("Account_LANAddress", ":")
	DB_Server 	= this.m_config.Get3("AccountDB", "DB_LANIP")
	DB_Name		= this.m_config.Get3("AccountDB","DB_Name");
	DB_UserId	= this.m_config.Get3("AccountDB", "DB_UserId");
	DB_Password	= this.m_config.Get3("AccountDB", "DB_Password")

	ShowMessage := func(){
		this.m_Log.Println("**********************************************************")
		this.m_Log.Printf("\tAccountServer Version:\t%s",base.BUILD_NO)
		this.m_Log.Printf("\tAccountServerIP(LAN):\t%s:%s", UserNetIP, UserNetPort)
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

	//初始化socket
	this.m_pService = new(network.ServerSocket)
	port := base.Int(UserNetPort)
	this.m_pService.Init(UserNetIP, port)
	this.m_pService.Start()

	//账号管理类
	this.m_AccountMgr = new(AccountMgr)
	this.m_AccountMgr.Init(1000)

	//本身账号集群管理
	this.m_pClusterMgr = new(ClusterManager)
	this.m_pClusterMgr.Init(1000)
	this.m_pClusterMgr.BindServer(this.m_pService)

	var packet EventProcess
	packet.Init(1000)
	this.m_pService.BindPacketFunc(packet.PacketFunc)
	this.m_pService.BindPacketFunc(this.m_AccountMgr.PacketFunc)
	this.m_pService.BindPacketFunc(this.m_pClusterMgr.PacketFunc)

	//snowflake
	this.m_SnowFlake = cluster.NewSnowflake(this.m_config.Get5("Etcd_SnowFlake_Cluster", ","))

	return  false
}

func (this *ServerMgr)InitDB() bool{
	this.m_pActorDB = db.OpenDB(DB_Server, DB_UserId, DB_Password, DB_Name)
	err := this.m_pActorDB.Ping()
	this.m_pActorDB.SetMaxOpenConns(base.Int(this.m_config.Get3("AccountDB", "DB_MaxOpenConns")))
	this.m_pActorDB.SetMaxIdleConns(base.Int(this.m_config.Get3("AccountDB", "DB_MaxIdleConns")))
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

func (this *ServerMgr) GetAccountMgr() *AccountMgr{
	return this.m_AccountMgr
}