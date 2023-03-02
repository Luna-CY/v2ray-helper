import '@/style.scss';
import "reset-css";

import {createApp} from 'vue';
import App from '@/App.vue';
import {pinia} from "@/pinia/main";
import {router} from "@/router/main";

createApp(App).use(pinia).use(router).mount('#app')
