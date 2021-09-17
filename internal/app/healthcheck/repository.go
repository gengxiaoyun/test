package healthcheck

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-version"
	"github.com/romberli/das/global"
	"github.com/romberli/das/internal/app/query"
	"github.com/romberli/das/internal/dependency/healthcheck"
	depquery "github.com/romberli/das/internal/dependency/query"
	"github.com/romberli/das/pkg/message"
	msghc "github.com/romberli/das/pkg/message/healthcheck"
	"github.com/romberli/go-util/common"
	"github.com/romberli/go-util/constant"
	"github.com/romberli/go-util/middleware"
	"github.com/romberli/go-util/middleware/clickhouse"
	"github.com/romberli/go-util/middleware/mysql"
	"github.com/romberli/go-util/middleware/prometheus"
	"github.com/romberli/log"
)

const (
	mysql57           = "5.7"
	performanceSchema = "performance_schema"
	informationSchema = "information_schema"
	deviceLabel       = "device"
	mountPointLabel   = "mountpoint"

	dataDirVariable   = "datadir"
	binlogDirVariable = "log_bin_base"

	minTableRows      = 30000000
	minRowsExamined   = 1
	SlowQueryNumLimit = 100
)

var (
	_ healthcheck.DASRepo              = (*DASRepo)(nil)
	_ healthcheck.ApplicationMySQLRepo = (*ApplicationMySQLRepo)(nil)
	_ healthcheck.PrometheusRepo       = (*PrometheusRepo)(nil)
	_ healthcheck.QueryRepo            = (*MySQLQueryRepo)(nil)
	_ healthcheck.QueryRepo            = (*ClickhouseQueryRepo)(nil)
)

// DASRepo for health check
type DASRepo struct {
	Database middleware.Pool
}

// NewDASRepo returns *DASRepo with given middleware.Pool
func NewDASRepo(db middleware.Pool) *DASRepo {
	return newDASRepo(db)
}

// NewDASRepoWithGlobal returns *DASRepo with global mysql pool
func NewDASRepoWithGlobal() *DASRepo {
	return NewDASRepo(global.DASMySQLPool)
}

// newDASRepo returns *DASRepo with given middleware.Pool
func newDASRepo(db middleware.Pool) *DASRepo {
	return &DASRepo{Database: db}
}

// Execute executes given command and placeholders on the middleware
func (dr *DASRepo) Execute(command string, args ...interface{}) (middleware.Result, error) {
	conn, err := dr.Database.Get()
	if err != nil {
		return nil, err
	}
	defer func() {
		err = conn.Close()
		if err != nil {
			log.Errorf("healthcheck DASRepo.Execute(): close database connection failed.\n%s", err.Error())
		}
	}()

	return conn.Execute(command, args...)
}

// Transaction returns a middleware.Transaction that could execute multiple commands as a transaction
func (dr *DASRepo) Transaction() (middleware.Transaction, error) {
	return dr.Database.Transaction()
}

// LoadEngineConfig loads engine config from the middleware
func (dr *DASRepo) LoadEngineConfig() (healthcheck.EngineConfig, error) {
	// load config
	sql := `
		select id, item_name, item_weight, low_watermark, high_watermark, unit, score_deduction_per_unit_high, max_score_deduction_high,
		score_deduction_per_unit_medium, max_score_deduction_medium, del_flag, create_time, last_update_time
		from t_hc_default_engine_config
		where del_flag = 0;
	`
	log.Debugf("healthcheck DASRepo.loadEngineConfig() sql: \n%s\n", sql)
	result, err := dr.Execute(sql)
	if err != nil {
		return nil, err
	}
	// init []*DefaultItemConfig
	defaultItemConfigs := make([]*DefaultItemConfig, result.RowNumber())
	for i := range defaultItemConfigs {
		defaultItemConfigs[i] = NewEmptyDefaultItemConfig()
	}
	// map to struct
	err = result.MapToStructSlice(defaultItemConfigs, constant.DefaultMiddlewareTag)
	if err != nil {
		return nil, err
	}
	defaultEngineConfig := NewEmptyDefaultEngineConfig()
	for _, defaultItemConfig := range defaultItemConfigs {
		defaultEngineConfig[defaultItemConfig.ItemName] = defaultItemConfig
	}

	return defaultEngineConfig, nil
}

