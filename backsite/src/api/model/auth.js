import config from "@/config"
import http from "@/utils/request"

export default {
	token: {
		url: `${config.BACKEND_API_URL}/login`,
		name: "登录获取TOKEN",
		post: async function(data={}){
			return await http.post(this.url, data);
		}
	},
	pin: {
		url: `${config.BACKEND_API_URL}/pin`,
		name: "登录获取Sncode",
		getp: async function(params){
			return await http.getp(this.url, params);
		}
	}
}
