package model

/**
 *  @Classname Userinfo
 *  @author: lf6128@163.com
 *  @Date: 2024/4/26 15:28
 *  @Description:
 */

type UserInfo struct {
	UserId    int64    `json:"userId"`
	Username  string   `json:"username"`
	Showname  string   `json:"showname"`
	Dashboard string   `json:"dashboard"`
	Roles     []string `json:"roles"`
	Avatar    string   `json:"avatar"`
}
