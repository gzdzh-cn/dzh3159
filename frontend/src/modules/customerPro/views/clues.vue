<template>
	<cl-crud ref="Crud">
		<cl-row>
			<!-- 刷新按钮 -->
			<cl-refresh-btn />
			<!-- 新增按钮 -->
			<cl-add-btn />
			<!-- 删除按钮 -->
			<cl-multi-delete-btn />

			<el-button plain @click="toggleRowExpansion()" style="margin-right: 10px">
				{{ isExtend ? "全部收起" : "全部展开" }}
			</el-button>
			<el-button
				plain
				@click="openExcelDialog()"
				style="margin-right: 10px"
				v-if="service.customer_pro.clues._permission.excel"
			>
				导出Excel
			</el-button>

			<el-button
				type="warning"
				:disabled="Table?.selection.length == 0"
				@click="openDistribute"
				v-permission="service.customer_pro.clues._permission.distribute"
				>分配</el-button
			>
			<!-- <el-popconfirm title="会清空当前系统的数据，确定迁移数据吗?" @confirm="migrateData()">
				<template #reference>
					<el-button v-if="isAdmin" bg type="danger">迁移数据</el-button>
				</template>
			</el-popconfirm> -->

			<!-- <el-popconfirm title="清空当前系统全部数据吗?" @confirm="clearData()">
				<template #reference>
					<el-button v-if="isAdmin" bg type="info">清除数据</el-button>
				</template>
			</el-popconfirm> -->

			<!-- <el-button plain @click="randomDialog = true" style="margin-right: 10px" v-if="isAdmin">
				生成随机数据
			</el-button> -->

			<!-- <cl-export-btn
				:columns="Table?.columns"
				:data="onExportData"
				:disabled="Table?.selection.length == 0"
			/> -->

			<cl-flex1 />

			<el-button type="info" text bg :icon="Search" v-show="searchStatus">
				正在搜索中
			</el-button>

			<!-- 高级按钮 -->
			<cl-adv-btn />
			<!-- 关键字搜索 -->
			<cl-search-key ref="searchRef" />
		</cl-row>

		<cl-row>
			<!-- 数据表格 -->
			<cl-table
				ref="Table"
				:border="false"
				:default-expand-all="true"
				:row-class-name="tableRowClassName"
			>
				<template #column-detail="{ scope }">
					<div style="padding: 0 30px">
						<p v-if="isAdmin">唯一码: {{ scope.row?.id }}</p>
						<p>guestId: {{ scope.row?.guestId }}</p>
						<p v-if="scope.row?.source_from == 1 && scope.row?.createdName">
							创建者: {{ scope.row?.createdName }}
						</p>
						<p>账户: {{ scope.row?.account_name }}</p>
						<p v-if="scope.row?.services_names" style="color: #d83b01">
							分配过的客服: {{ scope.row?.services_names }}
						</p>
						<p>来源: {{ sourceFormatter(scope.row?.source_from) }}</p>
						<p v-if="scope.row?.level">
							线索等级: {{ levelFormatter(scope.row?.level) }}
						</p>
						<p>IP归属地: {{ scope.row?.guestIpInfo }}</p>
						<p>
							最后跟进时间:
							{{
								scope.row?.last_followup_time
									? scope.row?.last_followup_time
									: scope.row?.createTime
							}}
						</p>
						<p>
							创建时间:
							{{ scope.row?.createTime }}
						</p>
						<p v-if="scope.row?.remark">备注: {{ scope.row?.remark }}</p>
						<p v-if="scope.row?.filterRemark">
							过滤原因: {{ scope.row?.filterRemark }}
						</p>
					</div>
				</template>

				<template #column-status="{ scope }">
					<span style="color: #d83b01" v-if="scope.row.status == 0"> 未成交 </span>
					<span style="color: #00b294" v-if="scope.row.status == 1"> 已成交 </span>
				</template>

				<template #slot-op="{ scope }">
					<div
						style="
							display: flex;
							flex-direction: row;
							flex-wrap: wrap;
							align-items: center;
							gap: 12px;
						"
					>
						<!-- <el-button
							text
							bg
							type="primary"
							@click="edit(scope.row)"
							v-if="
								scope.row.status == 0 &&
								service.customer_pro.clues._permission.update
							"
							>审核</el-button
						> -->

						<el-button text bg type="info" @click="openFollow(scope.row)"
							>跟进</el-button
						>
						<el-button
							text
							bg
							type="warning"
							@click="openTracks(scope.row)"
							v-if="service.customer_pro.clues._permission.getTrackList"
							>轨迹</el-button
						>
						<el-button
							text
							bg
							type="success"
							@click="openOrderAdd(scope.row)"
							v-if="
								scope.row.status == 0 && service.customer_pro.order._permission.add
							"
							>成交</el-button
						>
					</div>
				</template>
			</cl-table>
		</cl-row>

		<cl-row>
			<cl-flex1 />
			<!-- 分页控件 -->
			<cl-pagination ref="PaginationRef" :slotDom="slotDom" />
		</cl-row>

		<!-- 新增、编辑 -->
		<cl-upsert ref="Upsert">
			<template #slot-school_id="{ scope }">
				<!-- 学校 -->
				<el-select v-model="scope.school_id" @change="schoolChange">
					<el-option
						v-for="item in schoolList"
						:key="item.value"
						:label="item.name"
						:value="item.id"
					/>
				</el-select>
			</template>

			<!-- 专业 -->
			<template #slot-majors_id="{ scope }">
				<el-select v-model="scope.majors_id">
					<el-option
						v-for="item in majorsList"
						:key="item.value"
						:label="item.name"
						:value="item.id"
					/>
				</el-select>
			</template>
		</cl-upsert>

		<!-- 项目分配 -->
		<cl-form ref="distributeFormRef">
			<template #slot-project_id="{ scope }">
				<el-select v-model="scope.project_id" @change="projectChange">
					<el-option
						v-for="item in projectList"
						:key="item.value"
						:label="item.name"
						:value="item.id"
					/>
				</el-select>
			</template>

			<template #slot-group_id="{ scope }">
				<el-select v-model="scope.group_id" @change="groupChange">
					<el-option
						v-for="item in groupList"
						:key="item.value"
						:label="item.name"
						:value="item.id"
					/>
				</el-select>
			</template>

			<template #slot-services_id="{ scope }">
				<el-select v-model="scope.services_id">
					<el-option
						v-for="item in kfList"
						:key="item.value"
						:label="item.name"
						:value="item.userId"
					/>
				</el-select>
			</template>
		</cl-form>

		<!-- 跟进弹窗 -->
		<cl-dialog title="跟进" v-model="visible">
			<sub-follow ref="FollowRef" :id="cluesId" :status="cluesStatus" @cancel="cancel" />
		</cl-dialog>

		<!-- 轨迹弹窗 -->
		<cl-dialog title="轨迹" v-model="trackVisible">
			<sub-track ref="TrackRef" :id="cluesId" />
			<template #footer>
				<div class="dialog-footer">
					<el-button @click="trackVisible = false">关闭</el-button>
				</div>
			</template>
		</cl-dialog>

		<!-- 成交弹窗 -->
		<cl-form ref="OrderFormRef">
			<!-- 学校 -->
			<template #slot-school_id="{ scope }">
				<el-select v-model="scope.school_id" @change="schoolChange">
					<el-option
						v-for="item in schoolList"
						:key="item.value"
						:label="item.name"
						:value="item.id"
					/>
				</el-select>
			</template>

			<!-- 专业 -->
			<template #slot-majors_id="{ scope }">
				<el-select v-model="scope.majors_id">
					<el-option
						v-for="item in majorsList"
						:key="item.value"
						:label="item.name"
						:value="item.id"
					/>
				</el-select>
			</template>
		</cl-form>

		<!-- 高级搜索 -->
		<cl-adv-search ref="AdvSearch" />

		<cl-dialog title="导出Excel" v-model="openExcel">
			<div class="exportBox">
				<el-button plain @click="toexcel(true)" style="margin-right: 10px">
					导出当前页
				</el-button>
				<el-button
					plain
					@click="toexcel(false)"
					style="margin-right: 10px"
					v-loading.fullscreen.lock="fullscreenLoading"
				>
					导出全部
				</el-button>
			</div>

			<div class="excel-table">
				<excel-down ref="excelDownRef" :clues_status="clues_status" :isAdmin="isAdmin" />
			</div>
		</cl-dialog>

		<cl-dialog title="生成随机数据" v-model="randomDialog">
			<el-form :model="randomForm" label-width="auto" style="max-width: 600px">
				<el-form-item label="生成数量">
					<el-input v-model="randomForm.randomNum" />
				</el-form-item>
				<el-form-item label="年月份">
					<el-date-picker
						v-model="randomForm.dateTime"
						type="month"
						placeholder="日期"
						value-format="YYYY-MM"
					/>
				</el-form-item>
				<el-form-item>
					<el-button type="primary" @click="randomData()">提交</el-button>
				</el-form-item>
			</el-form>
		</cl-dialog>
	</cl-crud>
