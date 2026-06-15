import { sin, cos } from '@/utils/translate';

export function getStyle(style, filter = []) {
    const needUnit = ['fontSize', 'width', 'height', 'top', 'left', 'borderWidth', 'letterSpacing', 'borderRadius'];

    const result = {};
    Object.keys(style).forEach((key) => {
        if (!filter.includes(key)) {
            if (key != 'rotate') {
                result[key] = style[key];

                if (needUnit.includes(key)) {
                    result[key] += 'px';
                }
            } else {
                result.transform = key + '(' + style[key] + 'deg)';
            }
        }
    });

    return result;
}

// 获取一个组件旋转 rotate 后的样式
export function getComponentRotatedStyle(style,scale) {
    style = { ...style };
    const zoom = scale/100
    if (style.transform != 0)
    {
        style.position.b = parseInt(style.position.y)+parseInt(style.position.h);
        style.position.r = parseInt(style.position.x) + parseInt((style.position.w));
        return style
    }
    // if (style.transform != 0) {
    //     const newWidth = parseInt(style.position.w) * cos( parseInt(style.transform)) + parseInt(style.position.h) * sin( parseInt(style.transform));
    //     const diffX = (parseInt(style.position.w) - newWidth) / 2; // 旋转后范围变小是正值，变大是负值
    //     style.position.x = parseInt(style.position.x)+parseInt(diffX)
    //     style.position.r = parseInt(style.position.x) + parseInt(newWidth);
    //
    //     const newHeight = parseInt(style.position.h) * cos( parseInt(style.transform)) + parseInt(style.position.w) * sin( parseInt(style.transform));
    //     const diffY = (newHeight - parseInt(style.position.h)) / 2; // 始终是正
    //     style.position.y = parseInt(style.position.y)+parseInt(diffY);
    //     style.position.b = parseInt(style.position.y) + parseInt(newHeight);
    //
    //     style.position.w = parseInt(newWidth);
    //     style.position.h = parseInt(newHeight);
    //
    // }
    else {
        style.position.b = parseInt(style.position.y)+parseInt(style.position.h);
        style.position.r = parseInt(style.position.x) + parseInt(style.position.w);
    }

    return style;
}
