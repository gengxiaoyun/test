package router

import (
	"github.com/gin-gonic/gin"
	"github.com/romberli/das/api/v1/query"
)

// RegisterQuery is the sub-router of das for query
func RegisterQuery(group *gin.RouterGroup) {
	queryGroup := group.Group("/query")
	{
		queryGroup.GET("/cluster/:mysqlClusterID", query.GetByMySQLClusterID)
		queryGroup.GET("/server/:mysqlServerID", query.GetByMySQLServerID)
		queryGroup.GET("/db/:dbID", query.GetByDBID)
		queryGroup.GET("/:id", query.GetBySQLID)
	}
}