</template>

<script lang="ts" name="customer_pro-clues" setup>
import { useCrud, useForm, useTable, useUpsert, useAdvSearch } from "@cool-vue/crud";
import { useCool } from "/@/cool";
import { ElMessage, TabsPaneContext, ElLoading } from "element-plus";
import { Search } from "@element-plus/icons-vue";
import { computed, onMounted, ref, watch } from "vue";
import SubFollow from "../components/clues/subFollow.vue";
import SubTrack from "../components/clues/subTrack.vue";
import excelDown from "../components/clues/excelDown.vue";
import { useDict } from "/$/dict";
import { useBase } from "/$/base";

const { service } = useCool();
const { user } = useBase();
const FollowRef = ref(); //跟进
const cluesId = ref(); //线索id
const cluesStatus = ref(0); //线索状态
const projectList = ref(); // 项目列表
const searchStatus = ref(false); // 搜索状态
const searchData = ref(); //搜索条件
const isExtend = ref(true); //展开
const serviceGroup = ref(); // 客服组
const openExcel = ref(false); //打开导出弹窗
const { dict } = useDict();
const clues_status = ref(0);
const searchRef = ref();
const keyWord = ref();
// cl-upsert 配置
const Upsert = useUpsert({
	items: [
		() => {
			return {
				label: "53标识",
				prop: "guestId",
				span: 12,
				hidden: Upsert.value?.mode == "add",
				component: { name: "el-input", props: { disabled: true } }
			};
		},
		() => {
			return {
				label: "项目",
				prop: "projectId",
				span: 12,
				required: true,
				component: { name: "el-select", props: { disabled: Upsert.value?.mode !== "add" } }
			};
		},
		() => {
			return {
				label: "IP",
				prop: "ip",
				span: 12,
				hidden: Upsert.value?.mode == "add",
				component: { name: "el-input", props: { disabled: true } }
			};
		},
		() => {
			return {
				label: "IP归属地",
				prop: "guestIpInfo",
				span: 12,
				hidden: Upsert.value?.mode == "add",
				component: { name: "el-input", props: { disabled: true } }
			};
		},
		{
			label: "姓名",
			prop: "name",
			span: 12,
			component: { name: "el-input" }
		},
		{
			label: "来源",
			prop: "source_from",
			span: 12,
			required: true,
			component: { name: "el-select" }
		},

		{ label: "关键字", prop: "keywords", span: 12, component: { name: "el-input" } },

		{ label: "手机号", prop: "mobile", span: 12, component: { name: "el-input" } },
		{ label: "微信号", prop: "wechat", span: 12, component: { name: "el-input" } },

		{ label: "学历", prop: "education", span: 12, component: { name: "el-select" } },
		{ label: "毕业院校", prop: "graduatedSchool", span: 12, component: { name: "el-input" } },

		{ label: "意向院校", prop: "schoolId", span: 12, component: { name: "slot-school_id" } },
		{ label: "意向专业", prop: "majorsId", span: 12, component: { name: "slot-majors_id" } },

		{ label: "报读类型", prop: "majorsType", span: 12, component: { name: "el-select" } },
		{ label: "报读层次", prop: "degreeId", span: 12, component: { name: "el-select" } },

		{ label: "户口类型", prop: "householdType", span: 12, component: { name: "el-select" } },
		{ label: "户籍地址", prop: "householdAddress", span: 12, component: { name: "el-input" } },

		{
			label: "线索等级",
			prop: "level",
			span: 12,
			component: {
				name: "el-select",
				options: []
			}
		},
		{
			label: "性别",
			prop: "gender",
			span: 12,
			component: {
				name: "el-select"
			}
		},
		{
			label: "紧急联系人电话",
			prop: "emergencyMobile",
			props: {
				labelWidth: "130px"
			},
			span: 12,
			component: { name: "el-input" }
		},
		{
			label: "备注",
			prop: "remark",
			component: { name: "el-input", props: { type: "textarea", rows: 4 } }
		}
	],
	async onOpen() {
		// 学校列表
		getSchoolList();

		// 项目
		const projectList = await service.customer_pro.project.list();
		Upsert.value?.setOptions(
			"projectId",
			projectList.map((e) => {
				return {
					label: e.name,
					value: e.id
				};
			})
		);

		// 报读类型
		const majorsTypeList = await service.customer_pro.readtypes.list();
		Upsert.value?.setOptions(
			"majorsType",
			majorsTypeList.map((e) => {
				return {
					label: e.name,
					value: e.id
				};
			})
		);

		// 报读层次
		const degreeList = await service.customer_pro.readdegree.list();
		Upsert.value?.setOptions(
			"degreeId",
			degreeList.map((e) => {
				return {
					label: e.name,
					value: e.id
				};
			})
		);

		// 户口性质
		Upsert.value?.setOptions("householdType", [
			{
				label: "城镇",
				value: "1"
			},
			{
				label: "农村",
				value: "2"
			}
		]);

		// 线索级别
		Upsert.value?.setOptions("level", dict.get("cluesLevel").value);

		// 性别
		Upsert.value?.setOptions("gender", [
			{
				label: "保密",
				value: "0"
			},
			{
				label: "男",
				value: "1"
			},
			{
				label: "女",
				value: "2"
			}
		]);

		// 来源
		Upsert.value?.setOptions("source_from", [
			{
				label: "手动录入",
				value: "1"
			},
			{
				label: "百度",
				value: "2"
			},
			{
				label: "抖音",
				value: "3"
			},
			{
				label: "53客服",
				value: "4"
			},
			{
				label: "小红书",
				value: "5"
			}
		]);

		// 学历
		Upsert.value?.setOptions("education", [
			{
				label: "未知",
				value: "1"
			},
			{
				label: "初中",
				value: "2"
			},
			{
				label: "高中/中专/中技",
				value: "3"
			},
			{
				label: "大专/高技",
				value: "4"
			},
			{
				label: "本科",
				value: "5"
			}
		]);
	}
});

