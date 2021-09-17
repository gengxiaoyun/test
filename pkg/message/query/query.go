package query

import (
	"github.com/romberli/das/pkg/message"
	"github.com/romberli/go-util/config"
)

func init() {
	initQueryDebugMessage()
	initQueryInfoMessage()
	initQueryErrorMessage()
}

const (
	// debug
	DebugQueryGetByMySQLClusterID = 103001
	DebugQueryGetByMySQLServerID  = 103002
	DebugQueryGetByDBID           = 103003
	DebugQueryGetBySQLID          = 103004
	// info
	InfoQueryGetByMySQLClusterID = 203001
	InfoQueryGetByMySQLServerID  = 203002
	InfoQueryGetByDBID           = 203003
	InfoQueryGetBySQLID          = 203004
	// error
	ErrQueryGetByMySQLClusterID = 403001
	ErrQueryGetByMySQLServerID  = 403002
	ErrQueryGetByDBID           = 403003
	ErrQueryGetBySQLID          = 403004
	ErrQueryConfigNotValid      = 403005
	ErrQueryMonitorSystemType   = 403006
	ErrQueryCloseMonitorRepo    = 403007
)

func initQueryDebugMessage() {
	message.Messages[DebugQueryGetByMySQLClusterID] = config.NewErrMessage(message.DefaultMessageHeader, DebugQueryGetByMySQLClusterID, "get by mysql cluster id completed. mysql_cluster_id: %d.\n%s")
	message.Messages[DebugQueryGetByMySQLServerID] = config.NewErrMessage(message.DefaultMessageHeader, DebugQueryGetByMySQLServerID, "get by mysql server id completed. mysql_server_id: %d.\n%s")
	message.Messages[DebugQueryGetByDBID] = config.NewErrMessage(message.DefaultMessageHeader, DebugQueryGetByDBID, "get by db id completed. db_id: %d.\n%s")
	message.Messages[DebugQueryGetBySQLID] = config.NewErrMessage(message.DefaultMessageHeader, DebugQueryGetBySQLID, "get by sql id completed. mysql_server_id: %d, sql_id: %s.\n%s")
}

func initQueryInfoMessage() {
	message.Messages[InfoQueryGetByMySQLClusterID] = config.NewErrMessage(message.DefaultMessageHeader, InfoQueryGetByMySQLClusterID, "get by mysql cluster id completed. mysql_cluster_id: %d.")
	message.Messages[InfoQueryGetByMySQLServerID] = config.NewErrMessage(message.DefaultMessageHeader, InfoQueryGetByMySQLServerID, "get by mysql server id completed. mysql_server_id: %d.")
	message.Messages[InfoQueryGetByDBID] = config.NewErrMessage(message.DefaultMessageHeader, InfoQueryGetByDBID, "get by db id completed. db_id: %d.")
	message.Messages[InfoQueryGetBySQLID] = config.NewErrMessage(message.DefaultMessageHeader, InfoQueryGetBySQLID, "get by sql id completed. mysql_server_id: %d, sql_id: %s.")
}

func initQueryErrorMessage() {
	message.Messages[ErrQueryGetByMySQLClusterID] = config.NewErrMessage(message.DefaultMessageHeader, ErrQueryGetByMySQLClusterID, "get by mysql cluster id failed. mysql_cluster_id: %d.\n%s")
	message.Messages[ErrQueryGetByMySQLServerID] = config.NewErrMessage(message.DefaultMessageHeader, ErrQueryGetByMySQLServerID, "get by mysql server id failed. mysql_server_id: %d.\n%s")
	message.Messages[ErrQueryGetByDBID] = config.NewErrMessage(message.DefaultMessageHeader, ErrQueryGetByDBID, "get by db id failed. db_id: %d.\n%s")
	message.Messages[ErrQueryGetBySQLID] = config.NewErrMessage(message.DefaultMessageHeader, ErrQueryGetBySQLID, "get by sql id failed. mysql_server_id: %d, sql_id: %s.\n%s")
	message.Messages[ErrQueryConfigNotValid] = config.NewErrMessage(message.DefaultMessageHeader, ErrQueryConfigNotValid, "config is not valid. start_time: %s, end_time: %s, limit: %d")
	message.Messages[ErrQueryMonitorSystemType] = config.NewErrMessage(message.DefaultMessageHeader, ErrQueryMonitorSystemType, "monitor system type version should be either 1 or 2, %d is not valid")
	message.Messages[ErrQueryCloseMonitorRepo] = config.NewErrMessage(message.DefaultMessageHeader, ErrQueryCloseMonitorRepo, "close monitor repo failed.\n%s")
}
