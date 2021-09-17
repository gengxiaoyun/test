package healthcheck

import (
	"time"

	"github.com/romberli/das/internal/dependency/healthcheck"

	"github.com/romberli/go-util/common"
	"github.com/romberli/go-util/constant"
)

// Result include all data needed in healthcheck
type Result struct {
	healthcheck.DASRepo
	ID                                int       `middleware:"id" json:"id"`
	OperationID                       int       `middleware:"operation_id" json:"operation_id"`
	WeightedAverageScore              int       `middleware:"weighted_average_score" json:"weighted_average_score"`
	DBConfigScore                     int       `middleware:"db_config_score" json:"db_config_score"`
	DBConfigData                      string    `middleware:"db_config_data" json:"db_config_data"`
	DBConfigAdvice                    string    `middleware:"db_config_advice" json:"db_config_advice"`
	CPUUsageScore                     int       `middleware:"cpu_usage_score" json:"cpu_usage_score"`
	CPUUsageData                      string    `middleware:"cpu_usage_data" json:"cpu_usage_data"`
	CPUUsageHigh                      string    `middleware:"cpu_usage_high" json:"cpu_usage_high"`
	IOUtilScore                       int       `middleware:"io_util_score" json:"io_util_score"`
	IOUtilData                        string    `middleware:"io_util_data" json:"io_util_data"`
	IOUtilHigh                        string    `middleware:"io_util_high" json:"io_util_high"`
	DiskCapacityUsageScore            int       `middleware:"disk_capacity_usage_score" json:"disk_capacity_usage_score"`
	DiskCapacityUsageData             string    `middleware:"disk_capacity_usage_data" json:"disk_capacity_usage_data"`
	DiskCapacityUsageHigh             string    `middleware:"disk_capacity_usage_high" json:"disk_capacity_usage_high"`
	ConnectionUsageScore              int       `middleware:"connection_usage_score" json:"connection_usage_score"`
	ConnectionUsageData               string    `middleware:"connection_usage_data" json:"connection_usage_data"`
	ConnectionUsageHigh               string    `middleware:"connection_usage_high" json:"connection_usage_high"`
	AverageActiveSessionPercentsScore int       `middleware:"average_active_session_percents_score" json:"average_active_session_percents_score"`
	AverageActiveSessionPercentsData  string    `middleware:"average_active_session_percents_data" json:"average_active_session_percents_data"`
	AverageActiveSessionPercentsHigh  string    `middleware:"average_active_session_percents_high" json:"average_active_session_percents_high"`
	CacheMissRatioScore               int       `middleware:"cache_miss_ratio_score" json:"cache_miss_ratio_score"`
	CacheMissRatioData                string    `middleware:"cache_miss_ratio_data" json:"cache_miss_ratio_data"`
	CacheMissRatioHigh                string    `middleware:"cache_miss_ratio_high" json:"cache_miss_ratio_high"`
	TableRowsScore                    int       `middleware:"table_rows_score" json:"table_rows_score"`
	TableRowsData                     string    `middleware:"table_rows_data" json:"table_rows_data"`
	TableRowsHigh                     string    `middleware:"table_rows_high" json:"table_rows_high"`
	TableSizeScore                    int       `middleware:"table_size_score" json:"table_size_score"`
	TableSizeData                     string    `middleware:"table_size_data" json:"table_size_data"`
	TableSizeHigh                     string    `middleware:"table_size_high" json:"table_size_high"`
	SlowQueryScore                    int       `middleware:"slow_query_score" json:"slow_query_score"`
	SlowQueryData                     string    `middleware:"slow_query_data" json:"slow_query_data"`
	SlowQueryAdvice                   string    `middleware:"slow_query_advice" json:"slow_query_advice"`
	AccuracyReview                    int       `middleware:"accuracy_review" json:"accuracy_review"`
	DelFlag                           int       `middleware:"del_flag" json:"del_flag"`
	CreateTime                        time.Time `middleware:"create_time" json:"create_time"`
	LastUpdateTime                    time.Time `middleware:"last_update_time" json:"last_update_time"`
}

