<script setup>
import { useProfileStore } from '@/stores/profile';
import { onMounted, ref } from 'vue';

import AppMenuItem from './AppMenuItem.vue';

const profileStore = useProfileStore();

const model = ref([
  {
    label: '实例管理',
    items: [
      {
        label: '实例列表',
        icon: 'pi pi-fw pi-desktop text-blue-600',
        to: '/instances'
      },
      {
        label: '创建实例',
        icon: 'pi pi-fw pi-plus-circle text-purple-600',
        to: '/instances/create'
      },
    ]
  },
  {
    label: '费用管理',
    items: [
      {
        label: '历史订单',
        icon: 'pi pi-fw pi-ticket text-emerald-500',
        // disabled: true,
        // to: '/landing'
      }
    ]
  },
  {
    label: '设置',
    items: [
      {
        label: '个人信息',
        icon: 'pi pi-fw pi-user text-cyan-500',
        // to: '/settings'
      },
      {
        label: '安全设置',
        icon: 'pi pi-fw pi-cog text-emerald-500',
        to: '/settings',
        // disabled: true,
        // to: '/landing'
      },
    ]
  },
  {
    label: '相关信息',
    items: [
      {
        label: '使用文档',
        icon: 'pi pi-fw pi-book text-amber-500',
        url: 'http://docs.megrez.xsheng-ai.com/guide/usage/',
        target: '_blank'
      },
    ]
  }
]);

const adminModel = ref({
  label: '系统设置',
  items: [
    {
      label: '管理后台',
      icon: 'pi pi-fw pi-sliders-h text-yellow-500',
      to: '/admin/instances'
    },
  ]
})

const superAdminModel = ref({
  label: '系统设置',
  items: [
    {
      label: '管理后台',
      icon: 'pi pi-fw pi-sliders-h text-yellow-400',
      to: '/admin/servers'
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
