<template>
  <SectionBanner label="镜像管理" icon="pi pi-images text-teal-500">
    <Button icon="pi pi-plus" label="添加镜像" rounded @click="addImage" />
    <Button severity="success" icon="pi pi-save" label="保存" rounded @click="saveImages" />
  </SectionBanner>

  <div class="card mt-7 rounded-2xl">
    <!-- <Fluid v-for="(image, index) in data" :key="image.name">
      <Fieldset :legend="image.name" :toggleable="true">
        <InputText v-model="data[index].tag" type="text" size="large" placeholder="Large" />
      </Fieldset>
    </Fluid> -->
    <Accordion value="0" expandIcon="pi pi-plus" collapseIcon="pi pi-minus">
      <AccordionPanel v-for="(image, index) in data" :key="image.name" :value="image.name">
        <AccordionHeader>
          <span v-if="!image.edit">{{ image.name }}</span>
          <InputText v-else v-model="tmp.name" class="flex-auto mr-6" type="text" placeholder="备注名" />
        </AccordionHeader>
        <AccordionContent>
          <span class="flex items-center gap-2 w-full">
            <span class="whitespace-nowrap">镜像名:
              <Chip v-if="!image.edit" :label="image.tag" />
              <InputText v-else v-model="tmp.tag" class="flex-auto w-full" type="text" placeholder="镜像名" />
            </span>
            <!-- <Button class="ml-auto mr-2" icon="pi pi-check" size="small" aria-label="Filter" /> -->
            <div class="ml-auto mr-0 flex gap-2">
              <Tag v-if="!image.edit" class="ml-auto" icon="pi pi-pencil" value="编辑" @click="editImage(index)" />
              <Tag v-else class="ml-auto" severity="success" icon="pi pi-save" value="保存"
                @click="saveEditImage(index)" />
              <Tag class="ml-auto" severity="danger" icon="pi pi-trash" value="删除" @click="deleteImage(index)"></Tag>
            </div>
          </span>
        </AccordionContent>
      </AccordionPanel>
    </Accordion>
  </div>
</template>


<script setup>

import api from '@/api';
import { useToast } from 'primevue/usetoast';
import { onMounted, ref } from 'vue';

const toast = useToast()

const data = ref([])
const tmp = ref(null)

const refreshIcon = ref('pi pi-refresh')

const getImages = async () => {
  refreshIcon.value = 'pi pi-spin pi-refresh'
  data.value = []
  await api.AdminImagesList().then(res => {
    Object.keys(res.data.data.result).forEach(key => {
      data.value.push({
        name: key,
        tag: res.data.data.result[key],
        edit: false
      })
    })
    data.value.sort((a, b) => a.name.localeCompare(b.name))
  }).catch(err => {
    toast.add({ severity: 'error', summary: '获取镜像列表错误', detail: err.response.data.msg, life: 3000 })
  }).finally(() => {
    refreshIcon.value = 'pi pi-refresh'
  })
}

const saveImages = async () => {
  let images = {}
  data.value.forEach(item => {
    images[item.name] = item.tag
  })
  await api.AdminImagesModify(images).then(res => {
    toast.add({ severity: 'success', summary: '保存成功', detail: res.data.msg, life: 3000 })
    getImages()
  }).catch(err => {
    toast.add({ severity: 'error', summary: '保存失败', detail: err.response.data.msg, life: 3000 })
  })
}

const addImage = () => {
  data.value.unshift({
    name: '',
    tag: '',
    edit: true
  })
  tmp.value = {
    name: '',
    tag: ''
  }
}

const editImage = (index) => {
  tmp.value = {
    name: data.value[index].name,
    tag: data.value[index].tag
  }
  data.value[index].edit = true
}

const saveEditImage = (index) => {
  data.value[index].name = tmp.value.name
  data.value[index].tag = tmp.value.tag
  data.value[index].edit = false
}

const deleteImage = (index) => {
  data.value.splice(index, 1)
}

onMounted(() => {
  getImages()
})
</script>