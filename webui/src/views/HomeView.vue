<script>
import {BASE_URL, getConversations, getGroups} from "../services/axios";
import router from "../router";
import {nextTick} from "vue";
import LoadingSpinner from "../components/LoadingSpinner.vue";

export default {
	components: {LoadingSpinner},
	data() {
		return {
			errormsg: null,
			loading: false,

			conversations: [],
			groups: [],

			username: null,
			userId: null,
			token: null,

			activeFilter: "all",
			sidebarOpen: false,

			searchQuery: ""
		}
	},
	methods: {
		router() {
			return router;
		},
		BASE_URL() {
			return BASE_URL;
		},
		async refresh() {
			this.loading = true;
			this.errormsg = null;
			if(!this.userId) this.userId = sessionStorage.getItem('userId');
			if(!this.token) this.token = sessionStorage.getItem('token');

			if(!this.userId || !this.token){
				this.errormsg="Do the Login"
				this.loading=false;
				return;
			}
			try{
				let convResponse = await getConversations(this.userId);
				this.conversations = convResponse || [];

				let groupResponse = await getGroups(this.userId);
				this.groups = groupResponse || [];
			}catch(err){
				this.errormsg=err.toString();
			}
			this.loading = false;
		},
		openChat(id, type){
			if (type === "direct"){
				router.push({ name: 'conversation', params: { conversation_id: id } });
			} else if(type === "group"){
				router.push({ name: 'group', params: { group_id: id } });
			} else console.error("type must be direct or group");
		},
		toggleSidebar() {
			this.sidebarOpen = !this.sidebarOpen;
		},
		Logout(){
			router.push('/');
		},
		UserList(){
			router.push('/users');
		},
		goToParticipants(groupId) {
			router.push({name: 'participants', params: {group_id: groupId}});
		}
	},
	created() {
		this.username = sessionStorage.getItem("username");
		this.token = sessionStorage.getItem("token");
		this.userId = sessionStorage.getItem("userId");
		this.sidebarOpen= false;
		if (this.token && this.userId) {
			this.refresh();
		}
	},
	computed: {
		filteredChats() {
			// unisce conversations e groups
			let result = [...this.conversations, ...this.groups];

			// filtro per tipo
			if(this.activeFilter === "direct") {
				result = result.filter(chat => chat.conversation_id);
			} else if(this.activeFilter === "group") {
				result = result.filter(chat => chat.group_id);
			}

			// filtro per ricerca
			if(this.searchQuery.trim() !== "") {
				const q = this.searchQuery.toLowerCase();
				result = result.filter(chat =>
					chat.name.toLowerCase().startsWith(q)
				);
			}

			return result;
		}
	},
	watch: {
		sidebarOpen(newVal) {
			if (newVal) {
				nextTick(() => {
					const navbarHeight = document.querySelector('header.navbar')?.offsetHeight || 0;
					const sidebar = document.querySelector('.sidebar');
					if(sidebar){
						sidebar.style.top = `${navbarHeight}px`;
						sidebar.style.height = `calc(100% - ${navbarHeight}px)`;
					}
				});
			}
		}
	}
}
</script>

