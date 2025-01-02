import {createRouter, createWebHashHistory} from 'vue-router'
import HomeView from '../views/HomeView.vue'
import ChatList from '../views/chatList.vue'

const router = createRouter({
	history: createWebHashHistory(import.meta.env.BASE_URL),
	routes: [
		{path: '/chat/20', component: ChatList},
		{path: '/link1', component: ChatList},
		{path: '/link2', component: ChatList},
		{path: '/some/:id/link', component: ChatList},
	]
})

export default router
