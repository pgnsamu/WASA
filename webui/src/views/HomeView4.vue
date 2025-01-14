<template>
    <div class="dashboard container-fluid d-flex flex-column" style="min-height: 100vh;">
        <div class="row flex-grow-1">
            <!-- Sidebar -->
            <div id="app" class="col-md-3 col-lg-3 bg-light border-end d-flex flex-column">
                <div class="sidebar-header py-3">
                    <div class="d-flex align-items-center">
                        <button class="btn btn-sm me-2" @click="toggleView()">
                            <img :src="convertBlobToBase64(userInfo.photo)" alt="avatar" class="rounded-circle image">
                        </button>
                        <div>
                            <h5 v-if="isChatView" class="mb-0">Chats</h5>
                            <h5 v-else class="mb-0">Profile Settings</h5>
                        </div>
                    </div>

                </div>
                <div v-if="isChatView" class="flex-grow-1">
                    <ul class="list-group list-group-flush">
                        <li class="list-group-item" v-for="chat in chats" :key="chat.id" @click="selectChat(chat)">
                            <div class="d-flex align-items-center">
                                <img :src="convertBlobToBase64(chat.photo)" alt="avatar"
                                    class="rounded-circle me-2 image">
                                <div>
                                    <h6 class="mb-0">{{ chat.name }}</h6>
                                    <small class="text-muted">{{ chat.lastMessage }}</small>
                                </div>
                            </div>
                        </li>
                    </ul>
                </div>
                <div v-else class="flex-grow-1">
                    <!-- Profile Settings Content -->
                    <form>
                        <div class="mb-3">
                            <label for="photo" class="form-label mt-4">Profile Photo</label>
                            <div class="d-flex align-items-center">
                                <input type="file" id="photo" class="form-control me-2" @change="handlePhotoUpload">
                                <button type="submit" class="btn btn-primary btn-sm"
                                    style="height: 35px; font-size: 0.8rem;" @click="setMyPhoto">SAVVVVA </button>
                            </div>
                            <label for="username" class="form-label mt-4">Username</label>
                            <div class="d-flex align-items-center">
                                <input type="text" id="username" class="form-control me-2" v-model="userInfo.username">
                                <button type="submit" class="btn btn-primary btn-sm"
                                    style="height: 35px; font-size: 0.8rem;" @click="putUsername">Save </button>
                            </div>
                        </div>
                    </form>
                </div>
            </div>

            <!-- Main Chat Area -->
            <div class="col-md-9 col-lg-9 d-flex flex-column">
                <div v-if="selectedChat"
                    class="chat-header border-bottom py-3 d-flex justify-content-between align-items-center">
                    <h5>{{ selectedChat.name }}</h5>
                </div>

                <!-- Chat Body -->
                <div v-if="selectedChat" class="chat-body p-3 flex-grow-1" style="overflow-y: auto;">
                    <div v-for="message in messages" :key="message.id" class="mb-3">
                        <div :class="['p-2', message.senderId == userId ? 'bg-primary text-white' : 'bg-light']">
                            <div v-if="message.photoContent">
                                <img :src="message.photoContent" alt="photo" class="img-fluid" />
                            </div>
                            <p>{{ message.content }}</p>
                        </div>
                        <small class="text-muted">{{ convertUnixToTime(message.sentAt) }}</small>
                    </div>
                </div>

                <!-- Chat Footer -->
                <div v-if="selectedChat && isChatView" class="chat-footer border-top p-3">
                    <textarea v-model="newMessage" class="form-control w-100" placeholder="Type a message"
                        rows="2"></textarea>
                    <button @click="sendMessage" class="btn btn-primary mt-2 w-100">Send</button>
                </div>
            </div>
        </div>
    </div>
</template>

