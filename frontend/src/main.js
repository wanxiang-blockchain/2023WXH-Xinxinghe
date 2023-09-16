import Vue from 'vue'
import App from './App.vue'
import router from './router'
import store from './store'
import ElementUI from 'element-ui'; //引入elementUI
import 'element-ui/lib/theme-chalk/index.css'

Vue.config.productionTip = false
Vue.use(ElementUI)

import axios from 'axios';
import VueAxios from 'vue-axios';
// VueAxios 与 axios 的位置不能交换，否则会出现TypeError错误
Vue.use(VueAxios, axios)

axios.defaults.baseURL = 'http://localhost:8890/v1/triple_star';
axios.interceptors.request.use(config => {
    let token = window.localStorage.getItem("X-Token");
    let expiry = window.localStorage.getItem("X-Expiry");
    config.headers['X-Token'] = token;
    config.headers['X-Expiry'] = expiry;
    return config
})
axios.interceptors.response.use(config => {
    window.localStorage.setItem("X-Token", config.headers['X-Token'])
    window.localStorage.setItem("X-Expiry", config.headers['X-Expiry'])
    return config
})


Vue.prototype.$logon = window.localStorage.getItem("X-Token") == null;

new Vue({
    router,
    store,
    render: h => h(App)
}).$mount('#app')
