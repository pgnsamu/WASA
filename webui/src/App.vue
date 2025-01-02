<script setup>
import { RouterLink, RouterView } from 'vue-router'
</script>
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
}
</script>

<template>
	<header class="navbar navbar-dark sticky-top bg-dark flex-md-nowrap p-0 shadow">
		<a class="navbar-brand col-md-3 col-lg-2 me-0 px-3 fs-6" href="#/">WASACHAT</a>
		<button class="navbar-toggler position-absolute d-md-none collapsed" type="button" data-bs-toggle="collapse" data-bs-target="#sidebarMenu" aria-controls="sidebarMenu" aria-expanded="false" aria-label="Toggle navigation">
			<span class="navbar-toggler-icon"></span>
		</button>
	</header>
	<div class="container-fluid">
		<div class="row">
			<nav id="sidebarMenu" class="col-md-3 col-lg-2 d-md-block bg-light sidebar collapse">
				<div class="position-sticky pt-3 sidebar-sticky">
					<h6 class="sidebar-heading d-flex justify-content-between align-items-center px-3 mt-4 mb-1 text-muted text-uppercase">
						<span>Chats</span>
					</h6>
					<ul class="nav flex-column">
						<li v-for="chat in chats" :key="chat.id" @click="selectChat(chat.id)" class="nav-item">
							<!--<a href="#" class="nav-link" @click.prevent="selectChat(chat.id)">-->
							<RouterLink :to="'/chat/' + chat.id" class="nav-link">
								<svg class="feather"><use href="/feather-sprite-v4.29.0.svg#message-circle"/></svg>
								{{ chat.name }}
							</RouterLink>
						</li>
					</ul>
				</div>
			</nav>
			<main class="col-md-9 ms-sm-auto col-lg-10 px-md-4">
				<RouterView />
			</main>
		</div>
	</div>
</template>
<style>
</style>
