
import { Dnd } from '@antv/x6'
import { uuid } from 'vue-uuid';
export function initDnd(graph, container) {
    const dnd = new Dnd({
        target: graph,
        getDragNode: (node) => node.clone(),
        getDropNode: (node) => node.clone()
    })

    return {

        startDrag: (component, e) => {
            const detail = JSON.parse(JSON.stringify(component || {}))
            detail.identifier = uuid.v1()
            if((typeof detail.animate!="undefined")&&(typeof detail.animate.move=="undefined"))
            {
                detail.animate.move = {
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
            detail.name = detail.type
            if (!detail.style) {
                detail.style = {}
            }
            detail.style.visible = 1
            detail.style.borderWidth = detail.style.borderWidth
                ? detail.style.borderWidth
                : 0
            detail.style.BorderEdges = detail.style.BorderEdges
                ? detail.style.BorderEdges
                : 0
            detail.style.opacity = detail.style.opacity
                ? detail.style.opacity
                : 1
            detail.style.borderStyle = detail.style.borderStyle
                ? detail.style.borderStyle
                : "solid"
            detail.style.borderColor = detail.style.borderColor
                ? detail.style.borderColor
                : "#ccccff"
            const ports = {
                groups: {
                    top: {
                        position: 'top',
                        attrs: {
                            circle: {
                                r: 4,
                                magnet: true,
                                stroke: '#5F95FF',
                                strokeWidth: 1,
                                fill: '#fff',
                                style: {
                                    visibility: 'hidden',
                                },
                            },
                        },
                    },
                    right: {
                        position: 'right',
                        attrs: {
                            circle: {
                                r: 4,
                                magnet: true,
                                stroke: '#5F95FF',
                                strokeWidth: 1,
                                fill: '#fff',
                                style: {
                                    visibility: 'hidden',
                                },
                            },
                        },
                    },
                    bottom: {
                        position: 'bottom',
                        attrs: {
                            circle: {
                                r: 4,
                                magnet: true,
                                stroke: '#5F95FF',
                                strokeWidth: 1,
                                fill: '#fff',
                                style: {
                                    visibility: 'hidden',
                                },
                            },
                        },
                    },
                    left: {
                        position: 'left',
                        attrs: {
                            circle: {
                                r: 4,
                                magnet: true,
                                stroke: '#5F95FF',
                                strokeWidth: 1,
                                fill: '#fff',
                                style: {
                                    visibility: 'hidden',
                                },
                            },
                        },
                    },
                },
                items: [
                    {
                        group: 'top',
                    },
                    {
                        group: 'right',
                    },
                    {
                        group: 'bottom',
                    },
                    {
                        group: 'left',
                    },
                ],
            }
            const node = graph.createNode({
                shape:detail.type,
                width: detail.style.position.w,
                height: detail.style.position.h,
                zIndex: parseInt(detail.style.zIndex),
                attrs: {
                    body: { fill: '#1890ff', opacity: 1 }, // 初始透明度为1
                },
                data: {
                    locked:false,
                    UpdateNodeFlag:true,
                    editMode: true,
                    showDeviceUuid:"",
                    IsToolBox:false,
                    detail:detail
                },
                ports: { ...ports },
            })
            dnd.start(node, e)
        }
    }
}
