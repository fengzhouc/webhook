package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"webhook/config"
	"webhook/util"
)

/*
适配企业微信的消息体
构造透传，再处理的webhook服务器
*/

// 用于企微机器人消息的透传
type WxMsgModel struct {
	MsgType  string       `json:"msgtype" binding:"required"`
	Markdown ContentModel `json:"markdown"`
	Text     ContentModel `json:"text"`
}

type ContentModel struct {
	Content             string   `json:"content"`
	MentionedList       []string `json:"mentioned_list"`        //@谁，需要指定userid
	MentionedMobileList []string `json:"mentioned_mobile_list"` // 如果获取不到userid，可以使用手机号
}

// 企微返回的结构体
type ErrMsgModel struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

// 根据json数据初始化medol对象
func (medol *ErrMsgModel) Init(json string) {
	// 将json数据反序列化
	util.JsonToAnyModel(medol, json)
}

// 封装成json字符串
func (medol *WxMsgModel) String() string {
	jsons, _ := json.Marshal(medol)
	return string(jsons)
}

// 发送webhook消息
func (text *WxMsgModel) Send() error {
	resp, err := http.Post(config.Config.WebHookServerSetting.Wx, "application/json", strings.NewReader(text.String()))
	if err != nil {
		fmt.Println("[wx.Send.Post] ", err)
		return err
	}
	defer resp.Body.Close()                //设置最后关闭resp
	body, err := ioutil.ReadAll(resp.Body) //获取响应数据
	if err != nil {
		fmt.Println("[wx.Send.ReadAll] ", err)
		return err
	}
	errmsg := ErrMsgModel{}
	errmsg.Init(string(body))
	if errmsg.ErrCode != 0 {
		fmt.Println("[wx.Send.Err] send faild!! ", string(body))
		return errors.New(string(body))
	}

	return nil
}

// 设置消息内容
func (text *WxMsgModel) SetContent(msg string, msgtype string) {
	if msgtype == "text" {
		text.MsgType = "text"
		text.Text.Content = msg
	} else {
		text.MsgType = "markdown"
		text.Markdown.Content = msg
	}
}
