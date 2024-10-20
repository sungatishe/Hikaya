import { createRouter, createWebHistory } from 'vue-router';
import Login from '../views/Login.vue';
import Register from '../views/Register.vue';
import Movies from '../views/Movies.vue';
import User from '../views/User.vue';
import MovieDetails from "@/views/MovieDetails.vue";

const routes = [
    {
        path: '/login',
        name: 'Login',
        component: Login,
        beforeEnter: async (to, from, next) => {
            const isAuthenticated = await fetchUserData();
            if (isAuthenticated) {
                next({ name: 'Movies' }); // Перенаправление на Movies или другой маршрут
            } else {
                next(); // Разрешить доступ к странице логина
            }
        },
    },
    {
        path: '/register',
        name: 'Register',
        component: Register,
        beforeEnter: async (to, from, next) => {
            const isAuthenticated = await fetchUserData();
            if (isAuthenticated) {
                next({ name: 'Movies' }); // Перенаправление на Movies или другой маршрут
            } else {
                next(); // Разрешить доступ к странице регистрации
            }
        },
    },
    {
        path: '/movies',
        name: 'Movies',
        component: Movies,
    },
    {
        path: '/user',
        name: 'User',
        component: User,
        beforeEnter: async (to, from, next) => {
            const isAuthenticated = await fetchUserData();
            if (!isAuthenticated) {
                next({ name: 'Movies' }); // Перенаправление на Movies или другой маршрут
            } else {
                next(); // Разрешить доступ к странице логина
            }
        },
    },
    {
        path: '/movie/:id',
        name: 'MovieDetails',
        component: MovieDetails,
        props: true,
    },
    // Обработка всех остальных маршрутов
    {
        path: '/:catchAll(.*)', // захватывает все несуществующие маршруты
        redirect: { name: 'Movies' }, // перенаправление на Movies
    },
];

const router = createRouter({
    history: createWebHistory(),
    routes,
});

export default router;

// Функция проверки аутентификации
async function fetchUserData() {
    try {
        const response = await fetch('http://api-gateway-service:8080/user', {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json',
            },
            credentials: 'include',
        });

        if (response.ok) {
            const data = await response.json();
            return data.id !== null; // Возвращаем true, если userId существует
        } else {
            return false; // Возвращаем false, если запрос не удался
        }
    } catch (error) {
        console.error('Error fetching user data:', error);
        return false; // Возвращаем false в случае ошибки
    }
}
