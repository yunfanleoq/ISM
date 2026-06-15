import Vue from "vue";
import { uuid } from 'vue-uuid';
import { deepCopy } from "@/assets/libs/utils";
import store from "../index";
/**
 * 执行编辑命令
 * 注意：这里不要用箭头函数，防止this无法调用
 * @param {*} state
 * @param {*} command 命令对象
 */
export const execute = function(state, command) {
  //暂时不做参数校验
  //在这里分发命令--这里暂时先用switch分发，应该用表格分发
  switch (command.op) {
    case "add": {
      let component = command.component;
      component.identifier = uuid.v1();
      if((typeof component.animate!="undefined")&&(typeof component.animate.move=="undefined"))
      {
        component.animate.move = {
          x:{
            deviceSN:"",
            selectVideoType:0,
            isBandDevice:false,
            bandType:1,
            dataID: "",
            dataName: "",
          },
          y:{
            deviceSN:"",
            selectVideoType:0,
            isBandDevice:false,
            bandType:1,
            dataID: "",
            dataName: "",
          },
        }
      }
      component.name = component.type + state.LayerData.components.length;
      component.style.visible = 1;
      component.style.borderWidth = component.style.borderWidth
          ? component.style.borderWidth
          : 0;
      component.style.BorderEdges = component.style.BorderEdges
          ? component.style.BorderEdges
          : 0;
      component.style.opacity = component.style.opacity
          ? component.style.opacity
          : 1;
      component.style.borderStyle = component.style.borderStyle
          ? component.style.borderStyle
          : "solid";
      component.style.borderColor = component.style.borderColor
          ? component.style.borderColor
          : "#ccccff";
      //component.style.fontFamily = "Arial";
      state.LayerData.components.push(component);
    }
      break;
    case "del": {
      let keys = [];
      for (let i = 0; i < state.LayerData.components.length; i++) {
        let identifier = state.LayerData.components[i].identifier;
        if (state.selectedComponentMap[identifier] != undefined) {
          keys.push(i);
        }
      }
      //排序
      keys.sort((a, b) => {
        return a - b;
      });
      //逆向循环删除
      for (let i = keys.length - 1; i >= 0; i--) {
        state.LayerData.components.splice(keys[i], 1);
      }
      break;
    }
    case "Lock": {
      for (let i = 0; i < state.LayerData.components.length; i++) {
        let identifier = state.LayerData.components[i].identifier;
        if (state.selectedComponentMap[identifier] != undefined) {
          state.LayerData.components[i].lock = command.value
          state.selectedComponentMap[identifier]={}
          Vue.set(state.selectedComponentMap, identifier, state.LayerData.components[i]);
        }
      }
      break;
    }
    case "revolve": {
      let keys = [];
      for (let i = 0; i < state.LayerData.components.length; i++) {
        let identifier = state.LayerData.components[i].identifier;
        if (state.selectedComponentMap[identifier] != undefined) {
          keys.push(i);
        }
      }
      //排序
      keys.sort((a, b) => {
        return a - b;
      });
      //逆向循环删除
      for (let i = keys.length - 1; i >= 0; i--) {
        state.LayerData.components[keys[i]].style.transform=state.LayerData.components[keys[i]].style.transform+90
        if(state.LayerData.components[keys[i]].style.transform>=360)
        {
          state.LayerData.components[keys[i]].style.transform=0
        }
      }
      break;
    }
    case "reverse": {
      let keys = [];
      for (let i = 0; i < state.LayerData.components.length; i++) {
        let identifier = state.LayerData.components[i].identifier;
        if (state.selectedComponentMap[identifier] != undefined) {
          keys.push(i);
        }
      }
      //排序
      keys.sort((a, b) => {
        return a - b;
      });
      //逆向循环删除
      for (let i = keys.length - 1; i >= 0; i--) {
        state.LayerData.components[keys[i]].style.transform=state.LayerData.components[keys[i]].style.transform-90
        if(state.LayerData.components[keys[i]].style.transform<=-360)
        {
          state.LayerData.components[keys[i]].style.transform=0
        }
      }
      break;
    }
    case "FlipVertical": {
      let keys = [];
      for (let i = 0; i < state.LayerData.components.length; i++) {
        let identifier = state.LayerData.components[i].identifier;
        if (state.selectedComponentMap[identifier] != undefined) {
          keys.push(i);
        }
      }
      //排序
      keys.sort((a, b) => {
        return a - b;
      });
      //逆向循环删除
      for (let i = keys.length - 1; i >= 0; i--) {
        if(state.LayerData.components[keys[i]].style.transform==-1099)
        {
          state.LayerData.components[keys[i]].style.transform = 0
        }
        else {
          state.LayerData.components[keys[i]].style.transform = -1099
        }
      }
      break;
    }
    case "FlipHorizontally": {
      let keys = [];
      for (let i = 0; i < state.LayerData.components.length; i++) {
        let identifier = state.LayerData.components[i].identifier;
        if (state.selectedComponentMap[identifier] != undefined) {
          keys.push(i);
        }
      }
      //排序
      keys.sort((a, b) => {
        return a - b;
      });
      //逆向循环删除
      for (let i = keys.length - 1; i >= 0; i--) {
        if(state.LayerData.components[keys[i]].style.transform==-1098)
        {
          state.LayerData.components[keys[i]].style.transform=0
        }
        else
        {
          state.LayerData.components[keys[i]].style.transform=-1098
        }

      }
      break;
    }
    case "move": {
      let dx = command.dx,
          dy = command.dy;
      for (let key in command.items) {

        let component = command.items[key];
        if(!component.lock) {
          component.style.position.x = parseInt(component.style.position.x) + parseInt(dx);
          component.style.position.y = parseInt(component.style.position.y) + parseInt(dy);
        }
      }
    }
      break;
    case "adsorption": {
      let keyIndex = command.key
      let keyValue = command.value
      for (let key in command.items) {
        let component = command.items[key];
        if(!component.lock) {
          if(component.groupID!=undefined&&component.groupID!=="")
          {
                //成组就不再对其
          }
          else {
            if (keyIndex == "left") {
             component.style.position.x = parseInt(keyValue);
            } else {
             component.style.position.y = parseInt(keyValue);
            }
          }
        }
      }
    }
      break;
    case "newStyle": {
        console.log("更新样式")
    }
      break;
    case "copy-add": {
      this.commit("ISMDisPlayEditorTool/clearSelectedComponent");
      let te =  command.items
      command.items = []
      let groupID = uuid.v1();
      for (let i = 0; i < te.length; i++) {
        let t =te[i];
        let component = deepCopy(t);
        component.identifier = uuid.v1();
        component.name = "copy_"+component.type + state.LayerData.components.length;
        component.style.visible = 1;
        if((typeof component.groupID!=="undefined")&&(component.groupID!==""))
        {
          component.groupID = groupID
        }
        if((typeof  command.type=="undefined" || command.type!="undo")&&(state.curPageUuid==state.prePageUuid))
        {
          // component.style.position.x =
          //     component.style.position.x + 25 * (state.copyCount + 1);
          // component.style.position.y =
          //     component.style.position.y + 25 * (state.copyCount + 1);
          //
          // component.style.position.x =
          //     component.style.position.x + 50
          // component.style.position.y =
          //     component.style.position.y + 25 * (state.copyCount + 1);
          component.style.position.y =
              parseInt(component.style.position.y) + 10 ;
        }
        // component.style.position.y =
        //     component.style.position.y + 25 ;
        state.LayerData.components.push(component);
        command.items.push(component)
        this.commit("ISMDisPlayEditorTool/addSelectedComponent", component);
        this.commit("ISMDisPlayEditorTool/increaseCopyCount");
      }
    }
      break;
    case "AllSelect": {
      state.selectedComponentMap = {};
      if( command.state==1)
      {
        for (let i = 0; i < state.LayerData.components.length; i++) {
          Vue.set(state.selectedComponentMap, state.LayerData.components[i].identifier, state.LayerData.components[i]);
        }
      }
    }
      break;
    default:
      console.warn("不支持的命令.");
      break;
  }
  //记录操作
  // state.undoStack.push(command);
};


