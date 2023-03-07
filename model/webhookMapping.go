package model

// 根据form获取对应的消息模型
func GetWebhookByForm(form string) MsgModel {
	switch form {
	case "wx":
		return &WxMsgModel{}
	}
	return nil
}
