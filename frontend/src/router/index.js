import AppLayout from '@/layout/AppLayout.vue';
import { useProfileStore } from '@/stores/profile';
import { createRouter, createWebHistory } from 'vue-router';

import Forget from '@/views/Forget.vue';
import Login from '@/views/Login.vue';
import Register from '@/views/Register.vue';
import Reset from '@/views/Reset.vue';
import Verify from '@/views/Verify.vue';

import InstanceCreate from '@/views/users/InstanceCreate.vue';
import InstanceList from '@/views/users/InstanceList.vue';
import Settings from '@/views/users/Settings.vue';

import Images from '@/views/admin/Images.vue';
import Instances from '@/views/admin/Instances.vue';
import Servers from '@/views/admin/Servers.vue';
import Users from '@/views/admin/Users.vue';

import AppAdminLayout from '@/layout/AppAdminLayout.vue';
import NotFound from '@/views/NotFound.vue';

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      redirect: '/instances'
    },
    {
      path: '/login',
      name: 'login',
      component: Login
    },
    {
      path: '/register',
      name: 'register',
      component: Register
    },
    {
      path: '/forget',
      name: 'forget',
      component: Forget
    },
    {
      path: '/reset/:code',
      name: 'forget-reset',
      component: Reset
    },
    {
      path: '/verify/:code',
      name: 'verify',
      component: Verify
    },
    {
      path: '/',
      name: 'dashboard',
      component: AppLayout,
      children: [
        {
          path: 'instances',
          name: 'instances',
          component: InstanceList
        },
        {
          path: 'instances/create',
          name: 'instances-create',
          component: InstanceCreate
        },
        {
          path: 'settings',
          name: 'settings',
          component: Settings
        }
      ]
    },
    {
      path: '/admin/',
      name: 'admin-dashboard',
      component: AppAdminLayout,
      children: [
        {
          path: 'images',
          name: 'admin-images',
          component: Images
        },
        {
          path: 'instances',
          name: 'admin-instances',
          component: Instances
        },
        {
          path: 'servers',
          name: 'admin-servers',
          component: Servers
        },
        {
          path: 'users',
          name: 'admin-users',
          component: Users
        }
      ]
    },
    {
      path: '/:pathMatch(.*)*',
      name: 'not-found',
      component: NotFound
    }
  ]
});

router.beforeEach((to, from, next) => {
  const profileStore = useProfileStore();
  if (to.path.startsWith('/admin')) {
    if (!profileStore.isAdmin && !profileStore.isSuperAdmin) {
      next({ name: 'not-found' });
      return;
    }
    if ((to.path.endsWith('servers') || to.path.endsWith('settings')) && !profileStore.isSuperAdmin) {
      next({ name: 'not-found' });
      return;
    }
  }
  next();
});

export default router;
