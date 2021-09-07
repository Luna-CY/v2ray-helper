import {createRouter, createWebHistory, RouteRecordRaw} from 'vue-router'
import store from '@/store'

import Home from '../views/Home.vue'
import {ElMessage} from "element-plus"

const routes: Array<RouteRecordRaw> = [
  {
    path: '/',
    name: 'Home',
    component: Home,
    meta: {
      title: "节点列表"
    }
  },
  {
    path: '/auth/login',
    name: 'Login',
    component: () => import('@/views/auth/Login.vue'),
    meta: {
      title: '验证'
    }
  },
  {
    name: '401',
    path: '/401',
    component: () => import('@/views/error/401.vue'),
    meta: {
      title: '401',
    },
  },
  {
    path: '/404',
    name: '404',
    component: () => import('@/views/error/404.vue'),
    meta: {
      title: '404'
    }
  },
  {
    path: '/:catchAll(.*)',
    name: 'NotFound',
    redirect: '/404',
  },
]

const router = createRouter({
  history: createWebHistory(process.env.BASE_URL),
  routes
})

router.beforeEach((to, from, next) => {
  // @ts-ignore
  const timestamp = Date.parse(new Date()) / 1000
  if ('Login' != to.name && ('' == store.getters.token.token || store.getters.token.expired < timestamp)) {
    if (store.getters.token.expired < timestamp) {
      ElMessage.info("登录信息已失效，请重新登录")
    }

    store.dispatch('logout').then(() => {
      router.push({name: 'Login', query: {"redirect": to.path}}).then()
    })

    return
  }

  if (to.meta.hasOwnProperty('title')) {
    document.title = to.meta.title + ' - V2ray Subscription'
  }

  next()
})

export default router
