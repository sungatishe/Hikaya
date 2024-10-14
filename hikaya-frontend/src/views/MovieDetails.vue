<template>
  <div class="container mt-5">
    <div v-if="movie" class="row">
      <div class="col-md-4">
        <img v-if="movie.poster" :src="movie.poster" class="img-fluid" :alt="movie.title" />
        <div v-else class="no-poster">No Poster Available</div>
      </div>
      <div class="col-md-8">
        <h2>{{ movie.title }} ({{ movie.year }})</h2>
        <p class="lead">{{ movie.description }}</p>

        <div v-if="isAuthenticated" class="mt-4">
          <h4>Manage List</h4>
          <select v-model="selectedList" class="form-select" aria-label="Select list type" @change="updateListType">
            <option value="watched">Watched</option>
            <option value="planned">Planned</option>
            <option value="abandoned">Abandoned</option>
          </select>
          <button v-if="!isInList" class="btn btn-primary mt-2" @click="addToList">Add Movie to List</button>
          <button v-if="isInList" class="btn btn-danger mt-2" @click="removeFromList">Remove from List</button>
        </div>
      </div>
    </div>
    <div v-else class="loading">
      <p>Loading movie details...</p>
    </div>
  </div>
</template>

<script>
export default {
  props: ['id'], // ID фильма из маршрута
  data() {
    return {
      movie: null, // Данные о фильме
      selectedList: 'watched', // Выбранный тип списка (по умолчанию)
      isAuthenticated: false, // Статус авторизации
      isInList: false, // Проверка, добавлен ли фильм в список
      userId: null, // ID пользователя
    };
  },
  async created() {
    await this.fetchUserData(); // Получение данных пользователя
    await this.fetchMovieDetails(); // Получение деталей фильма
    await this.checkIfMovieInList(); // Проверка, добавлен ли фильм в список пользователя
  },
  methods: {
    async fetchUserData() {
      try {
        const response = await fetch('http://localhost:8080/user', {
          method: 'GET',
          headers: {
            'Content-Type': 'application/json',
          },
          credentials: 'include',
        });

        if (response.ok) {
          const data = await response.json();
          this.userId = data.id; // Получаем userId из ответа
          this.isAuthenticated = true; // Устанавливаем статус авторизации
        } else {
          console.error('Failed to fetch user data');
        }
      } catch (error) {
        console.error('Error fetching user data:', error);
      }
    },
    async fetchMovieDetails() {
      try {
        // Получение деталей фильма
        const response = await fetch(`http://localhost:8080/movies/${this.id}`, {
          method: 'GET',
          headers: {
            'Content-Type': 'application/json',
          },
          credentials: 'include',
        });

        if (response.ok) {
          const data = await response.json();
          this.movie = data; // Сохраняем данные о фильме
        } else {
          console.error('Failed to fetch movie details');
        }
      } catch (error) {
        console.error('Error fetching movie details:', error);
      }
    },
    async checkIfMovieInList() {
      if (!this.userId) return; // Если нет userId, выходим

      try {
        const listResponse = await fetch(`http://localhost:8080/userList/${this.userId}`, {
          method: 'GET',
          headers: {
            'Content-Type': 'application/json',
          },
          credentials: 'include',
        });

        const userMovies = await listResponse.json();
        const userMovie = userMovies.find(movie => movie.movie_id === parseInt(this.id));

        if (userMovie) {
          this.isInList = true;
          this.selectedList = userMovie.list_type; // Устанавливаем дефолтное значение списка
        }
      } catch (error) {
        console.error('Error checking movie list:', error);
      }
    },
    async addToList() {
      if (!this.userId) {
        alert('Please log in first!');
        return;
      }

      const response = await fetch('http://localhost:8080/userList', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        credentials: 'include',
        body: JSON.stringify({
          user_id: this.userId, // Используем userId
          movie_id: parseInt(this.id),
          list_type: this.selectedList,
        }),
      });

      if (response.ok) {
        alert('Movie added to list successfully!');
        this.isInList = true; // Обновляем статус добавления в список
      } else {
        alert('Failed to add movie to list');
      }
    },
    async removeFromList() {
      if (!this.userId) {
        alert('Please log in first!');
        return;
      }

      const response = await fetch(`http://localhost:8080/userList/${this.userId}/${this.id}`, {
        method: 'DELETE',
        headers: {
          'Content-Type': 'application/json',
        },
        credentials: 'include',
      });

      if (response.ok) {
        alert('Movie removed from list successfully!');
        this.isInList = false;
      } else {
        alert('Failed to remove movie from list');
      }
    },
    async updateListType() {
      // Обновляем тип списка при изменении селектора
      if (!this.userId || !this.isInList) {
        return;
      }

      try {
        const response = await fetch(`http://localhost:8080/userList/${this.userId}/${this.id}?listType=${this.selectedList}`, {
          method: 'PUT', // Отправляем PUT запрос для обновления списка
          headers: {
            'Content-Type': 'application/json',
          },
          credentials: 'include',
        });

        if (response.ok) {
          alert('List type updated successfully!');
        } else {
          alert('Failed to update list type');
        }
      } catch (error) {
        console.error('Error updating list type:', error);
        alert('An error occurred while updating the list type.');
      }
    },
  },
};
</script>

<style scoped>
.no-poster {
  width: 100%;
  height: 400px;
  background-color: #f0f0f0;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.2rem;
  color: #888;
}
.loading {
  text-align: center;
  font-size: 1.2rem;
}
</style>
