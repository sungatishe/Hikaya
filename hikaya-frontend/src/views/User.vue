<template>
  <div>
    <h1>User Information</h1>
    <p v-if="user">Name: {{ user.name }}</p>
    <p v-if="user">Email: {{ user.email }}</p>
    <p v-if="errorMessage" class="text-danger">{{ errorMessage }}</p>

    <div>
      <h1>Your Movie Lists</h1>
      <div v-if="!movieLists.watched.length && !movieLists.planned.length && !movieLists.abandoned.length">
        No movies found in your lists.
      </div>
      <div v-else>
        <h2>Watched Movies</h2>
        <div class="movie-list">
          <div v-for="movie in movies.watched" :key="movie.ID" class="movie-card" @click="goToMovieDetails(movie.movie_id)">
            <img :src="movie.poster" alt="Movie Poster" class="poster" />
          </div>
        </div>

        <h2>Planned Movies</h2>
        <div class="movie-list">
          <div v-for="movie in movies.planned" :key="movie.ID" class="movie-card" @click="goToMovieDetails(movie.movie_id)">
            <img :src="movie.poster" alt="Movie Poster" class="poster" />
          </div>
        </div>

        <h2>Abandoned Movies</h2>
        <div class="movie-list">
          <div v-for="movie in movies.abandoned" :key="movie.ID" class="movie-card" @click="goToMovieDetails(movie.movie_id)">
            <img :src="movie.poster" alt="Movie Poster" class="poster" />
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  data() {
    return {
      user: null,
      errorMessage: null,
      userId: null, // Добавляем переменную для хранения userId
      movieLists: {
        watched: [],
        planned: [],
        abandoned: [],
      },
      movies: {
        watched: [],
        planned: [],
        abandoned: [],
      },
    };
  },
  mounted() {
    this.fetchUserData();
  },
  methods: {
    async fetchUserData() {
      try {
        const response = await fetch('http://localhost:8080/user', {
          method: 'GET',
          credentials: 'include',
        });

        if (!response.ok) {
          throw new Error(`Error: ${response.statusText}`);
        }

        const data = await response.json();
        this.user = data;
        this.userId = data.id; // Предполагаем, что userId находится в данных пользователя
        this.fetchUserMovieLists(); // Вызываем fetchUserMovieLists после получения userId
      } catch (error) {
        this.errorMessage = error.message;
        console.error('Error fetching user data:', error);
      }
    },
    async fetchUserMovieLists() {
      try {
        if (!this.userId) {
          throw new Error('User not authenticated');
        }

        const response = await fetch(`http://localhost:8080/userList/${this.userId}`, {
          method: 'GET',
          credentials: 'include',
        });

        if (!response.ok) {
          throw new Error(`Error: ${response.statusText}`);
        }

        const data = await response.json();
        console.log('Fetched movie lists:', data);

        // Сортируем фильмы по типам
        this.movieLists.watched = data.filter(movie => movie.list_type === 'watched');
        this.movieLists.planned = data.filter(movie => movie.list_type === 'planned');
        this.movieLists.abandoned = data.filter(movie => movie.list_type === 'abandoned');

        // Получаем информацию о фильмах
        await this.fetchMoviesDetails();

        console.log('Sorted movie lists:', this.movieLists);
      } catch (error) {
        this.errorMessage = error.message;
        console.error('Error fetching user movie lists:', error);
      }
    },
    async fetchMoviesDetails() {
      try {
        const movieIds = [
          ...this.movieLists.watched.map(movie => movie.movie_id),
          ...this.movieLists.planned.map(movie => movie.movie_id),
          ...this.movieLists.abandoned.map(movie => movie.movie_id),
        ];

        const uniqueMovieIds = [...new Set(movieIds)];

        console.log('Fetching details for the following unique movie IDs:', uniqueMovieIds);

        const movieFetchPromises = uniqueMovieIds.map(id => {
          console.log(`Fetching details for movie ID: ${id}`);
          return fetch(`http://localhost:8080/movies/${id}`, {
            method: 'GET',
            credentials: 'include',
          }).then(res => {
            if (!res.ok) {
              throw new Error(`Error fetching movie ID ${id}: ${res.statusText}`);
            }
            return res.json();
          });
        });

        const moviesData = await Promise.all(movieFetchPromises);

        console.log('Fetched movie details:', moviesData);

        this.movies.watched = this.movieLists.watched.map(movie => ({
          ...movie,
          poster: moviesData.find(m => m.ID === movie.movie_id)?.poster || 'default-poster.jpg',
        }));

        this.movies.planned = this.movieLists.planned.map(movie => ({
          ...movie,
          poster: moviesData.find(m => m.ID === movie.movie_id)?.poster || 'default-poster.jpg',
        }));

        this.movies.abandoned = this.movieLists.abandoned.map(movie => ({
          ...movie,
          poster: moviesData.find(m => m.ID === movie.movie_id)?.poster || 'default-poster.jpg',
        }));
      } catch (error) {
        this.errorMessage = error.message;
        console.error('Error fetching movie details:', error);
      }
    },
    goToMovieDetails(movieId) {
      this.$router.push({ name: 'MovieDetails', params: { id: movieId } });
    },
  },
};
</script>

<style scoped>
.text-danger {
  color: red;
}

.movie-list {
  display: flex;
  flex-wrap: wrap;
  gap: 10px; /* Пробел между карточками */
}

.movie-card {
  cursor: pointer; /* Курсор указывает на кликабельность */
}

.poster {
  width: 100px; /* Ширина постера */
  height: auto; /* Высота пропорционально ширине */
  border-radius: 5px; /* Скругление углов */
  transition: transform 0.2s; /* Анимация при наведении */
}

.poster:hover {
  transform: scale(1.05); /* Увеличение постера при наведении */
}
</style>
