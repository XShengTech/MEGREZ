<template>
  <SectionBanner label="安全设置" icon="pi pi-cog text-emerald-500"></SectionBanner>
  <Fluid>
    <div class="flex flex-col md:flex-row gap-8">
      <div class="md:w-1/2">
        <div class="card rounded-2xl flex flex-col gap-4">
          <div class="font-semibold text-xl">邮箱验证</div>
          <div class="flex flex-col gap-2">
            <label for="email">邮箱</label>
            <InputText v-model="userProfile.email" id="email" type="text" disabled />
          </div>
          <div class="mt-2">
            <Button v-if="!userProfile.verify && !verifyRequesting" label="验证" severity="warn"
              style="width: 5.6rem; float: right;" @click="emailVerify" />
            <Button v-else-if="verifyRequesting" label="发送中" severity="warn" icon="pi pi-spin pi-spinner" disabled
              style="width: 7.2rem; float: right;" />
            <Button v-else label="已验证" severity="secondary" disabled style="width: 5.6rem; float: right;" />
          </div>
        </div>

        <div class="card rounded-2xl flex flex-col gap-4">
          <div class="font-semibold text-xl">修改邮箱</div>
          <div class="flex flex-col gap-2">
            <label for="email">邮箱</label>
            <InputText v-if="!emailModifyStatus" v-model="userProfile.email" id="email" type="text" disabled />
            <InputText v-else v-model="emailModifyValue" id="email" type="text" />
          </div>
          <div v-if="!emailModifyStatus" class="mt-2">
            <Button label="修改" style="width: 5.6rem; float: right;" @click="emailModifyStatus = true" />
          </div>
          <div v-else>
            <Button label="保存" style="width: 5.6rem; float: right;" />
            <Button class="mr-2" label="取消" severity="secondary" style="width: 5.6rem; float: right;"
              @click="emailModifyStatus = false" />
          </div>
        </div>
      </div>

      <div class="md:w-1/2">
        <div class="card rounded-2xl flex flex-col gap-4">
          <div class="font-semibold text-xl">修改密码</div>
          <div class="flex flex-col gap-2">
            <label for="old_password">原密码</label>
            <InputText id="old_password" type="text" />
          </div>
          <div class="flex flex-col gap-2">
            <label for="new_password">新密码</label>
            <InputText id="new_password" type="text" />
          </div>
          <div class="flex flex-col gap-2">
            <label for="re_password">重复密码</label>
            <InputText id="re_password" type="text" />
          </div>
          <div class="mt-2">
            <Button label="保存" style="width: 5.6rem; float: right;" />
          </div>
        </div>
      </div>
    </div>
  </Fluid>
</template>

<script setup>
import { useToast } from 'primevue/usetoast';
import { onMounted, ref } from 'vue';

import api from '@/api';

const toast = useToast()

const userProfile = ref({})

const verifyRequesting = ref(false)

const emailModifyStatus = ref(false)
const emailModifyValue = ref('')

const getProfile = () => {
  api.GetUserProfile().then((res) => {
    userProfile.value = res.data.data.result
    emailModifyValue.value = userProfile.value.email
    console.log(userProfile.value)
  }).catch((err) => {
    toast.add({ severity: 'error', summary: '获取用户信息失败', detail: err.response.data.msg, life: 3000 })
    console.error(err)
  })
}

const emailVerify = () => {
  verifyRequesting.value = true
  api.UserVerifyRequest().then((res) => {
    toast.add({ severity: 'success', summary: '验证邮件已发送', detail: '请前往邮箱查看', life: 3000 })
    verifyRequesting.value = false
  }).catch((err) => {
    toast.add({ severity: 'error', summary: '发送失败', detail: err.response.data.msg, life: 3000 })
    verifyRequesting.value = false
    console.error(err)
  })
}

onMounted(() => {
  getProfile()
})

</script>