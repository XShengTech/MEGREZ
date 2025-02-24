import AppLayout from '@/layout/AppLayout.vue';
import { createRouter, createWebHistory } from 'vue-router';

import Login from '@/views/Login.vue';
import Register from '@/views/Register.vue';
import Verify from '@/views/Verify.vue';

import InstanceCreate from '@/views/users/InstanceCreate.vue';
import InstanceList from '@/views/users/InstanceList.vue';
import Settings from '@/views/users/Settings.vue';

import Images from '@/views/admin/Images.vue';
import Instances from '@/views/admin/Instances.vue';
import Servers from '@/views/admin/Servers.vue';
import Users from '@/views/admin/Users.vue';

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
        },
        {
          path: 'admin/images',
          name: 'admin-images',
          component: Images
        },
        {
          path: 'admin/instances',
          name: 'admin-instances',
          component: Instances
        },
        {
          path: 'admin/servers',
          name: 'admin-servers',
          component: Servers
        },
        {
          path: 'admin/users',
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

export default router;
