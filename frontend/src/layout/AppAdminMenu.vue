<script setup>
import { useProfileStore } from '@/stores/profile';
import { onMounted, ref } from 'vue';

import AppMenuItem from './AppMenuItem.vue';

const profileStore = useProfileStore();


const model = ref([]);

const defaultModel = ref({
  label: '用户中心',
  items: [
    {
      label: '实例管理',
      icon: 'pi pi-fw pi-desktop text-lime-500',
      to: '/'
    },
  ]
})

const adminModel = ref({
  label: '系统设置',
  items: [
    {
      label: '实例管理',
      icon: 'pi pi-fw pi-desktop text-lime-500',
      to: '/admin/instances'
    },
    {
      label: '用户管理',
      icon: 'pi pi-fw pi-users text-indigo-500',
      to: '/admin/users'
    },
  ]
})

const superAdminModel = ref({
  label: '系统设置',
  items: [
    {
      label: '节点管理',
      icon: 'pi pi-fw pi-server text-yellow-400',
      to: '/admin/servers'
    },
    {
      label: '实例管理',
      icon: 'pi pi-fw pi-desktop text-lime-500',
      to: '/admin/instances'
    },
    {
      label: '用户管理',
      icon: 'pi pi-fw pi-users text-indigo-500',
      to: '/admin/users'
    },
    {
      label: '镜像管理',
      icon: 'pi pi-fw pi-images text-teal-500',
      to: '/admin/images'
    },
  ]
})

onMounted(() => {
  const superAdmin = profileStore.isSuperAdmin
  const admin = profileStore.isAdmin

  if (superAdmin) {
    model.value.push(superAdminModel.value)
  } else if (admin) {
    model.value.push(adminModel.value)
  }

  model.value.push(defaultModel.value);
});
</script>

<template>
  <ul class="layout-menu">
    <template v-for="(item, i) in model" :key="item">
      <app-menu-item v-if="!item.separator" :item="item" :index="i"></app-menu-item>
      <li v-if="item.separator" class="menu-separator"></li>
    </template>
  </ul>
</template>

<style lang="scss" scoped></style>