// cl-table 配置
const Table = useTable({
	columns: [
		{ type: "selection" },
		() => {
			return {
				label: "#",
				type: "expand",
				prop: "detail"
			};
		},

		{
			label: "序号",
			prop: "id",
			width: 100,
			formatter(row, column, value, index) {
				if (row.id > 67114) {
					return row.serialId;
				} else {
					return row.id;
				}
			}
		},
		// { label: "序号", prop: "serialId" },
		{ label: "项目", prop: "project_name" },
		{
			label: "姓名",
			prop: "name",
			width: 150
		},
		{ label: "手机号", prop: "mobile" },
		{ label: "微信号", prop: "wechat" },
		{ label: "关键词", prop: "keywords" },
		{
			label: "跟进状态",
			prop: "followupType",
			width: 80,
			dict: [
				{ label: "待跟进", value: "1", type: "danger" },
				{ label: "电话访谈", value: "2", type: "warning" },
				{ label: "微信沟通", value: "3", type: "warning" },
				{ label: "视频参观", value: "4", type: "info" },
				{ label: "预约参观", value: "5", type: "info" },
				{ label: "已参观", value: "6", type: "success" }
			]
		},
		{
			label: "状态",
			prop: "status",
			width: 80
		},

		{ type: "op", width: 160, buttons: ["slot-op"] }
	]
});

