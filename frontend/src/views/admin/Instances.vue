<template>
  <SectionBanner label="实例管理" icon="pi pi-desktop text-lime-500">
    <Button severity="contrast" :icon="refreshIcon" label="刷新" rounded @click="getInstances" />
  </SectionBanner>

  <div class="card mt-7 rounded-2xl">
    <DataTable :value="data" pt:header:class="font-black">
      <template #empty>
        <div class="text-center mt-2 mb-2 text-secondary">
          未创建实例
        </div>
      </template>
      <Column field="id" frozen class="font-bold"></Column>
      <Column header="主机名称 / 用户" class="min-w-40">
        <template #body="{ data }">
          <span class="text-primary font-bold" @click="getServerDetailEvent($event, data.server_id)"
            @mouseenter="getServerDetailEvent($event, data.server_id)">{{ data.server_name }}</span>
          <br />
          <span class="font-bold">{{ data.username }}</span>
          <br />
          <EditLabel v-model="data.label" class="mt-1" emptyText="设置备注"
            :save="(label) => instanceModifyLabel(data.id, label)" />
        </template>
      </Column>
      <Column header="状态" class="min-w-28 font-bold">
        <template #body="{ data }">
          <Tag v-if="data.status == statusRunning" icon="pi pi-play" value="运行中" />
          <Tag v-else-if="data.status == statusPaused" severity="warn" icon="pi pi-pause" value="暂停" />
          <Tag v-else-if="data.status == statusStoped" severity="secondary" icon="pi pi-stop-circle" value="已停止" />

          <Tag v-else-if="data.status == statusReady" severity="success" icon="pi pi-spin pi-spinner" value="准备中" />
          <Tag v-else-if="data.status == statusStarting" severity="info" icon="pi pi-spin pi-spinner" value="启动中" />
          <Tag v-else-if="data.status == statusStopping" severity="info" icon="pi pi-spin pi-spinner" value="停止中" />
          <Tag v-else-if="data.status == statusPausing" severity="info" icon="pi pi-spin pi-spinner" value="暂停中" />
          <Tag v-else-if="data.status == statusRestarting" severity="info" icon="pi pi-spin pi-spinner" value="重启中" />
          <Tag v-else-if="data.status == statusModifying" severity="success" icon="pi pi-spin pi-spinner" value="调整中" />
          <Tag v-else-if="data.status == statusDeleting" severity="danger" icon="pi pi-spin pi-spinner" value="删除中" />

          <Tag v-else-if="data.status == statusFail" severity="danger" icon="pi pi-exclamation-triangle" value="错误" />
        </template>
      </Column>
      <Column header="规格详情" class="min-w-36 font-bold">
        <template #body="{ data }">
          <span v-if="data.gpu_count !== 0">{{ data.gpu_type }} * {{ data.gpu_count }}</span>
          <span v-else>无卡模式</span>
          <br />
          <span class="text-primary" @click="getInstanceDetail($event, data)"
            @mouseenter="getInstanceDetail($event, data)">查看详情</span>
        </template>
      </Column>
      <Column header="SSH 信息">
        <template #body="{ data }">
          <span class="text-surface-900 dark:text-surface-0 font-medium mr-2 mb-1 md:mb-0">地址:</span>
          <div v-if="data.status === statusRunning" class="mt-1 text-muted-color">{{ data.ssh_address }}
            <CopyIcon class="ml-1" :text="data.ssh_address" />
          </div>
          <div v-else class="mt-1 text-muted-color">暂无</div>
          <span class="text-surface-900 dark:text-surface-0 font-medium mr-2 mb-1 md:mb-0">密码:</span>
          <div v-if="data.status === statusRunning" class="mt-1 text-muted-color">********
            <CopyIcon class="ml-1" :text="data.ssh_passwd" />
          </div>
          <div v-else class="mt-1 text-muted-color">暂无</div>
        </template>
      </Column>
      <Column field="create_at" header="创建时间"></Column>
      <Column>
        <template #body="{ data }">
          <div class="flex gap-2">
            <Button v-if="data.status == statusRunning" icon="pi pi-code" aria-label="Filter" as="a"
              :href="'http://' + data.code_server_address" target="_blank" v-tooltip.top="'VSCode Web'" />
            <Button v-else icon="pi pi-code" aria-label="Filter" v-tooltip.top="'VSCode Web'" disabled />
            <Button v-if="data.status == statusRunning" severity="info" icon="pi pi-inbox" aria-label="Filter" as="a"
              :href="'http://' + data.jupyter_address + '/lab'" target="_blank" v-tooltip.top="'Jupyter Lab'" />
            <Button v-else severity="info" icon="pi pi-inbox" aria-label="Filter" v-tooltip.top="'Jupyter Lab'"
              disabled />
            <Button v-if="data.status == statusRunning" severity="contrast" icon="pi pi-chart-bar" as="a"
              :href="'http://' + data.grafana_address + '/public-dashboards/2c510f203876465ba76617510ce3e219'"
              target="_blank" v-tooltip.top="'监控'" />
            <Button v-else severity="contrast" icon="pi pi-chart-bar" v-tooltip.top="'监控'" disabled />
            <Button v-if="!isAdmin" icon="pi pi-ellipsis-h" severity="secondary" aria-label="Bookmark"
              @click="showMenu($event, data)" />
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

    <Menu ref="instanceMenu" :model="instanceMenuItems" :popup="true" />

    <Popover ref="serverDetailRef">
      <Fluid v-if="serverDetail" class="rounded flex flex-col gap-4 w-96">
        <Fieldset legend="主机名称">
          <span>{{ serverDetail.name }}</span>
        </Fieldset>
        <Fieldset legend="GPU空闲/总量">
          <span>{{ serverDetail.gpu_num - serverDetail.gpu_used }} / {{ serverDetail.gpu_num }}</span>
        </Fieldset>
        <Fieldset legend="数据盘可扩容">
          <span>{{ serverDetail.volume_total - serverDetail.volume_used }} GB</span>
        </Fieldset>
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

    <Popover ref="instanceDetailRef">
      <Fluid v-if="instanceDetail" class="rounded flex flex-col gap-4 w-full">
        <Fieldset legend="镜像">
          <span>{{ imagesValueMap[instanceDetail.image_name] || instanceDetail.image_name }}</span>
        </Fieldset>
        <Fieldset legend="GPU">
          <span v-if="instanceDetail.gpu_count !== 0">{{ instanceDetail.gpu_type }} * {{ instanceDetail.gpu_count
          }}</span>
          <span v-else>无卡模式</span>
        </Fieldset>
        <div class="flex flex-col md:flex-row gap-4">
          <Fieldset class="flex flex-wrap gap-2 w-full" legend="CPU">
            <span v-if="instanceDetail.gpu_count !== 0">{{ instanceDetail.cpu_count_per_gpu * instanceDetail.gpu_count
            }}
              核</span>
            <span v-else>1 核</span>
          </Fieldset>
          <Fieldset class="flex flex-wrap gap-2 w-full" legend="内存">
            <span v-if="instanceDetail.gpu_count !== 0">{{ instanceDetail.memory_per_gpu * instanceDetail.gpu_count }}
              GB</span>
            <span v-else>2 GB</span>
          </Fieldset>
        </div>

        <div class="flex flex-col md:flex-row gap-4">
          <Fieldset class="flex flex-wrap gap-2 w-full" legend="系统盘">
            <div class="flex gap-2">
              <span>30G</span>
              <i class="pi pi-info-circle mt-1" v-tooltip.top="'本地盘,快速'"></i>
            </div>
          </Fieldset>
          <Fieldset class="flex flex-wrap gap-2 w-full" legend="数据盘">
            <div class="flex gap-2">
              <span> 50 + {{ instanceDetail.volume_size - 50 }} GB</span>
              <i class="pi pi-info-circle mt-1" v-tooltip.top="'本地盘,快速,免费 50G,可扩容/缩容'"></i>
            </div>
          </Fieldset>
        </div>
      </Fluid>
    </Popover>
  </div>

  <Drawer v-model:visible="instanceModifyVisible" header="修改实例配置" position="right" :dismissable="false"
    :showCloseIcon="false" class="!w-96">
    <Fluid class="flex flex-col gap-4 w-full">
      <div class="flex flex-col gap-2">
        <label>实例 ID</label>
        <InputText v-model="instanceDetail.id" type="text" disabled />
      </div>
      <div class="flex flex-col gap-2">
        <label>主机名称</label>
        <InputText v-model="instanceDetail.server_name" type="text" disabled />
      </div>
      <div class="flex flex-col gap-2">
        <label>GPU类型</label>
        <InputText v-model="serverDetail.gpu_type" type="text" disabled />
      </div>
      <div class="flex flex-col gap-2">
        <label>CPU 核数</label>
        <InputGroup>
          <InputNumber v-model="instanceConfiguration.cpu_count" disabled />
          <InputGroupAddon>核</InputGroupAddon>
        </InputGroup>
      </div>
      <div class="flex flex-col gap-2">
        <label>内存</label>
        <InputGroup>
          <InputNumber v-model="instanceConfiguration.memory" disabled />
          <InputGroupAddon>GB</InputGroupAddon>
        </InputGroup>
      </div>
      <div class="flex flex-col gap-2">
        <label>无卡模式</label>
        <ToggleSwitch v-model="cpu_only_mode" @change="instanceModifyCpuonlyChange" />
      </div>
      <div v-show="!instanceConfiguration.cpu_only" class="flex flex-col gap-2">
        <label>GPU 数量</label>
        <SelectButton v-model="instanceConfiguration.gpu_count" :options="instanceConfiguration.options"
          optionDisabled="disabled" optionLabel="label" optionValue="value" aria-labelledby="basic" :allowEmpty="false"
          @change="instanceModifyGpuChange" />
      </div>
      <div class="flex flex-col gap-2">
        <label>数据盘</label>
        <InputNumber v-model="instanceConfiguration.volume_size" showButtons buttonLayout="horizontal" :step="1"
          suffix=" GB" :min="50"
          :max="serverDetail.volume_total - serverDetail.volume_used + instanceDetail.volume_size">
          <template #incrementbuttonicon>
            <span class="pi pi-plus" />
          </template>
          <template #decrementbuttonicon>
            <span class="pi pi-minus" />
          </template>
        </InputNumber>
      </div>
      <div class="flex flex-col gap-2">
        <label>GPU 驱动版本</label>
        <InputText v-model="serverDetail.gpu_driver_version" type="text" disabled />
      </div>
      <div class="flex flex-col gap-2">
        <label>CUDA 版本</label>
        <InputText v-model="serverDetail.gpu_cuda_version" type="text" disabled />
      </div>
    </Fluid>
    <template #footer>
      <div class="flex items-center gap-4 w-52 ml-auto">
        <Button label="取消" class="flex-auto" severity="secondary" @click="instanceModifyVisible = false"></Button>
        <Button label="确定" class="flex-auto" @click="instanceModify"></Button>
      </div>
    </template>
  </Drawer>

  <ConfirmDialog></ConfirmDialog>