// GetResultByOperationID gets a Result by the operationID from the middleware
func (dr *DASRepo) GetResultByOperationID(operationID int) (healthcheck.Result, error) {
	sql := `
		select id, operation_id, weighted_average_score,
		db_config_score, db_config_data, db_config_advice,
		cpu_usage_score, cpu_usage_data, cpu_usage_high,
		io_util_score, io_util_data, io_util_high,
		disk_capacity_usage_score, disk_capacity_usage_data, disk_capacity_usage_high,
		connection_usage_score, connection_usage_data, connection_usage_high,
		average_active_session_percents_score, average_active_session_percents_data, average_active_session_percents_high,
		cache_miss_ratio_score, cache_miss_ratio_data, cache_miss_ratio_high,
		table_rows_score, table_rows_data, table_rows_high,
		table_size_score, table_size_data, table_size_high,
		slow_query_score, slow_query_data, slow_query_advice,
		accuracy_review, del_flag, create_time, last_update_time
		from t_hc_result
		where del_flag = 0
		and operation_id = ? 
		order by id;
	`
	log.Debugf("healthCheck DASRepo.GetResultByOperationID select sql: \n%s\nplaceholders: %s", sql, operationID)

	result, err := dr.Execute(sql, operationID)
	if err != nil {
		return nil, err
	}
	switch result.RowNumber() {
	case 0:
		return nil, fmt.Errorf("healthCheck DASRepo.GetResultByOperationID(): data does not exists, operation_id: %d", operationID)
	case 1:
		hcInfo := NewEmptyResultWithRepo(dr)
		// map to struct
		err = result.MapToStructByRowIndex(hcInfo, constant.ZeroInt, constant.DefaultMiddlewareTag)
		if err != nil {
			return nil, err
		}

		return hcInfo, nil
	default:
		return nil, fmt.Errorf("healthCheck DASRepo.GetResultByOperationID(): duplicate key exists, operation_id: %d", operationID)
	}
}

// IsRunning gets status by the mysqlServerID from the middleware
func (dr *DASRepo) IsRunning(mysqlServerID int) (bool, error) {
	sql := `select count(1) from t_hc_operation_info where del_flag = 0 and mysql_server_id = ? and status = 1;`
	log.Debugf("healthCheck DASRepo.IsRunning() select sql: \n%s\nplaceholders: %s", sql, mysqlServerID)

	result, err := dr.Execute(sql, mysqlServerID)
	if err != nil {
		return false, err
	}
	count, _ := result.GetInt(constant.ZeroInt, constant.ZeroInt)

	return count != 0, nil
}

// InitOperation creates a testOperationInfo in the middleware
func (dr *DASRepo) InitOperation(mysqlServerID int, startTime, endTime time.Time, step time.Duration) (int, error) {
	startTimeStr := startTime.Format(constant.TimeLayoutSecond)
	endTimeStr := endTime.Format(constant.TimeLayoutSecond)
	stepInt := int(step.Seconds())

	sql := `insert into t_hc_operation_info(mysql_server_id, start_time, end_time, step) values(?, ?, ?, ?);`
	log.Debugf("healthCheck DASRepo.InitOperation() insert sql: \n%s\nplaceholders: %s, %s, %s, %s",
		sql, mysqlServerID, startTimeStr, endTimeStr, stepInt)

	_, err := dr.Execute(sql, mysqlServerID, startTimeStr, endTimeStr, stepInt)
	if err != nil {
		return constant.ZeroInt, err
	}

	sql = `
		select id from t_hc_operation_info where del_flag = 0 and 
		mysql_server_id = ? and start_time = ? and end_time = ? and step = ?;
	`
	log.Debugf("healthCheck DASRepo.InitOperation() select sql: \n%s\nplaceholders: %s, %s, %s, %s",
		sql, mysqlServerID, startTimeStr, endTimeStr, stepInt)

	result, err := dr.Execute(sql, mysqlServerID, startTimeStr, endTimeStr, stepInt)
	if err != nil {
		return constant.ZeroInt, err
	}

	return result.GetInt(constant.ZeroInt, constant.ZeroInt)
}

