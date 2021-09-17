package healthcheck

import (
	"testing"
	"time"

	"github.com/romberli/das/config"
	"github.com/romberli/go-util/common"
	"github.com/romberli/go-util/constant"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

const (
	idStruct             = "id"
	clusterIDStruct      = "ClusterID"
	serverNameStruct     = "ServerName"
	hostIPStruct         = "HostIP"
	portNumStruct        = "PortNum"
	deploymentTypeStruct = "DeploymentType"
	versionStruct        = "Version"

	defaultMySQLServerInfoClusterID      = 1
	defaultMySQLServerInfoServerName     = "server1"
	defaultMySQLServerInfoHostIP         = "127.0.0.1"
	defaultMySQLServerInfoPortNum        = 3307
	defaultMySQLServerInfoDeploymentType = 1
	defaultMySQLServerInfoVersion        = "1.1.1"
)

func createService() (*Service, error) {
	initViper()

	var result = NewResult(testDASRepo,
		defaultResultOperationID,
		defaultResultWeightedAverageScore,
		defaultResultDBConfigScore,
		defaultResultDBConfigData,
		defaultResultDBConfigAdvice,
		defaultResultCPUUsageScore,
		defaultResultCPUUsageData,
		defaultResultCPUUsageHigh,
		defaultResultIOUtilScore,
		defaultResultIOUtilData,
		defaultResultIOUtilHigh,
		defaultResultDiskCapacityUsageScore,
		defaultResultDiskCapacityUsageData,
		defaultResultDiskCapacityUsageHigh,
		defaultResultConnectionUsageScore,
		defaultResultConnectionUsageData,
		defaultResultConnectionUsageHigh,
		defaultResultAverageActiveSessionPercentsScore,
		defaultResultAverageActiveSessionPercentsData,
		defaultResultAverageActiveSessionPercentsHigh,
		defaultResultCacheMissRatioScore,
		defaultResultCacheMissRatioData,
		defaultResultCacheMissRatioHigh,
		defaultResultTableRowsScore,
		defaultResultTableRowsData,
		defaultResultTableRowsHigh,
		defaultResultTableSizeScore,
		defaultResultTableSizeData,
		defaultResultTableSizeHigh,
		defaultResultSlowQueryScore,
		defaultResultSlowQueryData,
		defaultResultSlowQueryAdvice)
	err := testDASRepo.SaveResult(result)
	if err != nil {
		return nil, err
	}
	return &Service{
		DASRepo: testDASRepo,
		Result:  result,
	}, nil

}

func initViper() {
	viper.Set(config.DBApplicationMySQLUserKey, config.DefaultDBApplicationMySQLUser)
	viper.Set(config.DBApplicationMySQLPassKey, config.DefaultDBApplicationMySQLPass)

	viper.Set(config.DBMonitorPrometheusUserKey, config.DefaultDBMonitorPrometheusUser)
	viper.Set(config.DBMonitorPrometheusPassKey, config.DefaultDBMonitorPrometheusPass)

	viper.Set(config.DBMonitorMySQLUserKey, config.DefaultDBMonitorMySQLUser)
	viper.Set(config.DBMonitorMySQLPassKey, config.DefaultDBMonitorMySQLPass)

	viper.Set(config.DBMonitorClickhouseUserKey, config.DefaultDBMonitorClickhouseUser)
	viper.Set(config.DBMonitorClickhousePassKey, config.DefaultDBMonitorClickhousePass)
}

func deleteHCResultByOperationID(operationID int) error {
	sql := `delete from t_hc_result where operation_id = ?`
	_, err := testDASRepo.Execute(sql, operationID)
	return err
}

func TestServiceAll(t *testing.T) {
	TestService_GetResult(t)
	TestService_GetResultByOperationID(t)
	TestService_Check(t)
	TestService_ReviewAccuracy(t)
	TestService_MarshalJSON(t)
	TestService_MarshalJSONWithFields(t)
}

func TestService_GetResult(t *testing.T) {
	asst := assert.New(t)

	service, err := createService()
	asst.Nil(err, common.CombineMessageWithError("test GetResult() failed", err))
	result := service.GetResult()
	asst.Equal(defaultResultOperationID, result.GetOperationID(), common.CombineMessageWithError("test GetResult() failed", err))
	asst.Equal(defaultResultWeightedAverageScore, result.GetWeightedAverageScore(), common.CombineMessageWithError("test GetResult() failed", err))
	// delete
	err = deleteHCResultByOperationID(defaultResultOperationID)
	asst.Nil(err, common.CombineMessageWithError("test GetResult() failed", err))
}

func TestService_GetResultByOperationID(t *testing.T) {
	asst := assert.New(t)

	service, err := createService()
	asst.Nil(err, common.CombineMessageWithError("test GetResultByOperationID() failed", err))
	err = service.GetResultByOperationID(defaultResultOperationID)
	asst.Nil(err, common.CombineMessageWithError("test GetResultByOperationID() failed", err))
	result := service.GetResult()
	asst.Equal(defaultResultOperationID, result.GetOperationID(), common.CombineMessageWithError("test GetResultByOperationID() failed", err))
	asst.Equal(defaultResultWeightedAverageScore, result.GetWeightedAverageScore(), common.CombineMessageWithError("test GetResultByOperationID() failed", err))
	// delete
	err = deleteHCResultByOperationID(defaultResultOperationID)
	asst.Nil(err, common.CombineMessageWithError("test GetResultByOperationID() failed", err))
}

// bug
func TestService_Check(t *testing.T) {
	asst := assert.New(t)

	err := initGlobalMySQLPool()
	asst.Nil(err, common.CombineMessageWithError("test GetResultByOperationID() failed", err))

	service, err := createService()
	asst.Nil(err, common.CombineMessageWithError("test GetResultByOperationID() failed", err))

	err = service.Check(defaultMysqlServerID, time.Now().Add(-constant.Week), time.Now(), defaultStep)
	asst.Nil(err, common.CombineMessageWithError("test Check(mysqlServerID int, startTime, endTime time.Time, step time.Duration) failed", err))

	// delete
	err = deleteHCResultByOperationID(defaultResultOperationID)
	asst.Nil(err, common.CombineMessageWithError("test GetResultByOperationID() failed", err))
}

// bug
func TestService_CheckByHostInfo(t *testing.T) {
	// asst := assert.New(t)

	// service, err := createService()
	// asst.Nil(err, common.CombineMessageWithError("test CheckByHostInfo(hostIP string, portNum int, startTime, endTime time.Time, step time.Duration) failed", err))

	// startTime, _ := now.Parse(serviceStartTime)
	// endTime, _ := now.Parse(serviceEndTime)
	// step := time.Duration(serviceStep) * time.Second

	// err = service.CheckByHostInfo(serviceHostIP, servicePortNum, startTime, endTime, step)
	// asst.Nil(err, common.CombineMessageWithError("test CheckByHostInfo(hostIP string, portNum int, startTime, endTime time.Time, step time.Duration) failed", err))

	// // delete
	// err = deleteHCResultByOperationID(serviceOperationID)
	// asst.Nil(err, common.CombineMessageWithError("test CheckByHostInfo(hostIP string, portNum int, startTime, endTime time.Time, step time.Duration) failed", err))
}

func TestService_ReviewAccuracy(t *testing.T) {
	asst := assert.New(t)

	service, err := createService()
	asst.Nil(err, common.CombineMessageWithError("test ReviewAccuracy(id, review int) failed", err))
	review := 2
	err = service.ReviewAccuracy(defaultResultOperationID, review)
	asst.Nil(err, common.CombineMessageWithError("test ReviewAccuracy(id, review int) failed", err))
	err = service.GetResultByOperationID(defaultResultOperationID)
	result := service.GetResult()
	reviewed := result.GetAccuracyReview()
	asst.Equal(review, reviewed, common.CombineMessageWithError("test ReviewAccuracy(id, review int) failed", err))
	// delete
	err = deleteHCResultByOperationID(defaultResultOperationID)
	asst.Nil(err, common.CombineMessageWithError("test ReviewAccuracy(id, review int) failed", err))
}

func TestService_MarshalJSON(t *testing.T) {
	asst := assert.New(t)

	service, err := createService()
	asst.Nil(err, common.CombineMessageWithError("test MarshalJSON() failed", err))
	_, err = service.MarshalJSON()
	asst.Nil(err, common.CombineMessageWithError("test MarshalJSON() failed", err))
	// delete
	err = deleteHCResultByOperationID(defaultResultOperationID)
	asst.Nil(err, common.CombineMessageWithError("test MarshalJSON() failed", err))
}

func TestService_MarshalJSONWithFields(t *testing.T) {
	asst := assert.New(t)

	service, err := createService()
	asst.Nil(err, common.CombineMessageWithError("test MarshalJSONWithFields(fields ...string) failed", err))
	_, err = service.MarshalJSONWithFields("ID", "operationID", "WeightedAverageScore", "DBConfigScore", "DBConfigData", "DBConfigAdvice", "CPUUsageScore", "CPUUsageData", "CPUUsageHigh", "IOUtilScore", "IOUtilData", "IOUtilHigh", "DiskCapacityUsageScore", "DiskCapacityUsageData", "DiskCapacityUsageHigh", "ConnectionUsageScore", "ConnectionUsageData", "ConnectionUsageHigh", "AverageActiveSessionPercentsScore", "AverageActiveSessionPercentsData", "AverageActiveSessionPercentsHigh", "CacheMissRatioScore", "CacheMissRatioData", "CacheMissRatioHigh", "TableSizeScore", "TableSizeData", "TableSizeHigh", "SlowQueryScore", "SlowQueryData", "SlowQueryAdvice")
	asst.Nil(err, common.CombineMessageWithError("test MarshalJSONWithFields(fields ...string) failed", err))
	// delete
	err = deleteHCResultByOperationID(defaultResultOperationID)
	asst.Nil(err, common.CombineMessageWithError("test MarshalJSONWithFields(fields ...string) failed", err))
}

// go test ./service_test.go ./service.go ./query.go ./default_engine.go ./result.go