// cl-crud 配置
const Crud = useCrud(
	{
		service: service.customer_pro.clues,
		async onRefresh(params, { render }) {
			searchData.value = params;
			params.status = 0;
			if (params?.keyWord != "") {
				keyWord.value = params?.keyWord;
			}
			const { list, pagination } = await service.customer_pro.clues.page(params);
			render(list, pagination);
		}
	},
	async (app) => {
		app.refresh();
	}
);
// 展开按钮
const toggleRowExpansion = () => {
	isExtend.value = !isExtend.value;
	Table.value?.data.map((item: any) => {
		Table.value?.toggleRowExpansion(item, isExtend.value);
	});
};

// 分页自定义插槽
const selectNum = computed(() => Crud.value?.selection.length);
const slotDom = computed(() => {
	return {
		style: {
			width: "200px",
			fontWeight: "400",
			float: "left",
			marginRight: "auto",
			paddingLeft: "20px"
		},
		htmlText: `已选中 ${selectNum.value} 条`,
		selectNum: selectNum.value
	};
});

// 刷新
const refresh = (params?: any) => {
	Crud.value?.refresh(params);
};

const randomDialog = ref(false);
const randomForm = ref({
	randomNum: 1000,
	dateTime: ""
});
// 生成随机数据
const randomData = async () => {
	await service.customer_pro.clues.randomData(randomForm.value);
	randomDialog.value = false;
	ElMessage.success("数据生成中，请稍等10分钟再刷新页面");
};