// UpdateOperationStatus updates the status and message by the operationID in the middleware
func (dr *DASRepo) UpdateOperationStatus(operationID int, status int, message string) error {
	sql := `update t_hc_operation_info set status = ?, message = ? where id = ?;`
	log.Debugf("healthCheck DASRepo.UpdateOperationStatus() update sql: \n%s\nplaceholders: %s, %s, %s", sql, operationID, status, message)
	_, err := dr.Execute(sql, status, message, operationID)

	return err
}

// SaveResult saves the result in the middleware
func (dr *DASRepo) SaveResult(result healthcheck.Result) error {
	sql := `insert into t_hc_result(operation_id, weighted_average_score, db_config_score, db_config_data, 
		db_config_advice, cpu_usage_score, cpu_usage_data, cpu_usage_high, io_util_score,
		io_util_data, io_util_high, disk_capacity_usage_score, disk_capacity_usage_data, 
		disk_capacity_usage_high, connection_usage_score, connection_usage_data, 
		connection_usage_high, average_active_session_percents_score, average_active_session_percents_data,
		average_active_session_percents_high, cache_miss_ratio_score, cache_miss_ratio_data, 
		cache_miss_ratio_high, table_rows_score, table_rows_data, table_rows_high,
		table_size_score, table_size_data, table_size_high, slow_query_score,
		slow_query_data, slow_query_advice, accuracy_review)
		values(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);
	`
	log.Debugf("healthCheck DASRepo.SaveResult() insert sql: \n%s\nplaceholders: %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %ss, %s, %s, %s",
		sql, result.GetOperationID(), result.GetWeightedAverageScore(),
		result.GetDBConfigScore(), result.GetDBConfigData(), result.GetDBConfigAdvice(),
		result.GetCPUUsageScore(), result.GetCPUUsageData(), result.GetCPUUsageHigh(),
		result.GetIOUtilScore(), result.GetIOUtilData(), result.GetIOUtilHigh(),
		result.GetDiskCapacityUsageScore(), result.GetDiskCapacityUsageData(), result.GetDiskCapacityUsageHigh(),
		result.GetConnectionUsageScore(), result.GetConnectionUsageData(), result.GetConnectionUsageHigh(),
		result.GetAverageActiveSessionPercentsScore(), result.GetAverageActiveSessionPercentsData(), result.GetAverageActiveSessionPercentsHigh(),
		result.GetCacheMissRatioScore(), result.GetCacheMissRatioData(), result.GetCacheMissRatioHigh(),
		result.GetTableRowsScore(), result.GetTableRowsData(), result.GetTableRowsHigh(),
		result.GetTableSizeScore(), result.GetTableSizeData(), result.GetTableSizeHigh(),
		result.GetSlowQueryScore(), result.GetSlowQueryData(), result.GetSlowQueryAdvice(), result.GetAccuracyReview())

	// execute
	_, err := dr.Execute(sql, result.GetOperationID(), result.GetWeightedAverageScore(),
		result.GetDBConfigScore(), result.GetDBConfigData(), result.GetDBConfigAdvice(),
		result.GetCPUUsageScore(), result.GetCPUUsageData(), result.GetCPUUsageHigh(),
		result.GetIOUtilScore(), result.GetIOUtilData(), result.GetIOUtilHigh(),
		result.GetDiskCapacityUsageScore(), result.GetDiskCapacityUsageData(), result.GetDiskCapacityUsageHigh(),
		result.GetConnectionUsageScore(), result.GetConnectionUsageData(), result.GetConnectionUsageHigh(),
		result.GetAverageActiveSessionPercentsScore(), result.GetAverageActiveSessionPercentsData(), result.GetAverageActiveSessionPercentsHigh(),
		result.GetCacheMissRatioScore(), result.GetCacheMissRatioData(), result.GetCacheMissRatioHigh(),
		result.GetTableRowsScore(), result.GetTableRowsData(), result.GetTableRowsHigh(),
		result.GetTableSizeScore(), result.GetTableSizeData(), result.GetTableSizeHigh(),
		result.GetSlowQueryScore(), result.GetSlowQueryData(), result.GetSlowQueryAdvice(), result.GetAccuracyReview())

	return err
}

