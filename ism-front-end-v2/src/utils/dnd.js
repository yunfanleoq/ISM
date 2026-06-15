
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
            component.identifier = uuid.v1()
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
            component.name = component.type
            component.style.visible = 1
            component.style.borderWidth = component.style.borderWidth
                ? component.style.borderWidth
                : 0
            component.style.BorderEdges = component.style.BorderEdges
                ? component.style.BorderEdges
                : 0
            component.style.opacity = component.style.opacity
                ? component.style.opacity
                : 1
            component.style.borderStyle = component.style.borderStyle
                ? component.style.borderStyle
                : "solid"
            component.style.borderColor = component.style.borderColor
                ? component.style.borderColor
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
                shape:component.type,
                width: component.style.position.w,
                height: component.style.position.h,
                zIndex: parseInt(component.style.zIndex),
                attrs: {
                    body: { fill: '#1890ff', opacity: 1 }, // 初始透明度为1
                },
                data: {
                    locked:false,
                    UpdateNodeFlag:true,
                    editMode: true,
                    showDeviceUuid:"",
                    IsToolBox:false,
                    detail:component
                },
                ports: { ...ports },
            })
            dnd.start(node, e)
        }
    }
}
