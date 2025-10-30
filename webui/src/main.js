import {createApp} from 'vue'
import App from './App.vue'
import router from './router'
import axios from './services/axios.js';
import ErrorMsg from './components/ErrorMsg.vue'
import './assets/main.css'

// Previous version of the app used localStorage but we need to mantain
// the state only between reload so now the app uses sessionStorage.
// Let's remove any local storage if they were still there
localStorage.clear();

// Create and config the App
const app = createApp(App)

app.config.globalProperties.$version = "1.0"
app.config.globalProperties.$axios = axios;

app.component("ErrorMsg", ErrorMsg);
app.use(router)
app.mount('#app')
