<template>
  <FloatingConfigurator />
  <div class="flex items-center justify-center min-h-screen overflow-hidden">
    <div class="flex flex-col items-center justify-center">
      <img :src="logo" width="64" height="32" class="-mt-6 mb-4" />
      <div
        style="border-radius: 56px; padding: 0.3rem; background: linear-gradient(180deg, color-mix(in srgb, var(--primary-color), transparent 60%) 10%, var(--surface-ground) 30%)">
        <div class="w-full bg-surface-0 dark:bg-surface-900 py-20 px-8 sm:px-20 flex flex-col items-center"
          style="border-radius: 53px">
          <h1 class="text-surface-900 dark:text-surface-0 font-bold text-3xl lg:text-5xl mb-2">邮箱验证</h1>
          <div class="text-surface-600 dark:text-surface-200 mb-8">验证邮箱地址</div>

          <router-link v-if="status == 0"
            class="w-full flex items-center mb-8 py-8 border-surface-300 dark:border-surface-500 border-b">
            <span class="flex justify-center items-center border-2 border-primary text-primary rounded-border"
              style="height: 3.5rem; width: 3.5rem">
              <i class="pi pi-fw pi-spin pi-spinner !text-2xl"></i>
            </span>
            <span class="ml-6 flex flex-col">
              <span class="text-surface-900 dark:text-surface-0 lg:text-xl font-medium mb-0">
                验证中
              </span>
              <span class="text-surface-600 dark:text-surface-200 lg:text-xl">
                请稍后
              </span>
            </span>
          </router-link>

          <router-link v-else-if="status == 1" to="/"
            class="w-full flex items-center mb-8 py-8 border-surface-300 dark:border-surface-500 border-b">
            <span class="flex justify-center items-center border-2 border-emerald-500 text-emerald-500 rounded-border"
              style="height: 3.5rem; width: 3.5rem">
              <i class="pi pi-fw pi-check !text-2xl"></i>
            </span>
            <span class="ml-6 flex flex-col">
              <span class="text-surface-900 dark:text-surface-0 lg:text-xl font-medium mb-0">
                验证成功
              </span>
              <span class="text-surface-600 dark:text-surface-200 lg:text-xl">
                你可以开始使用了
              </span>
            </span>
          </router-link>

          <router-link v-else-if="status == 2" to="/"
            class="w-full flex items-center mb-8 py-8 border-surface-300 dark:border-surface-500 border-b">
            <span class="flex justify-center items-center border-2 border-red-500 text-red-500 rounded-border"
              style="height: 3.5rem; width: 3.5rem">
              <i class="pi pi-fw pi-times !text-2xl"></i>
            </span>
            <span class="ml-6 flex flex-col">
              <span class="text-surface-900 dark:text-surface-0 lg:text-xl font-medium mb-0">
                验证失败
              </span>
              <span class="text-surface-600 dark:text-surface-200 lg:text-xl">
                请重新验证
              </span>
            </span>
          </router-link>
          <Button as="router-link" label="回到仪表盘" to="/" />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import logo from '@/assets/logo.svg';
import FloatingConfigurator from '@/components/FloatingConfigurator.vue';

import api from '@/api';

import { useToast } from 'primevue/usetoast';
import { ref } from 'vue';
import { useRoute } from 'vue-router';

const route = useRoute()
const toast = useToast()

const status = ref(0) // 0: verifing, 1: success, 2: fail
const code = route.params.code

api.UserVerify(code).then(res => {
  status.value = 1
  toast.add({ severity: 'success', summary: '验证成功', detail: '你可以开始使用了', life: 3000 })
}).catch(err => {
  status.value = 2
  toast.add({ severity: 'error', summary: '验证失败', detail: err.response.data.msg, life: 3000 })
  console.error(err)
})
</script>