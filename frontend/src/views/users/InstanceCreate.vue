<template>
  <SectionBanner label="创建实例" icon="pi pi-plus text-purple-600">
    <Button severity="contrast" :icon="refreshIcon" label="刷新" rounded @click="getServersList" />
  </SectionBanner>

  <div class="card rounded-2xl">
    <DataView :value="data">
      <template #list="slotProps">
        <div class="flex flex-col">
          <div v-for="(item, index) in slotProps.items" :key="index">
            <div class="flex flex-col sm:flex-row sm:items-center p-6 gap-4"
              :class="{ 'border-t border-surface-200 dark:border-surface-700': index !== 0 }">
              <div class="flex flex-col md:flex-row justify-between md:items-center flex-1 gap-6">
                <div class="flex flex-row md:flex-col justify-between items-start gap-2">
                  <div>
                    <span class="font-medium text-surface-500 dark:text-surface-400 text-sm">{{ item.name }}</span>
                    <div class="text-lg font-medium mt-2">{{ item.gpu_type }}</div>
                  </div>

                  <div v-if="item.gpu_num - item.gpu_used > 0" class="flex felx-col space-x-2">
                    <Tag severity="success" value="可用"></Tag>
                    <Tag :value="'剩余 ' + (item.gpu_num - item.gpu_used) + ' 卡'"></Tag>
                    <Tag severity="info" :value="item.price + ' 计点/小时/GPU'"></Tag>
                  </div>
                  <div v-else class="flex felx-col">
                    <Tag severity="danger" value="售罄"></Tag>
                  </div>
                </div>

                <div class="grid grid-cols-3 w-[480px]">
                  <div class="flex flex-row md:flex-col justify-end items-start gap-2 mt-5">
                    <div>
                      <span class="font-medium text-surface-500 dark:text-surface-400 text-sm">每 GPU 分配</span>
                      <div class="text-sm font-medium mt-2">CPU: {{ item.cpu_count_per_gpu }} 核</div>
                      <div class="text-sm font-medium mt-2">内存: {{ item.memory_per_gpu }} GB</div>
                    </div>
                  </div>

                  <div class="flex flex-row md:flex-col justify-end items-start gap-2 mt-5">
                    <div>
                      <span class="font-medium text-surface-500 dark:text-surface-400 text-sm">硬盘</span>
                      <div class="text-sm font-medium mt-2">数据盘: 50 GB</div>
                      <div class="text-sm font-medium mt-2">
                        可扩容: {{ item.volume_total - item.volume_used }} GB
                      </div>
                    </div>
                  </div>

                  <div class="flex flex-row md:flex-col justify-end items-start gap-2 mt-5">
                    <div>
                      <span class="font-medium text-surface-500 dark:text-surface-400 text-sm">其他</span>
                      <div class="text-sm font-medium mt-2">GPU驱动版本: {{ item.gpu_driver_version }}</div>
                      <div class="text-sm font-medium mt-2">CUDA版本: ≤ {{ item.gpu_cuda_version }}</div>
                    </div>
                  </div>
                </div>

                <div class="flex flex-col md:items-end gap-8">
                  <span class="text-xl font-medium mt-2">
                    空闲GPU:
                    <span class="!font-extrabold text-rose-500">
                      {{ item.gpu_num - item.gpu_used }}
                    </span>
                    / {{ item.gpu_num }}
                  </span>
                  <div class="flex flex-row-reverse md:flex-row gap-2 -mt-2">
                    <Button icon="pi pi-plus" label="创建实例" :disabled="item.disabled" @click="clickCreate($event, item)"
                      class="flex-auto md:flex-initial whitespace-nowrap"></Button>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </template>
    </DataView>

    <Paginator class="ml-auto mt-4" :rows="10" :totalRecords="total" :rowsPerPageOptions="[10, 20, 30]"
      @page="changePage">
      <template #start>
        共 {{ total }} 个
      </template>
    </Paginator>
  </div>

  <Drawer v-model:visible="instanceCreateVisible" header="创建实例" position="right" :showCloseIcon="false" class="!w-96">
    <Fluid class="flex flex-col gap-4 w-full">
      <div class="flex flex-col gap-2">
        <label>主机名称</label>
        <InputText v-model="selectServer.name" type="text" disabled />
      </div>
      <div class="flex flex-col gap-2">
        <label>GPU类型</label>
        <InputText v-model="selectServer.gpu_type" type="text" disabled />
      </div>
      <div class="flex flex-col gap-2">
        <label>CPU 核数</label>
        <InputGroup>
          <InputNumber v-model="instanceCreateForm.cpu_count" disabled />
          <InputGroupAddon>核</InputGroupAddon>
        </InputGroup>
      </div>
      <div class="flex flex-col gap-2">
        <label>内存</label>
        <InputGroup>
          <InputNumber v-model="instanceCreateForm.memory" disabled />
          <InputGroupAddon>GB</InputGroupAddon>
        </InputGroup>
      </div>
      <div class="flex flex-col gap-2">
        <label>GPU 数量</label>
        <SelectButton v-model="instanceCreateForm.gpu_count" :options="instanceCreateForm.options"
          optionDisabled="disabled" optionLabel="label" optionValue="value" aria-labelledby="basic" :allowEmpty="false"
          @change="instanceCreateGpuChange" />
      </div>
      <div class="flex flex-col gap-2">
        <label>数据盘</label>
        <InputNumber v-model="instanceCreateForm.volume_size" showButtons buttonLayout="horizontal" :step="1"
          suffix=" GB" :min="50" :max="selectServer.volume_total - selectServer.volume_used">
          <template #incrementbuttonicon>
            <span class="pi pi-plus" />
          </template>
          <template #decrementbuttonicon>
            <span class="pi pi-minus" />
          </template>
        </InputNumber>
      </div>
      <div class="flex flex-col gap-2">
        <label>镜像</label>
        <Select v-model="instanceCreateForm.image_name" :options="imagesList" placeholder="请选择镜像" class="w-full" />
      </div>
      <div class="flex flex-col gap-2">
        <label>GPU 驱动版本</label>
        <InputText v-model="selectServer.gpu_driver_version" type="text" disabled />
      </div>
      <div class="flex flex-col gap-2">
        <label>CUDA 版本</label>
        <InputText v-model="selectServer.gpu_cuda_version" type="text" disabled />
      </div>
    </Fluid>
    <template #footer>
      <div class="flex items-center gap-4 w-52 ml-auto">
        <Button label="取消" class="flex-auto" severity="secondary" @click="instanceCreateVisible = false"></Button>
        <Button label="确定" class="flex-auto" @click="instanceCreate"></Button>
      </div>
    </template>
  </Drawer>
