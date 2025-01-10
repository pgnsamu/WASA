import {createRouter, createWebHashHistory} from 'vue-router'
import LoginView from '../views/loginView.vue'
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