/**
 * 执行历史记录
 * 注意：这里不要用箭头函数，防止this无法调用
 * @param {*} state
 * @param {*} 组件内容
 */
export const executeUndoStack = function(state, data) {

  if((data.components==undefined)||(data.components.length<=0))
  {
    return
  }
  // console.log("data",data.components[0].style.position)
  //记录操作
  if (state.undoStack.length>20){
    state.undoStack.shift()
    state.undoStack.push(data.components);
  }
  else {
    state.undoStack.push(data.components);
  }
};


export const ClearUndo = state => {
  state.undoStack=[]
};

export const ClearRedo = state => {
  state.redoStack=[]
};
export const ClearisHistoryOp = function(state, data) {
  state.isHistoryOp=data
};
export const SetFormatPainterState = function(state, data) {
  state.isFormatPainterActive=data
};
export const SetFormatPainterCell = function(state, data) {
  state.FormatPainterCell=data
};
export const ClearOpFlag = state => {
  state.isHistoryOp=false
};

export const undo = state => {
  state.selectedComponents = [];
  for (let key in state.selectedComponentMap) {
    Vue.delete(state.selectedComponentMap, key);
  }
  Vue.set(state, "selectedComponent", null);
  let components = state.undoStack.pop();
  if (components == undefined) {
    state.isHistoryOp=false
    return;
  }
  state.isHistoryOp=true
  state.LayerData.components = components
  state.redoStack.push(components);
};

