<template>
	<div class="p-4 text-white" style="min-height: calc(100vh - 3.5rem)">
		<div class="mb-4">
			<h2 class="fw-semibold">Find Users</h2>
			<div class="input-group">
				<input
					type="text"
					class="form-control bg-dark text-white"
					placeholder="Search by username..."
					v-model="searchQuery"
				/>
			</div>
		</div>

		<div v-if="errorMsg" class="alert alert-danger">{{ errorMsg }}</div>

		<div v-if="isLoading" class="text-muted">Loading users...</div>

		<div
			v-if="!isLoading && filteredUsers.length === 0 && searchQuery"
			class="text-muted"
		>
			No users found matching your search
		</div>
		<div
			v-if="!isLoading && filteredUsers.length === 0 && !searchQuery"
			class="text-muted"
		>
			No users found
		</div>

		<div class="user-list">
			<UserSearchEntry
				v-for="user in filteredUsers"
				:key="user.username"
				:username="user.username"
				:profilePicture="user.photo"
				@start-chat="handleStartChat"
			/>
		</div>
	</div>
</template>

<script>
import UserSearchEntry from "../components/UserSearchEntry.vue";
import { inject } from "vue";

export default {
	inject: ["refreshTrigger"],
	components: { UserSearchEntry },
	data() {
		return {
			searchQuery: "",
			allUsers: [],
			isLoading: false,
			errorMsg: "",
			defaultProfileImage:
				"https://cdn.pixabay.com/photo/2015/10/05/22/37/blank-profile-picture-973460_1280.png",
		};
	},
	computed: {
		filteredUsers() {
			if (!this.searchQuery.trim()) return this.allUsers;
			const query = this.searchQuery.toLorCase();
			return this.allUsers.filter((user) =>
				user.username.toLowerCase().includes(query)
			);
		},
	},
	methods: {
		async loadAllUsers() {
			this.isLoading = true;
			this.errorMsg = "";
			try {
				// Fetch all usernames
				const response = await this.$axios.get("/users", {
					headers: { Authorization: `Bearer ${this.$token.value}` },
				});

				// Fetch photos for all users in parallel
				this.allUsers = await Promise.all(
					response.data.map(async (username) => {
						try {
							const photoResponse = await this.$axios.get(
								`/users/photo/${username}`,
								{
									headers: {
										Authorization: `Bearer ${this.$token.value}`,
									},
								}
							);
							return {
								username,
								photo:
									photoResponse.data.photo ||
									this.defaultProfileImage,
							};
						} catch (error) {
							return {
								username,
								photo: this.defaultProfileImage,
							};
						}
					})
				);
			} catch (error) {
				this.errorMsg = "Failed to load users";
				console.error(error);
			} finally {
				this.isLoading = false;
			}
		},

		async handleStartChat(username) {
			try {
				// Clear previous errors
				this.errorMsg = "";

				// 1. Create conversation
				const response = await this.$axios.post(
					"/conversations/create",
					{ participants: [username] },
					{
						headers: {
							Authorization: `Bearer ${this.$token.value}`,
						},
						validateStatus: (status) => status < 500, // Handle all status codes < 500 as non-errors
					}
				);

				// Handle non-2xx responses
				if (response.status < 200 || response.status >= 300) {
					throw new Error(
						response.data?.message ||
							"Failed to create conversation"
					);
				}

				// 2. Refresh conversations list
				if (this.$root.refreshTrigger !== undefined) {
					this.$root.refreshTrigger++;
				}

				console.log(
					"Conversation created with ID: ",
					response.data.conversationId
				);
				let conversationId = response.data.conversationId;

				// 3. Navigate to chat with forced reload
				try {
					await this.$router.push(`/chat/${conversationId}`);
				} catch (error) {
					console.error("Failed to navigate to chat", error);
					this.errorMsg = "Failed to navigate to chat";
				}
			} catch (error) {
				// Unified error handling
				this.errorMsg = error.message;
				console.error("Error chain:", {
					message: error.message,
					response: error.response?.data,
					stack: error.stack,
				});

				// Rethrow to Vue's error handler
				console.error(error);
				return Promise.reject(error);
			}
		},
	},
	mounted() {
		this.loadAllUsers();
	},
};
</script>

<style scoped>
.user-list {
	display: flex;
	flex-direction: column;
	gap: 0.5rem;
}

.form-control {
	border-color: #495057;
}

.form-control:focus {
	background-color: #212529;
	border-color: #6c757d;
	box-shadow: none;
}
</style>
