package middleware

import (
	log "backos/logger"
	"backos/utils"
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"net/http"
)

// NewAuthorizer returns the authorizer, uses a Casbin enforcer as input
func MustAuthorizer(e *casbin.Enforcer) gin.HandlerFunc {
	return func(c *gin.Context) {
		a := &BasicAuthorizer{enforcer: e}
		mc, ok := c.MustGet("mc").(*utils.MyClaim)
		if !ok {
			//c.Redirect(http.StatusMovedPermanently, "/login")
			c.JSON(http.StatusForbidden, gin.H{"code": 403, "msg": "参数错误"})
			c.Abort()
		}

		if !a.CheckPermission(c.Request, mc.Username, mc.Domain) {
			c.JSON(http.StatusForbidden, gin.H{"code": 403, "msg": "无此权限"})
			c.Abort()
		}
	}
}

// BasicAuthorizer stores the casbin handler
type BasicAuthorizer struct {
	enforcer *casbin.Enforcer
}

// GetUserName gets the user name from the request.
// Currently, only HTTP basic authentication is supported
func (a *BasicAuthorizer) GetUserName(r *http.Request) string {
	//username, _, _ := r.BasicAuth()
	username := "abcabcefgefgagagefgagagefg"
	return username
}

// CheckPermission checks the user/method/path combination from the request.
// Returns true (permission granted) or false (permission forbidden)
func (a *BasicAuthorizer) CheckPermission(r *http.Request, uname string, domain string) bool {
	user := uname
	method := r.Method
	path := r.URL.Path

	//fmt.Println(a.enforcer.Enforce("jdadmin01", "/zjback/proMgr/declare/list", "POST"))
	//fmt.Println(user)
	//fmt.Println(path)
	//fmt.Println(method)
	//fmt.Println(a.enforcer.Enforce(user, path, method))
	ok, err := a.enforcer.Enforce(user, domain, path, method)
	if err != nil {
		log.Error(err.Error())
		return false
	}
	return ok

}

// RequirePermission returns the 403 Forbidden to the client
func (a *BasicAuthorizer) RequirePermission(w http.ResponseWriter) {
	w.WriteHeader(403)
	w.Write([]byte("403 Forbidden\n"))
}
