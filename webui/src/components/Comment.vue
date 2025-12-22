<script setup>
import { ref , onMounted,watch} from 'vue'
import {commentMessage,getComments} from "../services/axios";

const props = defineProps({
	userId: Number,
	show: Boolean,
	messageId: Number,
	chatId: Number,
	type: String
})

const emit = defineEmits(['close', 'confirm'])
const emojiInput = ref('')
const userId = props.userId

const comments = ref([])

// Funzione per lasciare solo emoji
function filterEmoji() {
	// Regex che prende tutti i simboli emoji
	emojiInput.value = Array.from(emojiInput.value)
		.filter(char => /\p{Emoji}/u.test(char))
		.join('')
}

function confirmComment() {
	if (emojiInput.value.trim()) {
		commentMessage(userId,props.chatId,props.messageId,props.type,emojiInput.value)
		emojiInput.value = ''
	}
}

function insertEmoji(emoji) {
	emojiInput.value = emoji
	filterEmoji()
}


async function loadComments()  {
	if (!userId || !props.chatId || !props.messageId) return;
	const comm = await getComments(userId, props.chatId, props.messageId, props.type) || []
	comments.value = [...comm]
}

onMounted(() => {
	if (props.messageId) loadComments()
})

watch(() => props.messageId, (newId) => {
	if (newId) loadComments()
})
</script>

<template>
<div v-if="show" class="overlay">
	<div class="action-box">
		<h2 class="text-center">Send Comment</h2>
		<div class="comment-box">
			<div v-for="comment in comments" :key="comment.id" class=" d-flex justify-content-between">
				<p >{{ comment.sender.username }}</p>
				<p >{{comment.emoji}}</p>
			</div>
		</div>
		<input
			type="text"
			v-model="emojiInput"
			@input="filterEmoji"
			placeholder="Comment with an emoji..."
		/>
		<div class="d-flex justify-content-center">
			<button class="btn" @click="insertEmoji(`🕶`)">🕶</button>
			<button class="btn" @click="insertEmoji(`✨`)">✨</button>
			<button class="btn" @click="insertEmoji(`😍`)">😍</button>
			<button class="btn" @click="insertEmoji(`💕`)">💕</button>
		</div>
		<div class="actions">
			<button class="btn btn-secondary" @click="emit('close')">Cancel</button>
			<button class="btn btn-primary" @click="confirmComment">Comment</button>
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
	overflow-y: auto;      /* scroll se troppe chat */
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

.comment-box {
	background: rgba(0, 0, 0, 0.11);
	padding: 20px;
	overflow-y: auto;
	display: flex;
	flex-direction: column;
	border-top-right-radius: 10px;
	border-top-left-radius: 10px;
}

</style>
