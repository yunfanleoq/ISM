import Vue from 'vue'

// 自定义指令使弹窗可拖动
// 使用方式: <a-modal v-drag-modal destroyOnClose V-model="a"></a-modal>
// 加上v-drag-modal 一定要加上destroyoncLose属性，否则弹窗不会回到初始位置
Vue.directive('drag-modal', (el, bindings, vnode) => {
    Vue.nextTick(() => {
        let { visible, destroyOnClose } = vnode.componentInstance
        // 防止未定义 destroyOnClose 关闭弹窗时dom未被销毁，指令被重复调用
        if (!visible) return
        let modal = el.getElementsByClassName('ant-modal')[0]
        let header = el.getElementsByClassName('ant-modal-header')[0]
        let footer = el.getElementsByClassName('ant-modal-footer')[0]

        let left = 0
        let top = 0

        // 未定义 destroyOnClose 时，dom未被销毁，关闭弹窗再次打开，弹窗会停留在上一次拖动的位置
        if (!destroyOnClose) {
            left = modal.left || 0
            top = modal.top || 0
        }
        // top 初始值为 offsetTop
        top = top || modal.offsetTop
        header.onmousedown = e => {
            let startX = e.clientX;
            let startY = e.clientY;
            header.left = header.offsetLeft;
            header.top = header.offsetTop-42;
            el.onmousemove = event => {
                let endX = event.clientX;
                let endY = event.clientY;
                modal.left = header.left + (endX - startX) + left;
                modal.top = header.top + (endY - startY) + top;
                modal.style.left = modal.left + 'px'
                modal.style.top = modal.top + 'px'
            }
            el.onmouseup = event => {
                left = modal.left
                top = modal.top
                el.onmousemove = null;
                el.onmouseup = null;
                header.releaseCapture && header.releaseCapture();
            }
            header.setCapture && header.setCapture();
        }

    })
})