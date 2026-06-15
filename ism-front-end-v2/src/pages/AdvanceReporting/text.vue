<!--// 节点支持自适应 -->
<template>
  <div>
    <div class="tool">
      <a-icon type="save" @click="saveEvent" />
    </div>
    <a-layout>
      <a-layout-sider :theme="'light'">
        <div id="stencil">
          <div>
            <div
                class="dnd-circle dnd-start"
                @mousedown="startDrag('start', $event)"
            ></div>
            <span>开始</span>
          </div>
          <div>
            <div class="dnd-rect" @mousedown="startDrag('rect', $event)"></div>
            <span>节点1</span>
          </div>
          <div>
            <div
                class="dnd-polygon"
                @mousedown="startDrag('polygon', $event)"
            ></div>
            <span>节点2</span>
          </div>
          <div>
            <div class="dnd-circle" @mousedown="startDrag('end', $event)"></div>
            <span>结束</span>
          </div>
        </div>
      </a-layout-sider>
      <a-layout-content>
        <div ref="graphContainer"></div>
      </a-layout-content>
    </a-layout>
    <a-drawer
        title="节点属性编辑"
        :visible="drawer"
        :closable="false"
        @close="handleClose"
    >
      <a-form :data="editNode">
        <a-form-item label="节点名称">
          <a-input v-model="editNode.label"></a-input>
        </a-form-item>
        <a-button type="primary" @click="saveNode">保存</a-button>
      </a-form>
    </a-drawer>
  </div>
</template>
<script>
import { Graph, Shape } from "@antv/x6";
import "@antv/x6-vue-shape";
// 插件 键盘监听事件
import { Keyboard } from "@antv/x6";
// 拖拽事件
import { Dnd } from "@antv/x6";

// import {Stencil} from '@antv/x6-plugin-stencil'
import { Transform } from "@antv/x6";
import { Selection } from "@antv/x6";
import { Snapline } from "@antv/x6";
import { Clipboard } from "@antv/x6";
import { History } from "@antv/x6";
// import {register} from '@antv/x6-vue-shape'
// import insertCss from 'insert-css'
// import {nodeData} from '@/config/liucheng/nodeConfig';

