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

        <div class="rating">
          <h5>Average Rating: {{ movieRating ? movieRating : 'Not Rated' }}</h5>
        </div>

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

        <div class="mt-5">
          <h4>Reviews</h4>
          <div v-for="review in reviews" :key="review.id" class="review">
            <p><strong>User {{ review.user_id }}:</strong> {{ review.review }} (Rating: {{ review.rating }})</p>
          </div>
          <div v-if="isAuthenticated" class="mt-4">
            <h5>Write a Review</h5>
            <form @submit.prevent="submitReview">
              <div class="mb-3">
                <label for="rating" class="form-label">Rating (1-10)</label>
                <input type="number" v-model="newReview.rating" id="rating" min="1" max="10" required class="form-control" />
              </div>
              <div class="mb-3">
                <label for="review" class="form-label">Review</label>
                <textarea v-model="newReview.review" id="review" rows="3" required class="form-control"></textarea>
              </div>
              <button type="submit" class="btn btn-success">Submit Review</button>
            </form>
          </div>
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
      reviews: [], // Отзывы о фильме
      newReview: { // Новые отзывы
        rating: 1,
        review: '',
      },
      movieRating: null, // Рейтинг фильма
    };
  },
  async created() {
    await this.fetchUserData(); // Получение данных пользователя
    await this.fetchMovieDetails(); // Получение деталей фильма
    await this.checkIfMovieInList(); // Проверка, добавлен ли фильм в список пользователя
    await this.fetchMovieReviews(); // Получение отзывов о фильме
    await this.fetchMovieRating(); // Получение рейтинга фильма
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
    async fetchMovieReviews() {
      try {
        const response = await fetch(`http://localhost:8080/movie/${this.id}/reviews`, {
          method: 'GET',
          headers: {
            'Content-Type': 'application/json',
          },
          credentials: 'include',
        });

        if (response.ok) {
          this.reviews = await response.json(); // Сохраняем отзывы
        } else {
          console.error('Failed to fetch movie reviews');
        }
      } catch (error) {
        console.error('Error fetching movie reviews:', error);
      }
    },
    async fetchMovieRating() {
      try {
        const response = await fetch(`http://localhost:8080/movie/${this.id}/rating`, {
          method: 'GET',
          headers: {
            'Content-Type': 'application/json',
          },
          credentials: 'include',
        });

        if (response.ok) {
          const data = await response.json();
          this.movieRating = data.avr_rating; // Сохраняем рейтинг фильма
        } else {
          console.error('Failed to fetch movie rating');
        }
      } catch (error) {
        console.error('Error fetching movie rating:', error);
      }
    },
    async submitReview() {
      if (!this.userId) {
        alert('Please log in first!');
        return;
      }

      const response = await fetch('http://localhost:8080/reviews', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        credentials: 'include',
        body: JSON.stringify({
          user_id: this.userId,
          movie_id: parseInt(this.id),
          rating: this.newReview.rating,
          review: this.newReview.review,
        }),
      });

      if (response.ok) {
        alert('Review submitted successfully!');
        this.newReview = { rating: 1, review: '' }; // Сброс формы
        await this.fetchMovieReviews(); // Обновляем список отзывов
        await this.fetchMovieRating(); // Обновляем рейтинг фильма
      } else {
        alert('Failed to submit review');
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
        this.isInList = false; // Обновляем статус
      } else {
        alert('Failed to remove movie from list');
      }
    },
    async updateListType() {
      if (!this.userId) return;

      const response = await fetch(`http://localhost:8080/userList/${this.userId}/${this.id}`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
        },
        credentials: 'include',
        body: JSON.stringify({ list_type: this.selectedList }), // Обновляем тип списка
      });

      if (response.ok) {
        alert('List type updated successfully!');
      } else {
        console.error('Failed to update list type');
      }
    },
  },
};
</script>

<style scoped>
.no-poster {
  width: 100%;
  height: 400px;
  background-color: #ccc;
  display: flex;
  justify-content: center;
  align-items: center;
  font-size: 1.5rem;
  color: #555;
}
.review {
  border-bottom: 1px solid #ccc;
  margin-bottom: 10px;
  padding-bottom: 10px;
}
.loading {
  text-align: center;
  margin-top: 20px;
}
.rating {
  margin-top: 10px;
  font-weight: bold;
}
</style>
