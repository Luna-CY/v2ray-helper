import {createRouter, createWebHistory, RouteRecordRaw} from "vue-router";
import {UseGlobalState} from "@/pinia/main";

const routers: RouteRecordRaw[] = []

routers.push({name: "home", path: "/", component: () => import("@/pages/Home.vue")})
routers.push({name: "guard", path: "/guard", component: () => import("@/pages/Guard.vue")})

export const router = createRouter({history: createWebHistory(), routes: routers})

router.beforeEach((to, from,  next) => {
    const PGS = UseGlobalState()
    if ('guard' !== to.name && !PGS.authorized) {
        next({name: "guard"})

        return
    }

    next()
})