// NewResult returns a new *Result
func NewResult(repo healthcheck.DASRepo, operationID int, weightedAverageScore int, dbConfigScore int, dbConfigData string, dbConfigAdvice string,
	cpuUsageScore int, cpuUsageData string, cpuUsageHigh string, ioUtilScore int, ioUtilData string, ioUtilHigh string,
	diskCapacityUsageScore int, diskCapacityUsageData string, diskCapacityUsageHigh string,
	connectionUsageScore int, connectionUsageData string, connectionUsageHigh string,
	averageActiveSessionPercentsScore int, averageActiveSessionPercentsData string, averageActiveSessionPercentsHigh string,
	cacheMissRatioScore int, cacheMissRatioData string, cacheMissRatioHigh string,
	tableRowsScore int, tableRowsData string, tableRowsHigh string,
	tableSizeScore int, tableSizeData string, tableSizeHigh string,
	slowQueryScore int, slowQueryData string, slowQueryAdvice string) *Result {
	return &Result{
		DASRepo:                           repo,
		OperationID:                       operationID,
		WeightedAverageScore:              weightedAverageScore,
		DBConfigScore:                     dbConfigScore,
		DBConfigData:                      dbConfigData,
		DBConfigAdvice:                    dbConfigAdvice,
		CPUUsageScore:                     cpuUsageScore,
		CPUUsageData:                      cpuUsageData,
		CPUUsageHigh:                      cpuUsageHigh,
		IOUtilScore:                       ioUtilScore,
		IOUtilData:                        ioUtilData,
		IOUtilHigh:                        ioUtilHigh,
		DiskCapacityUsageScore:            diskCapacityUsageScore,
		DiskCapacityUsageData:             diskCapacityUsageData,
		DiskCapacityUsageHigh:             diskCapacityUsageHigh,
		ConnectionUsageScore:              connectionUsageScore,
		ConnectionUsageData:               connectionUsageData,
		ConnectionUsageHigh:               connectionUsageHigh,
		AverageActiveSessionPercentsScore: averageActiveSessionPercentsScore,
		AverageActiveSessionPercentsData:  averageActiveSessionPercentsData,
		AverageActiveSessionPercentsHigh:  averageActiveSessionPercentsHigh,
		CacheMissRatioScore:               cacheMissRatioScore,
		CacheMissRatioData:                cacheMissRatioData,
		CacheMissRatioHigh:                cacheMissRatioHigh,
		TableRowsScore:                    tableRowsScore,
		TableRowsData:                     tableRowsData,
		TableRowsHigh:                     tableRowsHigh,
		TableSizeScore:                    tableSizeScore,
		TableSizeData:                     tableSizeData,
		TableSizeHigh:                     tableSizeHigh,
		SlowQueryScore:                    slowQueryScore,
		SlowQueryData:                     slowQueryData,
		SlowQueryAdvice:                   slowQueryAdvice,
	}
}

// NewEmptyResultWithRepo return a new Result
func NewEmptyResultWithRepo(repository healthcheck.DASRepo) *Result {
	return &Result{DASRepo: repository}
}

// NewEmptyResultWithGlobal return a new Result
func NewEmptyResultWithGlobal() *Result {
	return NewEmptyResultWithRepo(NewDASRepoWithGlobal())
}

