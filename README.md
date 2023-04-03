# webhook
造个webhook中间服务器，可以实现写其他东西


## issues表
```
CREATE TABLE `issues` (`id` INTEGER PRIMARY KEY AUTOINCREMENT,`issueId` VARCHAR(64) NULL,`form` VARCHAR(64) NULL,`desc` VARCHAR(64) NULL,`issueType` VARCHAR(64) NULL,`handle` VARCHAR(64) NULL,`handleDesc` VARCHAR(64) NULL,`status` VARCHAR(64) NULL,`owner` VARCHAR(64) NULL,`update` DATETIME DEFAULT CURRENT_TIMESTAMP);
```
字段描述
- id: 事件id,递增
- issueId: 事件id，作为唯一标识，使用uuid
- desc: 告警详情
- issueType: 事件类型（信息泄漏/入侵告警）
- handle: 处置类型（观察/误报/阻断）
- handleDesc: 处置描述
- status: 事件状态（进行中/关闭）
- owner: 责任人（公司/部门/姓名）
- form: 记录下来自哪个webhook，后面重发的时候，可以知道