export const redo = function(state) {
  let components = state.redoStack.pop();
  if (components == undefined) {
    state.isHistoryOp=false
    return;
  }
  state.isHistoryOp=true
  state.undoStack.push(components);
  state.LayerData.components = components
};
/**
 * 设置 组件列表
 * @param {*} state
 * @param {*} DataList
 */
export const setToolBoxList = (state,DataList) => {
  Vue.set(state, "toolBoxList", DataList);
};

/**
 * 设置 页面装饰列表
 * @param {*} state
 * @param {*} DataList
 */
export const setPageCanVasList = (state,DataList) => {
  Vue.set(state, "PageCanVasList", DataList);
};


/**
 * 清空当前选择的组件
 * @param {*} state
 * @param {*} DataList
 */
export const ClearSelectedComponent = (state) => {
  Vue.set(state, "selectedComponent", null);
};

/**
 * 设置 Mes系统列表
 * @param {*} state
 * @param {*} DataList
 */
export const setMesComponentsList = (state,DataList) => {
  Vue.set(state, "MesComponentsList", DataList);
};

/**
 * 设置 Diy 图库
 * @param {*} state
 * @param {*} DataList
 */
export const setDiyComponentsList = (state,DataList) => {
  Vue.set(state, "DiyComponentsList", DataList);
};

const formatDateTime = function () {
  let dateGet = new Date()
  let y = dateGet.getFullYear();
  let m = dateGet.getMonth() + 1;
  m = m < 10 ? ('0' + m) : m;
  let d = dateGet.getDate();
  d = d < 10 ? ('0' + d) : d;
  let h = dateGet.getHours();
  h=h < 10 ? ('0' + h) : h;
  let minute = dateGet.getMinutes();
  minute = minute < 10 ? ('0' + minute) : minute;
  let second=dateGet.getSeconds();
  second=second < 10 ? ('0' + second) : second;
  return y + '-' + m + '-' + d+' '+h+':'+minute+':'+second;
};

/**
 * 设置 图层的放大或者缩小
 * @param {*} state
 * @param {*} zoomValue
 */
export const setlayerZoom = (state,zoomValue) => {
  Vue.set(state, "selectedValue", zoomValue);
};

/**
 * 设置 日志列表
 * @param {*} state
 * @param {*} loggerList
 */
export const setLoggerList = (state,conent) => {
  const haved_list = store.state.ISMDisPlayEditorTool.loggerList
  const writeLoger = {
    id:state.loggerIndex++,
    level:conent.level,
    content:conent.content,
    time:formatDateTime()
  }
  haved_list.push(writeLoger)
  Vue.set(state, "loggerList", haved_list);
};

