<script>

import { getConversations, getGroups, forwardMessage } from '../services/axios'

export default {
	name: 'ForwardMessage',

	props: {
		userId: Number,
		show: Boolean,
		messageId: Number,
		chatId: Number,
		type: String
	},

	data() {
		return {
			chats: [],
			selected: new Set()
		}
	},

	methods: {
		async loadChats() {
			const conversations = await getConversations(this.userId) || []
			const groups = await getGroups(this.userId) || []
			this.chats = [...conversations, ...groups]
		},

		toggle(chat) {
			const id = chat.conversation_id || `g-${chat.group_id}`
			this.selected.has(id)
				? this.selected.delete(id)
				: this.selected.add(id)
		},

		confirmForward() {
			const conversations = []
			const groups = []

			this.chats.forEach(chat => {
				const id = chat.conversation_id || `g-${chat.group_id}`
				if (this.selected.has(id)) {
					chat.conversation_id
						? conversations.push(chat.conversation_id)
						: groups.push(chat.group_id)
				}
			})

			forwardMessage(
				this.userId,
				this.chatId,
				this.messageId,
				this.type,
				conversations,
				groups
			)

			this.$emit('close')
		}
	},

	mounted() {
		this.loadChats()
	},

	watch: {
		show(newVal) {
			if (newVal) {
				this.selected.clear()
			}
		}
	}
}
</script>

<template>
	<div v-if="show" class="overlay">
		<div class="action-box">
			<h4 class="text-center">Forward message</h4>
			<div
				v-for="chat in chats"
				:key="chat.conversation_id || chat.group_id"
				class="chat-item"
				@click="toggle(chat)"
			>
				<input
					type="checkbox"
					:checked="selected.has(chat.conversation_id || `g-${chat.group_id}`)"
				>
				<span class="ms-2">{{ chat.name }}</span>
			</div>

			<div class="actions">
				<button class="btn btn-secondary" @click="$emit('close')">Cancel</button>
				<button class="btn btn-primary" @click="confirmForward">Forward</button>
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

.chat-item {
	cursor: pointer;
	padding: 8px 4px;
	display: flex;
	align-items: center;
}

.chat-item:hover {
	background-color: #f0f0f0;
}

.actions {
	display: flex;
	justify-content: space-between;
	gap: 10px;
	margin-top: 10px;
}
</style>
