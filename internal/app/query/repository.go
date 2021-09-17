package query

import (
	"fmt"
	"time"

	"github.com/romberli/das/global"
	"github.com/romberli/das/internal/app/metadata"
	demetadata "github.com/romberli/das/internal/dependency/metadata"
	"github.com/romberli/das/internal/dependency/query"
	"github.com/romberli/go-util/common"
	"github.com/romberli/go-util/constant"
	"github.com/romberli/go-util/middleware"
	"github.com/romberli/go-util/middleware/clickhouse"
	"github.com/romberli/go-util/middleware/mysql"
	"github.com/romberli/log"
)

const (
	mysqlQueryWithServiceNames = `
        select qc.checksum as sql_id,
               qc.fingerprint,
               qe.query    as example,
               qe.db       as db_name,
               m.exec_count,
               m.total_exec_time,
               m.avg_exec_time,
               m.rows_examined_max
        from (
                 select qcm.query_class_id,
                        sum(qcm.query_count)                                        as exec_count,
                        truncate(sum(qcm.query_time_sum), 2)                        as total_exec_time,
                        truncate(sum(qcm.query_time_sum) / sum(qcm.query_count), 2) as avg_exec_time,
                        qcm.rows_examined_max
                 from query_class_metrics qcm
                          inner join instances i on qcm.instance_id = i.instance_id
                 where i.name in (%s)
                   and qcm.start_ts >= ?
                   and qcm.start_ts < ?
				   and qcm.rows_examined_max >= ?
                 group by qcm.query_class_id
                 order by qcm.rows_examined_max desc
                  limit ? offset ?) m
                 inner join query_classes qc on m.query_class_id = qc.query_class_id
                 left join query_examples qe on m.query_class_id = qe.query_class_id;
    `
	mysqlQueryWithDBName = `
        select qc.checksum as sql_id,
               qc.fingerprint,
               qe.query    as example,
               qe.db       as db_name,
               m.exec_count,
               m.total_exec_time,
               m.avg_exec_time,
               m.rows_examined_max
        from (
                 select qcm.query_class_id,
                        sum(qcm.query_count)                                        as exec_count,
                        truncate(sum(qcm.query_time_sum), 2)                        as total_exec_time,
                        truncate(sum(qcm.query_time_sum) / sum(qcm.query_count), 2) as avg_exec_time,
                        qcm.rows_examined_max
                 from query_class_metrics qcm
                          inner join instances i on qcm.instance_id = i.instance_id
                 		  inner join query_examples qe on qcm.query_class_id = qe.query_class_id
                 where i.name in (%s)
				   and qe.db = ?
                   and qcm.start_ts >= ?
                   and qcm.start_ts < ?
				   and qcm.rows_examined_max >= ?
                 group by qcm.query_class_id
                 order by qcm.rows_examined_max desc
				  limit ? offset ?) m
                 inner join query_classes qc on m.query_class_id = qc.query_class_id
                 left join query_examples qe on m.query_class_id = qe.query_class_id;
    `
	mysqlQueryWithSQLID = `
        select qc.checksum as sql_id,
               qc.fingerprint,
               qe.query    as example,
               qe.db       as db_name,
               m.exec_count,
               m.total_exec_time,
               m.avg_exec_time,
               m.rows_examined_max
        from (
                 select qcm.query_class_id,
                        sum(qcm.query_count)                                        as exec_count,
                        truncate(sum(qcm.query_time_sum), 2)                        as total_exec_time,
                        truncate(sum(qcm.query_time_sum) / sum(qcm.query_count), 2) as avg_exec_time,
                        qcm.rows_examined_max
                 from query_class_metrics qcm
                          inner join instances i on qcm.instance_id = i.instance_id
						  inner join query_classes qc on qcm.query_class_id = qc.query_class_id
                 where i.name in (%s)
				   and qc.checksum = ?
                   and qcm.start_ts >= ?
                   and qcm.start_ts < ?
                 group by query_class_id) m
                 inner join query_classes qc on m.query_class_id = qc.query_class_id
                 left join query_examples qe on m.query_class_id = qe.query_class_id
        limit 1;
    `
	clickhouseQueryWithServiceNames = `
        select sm.sql_id,
               m.fingerprint,
               m.example,
               m.db_name,
               sm.exec_count,
               sm.total_exec_time,
               sm.avg_exec_time,
               sm.rows_examined_max
        
        from (
                 select queryid                                               as sql_id,
                        sum(num_queries)                                      as exec_count,
                        truncate(sum(m_query_time_sum), 2)                    as total_exec_time,
                        truncate(sum(m_query_time_sum) / sum(num_queries), 2) as avg_exec_time,
                        max(m_rows_examined_max)                              as rows_examined_max
                 from metrics
                 where service_type = 'mysql'
                   and service_name in (%s)
                   and period_start >= ?
                   and period_start < ?
                   and m_rows_examined_max >= ?
                 group by queryid
                 order by rows_examined_max desc
                 limit ? offset ? ) sm
                 left join (select queryid          as sql_id,
                                   max(fingerprint) as fingerprint,
                                   max(example)     as example,
                                   max(database)    as db_name
                            from metrics
                            where service_type = 'mysql'
                              and service_name in (%s)
                              and period_start >= ?
                              and period_start < ?
                              and m_rows_examined_max >= ?
                            group by queryid) m
                           on sm.sql_id = m.sql_id;
    `
	clickhouseQueryWithDBName = `
        select sm.sql_id,
               m.fingerprint,
               m.example,
               m.db_name,
               sm.exec_count,
               sm.total_exec_time,
               sm.avg_exec_time,
               sm.rows_examined_max
        
        from (
                 select queryid                                               as sql_id,
                        sum(num_queries)                                      as exec_count,
                        truncate(sum(m_query_time_sum), 2)                    as total_exec_time,
                        truncate(sum(m_query_time_sum) / sum(num_queries), 2) as avg_exec_time,
                        max(m_rows_examined_max)                              as rows_examined_max
                 from metrics
                 where service_type = 'mysql'
                   and service_name in (%s)
                   and (database = ? or schema = ?)
                   and period_start >= ?
                   and period_start < ?
                   and m_rows_examined_max >= ?
                 group by queryid
                 order by rows_examined_max desc
                 limit ? offset ? ) sm
                 left join (select queryid          as sql_id,
                                   max(fingerprint) as fingerprint,
                                   max(example)     as example,
                                   max(database)    as db_name
                            from metrics
                            where service_type = 'mysql'
                              and service_name in (%s)
                              and (database = ? or schema = ?)
                              and period_start >= ?
                              and period_start < ?
                              and m_rows_examined_max >= ?
                            group by queryid) m
                           on sm.sql_id = m.sql_id;
    `
	clickhouseQueryWithSQLID = `
        select sm.sql_id,
               m.fingerprint,
               m.example,
               m.db_name,
               sm.exec_count,
               sm.total_exec_time,
               sm.avg_exec_time,
               sm.rows_examined_max
        
        from (
                 select queryid                                               as sql_id,
                        sum(num_queries)                                      as exec_count,
                        truncate(sum(m_query_time_sum), 2)                    as total_exec_time,
                        truncate(sum(m_query_time_sum) / sum(num_queries), 2) as avg_exec_time,
                        max(m_rows_examined_max)                              as rows_examined_max
                 from metrics
                 where service_type = 'mysql'
                   and service_name in (%s)
                   and queryid = ?
                   and period_start >= ?
                   and period_start < ?
                   and m_rows_examined_max >= ?
                 group by queryid
                 order by rows_examined_max desc
                 limit ? offset ? ) sm
                 left join (select queryid          as sql_id,
                                   max(fingerprint) as fingerprint,
                                   max(example)     as example,
                                   max(database)    as db_name
                            from metrics
                            where service_type = 'mysql'
                              and service_name in (%s)
                              and queryid = ?
                              and period_start >= ?
                              and period_start < ?
                              and m_rows_examined_max >= ?
                            group by queryid) m
                           on sm.sql_id = m.sql_id;
    `
)