function deepMerge(array1, array2) {
  if(!Array.isArray(array1)){
    return []
  }
  if(!Array.isArray(array2))
  {
    return array1
  }
  let isExist = 0
  for(let i = 0;i<array2.length;i++)
  {
    isExist=0
    for(let k=0;k<array1.length;k++)
    {
      if(array2[i].key==array1[k].key)
      {
        isExist=1
        break
      }
    }
    if(isExist==0)
    {
      array1.push(array2[i])
    }
  }
  return array1
}
/**
 * 设置 当前选中的组件-单选
 * @param {*} state
 * @param {*} component
 */
export const setSelectedComponent = (state, component) => {
  if (!component.identifier) {
    component.identifier = uuid.v1();
  }
  for(let i=0;i<state.toolBoxList.length;i++)
  {
    for(let k=0;k<state.toolBoxList[i].items.length;k++)
    {
      if(component.type==state.toolBoxList[i].items[k].info.type)
      {
        const tempDiy = state.toolBoxList[i].items[k].info.style.diy
        // const tempAction = state.toolBoxList[i].items[k].info.style.action
        // const tempDataBind = state.toolBoxList[i].items[k].info.style.dataBind

        const diy = deepMerge( component.style.diy,tempDiy);
        // const action = deepMerge(tempAction, component.action);
        // const dataBind = deepMerge(tempDataBind, component.dataBind);
        component.style.diy = diy
        // component.action = action
        // component.dataBind = dataBind
        state.selectedComponents = [component.identifier];

        if(state.isFormat) {
          for (let i = 0; i < state.LayerData.components.length; i++) {
            if(state.LayerData.components[i].identifier==component.identifier) {
              state.LayerData.components[i].style.position.w = state.formatSrcItems.style.position.w
              state.LayerData.components[i].style.position.h = state.formatSrcItems.style.position.h
              state.LayerData.components[i].style.borderWidth = state.formatSrcItems.style.borderWidth
              state.LayerData.components[i].style.BorderEdges = state.formatSrcItems.style.BorderEdges
              state.LayerData.components[i].style.opacity = state.formatSrcItems.style.opacity
              state.LayerData.components[i].style.borderStyle = state.formatSrcItems.style.borderStyle
              state.LayerData.components[i].style.borderColor = state.formatSrcItems.style.borderColor
              state.LayerData.components[i].style.transform = state.formatSrcItems.style.transform
              state.LayerData.components[i].style.backColor = state.formatSrcItems.style.backColor
              state.LayerData.components[i].style.foreColor = state.formatSrcItems.style.foreColor
              state.LayerData.components[i].style.textAlign = state.formatSrcItems.style.textAlign
              state.LayerData.components[i].style.fontFamily = state.formatSrcItems.style.fontFamily
              state.LayerData.components[i].style.fontWeight = state.formatSrcItems.style.fontWeight
              state.LayerData.components[i].style.italic = state.formatSrcItems.style.italic
              state.LayerData.components[i].style.fontSize = state.formatSrcItems.style.fontSize

              state.isFormat = false
              state.formatSrcItems = {}
              break
            }
          }
        }

        state.selectedComponentMap = {};
        Vue.set(state.selectedComponentMap, component.identifier, component);
        for(let i=0;i<state.LayerData.components.length;i++)
        {
          if((typeof component.groupID!=="undefined")&&(component.groupID!=="")&&(component.groupID==state.LayerData.components[i].groupID))
          {
            Vue.set(state.selectedComponentMap, state.LayerData.components[i].identifier, state.LayerData.components[i]);
          }
        }

        Vue.set(state, "selectedComponent", component);
        return
      }
    }
  }
  for(let i=0;i<state.PageCanVasList.length;i++)
  {
    for(let k=0;k<state.PageCanVasList[i].items.length;k++)
    {
      if(component.type==state.PageCanVasList[i].items[k].info.type)
      {
        const tempDiy = state.PageCanVasList[i].items[k].info.style.diy
        // const tempAction = state.toolBoxList[i].items[k].info.style.action
        // const tempDataBind = state.toolBoxList[i].items[k].info.style.dataBind

        const diy = deepMerge( component.style.diy,tempDiy);
        // const action = deepMerge(tempAction, component.action);
        // const dataBind = deepMerge(tempDataBind, component.dataBind);
        component.style.diy = diy
        // component.action = action
        // component.dataBind = dataBind
        state.selectedComponents = [component.identifier];

        if(state.isFormat) {
          for (let i = 0; i < state.LayerData.components.length; i++) {
            if(state.LayerData.components[i].identifier==component.identifier) {
              state.LayerData.components[i].style.position.w = state.formatSrcItems.style.position.w
              state.LayerData.components[i].style.position.h = state.formatSrcItems.style.position.h
              state.LayerData.components[i].style.borderWidth = state.formatSrcItems.style.borderWidth
              state.LayerData.components[i].style.BorderEdges = state.formatSrcItems.style.BorderEdges
              state.LayerData.components[i].style.opacity = state.formatSrcItems.style.opacity
              state.LayerData.components[i].style.borderStyle = state.formatSrcItems.style.borderStyle
              state.LayerData.components[i].style.borderColor = state.formatSrcItems.style.borderColor
              state.LayerData.components[i].style.transform = state.formatSrcItems.style.transform
              state.LayerData.components[i].style.backColor = state.formatSrcItems.style.backColor
              state.LayerData.components[i].style.foreColor = state.formatSrcItems.style.foreColor
              state.LayerData.components[i].style.textAlign = state.formatSrcItems.style.textAlign
              state.LayerData.components[i].style.fontFamily = state.formatSrcItems.style.fontFamily
              state.LayerData.components[i].style.fontWeight = state.formatSrcItems.style.fontWeight
              state.LayerData.components[i].style.italic = state.formatSrcItems.style.italic
              state.LayerData.components[i].style.fontSize = state.formatSrcItems.style.fontSize
              state.isFormat = false
              state.formatSrcItems = {}
              break
            }
          }
        }

        state.selectedComponentMap = {};
        Vue.set(state.selectedComponentMap, component.identifier, component);
        for(let i=0;i<state.LayerData.components.length;i++)
        {
          if((typeof component.groupID!=="undefined")&&(component.groupID!=="")&&(component.groupID==state.LayerData.components[i].groupID))
          {
            Vue.set(state.selectedComponentMap, state.LayerData.components[i].identifier, state.LayerData.components[i]);
          }
        }
        Vue.set(state, "selectedComponent", component);
        return
      }
    }
  }

  for(let i=0;i<state.MesComponentsList.length;i++)
  {
    for(let k=0;k<state.MesComponentsList[i].items.length;k++)
    {
      if(component.type==state.MesComponentsList[i].items[k].info.type)
      {
        const tempDiy = state.MesComponentsList[i].items[k].info.style.diy
        // const tempAction = state.toolBoxList[i].items[k].info.style.action
        // const tempDataBind = state.toolBoxList[i].items[k].info.style.dataBind

        const diy = deepMerge( component.style.diy,tempDiy);
        // const action = deepMerge(tempAction, component.action);
        // const dataBind = deepMerge(tempDataBind, component.dataBind);
        component.style.diy = diy
        // component.action = action
        // component.dataBind = dataBind
        state.selectedComponents = [component.identifier];

        if(state.isFormat) {
          for (let i = 0; i < state.LayerData.components.length; i++) {
            if(state.LayerData.components[i].identifier==component.identifier) {
              state.LayerData.components[i].style.position.w = state.formatSrcItems.style.position.w
              state.LayerData.components[i].style.position.h = state.formatSrcItems.style.position.h
              state.LayerData.components[i].style.borderWidth = state.formatSrcItems.style.borderWidth
              state.LayerData.components[i].style.BorderEdges = state.formatSrcItems.style.BorderEdges
              state.LayerData.components[i].style.opacity = state.formatSrcItems.style.opacity
              state.LayerData.components[i].style.borderStyle = state.formatSrcItems.style.borderStyle
              state.LayerData.components[i].style.borderColor = state.formatSrcItems.style.borderColor
              state.LayerData.components[i].style.transform = state.formatSrcItems.style.transform
              state.LayerData.components[i].style.backColor = state.formatSrcItems.style.backColor
              state.LayerData.components[i].style.foreColor = state.formatSrcItems.style.foreColor
              state.LayerData.components[i].style.textAlign = state.formatSrcItems.style.textAlign
              state.LayerData.components[i].style.fontFamily = state.formatSrcItems.style.fontFamily
              state.LayerData.components[i].style.fontWeight = state.formatSrcItems.style.fontWeight
              state.LayerData.components[i].style.italic = state.formatSrcItems.style.italic
              state.LayerData.components[i].style.fontSize = state.formatSrcItems.style.fontSize
              state.isFormat = false
              state.formatSrcItems = {}
              break
            }
          }
        }

        state.selectedComponentMap = {};
        Vue.set(state.selectedComponentMap, component.identifier, component);
        for(let i=0;i<state.LayerData.components.length;i++)
        {
          if((typeof component.groupID!=="undefined")&&(component.groupID!=="")&&(component.groupID==state.LayerData.components[i].groupID))
          {
            Vue.set(state.selectedComponentMap, state.LayerData.components[i].identifier, state.LayerData.components[i]);
          }
        }
        Vue.set(state, "selectedComponent", component);
        return
      }
    }
  }
};

