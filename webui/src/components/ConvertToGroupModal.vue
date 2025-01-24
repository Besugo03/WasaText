<template>
	<div class="convert-modal text-white" v-if="show">
		<h4>Add Participants</h4>
		<input
			type="text"
			class="form-control bg-dark text-white mb-3"
			placeholder="Search users..."
			v-model="searchQuery"
		/>

		<div class="user-list">
			<UserSearchEntry
				v-for="user in filteredUsers"
				:key="user.username"
				:username="user.username"
				:profilePicture="user.photo"
				@start-chat="toggleSelection"
			/>
		</div>

		<button
			class="btn btn-primary mt-3"
			@click="convert"
			:disabled="selectedParticipants.length === 0"
		>
			Add {{ selectedParticipants.length }} Participants
		</button>
		<button class="btn btn-secondary mt-3 ms-2" @click="$emit('close')">
			Cancel
		</button>
	</div>
</template>

<script>
import UserSearchEntry from "./UserSearchEntry.vue";

export default {
	components: { UserSearchEntry },
	props: ["show"],
	data() {
		return {
			searchQuery: "",
			allUsers: [],
			selectedParticipants: [],
		};
	},
	computed: {
		filteredUsers() {
			return this.allUsers.filter(
				(u) =>
					u.username
						.toLowerCase()
						.includes(this.searchQuery.toLowerCase()) &&
					!this.selectedParticipants.includes(u.username)
			);
		},
	},
	methods: {
		toggleSelection(username) {
			if (this.selectedParticipants.includes(username)) {
				this.selectedParticipants = this.selectedParticipants.filter(
					(u) => u !== username
				);
			} else {
				this.selectedParticipants.push(username);
			}
		},
		async convert() {
			try {
				await this.$axios.put(
					"/conversations/adduser",
					{
						conversationId: this.$route.params.id,
						participants: this.selectedParticipants,
					},
					{
						headers: {
							Authorization: `Bearer ${this.$token.value}`,
						},
					}
				);
				this.$emit("close");
				this.$parent.refreshTrigger++;
			} catch (error) {
				console.error("Failed to add participants:", error);
			}
		},
	},
	async mounted() {
		const response = await this.$axios.get("/users", {
			headers: { Authorization: `Bearer ${this.$token.value}` },
		});
		this.allUsers = response.data;
	},
};
</script>
