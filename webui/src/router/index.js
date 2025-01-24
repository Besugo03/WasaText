import {createRouter, createWebHashHistory} from 'vue-router'
import HomeView from '../views/HomeView.vue'
import ChatView from '../views/ChatView.vue'
import WelcomeView from '../views/WelcomeView.vue'
import ProfileView from '../views/ProfileView.vue'
import FindUsersView from '../views/FindUsersView.vue'
import CreateConversation from '../views/CreateConversationView.vue'

const router = createRouter({
	history: createWebHashHistory(import.meta.env.BASE_URL),
	routes: [
		{path: '/', component: HomeView},
		{path: '/settings', component: ProfileView},
		{path: "/group/settings/:id", component: ProfileView},
		{path: '/link2', component: FindUsersView},
		{path: '/create-conversation', component: CreateConversation},
		{path: '/some/:id/link', component: HomeView},
		{path: '/welcomescreen', component: WelcomeView},
		{path: '/chat/:id', component: ChatView},
	]
})

export default router