// UpdateAccuracyReviewByOperationID updates the accuracyReview by the operationID in the middleware
func (dr *DASRepo) UpdateAccuracyReviewByOperationID(operationID int, review int) error {
	sql := `update t_hc_result set accuracy_review = ? where operation_id = ?;`
	log.Debugf("healthCheck DASRepo.UpdateAccuracyReviewByOperationID() update sql: \n%s\nplaceholders: %s, %s", sql, operationID, review)

	_, err := dr.Execute(sql, review, operationID)
	return err
}

// loadEngineConfig loads engine config from the middleware
func (dr *DASRepo) loadEngineConfig() (DefaultEngineConfig, error) {
	// load config
	sql := `
		select id, item_name, item_weight, low_watermark, high_watermark, unit, score_deduction_per_unit_high, max_score_deduction_high,
		score_deduction_per_unit_medium, max_score_deduction_medium, del_flag, create_time, last_update_time
		from t_hc_default_engine_config
		where del_flag = 0;
	`
	log.Debugf("healthcheck DASRepo.loadEngineConfig() sql: \n%s\n", sql)
	result, err := dr.Execute(sql)
	if err != nil {
		return nil, err
	}
	// init []*DefaultItemConfig
	defaultItemConfigs := make([]*DefaultItemConfig, result.RowNumber())
	for i := range defaultItemConfigs {
		defaultItemConfigs[i] = NewEmptyDefaultItemConfig()
	}
	// map to struct
	err = result.MapToStructSlice(defaultItemConfigs, constant.DefaultMiddlewareTag)
	if err != nil {
		return nil, err
	}
	defaultEngineConfig := NewEmptyDefaultEngineConfig()
	for _, defaultItemConfig := range defaultItemConfigs {
		defaultEngineConfig[defaultItemConfig.ItemName] = defaultItemConfig
	}

	return defaultEngineConfig, nil
}

type ApplicationMySQLRepo struct {
	operationInfo *OperationInfo
	conn          *mysql.Conn
}

// NewApplicationMySQLRepo returns a new *ApplicationMySQLRepo
func NewApplicationMySQLRepo(operationInfo *OperationInfo, conn *mysql.Conn) *ApplicationMySQLRepo {
	return &ApplicationMySQLRepo{
		operationInfo: operationInfo,
		conn:          conn,
	}
}

// GetOperationInfo returns the operation information
func (amr *ApplicationMySQLRepo) GetOperationInfo() *OperationInfo {
	return amr.operationInfo
}

// getConnection returns the connection
func (amr *ApplicationMySQLRepo) getConnection() *mysql.Conn {
	return amr.conn
}

// Close closes the application mysql connection
func (amr *ApplicationMySQLRepo) Close() error {
	return amr.getConnection().Close()
}

