package query

import (
	"testing"

	"github.com/romberli/das/internal/app/metadata"
	"github.com/romberli/go-util/constant"
	"github.com/stretchr/testify/assert"
)

const (
	// modify the connection information
	dbAddr   = "192.168.10.220:3306"
	dbDBName = "das"
	dbDBUser = "root"
	dbDBPass = "root"
)

const (
	testServiceName = "test"
	testDbName      = "test"
	testSQLID       = "0"
)

func TestQueryRepositoryAll(t *testing.T) {
	TestDASRepo_GetMonitorSystemByDBID(t)
	TestDASRepo_GetMonitorSystemByClusterID(t)
	TestDASRepo_GetMonitorSystemByMySQLServerID(t)
	TestMySQLRepo_GetByServiceNames(t)
	TestMySQLRepo_GetByDBName(t)
	TestMySQLRepo_GetBySQLID(t)
	TestClickhouseRepo_GetByDBName(t)
	TestClickhouseRepo_GetByServiceNames(t)
	TestClickhouseRepo_GetBySQLID(t)
}

func TestDASRepo_GetMonitorSystemByDBID(t *testing.T) {
	asst := assert.New(t)

	dbInfo := metadata.NewDBServiceWithDefault()
	dbInfo.GetAll()
	dbs := dbInfo.DBs[constant.ZeroInt]
	dbID := dbs.Identity()

	var r *DASRepo
	ms, err := r.GetMonitorSystemByDBID(dbID)

	asst.Equal(nil, err, "test GetMonitorSystemByDBID Failed")
	asst.Equal(true, ms != nil, "test GetMonitorSystemByDBID Failed")
}

func TestDASRepo_GetMonitorSystemByClusterID(t *testing.T) {
	asst := assert.New(t)

	clusterInfo := metadata.NewMySQLClusterServiceWithDefault()
	clusterInfo.GetAll()
	mcs := clusterInfo.MySQLClusters[constant.ZeroInt]
	clusterID := mcs.Identity()

	var r *DASRepo
	ms, err := r.GetMonitorSystemByClusterID(clusterID)

	asst.Equal(nil, err, "test GetMonitorSystemByClusterID Failed")
	asst.Equal(true, ms != nil, "test GetMonitorSystemByClusterID Failed")
}

func TestDASRepo_GetMonitorSystemByMySQLServerID(t *testing.T) {
	asst := assert.New(t)

	serverInfo := metadata.NewMySQLServerServiceWithDefault()
	serverInfo.GetAll()
	ss := serverInfo.MySQLServers[constant.ZeroInt]
	serverID := ss.Identity()

	var r *DASRepo
	ms, err := r.GetMonitorSystemByMySQLServerID(serverID)

	asst.Equal(nil, err, "test GetMonitorSystemByMySQLServerID Failed")
	asst.Equal(true, ms != nil, "test GetMonitorSystemByMySQLServerID Failed")
}

func TestMySQLRepo_GetByServiceNames(t *testing.T) {
	asst := assert.New(t)

	serverInfo := metadata.NewMySQLServerServiceWithDefault()
	serverInfo.GetAll()
	ss := serverInfo.MySQLServers[constant.ZeroInt]
	serverName := ss.GetServerName()

	var mr *MySQLRepo
	qu, err := mr.GetByServiceNames([]string{serverName})
	asst.Equal(nil, err, "test MySQLRepo_GetByServiceNames Failed")
	asst.Equal(true, qu != nil, "test MySQLRepo_GetByServiceNames Failed")
}

func TestMySQLRepo_GetByDBName(t *testing.T) {
	asst := assert.New(t)

	var mr *MySQLRepo
	qu, err := mr.GetByDBName(testServiceName, testDbName)
	asst.Equal(nil, err, "test MySQLRepo_GetByDBName Failed")
	asst.Equal(true, qu != nil, "test MySQLRepo_GetByDBName Failed")
}

func TestMySQLRepo_GetBySQLID(t *testing.T) {
	asst := assert.New(t)

	var mr *MySQLRepo
	qu, err := mr.GetBySQLID(testServiceName, testSQLID)
	asst.Equal(nil, err, "test MySQLRepo_GetBySQLID Failed")
	asst.Equal(true, qu != nil, "test MySQLRepo_GetBySQLID Failed")
}

func TestClickhouseRepo_GetByDBName(t *testing.T) {
	asst := assert.New(t)

	var cr *ClickhouseRepo
	qu, err := cr.GetByDBName(testServiceName, testDbName)
	asst.Equal(nil, err, "test ClickhouseRepo_GetByDBName Failed")
	asst.Equal(true, qu != nil, "test ClickhouseRepo_GetByDBName Failed")

}

func TestClickhouseRepo_GetByServiceNames(t *testing.T) {
	asst := assert.New(t)

	var cr *ClickhouseRepo
	qu, err := cr.GetByServiceNames([]string{testServiceName})
	asst.Equal(nil, err, "test ClickhouseRepo_GetByServiceNamesFailed")
	asst.Equal(true, qu != nil, "test ClickhouseRepo_GetByServiceNames Failed")
}

func TestClickhouseRepo_GetBySQLID(t *testing.T) {
	asst := assert.New(t)

	var cr *ClickhouseRepo
	qu, err := cr.GetBySQLID(testServiceName, testSQLID)
	asst.Equal(nil, err, "test ClickhouseRepo_GetBySQLID Failed")
	asst.Equal(true, qu != nil, "test ClickhouseRepo_GetBySQLID Failed")
}
