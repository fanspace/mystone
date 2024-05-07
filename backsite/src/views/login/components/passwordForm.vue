<template>
	<el-form ref="loginForm" :model="form" :rules="rules" label-width="0" size="large" @keyup.enter="login">
		<el-form-item prop="user">
			<el-input v-model="form.user" prefix-icon="el-icon-user" clearable :placeholder="$t('login.userPlaceholder')">
				<template #append>
					<el-select v-model="userType" style="width: 130px;">
						<el-option :label="$t('login.admin')" value="admin"></el-option>
						<el-option :label="$t('login.user')" value="user"></el-option>
					</el-select>
				</template>
			</el-input>
		</el-form-item>
		<el-form-item prop="password">
			<el-input v-model="form.password" prefix-icon="el-icon-lock" clearable show-password :placeholder="$t('login.PWPlaceholder')"></el-input>
		</el-form-item>
		<el-form-item prop="pin">
			<el-row v-if="identifyCode">
			  <el-col :span="16">
				<el-input
				  v-model="form.pin"
				  auto-complete="off"
				  :placeholder="$t('login.PINError')"
				/>
			  </el-col>
			  <el-col :span="8">
				<s-identify v-if="!vsnLoading" :identify-code="identifyCode" @click="refreshCode" />
				<el-button
				  v-if="vsnLoading"
				  :loading="vsnLoading"
				  style="
				  width: 120px;
				  height: 45px;
				  float: right;
				  margin: 0px;
				  border: none;
				"
				/>
			  </el-col>
			</el-row>
			<el-input
			  v-else
			  v-model="form.pin"
			  auto-complete="off"
			  :placeholder="$t('login.PINError')"
			/>
		  </el-form-item>
		<el-form-item style="margin-bottom: 10px;">
				<el-col :span="12">
					<el-checkbox :label="$t('login.rememberMe')" v-model="form.autologin"></el-checkbox>
				</el-col>
				<el-col :span="12" class="login-forgot">
					<router-link to="/reset_password">{{ $t('login.forgetPassword') }}？</router-link>
				</el-col>
		</el-form-item>
		<el-form-item>
			<el-button type="primary" style="width: 100%;" :loading="islogin" round @click="login">{{ $t('login.signIn') }}</el-button>
		</el-form-item>
		<div class="login-reg">
			{{$t('login.noAccount')}} <router-link to="/user_register">{{$t('login.createAccount')}}</router-link>
		</div>
	</el-form>
</template>

<script>
import SIdentify from './sidentify'
	export default {
		components: {
			SIdentify
		},
		data() {
			return {
				userType: 'admin',
				vsnLoading: false,
				identifyCodes: '1234567890',
                identifyCode: '',
                showpin: false,
				form: {
					user: "admin",
					password: "admin",
					pin: '',
					autologin: false,
					sncode: '',
					ip: '',
					device: '',
					loginType: 'backend'
				},
				rules: {
					user: [
						{required: true, message: this.$t('login.userError'), trigger: 'blur'}
					],
					password: [
						{required: true, message: this.$t('login.PWError'), trigger: 'blur'}
					],
					pin: [
						{required: true, message: this.$t('login.PINError'), trigger: 'blur'}
					]
				},
				islogin: false,
			}
		},
		watch:{
			userType(val){
				if(val == 'admin'){
					this.form.user = 'admin'
					this.form.password = 'admin'
				}else if(val == 'user'){
					this.form.user = 'user'
					this.form.password = 'user'
				}
			}
		},
		mounted() {
			this.refreshCode()
		},
		methods: {
			async login(){

				var validate = await this.$refs.loginForm.validate().catch(()=>{})
				if(!validate){ return false }

				this.islogin = true
				var data = {
					loginType: this.form.loginType,
					sncode: this.form.sncode,
					pin: this.form.pin,
					username: this.form.user,
					password: encodeURIComponent(this.$TOOL.crypto.BASE64.encrypt(this.form.password))
				}
				//获取token
				var user = await this.$API.auth.token.post(data)
				if(user.success ){
					this.$TOOL.cookie.set("TOKEN", user.data.token, {
						expires: this.form.autologin? 4*60*60 : 0
					})
					this.$TOOL.data.set("USER_INFO", user.data.userInfo)
				}else{
					this.islogin = false
					this.$message.warning(user.msg)
					this.refreshCode()
					return false
				}
				//获取菜单
				var menu = null
				if(this.form.user == 'admin'){
					menu = await this.$API.system.menu.myMenus.get()
				}else{
					menu = await this.$API.demo.menu.get()
				}
				if(menu.code == 200){
					if(menu.data.menu.length==0){
						this.islogin = false
						this.$alert("当前用户无任何菜单权限，请联系系统管理员", "无权限访问", {
							type: 'error',
							center: true
						})
						return false
					}
					this.$TOOL.data.set("MENU", menu.data.menu)
					this.$TOOL.data.set("PERMISSIONS", menu.data.permissions)
					this.$TOOL.data.set("DASHBOARDGRID", menu.data.dashboardGrid)
				}else{
					this.islogin = false
					this.$message.warning(menu.message)
					return false
				}

				this.$router.replace({
					path: '/'
				})
				this.$message.success("Login Success 登录成功")
				this.islogin = false
			},
			refreshCode() {
      this.identifyCode = ''
      this.form.pin = ''
      if (this.form.user && this.form.user !== '') {
        this.makeCode()
      } else {
		this.identifyCode = '0000'
		this.form.sncode = ''
	  }
    },
    async makeCode() {
      this.vsnLoading = true
	 
      const data = await this.$API.auth.pin.getp([this.form.user])
	  console.log(JSON.stringify(data))
      if (data.success) {
        this.identifyCode = data.pin
		this.form.sncode = data.sncode
	
      } else {
        this.$message.error(data.msg)
      }
      this.vsnLoading = false
    }
		},
    
	}
</script>

<style>
</style>
