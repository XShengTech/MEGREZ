<template>
  <SectionBanner label="节点管理" icon="pi pi-server text-yellow-400">
    <Button icon="pi pi-plus" label="添加节点" rounded @click="openServerAdd" />
    <Button severity="contrast" :icon="refreshIcon" label="刷新" rounded @click="getServers" />
  </SectionBanner>

  <div class="card mt-7 rounded-2xl">
    <DataTable :value="data" pt:header:class="font-black">
      <template #empty>
        <div class="text-center mt-2 mb-2 text-secondary">
          未添加节点
        </div>
      </template>
      <Column field="id" frozen class="font-bold"></Column>
      <Column header="主机 / IP" class="min-w-32">
        <template #body="{ data }">
          <span class="text-primary font-bold">{{ data.name }}</span>
          <br />
          <span class="text-muted-color">{{ data.ip }}</span>
        </template>
      </Column>
      <Column header="规格" class="min-w-64 font-bold">
        <template #body="{ data }">
          <span class="mt-2 mb-2 block">显卡:
            <Tag :value="data.gpu_type" class="mr-1" rounded></Tag>
            <Tag v-if="data.gpu_num > data.gpu_used" severity="success"
              :value="'剩余 ' + Number(data.gpu_num - data.gpu_used) + '/' + data.gpu_num + ' 张'" rounded></Tag>
            <Tag v-else severity="danger"
              :value="'剩余 ' + Number(data.gpu_num - data.gpu_used) + '/' + data.gpu_num + ' 张'" rounded></Tag>
          </span>
          <span class="mt-2 mb-2 block">数据盘:
            <Tag v-if="data.volume_total > data.volume_used"
              :value="Number(data.volume_total - data.volume_used) + '/' + data.volume_total + ' GB'" rounded></Tag>
            <Tag v-else severity="danger"
              :value="Number(data.volume_total - data.volume_used) + '/' + data.volume_total + ' GB'" rounded></Tag>
          </span>
        </template>
      </Column>
      <Column header="配置详情" class="min-w-32">
        <template #body="{ data }">
          <div class="flex items-center gap-2 px-1">
            <div class="text-surface-900 dark:text-surface-0 font-medium -ml-1 mb-1 md:mb-0">CPU:</div>
            <span class="text-muted-color">{{ data.cpu_count_per_gpu }} 核/卡</span>
          </div>
          <div class="flex items-center gap-2 px-1">
            <div class="text-surface-900 dark:text-surface-0 font-medium -ml-1 mb-1 md:mb-0">内存:</div>
            <span class="text-muted-color">{{ data.memory_per_gpu }} GB/卡</span>
          </div>
          <span class="text-primary font-bold" @click="getServerDetailEvent($event, data.id)"
            @mouseenter="getServerDetailEvent($event, data.id)">查看详情</span>
        </template>
      </Column>
      <Column field="create_at" header="创建时间"></Column>
      <Column>
        <template #body="{ data }">
          <div class="flex gap-2">
            <Button icon="pi pi-pencil" v-tooltip.top="'编辑'" @click="openServerModify(data)" />
            <Button severity="danger" icon="pi pi-trash" v-tooltip.top="'删除'" @click="openServerDelete(data)" />
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

    <Popover ref="serverDetailRef">
      <Fluid v-if="serverDetail" class="rounded flex flex-col gap-4 w-96">
        <Fieldset legend="GPU 空闲/总量">
          <span>{{ serverDetail.gpu_num - serverDetail.gpu_used }} / {{ serverDetail.gpu_num }} 张</span>
        </Fieldset>
        <Fieldset legend="数据盘 空闲/总量">
          <span>{{ serverDetail.volume_total - serverDetail.volume_used }} / {{ serverDetail.volume_total }} GB</span>
        </Fieldset>
        <div class="flex flex-col md:flex-row gap-4">
          <Fieldset class="flex flex-wrap gap-2 w-full" legend="CPU">
            <span>{{ serverDetail.cpu_count_per_gpu }} 核/卡</span>
          </Fieldset>
          <Fieldset class="flex flex-wrap gap-2 w-full" legend="内存">
            <span>{{ serverDetail.memory_per_gpu }} GB/卡</span>
          </Fieldset>
        </div>
        <div class="flex flex-col md:flex-row gap-4">
          <Fieldset class="flex flex-wrap gap-2 w-full" legend="GPU驱动版本">
            <span>{{ serverDetail.gpu_driver_version }}</span>
          </Fieldset>
          <Fieldset class="flex flex-wrap gap-2 w-full" legend="CUDA版本">
            <span>≤ {{ serverDetail.gpu_cuda_version }}</span>
          </Fieldset>
        </div>
      </Fluid>
    </Popover>
  </div>

  <Dialog v-model:visible="serverModifyVisible" modal :header="'编辑节点 - ' + serverDetail.name"
    :style="{ width: '42rem' }">
    <span class="text-surface-500 dark:text-surface-400 block mb-6">编辑节点配置</span>
    <div class="flex items-center gap-0 mb-4">
      <label class="font-semibold w-20">名称:</label>
      <InputText v-model="serverDetail.name" class="flex-auto" type="text" />
    </div>
    <div class="grid grid-cols-2 md:flex-row gap-8 mb-4">
      <div class="flex gap-0">
        <label class="font-semibold w-28 mt-2">IP:</label>
        <InputGroup class="flex-auto">
          <InputText v-model="serverDetail.ip" type="text" />
        </InputGroup>
      </div>
      <div class="flex gap-0">
        <label class="font-semibold w-28 mt-2">端口:</label>
        <InputGroup class="flex-auto">
          <InputNumber v-model="serverDetail.port" :useGrouping="false" :min="1" :max="65535" />
        </InputGroup>
      </div>
    </div>
    <div class="grid grid-cols-2 md:flex-row gap-8 mb-4">
      <div class="flex gap-0">
        <label class="font-semibold w-28 mt-2">APIKEY:</label>
        <InputGroup class="flex-auto">
          <InputText v-model="serverDetail.apikey" />
        </InputGroup>
      </div>
      <div class="flex gap-0">
        <label class="font-semibold w-28 mt-2">CPU:</label>
        <InputGroup class="flex-auto">
          <InputNumber v-model="serverDetail.cpu_count_per_gpu" :useGrouping="false" :min="1" />
          <InputGroupAddon>核 / 卡</InputGroupAddon>
        </InputGroup>
      </div>
    </div>
    <div class="grid grid-cols-2 md:flex-row gap-8 mb-4">
      <div class="flex gap-0">
        <label class="font-semibold w-28 mt-2">内存:</label>
        <InputGroup class="flex-auto">
          <InputNumber v-model="serverDetail.memory_per_gpu" :useGrouping="false" :min="1" />
          <InputGroupAddon>GB / 卡</InputGroupAddon>
        </InputGroup>
      </div>
      <div class="flex gap-0">
        <label class="font-semibold w-28 mt-2">数据盘:</label>
        <InputGroup class="flex-auto">
          <InputNumber v-model="serverDetail.volume_total" :useGrouping="false" :min="1" />
          <InputGroupAddon>GB</InputGroupAddon>
        </InputGroup>
      </div>
    </div>
    <div class="grid grid-cols-2 md:flex-row gap-8 mb-4">
      <div class="flex gap-0">
        <label class="font-semibold w-28 mt-2">GPU 类型:</label>
        <InputGroup class="flex-auto">
          <InputText v-model="serverDetail.gpu_type" />
        </InputGroup>
      </div>
      <div class="flex gap-0">
        <label class="font-semibold w-28 mt-2">GPU 数量:</label>
        <InputGroup class="flex-auto">
          <InputNumber v-model="serverDetail.gpu_num" :useGrouping="false" :min="0" />
          <InputGroupAddon>张</InputGroupAddon>
        </InputGroup>
      </div>
    </div>
    <div class="grid grid-cols-2 md:flex-row gap-8 mb-4">
      <div class="flex gap-0">
        <label class="font-semibold w-28 mt-2">驱动版本:</label>
        <InputGroup>
          <InputText v-model="serverDetail.gpu_driver_version" class="flex-auto" type="text" />
        </InputGroup>
      </div>
      <div class="flex gap-0">
        <label class="font-semibold w-28 mt-2">CUDA:</label>
        <InputGroup>
          <InputText v-model="serverDetail.gpu_cuda_version" class="flex-auto" />
        </InputGroup>
      </div>
    </div>
    <div class="flex justify-end gap-2">
      <Button type="button" label="取消" severity="secondary" @click="serverModifyVisible = false"></Button>
      <Button type="button" label="保存" @click="serverModify"></Button>
    </div>
  </Dialog>

  <Dialog v-model:visible="serverAddVisible" modal :header="'添加节点'" :style="{ width: '42rem' }">
    <span class="text-surface-500 dark:text-surface-400 block mb-6">添加节点配置</span>
    <div class="flex items-center gap-0 mb-4">
      <label class="font-semibold w-20">名称:</label>
      <InputText v-model="serverDetail.name" class="flex-auto" type="text" />
    </div>
    <div class="flex flex-col md:flex-row gap-8 mb-4">
      <div class="flex gap-0">
        <label class="font-semibold w-28 mt-2">IP:</label>
        <InputGroup class="flex-auto">
          <InputText v-model="serverDetail.ip" type="text" />
        </InputGroup>
      </div>
      <div class="flex gap-0">
        <label class="font-semibold w-28 mt-2">端口:</label>
        <InputGroup class="flex-auto">
          <InputNumber v-model="serverDetail.port" :useGrouping="false" :min="1" :max="65535" />
        </InputGroup>
      </div>
    </div>
    <div class="grid grid-cols-2 md:flex-row gap-8 mb-4">
      <div class="flex gap-0">
        <label class="font-semibold w-28 mt-2">APIKEY:</label>
        <InputGroup class="flex-auto">
          <InputText v-model="serverDetail.apikey" />
        </InputGroup>
      </div>
      <div class="flex gap-0">
        <label class="font-semibold w-28 mt-2">CPU:</label>
        <InputGroup class="flex-auto">
          <InputNumber v-model="serverDetail.cpu_count_per_gpu" :useGrouping="false" :min="1" />
          <InputGroupAddon>核 / 卡</InputGroupAddon>
        </InputGroup>
      </div>
    </div>
    <div class="grid grid-cols-2 md:flex-row gap-8 mb-4">
      <div class="flex gap-0">
        <label class="font-semibold w-28 mt-2">内存:</label>
        <InputGroup class="flex-auto">
          <InputNumber v-model="serverDetail.memory_per_gpu" :useGrouping="false" :min="1" />
          <InputGroupAddon>GB / 卡</InputGroupAddon>
        </InputGroup>
      </div>
      <div class="flex gap-0">
        <label class="font-semibold w-28 mt-2">数据盘:</label>
        <InputGroup class="flex-auto">
          <InputNumber v-model="serverDetail.volume_total" :useGrouping="false" :min="1" />
          <InputGroupAddon>GB</InputGroupAddon>
        </InputGroup>
      </div>
    </div>
    <div class="grid grid-cols-2 md:flex-row gap-8 mb-4">
      <div class="flex gap-0">
        <label class="font-semibold w-28 mt-2">GPU 类型:</label>
        <InputGroup class="flex-auto">
          <InputText v-model="serverDetail.gpu_type" />
        </InputGroup>
      </div>
      <div class="flex gap-0">
        <label class="font-semibold w-28 mt-2">GPU 数量:</label>
        <InputGroup class="flex-auto">
          <InputNumber v-model="serverDetail.gpu_num" :useGrouping="false" :min="0" />
          <InputGroupAddon>张</InputGroupAddon>
        </InputGroup>
      </div>
    </div>
    <div class="grid grid-cols-2 md:flex-row gap-8 mb-4">
      <div class="flex gap-0">
        <label class="font-semibold w-28 mt-2">驱动版本:</label>
        <InputGroup>
          <InputText v-model="serverDetail.gpu_driver_version" class="flex-auto" type="text" />
        </InputGroup>
      </div>
      <div class="flex gap-0">
        <label class="font-semibold w-28 mt-2">CUDA:</label>
        <InputGroup>
          <InputText v-model="serverDetail.gpu_cuda_version" class="flex-auto" />
        </InputGroup>
      </div>
    </div>
    <div class="flex justify-end gap-2">
      <Button type="button" label="取消" severity="secondary" @click="serverAddVisible = false"></Button>
      <Button type="button" label="添加" @click="serverAdd"></Button>
    </div>
  </Dialog>

  <ConfirmDialog></ConfirmDialog>
