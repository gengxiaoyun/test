package query

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/romberli/das/internal/app/query"
	"github.com/romberli/das/pkg/message"
	msgquery "github.com/romberli/das/pkg/message/query"
	"github.com/romberli/das/pkg/resp"
	util "github.com/romberli/das/pkg/util/query"
	"github.com/romberli/go-util/constant"
	"github.com/romberli/log"
)

const (
	mysqlClusterIDJSON = "mysql_cluster_id"
	mysqlServerIDJSON  = "mysql_server_id"
	dbIDJSON           = "db_id"
	sqlIDJSON          = "sql_id"
)

// @Tags query
// @Summary get slow queries by mysql server id
// @Produce  application/json
// @Success 200 {string} string "{"code": 200, "data": []}"
// @Router /api/v1/query/cluster/:mysqlClusterID [get]
func GetByMySQLClusterID(c *gin.Context) {
	// get data
	mysqlClusterIDStr := c.Param(mysqlClusterIDJSON)
	if mysqlClusterIDStr == constant.EmptyString {
		resp.ResponseNOK(c, message.ErrFieldNotExists, mysqlClusterIDJSON)
		return
	}
	mysqlClusterID, err := strconv.Atoi(mysqlClusterIDStr)
	if err != nil {
		resp.ResponseNOK(c, message.ErrTypeConversion, err)
		return
	}

	data, err := c.GetRawData()
	if err != nil {
		resp.ResponseNOK(c, message.ErrGetRawData, err.Error())
		return
	}
	dataMap := make(map[string]string)
	err = json.Unmarshal(data, &dataMap)
	if err != nil {
		resp.ResponseNOK(c, message.ErrUnmarshalRawData, err.Error())
		return
	}
	dataMap = make(map[string]string)
	err = json.Unmarshal(data, &dataMap)
	if err != nil {
		resp.ResponseNOK(c, message.ErrUnmarshalRawData, err.Error())
		return
	}

	// get config
	config, err := util.GetConfig(dataMap)
	if err != nil {
		resp.ResponseNOK(c, message.ErrUnmarshalRawData, err.Error())
		return
	}

	// init server
	service := query.NewServiceWithDefault(config)
	err = service.GetByMySQLClusterID(mysqlClusterID)
	if err != nil {
		resp.ResponseNOK(c, msgquery.ErrQueryGetByMySQLClusterID, mysqlClusterID, err.Error())
		return
	}

	// marshal
	jsonBytes, err := service.Marshal()
	if err != nil {
		resp.ResponseNOK(c, message.ErrMarshalData, err.Error())
	}
	jsonStr := string(jsonBytes)
	log.Debug(message.NewMessage(msgquery.DebugQueryGetByMySQLClusterID, mysqlClusterID, jsonStr).Error())

	// response
	resp.ResponseOK(c, jsonStr, msgquery.InfoQueryGetByMySQLClusterID, mysqlClusterID)
}

// @Tags query
// @Summary get slow queries by mysql server id
// @Produce  application/json
// @Success 200 {string} string "{"code": 200, "data": []}"
// @Router /api/v1/query/server/:mysqlServerID [get]
func GetByMySQLServerID(c *gin.Context) {
	// get data
	mysqlServerIDStr := c.Param(mysqlServerIDJSON)
	if mysqlServerIDStr == constant.EmptyString {
		resp.ResponseNOK(c, message.ErrFieldNotExists, mysqlServerIDJSON)
		return
	}
	mysqlServerID, err := strconv.Atoi(mysqlServerIDStr)
	if err != nil {
		resp.ResponseNOK(c, message.ErrTypeConversion, err)
		return
	}

	data, err := c.GetRawData()
	if err != nil {
		resp.ResponseNOK(c, message.ErrGetRawData, err.Error())
		return
	}
	dataMap := make(map[string]string)
	err = json.Unmarshal(data, &dataMap)
	if err != nil {
		resp.ResponseNOK(c, message.ErrUnmarshalRawData, err.Error())
		return
	}

	// get config
	config, err := util.GetConfig(dataMap)
	if err != nil {
		resp.ResponseNOK(c, message.ErrUnmarshalRawData, err.Error())
		return
	}

	// init service
	service := query.NewServiceWithDefault(config)
	err = service.GetByMySQLServerID(mysqlServerID)
	if err != nil {
		resp.ResponseNOK(c, msgquery.ErrQueryGetByMySQLServerID, mysqlServerID, err.Error())
		return
	}

	// marshal
	jsonBytes, err := service.Marshal()
	if err != nil {
		resp.ResponseNOK(c, message.ErrMarshalData, err.Error())
	}
	jsonStr := string(jsonBytes)
	log.Debug(message.NewMessage(msgquery.DebugQueryGetByMySQLServerID, mysqlServerID, jsonStr).Error())

	// response
	resp.ResponseOK(c, jsonStr, msgquery.InfoQueryGetByMySQLServerID, mysqlServerID)
}

