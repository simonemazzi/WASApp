<script>


import {BASE_URL, getConversations, getGroups} from "../services/axios";

export default {
	data: function() {
		return {
			errormsg: null,
			loading: false,

			conversations: [],
			groups: [],

			username:null,
			userId: null,
			token:null
		}
	},
	methods: {
		BASE_URL() {
			return BASE_URL
		},
		async refresh() {
			this.loading = true;
			this.errormsg = null;
			if(!this.userId) this.userId = localStorage.getItem('userId');
			if(!this.token) this.token = localStorage.getItem('token');

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
		openChat(id,type){
			console.log("CHAT :", id,type);
			// TODO: FARE CHAT
		}
	},
	created() {
		this.username = localStorage.getItem("username");
		this.token = localStorage.getItem("token");
		this.userId = localStorage.getItem("userId");
		if (this.token && this.userId) {
			this.refresh();
		}
	}
}


</script>

<template>
	<div>
		<div class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
			<h1 class="h2">Bentornato, {{ username }}!</h1>
			<div class="btn-toolbar mb-2 mb-md-0">
				<button type="button" class="btn btn-sm btn-outline-secondary" @click="refresh">
					Aggiorna
				</button>
			</div>
		</div>

		<ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>

		<LoadingSpinner v-if="loading"></LoadingSpinner>

		<div v-if="!loading" class="row">

			<div class="col-md-6">
				<h3 class="mt-4 mb-3 text-primary">Chat Private</h3>

				<div v-if="conversations.length === 0" class="alert alert-secondary">
					Nessuna conversazione attiva.
				</div>

				<ul class="list-group shadow-sm">
					<li v-for="conv in conversations" :key="conv.conversation_id" class="list-group-item d-flex justify-content-between align-items-center">
						<div>
							<img
								:src="`${BASE_URL()}/file?file=${conv.avatar.url}`"
								alt="Avatar"
								class="rounded-circle"
								width="40"
								height="40"
							/>
							<span class="fw-bold">{{ conv.name }}</span>

							<br>
							<small class="text-muted">ID: {{ conv.conversation_id }}</small>
						</div>
						<button class="btn btn-outline-primary btn-sm" @click="openChat(conv.conversation_id, 'direct')">
							Apri
						</button>
					</li>
					<!--
					<div v-if="groups.length === 0" class="alert alert-secondary">
						Non sei in nessun gruppo.
					</div>
					-->
					<ul class="list-group shadow-sm">
						<li v-for="grp in groups" :key="grp.groupId" class="list-group-item d-flex justify-content-between align-items-center">
							<div>
								<span class="fw-bold">{{ grp.name }}</span>
								<br>
								<small class="text-muted">ID: {{ grp.groupId }}</small>
							</div>
							<button class="btn btn-outline-success btn-sm" @click="openChat(grp.groupId, 'group')">
								Entra
							</button>
						</li>
					</ul>
				</ul>
			</div>



		</div>
	</div>
</template>

<style>
</style>