</template>

<script setup>
import SectionBanner from '@/components/SectionBanner.vue';

import api from '@/api';
import { useToast } from 'primevue/usetoast';
import { onMounted, ref } from 'vue';
import { useRouter } from 'vue-router';

const instanceCreateRef = ref(null)

const router = useRouter()
const toast = useToast()

const offset = ref(0)
const limit = ref(10)

const data = ref([])
const total = ref(0)

const imagesMap = ref({})
const imagesList = ref([])

const refreshIcon = ref('pi pi-refresh')

const instanceCreateVisible = ref(false)
const selectServer = ref({
  price: 0,
  price_volume: 0
})
const instanceCreateForm = ref({
  server_id: 0,
  image_name: '',
  gpu_count: 1,
  volume_size: 50
})

const changePage = async (event) => {
  limit.value = event.rows
  offset.value = event.first
  await getServersList()
}

const getServersList = async () => {
  refreshIcon.value = 'pi pi-spin pi-refresh'
  await api.UserServerList({ offset: offset.value, limit: limit.value }).then((res) => {
    data.value = res.data.data.result
    for (let i = 0; i < data.value.length; i++) {
      data.value[i].key = i
    }
    total.value = res.data.data.total
    for (let i = 0; i < data.value.length; i++) {
      if (data.value[i].gpu_used >= data.value[i].gpu_num) {
        data.value[i].disabled = true
      } else {
        data.value[i].disabled = false
      }
    }
    refreshIcon.value = 'pi pi-refresh'
  }).catch((err) => {
    toast.add({ severity: 'error', summary: '获取主机列表失败', detail: err.response.data.msg, life: 3000 })
    refreshIcon.value = 'pi pi-refresh'
  })
}

const getImages = async () => {
  await api.UserImages().then((res) => {
    imagesMap.value = res.data.data.result
    imagesList.value = Object.keys(imagesMap.value).map(key => (key)).sort()
  }).catch((err) => {
    toast.add({ severity: 'error', summary: '获取镜像列表失败', detail: err.response.data.msg, life: 3000 })
  })
}

const clickCreate = (event, record) => {
  selectServer.value = record
  instanceCreateForm.value = {
    server_id: selectServer.value.id,
    image_name: '',
    gpu_count: 1,
    volume_size: 50,
    cpu_count: selectServer.value.cpu_count_per_gpu,
    memory: selectServer.value.memory_per_gpu
  }
  instanceCreateForm.value.options = []
  for (let i = 0; i < selectServer.value.gpu_num; i++) {
    const tmp = {
      label: i + 1,
      value: i + 1,
      disabled: i + 1 > selectServer.value.gpu_num - selectServer.value.gpu_used
    }
    instanceCreateForm.value.options.push(tmp)
  }
  console.log(selectServer.value)

  instanceCreateVisible.value = true
}

const instanceCreateGpuChange = (event) => {
  instanceCreateForm.value.cpu_count = selectServer.value.cpu_count_per_gpu * event.value
  instanceCreateForm.value.memory = selectServer.value.memory_per_gpu * event.value
}

const instanceCreate = async () => {
  // console.log(instanceCreateForm.value)
  if (instanceCreateForm.value.server_id === 0 || instanceCreateForm.value.server_id === null || instanceCreateForm.value.server_id === undefined) {
    toast.add({ severity: 'error', summary: '请选择主机', life: 3000 })
    return
  }
  if (instanceCreateForm.value.volume_size < 0) {
    toast.add({ severity: 'error', summary: '数据盘大小不能小于0', life: 3000 })
    return
  }
  if (instanceCreateForm.value.image_name === '') {
    toast.add({ severity: 'error', summary: '请选择镜像', life: 3000 })
    return
  } else {
    instanceCreateForm.value.image_name = imagesMap.value[instanceCreateForm.value.image_name]
  }

  api.UserInstancesCreate(instanceCreateForm.value).then((res) => {
    toast.add({ severity: 'success', summary: '创建成功', detail: '实例创建成功', life: 3000 })
    instanceCreateVisible.value = false
    router.push('/instances')
  }).catch((error) => {
    toast.add({ severity: 'error', summary: '创建失败', detail: error.response.data.msg, life: 3000 })
  })
}

onMounted(async () => {
  await getServersList()
  await getImages()
})
</script>