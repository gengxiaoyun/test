package query

import (
	"fmt"

	"github.com/romberli/das/config"
	"github.com/romberli/das/internal/app/metadata"
	depmeta "github.com/romberli/das/internal/dependency/metadata"
	"github.com/romberli/das/internal/dependency/query"
	"github.com/romberli/das/pkg/message"
	msgquery "github.com/romberli/das/pkg/message/query"
	"github.com/romberli/go-util/constant"
	"github.com/romberli/go-util/middleware/clickhouse"
	"github.com/romberli/go-util/middleware/mysql"
	"github.com/romberli/log"
	"github.com/spf13/viper"
)

const (
	pmmMySQLDBName      = "pmm"
	pmmClickhouseDBName = "pmm"
)

var _ query.Query = (*Query)(nil)

// Query include several members of a query
type Query struct {
	SQLID           string  `middleware:"sql_id" json:"sql_id"`
	Fingerprint     string  `middleware:"fingerprint" json:"fingerprint"`
	Example         string  `middleware:"example" json:"example"`
	DBName          string  `middleware:"db_name" json:"db_name"`
	ExecCount       int     `middleware:"exec_count" json:"exec_count"`
	TotalExecTime   float64 `middleware:"total_exec_time" json:"total_exec_time"`
	AvgExecTime     float64 `middleware:"avg_exec_time" json:"avg_exec_time"`
	RowsExaminedMax int     `middleware:"rows_examined_max" json:"rows_examined_max"`
}

// NewEmptyQuery return *Query
func NewEmptyQuery() *Query {
	return &Query{}
}

// GetSQLID returns the sql identity
func (q *Query) GetSQLID() string {
	return q.SQLID
}

// GetFingerprint returns the fingerprint
func (q *Query) GetFingerprint() string {
	return q.Fingerprint
}

// GetExample returns the example
func (q *Query) GetExample() string {
	return q.Example
}

// GetDBName returns the db name
func (q *Query) GetDBName() string {
	return q.DBName
}

// GetExecCount returns the execution count
func (q *Query) GetExecCount() int {
	return q.ExecCount
}

// GetTotalExecTime returns the total execution time
func (q *Query) GetTotalExecTime() float64 {
	return q.TotalExecTime
}

// GetAvgExecTime returns the average execution timey
func (q *Query) GetAvgExecTime() float64 {
	return q.AvgExecTime
}

// GetRowsExaminedMax returns the maximum row examined
func (q *Query) GetRowsExaminedMax() int {
	return q.RowsExaminedMax
}

// Querier include config of query and connection pool of DAS repo
type Querier struct {
	config  *Config
	dasRepo *DASRepo
}

// NewQuerier return *Querier
func NewQuerier(config *Config, dasRepo *DASRepo) *Querier {
	return newQuerier(config, dasRepo)
}

// NewQuerierWithGlobal return *Querier with global DASRepo
func NewQuerierWithGlobal(config *Config) *Querier {
	return newQuerier(config, NewDASRepoWithGlobal())
}

func newQuerier(config *Config, dasRepo *DASRepo) *Querier {
	return &Querier{
		config:  config,
		dasRepo: dasRepo,
	}
}

func (q *Querier) getConfig() *Config {
	return q.config
}

// GetByMySQLClusterID get queries by mysql cluster id
func (q *Querier) GetByMySQLClusterID(mysqlClusterID int) ([]query.Query, error) {
	mysqlServers, err := q.getMySQLServersByClusterID(mysqlClusterID)
	if err != nil {
		return nil, err
	}

	var queries []query.Query
	for _, mysqlServer := range mysqlServers {
		mysqlServerID := mysqlServer.Identity()
		// dispatch to GetByMySQLServerID()
		qs, err := q.GetByMySQLServerID(mysqlServerID)
		if err != nil {
			return nil, err
		}
		queries = append(queries, qs...)
	}
	return queries, nil
}

// GetByMySQLServerID get queries by mysql server id
func (q *Querier) GetByMySQLServerID(mysqlServerID int) ([]query.Query, error) {
	// init monitor repos
	monitorSystem, err := q.getMonitorSystemByMySQLServerID(mysqlServerID)
	if err != nil {
		return nil, err
	}
	monitorRepo, err := q.getMonitorRepo(monitorSystem)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = monitorRepo.Close()
		if err != nil {
			log.Error(message.NewMessage(msgquery.ErrQueryCloseMonitorRepo, err.Error()).Error())
		}
	}()
	// get mysql server
	mysqlServer, err := q.getMySQLServerByID(mysqlServerID)
	if err != nil {
		return nil, err
	}

	return monitorRepo.GetByServiceNames([]string{mysqlServer.GetServiceName()})
}

// GetByDBID get queries by db id
func (q *Querier) GetByDBID(mysqlServerID int, dbID int) ([]query.Query, error) {
	// init monitor repos
	monitorSystem, err := q.getMonitorSystemByMySQLServerID(mysqlServerID)
	if err != nil {
		return nil, err
	}
	monitorRepo, err := q.getMonitorRepo(monitorSystem)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = monitorRepo.Close()
		if err != nil {
			log.Error(message.NewMessage(msgquery.ErrQueryCloseMonitorRepo, err.Error()).Error())
		}
	}()
	// get mysql server
	mysqlServer, err := q.getMySQLServerByID(mysqlServerID)
	if err != nil {
		return nil, err
	}
	// get db
	db, err := q.getDBByID(dbID)
	if err != nil {
		return nil, err
	}

	return monitorRepo.GetByDBName(mysqlServer.GetServiceName(), db.GetDBName())
}

