package logic

import (
	"fmt"
	"strings"
	"sync"

	"github.com/Sirupsen/logrus"
	"github.com/reechou/wxrobot/wxweb"
)

type WxManager struct {
	sync.Mutex
	wxs map[string]*wxweb.WxWeb
}

func NewWxManager() *WxManager {
	wm := &WxManager{
		wxs: make(map[string]*wxweb.WxWeb),
	}
	return wm
}

func (self *WxManager) RegisterWx(wx *wxweb.WxWeb) {
	self.Lock()
	defer self.Unlock()

	nickName, ok := wx.User["NickName"]
	if ok {
		nick := nickName.(string)
		self.wxs[nick] = wx
		logrus.Infof("wx manager register wx[%s] success.", nick)
	}
}

func (self *WxManager) UnregisterWx(wx *wxweb.WxWeb) {
	self.Lock()
	defer self.Unlock()

	nickName, ok := wx.User["NickName"]
	if ok {
		nick := nickName.(string)
		_, ok := self.wxs[nick]
		if ok {
			delete(self.wxs, nick)
			logrus.Infof("wx manager unregister wx[%s] success.", nick)
		}
	}
}

func (self *WxManager) SendMsg(msg *SendMsgInfo, msgStr string) {
	wx := self.wxs[msg.WeChat]
	if wx == nil {
		logrus.Errorf("send msg unknown this wechat[%s].", msg.WeChat)
		return
	}
	switch msg.ChatType {
	case CHAT_TYPE_PEOPLE:
		var userName string
		if msg.UserName != "" {
			userName = msg.UserName
			uf := wx.Contact.Friends[userName]
			if uf == nil {
				uf := wx.Contact.NickFriends[msg.Name]
				if uf == nil {
					logrus.Errorf("unkown this friend[%s]", msg.Name)
					return
				}
				userName = uf.UserName
				logrus.Debugf("send msg to people find username[%s] from name[%s]", userName, msg.Name)
			} else {
				logrus.Debugf("send msg to people find username[%s] from request", userName)
			}
		} else {
			uf := wx.Contact.NickFriends[msg.Name]
			if uf == nil {
				logrus.Errorf("unkown this friend[%s]", msg.Name)
				return
			}
			userName = uf.UserName
		}
		if msg.MsgType == MSG_TYPE_TEXT {
			wx.Webwxsendmsg(msgStr, userName)
		} else if msg.MsgType == MSG_TYPE_IMG {
			mediaId, ok := wx.Webwxuploadmedia(userName, msg.Msg)
			if ok {
				wx.Webwxsendmsgimg(userName, mediaId)
			}
		}
	case CHAT_TYPE_GROUP:
		var userName string
		group := wx.Contact.NickGroups[msg.Name]
		if group == nil {
			logrus.Errorf("unkown this group[%s]", msg.Name)
			return
		}
		userName = group.UserName
		if msg.MsgType == MSG_TYPE_TEXT {
			wx.Webwxsendmsg(msgStr, userName)
		}
	}
}

func (self *WxManager) SendImgMsg(msg *SendImgInfo) {
	wx := self.wxs[msg.WeChat]
	if wx == nil {
		logrus.Errorf("send img msg unknown this wechat[%s].", msg.WeChat)
		return
	}
	mediaId, ok := wx.Webwxuploadmedia(msg.UserName, msg.ImgPath)
	if ok {
		wx.Webwxsendmsgimg(msg.UserName, mediaId)
	}
}

func (self *WxManager) VerifyUser(msg *wxweb.ReceiveMsgInfo) bool {
	wx := self.wxs[msg.BaseInfo.WechatNick]
	if wx == nil {
		logrus.Errorf("unknown this wechat[%s].", msg.BaseInfo.WechatNick)
		return false
	}
	ok := wx.Webwxverifyuser(wxweb.WX_VERIFY_USER_OP_CONFIRM, "", msg.AddFriend.Ticket, msg.BaseInfo.FromUserName)
	if ok {
		logrus.Infof("verigy user[%s] success.", msg.BaseInfo.FromNickName)
	} else {
		logrus.Infof("verigy user[%s] error.", msg.BaseInfo.FromNickName)
	}
	return ok
}

func (self *WxManager) CheckGroup() {

}

func (self *WxManager) StateGroupNum(wechat, g string) string {
	wx := self.wxs[wechat]
	if wx == nil {
		logrus.Errorf("unknown this wechat[%s].", wechat)
		return ""
	}
	g = strings.Replace(g, "\n", "", -1)
	result := "[奸笑][奸笑][奸笑]【%s | 群】群总数-%d 去重群成员数-%d 重复成员数-%d"
	allGroupNum := 0
	cfNum := 0
	members := make(map[string]int)
	for _, v := range wx.Contact.Groups {
		if !ExecCheckFunc(g, v.NickName) {
			continue
		}
		allGroupNum++
		for _, v2 := range v.MemberList {
			_, ok := members[v2.UserName]
			if ok {
				cfNum++
				continue
			}
			members[v2.UserName] = 1
		}
	}
	g = ExecGetArgvFunc(g)
	return fmt.Sprintf(result, g, allGroupNum, len(members), cfNum)
}

func (self *WxManager) CheckGroupChat(info *CheckGroupChatInfo) {

}