// GetVariables gets db config with given items
func (amr *ApplicationMySQLRepo) GetVariables(items []string) ([]healthcheck.Variable, error) {
	// prepare args
	interfaces, err := common.ConvertInterfaceToSliceInterface(items)
	if err != nil {
		return nil, err
	}
	inClause, err := middleware.ConvertSliceToString(interfaces...)
	if err != nil {
		return nil, err
	}
	// check mysql version
	mysqlVersion, err := version.NewVersion(amr.GetOperationInfo().GetMySQLServer().GetVersion())
	if err != nil {
		return nil, err
	}
	defaultVersion, err := version.NewVersion(mysql57)
	if err != nil {
		return nil, err
	}
	// prepare sql
	sql := fmt.Sprintf(applicationMySQLVariables, performanceSchema, inClause)
	if mysqlVersion.LessThan(defaultVersion) {
		sql = fmt.Sprintf(applicationMySQLVariables, informationSchema, inClause)
	}
	// get result
	result, err := amr.getConnection().Execute(sql)
	if err != nil {
		return nil, err
	}
	variables := make([]healthcheck.Variable, result.RowNumber())
	for i := range variables {
		variables[i] = NewEmptyGlobalVariable()
	}
	err = result.MapToStructSlice(variables, constant.DefaultMiddlewareTag)
	if err != nil {
		return nil, err
	}

	return variables, nil
}

// GetMySQLDirs gets the mysql directories
func (amr *ApplicationMySQLRepo) GetMySQLDirs() ([]string, error) {
	config, err := amr.GetVariables([]string{dataDirVariable, binlogDirVariable})
	if err != nil {
		return nil, err
	}

	var dirs []string
	for _, variable := range config {
		dirs = append(dirs, variable.GetValue())
	}

	return dirs, nil
}

// GetLargeTables gets the large tables
func (amr *ApplicationMySQLRepo) GetLargeTables() ([]healthcheck.Table, error) {
	result, err := amr.getConnection().Execute(applicationMySQLTableSize, minTableRows)
	if err != nil {
		return nil, err
	}
	tables := make([]healthcheck.Table, result.RowNumber())
	for i := range tables {
		tables[i] = NewEmptyTable()
	}
	err = result.MapToStructSlice(tables, constant.DefaultMiddlewareTag)
	if err != nil {
		return nil, err
	}

	return tables, nil
}

type PrometheusRepo struct {
	operationInfo *OperationInfo
	conn          *prometheus.Conn
}

// NewPrometheusRepo returns a new *PrometheusRepo
func NewPrometheusRepo(operationInfo *OperationInfo, conn *prometheus.Conn) *PrometheusRepo {
	return &PrometheusRepo{
		operationInfo: operationInfo,
		conn:          conn,
	}
}

// GetOperationInfo returns the operation information
func (pr *PrometheusRepo) GetOperationInfo() *OperationInfo {
	return pr.operationInfo
}

// getConnection returns the connection
func (pr *PrometheusRepo) getConnection() *prometheus.Conn {
	return pr.conn
}

// GetFileSystems gets the file systems from the prometheus
func (pr *PrometheusRepo) GetFileSystems() ([]healthcheck.FileSystem, error) {
	var prometheusQuery string

	// prepare query
	switch pr.getPMMVersion() {
	case 1:
		// pmm 1.x
		prometheusQuery = PrometheusFileSystemV1
	case 2:
		// pmm 2.x
		prometheusQuery = PrometheusFileSystemV2
	default:
		return nil, message.NewMessage(msghc.ErrPmmVersionInvalid)
	}

	prometheusQuery = fmt.Sprintf(prometheusQuery, pr.getNodeName())
	log.Debugf("healthcheck PrometheusRepo.GetFileSystems() query: \n%s\n", prometheusQuery)
	// get data
	result, err := pr.getConnection().Execute(prometheusQuery)
	if err != nil {
		return nil, err
	}
	// parse result
	vector, err := result.Raw.GetVector()
	if err != nil {
		return nil, err
	}

	var fileSystems []healthcheck.FileSystem
	for _, sample := range vector {
		fileSystems = append(fileSystems, NewFileSystem(string(sample.Metric[mountPointLabel]), string(sample.Metric[deviceLabel])))
	}

	return fileSystems, nil
}

