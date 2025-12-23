<script>
import router from "../router";
import {BASE_URL, createConversation, createGroup, getConversations, getGroups, getUsers} from "../services/axios";


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

			selectedUsers: new Set(), //per i selezionati del gruppo
			groupName:"",

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
			if(!this.isGroup){
				this.selectedUsers.clear();
			}
		},
		toggleUser(user){
			const id = user.username;
			this.selectedUsers.has(id) ? this.selectedUsers.delete(id) : this.selectedUsers.add(id);
		},

		async newGroup(){
			const members=[...this.selectedUsers];
			const name = this.groupName;
			if (name !== ""){
				try{
					const group= await createGroup(sessionStorage.getItem('userId'), members, name);
					await router.push({
						name: 'group',
						params: {group_id: group.group_id}
					});
				}catch(e){
					console.error(e);
				}
			}else{
				this.showError("Name can't be empty");
			}

		}

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
			<ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
			<h4 class="text-center">New Chat</h4>
			<input type="text" v-model="this.searchUsers" placeholder="Search user..." class="input-group">
			<div class="top-controls pt-2">
				<button
					:class="['btn', 'h-25', isGroup ? 'btn-danger' : 'btn-outline-primary']"
					:style="{ width: !isGroup ? '100%' : '15%' }"
					@click="createGroupBegin">
					<span v-if="!isGroup">Create Group</span>
					<img v-else src="../icons/reject.png" alt="Cancel" width="25" height="25" />
				</button>
				<input v-if="this.isGroup" type="text" placeholder="Name group" class="w-100" v-model="groupName"/>
			</div>
			<div class="users-box">
				<div v-for="user in filteredUsers" :key="user.user_id" class="user-row" @click="isGroup && toggleUser(user)">
					<div class="user-left">
						<input v-if="isGroup" type="checkbox" class="selected" :checked="selectedUsers.has(user.username)">
					</div>
					<div class="user-name">{{ user.username }}</div>
					<div class="user-right">
						<button v-if="!isGroup" class="btn btn-success btn-sm" @click.stop="goToConversation(user.username)">{{ getConversationWith(user.username) ? 'Open' : 'Create' }}</button>
					</div>

				</div>
			</div>

			<div class="actions">
				<button class="btn btn-secondary" @click="$emit('close')">Cancel</button>
				<button v-if="isGroup" class="btn btn-success" @click="newGroup">Create</button>
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

.selected {
	appearance: none;
	-webkit-appearance: none;
	width: 18px;
	height: 18px;
	border: 2px solid #0d6efd;
	border-radius: 4px;
	cursor: pointer;
	display: flex;
	align-items: center;
	justify-content: center;
	transition: all 0.2s ease;
	background-color: white;
}

.selected:checked {
	background-color: #0d6efd;
	border-color: #0d6efd;
}

.selected:checked::after {
	content: "✓";
	color: white;
	font-size: 14px;
	font-weight: bold;
	line-height: 1;
}

.user-row {
	display: flex;
	align-items: center;
	padding: 8px 6px;
	border-radius: 6px;
	cursor: pointer;
}

.user-row:hover {
	background-color: #f0f4ff;
}

.user-left {
	width: 26px;
	display: flex;
	justify-content: center;
}

.user-name {
	flex: 1;
	min-width: 0;
	white-space: nowrap;
	overflow: hidden;
	text-overflow: ellipsis;
}

.user-right {
	width: 80px;
	display: flex;
	justify-content: flex-end;
}


.top-controls {
	display: flex;
	gap: 8px;
	width: 100%;
}


.top-controls > input {
	flex: 1; /* prende lo spazio rimanente */
}

</style>