// 迁移数据
// const migrateData = async () => {
// 	const info = await service.customer_pro.config.info({ id: 1 });
// 	service.customer_pro.config
// 		.migrateData()
// 		.then(() => {
// 			ElMessage.success("迁移执行中，请稍等5分钟再刷新页面");
// 		})
// 		.catch((e) => {
// 			ElMessage.error(e.message);
// 		});
// };

//清除数据
// const clearData = async () => {
// 	service.customer_pro.config.clearTable().then(() => {
// 		ElMessage.success("清除数据完成");
// 		refresh();
// 	});
// };

// 导出
// const onExportData = async (params: any) => {
// 	return Table.value?.selection;
// };

// 切换类型
// const handleClick = (tab: TabsPaneContext) => {
// 	refresh({ followType: tab.index });
// };

// 编辑
// const edit = (row: any) => {
// 	if (row.status == 0) {
// 		Crud.value?.rowEdit(row);
// 	} else {
// 		Crud.value?.rowInfo(row);
// 	}
// };

// 轨迹弹窗
const trackVisible = ref(false);
const openTracks = (row: any) => {
	cluesId.value = row.id;
	trackVisible.value = true;
};

// 跟进弹窗
const visible = ref(false);
const openFollow = (row: any) => {
	cluesId.value = row.id;
	cluesStatus.value = row.status;
	visible.value = true;
};

// 保存跟进
// const followSave = () => {
// 	FollowRef.value.sub();
// };

// 取消
const cancel = () => {
	visible.value = false;
	refresh();
};

// 保存到公海
// const pushCommonClause = () => {
// 	FollowRef.value.pushCommonClause();
// 	refresh();
// };

