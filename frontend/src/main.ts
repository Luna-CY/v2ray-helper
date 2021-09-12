import {createApp} from 'vue'
import App from './App.vue'
import router from './router'
import store from './store'
import axios, {AxiosResponse} from "axios"
import {BaseResponse} from "@/api/base"
import ElementPlus, {ElMessage} from 'element-plus'

import 'reset-css'
import 'element-plus/dist/index.css'
import 'element-theme-chalk/lib/display.css'
import '@/assets/scss/global.scss'

declare global {
  interface Date {
    format(fmt: string): string
  }
}

Date.prototype.format = function (fmt) {
  const o = {
    "M+": this.getMonth() + 1,                 //月份
    "d+": this.getDate(),                    //日
    "h+": this.getHours(),                   //小时
    "m+": this.getMinutes(),                 //分
    "s+": this.getSeconds(),                 //秒
    "q+": Math.floor((this.getMonth() + 3) / 3), //季度
    "S": this.getMilliseconds()             //毫秒
  };

  if (/(y+)/.test(fmt)) {
    fmt = fmt.replace(RegExp.$1, (this.getFullYear() + "").substr(4 - RegExp.$1.length))
  }

  for (const k in o) {
    if (new RegExp("(" + k + ")").test(fmt)) {
      // @ts-ignore
      fmt = fmt.replace(RegExp.$1, (RegExp.$1.length == 1) ? (o[k]) : (("00" + o[k]).substr(("" + o[k]).length)))
    }
  }
  return fmt;
}

store.dispatch('load').then()

// 请求前的拦截器
axios.interceptors.request.use(function (request) {
  if ("" != store.getters.token.token) {
    request.headers["Authorization"] = "Bearer " + store.getters.token.token
  }

  if ('development' == process.env.NODE_ENV) {
    request.url = `http://127.0.0.1:8888${request.url}`
  }

  return request
})

// 注册统一响应拦截器
axios.interceptors.response.use(function (response: AxiosResponse<BaseResponse>) {
  if (200 == response.status && (20001 == response.data.code || 20002 == response.data.code)) {
    store.dispatch('logout').then()

    router.push({name: "Login"}).then()
  }

  return response
}, (error) => {
  if (error.response) {
    if (400 <= error.response.status && error.response.status < 500) {
      ElMessage.error('无效的数据请求')

      return
    }

    if (500 <= error.response.status) {
      ElMessage.error('服务器异常, 请稍后重试')

      return
    }

    console.log('请求失败, 未知异常, 请联系管理员')
  }
})

createApp(App).use(store).use(router).use(ElementPlus, {size: 'mini'}).mount('#app')
