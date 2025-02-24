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
            <span class="text-muted-color font-medium">重置密码</span>
          </div>

          <div>
            <label for="password"
              class="block text-surface-900 dark:text-surface-0 font-medium text-xl mb-2">新密码</label>
            <Password id="password" v-model="form.password" placeholder="新密码" :toggleMask="true"
              class="w-full md:w-[30rem] mb-4" fluid :feedback="false" />

            <label for="repassword"
              class="block text-surface-900 dark:text-surface-0 font-medium text-xl mb-2">重复密码</label>
            <Password id="repassword" v-model="form.repassword" placeholder="重复密码" :toggleMask="true" class="mb-4" fluid
              :feedback="false" @keydown.enter="handleSubmit" />

            <div class="flex items-center justify-between mt-2 mb-8 gap-8">
              <span class="font-medium no-underline ml-2 text-right cursor-pointer text-slate-600">记起密码？<span
                  class="text-primary" @click="handleLogin">立即登入</span></span>
            </div>
            <Button label="重置密码" class="w-full" @click="handleSubmit"></Button>
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
import { useRoute, useRouter } from 'vue-router';

const router = useRouter()
const route = useRoute()
const toast = useToast()

const form = ref({
  code: route.params.code,
  password: '',
  repassword: ''
})

const handleSubmit = () => {
  if (!form.value.code) {
    toast.add({ severity: 'error', summary: '验证码错误', detail: '请检查后重新尝试', life: 3000 })
    return
  }

  if (!form.value.password) {
    toast.add({ severity: 'error', summary: '密码不能为空', detail: '请检查后重新尝试', life: 3000 })
    return
  }

  if (!form.value.repassword) {
    toast.add({ severity: 'error', summary: '重复密码不能为空', detail: '请检查后重新尝试', life: 3000 })
    return
  }

  if (form.value.password !== form.value.repassword) {
    toast.add({ severity: 'error', summary: '两次密码不一致', detail: '请检查后重新尝试', life: 3000 })
    return
  }

  api.UserForgerPassword(form.value).then(res => {
    toast.add({ severity: 'success', summary: '重置密码成功', detail: '请登录', life: 3000 })
    router.push('/login')
  }).catch(err => {
    console.error(err)
    toast.add({ severity: 'error', summary: '重置密码失败', detail: '请检查后重新尝试', life: 3000 })
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