// 线索成交
const OrderFormRef = useForm(); //成交表单
const openOrderAdd = async (row: any) => {
	cluesId.value = row.id;
	const cluesOne = await service.customer_pro.clues.info({ id: cluesId.value });
	OrderFormRef.value?.open({
		title: "线索成交",
		items: [
			{
				type: "tabs",
				props: {
					labels: [
						{
							label: "基础信息",
							value: "base"
						},
						{
							label: "个人信息",
							value: "person"
						},
						{
							label: "收款信息",
							value: "financial"
						}
					]
				}
			},
			{
				label: "学生名称",
				prop: "name",
				span: 8,
				component: {
					name: "el-input"
				},
				required: true,
				group: "base"
			},
			{
				label: "学生电话",
				prop: "mobile",
				span: 8,
				component: {
					name: "el-input"
				},
				required: true,
				group: "base"
			},
			{
				label: "接待人员",
				prop: "receiver",
				span: 8,
				component: {
					name: "el-input"
				},
				required: true,
				group: "base"
			},

			{
				label: "身份证",
				prop: "idcard_number",
				span: 8,
				component: {
					name: "el-input"
				},
				required: true,
				group: "base"
			},
			{
				label: "性别",
				prop: "gender",
				span: 8,
				component: {
					name: "el-select",
					options: [
						{
							label: "保密",
							value: "0"
						},
						{
							label: "男",
							value: "1"
						},
						{
							label: "女",
							value: "2"
						}
					]
				},
				required: true,
				group: "base"
			},
			{
				label: "意向院校",
				prop: "school_id",
				span: 8,
				component: { name: "slot-school_id" },
				group: "base"
			},
			{
				label: "意向专业",
				prop: "majors_id",
				span: 8,
				component: { name: "slot-majors_id" },
				group: "base"
			},
			{
				label: "报读类型",
				prop: "majors_type",
				span: 8,
				component: {
					name: "el-select",
					options: []
				},
				required: true,
				group: "base"
			},
			{
				label: "报读层次",
				prop: "degree_id",
				span: 8,
				component: {
					name: "el-select",
					options: []
				},
				required: true,
				group: "base"
			},

			{
				label: "通讯地址",
				prop: "address",
				span: 24,
				component: {
					name: "el-input"
				},
				group: "person"
			},
			{
				label: "紧急联系人",
				prop: "emergency_contact",
				span: 12,
				component: {
					name: "el-input"
				},
				group: "person"
			},

			{
				label: "紧急联系人电话",
				prop: "emergency_mobile",
				props: {
					labelWidth: "130px"
				},
				span: 12,
				component: {
					name: "el-input"
				},
				group: "person"
			},
			{
				label: "微信",
				prop: "wechat",
				span: 8,
				component: {
					name: "el-input"
				},
				group: "person"
			},
			{
				label: "民族",
				prop: "nation",
				span: 8,
				component: {
					name: "el-input"
				},
				group: "person"
			},
			{
				label: "籍贯",
				prop: "native_place",
				span: 8,
				component: {
					name: "el-input"
				},
				group: "person"
			},

			{
				label: "政治面貌",
				prop: "politics_status",
				span: 8,
				component: {
					name: "el-input"
				},
				group: "person"
			},
			{
				label: "户口性质",
				prop: "household_type",
				span: 8,
				component: {
					name: "el-select"
				},
				group: "person"
			},
			{
				label: "户口所在地",
				prop: "household_address",
				span: 8,
				component: {
					name: "el-input"
				},
				group: "person"
			},

			{
				label: "是否应届",
				prop: "freshman",
				span: 8,
				component: {
					name: "el-input"
				},
				group: "person"
			},
			{
				label: "当前学历",
				prop: "education",
				span: 8,
				component: {
					name: "el-select"
				},
				group: "person"
			},
			{
				label: "毕业学校",
				prop: "graduated_school",
				span: 8,
				component: {
					name: "el-input"
				},
				group: "person"
			},

			{
				label: "毕业时间",
				prop: "graduated_date",
				span: 8,
				component: {
					name: "el-input"
				},
				group: "person"
			},
			{
				label: "备注",
				prop: "remark",
				span: 24,
				component: {
					name: "el-input",
					props: {
						type: "textarea",
						rows: 4
					}
				},
				group: "base"
			},

			{
				label: "校方收定位金",
				prop: "school_payment",
				span: 8,
				value: 0.0,
				component: {
					name: "el-input-number",
					props: {
						precision: 2,
						step: 0.1
					}
				},
				group: "financial"
			},
			{
				label: "自收定位金",
				prop: "teams_payment",
				span: 8,
				value: 0.0,
				component: {
					name: "el-input-number",
					props: {
						precision: 2,
						step: 0.1
					}
				},
				group: "financial"
			},
			{
				label: "支付编号",
				prop: "serial",
				span: 8,
				component: {
					name: "el-input"
				},
				group: "financial"
			},
			{
				label: "缴费凭证",
				prop: "voucher",
				span: 8,
				component: {
					name: "cl-upload"
				},
				group: "financial"
			}
		],
		form: {
			...cluesOne
		},
		on: {
			async open() {
				// 学校列表
				getSchoolList();

				// 报读类型
				const majorsTypeList = await service.customer_pro.readtypes.list();
				OrderFormRef.value?.setOptions(
					"majors_type",
					majorsTypeList.map((e) => {
						return {
							label: e.name,
							value: e.id
						};
					})
				);

				// 报读层次
				const degreeList = await service.customer_pro.readdegree.list();
				OrderFormRef.value?.setOptions(
					"degree_id",
					degreeList.map((e) => {
						return {
							label: e.name,
							value: e.id
						};
					})
				);

				// 户口性质
				OrderFormRef.value?.setOptions("household_type", [
					{
						label: "城镇",
						value: "1"
					},
					{
						label: "农村",
						value: "2"
					}
				]);

				// 学历
				OrderFormRef.value?.setOptions("education", [
					{
						label: "未知",
						value: "1"
					},
					{
						label: "初中",
						value: "2"
					},
					{
						label: "高中/中专/中技",
						value: "3"
					},
					{
						label: "大专/高技",
						value: "4"
					},
					{
						label: "本科",
						value: "5"
					}
				]);
			},
			submit(data, { close, done }) {
				data.clues_id = cluesId.value;
				data.services_id = row.services_id;
				data.project_id = row.project_id;
				service.customer_pro.order.add({ ...data }).then((r) => {
					close();
					refresh();
				});
			}
		}
	});
};