/**
 * 增加选中的组件--多选模式
 * @param {*} state
 * @param {*} component
 */
export const addSelectedComponent = (state, component) => {
  if (!component.identifier) {
    component.identifier = uuid.v1();
  }
  if (state.selectedComponentMap[component.identifier]) {
    return;
  }
  if(state.isFormat) {
    for (let i = 0; i < state.LayerData.components.length; i++) {
      if(state.LayerData.components[i].identifier==component.identifier) {
        state.LayerData.components[i].style.position.w = state.formatSrcItems.style.position.w
        state.LayerData.components[i].style.position.h = state.formatSrcItems.style.position.h
        state.LayerData.components[i].style.borderWidth = state.formatSrcItems.style.borderWidth
        state.LayerData.components[i].style.BorderEdges = state.formatSrcItems.style.BorderEdges
        state.LayerData.components[i].style.opacity = state.formatSrcItems.style.opacity
        state.LayerData.components[i].style.borderStyle = state.formatSrcItems.style.borderStyle
        state.LayerData.components[i].style.borderColor = state.formatSrcItems.style.borderColor
        state.LayerData.components[i].style.transform = state.formatSrcItems.style.transform
        state.LayerData.components[i].style.backColor = state.formatSrcItems.style.backColor
        state.LayerData.components[i].style.foreColor = state.formatSrcItems.style.foreColor
        state.LayerData.components[i].style.textAlign = state.formatSrcItems.style.textAlign
        state.LayerData.components[i].style.fontFamily = state.formatSrcItems.style.fontFamily
        state.LayerData.components[i].style.fontWeight = state.formatSrcItems.style.fontWeight
        state.LayerData.components[i].style.italic = state.formatSrcItems.style.italic
        state.LayerData.components[i].style.fontSize = state.formatSrcItems.style.fontSize
        break
      }
    }
  }
  state.selectedComponents.push(component.identifier);
  Vue.set(state.selectedComponentMap, component.identifier, component);
  if (state.selectedComponents.length == 1) {
    Vue.set(state, "selectedComponent", component);
  } else {
    Vue.set(state, "selectedComponent", null);
  }
};