</template>

<script setup>
import CopyIcon from '@/components/CopyIcon.vue';
import EditLabel from '@/components/EditLabel.vue';
import SectionBanner from '@/components/SectionBanner.vue';

import api from '@/api';
import { useProfileStore } from '@/stores/profile';
import { formatDate } from '@/utils/time.js';
import { useConfirm } from "primevue/useconfirm";
import { useToast } from 'primevue/usetoast';
import { nextTick, onMounted, ref } from 'vue';

const statusFail = ref(-1)
const statusRunning = ref(0)
const statusPaused = ref(1)
const statusStoped = ref(2)
const statusReady = ref(3)
const statusStarting = ref(4)
const statusStopping = ref(5)
const statusPausing = ref(6)
const statusRestarting = ref(7)
const statusModifying = ref(8)
const statusDeleting = ref(9)

const statusIng = [statusReady.value, statusStarting.value, statusStopping.value, statusPausing.value, statusRestarting.value, statusModifying.value, statusDeleting.value]

const toast = useToast()
const confirm = useConfirm()
const profileStore = useProfileStore();

const isAdmin = ref(profileStore.isAdmin)

const instanceMenu = ref(null)
const serverDetailRef = ref(null)
const serverDetail = ref(null)
const instanceDetailRef = ref(null)
const instanceDetail = ref(null)

