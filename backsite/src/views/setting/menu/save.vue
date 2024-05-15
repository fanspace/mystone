<template>
	<el-row :gutter="40">
		<el-col v-if="!form.id">
			<el-empty description="请选择左侧菜单后操作" :image-size="100"></el-empty>
		</el-col>
		<template v-else>
			<el-col :lg="12">
				<h2>{{form.meta.title || "新增菜单"}}</h2>
				<el-form :model="form" :rules="rules" ref="menuForm" label-width="80px" label-position="left">
					<el-form-item label="显示名称" prop="meta.title">
						<el-input v-model="form.meta.title" clearable placeholder="菜单显示名字"></el-input>
					</el-form-item>
					<el-form-item label="组名" prop="groupName">
						<el-input v-model="form.groupName" clearable placeholder="菜单显示名字"></el-input>
					</el-form-item>
					<el-form-item label="上级菜单" prop="pid">
						<el-cascader v-model="form.pid" :options="menuOptions" :props="menuProps" :show-all-levels="false" placeholder="顶级菜单" clearable disabled></el-cascader>
					</el-form-item>
					<el-form-item label="类型" prop="meta.type">
						<el-radio-group v-model="form.meta.type">
							<el-radio-button label="menu">菜单</el-radio-button>
							<el-radio-button label="iframe">Iframe</el-radio-button>
							<el-radio-button label="link">外链</el-radio-button>
							<el-radio-button label="button">按钮</el-radio-button>
						</el-radio-group>
					</el-form-item>
					<el-form-item label="别名" prop="name">
						<el-input v-model="form.name" clearable placeholder="菜单别名"></el-input>
						<div class="el-form-item-msg">系统唯一且与内置组件名一致，否则导致缓存失效。如类型为Iframe的菜单，别名将代替源地址显示在地址栏</div>
					</el-form-item>
					<el-form-item label="菜单图标" prop="meta.icon">
						<sc-icon-select v-model="form.meta.icon" clearable></sc-icon-select>
					</el-form-item>
					<el-form-item label="路由地址" prop="path">
						<el-input v-model="form.path" clearable placeholder=""></el-input>
					</el-form-item>
					<el-form-item label="重定向" prop="redirect">
						<el-input v-model="form.redirect" clearable placeholder=""></el-input>
					</el-form-item>
					<el-form-item label="菜单高亮" prop="active">
						<el-input v-model="form.active" clearable placeholder=""></el-input>
						<div class="el-form-item-msg">子节点或详情页需要高亮的上级菜单路由地址</div>
					</el-form-item>
					<el-form-item label="视图" prop="component">
						<el-input v-model="form.component" clearable placeholder="">
							<template #prepend>views/</template>
						</el-input>
						<div class="el-form-item-msg">如父节点、链接或Iframe等没有视图的菜单不需要填写</div>
					</el-form-item>
					<el-form-item label="颜色" prop="color">
						<el-color-picker v-model="form.color" :predefine="predefineColors"></el-color-picker>

					</el-form-item>
					<el-form-item label="是否隐藏" prop="hidden">
						<el-checkbox v-model="form.hidden">隐藏菜单</el-checkbox>
						<el-checkbox v-model="form.hiddenBreadcrumb">隐藏面包屑</el-checkbox>
						<div class="el-form-item-msg">菜单不显示在导航中，但用户依然可以访问，例如详情页</div>
					</el-form-item>
					<el-form-item label="整页路由" prop="fullpage">
						<el-switch v-model="form.fullpage" />
					</el-form-item>
					<el-form-item label="标签" prop="tag">
						<el-input v-model="form.tag" clearable placeholder=""></el-input>
					</el-form-item>
					<el-form-item>
						<el-button type="primary" @click="save" :loading="loading">保 存</el-button>
					</el-form-item>
				</el-form>

			</el-col>
			<el-col :lg="12" class="apilist">
				<h2>接口权限</h2>
				<el-row>
					<el-col :span="21" >
						<el-tree-select
						v-model="valueStrictly"
						:data="cacheData"
						multiple
						
						:render-after-expand="false"
						show-checkbox
						check-strictly
						check-on-click-node
						lazy
    :load="loadNode"
    :props="apipros"
    
						style="width: 470px"
					  /></el-col>
  <el-col :span="3" >
	<el-button type="success" >ok</el-button>
  </el-col>
</el-row>
				<sc-form-table v-model="form.apiList" :addTemplate="apiListAddTemplate" placeholder="暂无匹配接口权限">
					<el-table-column prop="code" label="标识" width="150">
						<template #default="scope">
							<el-input v-model="scope.row.code" placeholder="请输入内容"></el-input>
						</template>
					</el-table-column>
					<el-table-column prop="url" label="Api url">
						<template #default="scope">
							<el-input v-model="scope.row.url" placeholder="请输入内容"></el-input>
						</template>
					</el-table-column>
				</sc-form-table>

			</el-col>
		</template>
	</el-row>

</template>

