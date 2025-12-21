import {createRouter, createWebHashHistory} from 'vue-router'
import HomeView from '../views/HomeView.vue'
import LoginView from '../views/LoginView.vue'
import ChatView from '../views/ChatView.vue'
import ParticipantsGroupView from "../views/ParticipantsGroupView.vue";
import GroupView from "../views/GroupView.vue";

const router = createRouter({
	history: createWebHashHistory(import.meta.env.BASE_URL),
	routes: [
		{path: '/', component: LoginView},
		{path: '/home', component: HomeView},
		{	path: '/conversation/:conversation_id',
			name: 'conversation',
			component: ChatView},
		{	path: '/group/:group_id/participants',
			name: 'participants',
			component: ParticipantsGroupView},
		{	path: '/group/:group_id',
			name: 'group',
			component: GroupView},
	]
})

export default router
