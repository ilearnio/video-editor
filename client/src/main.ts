import moshaToast from 'mosha-vue-toastify'
import 'mosha-vue-toastify/dist/style.css'
import naive from 'naive-ui'
import 'nprogress/nprogress.css'
import { createPinia } from 'pinia'
import 'vfonts/Lato.css'
import { createApp } from 'vue'

import i18n from '@/i18n'
import Auth from '@/layouts/auth.vue'
import Default from '@/layouts/default/default.vue'

import App from './App.vue'
import router from './router'

const app = createApp(App)
app.use(naive)
app.use(moshaToast)
app.use(createPinia())
app.use(i18n)
app.use(router)
app.component('DefaultLayout', Default)
app.component('AuthLayout', Auth)
app.mount('#app')