// 分配表单打开
const distributeFormRef = useForm(); //分配表单
const openDistribute = async () => {
	groupList.value = [];
	distributeFormRef.value?.setForm("group_id", null);
	kfList.value = [];
	distributeFormRef.value?.setForm("services_id", null);

	projectList.value = await service.customer_pro.project.list();
	const ids = Crud.value?.selection.map((e) => {
		return e.id;
	});

	// const item: any = [];
	distributeFormRef.value?.open({
		title: `分配`,
		items: [
			{
				label: "项目",
				prop: "project_id",
				component: {
					name: "slot-project_id"
				},
				required: true
			},
			{
				label: "客服组",
				prop: "group_id",
				component: {
					name: "slot-group_id"
				},
				required: true
			},
			{
				label: "接收人",
				prop: "services_id",
				component: {
					name: "slot-services_id"
				},
				required: true
			}
		],

		on: {
			async open() {},
			submit(data, { done, close }) {
				service.customer_pro.clues
					.distribute({
						ids: ids,
						servicesId: data.services_id
					})
					.then(() => {
						ElMessage.success("分配完成");
						close();
						refresh();
					})
					.catch((err) => {
						ElMessage.error(err.message);
						done();
					});
			}
		}
	});
};

// 项目id
const projectId = ref();
// 项目id改变
const projectChange = (v: any) => {
	groupList.value = [];
	distributeFormRef.value?.setForm("group_id", null);
	kfList.value = [];
	distributeFormRef.value?.setForm("services_id", null);
	projectId.value = v;
	getGroupList(v);
};

// 组别id改变
const groupChange = (v: any) => {
	kfList.value = [];
	distributeFormRef.value?.setForm("services_id", null);
	getKfList(v, projectId.value);
};

// 客服组列表
const groupList = ref();
const getGroupList = async (projectId: string) => {
	groupList.value = await service.customer_pro.project_group.list({ projectId });
};

// 客服人员列表
const kfList = ref();
const getKfList = async (groupId: string, projectId: string) => {
	kfList.value = await service.customer_pro.kf.list({ groupId, projectId });
};

// 学校列表
const schoolList = ref();
const majorsList = ref();
const getSchoolList = async () => {
	schoolList.value = await service.customer_pro.school.list();
	if (Upsert.value?.mode != "add") {
		getMajorList(schoolList.value[0].id);
	}
};

// 学校改变
const schoolChange = async (v: any) => {
	majorsList.value = [];
	OrderFormRef.value?.setForm("majors_id", null);
	getMajorList(v);
};

// 专业列表
const getMajorList = async (v: any) => {
	majorsList.value = await service.customer_pro.majors.list({ schoolId: v });
};

