<script>

import {addToGroup, BASE_URL, deleteFromGroup, getGroup, getUserInfo} from "../services/axios";
import router from "../router";

export default {
	data(){
		return {
			group:null,
			users:[],
			editMode:false,
			addUserMode:false,
		}
	},
	props:{
		show: Boolean,
		group_id:Number,
		user_id:Number,
	},
	emits:['close'],
	methods:{
		router(){
			return router;
		},
		BASE_URL(){
			return BASE_URL;
		},
		async refresh(){
			try{
				this.group = await getGroup(this.user_id,this.group_id);
				this.users = this.group.participants;
				this.users = await Promise.all(
					this.group.participants.map(async user => {
						const userInfo = await getUserInfo(String(user.userId));
						return {
							...user,
							avatar: userInfo.avatar
						};
					})
				);
			}catch(error){

				console.error(error);
			}

		},
		closePanel(){
			this.editMode = false;
			this.$emit('close');
		},
		leaveGroup(){
			try{
				deleteFromGroup(this.user_id,this.group_id);
			}catch(error){
				console.error(error);
				return;
			}
			router.push('/home');
		},
		addMember(username){
			try{
				addToGroup(this.user_id,this.group_id,username);
			}catch (error){
				console.error(error);

			}
		},
		addUserLayout(){
			this.addUserMode = !this.addUserMode;
		}
	},
	mounted(){
		this.refresh();
	}
}
//TODO: FARE ADD USER LAYOUT CON BOTTONI E FUNZIONI
</script>

<template>
	<div v-if="show" class="overlay">
		<div class="action-box">
			<div class="header d-flex align-items-center justify-content-between">
				<button class="btn btn-close" @click="closePanel"></button>
				<h4 class="h4">Info Group</h4>
				<button class="btn" @click="editMode=true" style="border-radius: 50px; padding: 0; height: 45px; width: 45px;">
					<img src="../icons/edit.png" alt="Edit" width="20" height="20" class="mb-1"/>
				</button>
			</div>
			<div class="d-flex align-items-center justify-content-center flex-column gap-3">
				<img class="avatar rounded-circle" :src="`${BASE_URL()}/file?file=${group.photo.url}`" width="200" height="200" alt="Photo"/>
				<h4 class="fw-bold text-center">{{group.name}}</h4>
			</div>

			<div class="d-flex flex-column gap-3 mx-2 participants-wrapper">
				<div class="d-flex justify-content-between align-items-center">
					<h2 class="h2">Participants</h2>
					<button class="btn d-flex align-items-center justify-content-center"
							@click="addUserLayout"
							style="border-radius: 50px; padding: 0; height: 45px; width: 45px;">
						<img src="../icons/plus.png" alt="Add" width="20" height="20"/>
					</button>
				</div>

				<div class="participants-list ">
					<div v-for="user in users" :key="user.userId" class="pb-3 d-flex gap-2  align-items-center">
						<img class="avatar rounded-circle" :src="`${BASE_URL()}/file?file=${user.avatar.Url}`" width="50" height="50" alt="Photo" />
						<span class="text-muted">{{ user.username }}</span>
					</div>
				</div>
			</div>

			<div class="d-flex justify-content-between">
				<button class="btn btn-danger" @click="leaveGroup">Leave Group</button>
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


.header {
	height: 40px; /* riferimento comune */
}




.avatar {
	object-fit: cover;    /* taglia l’immagine mantenendo proporzioni 100x100 */
}

.participants-wrapper {
	flex: 1;
	display: flex;
	flex-direction: column;
	overflow: hidden; /* evita overflow esterno */
}

.participants-list {
	overflow-y: auto;          /* scroll verticale */
	max-height: 300px;         /* altezza massima per scroll */
	padding-right: 5px;        /* evita che il scrollbar tocchi il testo */
}



</style>
