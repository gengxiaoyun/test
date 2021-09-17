package healthcheck

import (
	"fmt"
	"testing"
	"time"

	"github.com/romberli/das/global"
	"github.com/romberli/das/internal/app/metadata"
	"github.com/romberli/das/internal/dependency/healthcheck"
	"github.com/romberli/go-util/common"
	"github.com/romberli/go-util/constant"
	"github.com/romberli/go-util/middleware/clickhouse"
	"github.com/romberli/go-util/middleware/mysql"
	"github.com/romberli/go-util/middleware/prometheus"
	"github.com/stretchr/testify/assert"
)

const (
	defaultEngineConfigDBUser = "root"
	defaultEngineConfigDBPass = "root"

	applicationMysqlAddr   = "192.168.10.210:3306"
	applicationMysqlDBName = "performance_schema"
	applicationMysqlDBUser = "root"
	applicationMysqlDBPass = "mysql123"

	defaultMySQLServerID = 1
	defaultStep          = time.Minute
)

func TestDefaultEngineAll(t *testing.T) {
	TestDefaultEngineConfig_Validate(t)
	TestDefaultEngine_Run(t)
}

func TestDefaultEngineConfig_Validate(t *testing.T) {
	asst := assert.New(t)
	// load config
	sql := `
		select id, item_name, item_weight, low_watermark, high_watermark, unit, score_deduction_per_unit_high, max_score_deduction_high,
		score_deduction_per_unit_medium, max_score_deduction_medium, del_flag, create_time, last_update_time
		from t_hc_default_engine_config
		where del_flag = 0;
	`
	result, err := testDASRepo.Execute(sql)
	asst.Nil(err, common.CombineMessageWithError("test Validate() failed", err))
	defaultEngineConfigList := make([]*DefaultItemConfig, result.RowNumber())
	for i := range defaultEngineConfigList {
		defaultEngineConfigList[i] = NewEmptyDefaultItemConfig()
	}
	err = result.MapToStructSlice(defaultEngineConfigList, constant.DefaultMiddlewareTag)
	asst.Nil(err, common.CombineMessageWithError("test Validate() failed", err))
	entityList := NewEmptyDefaultEngineConfig()
	for i := range defaultEngineConfigList {
		itemName := defaultEngineConfigList[i].ItemName
		entityList[itemName] = defaultEngineConfigList[i]
	}
	// validate config
	validate := entityList.Validate()
	asst.Equal(nil, validate, "test Validate() failed")
}

func TestDefaultEngine_Run(t *testing.T) {
	asst := assert.New(t)
	startTime := time.Now().Add(-constant.Week)
	endTime := time.Now()

	id, err := testDASRepo.InitOperation(defaultMySQLServerID, startTime, endTime, defaultStep)
	asst.Nil(err, common.CombineMessageWithError("test Run() failed", err))

	mysqlServerService := metadata.NewMySQLServerService(metadata.NewMySQLServerRepo(global.DASMySQLPool))
	err = mysqlServerService.GetByID(1)
	asst.Nil(err, common.CombineMessageWithError("test Run() failed", err))
	mysqlServer := mysqlServerService.GetMySQLServers()[constant.ZeroInt]
	asst.Nil(err, common.CombineMessageWithError("test Run() failed", err))
	monitorSystem, err := mysqlServer.GetMonitorSystem()
	asst.Nil(err, common.CombineMessageWithError("test Run() failed", err))

	operationInfo := NewOperationInfo(id, mysqlServer, monitorSystem, startTime, endTime, defaultStep)

	// init application mysql connection
	applicationMySQLConn, err := mysql.NewConn(fmt.Sprintf("%s:%d", mysqlServer.GetHostIP(), mysqlServer.GetPortNum()), applicationMysqlDBName, applicationMysqlDBUser, applicationMysqlDBPass)
	asst.Nil(err, common.CombineMessageWithError("test Run() failed", err))
	// init application mysql repository
	applicationMySQLRepo := NewApplicationMySQLRepo(operationInfo, applicationMySQLConn)

	var (
		prometheusConfig prometheus.Config
		queryRepo        healthcheck.QueryRepo
	)

	prometheusAddr := fmt.Sprintf("%s:%d%s", monitorSystem.GetHostIP(), monitorSystem.GetPortNum(), monitorSystem.GetBaseURL())
	queryAddr := fmt.Sprintf("%s:%d", monitorSystem.GetHostIP(), monitorSystem.GetPortNumSlow())

	switch monitorSystem.GetSystemType() {
	case 1:
		// pmm 1.x
		// init prometheus config
		prometheusConfig = prometheus.NewConfig(prometheusAddr, prometheus.DefaultRoundTripper)
		asst.Nil(err, common.CombineMessageWithError("test Run() failed", err))
		// init mysql connection
		mysqlConn, err := mysql.NewConn(queryAddr, defaultMonitorMySQLDBName, defaultEngineConfigDBUser, defaultEngineConfigDBPass)
		asst.Nil(err, common.CombineMessageWithError("test Run() failed", err))
		// init mysql query repository
		queryRepo = NewMySQLQueryRepo(operationInfo, mysqlConn)
	case 2:
		// pmm 2.x
		// init prometheus config
		prometheusConfig = prometheus.NewConfigWithBasicAuth(prometheusAddr, defaultPrometheusUser, defaultPrometheusPass)
		asst.Nil(err, common.CombineMessageWithError("test Run() failed", err))
		// init clickhouse connection
		clickhouseConn, err := clickhouse.NewConnWithDefault(queryAddr, defaultMonitorClickhouseDBName, constant.EmptyString, constant.EmptyString)
		asst.Nil(err, common.CombineMessageWithError("test Run() failed", err))
		// init clickhouse query repository
		queryRepo = NewClickhouseQueryRepo(operationInfo, clickhouseConn)
	}

	// init prometheus connection
	prometheusConn, err := prometheus.NewConnWithConfig(prometheusConfig)
	asst.Nil(err, common.CombineMessageWithError("test Run() failed", err))
	// init prometheus repository
	prometheusRepo := NewPrometheusRepo(operationInfo, prometheusConn)

	defaultEngine := NewDefaultEngine(operationInfo, testDASRepo, applicationMySQLRepo, prometheusRepo, queryRepo)
	err = defaultEngine.run()
	asst.Nil(err, common.CombineMessageWithError("test Run() failed", err))
}
