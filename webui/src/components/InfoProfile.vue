<script>


export default {
	props: {
		show:Boolean,
		username: String,
		photo: String,
	},
	emits: ['close'],
}

</script>

<template>
	<div v-if="show" class="overlay">
		<div class="action-box">
			<button class="btn btn-close" @click="$emit('close')"></button>
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

.btn{
	border-radius: 6px;
}

</style>