// GetBySQLID get queries by sql id
func (q *Querier) GetBySQLID(mysqlServerID int, sqlID string) ([]query.Query, error) {
	// init monitor repos
	monitorSystem, err := q.getMonitorSystemByMySQLServerID(mysqlServerID)
	if err != nil {
		return nil, err
	}
	monitorRepo, err := q.getMonitorRepo(monitorSystem)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = monitorRepo.Close()
		if err != nil {
			log.Error(message.NewMessage(msgquery.ErrQueryCloseMonitorRepo, err.Error()).Error())
		}
	}()
	// get mysql server
	mysqlServer, err := q.getMySQLServerByID(mysqlServerID)
	if err != nil {
		return nil, err
	}
	log.Infof("%s", mysqlServer.GetServiceName())
	queryResult, err := monitorRepo.GetBySQLID(mysqlServer.GetServiceName(), sqlID)

	return []query.Query{queryResult}, err
}

func (q *Querier) getMySQLServersByClusterID(mysqlClusterID int) ([]depmeta.MySQLServer, error) {
	service := metadata.NewMySQLServerServiceWithDefault()
	err := service.GetByClusterID(mysqlClusterID)
	if err != nil {
		return nil, err
	}

	return service.GetMySQLServers(), nil
}

func (q *Querier) getMySQLServerByID(mysqlServerID int) (depmeta.MySQLServer, error) {
	service := metadata.NewMySQLServerServiceWithDefault()
	err := service.GetByID(mysqlServerID)
	if err != nil {
		return nil, err
	}

	return service.GetMySQLServers()[constant.ZeroInt], nil
}

func (q *Querier) getDBByID(dbID int) (depmeta.DB, error) {
	service := metadata.NewDBServiceWithDefault()
	err := service.GetByID(dbID)
	if err != nil {
		return nil, err
	}

	return service.GetDBs()[constant.ZeroInt], nil
}

func (q *Querier) getMonitorSystemByMySQLClusterID(mysqlClusterID int) (depmeta.MonitorSystem, error) {
	mysqlClusterService := metadata.NewMySQLClusterServiceWithDefault()
	err := mysqlClusterService.GetByID(mysqlClusterID)
	if err != nil {
		return nil, err
	}
	mysqlCluster := mysqlClusterService.GetMySQLClusters()[constant.ZeroInt]
	monitorSystemID := mysqlCluster.GetMonitorSystemID()

	monitorSystemService := metadata.NewMonitorSystemServiceWithDefault()
	err = monitorSystemService.GetByID(monitorSystemID)
	if err != nil {
		return nil, err
	}

	return monitorSystemService.GetMonitorSystems()[constant.ZeroInt], nil
}

func (q *Querier) getMonitorSystemByMySQLServerID(mysqlServerID int) (depmeta.MonitorSystem, error) {
	mysqlServerService := metadata.NewMySQLServerServiceWithDefault()
	err := mysqlServerService.GetByID(mysqlServerID)
	if err != nil {
		return nil, err
	}
	mysqlServer := mysqlServerService.GetMySQLServers()[constant.ZeroInt]
	mysqlClusterID := mysqlServer.GetClusterID()

	return q.getMonitorSystemByMySQLClusterID(mysqlClusterID)
}

func (q *Querier) getMonitorRepo(monitorSystem depmeta.MonitorSystem) (query.MonitorRepo, error) {
	var monitorRepo query.MonitorRepo

	addr := fmt.Sprintf("%s:%d", monitorSystem.GetHostIP(), monitorSystem.GetPortNumSlow())
	switch monitorSystem.GetSystemType() {
	case 1:
		// pmm 1.x
		mysqlConn, err := mysql.NewConn(addr, pmmMySQLDBName, q.getMonitorMySQLUser(), q.getMonitorMySQLPass())
		if err != nil {
			return nil, err
		}
		monitorRepo = NewMySQLRepo(q.getConfig(), mysqlConn)
	case 2:
		// pmm 2.x
		clickhouseConn, err := clickhouse.NewConnWithDefault(addr, pmmClickhouseDBName, q.getMonitorClickhouseUser(), q.getMonitorClickhousePass())
		if err != nil {
			return nil, err
		}
		monitorRepo = NewClickHouseRepo(q.getConfig(), clickhouseConn)
	default:
		return nil, message.NewMessage(msgquery.ErrQueryMonitorSystemType, monitorSystem.GetSystemType())
	}

	return monitorRepo, nil
}

// getMonitorMySQLUser returns mysql username of monitor system
func (q *Querier) getMonitorMySQLUser() string {
	return viper.GetString(config.DBMonitorMySQLUserKey)
}

// getMonitorMySQLPass returns mysql password of monitor system
func (q *Querier) getMonitorMySQLPass() string {
	return viper.GetString(config.DBMonitorMySQLPassKey)
}

// getMonitorClickhouseUser returns clickhouse username of monitor system
func (q *Querier) getMonitorClickhouseUser() string {
	return viper.GetString(config.DBMonitorClickhouseUserKey)
}

// getMonitorClickhousePass returns clickhouse password of monitor system
func (q *Querier) getMonitorClickhousePass() string {
	return viper.GetString(config.DBMonitorClickhousePassKey)
}
