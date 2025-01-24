<template>
	<div class="user-entry d-flex align-items-center p-3 bg-dark rounded">
		<img
			:src="profilePicture"
			class="rounded-circle me-3"
			alt="Profile picture"
			width="50"
			height="50"
			style="object-fit: cover"
		/>
		<div class="flex-grow-1">
			<h5 class="mb-0 text-light">{{ username }}</h5>
		</div>
		<button
			class="btn btn-primary"
			@click="$emit('start-chat', username)"
			:disabled="isCurrentUser"
		>
			{{ isCurrentUser ? "That's you!" : "Start Chatting" }}
		</button>
	</div>
</template>

<script>
export default {
	props: {
		username: {
			type: String,
			required: true,
		},
		profilePicture: {
			type: String,
			default:
				"https://cdn.pixabay.com/photo/2015/10/05/22/37/blank-profile-picture-973460_1280.png",
		},
	},
	computed: {
		isCurrentUser() {
			return this.username === this.$username.value;
		},
	},
};
</script>

<style scoped>
.user-entry {
	transition: background-color 0.2s ease;
}

.user-entry:hover {
	background-color: rgba(255, 255, 255, 0.05) !important;
}
</style>
