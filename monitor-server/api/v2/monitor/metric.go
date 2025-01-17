package monitor

import (
	"encoding/json"
	"fmt"
	"github.com/WeBankPartners/open-monitor/monitor-server/middleware"
	"github.com/WeBankPartners/open-monitor/monitor-server/models"
	"github.com/WeBankPartners/open-monitor/monitor-server/services/db"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"time"
)

func ListMetric(c *gin.Context) {
	guid := c.Query("guid")
	monitorType := c.Query("monitorType")
	serviceGroup := c.Query("serviceGroup")
	onlyService := c.Query("onlyService")
	result, err := db.MetricListNew(guid, monitorType, serviceGroup, onlyService)
	if err != nil {
		middleware.ReturnHandleError(c, err.Error(), err)
	} else {
		middleware.ReturnSuccessData(c, result)
	}
}

func GetSysMetricTemplate(c *gin.Context) {
	workspace := c.Query("workspace")
	result, err := db.GetSysMetricTemplateConfig(workspace)
	if err != nil {
		middleware.ReturnHandleError(c, err.Error(), err)
	} else {
		middleware.ReturnSuccessData(c, result)
	}
}

func ExportMetric(c *gin.Context) {
	serviceGroup := c.Query("serviceGroup")
	monitorType := c.Query("monitorType")
	result, err := db.MetricListNew("", monitorType, serviceGroup, "Y")
	if err != nil {
		middleware.ReturnHandleError(c, err.Error(), err)
		return
	}
	b, marshalErr := json.Marshal(result)
	if marshalErr != nil {
		middleware.ReturnHandleError(c, "export metric fail, json marshal object error", marshalErr)
		return
	}
	c.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s_%s_%s.json", "metric_", serviceGroup, time.Now().Format("20060102150405")))
	c.Data(http.StatusOK, "application/octet-stream", b)
}

func ImportMetric(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		middleware.ReturnValidateError(c, err.Error())
		return
	}
	f, err := file.Open()
	if err != nil {
		middleware.ReturnHandleError(c, "file open error ", err)
		return
	}
	var paramObj []*models.MetricTable
	b, err := ioutil.ReadAll(f)
	defer f.Close()
	if err != nil {
		middleware.ReturnHandleError(c, "read content fail error ", err)
		return
	}
	err = json.Unmarshal(b, &paramObj)
	if err != nil {
		middleware.ReturnHandleError(c, "json unmarshal fail error ", err)
		return
	}
	if len(paramObj) == 0 {
		middleware.ReturnValidateError(c, "can not import empty file")
		return
	}
	serviceGroup := c.Query("serviceGroup")
	if serviceGroup == "" {
		middleware.ReturnValidateError(c, "serviceGroup can not empty")
		return
	}
	if err = db.MetricImport(serviceGroup, paramObj); err != nil {
		middleware.ReturnHandleError(c, "import metric fail", err)
	} else {
		middleware.ReturnSuccess(c)
	}
}
