<template>
    <div>
      <ul>
        <li v-for="chat in chats" :key="chat.id" @click="selectChat(chat.id)" class="p-4 border-b">
          <h3>{{ chat.name }}</h3>
          <p>{{ chat.lastMessage }}</p>
        </li>
      </ul>
    </div>
</template>
  
  <script>
  export default {
    data() {
      return {
        chats: [],
      };
    },
    created() {
      this.fetchChats();
    },
    methods: {
      async fetchChats() {
        // const token = localStorage.getItem('jwtToken');
        const token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MzU5MjU0MzEsImlkIjoxLCJ1c2VybmFtZSI6ImJiYiJ9.2iT7vDzGFAnMYNKi1qmTMRgwecQOHdbIBCNsnd1tDAo";
        const response = await fetch('/users/1/conversations', {
          headers: {
            'Authorization': `Bearer ${token}`
          }
        });
        this.chats = await response.json();
      },
      selectChat(chatId) {
        this.$emit('chat-selected', chatId);
      },
    },
  };
  </script>