// NewResultWithDefault returns a new *Result with default DASRepo
func NewResultWithDefault(operationID int, weightedAverageScore int, dbConfigScore int,
	cpuUsageScore int, ioUtilScore int, diskCapacityUsageScore int, connectionUsageScore int,
	averageActiveSessionPercentsScore int, cacheMissRatioScore int, tableRowsScore int, tableSizeScore int,
	slowQueryScore int, accuracyReview int) *Result {
	return &Result{
		DASRepo:                           NewDASRepoWithGlobal(),
		OperationID:                       operationID,
		WeightedAverageScore:              weightedAverageScore,
		DBConfigScore:                     dbConfigScore,
		DBConfigData:                      constant.DefaultRandomString,
		DBConfigAdvice:                    constant.DefaultRandomString,
		CPUUsageScore:                     cpuUsageScore,
		CPUUsageData:                      constant.DefaultRandomString,
		CPUUsageHigh:                      constant.DefaultRandomString,
		IOUtilScore:                       ioUtilScore,
		IOUtilData:                        constant.DefaultRandomString,
		IOUtilHigh:                        constant.DefaultRandomString,
		DiskCapacityUsageScore:            diskCapacityUsageScore,
		DiskCapacityUsageData:             constant.DefaultRandomString,
		DiskCapacityUsageHigh:             constant.DefaultRandomString,
		ConnectionUsageScore:              connectionUsageScore,
		ConnectionUsageData:               constant.DefaultRandomString,
		ConnectionUsageHigh:               constant.DefaultRandomString,
		AverageActiveSessionPercentsScore: averageActiveSessionPercentsScore,
		AverageActiveSessionPercentsData:  constant.DefaultRandomString,
		AverageActiveSessionPercentsHigh:  constant.DefaultRandomString,
		CacheMissRatioScore:               cacheMissRatioScore,
		CacheMissRatioData:                constant.DefaultRandomString,
		CacheMissRatioHigh:                constant.DefaultRandomString,
		TableRowsScore:                    tableRowsScore,
		TableRowsData:                     constant.DefaultRandomString,
		TableRowsHigh:                     constant.DefaultRandomString,
		TableSizeScore:                    tableSizeScore,
		TableSizeData:                     constant.DefaultRandomString,
		TableSizeHigh:                     constant.DefaultRandomString,
		SlowQueryScore:                    slowQueryScore,
		SlowQueryData:                     constant.DefaultRandomString,
		SlowQueryAdvice:                   constant.DefaultRandomString,
		AccuracyReview:                    accuracyReview,
	}
}

// NewEmptyResult return a empty Result
func NewEmptyResult() *Result {
	return &Result{}
}

// Identity returns the identity
func (r *Result) Identity() int {
	return r.ID
}

// GetOperationID returns the OperationID
func (r *Result) GetOperationID() int {
	return r.OperationID
}

// GetWeightedAverageScore returns the WeightedAverageScore
func (r *Result) GetWeightedAverageScore() int {
	return r.WeightedAverageScore
}

// GetDBConfigScore returns the DBConfigScore
func (r *Result) GetDBConfigScore() int {
	return r.DBConfigScore
}

// GetDBConfigData returns the DBConfigData
func (r *Result) GetDBConfigData() string {
	return r.DBConfigData
}

// GetDBConfigAdvice returns the DBConfigAdvice
func (r *Result) GetDBConfigAdvice() string {
	return r.DBConfigAdvice
}

// GetCPUUsageScore returns the CPUUsageScore
func (r *Result) GetCPUUsageScore() int {
	return r.CPUUsageScore
}

// GetCPUUsageData returns the CPUUsageData
func (r *Result) GetCPUUsageData() string {
	return r.CPUUsageData
}

// GetCPUUsageHigh returns the CPUUsageHigh
func (r *Result) GetCPUUsageHigh() string {
	return r.CPUUsageHigh
}

// GetIOUtilScore returns the IOUtilScore
func (r *Result) GetIOUtilScore() int {
	return r.IOUtilScore
}

// GetIOUtilData returns the IOUtilData
func (r *Result) GetIOUtilData() string {
	return r.IOUtilData
}

// GetIOUtilHigh returns the IOUtilHigh
func (r *Result) GetIOUtilHigh() string {
	return r.IOUtilHigh
}

// GetDiskCapacityUsageScore returns the DiskCapacityUsageScore
func (r *Result) GetDiskCapacityUsageScore() int {
	return r.DiskCapacityUsageScore
}

// GetDiskCapacityUsageData returns the DiskCapacityUsageData
func (r *Result) GetDiskCapacityUsageData() string {
	return r.DiskCapacityUsageData
}

// GetDiskCapacityUsageHigh returns the DiskCapacityUsageHigh
func (r *Result) GetDiskCapacityUsageHigh() string {
	return r.DiskCapacityUsageHigh
}

// GetConnectionUsageScore returns the ConnectionUsageScore
func (r *Result) GetConnectionUsageScore() int {
	return r.ConnectionUsageScore
}

