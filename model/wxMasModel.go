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
	MsgType      string            `json:"msgtype" binding:"required"`
	Markdown     ContentModel      `json:"markdown"`
	Text         ContentModel      `json:"text"`
	Image        ImageModel        `json:"image"`
	ImageContent ImageContentModel `json:"news"`
	File         FileModel         `json:"file"`
}

// 文本及markdown的内容结构体
type ContentModel struct {
	Content             string   `json:"content"`
	MentionedList       []string `json:"mentioned_list"`        //@谁，需要指定userid
	MentionedMobileList []string `json:"mentioned_mobile_list"` // 如果获取不到userid，可以使用手机号
}

// 图片类型消息的内容结构体
type ImageModel struct {
	Base64 string `json:"base64"`
	Md5    string `json:"md5"`
}

// 图文类型消息的内容结构体
type ImageContentModel struct {
	Articles []ArticlesModel `json:"articles"`
}

type ArticlesModel struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Url         string `json:"url"`
	PicUrl      string `json:"picurl"`
}

// 文件类型消息的内容结构体
type FileModel struct {
	MediaId string `json:"media_id"`
}

// TODO: 待实现上传文件获取media_id，按道理应该不用实现，发送的平台应该会先发送，然后拿到id发webhook消息

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

// 根据msgtype转换成结构化文本，用于消息通知
func (medol *WxMsgModel) NoteString() string {
	switch medol.MsgType {
	case "text":
		return medol.Text.Content
	case "markdown":
		return medol.Markdown.Content
	case "image":
		return fmt.Sprintf("图片消息base64: %s", medol.Image.Base64)
	case "news":
		var news_title string
		for i, article := range medol.ImageContent.Articles {
			news_title += fmt.Sprintf("\n图文消息%d:%s", i, article.Title)
		}
		return news_title
	case "file":
		return fmt.Sprintf("文件消息media_id: %s", medol.File.MediaId)
	}
	return "Unsupported type, please contact the developer!!!"
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

// 用于重发消息时，将数据库中的消息内容重新设置,消息类型默认md格式，重发消息只是为了提醒，详情需要点击查看
func (text *WxMsgModel) SetContent(msg string, msgtype string) {
	switch msgtype {
	case "text":
		text.MsgType = "text"
		text.Text.Content = msg
	default:
		text.MsgType = "markdown"
		text.Markdown.Content = msg
	}
}
