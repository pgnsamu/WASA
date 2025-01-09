import {createRouter, createWebHashHistory} from 'vue-router'
import HomeView from '../views/HomeView.vue'
import ChatList from '../views/chatList.vue'
import LoginView from '../views/loginView.vue'
import HomeView3 from '../views/HomeView3.vue';
import HomeView4 from '../views/HomeView4.vue';

const routes = [
	{
		path: '/',
		name: 'LoginView',
		component: LoginView
	},
	{
		path: '/users/:userId',
		name: 'HomeView',
		component: HomeView4
	}
];
  
	const router = createRouter({
		history: createWebHashHistory(import.meta.env.BASE_URL),
		routes
	});

export default router
