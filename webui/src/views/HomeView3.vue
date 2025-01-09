<template>
    <div class="dashboard container-fluid">
        <div class="row">
            <!-- Sidebar -->
            <div class="col-md-3 col-lg-2 bg-light border-end">
                <div class="sidebar-header py-3">
                    <h5>User Name</h5>
                </div>
                <ul class="list-group list-group-flush">
                    <li class="list-group-item" v-for="chat in chats" :key="chat.id" @click="selectChat(chat)">
                        <div class="d-flex align-items-center">
                            <img :src="convertBlobToBase64(chat.photo)" alt="avatar" class="rounded-circle me-2 image"
                                width="40">
                            <div>
                                <h6 class="mb-0">{{ chat.name }}</h6>
                                <small class="text-muted">{{ chat.lastMessage }}</small>
                            </div>
                        </div>
                    </li>
                </ul>
            </div>

            <!-- Main Chat Area -->
            <div class="col-md-9 col-lg-10">
                <div v-if="selectedChat"
                    class="chat-header border-bottom py-3 d-flex justify-content-between align-items-center">
                    <h5>{{ selectedChat.name }}</h5>
                    <div>
                        <button class="btn btn-light me-2">Search</button>
                        <button class="btn btn-light">More</button>
                    </div>
                </div>
                <div v-if="selectedChat" class="chat-body p-3" style="height: calc(100vh - 120px); overflow-y: auto;">
                    <div v-for="message in messages" :key="message.id" class="mb-3">
                        <div :class="['p-2', message.senderId == userId ? 'bg-primary text-white' : 'bg-light']">
                            {{ message.content }}
                        </div>
                        <small class="text-muted">{{ convertUnixToTime(message.sentAt) }}</small>
                    </div>
                </div>
                <div v-if="selectedChat" class="chat-footer border-top p-3">
                    <textarea v-model="newMessage" class="form-control w-100" placeholder="Type a message2"
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
            userId: null,
            chats: [],
            selectedChat: null,
            messages: [],
        };
    },
    async mounted() {
        this.userId = this.$route.params.userId;
        const token = localStorage.getItem('authToken'); // Retrieve token from localStorage
        console.log('userId:', this.userId, 'token:', token);
        if (this.userId && token) {
            await this.fetchChats(token);
        }
    },
    methods: {
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
                const response = await fetch(`/users/${this.userId}/conversations`, {
                    headers: {
                        'Authorization': `Bearer ${token}`,
                    },
                });
                const data = await response.json();
                this.chats = data;
            } catch (error) {
                console.error('Error fetching chats:', error);
            }
        },
        async fetchMessages(chatId) {
            const token = localStorage.getItem('authToken');
            try {
                const response = await fetch(`/users/${this.userId}/conversations/${chatId}/messages`, {
                    headers: {
                        'Authorization': `Bearer ${token}`,
                    },
                });
                const data = await response.json();
                this.messages = data;
            } catch (error) {
                console.error('Error fetching messages:', error);
            }
        },
        selectChat(chat) {
            this.selectedChat = chat;
            this.fetchMessages(chat.id);
        },
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
  padding-bottom: 10px; /* Ensure some space at the bottom */
}

.list-group-item {
  cursor: pointer;
}

textarea {
  resize: none; /* Prevent resizing */
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