const offset = ref(0)
const limit = ref(10)

const data = ref([])
const total = ref(0)
const imagesValueMap = ref({})

const instanceModifyVisible = ref(false)
const instanceConfiguration = ref({})
const cpu_only_mode = ref(false)

const refreshIcon = ref('pi pi-refresh')

const instanceMenuItemsTemplate = ref([
  { label: '无卡模式开机', icon: 'pi pi-power-off !text-emerald-600', command: () => { cpuOnlyMode(instanceDetail.value.id) } },
  { label: '开机', icon: 'pi pi-play !text-blue-600', command: () => { instanceStart(instanceDetail.value.id) } },
  { label: '关机', icon: 'pi pi-stop', command: () => { instanceStop(instanceDetail.value.id) } },
  { label: '暂停', icon: 'pi pi-pause !text-amber-600', command: () => { instancePause(instanceDetail.value.id) } },
  { label: '重启实例', icon: 'pi pi-refresh !text-sky-500', command: () => { instanceRestart(instanceDetail.value.id) } },
  { label: '调整配置', icon: 'pi pi-sliders-h !text-indigo-500', command: () => { openInstanceModify() } },
  { label: '删除实例', icon: 'pi pi-trash !text-red-500', command: () => { openInstanceDelete() } },
  { label: '强制删除', icon: 'pi pi-exclamation-triangle !text-red-500', command: () => { openInstanceForceDelete() } }
])
const instanceMenuItems = ref([])