</template>

<script setup>
import { nextTick, onMounted, ref } from 'vue';

import api from '@/api';
import { formatDate } from '@/utils/time.js';
import { useConfirm } from "primevue/useconfirm";
import { useToast } from 'primevue/usetoast';

const toast = useToast()
const confirm = useConfirm()

const offset = ref(0)
const limit = ref(10)

const data = ref([])
const total = ref(0)

const serverDetailRef = ref(null)
const serverDetail = ref({})
const serverAddVisible = ref(false)
const serverModifyVisible = ref(false)

const refreshIcon = ref('pi pi-refresh')

const changePage = async (event) => {
  limit.value = event.rows
  offset.value = event.first
  await getServers()
}

const getServers = async () => {
  refreshIcon.value = 'pi pi-spin pi-refresh'
  await api.AdminServersList({ offset: offset.value, limit: limit.value }).then(res => {
    data.value = res.data.data.result.map(item => {
      item.create_at = formatDate(item.create_at)
      return item
    })
    total.value = res.data.data.total || 0
  }).catch(err => {
    toast.add({ severity: 'error', summary: '获取节点列表失败', detail: err.response.data.msg, life: 3000 })
    console.error(err)
  }).finally(() => {
    refreshIcon.value = 'pi pi-refresh'
  })
}