<script>
	import scIconSelect from '@/components/scIconSelect'

	export default {
		components: {
			scIconSelect
		},
		props: {
			menu: { type: Object, default: () => {} },
			formAction: {type: String, default: ''}
		},
		data(){
			return {
				valueStrictly: [],
				options: [{value:'001', label: 'agrp', options:[{value:'0001', label:'smgar1'}, {value:'00002', label:'sgmt21'}]},
{value:'0021', label: 'bgrp', options:[{value:'00021', label:'smgar21'}, {value:'000022', label:'sgmt221'}]},
{value:'003', label: 'cgrp', options:[{value:'00031', label:'smgar31'}, {value:'000032', label:'sgmt321'}]}
],
                apipros: {
  label: 'name',
  children: 'children',
  isLeaf: 'isLeaf',
  level: 'level',
  id: 'id',
  nameCn: 'nameCn'
},
cacheData:[],
				form: {
					id: "",
					pid: "",
					name: "",
					path: "",
					component: "",
					redirect: "",
					fullpage: false,
					tag: "",
					active: "",
					color: "",
					domain: 'backend',
					groupName: '',
					meta:{
						title: "",
						icon: "",
						type: "menu",
					},
					hidden: false,
					hiddenBreadcrumb: false,
					apiList: []
				},
				apiOptions: [],
				menuOptions: [],
				menuProps: {
					value: 'id',
					label: 'title',
					checkStrictly: true
				},
				predefineColors: [
					'#ff4500',
					'#ff8c00',
					'#ffd700',
					'#67C23A',
					'#00ced1',
					'#409EFF',
					'#c71585'
				],
				rules: {
					groupName: [
          { required: true, message: '请输菜单组名称', trigger: 'blur' },
          { min: 2, max: 50, message: '长度在 2 到 20 个字符', trigger: 'blur' }
        ],
		meta:{title: [
          { required: true, message: '请输显示名称', trigger: 'blur' },
          { min: 2, max: 50, message: '长度在 2 到 20 个字符', trigger: 'blur' }
        ]},
		name: [
          { required: true, message: '请输别名（组件名称）', trigger: 'blur' },
          { min: 2, max: 50, message: '长度在 2 到 20 个字符', trigger: 'blur' }
        ],
		path: [
          { required: true, message: '请输路由地址', trigger: 'blur' },
          { min: 2, max: 50, message: '长度在 2 到 50 个字符', trigger: 'blur' }
        ]
		},
				apiListAddTemplate: {
					code: "",
					url: ""
				},
				loading: false
			}
		},
		watch: {
			menu: {
				handler(){
					this.menuOptions = this.treeToMap(this.menu)
				},
				deep: true
			}
		},
		mounted() {
			this.apiTreeCacheLoader()
		},
		methods: {
			// api tree data
			async apiTreeCacheLoader() {
				this.loading = true
				var res = await this.$API.system.api.apiTree.getp(0)
				this.loading = false
				
				if(res.success){
						this.cacheData = res.data
						console.log(JSON.stringify(this.cacheData))
				}else{
					this.$message.warning(res.msg)
				}
			},
			//简单化菜单
			treeToMap(tree){
				const map = []
				tree.forEach(item => {
					var obj = {
						id: item.id,
						pid: item.pid,
						title: item.meta.title,
						children: item.children&&item.children.length>0 ? this.treeToMap(item.children) : null
					}
					map.push(obj)
				})
				return map
			},
			//保存
			save(){
				
				this.$refs['menuForm'].validate((valid) => {
        if (valid) {
			this.submitMenu()
		} else {
			this.$message.error('表单提交错误')
		}
	})
			
			},
			async submitMenu() {
				const dataform = {
					action: 2,
					menu: this.form
				}
				this.loading = true
				var res = await this.$API.system.menu.update.post(dataform)
				this.loading = false
				if(res.success){
					this.$message.success("操作成功")
				}else{
					this.$message.warning(res.msg)
				}

			},
			
			//表单注入数据
			setData(data, pid){
				this.form = data
				this.form.apiList = data.apiList || []
				this.form.pid = pid
			},
			loadNode(node, resolve) {
				if (node.level === 0) {
        return resolve(this.cacheData)
      } else {
        this.loacChildNode(node, resolve)
      }
			},
			async loacChildNode(node, resolve){
				this.loading = true
				var res = await this.$API.system.api.apiTree.getp(node.data.id)
				this.loading = false
				if(res.success){
					var ttm = res.data.map(x => {
						x.isLeaf = true
						return x
					})
					return resolve(ttm)
				}else{
					this.$message.warning(res.msg)
					return resolve([])
				}
			}

		}
	}
</script>

<style scoped>
	h2 {font-size: 17px;color: #3c4a54;padding:0 0 30px 0;}
	.apilist {border-left: 1px solid #eee;}

	[data-theme="dark"] h2 {color: #fff;}
	[data-theme="dark"] .apilist {border-color: #434343;}
</style>
