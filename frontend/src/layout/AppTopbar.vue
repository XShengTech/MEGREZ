<template>
  <div class="layout-topbar">
    <div class="layout-topbar-logo-container">
      <button class="layout-menu-button layout-topbar-action" @click="onMenuToggle">
        <i class="pi pi-bars"></i>
      </button>
      <router-link to="/" class="layout-topbar-logo">
        <img :src="logo" />

        <span>MEGREZ</span>
      </router-link>
    </div>

    <div class="layout-topbar-actions">
      <div class="layout-config-menu">
        <button type="button" class="layout-topbar-action" @click="toggleDarkMode">
          <i :class="['pi', { 'pi-moon': isDarkTheme, 'pi-sun': !isDarkTheme }]"></i>
        </button>
      </div>

      <button class="layout-topbar-menu-button layout-topbar-action"
        v-styleclass="{ selector: '@next', enterFromClass: 'hidden', enterActiveClass: 'animate-scalein', leaveToClass: 'hidden', leaveActiveClass: 'animate-fadeout', hideOnOutsideClick: true }">
        <i class="pi pi-ellipsis-v"></i>
      </button>

      <div class="layout-topbar-menu hidden lg:block">
        <div class="layout-topbar-menu-content">
          <!-- <button type="button" class="layout-topbar-action layout-topbar-action-highlight">
            <i class="pi pi-user"></i>
            <span>Profile</span>
          </button> -->
          <Button severity="secondary" rounded class="h-10" @click="showMenu($event)">
            <Button icon="pi pi-user" rounded class="!h-[2.3rem] !w-[2.3rem] -ml-3" />
            <span class="font-bold">{{ username }}</span>
            <i class="pi" :class="profileMenuActive ? 'pi-angle-up' : 'pi-angle-down'"></i>
          </Button>
        </div>
      </div>
    </div>
  </div>

  <Menu ref="profileMenu" :model="instanceMenuItems" :popup="true" @blur="profileMenuActive = false" />
</template>


<script setup>
import api from '@/api';
import logo from '@/assets/logo.svg';
import { useLayout } from '@/layout/composables/layout';
import { useProfileStore } from '@/stores/profile';
import { useToast } from 'primevue/usetoast';
import { onMounted, ref } from 'vue';
import { useRouter } from 'vue-router';

const { onMenuToggle, toggleDarkMode, isDarkTheme } = useLayout();

const router = useRouter()
const toast = useToast()
const profileStore = useProfileStore();
const username = ref('');

const profileMenu = ref(null);
const profileMenuActive = ref(false);
const instanceMenuItems = [
  {
    label: '个人信息',
    icon: 'pi pi-user !text-cyan-500',
    command: () => {
      profileMenu.value.hide();
      // router.push('/profile');
    }
  },
  {
    label: '安全设置',
    icon: 'pi pi-cog !text-emerald-500',
    command: () => {
      profileMenu.value.hide();
      router.push('/settings');
    }
  },
  {
    label: '退出登录',
    icon: 'pi pi-sign-out',
    command: () => {
      profileMenu.value.hide();
      logout();
    }
  }
];

const showMenu = (event) => {
  profileMenu.value.show(event);
  profileMenuActive.value = true;
}

const logout = () => {
  api.UserLogout().then(() => {
    profileStore.clearUserProfile()
    toast.add({ severity: 'success', summary: '退出登录成功', life: 3000 })
    router.push('/login')
  }).catch(_ => {
    toast.add({ severity: 'error', summary: '退出登录失败', life: 3000 })
  })
}

onMounted(() => {
  username.value = profileStore.username;
});

</script>