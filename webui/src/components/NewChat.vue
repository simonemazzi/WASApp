<script>
import router from "../router";
import {BASE_URL, createConversation, getConversations, getGroups, getUsers} from "../services/axios";


export default {
	name: "NewChat",
	props: {
		show: Boolean
	},
	data(){
		return {
			errormsg: null,

			users:[],
			conversations: [],

			isGroup: false,

			pollingInterval: null, // per il polling

			refreshDebounced: null, //per evitare troppe richieste al server
			searchUsers: "",

		}
	},
	methods:{
		router(){
			return router;
		},
		BASE_URL(){
			return BASE_URL;
		},
		async fetchData() {
			try {
				// fetch utenti e conversazioni
				this.users = await getUsers();
				this.conversations = await getConversations(sessionStorage.getItem('userId'));
			} catch (err) {
				this.showError("Errore fetching data: " + err.toString());
			}
		},
		getConversationWith(username) {
			// ritorna la conversazione con quell'utente, se esiste
			return this.conversations.find(conv => {
				return conv.name.includes(username);
			});
		},

		startPolling(){
			this.fetchData();
			this.pollingInterval = setInterval(this.fetchData,2000); // 2 secondi
		},
		showError(msg) {
			this.errormsg = msg;

			// cancella eventuale timeout precedente
			if (this.errorTimeout) clearTimeout(this.errorTimeout);

			// setta il nuovo timeout
			this.errorTimeout = setTimeout(() => {
				this.errormsg = null;
				this.errorTimeout = null;
			}, 5000);
		},

		stopPolling(){
			if (this.pollingInterval) clearInterval(this.pollingInterval);
		},
		async goToConversation(username) {
			let conversation = this.getConversationWith(username);

			if (!conversation) {
				try {
					conversation = await createConversation(
						sessionStorage.getItem('userId'),
						username
					);
					this.fetchData()
				} catch (e) {
					console.error(e);
					return;
				}
			}

			await router.push({
				name: 'conversation',
				params: {conversation_id: conversation.conversation_id}
			});
		},
		createGroupBegin(){
			this.isGroup = !this.isGroup;
		},

	},
	beforeUnmount() {
		this.stopPolling();
	},
	watch: {
		show(newVal) {
			if (newVal) {
				this.startPolling();
			} else {
				this.stopPolling();
			}
		}
	},
	computed: {
		filteredUsers() {
			if (!this.searchUsers) {
				return this.users.filter(
					u => u.user_id !== sessionStorage.getItem('userId')
				);
			}

			const search = this.searchUsers.toLowerCase();

			return this.users.filter(user =>
				user.user_id !== sessionStorage.getItem('userId') &&
				user.username.toLowerCase().includes(search)
			);
		}
	},
}

</script>

<template>

	<div v-if="this.$props.show" class="overlay">
		<div class="action-box">
			<h4 class="text-center">New Chat</h4>
			<input type="text" v-model="this.searchUsers" placeholder="Search user..." class="input-group">
			<div class="d-flex pt-2">
				<button
					:class="['btn', 'h-25', isGroup ? 'btn-danger' : 'btn-outline-primary']"
					:style="{ width: !isGroup ? '100%' : '15%' }"
					@click="createGroupBegin">
					<span v-if="!isGroup">Create Group</span>
					<img v-else src="../icons/reject.png" alt="Cancel" width="25" height="25" />
				</button>
				<input v-if="this.isGroup" type="text" placeholder="Name group" class="w-100"/>
			</div>
			<div v-if="!isGroup"class="users-box">
				<div v-for="user in filteredUsers" :key="user.user_id" class=" d-flex justify-content-between pt-2">
					<span>{{user.username}}</span>
					<button v-if="this.getConversationWith(user.username) === undefined" class="btn btn-success" @click="goToConversation(user.username)">Crea</button>
					<button v-else class="btn btn-success" @click="goToConversation(user.username)">Apri</button>
				</div>
			</div>
			<div v-else class="users-box">
				<div v-for="user in filteredUsers" :key="user.user_id" class=" d-flex justify-content-between pt-2">
					<span>{{user.username}}</span>

				</div>
			</div>

			<div class="actions">
				<button class="btn btn-secondary" @click="$emit('close')">Cancel</button>
				<button v-if="isGroup"class="btn btn-success" @click="createGroup">Create</button>
			</div>
		</div>
	</div>

</template>

<style scoped>

.overlay {
	position: fixed;
	inset: 0;
	background: rgba(0,0,0,0.5);
	display: flex;
	justify-content: center;
	align-items: center;
	z-index: 1200;
}

.action-box {
	background: white;
	padding: 20px;
	width: 400px;
	max-height: 80vh;
	overflow-y: auto;
	border-radius: 8px;
	box-shadow: 0 4px 15px rgba(0,0,0,0.3);
	z-index: 1201;
	display: flex;
	flex-direction: column;
}

.actions {
	display: flex;
	justify-content: space-between;
	gap: 10px;
	margin-top: 10px;
}

.users-box > div {
	display: flex;
	justify-content: space-between;
	align-items: center;
	overflow-y: auto;
}

</style>
