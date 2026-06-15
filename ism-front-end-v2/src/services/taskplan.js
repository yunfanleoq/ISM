import {
    TASKPLANAADD,
    TASKPLANDEL,
    TASKPLANEDIT,
    TASKPLANLIST,
} from '@/services/api'
import {request, METHOD} from '@/utils/request'

export async function AddTaskPlan(params) {
    return request(TASKPLANAADD, METHOD.POST,params)
}
export async function DelTaskPlan(params) {
    return request(TASKPLANDEL, METHOD.POST,params)
}
export async function EditTaskPlan(params) {
    return request(TASKPLANEDIT, METHOD.POST,params)
}
export async function GetTaskPlanList() {
    return request(TASKPLANLIST, METHOD.POST)
}

export default {
    AddTaskPlan,
    DelTaskPlan,
    EditTaskPlan,
    GetTaskPlanList
}