const changePage = async (event) => {
  limit.value = event.rows
  offset.value = event.first
  await getInstances()
}

const getInstances = async () => {
  // data.value = []
  refreshIcon.value = 'pi pi-spin pi-refresh'
  await api.AdminInstancesList({ offset: offset.value, limit: limit.value }).then(res => {
    data.value = res.data.data.result
    let refresh = false
    for (let i = 0; i < data.value.length; i++) {
      if (statusIng.indexOf(data.value[i].status) !== -1) {
        refresh = true
      }
      data.value[i].create_at = formatDate(data.value[i].create_at)
    }
    total.value = res.data.data.total || 0
    refreshIcon.value = 'pi pi-refresh'
    if (refresh) {
      setTimeout(() => {
        getInstances()
      }, 3000);
    }
  }).catch(err => {
    console.error(err)
    toast.add({ severity: 'error', summary: '获取实例列表失败', detail: err.response.data.msg, life: 3000 });
    refreshIcon.value = 'pi pi-refresh'
  })
}

const getImages = async () => {
  await api.UserImages().then((res) => {
    imagesValueMap.value = Object.entries(res.data.data.result).reduce((acc, [key, value]) => {
      acc[value] = key
      return acc
    }, {});
  }).catch((error) => {
    Message.error('获取镜像列表失败')
  })
}

const getInstanceDetail = (event, instance) => {
  instanceDetailRef.value.hide()

  if (instanceDetail.value?.id === instance.id) {
    instanceDetail.value = null
  } else {
    instanceDetail.value = instance

    nextTick(() => {
      instanceDetailRef.value.show(event)
    });
  }
}