// @Tags query
// @Summary get slow queries by db id
// @Produce  application/json
// @Success 200 {string} string "{"code": 200, "data": []}"
// @Router /api/v1/query/db/:dbID [get]
func GetByDBID(c *gin.Context) {
	// get data
	dbIDStr := c.Param(dbIDJSON)
	if dbIDStr == constant.EmptyString {
		resp.ResponseNOK(c, message.ErrFieldNotExists, dbIDJSON)
		return
	}
	dbID, err := strconv.Atoi(dbIDStr)
	if err != nil {
		resp.ResponseNOK(c, message.ErrTypeConversion, err)
		return
	}

	data, err := c.GetRawData()
	if err != nil {
		resp.ResponseNOK(c, message.ErrGetRawData, err.Error())
		return
	}
	dataMap := make(map[string]string)
	err = json.Unmarshal(data, &dataMap)
	if err != nil {
		resp.ResponseNOK(c, message.ErrUnmarshalRawData, err.Error())
		return
	}

	// get config
	config, err := util.GetConfig(dataMap)
	if err != nil {
		resp.ResponseNOK(c, message.ErrUnmarshalRawData, err.Error())
		return
	}

	// get mysqlServerID

	mysqlServerIDStr, exists := dataMap[mysqlServerIDJSON]
	if !exists {
		resp.ResponseNOK(c, message.ErrUnmarshalRawData, fmt.Errorf("%s not exists", mysqlClusterIDJSON))
	}

	mysqlServerID, err := strconv.Atoi(mysqlServerIDStr)
	if err != nil {
		resp.ResponseNOK(c, message.ErrTypeConversion, err)
		return
	}

	// init service
	service := query.NewServiceWithDefault(config)
	err = service.GetByDBID(mysqlServerID, dbID) //
	if err != nil {
		resp.ResponseNOK(c, msgquery.DebugQueryGetByDBID, dbID, err.Error())
		return
	}

	// marshal
	jsonBytes, err := service.Marshal()
	if err != nil {
		resp.ResponseNOK(c, message.ErrMarshalData, err.Error())
	}
	jsonStr := string(jsonBytes)
	log.Debug(message.NewMessage(msgquery.DebugQueryGetByDBID, dbID, jsonStr).Error())

	// response
	resp.ResponseOK(c, jsonStr, msgquery.DebugQueryGetByDBID, dbID)
}

// @Tags query
// @Summary get slow query by query id
// @Produce  application/json
// @Success 200 {string} string "{"code": 200, "data": []}"
// @Router /api/v1/query/:sqlID [get]
func GetBySQLID(c *gin.Context) {
	// get data
	sqlIDStr := c.Param(sqlIDJSON)
	if sqlIDStr == constant.EmptyString {
		resp.ResponseNOK(c, message.ErrFieldNotExists, sqlIDJSON)
		return
	}

	data, err := c.GetRawData()
	if err != nil {
		resp.ResponseNOK(c, message.ErrGetRawData, err.Error())
		return
	}
	dataMap := make(map[string]string)
	err = json.Unmarshal(data, &dataMap)
	if err != nil {
		resp.ResponseNOK(c, message.ErrUnmarshalRawData, err.Error())
		return
	}

	// get mysqlServerID

	mysqlServerIDStr, exists := dataMap[mysqlServerIDJSON]
	if !exists {
		resp.ResponseNOK(c, message.ErrUnmarshalRawData, fmt.Errorf("%s not exists", mysqlClusterIDJSON))
	}

	mysqlServerID, err := strconv.Atoi(mysqlServerIDStr)
	if err != nil {
		resp.ResponseNOK(c, message.ErrTypeConversion, err)
		return
	}

	// get config
	config, err := util.GetConfig(dataMap)
	if err != nil {
		resp.ResponseNOK(c, message.ErrUnmarshalRawData, err.Error())
	}

	// init service
	service := query.NewServiceWithDefault(config)
	err = service.GetBySQLID(mysqlServerID, sqlIDStr)
	if err != nil {
		resp.ResponseNOK(c, msgquery.DebugQueryGetBySQLID, sqlIDStr, err.Error())
		return
	}

	// marshal
	jsonBytes, err := service.Marshal()
	if err != nil {
		resp.ResponseNOK(c, message.ErrMarshalData, err.Error())
	}
	jsonStr := string(jsonBytes)
	log.Debug(message.NewMessage(msgquery.DebugQueryGetBySQLID, sqlIDStr).Error())

	// response
	resp.ResponseOK(c, jsonStr, msgquery.InfoQueryGetBySQLID, sqlIDStr)
}