// GetCPUUsage gets the cpu usage
func (pr *PrometheusRepo) GetCPUUsage() ([]healthcheck.PrometheusData, error) {
	var prometheusQuery string

	// prepare query
	switch pr.getPMMVersion() {
	case 1:
		// pmm 1.x
		prometheusQuery = PrometheusCPUUsageV1
	case 2:
		// pmm 2.x
		prometheusQuery = PrometheusCPUUsageV2
	default:
		return nil, message.NewMessage(msghc.ErrPmmVersionInvalid)
	}

	nodeName := pr.getNodeName()
	prometheusQuery = fmt.Sprintf(prometheusQuery, nodeName, nodeName, nodeName, nodeName, nodeName, nodeName)
	log.Debugf("healthcheck PrometheusRepo.GetCPUUsage() query: \n%s\n", prometheusQuery)

	return pr.execute(prometheusQuery)
}

// GetIOUtil gets the io util
func (pr *PrometheusRepo) GetIOUtil() ([]healthcheck.PrometheusData, error) {
	var prometheusQuery string

	nodeName := pr.getNodeName()

	// prepare query
	switch pr.getPMMVersion() {
	case 1:
		// pmm 1.x
		prometheusQuery = fmt.Sprintf(PrometheusIOUtilV1, nodeName, nodeName)
	case 2:
		// pmm 2.x
		prometheusQuery = fmt.Sprintf(PrometheusIOUtilV2, nodeName, nodeName, nodeName, nodeName)
	default:
		return nil, message.NewMessage(msghc.ErrPmmVersionInvalid)
	}

	log.Debugf("healthcheck PrometheusRepo.GetIOUtil() query: \n%s\n", prometheusQuery)
	// get data
	return pr.execute(prometheusQuery)
}

// GetDiskCapacityUsage gets the disk capacity usage
func (pr *PrometheusRepo) GetDiskCapacityUsage(mountPoints []string) ([]healthcheck.PrometheusData, error) {
	var prometheusQuery string

	mps := common.ConvertStringSliceToString(mountPoints, constant.VerticalBarString)
	nodeName := pr.getNodeName()

	// prepare query
	switch pr.getPMMVersion() {
	case 1:
		// pmm 1.x
		prometheusQuery = fmt.Sprintf(PrometheusDiskCapacityV1, nodeName, mps, nodeName, mps)
	case 2:
		// pmm 2.x
		prometheusQuery = fmt.Sprintf(PrometheusDiskCapacityV2, nodeName, mps, nodeName, mps, nodeName, mps, nodeName, mps)
	default:
		return nil, message.NewMessage(msghc.ErrPmmVersionInvalid)
	}

	log.Debugf("healthcheck PrometheusRepo.GetDiskCapacityUsage() query: \n%s\n", prometheusQuery)
	// get data
	return pr.execute(prometheusQuery)
}

// GetConnectionUsage gets the connection usage
func (pr *PrometheusRepo) GetConnectionUsage() ([]healthcheck.PrometheusData, error) {
	var prometheusQuery string

	// prepare query
	switch pr.getPMMVersion() {
	case 1:
		// pmm 1.x
		prometheusQuery = PrometheusConnectionUsageV1
	case 2:
		// pmm 2.x
		prometheusQuery = PrometheusConnectionUsageV2
	default:
		return nil, message.NewMessage(msghc.ErrPmmVersionInvalid)
	}

	serviceName := pr.getServiceName()
	prometheusQuery = fmt.Sprintf(prometheusQuery, serviceName, serviceName, serviceName, serviceName)
	log.Debugf("healthcheck PrometheusRepo.GetConnectionUsage() query: \n%s\n", prometheusQuery)
	// get data
	return pr.execute(prometheusQuery)
}