var _ query.DASRepo = (*DASRepo)(nil)
var _ query.MonitorRepo = (*MySQLRepo)(nil)
var _ query.MonitorRepo = (*ClickhouseRepo)(nil)

type DASRepo struct {
	Database middleware.Pool
}

// NewDASRepo returns *DASRepo
func NewDASRepo(db middleware.Pool) *DASRepo {
	return newDASRepo(db)
}

// NewDASRepoWithGlobal returns *DASRepo with global mysql pool
func NewDASRepoWithGlobal() *DASRepo {
	return NewDASRepo(global.DASMySQLPool)
}

// NewDASRepo returns *DASRepo
func newDASRepo(db middleware.Pool) *DASRepo {
	return &DASRepo{Database: db}
}

// Execute executes given command and placeholders on the middleware
func (r *DASRepo) Execute(command string, args ...interface{}) (middleware.Result, error) {
	conn, err := r.Database.Get()
	if err != nil {
		return nil, err
	}
	defer func() {
		err = conn.Close()
		if err != nil {
			log.Errorf("query DASRepo.Execute(): close database connection failed.\n%s", err.Error())
		}
	}()

	return conn.Execute(command, args...)
}

// Transaction returns a middleware.Transaction that could execute multiple commands as a transaction
func (r *DASRepo) Transaction() (middleware.Transaction, error) {
	return r.Database.Transaction()
}

