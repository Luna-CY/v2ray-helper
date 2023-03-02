import {createPinia, defineStore} from "pinia";

export const pinia = createPinia()

export const UseGlobalState = defineStore('pinia/global/state', {
    state: (): { _authorized: boolean } => {
        return {_authorized: false}
    },
    getters: {
        authorized(): boolean {
            return this._authorized
        }
    },
    actions: {
        authorize() {
            this._authorized = true
        }
    }
})
