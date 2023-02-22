package model

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"webhook/config"
)

/*
适配企业微信的消息体
构造透传，再处理的webhook服务器
*/
type MsgModel struct {
	Content             string   `json:"content" binding:"required"`
	MentionedList       []string `json:"mentioned_list"`        //@谁，需要指定userid
	MentionedMobileList []string `json:"mentioned_mobile_list"` // 如果获取不到userid，可以使用手机号
}

// 用于企微机器人消息的透传
type WxMsgModel struct {
	MsgType  string   `json:"msgtype" binding:"required"`
	Markdown MsgModel `json:"markdown"`
	Text     MsgModel `json:"text"`
}

// 封装成json字符串
func (medol *WxMsgModel) String() string {
	jsons, _ := json.Marshal(medol)
	return string(jsons)
}

func (text *WxMsgModel) Send() error {
	resp, err := http.Post(config.Config.WxServerSetting.Api, "application/json", strings.NewReader(text.String()))
	if err != nil {
		return err
	}
	defer resp.Body.Close()                //设置最后关闭resp
	body, err := ioutil.ReadAll(resp.Body) //获取响应数据
	if err != nil {
		return nil
	}
	fmt.Println(string(body))
	return nil
}

// markdown格式的内容
type WxMsgMdModel struct {
	MsgType  string   `json:"msgtype" binding:"required"`
	Markdown MsgModel `json:"markdown" binding:"required"`
}

// 封装成json字符串
func (medol *WxMsgMdModel) String() string {
	jsons, _ := json.Marshal(medol)
	return string(jsons)
}

func (text *WxMsgMdModel) Send() error {
	resp, err := http.Post(config.Config.WxServerSetting.Api, "application/json", strings.NewReader(text.String()))
	if err != nil {
		return err
	}
	defer resp.Body.Close()                //设置最后关闭resp
	body, err := ioutil.ReadAll(resp.Body) //获取响应数据
	if err != nil {
		return nil
	}
	fmt.Println(string(body))
	return nil
}

// text格式的内容
type WxMsgTextModel struct {
	MsgType string   `json:"msgtype"`
	Text    MsgModel `json:"text"`
}

// 封装成json字符串
func (medol *WxMsgTextModel) String() string {
	jsons, _ := json.Marshal(medol)
	return string(jsons)
}

func (text *WxMsgTextModel) Send() error {
	resp, err := http.Post(config.Config.WxServerSetting.Api, "application/json", strings.NewReader(text.String()))
	if err != nil {
		return err
	}
	defer resp.Body.Close()                //设置最后关闭resp
	body, err := ioutil.ReadAll(resp.Body) //获取响应数据
	if err != nil {
		return nil
	}
	fmt.Println(string(body))
	return nil
}
