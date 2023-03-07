# webhook
造个webhook中间服务器，可以实现写其他东西


## issues表
```
CREATE TABLE `issues` (`id` INTEGER PRIMARY KEY AUTOINCREMENT,`desc` VARCHAR(64) NULL,`handle` VARCHAR(64) NULL,`handleDesc` VARCHAR(64) NULL,`status` VARCHAR(64) NULL,`update` DATETIME DEFAULT CURRENT_TIMESTAMP);
```
字段描述
- id: 事件id,递增,作为唯一标识
- desc: 告警详情
- handle: 处置类型（观察/误报/阻断）
- handleDesc: 处置描述
- status: 事件状态（进行中/关闭）