export default {
  name: "MindMap",

  data() {
    return {
      graphOut: {},
      drawer: false,
      currentNode: {},
      editNode: {},
      dnd: {},
      // 节点形状
      shapeList: [
        {
          label: "矩形",
          value: "rect",
        },
        {
          label: "圆形",
          value: "circle",
        },
        {
          label: "椭圆",
          value: "ellipse",
        },
        {
          label: "多边形",
          value: "polygon",
        },
        {
          label: "折线",
          value: "polyline",
        },
        {
          label: "路径",
          value: "path",
        },
        {
          label: "图片",
          value: "image",
        },
      ],
      // 连接桩
      ports: {
        groups: {
          top: {
            position: "top",
            attrs: {
              circle: {
                magnet: true,
                stroke: "black",
                r: 4,
              },
            },
          },
          bottom: {
            position: "bottom",
            attrs: {
              circle: {
                magnet: true,
                stroke: "black",
                r: 4,
              },
            },
          },
          left: {
            position: "left",
            attrs: {
              circle: {
                magnet: true,
                stroke: "black",
                r: 4,
              },
            },
          },
          right: {
            position: "right",
            attrs: {
              circle: {
                magnet: true,
                stroke: "black",
                r: 4,
              },
            },
          },
        },
        items: [
          {
            id: "port_1",
            group: "bottom",
          },
          {
            id: "port_2",
            group: "top",
          },
          {
            id: "port_3",
            group: "left",
          },
          {
            id: "port_4",
            group: "right",
          },
        ],
      },
    };
  },

  mounted() {
    this.graphOut = this.initGraph();
    this.initData();
    this.graphOut.on('node:change:size', ({ node, options }) => {

      const { width, height } = node.store.data.size;
      node.attr('text/textWrap/width', width);
      node.attr('text/textWrap/height', height);
      // 确保文字始终居中于节点
      node.attr('text/refX', width / 2);
      node.attr('text/refY', height / 2);
    })

    this.graphOut.on('node:change:attrs', ({ node, options }) => {
      const { width, height } = node.store.data.size;
      node.attr('text/textWrap/width', width);
      node.attr('text/textWrap/height', height);
      // 确保文字始终居中于节点
      node.attr('text/refX', width / 2);
      node.attr('text/refY', height / 2);
    })
  },
  methods: {
    initGraph() {
      const graph = new Graph({
        container: this.$refs.graphContainer,
        // autoResize: true, // 大小自适应
        height: 700,
        width: "100%",
        grid: true,
        magnetThreshold: "onleave",
        panning: {
          enabled: true,
          modifiers: "shift",
          magnetThreshold: 1,
          // 鼠标画布移动
          eventTypes: ["leftMouseDown"],
        },
        // 开启自动吸附
        // connecting: {
        //   // 距离节点或者连接桩 50 px 触发自动吸附
        //   snap: true,
        //   // 是否允许连接到画布空白位置的点
        //   allowBlank: false,
        //   // 是否允许创建循环连线
        //   allowLoop: false,
        //   // 拖动边时,是否高亮显示所有可用连接桩或节点
        //   highlight: true,
        // },
        connecting: {
          router: "manhattan",
          connector: {
            name: "rounded",
            args: {
              radius: 8,
            },
          },
          anchor: "center",
          connectionPoint: "anchor",
          allowBlank: false,
          snap: {
            radius: 20,
          },
          createEdge() {
            return new Shape.Edge({
              attrs: {
                line: {
                  stroke: "#A2B1C3",
                  strokeWidth: 2,
                  targetMarker: {
                    name: "block",
                    width: 12,
                    height: 8,
                  },
                },
              },
              zIndex: 0,
            });
          },
          // validateConnection({ targetMagnet }) {
          //   return !!targetMagnet
          // },
        },

        modes: {
          default: ["drag-node"],
        },
        background: {
          color: "#F2F7FA",
        },
        mousewheel: {
          // 是否开启滚轮缩放交互
          enabled: true,
          // 滚动缩放因子 默认 1.2
          factor: 1.2,
          // 是否将鼠标位置作为中心缩放、默认为true
          zoomAtMousePosition: true,
          // 按下什么键 才会缩放
          modifiers: ["ctrl", "meta"],
          // 判断什么情况下 滚轮事件被处理
          // guard: false,
        },
        connector: {
          name: "rounded",
          args: {
            radius: 8,
          },
        },
      });

      // 支持拖拽
      this.dnd = new Dnd({
        target: graph,
        scaled: false,
      });

      Graph.registerNode(
          "custom-node-width-port",
          {
            inherit: "rect",
            width: 100,
            height: 40,
            attrs: {
              body: {
                stroke: "#8f8f8f",
                strokeWidth: 1,
                fill: "#fff",
                rx: 6,
                ry: 6,
              },
              textWrap: {
                text: '这是一段很长的文字内容，超过节点宽度时将会自动换行',
                width: 100, // 节点宽度
                height: 40, // 节点高度
                ellipsis: true, // 超出节点高度是否显示省略号
              },
            },
            // 上下左右 四条边都有连接桩
            ports: this.ports,
          },
          true
      );

      Graph.registerNode(
          "custom-circle-start",
          {
            inherit: "circle",
            ports: this.ports,
          },

          true
      );

      Graph.registerNode(
          "custom-polygon",
          {
            inherit: "polygon",
            points: "0,10 10,0 20,10 10,20",
            ports: this.ports,
          },
          true
      );

      Graph.registerNode(
          "custom-rect",
          {
            inherit: "rect",
            ports: this.ports,
          },
          true
      );

      // graph.addNode({
      //   x: 100,
      //   y: 40,
      //   width: 180,
      //   height: 30,
      //   label: "中心主题",
      //   shape: "custom-node-width-port", // 节点形状
      //   attrs: {
      //     body: {
      //       fill: "#f5f5f5",
      //       stroke: "#333",
      //     },
      //     type: "root",
      //   },
      //   tools: [
      //     {
      //       name: "boundary",
      //       args: {
      //         attrs: {
      //           fill: "#16B8AA",
      //           stroke: "#2F80EB",
      //           strokeWidth: 1,
      //           fillOpacity: 0.1,
      //         },
      //       },
      //     },
      //   ],
      // });

      // 添加 plugin 插件
      graph
          .use(
              new Transform({
                resizing: true,
                rotating: true,
              })
          )
          .use(new Keyboard()) // 键盘事件
          .use(
              new Selection({
                enabled: true,
                multiple: true,
                rubberband: true,
                movable: true,
                showEdgeSelectionBox: true,
                showNodeSelectionBox: true,
                pointerEvents: "none",
              })
          ) // 绑定框选
          .use(
              new Snapline({
                enabled: true,
                sharp: true,
              })
          ) // 对齐线
          .use(new Clipboard())
          .use(new History({ enabled: true })); // 绑定撤销

      // 鼠标事件
      this.mouseEvent(graph);

      // 键盘时间
      this.keyboardEvent(graph);

      // 添加子节点的逻辑...
      return graph;
    },
    initData() {
      // 获取回显数据
      // this.nodeJson = nodeData
      // this.graphOut.fromJSON(this.nodeJson);
    },
    addChildNode(nodeId, type) {
      console.log(nodeId, type);
    },
    handleClose() {
      // this.$confirm("确认关闭？")
      //   .then(() => {
      this.drawer = false;
      // })
      // .catch(() => {});
    },
    saveNode() {
      // this.$confirm("确认保存？")
      //   .then(() => {

      this.currentNode["label"] = this.editNode["label"];
      console.log(this.editNode);
      console.log(this.currentNode)
      // this.currentNode['shape'] = this.editNode['shape']
      // })
      // .catch(() => {});
      // 关闭当前 抽屉 el-drawer
      this.drawer = false;
    },

    startDrag(type, e) {
      this.startDragToGraph(this.graphOut, type, e);
    },

    startDragToGraph(graph, type, e) {
      console.log("graph", graph);
      const startNode = this.graphOut.createNode({
        shape: "custom-circle-start",
        width: 38,
        height: 38,
        attrs: {
          body: {
            strokeWidth: 1,
            stroke: "#000000",
            fill: "#ffffff",
            rx: 10,
            ry: 10,
          },
        },
      });
      const polygonNode = this.graphOut.createNode({
        shape: "custom-polygon",
        width: 80,
        height: 60,
        attrs: {
          body: {
            strokeWidth: 1,
            stroke: "#000000",
            fill: "#ffffff",
            rx: 10,
            ry: 10,
          },
          label: {
            fontSize: 13,
            fontWeight: "bold",
          },
        },
      });
      //   const rectNode = this.graphOut.createNode({
      //     shape: "custom-rect",
      //     width: 80,
      //     height: 60,
      //     attrs: {
      //       body: {
      //         strokeWidth: 1,
      //         stroke: "#000000",
      //         fill: "#ffffff",
      //         padding:5,
      //         rx: 10,
      //         ry: 10,
      //       },
      //       // text:{
      //       //   textWrap: {
      //       //     text: '',
      //       //     width: 100, // 节点宽度
      //       //     height: 40, // 节点高度
      //       //     // ellipsis: true, // 超出节点高度是否显示省略号
      //       //   },
      //       // },
      //       label: {
      //         fontSize: 13,
      //         fontWeight: "bold",
      //       },
      //     },
      //   });

      const rectNode = this.graphOut.createNode({
        shape: "custom-rect",
        width: 80,
        height: 60,
        attrs: {
          body: {
            strokeWidth: 1,
            stroke: "#000000",
            fill: "#ffffff",
            rx: 10,
            ry: 10,
          },
          label: {
            fontSize: 13,
            fontWeight: "bold",
          },
        },
      });

      const endNode = this.graphOut.createNode({
        shape: "custom-circle-start",
        width: 38,
        height: 38,
        key: "end",
        attrs: {
          body: {
            strokeWidth: 4,
            stroke: "#000000",
            fill: "#ffffff",
            rx: 10,
            ry: 10,
          },
          label: {
            text: "结束",
            fontSize: 13,
            fontWeight: "bold",
          },
        },
      });
      let dragNode;
      if (type === "start") {
        dragNode = startNode;
      } else if (type === "end") {
        dragNode = endNode;
      } else if (type === "rect") {
        dragNode = rectNode;
      } else if (type === "polygon") {
        dragNode = polygonNode;
      }
      console.log("dnd", dragNode, e, type);
      this.dnd.start(dragNode, e);
    },

    // 删除事件 节点
    removeNode(node) {
      this.graphOut.removeNode(node);
    },

    // 鼠标事件
    mouseEvent(graph) {
      // 鼠标事件

      // 鼠标 Hover 时添加按钮
      graph.on("node:mouseenter", ({ node }) => {
        node.addTools({
          name: "button",
          args: {
            x: 0,
            y: 0,
            offset: { x: 18, y: 18 },
            // onClick({ view }) { ... },
          },
        });
      });

      // 鼠标移开时删除按钮
      graph.on("node:mouseleave", ({ node }) => {
        node.removeTools(); // 删除所有的工具
      });

      graph.on("node:dblclick", ({ node }) => {
        // 添加连接桩
        // node.addPort({
        //   group: "top",
        //   attrs: {
        //     circle: {
        //       magnet: true,
        //       stroke: "#8f8f8f",
        //       r: 4,
        //     },
        //   },
        // });
        // 编辑node
        this.currentNode = node;
        this.drawer = true;
      });

      graph.on("edge:mouseenter", ({ cell }) => {
        cell.addTools([
          { name: "vertices" },
          {
            name: "button-remove",
            args: { distance: 20 },
          },
        ]);
      });

      graph.on("node:click", ({ node }) => {
        this.currentNode = node;
        console.log('this.currentNode',this.currentNode)
      });
    },

    // 键盘事件
    keyboardEvent(graph) {
      // 键盘事件
      graph.bindKey("tab", (e) => {
        e.preventDefault();

        const selectedNodes = graph.getCells().filter((item) => item.isNode());
        console.log(selectedNodes);
        if (selectedNodes.length) {
          const node = selectedNodes[0];
          const type = node.attrs["type"];
          this.addChildNode(node.id, type);
        }
      });

      graph.bindKey("delete", () => {
        this.removeNode(this.currentNode);
      });

      graph.bindKey("backspace", () => {
        this.removeNode(this.currentNode);
      });
    },
    // 保存流程
    saveEvent() {
      console.log(this.graphOut.toJSON(), "graphsavejson");
    },
  },
  watch: {
    // currentNode: {
    //   handler (nwVal, old) {
    //   },
    //   immediate: true,
    //   deep: true
    // }
  },
};
</script>

