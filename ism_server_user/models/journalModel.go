/**
 * @ Author: ISM Web组态软件
 * @ Create Time: 2023-01-09 08:53:25
 * @ Modified by: ISM Web组态软件
 * @ Modified time: 2023-04-03 08:59:00
 * @ Description: 此源码版权归 www.ismctl.com 所有,个人私自不得二次销售。
 */

package models

import (
	"ISMServer/utils/errmsg"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
)

type SystemJournal struct {
	gorm.Model
	Content      string    `gorm:"type:text;not null"  json:"Content" validate:"required,min=4,max=250" label:"日志内容"`
	Time         time.Time `gorm:"index;type:datetime;not null" json:"Time" validate:"required" label:"发生时间"`
	JournalType  int       `gorm:"index;type:int;not null" json:"JournalType" validate:"required" label:"日志类型"`
	JournalLevel int       `gorm:"type:int;not null" json:"JournalLevel" validate:"required" label:"日志等级"`
	Operator     string    `gorm:"index;type:varchar(250);not null" json:"Operator" validate:"required,min=2,max=250" label:"操作者"`
	ProjectUuid  string    `gorm:"index;type:varchar(250);" json:"ProjectUuid" validate:"required,min=2,max=250" label:"项目ID"`
	ClientInfo   string    `gorm:"index;type:varchar(250);" json:"ClientInfo" validate:"required,min=2,max=250" label:"客户端信息"`
	UserName     string    `gorm:"index;type:varchar(250);" json:"UserName" validate:"required,min=2,max=250" label:"操作者的ID"`
}

// 日志获取
func GetJournalList(params map[string]interface{}) ([]SystemJournal, int) {

	var getJournal []SystemJournal

	var queryStartTime string
	var queryEndTime string

	dateType := params["dateType"]
	var buildStart strings.Builder
	var buildEnd strings.Builder
	for k, v := range params {
		switch value := v.(type) {
		case string:
			if k == "dateRange" {
				if dateType == "Day" {
					buildStart.WriteString(value)
					buildStart.WriteString(" ")
					buildStart.WriteString("00:00:00")
					queryStartTime = buildStart.String()

					buildEnd.WriteString(value)
					buildEnd.WriteString(" ")
					buildEnd.WriteString("23:59:59")
					queryEndTime = buildEnd.String()

				} else if dateType == "Month" {
					buildStart.WriteString(value)
					buildStart.WriteString("-01 ")
					buildStart.WriteString("00:00:00")
					queryStartTime = buildStart.String()

					monthValue := strings.Split(value, "-")
					buildEnd.WriteString(value)
					if monthValue[1] == "01" || monthValue[1] == "03" || monthValue[1] == "05" || monthValue[1] == "07" || monthValue[1] == "08" || monthValue[1] == "10" || monthValue[1] == "12" {
						buildEnd.WriteString("-31 ")
					} else if monthValue[1] == "02" {
						year, yearerr := strconv.Atoi(monthValue[0])
						if yearerr != nil {
							return getJournal, errmsg.ERROR
						}
						if (year%4 == 0 && year%100 != 0) || year%400 == 0 {
							buildEnd.WriteString("-29 ")
						} else {
							buildEnd.WriteString("-28 ")
						}
					} else {
						buildEnd.WriteString("-30 ")
					}

					buildEnd.WriteString("23:59:59")
					queryEndTime = buildEnd.String()
				}
			}
		case []interface{}:
			if k == "dateRange" {
				var startTime string
				var endTime string
				for i, u := range value {
					if i == 0 {
						startTime = u.(string)
					} else if i == 1 {
						endTime = u.(string)
					}
				}
				if dateType == "Weekly" {
					buildStart.WriteString(startTime)
					buildStart.WriteString(" ")
					buildStart.WriteString("00:00:00")
					queryStartTime = buildStart.String()

					buildEnd.WriteString(endTime)
					buildEnd.WriteString(" ")
					buildEnd.WriteString("23:59:59")
					queryEndTime = buildEnd.String()
				} else if dateType == "Diy" {
					queryStartTime = startTime
					queryEndTime = endTime
				}
			}
		}
	}

	err := Db.Model(&SystemJournal{}).Where(" time>=? AND time<=?", queryStartTime, queryEndTime).Select("*").Order("time desc ").Limit(1000000).Find(&getJournal).Error
	if err != nil {
		return getJournal, errmsg.ERROR_DATABASE
	}
	return getJournal, errmsg.SUCCSECODE
}

// 日志写入
func WriteJournalModel(writeLog SystemJournal) int {

	err := Db.Model(&SystemJournal{}).Create(&writeLog).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCSECODE
}
