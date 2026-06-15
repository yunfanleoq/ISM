

//射线法-边界检测-只能检测凸边形（检测点在不在多边形里面）
function inRange(x, y, points) {
    // points表示多边形的顶点集合     
    let inside = false;
    for (let i = 0, j = points.length - 1; i < points.length; j = i++) {
        let xi = points[i][0],
            yi = points[i][1];
        let xj = points[j][0],
            yj = points[j][1];
        let intersect = ((yi > y) !== (yj > y)) && (x < (xj - xi) * (y - yi) / (yj - yi) + xi);
        if (intersect)
            inside = !inside;
    }
    return inside;
}
export function checkByPointInRectF2(rectA,rectB) {
    return !(rectB.y+rectB.height < rectA.y || rectB.x> rectA.x +rectA.width ||
        rectB.y > rectA.y + rectA.height|| rectB.x+rectB.width < rectA.x)
}
//碰撞
export function checkByPointInRectF4(rect1, rect2) {
    return (
        rect1.x <= rect2.x + rect2.width &&
        rect1.x + rect1.width >= rect2.x &&
        rect1.y <= rect2.y + rect2.height &&
        rect1.y + rect1.height >= rect2.y
    );
}
//包含
export function checkByPointInRectF5(rect1, rect2) {
    return (
        rect1.x <= rect2.x &&
        rect1.x + rect1.width >= rect2.x + rect2.width &&
        rect1.y <= rect2.y &&
        rect1.y + rect1.height >= rect2.y + rect2.height
    );
}

//判断两个矩形DIV是否重叠----------------------------------------
export function checkByPointInRectF1(targetDom, moveDom) {
    if (targetDom === moveDom) return false;
    let tx1 = parseInt(targetDom.x);
    let ty1 = parseInt(targetDom.y);
    let tx2 = tx1 + parseInt(targetDom.width);
    let ty2 = ty1 + parseInt(targetDom.height);
    let mx1 = parseInt(moveDom.x);
    let my1 = parseInt(moveDom.y);
    let mx2 = mx1 + parseInt(moveDom.width);
    let my2 = my1 + parseInt(moveDom.height);
    let width = Math.min(tx2, mx2) - Math.max(tx1, mx1);
    let height = Math.min(ty2, my2) - Math.max(ty1, my1);
    let stackArea = (width > 0 ? width : 0) * (height > 0 ? height : 0);
    // ----------------------------------------
    if (stackArea > 0) return true;
    return false;
}
/**
 * 矩形边界检测-中心点判断法
 * @param {*} frameRect 
 * @param {*} itemRect 
 */
export function checkByPointInRect(frameRect,itemRect) {
    let x = itemRect.x  ,
        y = itemRect.y ;
    if( x > frameRect.x && y > frameRect.y && x < (frameRect.x + frameRect.width) && y < (frameRect.y + frameRect.height))
        return true;
    else
        return false;
}

/**
 * 矩形边界检测-碰撞检测法
 * 两个矩形中心点在x方向的距离的绝对值小于等于矩形宽度和的二分之一，同时y方向的距离的绝对值小于等于矩形高度和的二分之一
 * @param {*} frameRect
 * @param {*} itemRect
 */
export function checkByRectCollisionDetection(frameRect,itemRect) {
    let x1 = frameRect.x + frameRect.width / 2,
        y1 = frameRect.y + frameRect.height / 2,
        w1 = frameRect.width,
        h1 = frameRect.height;
    let x2 = itemRect.x + itemRect.width / 2,
        y2 = itemRect.y + itemRect.height / 2,
        w2 = itemRect.width,
        h2 = itemRect.height;
    if(Math.abs(x1 - x2) < ((w1 + w2) / 2) && Math.abs(y1 - y2) < ((h1 + h2) / 2)) 
        return true;
    else
        return false;
}
