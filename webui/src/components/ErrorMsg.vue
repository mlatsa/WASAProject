<script setup>
</script>

<script>
export default {
	emits: ["errorDismissed"],
	props: ['errormsg']
}
</script>

<template>
	<div v-if="errormsg.length" class="alert-container alert-danger" role="alert">
		<div style="position: relative;top: 0.6em;">
			<span class="error-header">{{ errormsg }}</span>
			<button class="error-button" @click="$emit('errorDismissed', '')">Dismiss</button>
		</div>
	</div>
</template>