const getServerDetail = async (id) => {
  await api.UserServerDetail(id).then(res => {
    serverDetail.value = res.data.data.result
  }).catch(err => {
    console.error(err)
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

const showMenu = (event, instance) => {
  instanceDetail.value = instance
  instanceMenuItems.value = []
  instanceMenuItemsTemplate.value.forEach(item => {
    let newItem = { ...item }
    switch (item.label) {
      case '无卡模式开机':
        if (instanceDetail.value.cpu_only === true || instanceDetail.value.status !== statusStoped.value || instanceDetail.value.status === statusFail.value) {
          newItem.disabled = true
        }
        break
      case '开机':
        if (instanceDetail.value.status === statusRunning.value || statusIng.indexOf(instanceDetail.value.status) !== -1 || instanceDetail.value.status === statusFail.value) {
          newItem.disabled = true
        }
        break
      case '关机':
        if (instanceDetail.value.status === statusStoped.value || statusIng.indexOf(instanceDetail.value.status) !== -1 || instanceDetail.value.status === statusFail.value) {
          newItem.disabled = true
        }
        break
      case '暂停':
        if (instanceDetail.value.status === statusPaused.value || instanceDetail.value.status === statusStoped.value || statusIng.indexOf(instanceDetail.value.status) !== -1 || instanceDetail.value.status === statusFail.value) {
          newItem.disabled = true
        }
        break
      case '重启实例':
        if (statusIng.indexOf(instanceDetail.value.status) !== -1 || instanceDetail.value.status === statusFail.value) {
          newItem.disabled = true
        }
        break
      case '调整配置':
        if (instanceDetail.value.status !== statusStoped.value || instanceDetail.value.status === statusFail.value) {
          newItem.disabled = true
        }
        break
      case '删除实例':
        if (instanceDetail.value.status === statusDeleting.value || instanceDetail.value.status === statusFail.value) {
          newItem.disabled = true
        }
        break
      case '强制删除':
        if (instanceDetail.value.status !== statusFail.value) {
          newItem.disabled = true
        }
        break
    }
    instanceMenuItems.value.push(newItem)
  })

  instanceMenu.value.show(event)
}

const cpuOnlyMode = async (id) => {
  toast.add({ severity: 'info', summary: '切换为无卡模式', detail: '正在切换为无卡模式', life: 3000 });
  setTimeout(() => {
    getInstances()
  }, 100);
  await api.AdminInstancesModify(id, { cpu_only: true }).then(async (res) => {
    toast.add({ severity: 'success', summary: '切换为无卡模式', detail: '已切换为无卡模式', life: 3000 });
    await getInstances()
    // await instanceStart(id)
  }).catch(err => {
    console.error(err)
    toast.add({ severity: 'error', summary: '切换为无卡模式失败', detail: err.response.data.msg, life: 3000 });
  })
}

const instanceStart = async (id) => {
  toast.add({ severity: 'info', summary: '开机', detail: '正在开机', life: 3000 });
  setTimeout(() => {
    getInstances()
  }, 100);
  await api.AdminInstancesAction(id, { action: 1 }).then(async (res) => {
    toast.add({ severity: 'success', summary: '开机', detail: '实例已开机', life: 3000 });
    await getInstances()
  }).catch(err => {
    console.error(err)
    toast.add({ severity: 'error', summary: '开机失败', detail: err.response.data.msg, life: 3000 });
  })
}

const instancePause = async (id) => {
  toast.add({ severity: 'info', summary: '暂停实例', detail: '正在暂停实例', life: 3000 });
  setTimeout(() => {
    getInstances()
  }, 100);
  await api.AdminInstancesAction(id, { action: 2 }).then(async (res) => {
    toast.add({ severity: 'success', summary: '暂停实例', detail: '实例已暂停', life: 3000 });
    await getInstances()
  }).catch(err => {
    console.error(err)
    toast.add({ severity: 'error', summary: '暂停实例失败', detail: err.response.data.msg, life: 3000 });
  })
}

const instanceStop = async (id) => {
  toast.add({ severity: 'info', summary: '停止实例', detail: '正在停止实例', life: 3000 });
  setTimeout(() => {
    getInstances()
  }, 100);
  await api.AdminInstancesAction(id, { action: 3 }).then(async (res) => {
    toast.add({ severity: 'success', summary: '停止实例', detail: '实例已停止', life: 3000 });
    await getInstances()
  }).catch(err => {
    console.error(err)
    toast.add({ severity: 'error', summary: '停止实例失败', detail: err.response.data.msg, life: 3000 });
  })
}

const instanceRestart = async (id) => {
  toast.add({ severity: 'info', summary: '重启实例', detail: '正在重启实例', life: 3000 });
  setTimeout(() => {
    getInstances()
  }, 100);
  await api.AdminInstancesAction(id, { action: 4 }).then(async (res) => {
    toast.add({ severity: 'success', summary: '重启实例', detail: '实例已重启', life: 3000 });
    await getInstances()
  }).catch(err => {
    console.error(err)
    toast.add({ severity: 'error', summary: '重启实例失败', detail: err.response.data.msg, life: 3000 });
  })
}

const instanceModify = async () => {
  toast.add({ severity: 'info', summary: '调整配置', detail: '正在调整配置', life: 3000 });
  setTimeout(() => {
    getInstances()
  }, 100);
  await api.AdminInstancesModify(instanceDetail.value.id, instanceConfiguration.value).then(async (res) => {
    toast.add({ severity: 'success', summary: '调整配置', detail: '已调整配置', life: 3000 });
    instanceModifyVisible.value = false
    await getInstances()
  }).catch(err => {
    console.error(err)
    toast.add({ severity: 'error', summary: '调整配置失败', detail: err.response.data.msg, life: 3000 });
  })
}

const instanceDelete = async (id) => {
  toast.add({ severity: 'info', summary: '释放实例', detail: '正在释放实例', life: 3000 });
  setTimeout(() => {
    getInstances()
  }, 100);
  await api.AdminInstancesDelete(id).then(async (res) => {
    toast.add({ severity: 'success', summary: '释放实例', detail: '实例已释放', life: 3000 });
    await getInstances()
  }).catch(err => {
    console.error(err)
    toast.add({ severity: 'error', summary: '释放实例失败', detail: err.response.data.msg, life: 3000 });
  })
}

const instanceForceDelete = async (id) => {
  toast.add({ severity: 'info', summary: '强制释放实例', detail: '正在强制释放实例', life: 3000 });
  setTimeout(() => {
    getInstances()
  }, 100);
  await api.AdminInstancesForceDelete(id, { force: true }).then(async (res) => {
    toast.add({ severity: 'success', summary: '强制释放实例', detail: '实例已强制释放', life: 3000 });
    await getInstances()
  }).catch(err => {
    console.error(err)
    toast.add({ severity: 'error', summary: '强制释放实例失败', detail: err.response.data.msg, life: 3000 });
  })
}

const openInstanceModify = async () => {
  instanceConfiguration.value.gpu_count = instanceDetail.value.gpu_count
  instanceConfiguration.value.volume_size = instanceDetail.value.volume_size
  instanceConfiguration.value.cpu_only = instanceDetail.value.cpu_only
  await getServerDetail(instanceDetail.value.server_id)
  instanceConfiguration.value.options = []
  for (let i = 0; i < serverDetail.value.gpu_num; i++) {
    const tmp = {
      label: i + 1,
      value: i + 1,
      disabled: i + 1 > serverDetail.value.gpu_num - serverDetail.value.gpu_used
    }
    instanceConfiguration.value.options.push(tmp)
  }
  if (serverDetail.value.gpu_num - serverDetail.value.gpu_used < instanceConfiguration.value.gpu_count) {
    instanceConfiguration.value.gpu_count = serverDetail.value.gpu_num - serverDetail.value.gpu_used
  }
  instanceConfiguration.value.cpu_count = serverDetail.value.cpu_count_per_gpu * instanceConfiguration.value.gpu_count
  instanceConfiguration.value.memory = serverDetail.value.memory_per_gpu * instanceConfiguration.value.gpu_count
  if (serverDetail.value.gpu_num === serverDetail.value.gpu_used || instanceDetail.value.cpu_only) {
    instanceConfiguration.value.gpu_count = 0
    instanceConfiguration.value.cpu_only = true
    cpu_only_mode.value = true
    instanceConfiguration.value.cpu_count = 1
    instanceConfiguration.value.memory = 2
  }
  instanceModifyVisible.value = true
}

const instanceModifyGpuChange = (event) => {
  instanceConfiguration.value.cpu_count = serverDetail.value.cpu_count_per_gpu * event.value
  instanceConfiguration.value.memory = serverDetail.value.memory_per_gpu * event.value
}

const instanceModifyCpuonlyChange = () => {
  if (!cpu_only_mode.value && serverDetail.value.gpu_num === serverDetail.value.gpu_used) {
    toast.add({ severity: 'error', summary: '无法调整显卡资源', detail: '宿主机显卡资源不足', life: 3000 });
    cpu_only_mode.value = true
  }
  if (cpu_only_mode.value) {
    instanceConfiguration.value.cpu_only = true
    instanceConfiguration.value.cpu_count = 1
    instanceConfiguration.value.memory = 2
  } else {
    instanceConfiguration.value.cpu_only = false
    instanceConfiguration.value.gpu_count = instanceDetail.value.gpu_count
    if (instanceConfiguration.value.gpu_count === 0) {
      instanceConfiguration.value.gpu_count = 1
    }
    if (serverDetail.value.gpu_num - serverDetail.value.gpu_used < instanceConfiguration.value.gpu_count) {
      instanceConfiguration.value.gpu_count = serverDetail.value.gpu_num - serverDetail.value.gpu_used
    }
    instanceConfiguration.value.cpu_count = serverDetail.value.cpu_count_per_gpu * instanceConfiguration.value.gpu_count
    instanceConfiguration.value.memory = serverDetail.value.memory_per_gpu * instanceConfiguration.value.gpu_count
  }
}

const openInstanceDelete = () => {
  confirm.require({
    header: '确认删除 实例ID: ' + instanceDetail.value.id,
    message: '实例删除后，数据将无法恢复，请确认删除',
    icon: 'pi pi-info-circle',
    rejectProps: {
      label: '取消',
      severity: 'secondary',
      outlined: true
    },
    acceptProps: {
      label: '删除',
      severity: 'danger'
    },
    accept: async () => {
      await instanceDelete(instanceDetail.value.id)
    },
    reject: () => { }
  });
}

const openInstanceForceDelete = () => {
  confirm.require({
    header: '确认强制删除 实例ID: ' + instanceDetail.value.id,
    message: '强制实例删除后，数据将无法恢复，请确认删除',
    icon: 'pi pi-info-circle',
    rejectProps: {
      label: '取消',
      severity: 'secondary',
      outlined: true
    },
    acceptProps: {
      label: '强制删除',
      severity: 'danger'
    },
    accept: async () => {
      await instanceForceDelete(instanceDetail.value.id)
    },
    reject: () => { }
  });
}

const instanceModifyLabel = (id, label) => {
  api.AdminInstancesModifyLabel(id, { label: label }).then(async (res) => {
    toast.add({ severity: 'success', summary: '修改备注成功', detail: '已保存备注', life: 3000 });
    await getInstances()
  }).catch(err => {
    console.error(err)
    toast.add({ severity: 'error', summary: '修改备注失败', detail: err.response.data.msg, life: 3000 });
  })
}

onMounted(() => {
  getInstances()
  getImages()
})

</script>