package model

/**
 *  @Classname SomeModels
 *  @author: lf6128@163.com
 *  @Date: 2024/4/26 15:19
 *  @Description:  View models
 */

type AccountLoginReq struct {
	Username string `json:"username"`
	// 密码
	Password string `json:"password"`
	Pin      string `json:"pin"`
	// 验证码
	Ip     string `json:"ip"`
	Device int32  `json:"device"`
	// 为避免部分数据库中domain为关键字，此处使用domaim
	Domaim string `json:"domain"`
	// user type  考生：fore    单位用户   com     系统管理用户: back
	LoginType string `json:"loginType"`
	// 校验码 （系统生成的随机校验码，获取验证码时由系统分配，提交表单时回传）
	Sncode string `json:"sncode"`
}

type AccountLoginRes struct {
	Token    string    `json:"token"`
	UserInfo *UserInfo `json:"userInfo"`
}