// GetActiveSessionPercents gets the active session number
func (pr *PrometheusRepo) GetAverageActiveSessionPercents() ([]healthcheck.PrometheusData, error) {
	var prometheusQuery string

	// prepare query
	switch pr.getPMMVersion() {
	case 1:
		// pmm 1.x
		prometheusQuery = PrometheusAverageActiveSessionPercentsV1
	case 2:
		// pmm 2.x
		prometheusQuery = PrometheusAverageActiveSessionPercentsV2
	default:
		return nil, message.NewMessage(msghc.ErrPmmVersionInvalid)
	}

	serviceName := pr.getServiceName()
	prometheusQuery = fmt.Sprintf(prometheusQuery, serviceName, serviceName, serviceName, serviceName)
	log.Debugf("healthcheck PrometheusRepo.GetAverageActiveSessionPercents() query: \n%s\n", prometheusQuery)
	// get data
	return pr.execute(prometheusQuery)
}

// GetCacheMissRatio gets the cache miss ratio
func (pr *PrometheusRepo) GetCacheMissRatio() ([]healthcheck.PrometheusData, error) {
	var prometheusQuery string

	// prepare query
	switch pr.getPMMVersion() {
	case 1:
		// pmm 1.x
		prometheusQuery = PrometheusCacheMissRatioV1
	case 2:
		// pmm 2.x
		prometheusQuery = PrometheusCacheMissRatioV2
	default:
		return nil, message.NewMessage(msghc.ErrPmmVersionInvalid)
	}

	serviceName := pr.getServiceName()
	prometheusQuery = fmt.Sprintf(prometheusQuery, serviceName, serviceName, serviceName, serviceName)
	log.Debugf("healthcheck PrometheusRepo.getCacheMissRatio() query: \n%s\n", prometheusQuery)
	// get data
	return pr.execute(prometheusQuery)
}

// getServiceName returns the service name
func (pr *PrometheusRepo) getServiceName() string {
	return pr.GetOperationInfo().GetMySQLServer().GetServiceName()
}

// getNodeName returns the node name
func (pr *PrometheusRepo) getNodeName() string {
	serviceName := pr.getServiceName()
	strList := strings.Split(serviceName, constant.ColonString)
	if len(strList) > constant.ZeroInt {
		return strList[constant.ZeroInt]
	}

	return serviceName[:strings.LastIndex(serviceName, constant.DashString)]
}

// getPMMVersion returns the pmm version
func (pr *PrometheusRepo) getPMMVersion() int {
	return pr.GetOperationInfo().GetMonitorSystem().GetSystemType()
}

// execute executes the given query
func (pr *PrometheusRepo) execute(query string) ([]healthcheck.PrometheusData, error) {
	// execute query
	result, err := pr.getConnection().Execute(query, pr.GetOperationInfo().GetStartTime(),
		pr.GetOperationInfo().GetEndTime(), pr.GetOperationInfo().GetStep())
	if err != nil {
		return nil, err
	}
	// parse result
	var datas []healthcheck.PrometheusData

	matrix, err := result.Raw.GetMatrix()
	if err != nil {
		return nil, err
	}
	for _, sampleStream := range matrix {
		for _, samplePair := range sampleStream.Values {
			datas = append(datas, NewPrometheusData(samplePair.Timestamp.String(), float64(samplePair.Value)))
		}
	}

	return datas, nil
}

type MySQLQueryRepo struct {
	operationInfo *OperationInfo
	conn          *mysql.Conn
}

// NewMySQLQueryRepo returns the new *MySQLQueryRepo
func NewMySQLQueryRepo(operationInfo *OperationInfo, conn *mysql.Conn) *MySQLQueryRepo {
	return &MySQLQueryRepo{
		operationInfo: operationInfo,
		conn:          conn,
	}
}

