<template>
  <div id="app">

  </div>
</template>

<script>

export default {
  name: "AmisEditor",
  components: {

  },
  data() {
    return {
      preview: true,
      schema: {
        "type": "page",
        "title": "浏览器内核对 CSS 的支持情况",
        "body": [
          {
            "type": "crud",
            "draggable": true,
            "syncLocation": false,
            "api": "http://127.0.0.1:18089/RestApi",
            "keepItemSelectionOnPageChange": true,
            "autoGenerateFilter": true,
            "bulkActions": [
              {
                "type": "button",
                "label": "批量删除",
                "actionType": "ajax",
                "api": "delete:https://3xsw4ap8wah59.cfc-execute.bj.baidubce.com/api/amis-mock/mock2/sample/${ids|raw}",
                "confirmText": "确定要批量删除?",
                "id": "u:cdd3f6a89ad7"
              },
              {
                "type": "button",
                "label": "批量修改",
                "actionType": "dialog",
                "dialog": {
                  "title": "批量编辑",
                  "name": "sample-bulk-edit",
                  "body": {
                    "type": "form",
                    "api": "https://3xsw4ap8wah59.cfc-execute.bj.baidubce.com/api/amis-mock/mock2/sample/bulkUpdate2",
                    "body": [
                      {
                        "type": "hidden",
                        "name": "ids"
                      },
                      {
                        "type": "input-text",
                        "name": "engine",
                        "label": "Engine"
                      }
                    ]
                  }
                },
                "id": "u:1d09bed8de6e"
              }
            ],
            "quickSaveApi": "https://3xsw4ap8wah59.cfc-execute.bj.baidubce.com/api/amis-mock/mock2/sample/bulkUpdate",
            "quickSaveItemApi": "https://3xsw4ap8wah59.cfc-execute.bj.baidubce.com/api/amis-mock/mock2/sample/$id",
            "headerToolbar": [
              "bulkActions",
              {
                "type": "button",
                "label": "重置测试数据",
                "actionType": "ajax",
                "size": "sm",
                "api": "https://3xsw4ap8wah59.cfc-execute.bj.baidubce.com/api/amis-mock/mock2/sample/reset",
                "id": "u:3b4250586b32"
              },
              "export-excel",
              {
                "type": "tpl",
                "tpl": "一共有 ${count} 行数据。",
                "className": "v-middle",
                "id": "u:c13c4d1b4ee5"
              },
              {
                "type": "columns-toggler",
                "align": "right",
                "draggable": true
              },
              {
                "type": "drag-toggler",
                "align": "right"
              }
            ],
            "footerToolbar": [
              "statistics",
              "switch-per-page",
              "pagination"
            ],
            "columns": [
              {
                "name": "Id",
                "label": "ID",
                "width": 20,
                "sortable": true,
                "type": "text",
                "searchable": {
                  "type": "input-text",
                  "name": "id",
                  "label": "主键",
                  "placeholder": "输入id",
                  "id": "u:f5d4e8272960"
                },
                "id": "u:757100a0d226"
              },
              {
                "name": "Browser",
                "label": "Browser",
                "searchable": {
                  "type": "select",
                  "name": "browser",
                  "label": "浏览器",
                  "placeholder": "选择浏览器",
                  "options": [
                    {
                      "label": "Internet Explorer ",
                      "value": "ie"
                    },
                    {
                      "label": "AOL browser",
                      "value": "aol"
                    },
                    {
                      "label": "Firefox",
                      "value": "firefox"
                    }
                  ],
                  "id": "u:6cab195138de"
                },
                "id": "u:8672b62de77b"
              },
              {
                "name": "Platform",
                "label": "平台",
                "popOver": {
                  "trigger": "hover",
                  "body": {
                    "type": "tpl",
                    "tpl": "就是为了演示有个叫 popOver 的功能"
                  }
                },
                "sortable": true,
                "type": "text",
                "id": "u:5b2b605e5614"
              },
              {
                "name": "Grade",
                "label": "CSS 等级",
                "type": "select",
                "options": [
                  "A",
                  "B",
                  "C",
                  "D",
                  "X"
                ],
                "id": "u:abd42b710ce9"
              },
              {
                "type": "Operation",
                "label": "操作",
                "width": 100,
                "buttons": [
                  {
                    "type": "button",
                    "actionType": "ajax",
                    "label": "删除",
                    "confirmText": "您确认要删除?",
                    "api": "delete:https://3xsw4ap8wah59.cfc-execute.bj.baidubce.com/api/amis-mock/mock2/sample/$id",
                    "id": "u:536bacccc8a7"
                  }
                ],
                "id": "u:e7ced0ae0197"
              }
            ],
            "id": "u:be141af79f8e"
          }
        ],
        "remark": "嘿，不保证数据准确性",
        "id": "u:b02ec829e806"
      },
    }
  },
}
</script>

<style scoped>

</style>
