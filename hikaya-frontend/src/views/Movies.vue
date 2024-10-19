<template>
  <div class="container mt-5">
    <h2 class="mb-4">Movies</h2>

    <!-- Форма поиска -->
    <form @submit.prevent="searchMovies">
      <div class="input-group mb-4">
        <input
            type="text"
            v-model="searchQuery"
            class="form-control"
            placeholder="Search for a movie..."
        />
        <button type="submit" class="btn btn-primary">Search</button>
      </div>
    </form>

    <div class="row">
      <div v-for="movie in movies" :key="movie.ID" class="col-md-4">
        <!-- Оборачиваем карточку в router-link -->
        <router-link :to="'/movie/' + movie.ID" class="text-decoration-none">
          <div class="card mb-3 custom-card">
            <img
                v-if="movie.poster"
                :src="movie.poster"
                class="card-img-top"
                :alt="movie.title"
            />
            <div class="card-body p-2">
              <h5 class="card-title">{{ movie.title }} ({{ movie.year }})</h5>
              <p class="card-text">{{ movie.description }}</p>
            </div>
          </div>
        </router-link>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  data() {
    return {
      movies: [], // Список фильмов
      searchQuery: '', // Поле для ввода поискового запроса
    };
  },
  async created() {
    await this.fetchMovies(); // Загрузка всех фильмов при создании компонента
  },
  methods: {
    async fetchMovies(query = '') {
      try {
        // Если query есть, добавляем его в URL
        const url = query
            ? `http://localhost:8080/movies/search?q=${encodeURIComponent(query)}`
            : 'http://localhost:8080/movies';

        const response = await fetch(url, {
          method: 'GET',
          headers: {
            'Content-Type': 'application/json',
          },
          credentials: 'include', // Для передачи куки
        });

        if (response.ok) {
          const data = await response.json();
          this.movies = data; // Обновляем список фильмов
        } else {
          console.error('Failed to fetch movies');
        }
      } catch (error) {
        console.error('Error fetching movies:', error);
      }
    },
    async searchMovies() {
      // Выполняем поиск фильмов по введенному запросу
      await this.fetchMovies(this.searchQuery);
    },
  },
};
</script>

<style scoped>
.custom-card {
  transition: transform 0.2s;
}

.custom-card:hover {
  transform: scale(1.05); /* Эффект увеличения при наведении */
  cursor: pointer;
}

.text-decoration-none {
  text-decoration: none; /* Убираем подчеркивание у ссылок */
  color: inherit; /* Сохраняем цвет текста как в карточке */
}
</style>
