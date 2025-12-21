<script>
import { BASE_URL, getConversation, getConversations, getGroups, getMessages,postMessage } from "../services/axios";
import router from "../router";
import LoadingSpinner from "../components/LoadingSpinner.vue";
import Comment from "../components/Comment.vue";
import Delete from "../components/Delete.vue";
import Forward from "../components/Forward.vue";


//TODO:Fare bottone per azioni
export default {
	components: {Forward, Delete, Comment, LoadingSpinner},
	data() {
		return {
			errormsg: null,
			errorTimeout:null,
			loading: false,

			messages: [],
			currentConversation: null,

			username: null,
			userId: null,
			token: null,
			conversation_id: null,

			searchMessage: "",

			pollingInterval: null, // per il polling

			refreshDebounced: null, //per evitare troppe richieste al server

			openMessageOptions:null,

			selectedMessageId : null,

			showForward: false,
			forwardMessageId: null,
			deleteMessageId: null,
			deleteMessage: false,
			showComments: false,
			commentMessageId: null,
		};
	},
	methods: {
		router() {
			return router;
		},
		BASE_URL() {
			return BASE_URL;
		},

		async fetchMessages() {
			try {
				const msgs = await getMessages(this.userId, this.conversation_id) || [];

				const container = this.$refs.messageContainer;
				let isAtBottom = false;

				if (container) {
					// verifica se siamo già in fondo (tolleranza 20px)
					isAtBottom = container.scrollTop + container.clientHeight >= container.scrollHeight - 20;
				}

				this.messages = msgs;
				this.$nextTick(() => {
					if (container && isAtBottom) {
						container.scrollTop = container.scrollHeight;
					}
				});

			} catch (err) {
				if(err.response?.data?.includes("database is locked")) {
					setTimeout(() => this.fetchMessages(), 2000); // ritenta dopo 2 secondi
				} else {
					console.error("Errore fetching messages:", err);
				}

			}
		},

		async refresh() {
			if(this.refreshDebounced)   clearTimeout(this.refreshDebounced);

			this.refreshDebounced = setTimeout(async() => {
			this.loading = true;
			this.errormsg = null;


			this.userId = this.userId || localStorage.getItem("userId");
			this.token = this.token || localStorage.getItem("token");
			this.conversation_id = this.conversation_id || this.$route.params.conversation_id;

			if (!this.userId || !this.token) {
				this.errormsg = "Do login!";
				this.loading = false;
				return;
			}

			try {
				this.currentConversation = await getConversation(this.userId, this.conversation_id);

				await this.fetchMessages(); // primo caricamento messaggi
			} catch (err) {
				this.errormsg = err.toString();
			} finally {
				this.loading = false;
			}},300); //300ms tra una richiesta e l'altra
		},

		startPolling() {
			this.fetchMessages(); // fetch iniziale
			this.pollingInterval = setInterval(this.fetchMessages, 2000); // ogni 2 secondi
		},

		stopPolling() {
			if (this.pollingInterval) {
				clearInterval(this.pollingInterval);
				this.pollingInterval = null;
			}
		},
		async sendMessageButton() {
			const textInput = this.$refs.messageText;
			const photoInput = this.$refs.messagePhoto;

			const text = textInput.value;
			const photo = photoInput.files[0];

			// controllo estensione foto
			if (photo && photo.type !== "image/png" && photo.type !== "image/jpeg") {
				this.showError("Only PNG or JPEG!");
				photoInput.value = "";
				return;
			}

			if (text || photo) {
				try {
					this.messages = await postMessage(this.userId, this.conversation_id, text, photo,"direct");
					textInput.value = "";
					photoInput.value = "";
					await this.fetchMessages();
				} catch (e) {
					console.error("Error send message:", e);
					this.showError("Error sending message");
				}
			}
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
		showOptions(messageId){
			if (this.openMessageOptions === messageId) {
				this.openMessageOptions = null;
			} else {
				this.openMessageOptions = messageId;
			}
		},
		forward(messageId) {
			this.forwardMessageId = messageId;
			this.showForward = true;
		},
		comment(messageId){
			this.showComments = true;
			this.commentMessageId = messageId;
		},
		deleteM(messageId){
			this.deleteMessage = true;
			this.deleteMessageId = messageId;
		},
	},
	created() {
		this.username = localStorage.getItem("username");
		this.token = localStorage.getItem("token");
		this.userId = localStorage.getItem("userId");
		this.conversation_id = this.$route.params.conversation_id;

		if (this.token && this.userId) {
			this.refresh();
			this.startPolling();
		}
	},
	beforeUnmount() {
		this.stopPolling();
	},
	watch: {
		'$route.params.conversation_id'(newId) {
			this.stopPolling();
			this.conversation_id = newId;
			this.refresh();
			this.startPolling();
		}
	},
	computed: {
		filteredMessages() {
			let result = this.messages

			if(this.searchMessage.trim().length !== ""){
				const q = this.searchMessage.toLowerCase();
				result=result.filter(chat => chat.body.text.toLowerCase().includes(q));
			}
			return result;
		}
	}
};
</script>

<template>
	<LoadingSpinner v-if="loading ===true"></LoadingSpinner>
	<Comment
			:userId="userId"
			:show="showComments"
			:messageId="commentMessageId"
			:chatId="conversation_id"
			:type="`direct`"
			@close="showComments=false"
	/>
	<Delete
			:userId="userId"
			:show="deleteMessage"
			:messageId="deleteMessageId"
			:chatId="conversation_id"
			:type="`direct`"
			@close="deleteMessage=false"
	/>
	<Forward
		:userId="userId"
		:show="showForward"
		:messageId="forwardMessageId"
		:chatId="conversation_id"
		:type="`direct`"
		@close="showForward = false"
	/>
	<div class="chat-header-wrapper">
		<div class="d-flex justify-content-between align-items-center">
			<div class="d-flex align-items-center gap-2 p-lg-2">
				<img
					:src="`${BASE_URL()}/file?file=${this.currentConversation.avatar.url}`"
					alt="Avatar"
					class="rounded-circle align-self-center"
					width="45"
					height="45"
				/>
				<h1 class="mb-0 align-bottom">{{ this.currentConversation.name }}</h1>
			</div>
			<div>
				<input type="text" placeholder="Search message..." v-model="searchMessage">
			</div>
		</div>
	</div>

	<div class="chat-body">
		<div class="messages-container" ref="messageContainer">
			<ErrorMsg v-if="errormsg" :msg="errormsg" />
			<div
				v-for="msg in filteredMessages"
				:key="msg.messageId"
				:class="['message-row gap-3', msg.sender.userId === userId ? 'mine' : 'theirs']"
			>
				<img v-if="msg.sender.userId !== userId"
					:src="`${BASE_URL()}/file?file=${this.currentConversation.avatar.url}`"
					alt="Avatar"
					class="rounded-circle align-self-end"
					width="35"
					height="35"
				/>
				<div class="message-bubble">
					<div class="d-flex justify-content-end"
						 v-if="openMessageOptions === msg.messageId">
						<button class="btn btn-outline-secondary" @click="forward(msg.messageId)"><img src="../icons/share-icon_4662621.png" alt="Forward" width="25" height="25"></button>
						<button class="btn btn-outline-primary" @click="comment(msg.messageId)"><img src="../icons/chat-dots-fill.svg" alt="comment" width="23" height="23"></button>
						<button v-if="msg.sender.userId === userId" class="btn btn-outline-danger" @click="deleteM(msg.messageId)"><img src="../icons/trash3-fill.svg" alt="Delete" width="23" height="23"></button>
					</div>
					<div class="d-flex justify-content-between">
						<small v-if="msg.sender.userId !== userId" class="sender">{{ msg.sender.username }}</small>
						<small v-if="msg.sender.userId === userId" class="sender"></small>
						<div class="d-flex">
						<small v-if="msg.isForwarded" class="text-end">(forwarded)</small>
						<button class="icon-btn" @click="showOptions(msg.messageId)"><img src="../icons/dots_16164512.png" width="15" height="15" alt="Options"></button>
						</div>
					</div>
					<img v-if="msg.body.photo && msg.body.photo.url" :src="`${BASE_URL()}/file?file=${msg.body.photo.url}`" alt="PhotoMessage" class="message-photo">
					<p class="text">{{ msg.body.text }}</p>
					<div class="justify-content-between">
						<span>{{msg.time}}</span>
						<img v-if="msg.read==='received'" src="../icons/icons8-double-tick-100.png" width="15" height="15" alt="received">
						<img v-if="msg.read==='read'" src="../icons/icons8-double-tick-100-2.png" width="15" height="15" alt="read">
					</div>
				</div>
			</div>
		</div>
		<div class ="d-flex justify-content-between">
			<input type="file" ref="messagePhoto">
			<input class="input-group" type="text" placeholder="Write message..." ref="messageText">
			<button class="btn btn-outline-dark" @click="sendMessageButton">Send</button>
		</div>
	</div>


</template>

<style scoped>
.chat-header-wrapper {
	margin-top: 1rem;
}

.chat-body {
	height: calc(100vh - 140px); /* header + margini */
	display: flex;
	flex-direction: column;
}

.messages-container {
	flex: 1;
	overflow-y: auto;
	padding: 1rem;
	background-color: #f8f9fa;
	border-radius: 20px;
	flex-direction: column-reverse;
	display: flex;
}


.sender {
	font-weight: bold;
	font-size: 0.85rem;
}

.text {
	margin: 0;
}

.message-row {
	display: flex;
	margin-bottom: 0.75rem;
}

.message-row.mine {
	justify-content: flex-end;
}

.message-row.theirs {
	justify-content: flex-start;
}

.message-bubble {
	position: relative;
	max-width: 70%;
	padding: 0.6rem 0.9rem;
	border-radius: 14px;
	background-color: #ffffff;
	box-shadow: 0 1px 3px rgba(0,0,0,0.1);
	word-wrap: break-word;
}

/* colori diversi */
.message-row.mine .message-bubble {
	background-color: #7cd65a;
	color: white;
	border-bottom-right-radius: 4px;
}

.message-row.theirs .message-bubble {
	background-color: #9adf97;
	color: #212529;
	border-bottom-left-radius: 4px;
}

.sender {
	font-size: 0.7rem;
	opacity: 0.7;
	display: block;
	margin-bottom: 2px;
}

.message-photo {
	max-width: 300px;   /* larghezza massima */
	max-height: 300px;  /* altezza massima */
	width: auto;        /* scala proporzionalmente */
	height: auto;
	border-radius: 10px;
	object-fit: contain; /* mantiene proporzioni */
}

.icon-btn {
	padding-bottom: 5px;
	border: none;
	background: none;
	display: inline-flex;
	align-items: center;
	justify-content: center;
}



</style>