// 时间选择器起始
const defaultTime = new Date();
// 高级搜索
const AdvSearch = useAdvSearch({
	items: [
		{
			label: "线索等级",
			prop: "levelStatus",
			component: {
				name: "el-select",
				props: {
					clearable: true
				},
				options: dict.get("cluesLevel")
			}
		},
		() => {
			return {
				label: "来源",
				prop: "sourceStatus",
				component: {
					name: "el-select",
					props: {
						clearable: true
					},
					options: () => {
						let option = [
							{
								label: "手动录入",
								value: "1"
							},
							{
								label: "百度",
								value: "2"
							},
							{
								label: "抖音",
								value: "3"
							},
							{
								label: "53客服",
								value: "4"
							},
							{
								label: "小红书",
								value: "5"
							}
						];

						if (isAdmin.value) {
							option.push({
								label: "测试数据",
								value: "-1"
							});
						}

						return option;
					}
				}
			};
		},
		() => {
			return {
				label: "客服组",
				prop: "serviceGroup",
				component: {
					name: "el-select",
					props: {
						clearable: true
					},
					options: serviceGroup.value
				}
			};
		},
		{
			label: "客服分配状态",
			prop: "serviceStatus",
			component: {
				name: "el-select",
				props: {
					clearable: true
				},
				options: [
					{ label: "未分配", value: 1 },
					{ label: "已分配", value: 2 }
				]
			}
		},
		{
			label: "跟进状态",
			prop: "followType",
			component: {
				name: "el-select",
				props: {
					clearable: true
				},
				options: [
					{ label: "待跟进", value: 1 },
					{ label: "电话访谈", value: 2 },
					{ label: "微信沟通", value: 3 },
					{ label: "视频参观", value: 4 },
					{ label: "预约参观", value: 5 },
					{ label: "已参观", value: 6 }
				]
			}
		},
		{
			label: "创建时间",
			prop: "datetimerange",
			component: {
				name: "el-date-picker",
				props: {
					type: "datetimerange",
					startPlaceholder: "开始日期",
					endPlaceholder: "结束日期",
					defaultTime: defaultTime,
					value: "YYYY-MM-DD HH:mm",
					valueFormat: "YYYY-MM-DD HH:mm",
					timeFormat: "HH:mm"
				}
			}
		},
		{
			label: "更新时间",
			prop: "updatetimerange",
			component: {
				name: "el-date-picker",
				props: {
					type: "datetimerange",
					startPlaceholder: "开始日期",
					endPlaceholder: "结束日期",
					defaultTime: defaultTime,
					value: "YYYY-MM-DD HH:mm",
					valueFormat: "YYYY-MM-DD HH:mm",
					timeFormat: "HH:mm"
				}
			}
		}
	],
	op: ["reset", "close", "search"],
	onSearch(data, { next, close }) {
		next(data);
		searchStatus.value = false;
		searchStatus.value = Object.values(data).some((value) => {
			if (value) return true;
		});
	}
});

const userInfo = ref();
// 是否是管理员
const isAdmin = ref(false);
const getUserInfo = async () => {
	userInfo.value = await service.customer_pro.comm.person();
	isAdmin.value = user.info.roleIds?.split(",").includes("1");
};

// 来源
const sourceFormatter = (v: string) => {
	switch (v) {
		case "1":
			return "手动录入";
		case "2":
			return "百度";
		case "3":
			return "抖音";
		case "4":
			return "53客服";
		case "5":
			return "小红书";
		default:
			break;
	}
};

// 线索等级
const levelFormatter = (v: string) => {
	switch (v) {
		case "0":
			return "无意向";
		case "1":
			return "无效";
		case "2":
			return "有意向";

		default:
			break;
	}
};

// 客服组
const getServiceGroup = async () => {
	const list = await service.customer_pro.project_group.list();
	serviceGroup.value = list.map((item) => {
		return {
			label: item.name,
			value: item.id
		};
	});
};

// table行颜色
const tableRowClassName = () => {
	return "rowColor";
};

// 打开弹窗
const openExcelDialog = () => {
	openExcel.value = true;
};

const fullscreenLoading = ref(false);
const excelDownRef = ref();
// 导出excel
const toexcel = async (isCurrent: boolean) => {
	const loading = ElLoading.service({
		lock: true,
		text: "Loading",
		background: "rgba(0, 0, 0, 0.7)"
	});
	const result = await service.customer_pro.clues.excel({
		...searchData.value,
		keyWord: keyWord.value,
		isCurrentPage: isCurrent
	});

	loading.close();

	if (result.status == 1) {
		ElMessage.warning(result.msg);
	} else {
		ElMessage.success("导出任务已启动，请等待完成");
		excelDownRef.value.refresh();
	}

	return;
};
watch(
	() => openExcel.value,
	(val) => {
		if (!val) {
			console.log("停止任务");
			excelDownRef.value.stopInterval();
		}
	}
);
onMounted(async () => {
	await getUserInfo();
	await getServiceGroup();
});
</script>

<style lang="scss" scoped>
.dialog-footer {
	display: inline-flex;
	flex-direction: row;
	gap: 10px;
}
// 表格表头的背景色;
:deep(.el-table .rowColor) {
	background: #f0f0f1;
}

:deep(.el-pagination) {
	width: 100%;
	justify-content: end;
}
.el-button + .el-button {
	margin-left: 0;
}

.exportBox {
	display: flex;
	flex-direction: row;
	justify-content: center;
	align-items: center;
}
.excel-table {
	height: 600px;
}
</style>
