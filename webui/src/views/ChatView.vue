<script>
import { BASE_URL, getConversation, getConversations, getGroups, getMessages } from "../services/axios";
import router from "../router";
//TODO:Fare bottone per azioni e send
export default {
	data() {
		return {
			errormsg: null,
			loading: false,

			messages: [],
			currentConversation: null,
			conversations: [],
			groups: [],

			username: null,
			userId: null,
			token: null,
			conversation_id: null,

			searchQuery: "",

			pollingInterval: null, // per il polling
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
				console.error("Errore fetching messages:", err);
			}
		},

		async refresh() {
			this.loading = true;
			this.errormsg = null;

			this.userId = this.userId || localStorage.getItem("userId");
			this.token = this.token || localStorage.getItem("token");
			this.conversation_id = this.conversation_id || this.$route.params.conversation_id;

			if (!this.userId || !this.token) {
				this.errormsg = "Effettua il login";
				this.loading = false;
				return;
			}

			try {
				this.conversations = (await getConversations(this.userId)) || [];
				this.groups = (await getGroups(this.userId)) || [];
				this.currentConversation = await getConversation(this.userId, this.conversation_id);

				await this.fetchMessages(); // primo caricamento messaggi
			} catch (err) {
				this.errormsg = err.toString();
			} finally {
				this.loading = false;
			}
		},

		startPolling() {
			this.fetchMessages(); // fetch iniziale
			this.pollingInterval = setInterval(this.fetchMessages, 2000); // ogni 2 secondi
		},

		stopPolling() {
			if (this.pollingInterval) clearInterval(this.pollingInterval);
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
};
</script>

<template>
	<div class="chat-header-wrapper">
		<div class="d-flex justify-content-between align-items-center">
			<div class="d-flex align-items-center gap-2">
				<img
					:src="`${BASE_URL()}/file?file=${this.currentConversation.avatar.url}`"
					alt="Avatar"
					class="rounded-circle align-self-center"
					width="45"
					height="45"
				/>
				<h1 class="mb-0 align-bottom">{{ this.currentConversation.name }}</h1>
			</div>
		</div>
	</div>

	<div class="chat-body">
		<div class="messages-container" ref="messageContainer">
			<div
				v-for="msg in messages"
				:key="msg.message_id"
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
					<div class="d-flex justify-content-between">
						<small v-if="msg.sender.userId !== userId" class="sender">{{ msg.sender.username }}</small>
						<small v-if="msg.sender.userId === userId" class="sender"></small>
						<small v-if="!msg.isForwarded" class="text-end">(forwarded)</small>
					</div>
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
		<input class="input-group" type="text">
		<button class="btn">Send</button>
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



</style>
