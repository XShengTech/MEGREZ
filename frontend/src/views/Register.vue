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
            <span class="text-muted-color font-medium">注册新账号</span>
          </div>

          <div>
            <label for="username"
              class="block text-surface-900 dark:text-surface-0 text-xl font-medium mb-2">用户名</label>
            <InputText id="用户名" type="text" placeholder="用户名" class="w-full md:w-[30rem] mb-8"
              v-model="form.username" />

            <label for="email" class="block text-surface-900 dark:text-surface-0 text-xl font-medium mb-2">邮箱</label>
            <InputText id="email" type="text" placeholder="邮箱" class="w-full md:w-[30rem] mb-8" v-model="form.email" />

            <label for="password" class="block text-surface-900 dark:text-surface-0 font-medium text-xl mb-2">密码</label>
            <Password id="password" v-model="form.password" placeholder="密码" :toggleMask="true" class="mb-4" fluid
              :feedback="false" @keydown.enter="handleSubmit" />

            <div class="flex items-center justify-between mt-2 mb-8 gap-8">
              <span class="font-medium no-underline ml-2 text-right cursor-pointer text-slate-600">已有账号？<span
                  class="text-primary" @click="handleLogin">立即登入</span></span>
            </div>
            <Button label="注册" class="w-full" @click="handleSubmit"></Button>
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
import { useToast } from 'primevue/usetoast';
import { ref } from 'vue';
import { useRouter } from 'vue-router';

const router = useRouter()
const toast = useToast()

const form = ref({
  username: '',
  email: '',
  password: ''
})

const handleSubmit = () => {
  api.UserRegister(form.value).then(res => {
    toast.add({ severity: 'success', summary: '注册成功', life: 3000 })
    router.push('/login')
  }).catch(err => {
    console.error(err)
    toast.add({ severity: 'error', summary: '注册失败', detail: err.response.data.msg, life: 3000 })
  })
}

const handleLogin = () => {
  router.push('/login')
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