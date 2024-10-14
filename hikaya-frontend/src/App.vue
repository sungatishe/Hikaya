<template>
  <div id="app">
    <nav class="navbar navbar-expand-lg navbar-light bg-light">
      <div class="container">
        <a class="navbar-brand" href="#">Movie App</a>
        <button class="navbar-toggler" type="button" data-bs-toggle="collapse" data-bs-target="#navbarNav" aria-controls="navbarNav" aria-expanded="false" aria-label="Toggle navigation">
          <span class="navbar-toggler-icon"></span>
        </button>
        <div class="collapse navbar-collapse" id="navbarNav">
          <ul class="navbar-nav me-auto">
            <li class="nav-item">
              <router-link class="nav-link" to="/movies">Movies</router-link>
            </li>
          </ul>

          <ul class="navbar-nav ms-auto">
            <li v-if="!isAuthenticated" class="nav-item">
              <router-link class="nav-link" to="/login">Login</router-link>
            </li>
            <li v-if="!isAuthenticated" class="nav-item">
              <router-link class="nav-link" to="/register">Register</router-link>
            </li>
            <li v-if="isAuthenticated" class="nav-item">
              <a class="nav-link" href="/user" role="button">
                <img :src="userAvatar" alt="User Avatar" class="rounded-circle" width="30" height="30" />
              </a>
            </li>
            <li v-if="isAuthenticated" class="nav-item">
              <a class="nav-link" role="button" @click="logout">
                <img :src="logoutLogo" alt="User Avatar" class="rounded-circle" width="30" height="30" />
              </a>
            </li>
          </ul>
        </div>
      </div>
    </nav>

    <div class="container">
      <router-view @loggedIn="updateAuthStatus" />
    </div>

    <footer class="bg-light text-center py-3 mt-4">
      <p>&copy; 2024 Movie App. All rights reserved.</p>
    </footer>
  </div>
</template>

<script>
export default {
  data() {
    return {
      isAuthenticated: false,
      userAvatar: '',
      logoutLogo: 'https://cdn1.iconfinder.com/data/icons/heroicons-ui/24/logout-512.png'
    };
  },
  methods: {
    async checkAuthentication() {
      try {
        const response = await fetch('http://localhost:8080/user', {
          method: 'GET',
          credentials: 'include' // Включение куков для авторизации
        });

        if (response.ok) {
          const data = await response.json();
          this.isAuthenticated = true; // Если запрос успешен, пользователь авторизован
          this.userAvatar = data.avatarUrl || 'https://cdn-icons-png.flaticon.com/512/149/149071.png'; // Загрузите URL аватара
        } else {
          this.isAuthenticated = false; // Если ответ не успешен, считаем, что пользователь не авторизован
        }
      } catch (error) {
        console.error('Error:', error);
        this.isAuthenticated = false; // Устанавливаем статус авторизации в false при ошибке
      }
    },
    updateAuthStatus() {
      this.checkAuthentication(); // Обновляем статус аутентификации после входа
    },
    logout() {
      fetch('http://localhost:8080/logout', {
        method: 'POST',
        credentials: 'include' // Включение куков для авторизации
      }).then(() => {
        localStorage.clear();
        this.isAuthenticated = false;
        this.userAvatar = ''; // Очистить аватар при выходе
        this.$router.push('/login');
      }).catch(error => {
        console.error('Error during logout:', error);
      });
    }
  },
  created() {
    this.checkAuthentication(); // Проверяем статус авторизации при загрузке
  }
};
</script>

<style>

</style>
