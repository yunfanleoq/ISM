
import $ from 'jquery'
/**
 * 触发器添加
 */
export async function ComponentRestApi(type,url,params) {
    switch (type){
        case "Get":{
           return  $.ajax({
                url:url,
                dataType:"json",
                data:params,
                type:"get",
                success:function(resp){
                    return resp
                }
            })
        }
        case "Post":{
            return  $.ajax({
                url:url,
                dataType:"json",
                data:params,
                type:"post",
                success:function(resp){
                    return resp
                }
            })
        }
    }
}

export default {
    ComponentRestApi,
}