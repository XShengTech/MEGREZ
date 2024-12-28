<template>
  <SectionBanner label="用户管理" icon="pi pi-users text-indigo-500"></SectionBanner>

  <div class="card mt-7 rounded-2xl">
    <DataTable :value="data" pt:header:class="font-black">
      <template #empty>
        <div class="text-center mt-2 mb-2 text-secondary">
          暂无用户
        </div>
      </template>
      <Column field="id" frozen class="font-bold"></Column>
      <Column header="用户名" field="username" class="min-w-28 text-primary font-bold"></Column>
      <Column header="邮箱" field="email" class="min-w-40 font-bold"></Column>
      <Column v-if="!isAdmin" header="权限" class="min-w-28 font-bold">
        <template #body="{ data }">
          <Tag v-if="data.role == 0" severity="secondary" value="受限用户" @click="userVerify(data.id)"
            v-tooltip.top="'点击授权'" />
          <Tag v-else-if="data.role == 1" value="普通用户" />
          <Tag v-else-if="data.role == 2" severity="warn" value="管理员" />
          <Tag v-else-if="data.role == 3" severity="danger" value="超级管理员" />
        </template>
      </Column>
      <Column field="create_at" header="创建时间"></Column>
      <Column v-if="!isAdmin">
        <template #body="{ data }">
          <div class="flex gap-2">
            <Button icon="pi pi-pencil" v-tooltip.top="'编辑'" @click="openUserModify(data)" />
            <Button severity="danger" icon="pi pi-trash" v-tooltip.top="'删除'" @click="openUserDelete(data)" />
          </div>
        </template>
      </Column>
    </DataTable>

    <Paginator class="ml-auto mt-4" :rows="10" :totalRecords="total" :rowsPerPageOptions="[10, 20, 30]"
      @page="changePage">
      <template #start>
        共 {{ total }} 个
      </template>
    </Paginator>
  </div>

  <Dialog v-model:visible="userModifyVisible" modal :header="'编辑用户 - ' + userData.username" :style="{ width: '25rem' }">
    <span class="text-surface-500 dark:text-surface-400 block mb-6">编辑用户信息</span>
    <div class="flex items-center gap-0 mb-4">
      <label class="font-semibold w-20">邮箱:</label>
      <InputText v-model="userData.email" class="flex-auto" type="text" disabled />
    </div>
    <div class="flex items-center gap-0 mb-4">
      <label class="font-semibold w-20">密码:</label>
      <InputText v-model="userData.password" class="flex-auto" type="text" placeholder="留空不修改密码" />
    </div>
    <div class="flex items-center gap-0 mb-4">
      <label class="font-semibold w-20">权限:</label>
      <!-- <InputText v-model="userData.role" class="flex-auto" type="text" /> -->
      <Select v-model="userData.role" class="flex-auto" :options="permissionOptions" optionLabel="label"
        optionValue="value" />
    </div>
    <div class="flex justify-end gap-2">
      <Button type="button" label="取消" severity="secondary" @click="userModifyVisible = false"></Button>
      <Button type="button" label="保存" @click="userModify"></Button>
    </div>
  </Dialog>

  <ConfirmDialog></ConfirmDialog>
</template>

<script setup>
import api from '@/api';
import { useProfileStore } from '@/stores/profile';
import { useConfirm } from "primevue/useconfirm";
import { useToast } from 'primevue/usetoast';
import { onMounted, ref } from 'vue';


const toast = useToast()
const confirm = useConfirm()
const profileStore = useProfileStore();

const isAdmin = ref(profileStore.isAdmin)

const offset = ref(0)
const limit = ref(10)

const data = ref([])
const total = ref(0)

const userModifyVisible = ref(false)
const userData = ref({})

const permissionOptions = ref([
  { label: '受限用户', value: 0 },
  { label: '普通用户', value: 1 },
  { label: '管理员', value: 2 },
  { label: '超级管理员', value: 3 }
])

const formatDate = (dateString) => {
  const date = new Date(dateString);

  // 自定义格式化输出 (YYYY-MM-DD HH:mm:ss)
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0') // 月份从 0 开始
  const day = String(date.getDate()).padStart(2, '0')
  const hours = String(date.getHours()).padStart(2, '0')
  const minutes = String(date.getMinutes()).padStart(2, '0')
  const seconds = String(date.getSeconds()).padStart(2, '0')

  return `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`
}

const changePage = (event) => {
  offset.value = event.first
  limit.value = event.rows
  getUsers()
}

const getUsers = async () => {
  await api.AdminUserList({ offset: offset.value, limit: limit.value }).then(res => {
    data.value = res.data.data.result
    total.value = res.data.data.total
    for (let i = 0; i < data.value.length; i++) {
      data.value[i].create_at = formatDate(data.value[i].create_at)
    }
  }).catch(err => {
    toast.add({ severity: 'error', summary: '获取用户列表失败', detail: err.response.data.msg, life: 3000 })
  })
}

const userVerify = async (id) => {
  await api.AdminUserModify(id, { 'role': 1 }).then(res => {
    toast.add({ severity: 'success', summary: '用户授权成功', life: 3000 })
    getUsers()
  }).catch(err => {
    toast.add({ severity: 'error', summary: '用户授权失败', detail: err.response.data.msg, life: 3000 })
  })
}

const openUserModify = (data) => {
  userData.value = data
  userModifyVisible.value = true
}

const userModify = async () => {
  await api.AdminUserModify(userData.value.id, userData.value).then(res => {
    toast.add({ severity: 'success', summary: '用户编辑成功', life: 3000 })
    userModifyVisible.value = false
    getUsers()
  }).catch(err => {
    toast.add({ severity: 'error', summary: '用户编辑失败', detail: err.response.data.msg, life: 3000 })
  })
}

const openUserDelete = (data) => {
  confirm.require({
    header: '删除用户 - ' + data.username,
    message: '确定删除用户 ' + data.username + '(ID: ' + data.id + ') 吗？',
    icon: 'pi pi-exclamation-triangle',
    rejectProps: {
      label: '取消',
      severity: 'secondary',
      outlined: true
    },
    acceptProps: {
      label: '删除',
      severity: 'danger'
    },
    accept: () => {
      userDelete(data.id)
    },
    reject: () => { }
  })
}

const userDelete = async (id) => {
  await api.AdminUserDelete(id).then(res => {
    toast.add({ severity: 'success', summary: '用户删除成功', life: 3000 })
    getUsers()
  }).catch(err => {
    toast.add({ severity: 'error', summary: '用户删除失败', detail: err.response.data.msg, life: 3000 })
  })
}

onMounted(() => {
  getUsers()
})
</script>