// GetMonitorSystemByDBID returns a metadata.MonitorSystem by Database ID
func (r *DASRepo) GetMonitorSystemByDBID(dbID int) (demetadata.MonitorSystem, error) {
	dbInfo := metadata.NewDBServiceWithDefault()
	err := dbInfo.GetByID(dbID)
	if err != nil {
		return nil, err
	}
	dbs := dbInfo.DBs[constant.ZeroInt]
	clusterID := dbs.GetClusterID()
	return r.GetMonitorSystemByClusterID(clusterID)
}

// GetMonitorSystemByMySQLServerID returns a metadata.MonitorSystem by mysqlServerID
func (r *DASRepo) GetMonitorSystemByMySQLServerID(mysqlServerID int) (demetadata.MonitorSystem, error) {
	serverInfo := metadata.NewMySQLServerServiceWithDefault()
	err := serverInfo.GetByID(mysqlServerID)
	if err != nil {
		return nil, err
	}
	ss := serverInfo.MySQLServers[constant.ZeroInt]
	clusterID := ss.GetClusterID()
	return r.GetMonitorSystemByClusterID(clusterID)
}

// GetMonitorSystemByClusterID returns a metadata.MonitorSystem by clusterID
func (r *DASRepo) GetMonitorSystemByClusterID(clusterID int) (demetadata.MonitorSystem, error) {
	clusterInfo := metadata.NewMySQLClusterServiceWithDefault()
	err := clusterInfo.GetByID(clusterID)
	if err != nil {
		return nil, err
	}
	mcs := clusterInfo.MySQLClusters[constant.ZeroInt]
	monitorSystemID := mcs.GetMonitorSystemID()

	monitorSystemInfo := metadata.NewMonitorSystemServiceWithDefault()
	err = clusterInfo.GetByID(monitorSystemID)
	if err != nil {
		return nil, err
	}
	msi := monitorSystemInfo.MonitorSystems[constant.ZeroInt]
	return msi, nil
}

// Save save dasInfo into table
func (r *DASRepo) Save(mysqlClusterID, mysqlServerID, dbID int, sqlID string, startTime, endTime time.Time, limit, offset int) error {
	sql := "\t\tinsert into t_query_operation_info(mysql_cluster_id, mysql_server_id, db_id, sql_id, start_time, end_time, `limit`, offset) values(?, ?, ?, ?, ?, ?, ?, ?);"
	_, err := r.Execute(sql, mysqlClusterID, mysqlServerID, dbID, sqlID, startTime.Format(constant.DefaultTimeLayout), endTime.Format(constant.DefaultTimeLayout), limit, offset)

	return err
}

