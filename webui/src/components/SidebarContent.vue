<script setup>
import {ref, onMounted, computed} from 'vue'
import {getConversations, getGroups,BASE_URL} from '../services/axios'
import router from '../router'

const chats = ref([])
const userId = localStorage.getItem('userId')
const inputMarginTop = ref('0px')
const searchQuery = ref('')

async function loadChats() {
	try {
		const conversations = await getConversations(userId) || []
		const groups = await getGroups(userId) || []
		chats.value = [...conversations, ...groups]
	} catch (err) {
		console.error(err)
	}
}

function openChat(chat) {
	if(chat.conversation_id) {
		router.push({ name: 'conversation', params: { conversation_id: chat.conversation_id } });
	} else if(chat.group_id) {
		router.push({ name: 'group', params: { group_id: chat.group_id } });
	}
}

onMounted(() => {
	loadChats()
})

const filteredchats = computed(() => {
	if(!searchQuery.value) return chats.value
	return chats.value.filter(chat =>
		chat.name.toLowerCase().startsWith(searchQuery.value.toLowerCase()))
})

</script>

<template>
	<div>
		<h2 class="d-flex flex-column p-lg-2 ">Chats</h2>
		<input type="text" placeholder="Search..." class="input-group" :style="{ marginTop: inputMarginTop }" v-model="searchQuery">
	</div>
	<div>
		<div v-for="chat in filteredchats" :key="chat.conversation_id || chat.group_id" class="d-flex justify-content-start align-items-center mb-2">
			<button class="btn w-100 text-start" @click="openChat(chat)">
				<img
					:src="chat.conversation_id ? `${BASE_URL}/file?file=${chat.avatar.url}` : `${BASE_URL}/file?file=${chat.photo.url}`"
					alt="Avatar"
					class="rounded-circle"
					width="40"
					height="40"
				/>
				<span class="fw-bold ms-2 truncate-text">{{ chat.name }}</span>
				<br>
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