<template>
	<LoadingSpinner v-if="loading ===true"></LoadingSpinner>
	<div>
		<div class="row align-items-center pt-3 pb-2 mb-3 border-bottom">
			<div class="col-3 d-flex justify-content-start">
				<div>
					<button class="hamburger" @click="toggleSidebar">
						<span></span>
						<span></span>
						<span></span>
					</button>

					<div class="sidebar" :class="{ open: sidebarOpen }">
						<button class="btn btn-outline-primary w-100" @click="doSomething">Edit Profile</button>
						<button class="btn btn-outline-primary w-100" @click="UserList">User List</button>
						<button class="btn w-100" id="Logout" @click="Logout">Logout</button>
					</div>

					<div class="overlay" v-if="sidebarOpen" @click="toggleSidebar"></div>
				</div>
			</div>

			<div class="col-6 text-center">
				<h1 class="h2 mb-0">Welcome Back, {{ username }}!</h1>
			</div>

			<div class="col-3 d-flex justify-content-end">
				<div class="btn-toolbar mb-2 mb-md-0">
					<button type="button" class="btn btn-sm btn-outline-secondary" @click="refresh">
						Refresh
					</button>
				</div>
			</div>
		</div>

		<ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
		<LoadingSpinner v-if="loading"></LoadingSpinner>

		<div v-if="!loading" class="row">
			<div class="d-flex justify-content-between mb-3">
				<div class="d-flex gap-2">
					<button
						class="btn btn-outline-primary"
						:class="{ selected: activeFilter === 'all'}"
						@click="activeFilter = 'all'"
					>All</button>
					<button
						class="btn btn-outline-primary"
						:class="{ selected: activeFilter === 'direct'}"
						@click="activeFilter = 'direct'"
					>Chats</button>
					<button
						class="btn btn-outline-primary"
						:class="{ selected: activeFilter === 'group'}"
						@click="activeFilter = 'group'"
					>Groups</button>
				</div>
				<div class="d-flex justify-content-end gap-2">
					<input type="text" class="form-control" placeholder="Search..." v-model="searchQuery">
					<button class="btn btn-primary" @click="CreateChat">
						<img src="../icons/new-svgrepo-com.svg" width="25" height="25"  alt="New"/>
					</button>
				</div>
			</div>

			<div v-if="filteredChats.length === 0" class="alert alert-secondary">
				No chats found.
			</div>

			<ul class="list-group shadow-sm">
				<li v-for="chat in filteredChats" :key="chat.conversation_id || chat.group_id" class="list-group-item d-flex justify-content-between align-items-center py-3">
					<div>
						<img
							:src="chat.conversation_id ? `${BASE_URL()}/file?file=${chat.avatar.url}` : `${BASE_URL()}/file?file=${chat.photo.url}`"
							alt="Avatar"
							class="rounded-circle"
							width="40"
							height="40"
						/>
						<span class="fw-bold ms-2">{{ chat.name }}</span>
						<span v-if="chat.group_id" class="participant text-decoration-none" @click="goToParticipants(chat.group_id)" style="cursor: pointer;">
							(
							<span v-for="(user, index) in chat.participants.slice(0, 10)" :key="user.userId">
								{{ user.username }}<span v-if="index < Math.min(chat.participants.length, 10) - 1">, </span>
							</span>
							<span v-if="chat.participants.length > 10">, ...</span>
							)
						</span>
						<br>
						<small class="text-muted">ID: {{ chat.conversation_id || chat.group_id }}</small>
					</div>
					<button class="btn btn-outline-primary btn-sm" @click="openChat(chat.conversation_id || chat.group_id, chat.conversation_id ? 'direct' : 'group')">
						Open
					</button>
				</li>
			</ul>

		</div>
	</div>
</template>

<style scoped>
/* mantiene lo stile esistente */
.hamburger {
	display: flex;
	flex-direction: column;
	justify-content: space-around;
	width: 30px;
	height: 25px;
	background: transparent;
	border: none;
	cursor: pointer;
	padding: 0;
	z-index: 1100;
}
.selected {
	font-weight: bold;
	background-color: #0d55c7;
	color: white;
	border-color: #0d6efd;
}
.hamburger span {
	display: block;
	width: 100%;
	height: 3px;
	background-color: #333;
	border-radius: 2px;
}
.sidebar {
	position: fixed;
	left: -260px;
	width: 260px;
	height: 100%;
	background-color: #ffffff;
	padding: 1.5rem 1rem;
	box-shadow: 2px 0 10px rgba(0,0,0,0.2);
	transition: left 0.3s ease;
	z-index: 1000;
	display: flex;
	flex-direction: column;
	gap: 0.75rem;
}
.sidebar.open {
	left: 0;
}
.overlay {
	position: fixed;
	top: 0;
	left: 0;
	width: 100%;
	height: 100%;
	background: rgba(0,0,0,0.3);
	z-index: 999;
}
.btn {
	transition: background-color 0.2s ease, color 0.2s ease, border-color 0.2s ease;
}
#Logout {
	background-color: red;
	color: white;
	border-radius: 5px;
	font-weight: bold;
	margin-top: auto;
}
.participant{
	color: rgba(128, 128, 128, 0.38);
	transition: background-color 0.2s ease, color 0.2s ease, border-color 0.2s ease;
}
.participant:hover {
	color: rgb(108, 108, 108);
	transition: background-color 0.2s ease, color 0.2s ease, border-color 0.2s ease;
}
</style>
