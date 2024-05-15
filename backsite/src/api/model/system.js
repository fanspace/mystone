import config from "@/config"
import http from "@/utils/request"

export default {
	menu: {
		/*myMenus: {
			url: `${config.API_URL}/system/menu/my/1.6.1`,
			name: "获取我的菜单",
			get: async function(){
				return await http.get(this.url);
			}
		},*/
		myMenus: {
			url: `${config.BACKEND_API_URL}/auth/menu/mine`,
			name: "获取我的菜单",
			get: async function(){
				return await http.get(this.url);
			}
		},
		list: {
			url: `${config.BACKEND_API_URL}/menuMgr/list`,
			name: "获取菜单",
			get: async function(){
				return await http.get(this.url);
			}
		},
	    add: {
		    url: `${config.BACKEND_API_URL}/menuMgr/add`,
		    name: "add菜单",
		    post: async function(data={}){
			 return await http.post(this.url, data);
		    }
	    },
	    update: {
	        url: `${config.BACKEND_API_URL}/menuMgr/update`,
	        name: "edit菜单",
	        post: async function(data={}){
		    return await http.post(this.url, data);
	        }
        },
        del: {
            url: `${config.BACKEND_API_URL}/menuMgr/del`,
            name: "delete菜单",
            post: async function(data={}){
	        return await http.post(this.url, data);
            }
        }
	},
	api:{
		apiTree: {
			url: `${config.BACKEND_API_URL}/apiMgr/list`,
			name: "获取Api_Tree",
			getp: async function(params){
			return await http.getp(this.url, params);
		}
		}
	},
	dic: {
		tree: {
			url: `${config.API_URL}/system/dic/tree`,
			name: "获取字典树",
			get: async function(){
				return await http.get(this.url);
			}
		},
		list: {
			url: `${config.API_URL}/system/dic/list`,
			name: "字典明细",
			get: async function(params){
				return await http.get(this.url, params);
			}
		},
		get: {
			url: `${config.API_URL}/system/dic/get`,
			name: "获取字典数据",
			get: async function(params){
				return await http.get(this.url, params);
			}
		}
	},
	role: {
		list: {
			url: `${config.API_URL}/system/role/list2`,
			name: "获取角色列表",
			get: async function(params){
				return await http.get(this.url, params);
			}
		}
	},
	dept: {
		list: {
			url: `${config.API_URL}/system/dept/list`,
			name: "获取部门列表",
			get: async function(params){
				return await http.get(this.url, params);
			}
		}
	},
	user: {
		list: {
			url: `${config.API_URL}/system/user/list`,
			name: "获取用户列表",
			get: async function(params){
				return await http.get(this.url, params);
			}
		}
	},
	app: {
		list: {
			url: `${config.API_URL}/system/app/list`,
			name: "应用列表",
			get: async function(){
				return await http.get(this.url);
			}
		}
	},
	log: {
		list: {
			url: `${config.API_URL}/system/log/list`,
			name: "日志列表",
			get: async function(params){
				return await http.get(this.url, params);
			}
		}
	},
	table: {
		list: {
			url: `${config.API_URL}/system/table/list`,
			name: "表格列管理列表",
			get: async function(params){
				return await http.get(this.url, params);
			}
		},
		info: {
			url: `${config.API_URL}/system/table/info`,
			name: "表格列管理详情",
			get: async function(params){
				return await http.get(this.url, params);
			}
		}
	},
	tasks: {
		list: {
			url: `${config.API_URL}/system/tasks/list`,
			name: "系统任务管理",
			get: async function(params){
				return await http.get(this.url, params);
			}
		}
	}
}
