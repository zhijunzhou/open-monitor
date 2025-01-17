package models

import (
	"fmt"
	"reflect"
	"strings"
)

var GlobalSGDisplayNameMap = make(map[string]string)

type ServiceGroupTable struct {
	Guid        string `json:"guid" xorm:"guid"`
	DisplayName string `json:"display_name" xorm:"display_name"`
	Description string `json:"description" xorm:"description"`
	Parent      string `json:"parent" xorm:"parent"`
	ServiceType string `json:"service_type" xorm:"service_type"`
	UpdateTime  string `json:"update_time" xorm:"update_time"`
}

type EndpointServiceRelTable struct {
	Guid         string `json:"guid" xorm:"guid"`
	Endpoint     string `json:"endpoint" xorm:"endpoint"`
	ServiceGroup string `json:"service_group" xorm:"service_group"`
}

type ServiceGroupLinkNode struct {
	Guid     string                  `json:"guid"`
	Parent   *ServiceGroupLinkNode   `json:"parent"`
	Children []*ServiceGroupLinkNode `json:"children"`
}

func (s *ServiceGroupLinkNode) FetchChildGuid() []string {
	result := []string{s.Guid}
	for _, v := range s.Children {
		result = append(result, v.FetchChildGuid()...)
	}
	return result
}

func (s *ServiceGroupLinkNode) FetchParentGuid() []string {
	result := []string{s.Guid}
	if s.Parent != nil {
		result = append(result, s.Parent.FetchParentGuid()...)
	}
	return result
}

type ServiceGroupEndpointListObj struct {
	Guid        string `json:"guid"`
	DisplayName string `json:"display_name"`
	Type        string `json:"type"`
}

type ServiceGroupRoleRelTable struct {
	Guid         string `json:"guid" xorm:"guid"`
	ServiceGroup string `json:"service_group" xorm:"service_group"`
	Role         string `json:"role" xorm:"role"`
}

type PluginUpdateServicePathRequest struct {
	RequestId      string                               `json:"requestId"`
	DueDate        string                               `json:"dueDate"`
	AllowedOptions []string                             `json:"allowedOptions"`
	Inputs         []*PluginUpdateServicePathRequestObj `json:"inputs"`
}

type PluginUpdateServicePathRequestObj struct {
	CallbackParameter string      `json:"callbackParameter"`
	Guid              string      `json:"guid"`
	SystemName        string      `json:"systemName"`
	LogPathList       interface{} `json:"logPathList"`
	MonitorType       string      `json:"monitorType"`
}

type PluginUpdateServicePathResp struct {
	ResultCode    string                        `json:"resultCode"`
	ResultMessage string                        `json:"resultMessage"`
	Results       PluginUpdateServicePathOutput `json:"results"`
}

type PluginUpdateServicePathOutput struct {
	RequestId      string                              `json:"requestId"`
	AllowedOptions []string                            `json:"allowedOptions,omitempty"`
	Outputs        []*PluginUpdateServicePathOutputObj `json:"outputs"`
}

type PluginUpdateServicePathOutputObj struct {
	CallbackParameter string `json:"callbackParameter"`
	Guid              string `json:"guid"`
	ErrorCode         string `json:"errorCode"`
	ErrorMessage      string `json:"errorMessage"`
	ErrorDetail       string `json:"errorDetail,omitempty"`
}

func TransPluginMultiStringParam(input interface{}) []string {
	var result []string
	rn := reflect.TypeOf(input).String()
	if strings.Contains(rn, "[]interface") {
		for _, v := range input.([]interface{}) {
			tmpV := fmt.Sprintf("%s", v)
			if tmpV != "" {
				result = append(result, tmpV)
			}
		}
	} else {
		tmpString := fmt.Sprintf("%s", input)
		if strings.HasPrefix(tmpString, "[") {
			tmpString = tmpString[1:]
		}
		if strings.HasSuffix(tmpString, "]") {
			tmpString = tmpString[:len(tmpString)-1]
		}
		for _, v := range strings.Split(tmpString, ",") {
			if v != "" {
				result = append(result, v)
			}
		}
	}
	return result
}
