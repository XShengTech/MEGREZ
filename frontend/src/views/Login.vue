<template>
  <FloatingConfigurator />
  <div
    class="bg-surface-50 dark:bg-surface-950 flex items-center justify-center min-h-screen min-w-[100vw] overflow-hidden">
    <div class="flex flex-col items-center justify-center">
      <div
        style="border-radius: 56px; padding: 0.3rem; background: linear-gradient(180deg, var(--primary-color) 10%, rgba(33, 150, 243, 0) 30%)">
        <div class="w-full bg-surface-0 dark:bg-surface-900 py-20 px-8 sm:px-20" style="border-radius: 53px">
          <div class="text-center items-center flex flex-col mb-8">
            <img :src="logo" style="width: 16rem;" />
            <div class="text-surface-900 dark:text-surface-0 text-3xl font-medium mt-6 mb-4">天权 算能聚联计算平台</div>
            <span class="text-muted-color font-medium">登入以继续</span>
          </div>

          <div>
            <label for="email1" class="block text-surface-900 dark:text-surface-0 text-xl font-medium mb-2">账号</label>
            <InputText id="email1" type="text" placeholder="用户名或邮箱" class="w-full md:w-[30rem] mb-8"
              v-model="form.account" />

            <label for="password1"
              class="block text-surface-900 dark:text-surface-0 font-medium text-xl mb-2">密码</label>
            <Password id="password1" v-model="form.password" placeholder="密码" :toggleMask="true" class="mb-4" fluid
              :feedback="false" @keydown.enter="handleSubmit" />

            <div class="flex items-center justify-between mt-2 mb-8 gap-8">
              <span class="font-medium no-underline ml-2 text-right cursor-pointer text-slate-600">没有账号？<span
                  class="text-primary" @click="handleRegister">立即注册</span></span>
              <span class="font-medium no-underline ml-2 text-right cursor-pointer text-primary">忘记密码</span>
            </div>
            <Button label="登入" class="w-full" @click="handleSubmit"></Button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import logo from '@/assets/logo-text.webp';
import FloatingConfigurator from '@/components/FloatingConfigurator.vue';

import api from '@/api';
import { useProfileStore } from '@/stores/profile';
import { useToast } from 'primevue/usetoast';
import { ref } from 'vue';
import { useRouter } from 'vue-router';
const profileStore = useProfileStore()

const router = useRouter()
const toast = useToast()

const form = ref({
  account: '',
  password: ''
})

const handleSubmit = () => {
  api.UserLogin(form.value).then(res => {
    // console.log(res)
    toast.add({ severity: 'success', summary: '登入成功', detail: '欢迎回来！', life: 3000 })
    profileStore.setUserProfile(res.data.data.result)
    router.push('/')
  }).catch(err => {
    toast.add({ severity: 'error', summary: '登录失败', detail: err.response.data.msg, life: 3000 })
    console.error(err)
  })
}

const handleRegister = () => {
  router.push('/register')
}
</script>

<style scoped>
.pi-eye {
  transform: scale(1.6);
  margin-right: 1rem;
}

.pi-eye-slash {
  transform: scale(1.6);
  margin-right: 1rem;
}
</style>
