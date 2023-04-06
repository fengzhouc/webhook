# webhook
造个webhook中间服务器，可以实现写其他东西
## 告警生命周期管理
如果通过告警的接口，就会产生告警，并管理告警的生命周期
- 创建: 本地保存告警内容
- 处置记录: 会记录确认处理的结果（包含告警类型/处置动作/处置描述/责任人等信息）
- 关闭: 处理完后手动关闭，不关闭的话，会根据指定时间周期进行告警提醒，直到告警关闭
## 重复告警抑制
告警肯定是会有重复的时候的，常规处置就是告警平台上处置，如加白/忽略等

但如果发送到我平台上会怎么做
- 重复告警会正常告警提醒，不会入库，但会新增数量
- 如果是历史告警，然后已经关闭处理了，这时会激活重新告警


## issues表
```
CREATE TABLE `issues` (`id` INTEGER PRIMARY KEY AUTOINCREMENT,`issueId` VARCHAR(64) NULL,`form` VARCHAR(64) NULL,`desc` VARCHAR(64) NULL,`issueType` VARCHAR(64) NULL,`handle` VARCHAR(64) NULL,`handleDesc` VARCHAR(64) NULL,`status` VARCHAR(64) NULL,`owner` VARCHAR(64) NULL,`orgmsg` VARCHAR(64) NULL,`count` INTEGER 1,`update` DATETIME DEFAULT CURRENT_TIMESTAMP);
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
- count: 记录告警次数/天
- orgmsg: 源消息体
