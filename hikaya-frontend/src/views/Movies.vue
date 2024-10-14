<template>
  <div class="container mt-5">
    <h2 class="mb-4">Movies</h2>
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
      movies: [],
    };
  },
  async created() {
    try {
      const response = await fetch('http://localhost:8080/movies', {
        method: 'GET',
        headers: {
          'Content-Type': 'application/json',
        },
        credentials: 'include', // Для передачи куки
      });

      if (response.ok) {
        const data = await response.json();
        this.movies = data; // Сохраняем список фильмов
      } else {
        console.error('Failed to fetch movies');
      }
    } catch (error) {
      console.error('Error fetching movies:', error);
    }
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
