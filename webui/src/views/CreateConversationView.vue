<template>
	<div class="p-4 text-white" style="min-height: calc(100vh - 3.5rem)">
		<div class="mb-4">
			<h2 class="fw-semibold">Create New Conversation</h2>
			<div class="input-group mb-3">
				<input
					type="text"
					class="form-control bg-dark text-white"
					placeholder="Search users..."
					v-model="searchQuery"
				/>
			</div>
		</div>

		<div v-if="errorMsg" class="alert alert-danger">{{ errorMsg }}</div>
		<div v-if="successMsg" class="alert alert-success">
			{{ successMsg }}
		</div>

		<div v-if="isLoading" class="text-muted">Loading users...</div>

		<div class="selected-users mb-4">
			<h5>Selected Participants:</h5>
			<div v-if="selectedParticipants.length === 0" class="text-muted">
				No participants selected yet.
				<div class="fw-bold">
					Creating a conversation
					<div class="d-inline text-warning">now</div>
					will result in a conversation with yourself.
				</div>
			</div>
			<div
				v-if="selectedParticipants.length === 1"
				class="text-muted mb-2"
			>
				One participant selected.
				<div class="fw-bold">
					Creating a conversation
					<div class="d-inline text-warning">now</div>
					will result in a normal chat and not a group.
				</div>
			</div>
			<div class="d-flex flex-wrap gap-2">
				<div
					v-for="user in selectedParticipants"
					:key="user.username"
					class="badge bg-primary p-2 d-flex align-items-center gap-2"
				>
					{{ user.username }}
					<button
						class="btn btn-sm btn-danger"
						@click="toggleSelection(user)"
					>
						&times;
					</button>
				</div>
			</div>
		</div>

		<div
			v-if="!isLoading && filteredUsers.length === 0 && searchQuery"
			class="text-muted"
		>
			No users found matching your search
		</div>

		<div class="user-list">
			<div
				v-for="user in filteredUsers"
				:key="user.username"
				class="user-entry d-flex align-items-center p-3 bg-dark rounded mb-2"
			>
				<img
					:src="user.photo"
					class="rounded-circle me-3"
					alt="Profile picture"
					width="50"
					height="50"
					style="object-fit: cover"
				/>
				<div class="flex-grow-1">
					<h5 class="mb-0 text-light">{{ user.username }}</h5>
				</div>
				<button
					class="btn"
					:class="isSelected(user) ? 'btn-danger' : 'btn-primary'"
					@click="toggleSelection(user)"
				>
					{{ isSelected(user) ? "Remove" : "Add" }}
				</button>
			</div>
		</div>

		<div class="mt-4">
			<button
				class="btn btn-success"
				:class="selectedParticipants.length === 0 ? 'btn-warning' : ''"
				:disabled="isCreating"
				@click="createConversation"
			>
				<span
					v-if="isCreating"
					class="spinner-border spinner-border-sm"
				></span>
				{{
					isCreating
						? "Creating..."
						: selectedParticipants.length === 0
						? "Create Chat With Yourself"
						: selectedParticipants.length === 1
						? "Create Chat"
						: "Create Group Chat"
				}}
			</button>
		</div>
	</div>
</template>

<script>
export default {
	data() {
		return {
			searchQuery: "",
			allUsers: [],
			selectedParticipants: [],
			isLoading: false,
			isCreating: false,
			errorMsg: "",
			successMsg: "",
			defaultProfileImage:
				"https://cdn.pixabay.com/photo/2015/10/05/22/37/blank-profile-picture-973460_1280.png",
		};
	},
	watch: {
		"$username.value": {
			immediate: true,
			handler() {
				this.loadAllUsers();
			},
		},
	},
	computed: {
		filteredUsers() {
			return this.allUsers.filter(
				(user) =>
					user.username
						.toLowerCase()
						.includes(this.searchQuery.toLowerCase()) &&
					// exclude yourself from the list
					user.username.toLowerCase() !==
						this.$username.value.toLowerCase()
			);
		},
	},
	methods: {
		async loadAllUsers() {
			this.isLoading = true;
			this.errorMsg = "";
			try {
				const response = await this.$axios.get("/users", {
					headers: { Authorization: `Bearer ${this.$token.value}` },
				});

				this.allUsers = [];

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
				if (this.$username.value == "") {
					this.errorMsg =
						"You are not logged in. Please log in from the 'Change User' section to view your profile.";
					return;
				}
				this.errorMsg = "Failed to load users";
				console.error(error);
			} finally {
				this.isLoading = false;
			}
		},

		isSelected(user) {
			return this.selectedParticipants.some(
				(u) => u.username === user.username
			);
		},

		toggleSelection(user) {
			if (this.isSelected(user)) {
				this.selectedParticipants = this.selectedParticipants.filter(
					(u) => u.username !== user.username
				);
			} else {
				this.selectedParticipants.push(user);
			}
		},

		async createConversation() {
			this.isCreating = true;
			this.errorMsg = "";
			this.successMsg = "";

			try {
				const participants = this.selectedParticipants.map(
					(u) => u.username
				);

				const response = await this.$axios.post(
					"/conversations/create",
					{
						participants,
						isGroup: participants.length > 1,
					},
					{
						headers: {
							Authorization: `Bearer ${this.$token.value}`,
						},
					}
				);

				this.successMsg = "Conversation created successfully!";

				// Refresh conversations list
				if (this.$root.refreshTrigger !== undefined) {
					this.$root.refreshTrigger++;
				}

				// Navigate to the new conversation after short delay
				setTimeout(() => {
					this.$router.push(`/chat/${response.data.conversationId}`);
				}, 1);
			} catch (error) {
				this.errorMsg =
					error.response?.data?.message ||
					"Failed to create conversation";
				console.error(error);
			} finally {
				this.isCreating = false;
			}
		},
	},
	mounted() {
		this.loadAllUsers();
	},
};
</script>

<style scoped>
.user-entry {
	transition: background-color 0.2s ease;
	cursor: pointer;
}

.user-entry:hover {
	background-color: rgba(255, 255, 255, 0.05) !important;
}

.selected-users .badge {
	font-size: 1rem;
	padding: 0.5rem 1rem;
}
</style>
