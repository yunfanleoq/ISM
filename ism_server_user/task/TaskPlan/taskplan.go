/**
 * @ Author: ISM Web组态软件
 * @ Create Time: 2023-01-09 08:53:36
 * @ Modified by: ISM Web组态软件
 * @ Modified time: 2023-04-03 09:02:34
 * @ Description: 此源码版权归 www.ismctl.com 所有,个人私自不得二次销售。
 */

package taskplanpthread

import (
	"ISMServer/models"
	"ISMServer/utils/errmsg"

	"github.com/jakecoffman/cron"
)

var TaskPlanCron *cron.Cron

func getAllTaskPlan() (int, []models.TaskPlanList) {
	var GetTask []models.TaskPlanList
	err := models.Db.Unscoped().Model(&models.TaskPlanList{}).Where("ID >= 0").Find(&GetTask).Error
	if err != nil {
		return errmsg.ERROR, GetTask
	}
	return errmsg.SUCCSE, GetTask
}
func TaskPlanPthread() {
	TaskPlanCron = cron.New()
	TaskPlanCron.Start()

	code, taskList := getAllTaskPlan()
	if code == errmsg.SUCCSE && (len(taskList) > 0) {
		for _, task := range taskList {
			t := TaskJobPthread{Task: task}
			TaskPlanCron.AddJob(t.Task.CronExpression, &t, t.Task.TaskUuid)
		}
	}
}
