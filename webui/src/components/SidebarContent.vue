<script>
import { getConversations, getGroups, BASE_URL } from '../services/axios'
import router from '../router'
//TODO:Fare polling
export default {
	name: 'Chats',

	data() {
		return {
			chats: [],
			userId: Number(sessionStorage.getItem('userId')),
			inputMarginTop: '0px',
			searchQuery: '',
			editProfile: false,
			pollingInterval:null,
		}
	},

	computed: {
		filteredchats() {
			if (!this.searchQuery) return this.chats
			return this.chats.filter(chat =>
				chat.name.toLowerCase().startsWith(this.searchQuery.toLowerCase())
			)
		}
	},

	methods: {
		async loadChats() {
			try {
				const conversations = await getConversations(this.userId) || []
				const groups = await getGroups(this.userId) || []
				const newChats = [...conversations, ...groups];
				if (JSON.stringify(newChats) !== JSON.stringify(this.chats)) {
					this.chats = newChats;
				}

			} catch (err) {
				console.error(err)
			}
		},

		openChat(chat) {
			if (chat.conversationId) {
				router.push({ name: 'conversation', params: { conversationId: chat.conversationId } })
			} else if (chat.groupId) {
				router.push({ name: 'group', params: { groupId: chat.groupId } })
			}
		},

		getAvatarUrl(chat) {
			// funzione helper per gestire BASE_URL e avatar/photo
			if (chat.conversationId) {
				return `${BASE_URL}/file?file=${chat.avatar.url}`
			} else if (chat.groupId) {
				return `${BASE_URL}/file?file=${chat.photo.url}`
			}
			return ''
		},
		startPollingChats() {
			this.loadChats();
			this.pollingInterval = setInterval(() => {
				this.loadChats();
			}, 2000);
		},
		stopPollingChats() {
			if (this.pollingInterval) {
				clearInterval(this.pollingInterval);
				this.pollingInterval = null;
			}
		},
	},

	mounted() {
		this.loadChats();
		this.startPollingChats();
	},
	beforeUnmount() {
		this.stopPollingChats();
	}
}
</script>

<template>
		<div class="sticky-top">
			<h2 class="p-2">Chats</h2>
			<div class="search-bar">
				<input
					type="text"
					placeholder="Search..."
					class="input-group"
					v-model="searchQuery"
				/>
			</div>
		</div>
		<div class="container-chats">
			<div
				v-for="chat in filteredchats"
				:key="chat.conversationId || chat.groupId"
				class="d-flex justify-content-start align-items-center mb-2"
			>
				<button class="btn w-100 text-start" @click="openChat(chat)">
					<img
						:src="getAvatarUrl(chat)"
						alt="Avatar"
						class="rounded-circle avatar"
						width="40"
						height="40"
					/>
					<span class="fw-bold ms-2 truncate-text">{{ chat.name }}</span>
					<br /> <!--mettere messaggio-->
					<small class="text-muted">ID: {{ chat.conversationId || chat.groupId }}</small>
				</button>
			</div>
		</div>
</template>

<style scoped>
.truncate-text {
	white-space: nowrap;
	overflow: hidden;
	text-overflow: ellipsis;
}

.avatar {
	object-fit: cover;    /* taglia l’immagine mantenendo proporzioni 100x100 */
}

.search-bar {
	display: flex;
	margin-left: 10px;
	margin-right: 10px;
	padding-bottom: 10px;
}

.container-chats {
	overflow-y: auto;
	height: calc(100vh - 120px); /* header + search */

}

</style>
