<template>
  <div style="padding: 5px;height: 85vh;overflow-y: scroll;" class="TreeBox">
    <a-input-search style="margin-bottom: 8px" placeholder="Search" @change="onTreeChange" />
    <a-directory-tree
        :tree-data="treeData"
        :expanded-keys="expandedKeys"
        :auto-expand-parent="autoExpandParent"
        :replace-fields="{ value: 'key',title:'text'}"
        @expand="onExpand"
        @select="onSelect">

      <template slot="title" slot-scope="{ title }">
        {{title}}
        <span v-if="title.indexOf(searchValue) > -1">
          {{ title.substr(0, title.indexOf(searchValue)) }}
          <span style="color: #f50">{{ searchValue }}</span>
          {{ title.substr(title.indexOf(searchValue) + searchValue.length) }}
        </span>
        <span v-else>{{ title }}</span>
      </template>
    </a-directory-tree>
  </div>

</template>
<script>

import {getMonitorTree} from "../../services/device";

export default {
  name: 'deviceTree',
  data() {
    return {
      selection: null,
      defaultSelectKey:[],
      selectNode:null,
      selectKey:null,
      expandedKeys: [],
      searchValue: '',
      dataList: [],
      findResult:false,
      autoExpandParent: true,
      treeData:[]
    };
  },
  watch: {
    '$route'() {
      this.getMonitorTree()
    }
  },
  mounted(){
    this.getMonitorTree()
  },
  methods: {
    generateList(data) {
      for (let i = 0; i < data.length; i++) {
        const node = data[i];
        const key = node.key;
        this.dataList.push({ key, title: key });
        if (node.children) {
          this.generateList(node.children);
        }
      }
    },
    getMonitorTree(){
      let _t = this
      getMonitorTree().then(function (res){
        if(res.data.code==0)
        {
          if(res.data.list.length>0)
          {
           // _t.$refs.deviceTree.selectNode(res.data.list[0]);
          }
          _t.treeData =res.data.list
          if(_t.selectKey!=null)
          {
            _t.dataSource=[]
            _t.getTreeChildren(_t.selectKey,_t.treeData)
            _t.$emit("updateTree", _t.dataSource);
          }
          if(_t.treeData.length==0)
          {
            _t.selectKey=null
            _t.selectNode=null
          }
          _t.generateList(_t.treeData);
        }
      })
    },
    checkHavedDevice(key,treeNode){
      for (let i = 0; i < treeNode.length; i++)
      {
        const node = treeNode[i];
        if (node.key==key)
        {
          if((node.children)&&(node.children.length>0))
          {
            this.findResult = true
          }
        }
        else if(node.children)
        {
          this.checkHavedDevice(key, node.children);
        }
      }
      return this.findResult
    },
    checkZoneHavedDevice(key){
      this.findResult = false
      return this.checkHavedDevice(key,this.treeData)
    },
    checkIsEmpty(){
      return this.treeData.length?false:true
    },
    checkChildrenHavedDevice(children) {
      if (children&&children.length>0)
      {
        for (let i = 0; i < children.length; i++)
        {
          const node = children[i]
          if(node.value.type==1)
          {
            return true
          }
          if (node.children&&node.children.length>0)
          {
            this.getChildren(node.children)
          }
        }
      }
    },
    getTreeChildren(key,treeNode){
      for (let i = 0; i < treeNode.length; i++)
      {
        const node = treeNode[i];
        if (node.key==key)
        {
          let temp = {
            key:node.key,
            no:node.value.ID,
            nodeName:node.text,
            nodeType:node.value.type,
            Status:node.value.Status,
            deviceType:node.value.deviceType,
            IsEnable:node.value.IsEnable,
            extra:node.value
          }
          this.dataSource.push(temp)
          this.getChildren(node.children)
        }
        else if(node.children)
        {
          let flag = false
          for(let j=0;j<node.children.length;j++)
          {
            if(node.children[j].key==key)
            {
              flag = true
              let nodeInfo = node.children[j]
              let temp = {
                key:nodeInfo.key,
                no:nodeInfo.value.ID,
                nodeName:nodeInfo.text,
                nodeType:nodeInfo.value.type,
                Status:nodeInfo.value.Status,
                deviceType:node.value.deviceType,
                IsEnable:node.value.IsEnable,
                extra:nodeInfo.value
              }
              this.dataSource.push(temp)
              this.getChildren(node.children[j].children)
              break
            }
          }
          if(!flag)
          {
            this.getTreeChildren(key, node.children);
          }
        }
      }
    },
    getChildren(children) {
      if (children&&children.length>0)
      {
        for (let i = 0; i < children.length; i++)
        {
          const node = children[i]
          let temp = {
            key:node.key,
            no:node.value.ID,
            nodeName:node.text,
            Status:node.value.Status,
            nodeType:node.value.type,
            deviceType:node.value.deviceType,
            IsEnable:node.value.IsEnable,
            extra:node.value
          }
          this.dataSource.push(temp)
          if (node.children&&node.children.length>0)
          {
            this.getChildren(node.children)
          }
        }
      }
    },
    onSelect(keys,event) {
      if(event.node.value.type==0)
      {
        this.editIsDevice=false
        this.selectKey = keys[0]
        this.dataSource=[]
        this.getTreeChildren(this.selectKey,this.treeData)
      }

      const onSelectData = {
        key:keys[0],
        info:event.node,
        tableList:this.dataSource
      }
      this.$emit("onSelect", onSelectData);
    },
    onExpand(expandedKeys) {
      this.expandedKeys = expandedKeys;
      this.autoExpandParent = false;
    },
    getParentKey(key, tree) {
      let parentKey;
      for (let i = 0; i < tree.length; i++) {
        const node = tree[i];
        if (node.children) {
          if (node.children.some(item => item.key === key)) {
            parentKey = node.key;
          } else if (this.getParentKey(key, node.children)) {
            parentKey = this.getParentKey(key, node.children);
          }
        }
      }
      return parentKey;
    },
    onTreeChange(e) {
      let _t = this
      const value = e.target.value;
      const expandedKeys = _t.dataList
          .map(item => {
            if (item.title.indexOf(value) > -1) {
              return _t.getParentKey(item.key, _t.treeData);
            }
            return null;
          })
          .filter((item, i, self) => item && self.indexOf(item) === i);
      Object.assign(_t, {
        expandedKeys,
        searchValue: value,
        autoExpandParent: true,
      });
    },
    onChange(e) {
      const value = e.target.value;
      let _t = this
      const tempExpandedKeys = this.dataList
          .map(item => {
            if (item.title.indexOf(value) > -1) {
              return _t.getParentKey(item.key, _t.treeData);
            }
            return null;
          })
          .filter((item, i, self) => item && self.indexOf(item) === i);
      this.searchValue = value
      this.autoExpandParent = true
      this.expandedKeys = tempExpandedKeys
    },
  },
}
</script>
<style scoped>
::v-deep .tree-node-selected {
  background: #13c2c2;
  color: #fff;
}
::v-deep .tree-title {
  font-size: 14px;
  display: inline-block;
  text-decoration: none;
  vertical-align: middle;
  white-space: nowrap;
  padding: 0 2px;
  margin: 4px 0;
  height: 26px;
  line-height: 26px;
  font-family: Helvetica Neue,Helvetica,PingFang SC,Hiragino Sans GB,Microsoft YaHei,Arial,sans-serif;
}

::v-deep .textbox {
  border:1px solid #d9d9d9;
}
.TreeBox::-webkit-scrollbar {/*滚动条整体样式*/
  width:4px;/*高宽分别对应横竖滚动条的尺寸*/
  height:4px;
}

.TreeBox::-webkit-scrollbar-thumb {/*滚动条里面小方块*/
  /*滚动条里面小方块*/
  border-radius   : 10px;
  background-color: skyblue;
  background-image: -webkit-linear-gradient(
      45deg,
      rgba(255, 255, 255, 0.2) 25%,
      transparent 25%,
      transparent 50%,
      rgba(255, 255, 255, 0.2) 50%,
      rgba(255, 255, 255, 0.2) 75%,
      transparent 75%,
      transparent
  );
}
::v-deep .f-row {
  display: block;
}
.TreeBox::-webkit-scrollbar-track {/*滚动条里面轨道*/
  /*滚动条里面轨道*/
  box-shadow   : inset 0 0 5px rgba(0, 0, 0, 0.2);
  background   : #ededed;
  border-radius: 10px;
}
</style>