<script>
export default {
    data() {
        return {
            newMessage: '', // Model for the new message input
            isChatView: true,
            userId: null,
            chats: [],
            selectedChat: null,
            selectedFile: null,
            messages: [],
            userInfo: {
                id: null,
                name: null,
                username: '',
                surname: null,
                photo: null,
            },
        };
    },
    async mounted() {
        this.userId = this.$route.params.userId;
        const token = localStorage.getItem('authToken'); // Retrieve token from localStorage
        console.log('userId:', this.userId, 'token:', token);
        if (this.userId && token) {
            await this.fetchChats(token);
            this.fetchUserData(true);
        }
    },
    methods: {
        fetchUserData(firstTime = false) {
            const token = localStorage.getItem('authToken');
            // Call toggleView to switch the view before fetching the user data

            //if(!firstTime){
            //    this.toggleView();
            //} 

            // Make the API request with Authorization header
            this.$axios.get(`/users/${this.userId}`, {
                headers: {
                    'Authorization': `Bearer ${token}`
                }
            })
                .then(response => {
                    this.userInfo = response.data; // Set user info data
                    console.log('User info:', this.userInfo);
                })
                .catch(error => {
                    console.error("There was an error fetching user data:", error);
                });
        },
        toggleView() {
            this.isChatView = !this.isChatView;
        },
        convertBlobToBase64(blob) {
            return "data:image/jpeg;base64," + blob;
        },
        convertUnixToTime(unixTime) {
            // Multiply by 1000 to convert from seconds to milliseconds
            const date = new Date(unixTime);

            // Format the date as a human-readable string
            return date.toLocaleString(); // Example: "1/9/2025, 12:00:00 PM"
        },
        async fetchChats(token) {
            try {
                const response = await this.$axios.get(`/users/${this.userId}/conversations`, {
                    headers: {
                        'Authorization': `Bearer ${token}`,
                    },
                });
                this.chats = response.data;
            } catch (error) {
                console.error('Error fetching chats:', error);
            }
        },
        async fetchMessages(chatId) {
            const token = localStorage.getItem('authToken');
            try {
                const response = await this.$axios.get(`/users/${this.userId}/conversations/${chatId}/messages`, {
                    headers: {
                        'Authorization': `Bearer ${token}`,
                    },
                });
                this.messages = response.data;
            } catch (error) {
                console.error('Error fetching messages:', error);
            }
        },
        async putUsername(event) {
            event.preventDefault();  // Prevent the form from submitting
            try {
                const token = localStorage.getItem('authToken');
                const response = await this.$axios.put(`/users/${this.userId}/username`, {
                    username: this.userInfo.username,
                }, {
                    headers: {
                        'Authorization': `Bearer ${token}`,
                    },
                });
                this.userInfo = response.data;
            } catch (error) {
                console.error('Error updating username:', error);
            }
        },
        selectChat(chat) {
            this.selectedChat = chat;
            this.fetchMessages(chat.id);
        },
        handlePhotoUpload(event) {
            this.selectedFile = event.target.files[0];
        },
        async setMyPhoto(event) {   
            event.preventDefault();  // Prevent the form from submitting
            if (!this.selectedFile) {
                alert('Please select a file first.');
                return;
            }
            const token = localStorage.getItem('authToken');
            const formData = new FormData();
            formData.append('photo', this.selectedFile);

            try {
                const response = await this.$axios.post(`/users/${this.userId}/photo`, formData, {
                    headers: {
                        'Content-Type': 'multipart/form-data',
                        'Authorization': `Bearer ${token}`,
                    }
                });
                this.userInfo.photo = response.data.photo;
                console.log('File uploaded successfully:', response.data);
            } catch (error) {
                console.error('Error uploading file:', error);
            }
        }
    },
};
</script>

<style scoped>
.dashboard {
    height: 100vh;
    overflow: hidden;
}

.sidebar-header {
    background-color: #f8f9fa;
    border-bottom: 1px solid #dee2e6;
}

.chat-body {
    background-color: #f1f1f1;
}

.chat-footer {
    background-color: #ffffff;
    padding-bottom: 10px;
}

.list-group-item {
    cursor: pointer;
}

textarea {
    resize: none;
}

button {
    border-radius: 0.25rem;
}

.image {
    width: 40px;
    /* Set the width */
    height: 40px;
    /* Set the height to match the width */
    object-fit: cover;
    /* Ensures the image fills the circle without stretching */
}
</style>