<style>
/* 样式调整 */
#stencil {
  width: 100px;
  height: 100%;
  position: relative;
  display: flex;
  flex-direction: column;
  align-items: center;
  border-right: 1px solid #dfe3e8;
  text-align: center;
  font-size: 12px;
}

.dnd-rect {
  width: 50px;
  height: 30px;
  line-height: 40px;
  text-align: center;
  border: 2px solid #000000;
  border-radius: 6px;
  cursor: move;
  font-size: 12px;
  margin-top: 30px;
}

.dnd-polygon {
  width: 35px;
  height: 35px;
  border: 2px solid #000000;
  transform: rotate(45deg);
  cursor: move;
  font-size: 12px;
  margin-top: 30px;
  margin-bottom: 10px;
}

.dnd-circle {
  width: 35px;
  height: 35px;
  line-height: 45px;
  text-align: center;
  border: 5px solid #000000;
  border-radius: 100%;
  cursor: move;
  font-size: 12px;
  margin-top: 30px;
}

.dnd-start {
  border: 2px solid #000000;
}

.x6-widget-stencil {
  background-color: #f8f9fb;
}

.x6-widget-stencil-title {
  background: #eee;
  font-size: 1rem;
}

.x6-widget-stencil-group-title {
  font-size: 1rem !important;
  background-color: #fff !important;
  height: 40px !important;
}

.x6-widget-transform {
  margin: -1px 0 0 -1px;
  padding: 0px;
  border: 1px solid #239edd;
}

.x6-widget-transform > div {
  border: 1px solid #239edd;
}

.x6-widget-transform > div:hover {
  background-color: #3dafe4;
}

.x6-widget-transform-active-handle {
  background-color: #3dafe4;
}

.x6-widget-transform-resize {
  border-radius: 0;
}

.tool {
  height: 40px;
  position: absolute;
  z-index: 9;
  left: 200px;
  /* width: ; */
}
.anticon {
  font-size: 30px;
}
</style>