// GetConnectionUsageData returns the ConnectionUsageData
func (r *Result) GetConnectionUsageData() string {
	return r.ConnectionUsageData
}

// GetConnectionUsageHigh returns the ConnectionUsageHigh
func (r *Result) GetConnectionUsageHigh() string {
	return r.ConnectionUsageHigh
}

// GetAverageActiveSessionPercentsScore returns the AverageActiveSessionPercentsScore
func (r *Result) GetAverageActiveSessionPercentsScore() int {
	return r.AverageActiveSessionPercentsScore
}

// GetAverageActiveSessionPercentsData returns the AverageActiveSessionPercentsData
func (r *Result) GetAverageActiveSessionPercentsData() string {
	return r.AverageActiveSessionPercentsData
}

// GetAverageActiveSessionPercentsHigh returns the AverageActiveSessionPercentsHigh
func (r *Result) GetAverageActiveSessionPercentsHigh() string {
	return r.AverageActiveSessionPercentsHigh
}

// GetCacheMissRatioScore returns the CacheMissRatioScore
func (r *Result) GetCacheMissRatioScore() int {
	return r.CacheMissRatioScore
}

// GetCacheMissRatioData returns the CacheMissRatioData
func (r *Result) GetCacheMissRatioData() string {
	return r.CacheMissRatioData
}

// GetCacheMissRatioHigh returns the CacheMissRatioHigh
func (r *Result) GetCacheMissRatioHigh() string {
	return r.CacheMissRatioHigh
}

// GetTableRowsScore returns the TableRowsScore
func (r *Result) GetTableRowsScore() int {
	return r.TableRowsScore
}

// GetTableRowsData returns the TableRowsData
func (r *Result) GetTableRowsData() string {
	return r.TableRowsData
}

// GetTableRowsHigh returns the TableRowsHigh
func (r *Result) GetTableRowsHigh() string {
	return r.TableRowsHigh
}

// GetTableSizeScore returns the TableSizeScore
func (r *Result) GetTableSizeScore() int {
	return r.TableSizeScore
}

// GetTableSizeData returns the TableSizeData
func (r *Result) GetTableSizeData() string {
	return r.TableSizeData
}

// GetTableSizeHigh returns the TableSizeHigh
func (r *Result) GetTableSizeHigh() string {
	return r.TableSizeHigh
}

// GetSlowQueryScore returns the SlowQueryScore
func (r *Result) GetSlowQueryScore() int {
	return r.SlowQueryScore
}

// GetSlowQueryData returns the SlowQueryData
func (r *Result) GetSlowQueryData() string {
	return r.SlowQueryData
}

// GetSlowQueryAdvice returns the SlowQueryAdvice
func (r *Result) GetSlowQueryAdvice() string {
	return r.SlowQueryAdvice
}

// GetAccuracyReview returns the AccuracyReview
func (r *Result) GetAccuracyReview() int {
	return r.AccuracyReview
}

// GetDelFlag returns the delete flag
func (r *Result) GetDelFlag() int {
	return r.DelFlag
}

// GetCreateTime returns the create time
func (r *Result) GetCreateTime() time.Time {
	return r.CreateTime
}

// GetLastUpdateTime returns the last update time
func (r *Result) GetLastUpdateTime() time.Time {
	return r.LastUpdateTime
}

// Set sets health check with given fields, key is the field name and value is the relevant value of the key
func (r *Result) Set(fields map[string]interface{}) error {
	for fieldName, fieldValue := range fields {
		err := common.SetValueOfStruct(r, fieldName, fieldValue)
		if err != nil {
			return err
		}
	}

	return nil
}

// MarshalJSON marshals health check to json string
func (r *Result) MarshalJSON() ([]byte, error) {
	return common.MarshalStructWithTag(r, constant.DefaultMarshalTag)
}

// MarshalJSONWithFields marshals only specified field of the health check to json string
func (r *Result) MarshalJSONWithFields(fields ...string) ([]byte, error) {
	return common.MarshalStructWithFields(r, fields...)
}
