<script>
import { getConversations, getGroups, BASE_URL } from '../services/axios'
import router from '../router'

export default {
	name: 'Chats',

	data() {
		return {
			chats: [],
			userId: sessionStorage.getItem('userId'),
			inputMarginTop: '0px',
			searchQuery: ''
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
				this.chats = [...conversations, ...groups]
			} catch (err) {
				console.error(err)
			}
		},

		openChat(chat) {
			if (chat.conversation_id) {
				router.push({ name: 'conversation', params: { conversation_id: chat.conversation_id } })
			} else if (chat.group_id) {
				router.push({ name: 'group', params: { group_id: chat.group_id } })
			}
		},

		getAvatarUrl(chat) {
			// funzione helper per gestire BASE_URL e avatar/photo
			if (chat.conversation_id) {
				return `${BASE_URL}/file?file=${chat.avatar.url}`
			} else if (chat.group_id) {
				return `${BASE_URL}/file?file=${chat.photo.url}`
			}
			return ''
		}
	},

	mounted() {
		this.loadChats()
	}
}
</script>

<template>
	<div>
		<h2 class="d-flex flex-column p-lg-2">Chats</h2>
		<input
			type="text"
			placeholder="Search..."
			class="input-group"
			:style="{ marginTop: inputMarginTop }"
			v-model="searchQuery"
		/>
	</div>
	<div>
		<div
			v-for="chat in filteredchats"
			:key="chat.conversation_id || chat.group_id"
			class="d-flex justify-content-start align-items-center mb-2"
		>
			<button class="btn w-100 text-start" @click="openChat(chat)">
				<img
					:src="getAvatarUrl(chat)"
					alt="Avatar"
					class="rounded-circle"
					width="40"
					height="40"
				/>
				<span class="fw-bold ms-2 truncate-text">{{ chat.name }}</span>
				<br />
				<small class="text-muted">ID: {{ chat.conversation_id || chat.group_id }}</small>
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
</style>
