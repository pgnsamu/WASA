<template>
    <div>
      <ul>
        <li
          v-for="chat in chats"
          :key="chat.id"
          @click="selectChat(chat.id)"
          class="p-4 border-b"
        >
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
        const response = await fetch('/users/1/conversations'); // Replace with your backend API URL.
        this.chats = await response.json();
      },
      selectChat(chatId) {
        this.$emit('chat-selected', chatId);
      },
    },
  };
  </script>