import { createStore } from 'vuex';

const store = createStore({
    state: {
        userId: null,
        isAuthenticated: false,
    },
    mutations: {
        setUserId(state, userId) {
            state.userId = userId;
        },
        setAuthentication(state, status) {
            state.isAuthenticated = status;
        },
    },
    actions: {
        async fetchUserData({ commit }) {
            try {
                const response = await fetch('http://localhost:8080/user', {
                    method: 'GET',
                    credentials: 'include', // Для отправки куков
                });

                if (!response.ok) {
                    throw new Error(`Error: ${response.statusText}`);
                }

                const data = await response.json();
                commit('setUserId', data.id); // Сохраняем userId в состоянии
                commit('setAuthentication', true); // Устанавливаем статус авторизации
            } catch (error) {
                console.error('Error fetching user data:', error);
            }
        },
    },
});

export default store;
