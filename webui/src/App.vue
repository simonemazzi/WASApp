<script setup>
import { RouterLink, RouterView, useRoute } from 'vue-router'
const route = useRoute()
</script>
<script>
export default {}
</script>

<template>
	<header v-if="route.path !== '/' && route.path !== '/session'" class="navbar navbar-dark sticky-top bg-dark flex-md-nowrap p-0 shadow">
	</header>

	<div class="container-fluid">
		<div class="row">
			<nav v-if="route.path !== '/' && route.path !== '/session'" id="sidebarMenu" class="col-md-3 col-lg-2 d-md-block bg-light sidebar collapse">
			</nav>

			<main :class="['px-md-4', (route.path === '/' || route.path === '/session') ? 'col-12' : 'col-md-9 ms-sm-auto col-lg-10']">
				<RouterView />
			</main>
		</div>
	</div>
</template>

<style>
</style>
