import {createRouter, createWebHashHistory, RouteRecordRaw} from 'vue-router'
import store from '@/store'

import Home from '../views/Home.vue'
import {ElMessage} from "element-plus"

const routes: Array<RouteRecordRaw> = [
  {
    path: '/',
    name: 'Home',
    component: Home,
    meta: {
      title: "配置列表"
    }
  },
  {
    path: '/deploy',
    name: 'Deploy',
    component: () => import('@/views/deploy/Deploy.vue'),
    meta: {
      title: "部署服务器"
    }
  },
  {
    path: '/helper/v2ray-x',
    name: 'HelperV2rayX',
    component: () => import('@/views/help/V2rayX.vue'),
    meta: {
      title: "V2rayX配置帮助"
    }
  },
  {
    path: '/helper/v2ray-n',
    name: 'HelperV2rayN',
    component: () => import('@/views/help/V2rayN.vue'),
    meta: {
      title: "V2rayN配置帮助"
    }
  },
  {
    path: '/helper/v2ray-ng',
    name: 'HelperV2rayNG',
    component: () => import('@/views/help/V2rayNG.vue'),
    meta: {
      title: "V2rayNG配置帮助"
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
  history: createWebHashHistory(),
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
