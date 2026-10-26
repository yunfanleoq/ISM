<template>
  <div style="padding: 5px;height: 85vh;overflow-y: scroll;" class="TreeBox">
    <a-spin :spinning="treeLoading" tip="Loading...">
    <a-input-search style="margin-bottom: 8px" placeholder="Search" @change="onTreeChange" />
    <a-directory-tree
        :tree-data="treeData"
        :load-data="treeLazy ? onLoadTreeData : null"
        :expanded-keys="expandedKeys"
        :selected-keys="selectedKeys"
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
    </a-spin>
  </div>

</template>
<script>

import {getMonitorTree} from "../../services/device";
import {sortMonitorTreeByName} from "@/utils/naturalSort";

export default {
  name: 'deviceTree',
  data() {
    return {
      selection: null,
      defaultSelectKey:[],
      selectNode:null,
      selectKey:null,
      selectedKeys: [],
      expandedKeys: [],
      searchValue: '',
      dataList: [],
      findResult:false,
      autoExpandParent: true,
      treeData:[],
      treeLoading: false,
      treeLazy: true,
    };
  },
  watch: {
    // keep-alive 多页签下，$route 每次切菜单都会触发；若整树重载会丢掉懒加载子节点，
    // 回来只剩 RootZone，表现为「请选择区域」。仅在树为空时拉取。
    '$route'(to) {
      if (!to || !to.path || to.path.indexOf('DataWarehouse') === -1) {
        return
      }
      if (!this.treeData || this.treeData.length === 0) {
        this.getMonitorTree()
      }
    }
  },
  mounted(){
    this.getMonitorTree()
  },
  activated() {
    if (!this.treeData || this.treeData.length === 0) {
      this.getMonitorTree()
    }
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
      _t.treeLoading = true
      const params = _t.treeLazy ? {lazy: true, pid: 0} : {}
      getMonitorTree(params).then(function (res){
        if(res.data.code==0)
        {
          const list = res.data.list || []
          _t.treeData = _t.normalizeTreeNodes(list)
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
        else
        {
          _t.treeData = []
        }
      }).catch(function () {
        _t.treeData = []
      }).finally(function () {
        _t.treeLoading = false
      })
    },
    // 编辑/变更后刷新树：保留 expandedKeys 与选中态，懒加载路径逐级恢复，避免只剩 RootZone
    reloadKeepExpand() {
      const _t = this
      const savedExpanded = (_t.expandedKeys || []).slice()
      const savedSelectKey = _t.selectKey
      const savedSelectedKeys = (_t.selectedKeys || []).slice()
      _t.treeLoading = true
      const params = _t.treeLazy ? {lazy: true, pid: 0} : {}
      return getMonitorTree(params).then(async function (res) {
        if (res.data.code != 0) {
          _t.treeData = []
          return
        }
        _t.treeData = _t.normalizeTreeNodes(res.data.list || [])
        _t.dataList = []
        _t.generateList(_t.treeData)

        if (_t.treeLazy && savedExpanded.length) {
          for (let i = 0; i < savedExpanded.length; i++) {
            await _t.loadChildrenByKey(savedExpanded[i])
          }
          _t.expandedKeys = savedExpanded
          _t.autoExpandParent = false
        }

        if (savedSelectKey != null) {
          _t.selectKey = savedSelectKey
          _t.selectedKeys = savedSelectedKeys.length ? savedSelectedKeys : [savedSelectKey]
          _t.dataSource = []
          const selectNode = _t.findNodeByKey(savedSelectKey, _t.treeData)
          if (selectNode && selectNode.value && selectNode.value.type === 0) {
            if (!selectNode.children) {
              await _t.loadChildrenByKey(savedSelectKey)
            }
            _t.getTreeChildren(savedSelectKey, _t.treeData)
          }
          _t.$emit('updateTree', _t.dataSource)
        }
        if (_t.treeData.length === 0) {
          _t.selectKey = null
          _t.selectNode = null
        }
      }).catch(function () {
        _t.treeData = []
      }).finally(function () {
        _t.treeLoading = false
      })
    },
    findNodeByKey(key, tree) {
      if (!tree || !tree.length) return null
      for (let i = 0; i < tree.length; i++) {
        const node = tree[i]
        if (node.key === key) return node
        if (node.children && node.children.length) {
          const found = this.findNodeByKey(key, node.children)
          if (found) return found
        }
      }
      return null
    },
    loadChildrenByKey(key) {
      const _t = this
      const node = _t.findNodeByKey(key, _t.treeData)
      if (!node || !node.value || node.value.type === 1) {
        return Promise.resolve()
      }
      const parentSid = node.value.sid
      return getMonitorTree({lazy: true, pid: parentSid}).then(function (res) {
        if (res.data.code == 0) {
          node.children = _t.normalizeTreeNodes(res.data.list || [])
          _t.treeData = [..._t.treeData]
          _t.dataList = []
          _t.generateList(_t.treeData)
        }
      }).catch(function () { /* keep previous children */ })
    },
    normalizeTreeNodes(list) {
      if (!list || !list.length) {
        return []
      }
      let nodes = list
      if (this.treeLazy) {
        nodes = list.map(node => {
          const copy = Object.assign({}, node)
          if (copy.value && copy.value.type === 1) {
            copy.isLeaf = true
          } else if (typeof copy.isLeaf === 'boolean') {
            copy.isLeaf = copy.isLeaf
          } else {
            copy.isLeaf = false
          }
          if (!copy.children) {
            copy.children = undefined
          }
          return copy
        })
      }
      return sortMonitorTreeByName(nodes)
    },
    onLoadTreeData(treeNode) {
      const _t = this
      return new Promise(function (resolve) {
        const nodeValue = treeNode.dataRef && treeNode.dataRef.value
        const parentSid = nodeValue ? nodeValue.sid : 0
        getMonitorTree({lazy: true, pid: parentSid}).then(function (res) {
          if (res.data.code == 0) {
            treeNode.dataRef.children = _t.normalizeTreeNodes(res.data.list || [])
            _t.treeData = [..._t.treeData]
          }
          resolve()
        }).catch(function () {
          resolve()
        })
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
      this.selectedKeys = keys && keys.length ? [keys[0]] : []
      if(event.node.value.type==0)
      {
        this.editIsDevice=false
        this.selectKey = keys[0]
        this.dataSource=[]
        this.getTreeChildren(this.selectKey,this.treeData)
      } else {
        this.selectKey = keys[0]
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
