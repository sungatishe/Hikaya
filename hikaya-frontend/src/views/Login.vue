<template>
  <div class="container">
    <h2 class="mt-5">Login</h2>
    <form @submit.prevent="login">
      <div class="mb-3">
        <label for="email" class="form-label">Email address</label>
        <input type="email" v-model="email" class="form-control" id="email" required />
      </div>
      <div class="mb-3">
        <label for="password" class="form-label">Password</label>
        <input type="password" v-model="password" class="form-control" id="password" required />
      </div>
      <button type="submit" class="btn btn-primary">Login</button>
      <p class="mt-3">
        Don't have an account? <router-link to="/register">Register</router-link>
      </p>
    </form>
  </div>
</template>

<script>
import { mapActions } from 'vuex';

export default {
  data() {
    return {
      email: '',
      password: '',
    };
  },
  methods: {
    ...mapActions(['fetchUserData']), // Подключаем экшн для получения данных о пользователе

    async login() {
      try {
        const response = await fetch('http://localhost:8080/login', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify({
            email: this.email,
            password: this.password,
          }),
          credentials: 'include', // Для отправки куков
        });

        if (response.ok) {
          await this.fetchUserData(); // Получаем информацию о пользователе после успешного логина

          this.$router.push('/movies'); // Перенаправление на страницу с фильмами
          this.$emit('loggedIn'); // Эмитируем событие входа
        } else {
          alert('Invalid login credentials');
        }
      } catch (error) {
        console.error('Error during login:', error);
      }
    },
  },
};
</script>