// GetOperationInfo returns the operation information
func (mqr *MySQLQueryRepo) GetOperationInfo() *OperationInfo {
	return mqr.operationInfo
}

// getConnection returns the connection
func (mqr *MySQLQueryRepo) getConnection() *mysql.Conn {
	return mqr.conn
}

// Close closes the connection
func (mqr *MySQLQueryRepo) Close() error {
	return mqr.getConnection().Close()
}

// GetSlowQuery gets the slow query
func (mqr *MySQLQueryRepo) GetSlowQuery() ([]depquery.Query, error) {
	// get result
	result, err := mqr.getConnection().Execute(MonitorMySQLQuery, mqr.getServiceName(), mqr.GetOperationInfo().GetStartTime(),
		mqr.GetOperationInfo().GetEndTime(), minRowsExamined, SlowQueryNumLimit)
	if err != nil {
		return nil, err
	}
	// map result to slice
	queries := make([]depquery.Query, result.RowNumber())
	for i := constant.ZeroInt; i < result.RowNumber(); i++ {
		queries[i] = query.NewEmptyQuery()
	}
	err = result.MapToStructSlice(queries, constant.DefaultMiddlewareTag)
	if err != nil {
		return nil, err
	}

	return queries, nil
}

// getPMMVersion return the pmm version
func (mqr *MySQLQueryRepo) getServiceName() string {
	return mqr.GetOperationInfo().GetMySQLServer().GetServiceName()
}

// getPMMVersion return the pmm version
func (mqr *MySQLQueryRepo) getPMMVersion() int {
	return mqr.GetOperationInfo().GetMonitorSystem().GetSystemType()
}

type ClickhouseQueryRepo struct {
	operationInfo *OperationInfo
	conn          *clickhouse.Conn
}

// NewClickhouseQueryRepo returns the new *ClickhouseQueryRepo
func NewClickhouseQueryRepo(operationInfo *OperationInfo, conn *clickhouse.Conn) *ClickhouseQueryRepo {
	return &ClickhouseQueryRepo{
		operationInfo: operationInfo,
		conn:          conn,
	}
}

// GetOperationInfo returns the operation information
func (cqr *ClickhouseQueryRepo) GetOperationInfo() *OperationInfo {
	return cqr.operationInfo
}

// getConnection returns the connection
func (cqr *ClickhouseQueryRepo) getConnection() *clickhouse.Conn {
	return cqr.conn
}

// Close closes the connection
func (cqr *ClickhouseQueryRepo) Close() error {
	return cqr.getConnection().Close()
}

// GetSlowQuery gets the slow query
func (cqr *ClickhouseQueryRepo) GetSlowQuery() ([]depquery.Query, error) {
	// get result
	result, err := cqr.getConnection().Execute(MonitorClickhouseQuery, cqr.getServiceName(), cqr.GetOperationInfo().GetStartTime(),
		cqr.GetOperationInfo().GetEndTime(), minRowsExamined, SlowQueryNumLimit,
		cqr.getServiceName(), cqr.GetOperationInfo().GetStartTime(), cqr.GetOperationInfo().GetEndTime(), minRowsExamined)
	if err != nil {
		return nil, err
	}
	// map result to slice
	queries := make([]depquery.Query, result.RowNumber())
	for i := constant.ZeroInt; i < result.RowNumber(); i++ {
		queries[i] = query.NewEmptyQuery()
	}
	err = result.MapToStructSlice(queries, constant.DefaultMiddlewareTag)
	if err != nil {
		return nil, err
	}

	return queries, nil
}

// getServiceName returns the service name
func (cqr *ClickhouseQueryRepo) getServiceName() string {
	return cqr.GetOperationInfo().GetMySQLServer().GetServiceName()
}

// getPMMVersion returns the pmm version
func (cqr *ClickhouseQueryRepo) getPMMVersion() int {
	return cqr.GetOperationInfo().GetMonitorSystem().GetSystemType()
}