const getServerDetail = async (serverId) => {
  await api.AdminServersDetail(serverId).then(res => {
    serverDetail.value = res.data.data.result
  })
}

const getServerDetailEvent = (event, serverId) => {
  serverDetailRef.value.hide()

  if (serverDetail.value?.id === serverId) {
    serverDetail.value = null
  } else {
    nextTick(() => {
      serverDetailRef.value.show(event)
    });
    getServerDetail(serverId)
  }
}

const openServerAdd = () => {
  serverDetail.value = {}
  serverAddVisible.value = true
}

const serverAdd = async () => {
  await api.AdminServersAdd(serverDetail.value).then(res => {
    toast.add({ severity: 'success', summary: '节点添加成功', life: 3000 })
    serverAddVisible.value = false
    getServers()
  }).catch(err => {
    toast.add({ severity: 'error', summary: '节点添加失败', detail: err.response.data.msg, life: 3000 })
  })
}

const openServerModify = (data) => {
  serverDetail.value = data
  serverModifyVisible.value = true
  getServerDetail(data.id)
}

const serverModify = async () => {
  await api.AdminServersModify(serverDetail.value.id, serverDetail.value).then(res => {
    toast.add({ severity: 'success', summary: '节点编辑成功', life: 3000 })
    serverModifyVisible.value = false
    getServers()
  }).catch(err => {
    toast.add({ severity: 'error', summary: '节点编辑失败', detail: err.response.data.msg, life: 3000 })
  })
}

const openServerDelete = (data) => {
  confirm.require({
    header: '删除节点 - ' + data.name,
    message: '确定删除节点 ' + data.name + '(ID: ' + data.id + ') 吗？',
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
      serverDelete(data.id)
    }
  })
}

const serverDelete = async (id) => {
  await api.AdminServersDelete(id).then(res => {
    toast.add({ severity: 'success', summary: '节点删除成功', life: 3000 })
    getServers()
  }).catch(err => {
    toast.add({ severity: 'error', summary: '节点删除失败', detail: err.response.data.msg, life: 3000 })
  })
}

onMounted(() => {
  getServers()
})

</script>