package query

import (
	"github.com/romberli/das/config"
	"github.com/romberli/das/global"
	"github.com/romberli/go-util/common"
	"github.com/romberli/go-util/middleware/mysql"
	"github.com/romberli/log"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	defaultQuerySQLID           = "sql_id"
	defaultQueryFingerprint     = "fingerprint"
	defaultQueryExample         = "example"
	defaultQueryDBName          = "test"
	defaultQueryExecCount       = 1
	defaultQueryTotalExecTime   = 3.5
	defaultQueryAvgExecTime     = 1.5
	defaultQueryRowsExaminedMax = 10

	// modify the connection information
	// pmm1
	serviceTestDBAddr = "192.168.10.220:3306"
	// pmm2
	//serviceTestDBAddr   = "192.168.10.219:3306"
	serviceTestDBDBName = "das"
	serviceTestDBDBUser = "root"
	serviceTestDBDBPass = "root"
)

var pmmVersion = 0

func init() {
	// pmm1
	pmmVersion = 1
	// pmm2
	//pmmVersion = 2

	switch pmmVersion {
	case 1:
		viper.Set(config.DBMonitorMySQLUserKey, config.DefaultDBMonitorMySQLUser)
		viper.Set(config.DBMonitorMySQLPassKey, config.DefaultDBMonitorMySQLPass)
	case 2:
		viper.Set(config.DBMonitorClickhouseUserKey, "")
		viper.Set(config.DBMonitorClickhousePassKey, "")
	}
}

func initQueryRepo() *DASRepo {
	var err error
	dbAddr := serviceTestDBAddr
	dbDBName := serviceTestDBDBName
	dbDBUser := serviceTestDBDBUser
	dbDBPass := serviceTestDBDBPass
	global.DASMySQLPool, err = mysql.NewPoolWithDefault(dbAddr, dbDBName, dbDBUser, dbDBPass)
	if err != nil {
		log.Error(common.CombineMessageWithError("initQueryRepo() failed", err))
		return nil
	}
	return newDASRepo(global.DASMySQLPool)
}

func initNewQuery() *Query {
	return &Query{
		defaultQuerySQLID,
		defaultQueryFingerprint,
		defaultQueryExample,
		defaultQueryDBName,
		defaultQueryExecCount,
		defaultQueryTotalExecTime,
		defaultQueryAvgExecTime,
		defaultQueryRowsExaminedMax,
	}
}

var service = createService()

func createService() *Service {
	return newService(NewConfigWithDefault(), initQueryRepo())
}

func TestService_All(t *testing.T) {
	TestService_GetConfig(t)
	TestService_GetQueries(t)
	TestService_GetByMySQLClusterID(t)
	TestService_GetByMySQLServerID(t)
	TestService_GetByDBID(t)
	TestService_GetBySQLID(t)
	TestService_Marshal(t)
	TestService_MarshalWithFields(t)
}

func TestService_GetConfig(t *testing.T) {
	asst := assert.New(t)

	limit := service.GetConfig().GetLimit()
	asst.Equal(defaultLimit, limit, "test GetConfig() failed")
}

func TestService_GetQueries(t *testing.T) {
	asst := assert.New(t)

	service.queries = append(service.queries, initNewQuery())
	sqlID := service.GetQueries()[0].GetSQLID()
	asst.Equal(defaultQuerySQLID, sqlID, "test GetQueries() failed")
}

func TestService_GetByMySQLClusterID(t *testing.T) {
	asst := assert.New(t)

	switch pmmVersion {
	case 1:
		err := service.GetByMySQLClusterID(2)
		asst.Nil(err, common.CombineMessageWithError("test GetByMySQLClusterID() failed", err))
	case 2:
		err := service.GetByMySQLClusterID(1)
		asst.Nil(err, common.CombineMessageWithError("test GetByMySQLClusterID() failed", err))
	}
}

func TestService_GetByMySQLServerID(t *testing.T) {
	asst := assert.New(t)

	switch pmmVersion {
	case 1:
		err := service.GetByMySQLServerID(3)
		asst.Nil(err, common.CombineMessageWithError("test GetByMySQLServerID() failed", err))
	case 2:
		err := service.GetByMySQLServerID(1)
		asst.Nil(err, common.CombineMessageWithError("test GetByMySQLServerID() failed", err))
	}
}

func TestService_GetByDBID(t *testing.T) {
	asst := assert.New(t)

	switch pmmVersion {
	case 1:
		err := service.GetByDBID(3, 1)
		asst.Nil(err, common.CombineMessageWithError("test GetByDBID() failed", err))
	case 2:
		err := service.GetByDBID(2, 2)
		asst.Nil(err, common.CombineMessageWithError("test GetByDBID() failed", err))
	}
}

func TestService_GetBySQLID(t *testing.T) {
	asst := assert.New(t)

	switch pmmVersion {
	case 1:
		err := service.GetBySQLID(6, "999ECD050D719733")
		asst.Nil(err, common.CombineMessageWithError("test GetBySQLID() failed", err))
	case 2:
		err := service.GetBySQLID(2, "999ECD050D719733")
		asst.Nil(err, common.CombineMessageWithError("test GetBySQLID() failed", err))
	}
}

func TestService_Marshal(t *testing.T) {
	asst := assert.New(t)

	_, err := service.Marshal()
	asst.Nil(err, common.CombineMessageWithError("test Marshal() failed", err))
}

func TestService_MarshalWithFields(t *testing.T) {
	asst := assert.New(t)

	_, err := service.MarshalWithFields(queryQueriesStruct)
	asst.Nil(err, common.CombineMessageWithError("test MarshalWithFields(fields ...string) failed", err))
}