/**
 * 将一个组件从已选中当中移除
 * @param {*} state
 * @param {*} component
 */
export const removeSelectedComponent = (state, component) => {
  if (!component.identifier) return;
  let index = -1;
  for (let i = 0; i < state.selectedComponents.length; i++) {
    if (state.selectedComponents[i] == component.identifier) {
      index = i;
      break;
    }
  }
  if (index > -1) {
    state.selectedComponents.splice(index, 1);
  }
  Vue.delete(state.selectedComponentMap, component.identifier);
  //如果移除的是选中组件
  if (
    state.selectedComponent != null &&
    component.identifier == state.selectedComponent.identifier
  ) {
    Vue.set(state, "selectedComponent", null);
  }
  //如果只有一个组件，则默认选中
  if (state.selectedComponents.length == 1) {
    let _component = state.selectedComponentMap[state.selectedComponents[0]];
    Vue.set(state, "selectedComponent", _component);
  }
};

/**
 * 清理所有选中的组件
 * @param {*} state
 */
export const clearSelectedComponent = state => {
  state.selectedComponents = [];
  for (let key in state.selectedComponentMap) {
    Vue.delete(state.selectedComponentMap, key);
  }
  Vue.set(state, "selectedComponent", null);
};

export const setLayerSelected = (state, selected) => {
  state.selectedIsLayer = selected;
};

