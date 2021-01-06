package account

import (
	"context"
	"gonet/actor"
	"database/sql"
	"gonet/db"
	"fmt"
	"gonet/rpc"
)

type (
	AccountMgr struct {
		actor.Actor

		m_AccountMap map[int64] *Account
		m_AccountNameMap map[string] *Account
		m_db *sql.DB
	}

	IAccountMgr interface {
		actor.IActor

		GetAccount(int64) *Account
		AddAccount(int64) *Account
		RemoveAccount(int64, bool)
		KickAccount(int64)
	}
)

var (
	ACCOUNTMGR AccountMgr
)

func (this* AccountMgr) Init(num int){
	this.m_db = SERVER.GetDB()
	this.Actor.Init(1000)
	this.m_AccountMap = make(map[int64] *Account)
	this.m_AccountNameMap = make(map[string] *Account)
	//this.RegisterTimer(1000 * 1000 * 1000, this.Update)//定时器
	//账号登录处理
	this.RegisterCall("Account_Login", func(ctx context.Context, accountName string, accountId int64, socketId uint32, id uint32) {
		LoginAccount := func(pAccount *Account) {
			if pAccount != nil {
				SERVER.GetLog().Printf("帐号[%s]返回登录OK", accountName)
				SERVER.GetClusterMgr().SendMsg(rpc.RpcHead{SocketId:id}, "A_G_Account_Login", accountId, socketId)
			}
		}

		/*pAccount := this.GetAccount(accountId)
		if pAccount != nil {
			if pAccount.CheckLoginTime(){
				return
			}

			this.RemoveAccount(accountId)
		}*/
		//踢出其他账号服务器
		this.RemoveAccount(accountId, true)
		pAccount := this.AddAccount(accountId)
		LoginAccount(pAccount)
	})

	//账号断开连接
	this.RegisterCall("G_ClientLost", func(ctx context.Context, accountId int64) {
		SERVER.GetLog().Printf("账号[%d] 断开链接", accountId)
		this.RemoveAccount(accountId, false)
	})

	this.Actor.Start()
}

func (this *AccountMgr) GetAccount(accountId int64) *Account{
	pAccount, exist := this.m_AccountMap[accountId]
	if exist{
		return pAccount
	}
	return nil
}

func loadAccount(row db.IRow, a *AccountDB){
	a.AccountId = row.Int64("account_id")
	a.AccountName = row.String("account_name")
	a.LoginIp = row.String("login_ip")
	a.Status = row.Int("status")
	a.LoginTime = row.Time("login_time")
	a.LogoutTime = row.Time("logout_time")
}

func (this *AccountMgr) AddAccount(accountId int64) *Account{
	LoadAccountDB := func(accountId int64) *AccountDB {
		rows, err := this.m_db.Query(fmt.Sprintf("select account_id, account_name, status, login_time, logout_time, login_ip from tbl_account where account_id=%d", accountId))
		rs := db.Query(rows, err)
		if rs.Next() {
			pAccountDB := &AccountDB{}
			pAccountDB.AccountId = accountId
			loadAccount(rs.Row(), pAccountDB)
			return pAccountDB
		}
		return  nil
	}

	pAccountDB  := LoadAccountDB(accountId)
	if pAccountDB != nil{
		pAccount := &Account{}
		pAccount.AccountDB = *pAccountDB
		this.m_AccountMap[accountId] = pAccount
		this.m_AccountNameMap[pAccount.AccountName] = pAccount
		return pAccount
	}

	return nil
}

func (this *AccountMgr) RemoveAccount(accountId int64, bLogin bool){
	pAccount := this.GetAccount(accountId)
	if pAccount != nil{
		delete(this.m_AccountNameMap, pAccount.AccountName)
		delete(this.m_AccountMap, accountId)
		SERVER.GetLog().Printf("账号[%d]断开链接", accountId)
	}
	//假如账号服务器分布式，只要踢出world世界服务器即可
	//这里要登录的时候就同步到踢人world
	if bLogin || pAccount != nil{
		SERVER.GetClusterMgr().KickWorldPlayer(accountId)
	}
}

func (this *AccountMgr) KickAccount(accountId int64){

}
