import { createRouter, createWebHistory } from 'vue-router';
import Login from '../views/Login.vue';
import Register from '../views/Register.vue';
import Movies from '../views/Movies.vue';
import User from '../views/User.vue';
import MovieDetails from "@/views/MovieDetails.vue";
import store from "@/store/store.js";

const routes = [
    {
        path: '/login',
        name: 'Login',
        component: Login,
    },
    {
        path: '/register',
        name: 'Register',
        component: Register,
    },
    {
        path: '/movies',
        name: 'Movies',
        component: Movies,
        // Убрали проверку на аутентификацию
    },
    {
        path: '/user',
        name: 'User',
        component: User,
        // Убрали проверку на аутентификацию
    },
    {
        path: '/movie/:id',
        name: 'MovieDetails',
        component: MovieDetails,
        props: true,
        // Убрали проверку на аутентификацию
    },
];

const router = createRouter({
    history: createWebHistory(),
    routes,
});

export default router;
