package db

import (
	"fmt"
	"github.com/WeBankPartners/go-common-lib/guid"
	"github.com/WeBankPartners/open-monitor/monitor-server/models"
	"time"
)

func ListEndpointGroup(param *models.QueryRequestParam) (pageInfo models.PageInfo, rowData []*models.EndpointGroupTable, err error) {
	rowData = []*models.EndpointGroupTable{}
	filterSql, queryColumn, queryParam := transFiltersToSQL(param, &models.TransFiltersParam{IsStruct: true, StructObj: models.EndpointGroupTable{}, PrimaryKey: "guid"})
	baseSql := fmt.Sprintf("SELECT %s FROM endpoint_group WHERE 1=1 %s ", queryColumn, filterSql)
	if param.Paging {
		pageInfo.StartIndex = param.Pageable.StartIndex
		pageInfo.PageSize = param.Pageable.PageSize
		pageInfo.TotalRows = queryCount(baseSql, queryParam...)
		pageSql, pageParam := transPageInfoToSQL(*param.Pageable)
		baseSql += pageSql
		queryParam = append(queryParam, pageParam...)
	}
	err = x.SQL(baseSql, queryParam...).Find(&rowData)
	return
}

func CreateEndpointGroup(param *models.EndpointGroupTable) error {
	var actions []*Action
	nowTime := time.Now().Format(models.DatetimeFormat)
	if param.DisplayName == "" {
		param.DisplayName = param.Guid
	}
	if param.ServiceGroup == "" {
		actions = append(actions, &Action{Sql: "insert into endpoint_group(guid,display_name,description,monitor_type,update_time) value (?,?,?,?,?)", Param: []interface{}{param.Guid, param.DisplayName, param.Description, param.MonitorType, nowTime}})
	} else {
		actions = append(actions, &Action{Sql: "insert into endpoint_group(guid,display_name,description,monitor_type,service_group,update_time) value (?,?,?,?,?,?)", Param: []interface{}{param.Guid, param.DisplayName, param.Description, param.MonitorType, param.ServiceGroup, nowTime}})
	}
	actions = append(actions, &Action{Sql: "insert into grp(name,description,endpoint_type,create_at) value (?,?,?,NOW())", Param: []interface{}{param.Guid, param.Description, param.MonitorType}})
	return Transaction(actions)
}

func UpdateEndpointGroup(param *models.EndpointGroupTable) error {
	var actions []*Action
	nowTime := time.Now().Format(models.DatetimeFormat)
	if param.DisplayName == "" {
		param.DisplayName = param.Guid
	}
	if param.ServiceGroup == "" {
		actions = append(actions, &Action{Sql: "update endpoint_group set display_name=?,description=?,monitor_type=?,service_group=NULL,update_time=? where guid=?", Param: []interface{}{param.DisplayName, param.Description, param.MonitorType, nowTime, param.Guid}})
	} else {
		actions = append(actions, &Action{Sql: "update endpoint_group set display_name=?,description=?,monitor_type=?,service_group=?,update_time=? where guid=?", Param: []interface{}{param.DisplayName, param.Description, param.MonitorType, param.ServiceGroup, nowTime, param.Guid}})
	}
	actions = append(actions, &Action{Sql: "update grp set description=? where name=?", Param: []interface{}{param.Description, param.Guid}})
	return Transaction(actions)
}

func DeleteEndpointGroup(endpointGroupGuid string) error {
	var actions []*Action
	actions = append(actions, &Action{Sql: "delete from endpoint_group where guid=?", Param: []interface{}{endpointGroupGuid}})
	actions = append(actions, &Action{Sql: "delete from grp where name=?", Param: []interface{}{endpointGroupGuid}})
	return Transaction(actions)
}

func GetGroupEndpointRel(endpointGroupGuid string) (result []*models.EndpointGroupRelTable, err error) {
	result = []*models.EndpointGroupRelTable{}
	err = x.SQL("select endpoint from endpoint_group_rel where endpoint_group=?", endpointGroupGuid).Find(&result)
	return
}

func UpdateGroupEndpoint(param *models.UpdateGroupEndpointParam) error {
	var actions []*Action
	actions = append(actions, &Action{Sql: "delete from endpoint_group_rel where endpoint_group=?", Param: []interface{}{param.GroupGuid}})
	guidList := guid.CreateGuidList(len(param.EndpointGuidList))
	for i, v := range param.EndpointGuidList {
		actions = append(actions, &Action{Sql: "insert into endpoint_group_rel(guid,endpoint,endpoint_group) value (?,?,?)", Param: []interface{}{guidList[i], v, param.GroupGuid}})
	}
	err := Transaction(actions)
	return err
}

func GetGroupEndpointNotify(endpointGroupGuid string) (result []*models.NotifyObj, err error) {
	result = getNotifyList("", endpointGroupGuid, "")
	return result, nil
}

func UpdateGroupEndpointNotify(endpointGroupGuid string, param []*models.NotifyObj) error {
	actions := getNotifyListDeleteAction("", endpointGroupGuid, "")
	actions = append(actions, getNotifyListInsertAction(param)...)
	return Transaction(actions)
}

func ListEndpointGroupOptions(searchText string) (result []*models.OptionModel, err error) {
	result = []*models.OptionModel{}
	if searchText == "." {
		searchText = ""
	}
	searchText = "%" + searchText + "%"
	var endpointGroupTable []*models.EndpointGroupTable
	err = x.SQL("select guid,monitor_type from endpoint_group where guid like ?", searchText).Find(&endpointGroupTable)
	if err != nil {
		return
	}
	for _, v := range endpointGroupTable {
		result = append(result, &models.OptionModel{OptionValue: v.Guid, OptionText: v.Guid, OptionType: v.MonitorType, OptionTypeName: v.MonitorType})
	}
	return
}