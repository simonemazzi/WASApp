<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { doLogin} from "../services/axios";

const router = useRouter();
const username = ref('');
const errorMessage = ref(null);

const login = async () => {
	errorMessage.value = null;
	try{
		const userData = await doLogin(username.value);

		localStorage.setItem('username', username.value);
		localStorage.setItem('userId', userData.userId);
		localStorage.setItem('token', userData.token);

		router.push('/home');
	} catch(err){
		errorMessage.value = err.message;
	}
};


</script>

<template>
	<div class="d-flex justify-content-center align-items-center vh-100">
		<div class="card p-4 shadow-sm" style="max-width: 400px; width: 100%;">
			<h2 class="text-center mb-4">Login</h2>

			<div class="mb-3">
				<label for="username" class="form-label">Username</label>
				<input
					type="text"
					class="form-control"
					id="username"
					v-model="username"
					placeholder="Inserisci il tuo nome"
					@keyup.enter="login"
				>
			</div>

			<div v-if="errorMessage" class="alert alert-danger" role="alert">
				{{ errorMessage }}
			</div>

			<button
				type="button"
				class="btn btn-primary w-100"
				@click="login"
				:disabled="!username"
			>
				Entra
			</button>
		</div>
	</div>
</template>

<style>
</style>
