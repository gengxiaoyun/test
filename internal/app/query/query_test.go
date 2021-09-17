package query

import (
	"fmt"
	"testing"

	"github.com/romberli/das/config"
	"github.com/romberli/go-util/common"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

const (
	// modify the connection information
	queryTestDBAddr   = "192.168.10.219:3306"
	queryTestDBDBName = "das"
	queryTestDBDBUser = "root"
	queryTestDBDBPass = "root"
)

const (
	defaultQueryInfoSQLID           = "sql_id"
	defaultQueryInfoFingerprint     = "fingerprint"
	defaultQueryInfoExample         = "example"
	defaultQueryInfoDBName          = "db"
	defaultQueryInfoExecCount       = 1
	defaultQueryInfoTotalExecTime   = 2.1
	defaultQueryInfoAvgExecTime     = 3.2
	defaultQueryInfoRowsExaminedMax = 4

	defaultQuerierPMM1MySQLClusterID = 1
	defaultQuerierPMM1MySQLServerID  = 2
	defaultQuerierPMM1DBID           = 1
	defaultQuerierPMM1SQLID          = "999ECD050D719733"

	defaultQuerierPMM2MySQLClusterID = 3
	defaultQuerierPMM2MySQLServerID  = 4
	defaultQuerierPMM2DBID           = 2
	defaultQuerierPMM2SQLID          = "999ECD050D719733"
)

//var pmmVersion = 0

func init() {
	viper.Set(config.DBMonitorMySQLUserKey, config.DefaultDBMonitorMySQLUser)
	viper.Set(config.DBMonitorMySQLPassKey, config.DefaultDBMonitorMySQLPass)

	// viper.Set(config.DBMonitorClickhouseUserKey, config.DefaultDBMonitorClickhouseUser)
	// viper.Set(config.DBMonitorClickhousePassKey, config.DefaultDBMonitorClickhousePass)

	viper.Set(config.DBMonitorClickhouseUserKey, "")
	viper.Set(config.DBMonitorClickhousePassKey, "")

	// pmmVersion = 1
	pmmVersion = 2
	//
	//if err := initGlobalMySQLPool(); err != nil {
	//	panic(err)
	//}
}

//func initGlobalMySQLPool() error {
//	dbAddr := queryTestDBAddr
//	dbName := queryTestDBDBName
//	dbUser := queryTestDBDBUser
//	dbPass := queryTestDBDBPass
//	maxConnections := mysql.DefaultMaxConnections
//	initConnections := mysql.DefaultInitConnections
//	maxIdleConnections := mysql.DefaultMaxIdleConnections
//	maxIdleTime := mysql.DefaultMaxIdleTime
//	keepAliveInterval := mysql.DefaultKeepAliveInterval
//
//	config := mysql.NewConfig(dbAddr, dbName, dbUser, dbPass)
//	poolConfig := mysql.NewPoolConfigWithConfig(config, maxConnections, initConnections, maxIdleConnections, maxIdleTime, keepAliveInterval)
//	log.Debugf("pool config: %v", poolConfig)
//	var err error
//	global.DASMySQLPool, err = mysql.NewPoolWithPoolConfig(poolConfig)
//
//	return err
//}

func TestQueryAll(t *testing.T) {
	TestQuery_GetSQLID(t)
	TestQuery_GetFingerprint(t)
	TestQuery_GetExample(t)
	TestQuery_GetDBName(t)
	TestQuery_GetExecCount(t)
	TestQuery_GetTotalExecTime(t)
	TestQuery_GetAvgExecTime(t)
	TestQuery_GetRowsExaminedMax(t)

	// Test PMM1.x
	pmmVersion = 1
	TestQuerier_GetByMySQLClusterID(t)
	TestQuerier_GetByMySQLServerID(t)
	TestQuerier_GetByDBID(t)
	TestQuerier_GetBySQLID(t)

	// Test PMM2.x
	pmmVersion = 2
	TestQuerier_GetByMySQLClusterID(t)
	TestQuerier_GetByMySQLServerID(t)
	TestQuerier_GetByDBID(t)
	TestQuerier_GetBySQLID(t)
}

func initNewQueryInfo() *Query {
	return &Query{
		defaultQueryInfoSQLID,
		defaultQueryInfoFingerprint,
		defaultQueryInfoExample,
		defaultQueryInfoDBName,
		defaultQueryInfoExecCount,
		defaultQueryInfoTotalExecTime,
		defaultQueryInfoAvgExecTime,
		defaultQueryInfoRowsExaminedMax,
	}
}

func TestQuery_GetSQLID(t *testing.T) {
	asst := assert.New(t)

	queryInfo := initNewQueryInfo()
	asst.Equal(defaultQueryInfoSQLID, queryInfo.GetSQLID(), "test GetUserName() failed")
}
func TestQuery_GetFingerprint(t *testing.T) {
	asst := assert.New(t)

	queryInfo := initNewQueryInfo()
	asst.Equal(defaultQueryInfoFingerprint, queryInfo.GetFingerprint(), "test GetUserName() failed")
}
func TestQuery_GetExample(t *testing.T) {
	asst := assert.New(t)

	queryInfo := initNewQueryInfo()
	asst.Equal(defaultQueryInfoExample, queryInfo.GetExample(), "test GetUserName() failed")
}
func TestQuery_GetDBName(t *testing.T) {
	asst := assert.New(t)

	queryInfo := initNewQueryInfo()
	asst.Equal(defaultQueryInfoDBName, queryInfo.GetDBName(), "test GetUserName() failed")
}
func TestQuery_GetExecCount(t *testing.T) {
	asst := assert.New(t)

	queryInfo := initNewQueryInfo()
	asst.Equal(defaultQueryInfoExecCount, queryInfo.GetExecCount(), "test GetUserName() failed")
}
func TestQuery_GetTotalExecTime(t *testing.T) {
	asst := assert.New(t)

	queryInfo := initNewQueryInfo()
	asst.Equal(defaultQueryInfoTotalExecTime, queryInfo.GetTotalExecTime(), "test GetUserName() failed")
}
func TestQuery_GetAvgExecTime(t *testing.T) {
	asst := assert.New(t)

	queryInfo := initNewQueryInfo()
	asst.Equal(defaultQueryInfoAvgExecTime, queryInfo.GetAvgExecTime(), "test GetUserName() failed")
}
func TestQuery_GetRowsExaminedMax(t *testing.T) {
	asst := assert.New(t)

	queryInfo := initNewQueryInfo()
	asst.Equal(defaultQueryInfoRowsExaminedMax, queryInfo.GetRowsExaminedMax(), "test GetUserName() failed")
}

func TestQuerier_GetByMySQLClusterID(t *testing.T) {
	asst := assert.New(t)

	querier := NewQuerierWithGlobal(NewConfigWithDefault())

	querierMySQLClusterID := 0
	switch pmmVersion {
	case 1:
		querierMySQLClusterID = defaultQuerierPMM1MySQLClusterID
	case 2:
		querierMySQLClusterID = defaultQuerierPMM2MySQLClusterID
	default:
		err := fmt.Errorf("PMM with version:%d is not supported for now", pmmVersion)
		asst.Nil(err, common.CombineMessageWithError("test GetByMySQLClusterID() failed", err))
	}

	queries, err := querier.GetByMySQLClusterID(querierMySQLClusterID)
	asst.Nil(err, common.CombineMessageWithError("test GetByMySQLClusterID() failed", err))

	asst.NotZero(len(queries), "test GetByMySQLClusterID() failed")
}

func TestQuerier_GetByMySQLServerID(t *testing.T) {
	asst := assert.New(t)

	querierMySQLServerID := 0
	switch pmmVersion {
	case 1:
		querierMySQLServerID = defaultQuerierPMM1MySQLServerID
	case 2:
		querierMySQLServerID = defaultQuerierPMM2MySQLServerID
	default:
		err := fmt.Errorf("PMM with version:%d is not supported for now", pmmVersion)
		asst.Nil(err, common.CombineMessageWithError("test GetByMySQLServerID() failed", err))
	}

	querier := NewQuerierWithGlobal(NewConfigWithDefault())
	queries, err := querier.GetByMySQLServerID(querierMySQLServerID)
	asst.Nil(err, common.CombineMessageWithError("test GetByMySQLServerID() failed", err))

	asst.NotZero(len(queries), "test GetByMySQLServerID() failed")
}

func TestQuerier_GetByDBID(t *testing.T) {
	asst := assert.New(t)

	querierDBID := 0
	querierMySQLServerID := 0
	switch pmmVersion {
	case 1:
		querierMySQLServerID = defaultQuerierPMM1MySQLServerID
		querierDBID = defaultQuerierPMM1DBID
	case 2:
		querierMySQLServerID = defaultQuerierPMM1MySQLServerID
		querierDBID = defaultQuerierPMM2DBID
	default:
		err := fmt.Errorf("PMM with version:%d is not supported for now", pmmVersion)
		asst.Nil(err, common.CombineMessageWithError("test GetByDBID() failed", err))
	}

	querier := NewQuerierWithGlobal(NewConfigWithDefault())
	queries, err := querier.GetByDBID(querierMySQLServerID, querierDBID)
	asst.Nil(err, common.CombineMessageWithError("test GetByDBID() failed", err))

	asst.NotZero(len(queries), "test GetByDBID() failed")
}
func TestQuerier_GetBySQLID(t *testing.T) {
	asst := assert.New(t)

	querierSQLID := ""
	querierMySQLServerID := 0
	switch pmmVersion {
	case 1:
		querierMySQLServerID = defaultQuerierPMM1MySQLServerID
		querierSQLID = defaultQuerierPMM1SQLID
	case 2:
		querierMySQLServerID = defaultQuerierPMM1MySQLServerID
		querierSQLID = defaultQuerierPMM2SQLID
	default:
		err := fmt.Errorf("PMM with version:%d is not supported for now", pmmVersion)
		asst.Nil(err, common.CombineMessageWithError("test GetByDBID() failed", err))
	}

	querier := NewQuerierWithGlobal(NewConfigWithDefault())
	queries, err := querier.GetBySQLID(querierMySQLServerID, querierSQLID)
	asst.Nil(err, common.CombineMessageWithError("test GetBySQLID() failed", err))

	asst.NotZero(len(queries), "test GetBySQLID() failed")
}