export const clearSelectComponent = (state) => {
  state.selectedComponentMap = {};
};

export const setCopySrcItems = (state, items) => {
  state.copySrcItems = items;
  state.copyCount = 0;
};
export const setFormatSrcItems = (state, items) => {
  if(JSON.stringify(items)=="{}")
  {
    state.isFormat = false
  }
  else
  {
    state.isFormat = true
  }
  state.formatSrcItems = items;

};
export const increaseCopyCount = state => {
  state.copyCount++;
};
export const setCurrentISMCavasContainer = (state, CavasContainer) => {
  state.ISMCavasContainer=CavasContainer
};
export const setCurrentISMCavasDND = (state, ISMCavasDND) => {
  state.ISMCavasDND=ISMCavasDND
};
export const setCurrentSelectNode = (state, Node) => {
  state.selectedIsLayer = false
  if(state.isFormatPainterActive)
  {
    return
  }
  const NodeData = Node.prop().data
  if(typeof NodeData=="undefined"||typeof NodeData.detail=="undefined")
  {
    Node.setProp({
      position:{x:12, y:13},
      size:{width:100,height:100},
      type:"connect-line",
      attrs: {
        line: {
          connection: true,
          style: {
            stroke: "#13c2c2",
            strokeWidth: 10,
          },
          sourceMarker: "classic",
          targetMarker: "classic",
        },
        wrap: {
          connection: true,
          stroke: "none",
          strokeWidth: 25,
          strokeLinejoin: 'round',
        }
      },
      data: {
        locked:false,
        UpdateNodeFlag:true,
        editMode: true,
        showDeviceUuid:"",
        IsToolBox:false,
        detail:{
          identifier :uuid.v1(),
          name:"线段",
          "type": "image",
          isCanvas:true,
          "action": [],
          "dataBind":[],
          "active": [
            {
              id:"Forward",
              name:"component.ViewCanvasMoveLineArrow.Forward",
              result:"",
              isExpression:true,
              condition:{
                deviceSN:"",
                selectVideoType:0,
                isBandDevice:false,
                bandType:1,
                dataID: "",
                dataName: "",
                operator:"",
                OperatorValue:"",
                OperatorMaxValue:"",
              },
            },
            {
              id:"Reverse",
              name:"component.ViewCanvasMoveLineArrow.Reverse",
              result:"",
              isExpression:true,
              condition:{
                deviceSN:"",
                selectVideoType:0,
                isBandDevice:false,
                bandType:1,
                dataID: "",
                dataName: "",
                operator:"",
                OperatorValue:"",
                OperatorMaxValue:"",
              },
            },
          ],
          "animate": {
            "selected": [],
            "condition":{
              deviceSN:"",
              selectVideoType:0,
              isBandDevice:false,
              bandType:1,
              dataID: "",
              dataName: "",
              operator:"",
              OperatorValue:"",
              OperatorMaxValue:"",
            },
            "isExpression": false,
            "animateList": [

            ],
            "animateElement": [
              {
                id: "blink",
                elementList:[
                  {
                    "name":"component.public.animateSpeed",
                    "type":7,
                    "value":1,
                    "min":0.1,
                    "key":"blinkSpeed",
                  },
                ]
              },
              {
                id: "millcolorGrad",
                elementList:[
                  {
                    "name": "component.public.startColor",
                    "type": 2,
                    "value": "#74f808",
                    "key": "startColor",
                  },
                  {
                    "name": "component.public.stopColor",
                    "type": 2,
                    "value": "#f30b0b",
                    "key": "stopColor",
                  },
                  {
                    "name":"component.public.animateSpeed",
                    "type":7,
                    "value":1,
                    "min":0.1,
                    "key":"animateSpeed",
                  },
                ]
              },
              {
                id: "animateSpin",
                elementList:[
                  {
                    "name":"component.public.animateSpinSpeed",
                    "type":7,
                    "value":1,
                    "min":0.1,
                    "key":"spinSpeed",
                  },
                  {
                    name:"configComponent.bigScreen.border.border89Direction",
                    type:6,
                    value:0,
                    enumList:[
                      {
                        value:0,
                        option:"configComponent.bigScreen.border.border89DirectionForward"
                      },
                      {
                        value:1,
                        option:"configComponent.bigScreen.border.border89DirectionNegative"
                      }
                    ],
                    min:1,
                    key:"spinDirection",
                  }
                ]
              },
            ],
          },
          "style": {
            "position": {
              "x": 0,
              "y": 0,
              "w": 600,
              "h": 50
            },
            "points": [],
            "visible":1,
            "zIndex": -1,
            "transform": 0,
            "backColor": "#D9D9D9",
            foreColor:"#13c2c2",
            borderWidth:0,
            BorderEdges:0,
            opacity:1,
            borderStyle:"solid",
            borderColor:"#ccccff",
            "diy":[
              {
                "name":"component.public.strokeLinejoin",
                type:6,
                value:1,
                enumList:[
                  {
                    value:0,
                    option:"component.public.strokeLineMiter"
                  },
                  {
                    value:1,
                    option:"component.public.strokeLineRound"
                  },
                  {
                    value:2,
                    option:"component.public.strokeLineBevel"
                  }
                ],
                "key":"strokeLinejoin",
              },
              {
                "name":"component.public.strokeBgWidth",
                "type":1,
                "value":25,
                "min":1,
                "key":"strokeBgWidth",
              },
              {
                "name":"component.public.strokeWidth",
                "type":1,
                "value":10,
                "min":1,
                "key":"strokeWidth",
              },
              {
                "name":"component.public.strokeLength",
                "type":1,
                "value":15,
                "min":1,
                "key":"strokeLength",
              },
              {
                "name":"component.public.strokeSpace",
                "type":1,
                "value":10,
                "min":1,
                "key":"strokeSpace",
              },
              {
                "name":"component.public.animateSpeed",
                "type":1,
                "value":30,
                "min":1,
                "key":"MoveBrokenLineInterval",
              },
              {
                name:"displayConfig.ToolBox.Diagram.MoveBrokenLineConditionEnable",
                type:6,
                value:0,
                enumList:[
                  {
                    value:0,
                    option:"component.public.Forbidden"
                  },
                  {
                    value:1,
                    option:"component.public.Enable"
                  }
                ],
                min:1,
                key:"MoveBrokenLineConditionEnable",
              },
              {
                name:"configComponent.bigScreen.border.border89Direction",
                type:6,
                value:0,
                enumList:[
                  {
                    value:0,
                    option:"configComponent.bigScreen.border.border89DirectionForward"
                  },
                  {
                    value:1,
                    option:"configComponent.bigScreen.border.border89DirectionNegative"
                  }
                ],
                min:1,
                key:"spinDirection",
              }
            ]
          }
        }
      },
    })
  }

  Vue.set(state, "selectedComponent", Node.prop().data.detail);
  Vue.set(state, "selectedNodePops", Node.prop());
  Vue.set(state, "selectedNode", Node);
};

export const ClickUnSelectedComponent = (state, selected) => {
  Vue.set(state, "UnSelectedComponent", selected);
};
// 触发锁定
export const LOCK_SCREEN = (state) =>{
  Vue.set(state, "isLocked", true);
}
// 触发解锁
export const  UNLOCK_SCREEN = (state) => {
  Vue.set(state, "isLocked", false);
}
// 更新真实密码（如登录后同步）
export const  UPDATE_REAL_PASSWORD = (state, password) =>{
  Vue.set(state, "realPassword", false);
}