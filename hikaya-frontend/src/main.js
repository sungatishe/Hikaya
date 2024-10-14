import { createApp } from 'vue';
import App from './App.vue';
import router from './router'; // Импортируй router
import 'bootstrap/dist/css/bootstrap.min.css'; // Импортируй Bootstrap

createApp(App)
    .use(router) // Подключи router
    .mount('#app');
