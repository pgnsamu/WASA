<template>
    <div class="container d-flex justify-content-center align-items-center vh-100">
        <div class="card p-4" style="width: 400px;">
            <h3 class="card-title text-center mb-4">Login</h3>
            <form @submit.prevent="handleLogin">
                <div class="mb-3">
                    <label for="username" class="form-label">Username</label>
                    <input type="username" class="form-control" id="username" v-model="username" required>
                </div>
                <div class="d-grid">
                    <button type="submit" class="btn btn-primary">Login</button>
                </div>
            </form>
            <div v-if="errorMessage" class="alert alert-danger mt-3" role="alert">
                {{ errorMessage }}
            </div>
        </div>
    </div>
</template>

<script>
export default {
    data() {
        return {
            username: '',
            errorMessage: ''
        };
    },
    async mounted() {
        this.userId = this.$route.params.userId;
        localStorage.removeItem('authToken');
    },
    methods: {
        handleLogin() {
            // Placeholder for backend link
            const backendUrl = '/session';

            // Example login payload
            const loginData = {
                username: this.username
            };

            // Example POST request (Replace with actual AJAX request)
            this.$axios.post(backendUrl, loginData)
                .then(response => {
                    if (response.status !== 200) {
                        this.errorMessage = 'Invalid username. Please try again.';
                        throw new Error(`HTTP error! Status: ${response.status}`);
                    }
                    return response.data;
                })
                .then(data => {
                    if (data.user_id && data.token) {
                        // Navigate to UserView with userId and token
                        localStorage.setItem('authToken', data.token); // Store the token securely
                        this.$router.push({
                            name: 'HomeView',
                            params: { userId: data.user_id }
                        });
                    } else {
                        // Handle unsuccessful login, assuming an error message is present
                        this.errorMessage = data.message || 'Login failed. Please check your credentials.';
                    }
                })
                .catch(error => {
                    this.errorMessage = 'An error occurred. Please try again later.';
                    console.error('Error:', error);
                });
        }
    }
};
</script>

<style scoped>
/* Optional custom styles */
</style>