/**
 * amis v6.1.0
 * build time: 2024-01-31
 * Copyright 2018-2024 baidu
 */

var ICONS = [
    {
        name: 'Font Awesome 4.7',
        prefix: 'fa fa-',
        icons: [
            'slideshare',
            'snapchat',
            'snapchat-ghost',
            'snapchat-square',
            'soundcloud',
            'spotify',
            'stack-exchange',
            'stack-overflow'
        ]
    }
];
function setIconVendor(icons) {
    ICONS = icons;
}

export { ICONS, setIconVendor };
window.amisVersionInfo={version:'6.1.0',buildTime:'2024-01-31'};