type MySQLRepo struct {
	config *Config
	conn   *mysql.Conn
}

// NewMySQLRepo returns a new mysqlRepo
func NewMySQLRepo(config *Config, conn *mysql.Conn) *MySQLRepo {
	return &MySQLRepo{
		config: config,
		conn:   conn,
	}
}

// getConfig gets Config
func (mr *MySQLRepo) getConfig() *Config {
	return mr.config
}

// Close closes the connection
func (mr *MySQLRepo) Close() error {
	return mr.conn.Close()
}

// GetByServiceNames return query.query list by serviceName
func (mr *MySQLRepo) GetByServiceNames(serviceName []string) ([]query.Query, error) {
	interfaces, err := common.ConvertInterfaceToSliceInterface(serviceName)
	if err != nil {
		return nil, err
	}

	services, err := middleware.ConvertSliceToString(interfaces...)
	if err != nil {
		return nil, err
	}

	sql := fmt.Sprintf(mysqlQueryWithServiceNames, services)

	return mr.execute(sql,
		mr.getConfig().GetStartTime().Format(constant.DefaultTimeLayout),
		mr.getConfig().GetEndTime().Format(constant.DefaultTimeLayout),
		minRowsExamined,
		mr.getConfig().GetLimit(),
		mr.getConfig().GetOffset(),
	)
}

// GetByDBName returns query.query list by dbName
func (mr *MySQLRepo) GetByDBName(serviceName, dbName string) ([]query.Query, error) {
	interfaces, err := common.ConvertInterfaceToSliceInterface([]string{serviceName})
	if err != nil {
		return nil, err
	}

	services, err := middleware.ConvertSliceToString(interfaces...)
	if err != nil {
		return nil, err
	}

	sql := fmt.Sprintf(mysqlQueryWithDBName, services)

	return mr.execute(sql,
		dbName,
		mr.getConfig().GetStartTime().Format(constant.DefaultTimeLayout),
		mr.getConfig().GetEndTime().Format(constant.DefaultTimeLayout),
		minRowsExamined,
		mr.getConfig().GetLimit(),
		mr.getConfig().GetOffset())
}

// GetBySQLID return query.query by SQL ID
func (mr *MySQLRepo) GetBySQLID(serviceName, sqlID string) (query.Query, error) {
	interfaces, err := common.ConvertInterfaceToSliceInterface([]string{serviceName})
	if err != nil {
		return nil, err
	}

	services, err := middleware.ConvertSliceToString(interfaces...)
	if err != nil {
		return nil, err
	}

	sql := fmt.Sprintf(mysqlQueryWithSQLID, services)

	queries, err := mr.execute(sql,
		sqlID,
		mr.getConfig().GetStartTime().Format(constant.DefaultTimeLayout),
		mr.getConfig().GetEndTime().Format(constant.DefaultTimeLayout),
	)
	if len(queries) == 0 {
		return nil, fmt.Errorf("sql(id=%s) in service(name=%s) is not found", sqlID, serviceName)
	}
	return queries[constant.ZeroInt], err
}

// execute executes the SQL with args
func (mr *MySQLRepo) execute(command string, args ...interface{}) ([]query.Query, error) {
	log.Debugf("query MySQLRepo.execute() sql: %s, args: %v", command, args)

	// get slow queries from the monitor database
	result, err := mr.conn.Execute(command, args...)
	if err != nil {
		return nil, err
	}
	// init queries
	queries := make([]query.Query, result.RowNumber())
	for i := constant.ZeroInt; i < result.RowNumber(); i++ {
		queries[i] = NewEmptyQuery()
	}
	// map result to queries
	err = result.MapToStructSlice(queries, constant.DefaultMiddlewareTag)
	if err != nil {
		return nil, err
	}

	return queries, nil
}

type ClickhouseRepo struct {
	config *Config
	conn   *clickhouse.Conn
}

