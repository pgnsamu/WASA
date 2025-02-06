<template>
    <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/bootstrap-icons/font/bootstrap-icons.css">
    <div class="dashboard container-fluid d-flex flex-column" style="min-height: 100vh;">
        <div class="row flex-grow-1">
            <!-- Sidebar -->
            <div id="app" class="col-md-3 col-lg-3 bg-light border-end d-flex flex-column" style="overflow-y: auto; max-height: 100vh;">
                <div class="sidebar-header py-3">
                    <div class="d-flex align-items-center">
                        <button class="btn btn-sm me-2" @click="toggleView()">
                            <img v-if="userInfo.photo != null" :src="convertBlobToBase64(userInfo.photo)" alt="avatar" class="rounded-circle image">
                            <i v-else class="bi bi-person-circle" style="font-size: 2rem;"></i>
                        </button>
                        <div class="d-flex">
                            <h2 v-if="selectedView == 0" class="mb-0">Chats</h2>
                            <h2 v-if="selectedView == 1" class="mb-0">Profile Settings</h2>
                            <h2 v-if="selectedView == 2" class="mb-0">Creazione gruppo</h2>
                            <h2 v-if="selectedView == 3" class="mb-0">Impostazioni chat</h2>
                            <h2 v-if="selectedView == 4" class="mb-0">Inoltra a...</h2>
                        </div>
                        <button v-if="selectedView == 0" class="btn btn-sm ms-auto btn-primary" @click="changeToView(2)">nuovo gruppo</button>
                        <button v-if="selectedView == 2 || selectedView == 3" class="btn btn-sm ms-auto btn-primary" @click="changeToView(0)">x</button>
                    </div>
                    <!--create conversation-->
                    <div v-if="this.selectedView == 0 || this.selectedView == 4" class="d-flex align-items-center mt-3">
                        <input type="text" v-model="selectedUser" class="form-control me-2" placeholder="Username a cui scrivere"/>
                        <button type="button" class="btn btn-sm ms-auto btn-primary" style="border-radius: 0.4rem;" @click="newConversation(selectedUser)">Crea Chat</button>
                    </div>
                </div>
                
                <!--chat list-->
                <div v-if="selectedView == 0" class="flex-grow-1">
                    <ul class="list-group list-group-flush">
                        <li class="list-group-item" v-for="chat in chats" :key="chat.id" @click="selectChat(chat)">
                            <div class="d-flex align-items-center">
                                <img v-if="chat.photo != null" :src="convertBlobToBase64(chat.photo)" alt="avatar"
                                    class="rounded-circle me-2 image">
                                <div>
                                    <h6 class="mb-0">{{ chat.name }}</h6>
                                    <small class="text-muted">{{ chat.lastMessage }}</small>
                                </div>
                            </div>
                        </li>
                    </ul>
                </div>

                <!--profile settings-->
                <div v-if="selectedView == 1" class="flex-grow-1">
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
                            <div class="d-flex justify-content-start mt-4">
                                <button class="btn btn-danger btn-sm" @click="logout">Logout</button>
                            </div>
                        </div>
                    </form>
                </div>
                
                <!--TODO: aggiungere controlli new conversation--> 
                <div v-if="selectedView == 2" class="flex-grow-1 mb-2">
                    <div class="mb-2 mt-2 pt-2">
                        <label for="groupPhoto" class="form-label">Foto del Gruppo</label>
                        <input type="file" id="groupPhoto" class="form-control" @change="handleGroupPhotoUpload">
                    </div>
                    <div class="mb-2 mt-2 pt-2">
                        <label for="groupName" class="form-label">Nome Gruppo</label>
                        <input type="text" v-model="groupReqInfo.name" class="form-control" id="groupName" placeholder="Nome gruppo">
                    </div>
                    <div class="mb-2 mt-2 pt-2">
                        <label for="groupDescription" class="form-label">Descrizione Gruppo</label>
                        <textarea v-model="groupReqInfo.description" class="form-control" id="groupDescription" placeholder="Descrizione del gruppo"></textarea>
                    </div>
                    <!--
                    <div class="mb-2 mt-2 pt-2">
                        <label for="addMember" class="form-label">Username membro da aggiungere</label>
                        <input type="text" v-model="newParticipant" class="form-control" id="addMember" placeholder="username dell'utente">
                    </div>
                    -->
                    <div class="d-flex justify-content-between align-items-center">
                        <!--<button class="btn btn-secondary mb-3" @click="addParticipant">Aggiungi membro</button>-->
                        <button class="btn btn-primary mb-3" @click="newConversation">Crea Gruppo</button>
                    </div>
                    <ul class="list-group">
                        <li v-for="(participant, index) in participants" :key="index" class="list-group-item d-flex justify-content-between align-items-center">
                            {{ participant }}
                            <button class="btn btn-danger btn-sm" @click="removeParticipant(index)">
                                <i class="bi bi-trash"></i>
                            </button>
                        </li>
                    </ul>
                </div>
                
                <!--chat settings-->
                <div v-if="selectedView == 3" class="flex-grow-1 mb-2">
                    <div class="chat-info">
                        <div class="mb-2 mt-2 pt-2">
                            <label for="chatPhoto" class="form-label">Foto della Chat</label>
                            <input v-if="selectedChat.isGroup" type="file" id="chatPhoto" class="form-control" @change="handlePhotoUpload">
                        </div>
                        <div class="mb-2 mt-2 pt-2">
                            <img v-if="chatInfo.photo != null" :src="convertBlobToBase64(chatInfo.photo)" alt="chat photo" class="img-fluid" style="max-width: 100%; max-height: 300px;">
                        </div>
                        <div class="mb-2 mt-2 pt-2">
                            <label for="chatName" class="form-label">Nome Chat</label>
                            <input type="text" v-model="chatInfo.name" class="form-control" id="chatName" placeholder="Nome chat" :disabled="!selectedChat.isGroup">
                        </div>
                        <div v-if="selectedChat.isGroup" class="mb-2 mt-2 pt-2">
                            <label for="chatDescription" class="form-label">Descrizione Chat</label>
                            <textarea v-model="chatInfo.description" class="form-control" id="chatDescription" placeholder="Descrizione della chat" disabled></textarea>
                        </div>
                        <div v-if="selectedChat.isGroup" class="d-flex justify-content-between align-items-center">
                            <button class="btn btn-primary mb-3" @click="updateChatInfo">Salva Modifiche</button>
                        </div>
                        <div v-if="selectedChat.isGroup" class="d-flex justify-content-between align-items-center mt-3">
                            <div class="d-flex align-items-center">
                                <input type="text" v-model="newParticipant" class="form-control me-2" placeholder="Username da aggiungere">
                                <button class="btn btn-primary" @click="addMemberToGroup">Aggiungi</button>
                            </div>
                        </div>
                        <button v-if="selectedChat.isGroup" class="btn btn-danger mt-2" @click="leaveGroup">Esci dal Gruppo</button>
                        <div v-if="selectedChat.isGroup" class="mb-2 mt-2 pt-2">
                            <label class="form-label">Membri del Gruppo</label>
                            <ul class="list-group">
                                <li v-for="(member, index) in groupMembers" :key="index" class="list-group-item d-flex justify-content-between align-items-center">
                                    {{ member.username }}
                                </li>
                            </ul>
                        </div>
                        
                    </div>
                </div>

                <!-- List of chats the user participates in -->
                <div v-if="selectedView == 4" class="flex-grow-1">
                    <ul class="list-group list-group-flush">
                        <li class="list-group-item" v-for="chat in chats" :key="chat.id" @click="selectChatToForward(chat)">
                            <div class="d-flex align-items-center">
                                <img v-if="chat.photo != null" :src="convertBlobToBase64(chat.photo)" alt="avatar" class="rounded-circle me-2 image">
                                <div>
                                    <h6 class="mb-0">{{ chat.name }}</h6>
                                    <small class="text-muted">{{ chat.lastMessage }}</small>
                                </div>
                            </div>                            
                        </li>
                    </ul>
                </div>
            </div>

            <!-- Main Chat Area -->
            <div class="col-md-9 col-lg-9 chat-container">
                <div v-if="selectedChat"
                    class="chat-header border-bottom py-3 d-flex justify-content-between align-items-center">
                    <h5>{{ selectedChat.name }}</h5>
                    <button class="btn btn-sm btn-primary" @click="fetchChatInfo">Settings</button> <!--TODO: icona ingranaggio?-->
                </div>

                <!-- Chat Body -->
                <div v-if="selectedChat" class="chat-body p-3" style="overflow-y: auto;" ref="chatBody">
                    <div v-for="message in messages" :key="message.id" class="mb-3 d-flex align-items-start">
                        
                        <button v-if="message.senderId != userId" class="btn btn-sm btn-secondary me-2" @click="selectMessageToForward(message)">
                            <i class="bi bi-send" ></i>
                        </button> 
                        
                        <!--bottoni interni-->
                        <div v-if="message.senderId == userId" class="d-flex flex-column ms-auto">
                            <button class="btn btn-sm btn-secondary rounded mb-2" @click="message.answerTo != -1 ? uncommentMessage(message) : deleteMessage(message)">
                                <i v-if="message.answerTo == -1" class="bi bi-trash"></i>
                                <div v-else class="position-relative" style="width: 20px; height: 20px;">
                                    <!-- Icona di sfondo -->
                                    <i class="bi bi-reply text-white position-absolute" style="top: 50%; left: 50%; font-size: 20px; transform: translate(-50%, -50%);"></i>
                                    <!-- Icona sovrapposta -->
                                    <i class="bi bi-x text-red position-absolute" style="top: 55%; left: 50%; font-size: 30px; transform: translate(-50%, -50%);"></i>
                                </div>
                            </button> <!--delete-->
                            <button class="btn btn-sm btn-secondary" @click="selectMessage(message)">
                                <div class="position-relative" style="width: 20px; height: 20px;">
                                    <i class="bi bi-reply text-white position-absolute" style="top: 50%; left: 50%; font-size: 20px; transform: translate(-50%, -50%);"></i>
                                </div>
                            </button> <!--reply-->
                        </div>
                        <!-- MESSAGGIO -->
                        <div :id="message.id" :class="['p-2', (selectedMessage != null && message.id == selectedMessage.id) ? (message.senderId == userId ? 'bg-green-light text-black rounded ms-2' : 'bg-green-light text-black rounded')  : (message.senderId == userId ? 'bg-blue-light text-white rounded ms-2' : 'bg-light rounded')]" style="max-width: 40%;">
                            <div v-if="message.isForwarded" class="text-muted small">
                                <i class="bi bi-arrow-right"></i> Inoltrato
                            </div>
                            <div class="d-flex justify-content-between align-items-center">
                                <strong>{{ message.senderUsername }}</strong>
                            </div>
                            <div v-if="message.answerTo != -1">
                                <div class="bg-success text-white p-1 mb-2 rounded" @click="scrollToMessage(message.answerTo)">
                                    <button class="btn btn-link text-white p-0" style="text-decoration: none;">
                                        <small>{{ getSnippet(messages.find(msg => msg.id === message.answerTo)) }}</small> <!-- TODO: mmettere nel caso in cui non esista content un'icona della foto-->
                                    </button>
                                </div>
                            </div>
                            <div v-if="message.photoContent">
                                <img :src="convertBlobToBase64(message.photoContent)" alt="photo" class="img-fluid" style="max-width: 100%; max-height: 300px;" />
                            </div>
                            <p v-html="formatContent(message.content)"></p>
                            <div class="message-reactions mt-1">
                                <span v-for="reaction in message.reactions" :key="reaction.id" class="me-2 hover-container" @mouseover="showTooltip" @mouseleave="hideTooltip">
                                    {{ reaction.content }} 
                                    <div class="hover-text" :class="{ visible: isTooltipVisible }">
                                        <small>{{ reaction.sentBy }}</small>
                                    </div>
                                </span>
                                <button class="btn btn-sm btn-light" @click="handleReactionButton(message, '👍')">👍</button>
                                <button class="btn btn-sm btn-light" @click="handleReactionButton(message, '❤️')">❤️</button>
                                <button class="btn btn-sm btn-light" @click="handleReactionButton(message, '😂')">😂</button>
                            </div>
                            <div class="message-status mt-1">
                                <small class="text-muted me-2">{{ convertUnixToTime(message.sentAt) }}</small>
                                <i v-if="message.senderId == userId && message.status === 1" class="bi bi-check text-secondary"></i> <!-- Sent -->
                                <i v-if="message.senderId == userId && message.status === 2" class="bi bi-check-all text-primary"></i> <!-- Read -->
                            </div>
                        </div>

                        <button v-if="message.senderId == userId" class="btn btn-sm btn-secondary ms-2" @click="selectMessageToForward(message)">
                            <i class="bi bi-send"></i>
                        </button> <!--icona profilo?-->
                        
                        <!--bottoni interni-->
                        <div v-if="message.senderId != userId" class="d-flex flex-column ms-2">
                            <button class="btn btn-sm btn-secondary" @click="selectMessage(message)">
                                <div class="position-relative" style="width: 20px; height: 20px;">
                                    <i class="bi bi-reply text-white position-absolute" style="top: 50%; left: 50%; font-size: 20px; transform: translate(-50%, -50%);"></i>
                                </div>
                            </button> <!--reply-->
                        </div>
                    </div>
                </div>
                <!--
                <div v-if="(selectedChat != null) && (groupMembers.length <= 1) && groupMembers != null">
                    <div class="alert alert-warning" role="alert">
                        Questo gruppo è vuoto, aggiungi un membro per iniziare a chattare!
                    </div>
                </div>
                -->
                

                <!-- Chat Footer -->
                <div v-if="selectedChat" class="chat-footer border-top p-3">
                    <textarea v-model="newMessage" class="form-control w-100" placeholder="Type a message" rows="2"></textarea>
                    <div class="d-flex justify-content-between align-items-center">
                        <button v-if="selectedFile" class="btn btn-danger mt-2 me-2" @click="unselectFile">Rimuovi File</button>
                        <input type="file" id="photo" class="form-control mt-2 w-100 me-2" @change="handlePhotoUpload">
                        <button v-if="selectedMessage != null" class="btn btn-success mt-2 w-100" @click="commentMessage" :disabled="(!newMessage.trim() && selectedFile == null )">Rispondi a</button>
                        <!--<button v-else @click="sendMessage" class="btn btn-primary mt-2 w-100" :disabled="groupMembers.length <= 1">Invia</button>-->
                    <button v-else @click="sendMessage" class="btn btn-primary mt-2 w-100" :disabled="(!newMessage.trim() && selectedFile == null ) ">Invia</button>
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
            // environment variables
            selectedView: 0,
            selectedChat: null,
            chats: [],

            // newMessage variables
            answerTo: null, // TODO: servirebbe per cambiare il bottone finale nel caso viene scelta una chat a cui rispondere
            messages: [],
            newMessage: '', // Model for the new message input
            selectedFile: null,

            // userInfo variables
            userId: null,
            userInfo: {
                id: null,
                username: '',
                photo: null,
            },

            // newChat singola variables
            selectedUser: '',

            // newGroup variables
            isGroup: false,
            participants: [],
            newParticipant: '',
            groupReqInfo: {
                id: null,
                name: null,
                description: null,
                photo: null,
            },

            // risposta a messaggio
            replyToMode: false,
            selectedMessage: null,

            // fetchChatInfo
            chatInfo: null,

            // fetchGroupMembers
            groupMembers: [],

            // interval
            interval: null,

            // hovering
            isTooltipVisible: false,

            token: localStorage.getItem('authToken'),

        }
    },
    async mounted() {
        this.userId = this.$route.params.userId;
        const token = localStorage.getItem('authToken'); // Retrieve token from localStorage
        // console.log('userId:', this.userId, 'token:', token);
        if (this.userId && token) {
            await this.fetchChats();
            this.fetchUserData();
        }
        this.interval = setInterval(() => {
            const token = localStorage.getItem('authToken');
            if (this.selectedChat != null) {
            // console.log(this.userId, this.selectedChat.id);
            this.fetchMessages(this.selectedChat.id);
            //this.fetchGroupMembers();
            }
            this.fetchChats();
        }, 5000);
        
    },
    watch: {
        messages() {
            this.scrollToBottom();
        },
        selectedChat() {
            this.fetchGroupMembers();
            this.scrollToBottom();
            this.unselectFile();
        },
        selectedView(){
            if(this.selectedView == 2){
                this.isGroup = true;
            }else{
                this.isGroup = false;
            }
        }
        
    },
    methods: {
        unselectFile() {
            this.selectedFile = null;
            const fileInput = document.getElementById('photo');
            if (fileInput) {
                fileInput.value = '';
            }
        },
        showTooltip() {
            this.isTooltipVisible = true;
        },
        hideTooltip() {
            this.isTooltipVisible = false;
        },
        handleReactionButton(message, content) {
            // probabile bug se vado a rimuovere una reazione e si crea concorrenza con il fetching automatico
            if(message.reactions == null){
                this.addReaction(message, content);
                return;
            }
            const existingReaction = message.reactions.find(r => r.sentBy == this.userInfo.username);
            console.log('existingReaction:', this.userInfo.username);
            
            if (existingReaction) {
                if (existingReaction.content === content) {
                    this.removeReaction(message, existingReaction);
                } else {
                    this.removeReaction(message, existingReaction);
                    this.addReaction(message, content);
                }
            } else {
                this.addReaction(message, content);
            }
        },

        getSnippet(message) {
            return message.content.length > 20 ? message.content.substring(0, 20) + '...' : message.content;
            // TODO: cambiare nel caso sia solo foto 
        },
        scrollToMessage(messageId) {
            const messageElement = document.getElementById(messageId);
            if (messageElement) {
                const offset = messageElement.offsetTop - (this.$refs.chatBody.clientHeight / 2) + (messageElement.clientHeight / 2);
                this.$refs.chatBody.scrollTo({
                    top: offset,
                    behavior: 'smooth'
                });
            }
        },
        selectMessage(message){
            if(this.replyToMode && this.selectedMessage == message){
                this.replyToMode=false;
                this.selectedMessage=null;
            }else{
                this.replyToMode=true;
                this.selectedMessage=message
            }
        },
        selectMessageToForward(message){
            if(this.selectedView == 4 && this.selectedMessage == message){
                this.selectedView=0;
                this.selectedMessage=null;
            }else{
                this.selectedView=4;
                this.selectedMessage=message
            }
        },
        handleGroupPhotoUpload(event) {
            this.groupReqInfo.photo = event.target.files[0];
        },
        changeToView(tmpView){
            console.log(tmpView)
            if(this.selectedView != tmpView){
                this.selectedView = tmpView;
            }else{
                this.selectedView = 0;
                selectedMessage = null;
            }
        },
        addParticipant() {
            if (this.newParticipant.trim() !== '') {
                if (!this.participants.includes(this.newParticipant.trim())) { // forse far uscire un feedback a schermo
                    this.participants.push(this.newParticipant.trim()); 
                    this.newParticipant = '';
                }
            }
        },
        removeParticipant(index) {
            this.participants.splice(index, 1);
        },
        scrollToBottom() {
            this.$nextTick(() => {
                const chatBody = this.$refs.chatBody;
                if (chatBody) {
                    chatBody.scrollTo({
                        top: chatBody.scrollHeight,
                        behavior: 'smooth'
                    });
                    // chatBody.scrollTop = chatBody.scrollHeight; alternative senza animazione
                }
            });
        },
        formatContent(content) {
            return content.replace(/\n/g, '<br>');
        },
        fetchUserData() {
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
            if(this.selectedView != 0){
                this.selectedView = 0;
            } else {
                this.selectedView = 1;
            }
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
        async fetchChats() {
            const token = localStorage.getItem('authToken');
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
                console.log('Messages:', this.messages);
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
                if(error.response.data == "username già esistente\n"){
                    alert('username già in uso');
                }
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
        },
        async sendMessage() {
            const token = localStorage.getItem('authToken');
            const messageData = {
                content: this.newMessage,
                isPhoto: this.selectedFile != null ? true : false,
                photo: this.selectedFile,
            };
            if (this.selectedFile) {
                messageData.photoContent = this.selectedFile;
            }
            try {
                const response = await this.$axios.post(`/users/${this.userId}/conversations/${this.selectedChat.id}/messages`, messageData, {
                    headers: {
                        'Content-Type': 'multipart/form-data',
                        'Authorization': `Bearer ${token}`,
                    }
                });
                console.log('Message sent successfully:', response.data);
                this.newMessage = '';

                // Clear the file input
                const fileInput = document.getElementById('photo');
                if (fileInput) {
                    fileInput.value = '';
                }

                this.selectedFile = null;
                this.fetchMessages(this.selectedChat.id);
                this.fetchChats();
            } catch (error) {
                console.error('Error sending message:', error);
            }

        },
        async deleteMessage(message) {
            const token = localStorage.getItem('authToken');
            try {
                const response = await this.$axios.delete(`/users/${this.userId}/conversations/${this.selectedChat.id}/messages/${message.id}`, {
                    headers: {
                        'Authorization': `Bearer ${token}`,
                    },
                });
                console.log('Message deleted successfully:', response.data);
                this.fetchMessages(this.selectedChat.id);
            } catch (error) {
                console.error('Error deleting message:', error);
            }
        },
        async newConversation(username = '') {
            if (!this.isGroup && username.trim() === '') {
                alert('Please enter a valid username.');
                return;
            }
            if (this.isGroup && (this.groupReqInfo.name == '' ||this.groupReqInfo.photo == null)) { // possiamo aggiungere anche controllo sui partecipanti " this.participants.length < 2 || "
                alert('inserisci almeno un nome e una foto');
                return;
            }
            const token = localStorage.getItem('authToken');
            
            // Crea un nuovo oggetto FormData
            const formData = new FormData();

            // Aggiungi i dati a FormData
            formData.append('name', this.isGroup ? this.groupReqInfo.name : '');
            formData.append('isGroup', this.isGroup);

            if (this.groupReqInfo.photo) {
                formData.append('photo', this.groupReqInfo.photo);
            }

            if (this.groupReqInfo.description) {
                formData.append('description', this.groupReqInfo.description);
            }

            // Aggiungi i partecipanti se è un gruppo
            if (this.isGroup) {
                this.participants.forEach(participant => {
                    formData.append('partecipantsUsername', participant);
                });
            } else {
                formData.append('partecipantsUsername', username);
            }
            console.log('formData:', formData.get('name'), formData.get('isGroup'), formData.get('photo'), formData.get('description'), formData.getAll('partecipantsUsername[]'));


            try {
                const response = await this.$axios.post(`/users/${this.userId}/conversations`, formData, {
                    headers: {
                        'Content-Type': 'multipart/form-data',
                        'Authorization': `Bearer ${token}`,
                    }
                });
                username = '';
                this.selectedView = 0;
                this.selectedUser = ''
                if (this.isGroup) {
                    this.groupReqInfo = {
                        id: null,
                        name: null,
                        description: null,
                        photo: null,
                    };
                    this.participants = [];
                }
                this.fetchChats();
            } catch (error) {
                if(error.response.data == "utente non trovato\n"){
                    alert('utente non esistente');
                }
                console.error('Error creating chat:', error);
            }
        },
        async commentMessage(){
            const token = localStorage.getItem('authToken');
            const messageData = {
                content: this.newMessage,
                isPhoto: this.selectedFile != null ? true : false,
                photo: this.selectedFile,
            };
            if (this.selectedFile) {
                messageData.photoContent = this.selectedFile;
            }
            try {
                const response = await this.$axios.post(`/users/${this.userId}/conversations/${this.selectedChat.id}/messages/${this.selectedMessage.id}/comments`, messageData, {
                    headers: {
                        'Content-Type': 'multipart/form-data',
                        'Authorization': `Bearer ${token}`,
                    }
                });
                console.log('Message sent successfully:', response.data);
                this.newMessage = '';
                this.selectedFile = null;
                // Clear the file input
                const fileInput = document.getElementById('photo');
                if (fileInput) {
                    fileInput.value = '';
                }
                this.selectedMessage = null;
                this.fetchMessages(this.selectedChat.id);
                this.fetchChats();
            } catch (error) {
                console.error('Error sending message:', error);
            }
        },
        async uncommentMessage(selectedComment){
            const token = localStorage.getItem('authToken');
            try {
                const response = await this.$axios.delete(`/users/${this.userId}/conversations/${this.selectedChat.id}/messages/${selectedComment.answerTo}/comments/${selectedComment.id}`, {
                    headers: {
                        'Authorization': `Bearer ${token}`,
                    },
                });
                console.log('Message deleted successfully:', response.data);
                this.fetchMessages(this.selectedChat.id);
            } catch (error) {
                console.error('Error deleting message:', error);
            }
        },
        async fetchChatInfo(){
            if (this.selectedView == 3) {
                this.selectedView = 0;
                return
            }
            const token = localStorage.getItem('authToken');
            try {
                const response = await this.$axios.get(`/users/${this.userId}/conversations/${this.selectedChat.id}`, {
                    headers: {
                        'Authorization': `Bearer ${token}`,
                    },
                });
                this.chatInfo = response.data;
                this.fetchGroupMembers();
                this.changeToView(3);
                console.log('Chat info:', this.chatInfo);
            } catch (error) {
                console.error('Error fetching chat info:', error);
            }
        }, 
        async fetchGroupMembers(){
            const token = localStorage.getItem('authToken');
            try {
                const response = await this.$axios.get(`/users/${this.userId}/conversations/${this.selectedChat.id}/users`, {
                    headers: {
                        'Authorization': `Bearer ${token}`,
                    },
                });
                this.groupMembers = response.data;
                console.log('Chat info:', this.participants);
            } catch (error) {
                console.error('Error fetching chat info:', error);
            }
        },
        async leaveGroup(){
            const token = localStorage.getItem('authToken');
            try {
                const response = await this.$axios.delete(`/users/${this.userId}/conversations/${this.selectedChat.id}/users`, {
                    headers: {
                        'Authorization': `Bearer ${token}`,
                    },
                });
                console.log('Chat deleted successfully:', response.data);
                this.selectedChat = null;
                this.changeToView(0);   
                this.fetchChats();
            } catch (error) {
                console.error('Error deleting chat:', error);
            }
        },
        async addMemberToGroup(){
            const token = localStorage.getItem('authToken');
            try {
                const response = await this.$axios.post(`/users/${this.userId}/conversations/${this.selectedChat.id}/users`, {
                    username: this.newParticipant,
                }, {
                    headers: {
                        'Authorization': `Bearer ${token}`,
                    },
                });
                console.log('Chat info:', response.data);
                this.newParticipant = '';
                this.changeToView(0);   
                this.fetchGroupMembers();
            } catch (error) {
                console.log(error.response.data);
                if (error.response.data == "user not found\n" || error.response.data == "utente non trovato\n") {
                    alert('utente non esistente');
                }
                console.error('Error fetching chat info:', error);
            }
        },
        updateChatInfo(){
            if(this.chatInfo.name != this.selectedChat.name){
                this.setGroupName();
            }
            if(this.chatInfo.description != this.selectedChat.description){
                
            }
            if(this.selectedFile != null){
                this.setGroupPhoto();
            }
        },
        async setGroupName(){
            const token = localStorage.getItem('authToken');
            try {
                const response = await this.$axios.put(`/users/${this.userId}/conversations/${this.selectedChat.id}/group`, {
                    name: this.chatInfo.name,
                }, {
                    headers: {
                        'Authorization': `Bearer ${token}`,
                    },
                });
                const updatedChat = this.chats.find(chat => chat.id === response.data.id);
                if (updatedChat) {
                    updatedChat.name = response.data.name;
                }
                this.changeToView(0);
            } catch (error) {
                console.error('Error fetching chat info:', error);
            }
        },
        async setGroupPhoto(){
            const token = localStorage.getItem('authToken');
            const formData = new FormData();
            formData.append('photo', this.selectedFile);

            try {
                const response = await this.$axios.put(`/users/${this.userId}/conversations/${this.selectedChat.id}/photo`, formData, {
                    headers: {
                        'Content-Type': 'multipart/form-data',
                        'Authorization': `Bearer ${token}`,
                    },
                });
                this.changeToView(0);
                this.fetchChats();
                console.log('Chat info:', response.data);
            } catch (error) {
                console.error('Error fetching chat info:', error);
            }
        },
        selectChatToForward(chat) {
            if (this.selectedMessage) {
                this.forwardMessage(chat);
            } else {
                alert('Please select a message to forward.');
            }
        },
        async forwardMessage(chat) {
            const token = localStorage.getItem('authToken');
            try {
                const response = await this.$axios.post(`/users/${this.userId}/conversations/${this.selectedChat.id}/messages/${this.selectedMessage.id}`, {
                    targetConversationId: chat.id,
                }, {
                    headers: {
                        'Authorization': `Bearer ${token}`,
                    },
                });
                this.fetchChats();
                this.selectedMessage = null;
                this.changeToView(0);

            } catch (error) {
                console.error('Error fetching chat info:', error);
            }
        },
        async addReaction(message, reaction) {
            const token = localStorage.getItem('authToken');
            try {
                const response = await this.$axios.post(`/users/${this.userId}/conversations/${this.selectedChat.id}/messages/${message.id}/reactions`, {
                    content: reaction,
                }, {
                    headers: {
                        'Authorization': `Bearer ${token}`,
                    },
                });
                console.log('Reaction added successfully:', response.data);
                this.fetchMessages(this.selectedChat.id);
            } catch (error) {
                console.error('Error adding reaction:', error);
            }
        },
        async removeReaction(message, reaction){
            const token = localStorage.getItem('authToken');
            try {
                const response = await this.$axios.delete(`/users/${this.userId}/conversations/${this.selectedChat.id}/messages/${message.id}/reactions/${reaction.id}`, {
                    headers: {
                        'Authorization': `Bearer ${token}`,
                    },
                });
                console.log('Reaction removed successfully:', response.data);
                this.fetchMessages(this.selectedChat.id);
            } catch (error) {
                console.error('Error removing reaction:', error);
            }
        },
        async logout() {
            try {
                // Clear the token from localStorage
                localStorage.removeItem('authToken');
                clearInterval(this.interval);
                // Redirect to the login page
                this.$router.push('/');
            } catch (error) {
                console.error('Error during logout:', error);
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
.bg-blue-light {
    background-color: #77b2d9; /* Colore verde più chiaro */
}
.bg-green-light {
    background-color: #d4edda; /* Colore verde più chiaro */
}
.text-red {
    color: #ff0000; /* Colore rosso personalizzato */
  }
.bi-flip-vertical{ transform:scale(1, -1); }

.hover-container {
    position: relative;
    display: inline-block;
    cursor: pointer;
}

.hover-text {
    visibility: hidden;
    width: 150px;
    background-color: black;
    color: #fff;
    text-align: center;
    border-radius: 5px;
    padding: 5px;
    position: absolute;
    z-index: 1;
    bottom: 100%; /* Position above the element */
    left: 50%;
    transform: translateX(-50%);
    opacity: 0;
    transition: opacity 0.3s;
}

.hover-container:hover .hover-text {
    visibility: visible;
    opacity: 1;
}

/* Contenitore principale che include chat-header, chat-body e chat-footer */
.chat-container {
  display: flex;
  flex-direction: column;
  height: 100vh; /* oppure un’altezza relativa a ciò che ti serve */
}

/* Se vuoi che l'header abbia una dimensione fissa */
.chat-header {
  flex: 0 0 auto; /* oppure imposta un'altezza fissa, ad esempio height: 60px; */
}

/* Chat-body prende il 75% dello spazio residuo */
.chat-body {
  flex: 6;
  overflow-y: auto; /* in modo da gestire lo scroll */
}

/* Chat-footer prende il 25% dello spazio residuo */
.chat-footer {
  flex: 1;
}
</style>