// NewClickHouseRepo returns a new ClickHouseRepo
func NewClickHouseRepo(config *Config, conn *clickhouse.Conn) *ClickhouseRepo {
	return &ClickhouseRepo{
		config: config,
		conn:   conn,
	}
}

// getConfig returns the configuration
func (cr *ClickhouseRepo) getConfig() *Config {
	return cr.config
}

// Close closes the connection
func (cr *ClickhouseRepo) Close() error {
	return cr.conn.Close()
}

// GetByServiceNames returns query.Query list by serviceNames
func (cr *ClickhouseRepo) GetByServiceNames(serviceNames []string) ([]query.Query, error) {
	interfaces, err := common.ConvertInterfaceToSliceInterface(serviceNames)
	if err != nil {
		return nil, err
	}

	services, err := middleware.ConvertSliceToString(interfaces...)
	if err != nil {
		return nil, err
	}

	sql := fmt.Sprintf(clickhouseQueryWithServiceNames, services, services)

	return cr.execute(
		sql,
		cr.getConfig().GetStartTime(),
		cr.getConfig().GetEndTime(),
		minRowsExamined,
		cr.getConfig().GetLimit(),
		cr.getConfig().GetOffset(),
		cr.getConfig().GetStartTime(),
		cr.getConfig().GetEndTime(),
		minRowsExamined,
	)
}

// GetByDBName returns query.Query list by dbNameS
func (cr *ClickhouseRepo) GetByDBName(serviceName, dbName string) ([]query.Query, error) {
	interfaces, err := common.ConvertInterfaceToSliceInterface([]string{serviceName})
	if err != nil {
		return nil, err
	}

	services, err := middleware.ConvertSliceToString(interfaces...)
	if err != nil {
		return nil, err
	}

	sql := fmt.Sprintf(clickhouseQueryWithDBName, services, services)

	return cr.execute(sql,
		dbName,
		dbName,
		cr.getConfig().GetStartTime(),
		cr.getConfig().GetEndTime(),
		minRowsExamined,
		cr.getConfig().GetLimit(),
		cr.getConfig().GetOffset(),
		dbName,
		dbName,
		cr.getConfig().GetStartTime(),
		cr.getConfig().GetEndTime(),
		minRowsExamined,
	)
}

// GetBySQLID returns query.Query by SQL ID
func (cr *ClickhouseRepo) GetBySQLID(serviceName, sqlID string) (query.Query, error) {
	interfaces, err := common.ConvertInterfaceToSliceInterface([]string{serviceName})
	if err != nil {
		return nil, err
	}

	services, err := middleware.ConvertSliceToString(interfaces...)
	if err != nil {
		return nil, err
	}

	sql := fmt.Sprintf(clickhouseQueryWithSQLID, services, services)

	queries, err := cr.execute(sql,
		sqlID,
		cr.getConfig().GetStartTime(),
		cr.getConfig().GetEndTime(),
		minRowsExamined,
		cr.getConfig().GetLimit(),
		cr.getConfig().GetOffset(),
		sqlID,
		cr.getConfig().GetStartTime(),
		cr.getConfig().GetEndTime(),
		minRowsExamined,
	)
	if len(queries) == 0 {
		return nil, fmt.Errorf("sql(id=%s) in service(name=%s) is not found", sqlID, serviceName)
	}

	return queries[constant.ZeroInt], err
}

func (cr *ClickhouseRepo) execute(command string, args ...interface{}) ([]query.Query, error) {
	log.Debugf("query ClickhouseRepo.execute() sql: %s, args: %v", command, args)

	// get slow queries from the monitor database
	result, err := cr.conn.Execute(command, args...)
	if err != nil {
		return nil, err
	}
	// init queries
	queries := make([]query.Query, result.RowNumber())
	for i := constant.ZeroInt; i < result.RowNumber(); i++ {
		queries[i] = NewEmptyQuery()
	}
	// map result to queries
	err = result.MapToStructSlice(queries, constant.DefaultMiddlewareTag)
	if err != nil {
		return nil, err
	}

	return queries, nil
}
