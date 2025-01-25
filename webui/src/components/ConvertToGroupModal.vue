<template>
	<div
		class="convert-modal text-white p-4 rounded"
		style="width: 500px; background: #2d2d2d"
	>
		<div class="d-flex justify-content-between align-items-center mb-4">
			<h4 class="mb-0">Add Participants</h4>
			<button
				class="btn btn-sm btn-close btn-close-white"
				@click="$emit('close')"
			></button>
		</div>

		<div class="mb-3">
			<input
				type="text"
				class="form-control bg-dark text-white"
				placeholder="Search users..."
				v-model="searchQuery"
				@keyup.enter="searchUsers"
			/>
		</div>

		<div class="selected-users mb-3">
			<div class="d-flex flex-wrap gap-2">
				<div
					v-for="user in selectedParticipants"
					:key="user.username"
					class="badge bg-primary p-2 d-flex align-items-center gap-2"
				>
					{{ user.username }}
					<button
						class="btn btn-sm btn-danger p-0"
						@click="toggleSelection(user)"
					>
						&times;
					</button>
				</div>
			</div>
		</div>

		<div v-if="errorMsg" class="alert alert-danger">{{ errorMsg }}</div>

		<div class="user-list overflow-auto" style="max-height: 300px">
			<div
				v-for="user in filteredUsers"
				:key="user.username"
				class="user-entry d-flex align-items-center p-2 bg-dark rounded mb-2"
			>
				<img
					:src="user.photo"
					class="rounded-circle me-3"
					alt="Profile picture"
					width="40"
					height="40"
					style="object-fit: cover"
				/>
				<div class="flex-grow-1">
					<h5 class="mb-0 text-light">{{ user.username }}</h5>
				</div>
				<button
					class="btn btn-sm"
					:class="isSelected(user) ? 'btn-danger' : 'btn-primary'"
					@click="toggleSelection(user)"
				>
					{{ isSelected(user) ? "Remove" : "Add" }}
				</button>
			</div>
		</div>

		<div class="mt-4 d-flex justify-content-end gap-2">
			<button class="btn btn-secondary" @click="$emit('close')">
				Cancel
			</button>
			<button
				class="btn btn-primary"
				:disabled="selectedParticipants.length === 0 || isSaving"
				@click="addParticipants"
			>
				<span
					v-if="isSaving"
					class="spinner-border spinner-border-sm"
				></span>
				{{
					isSaving
						? "Adding..."
						: `Add ${selectedParticipants.length} Participants`
				}}
			</button>
		</div>
	</div>
</template>

<script>
export default {
	props: {
		currentParticipants: {
			type: Array,
			default: () => [], // Add default empty array
		},
	},
	data() {
		return {
			searchQuery: "",
			allUsers: [],
			selectedParticipants: [],
			isSaving: false,
			errorMsg: "",
			defaultProfileImage:
				"https://cdn.pixabay.com/photo/2015/10/05/22/37/blank-profile-picture-973460_1280.png",
		};
	},
	computed: {
		filteredUsers() {
			return this.allUsers.filter(
				(user) =>
					user.username
						.toLowerCase()
						.includes(this.searchQuery.toLowerCase()) &&
					!this.currentParticipants.includes(user.username) && // Add 'this.'
					user.username !== this.$username.value
			);
		},
	},
	methods: {
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

		async loadAllUsers() {
			try {
				const response = await this.$axios.get("/users", {
					headers: { Authorization: `Bearer ${this.$token.value}` },
				});

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
				console.error("Failed to load users:", error);
			}
		},

		async addParticipants() {
			this.isSaving = true;
			this.errorMsg = "";

			try {
				const participants = this.selectedParticipants.map(
					(u) => u.username
				);
				for (let i = 0; i < participants.length; i++) {
					await this.$axios.put(
						"/conversations/adduser",
						{
							conversationId: this.$route.params.id,
							username: participants[i],
						},
						{
							headers: {
								Authorization: `Bearer ${this.$token.value}`,
							},
						}
					);
				}

				this.$emit("added-participants", participants);
				this.$emit("close");
			} catch (error) {
				this.errorMsg =
					error.response?.data?.message ||
					"Failed to add participants";
				console.error("Add participants error:", error);
			} finally {
				this.isSaving = false;
			}
		},
	},
	async mounted() {
		await this.loadAllUsers();
	},
};
</script>

<style scoped>
.user-entry {
	transition: background-color 0.2s ease;
	cursor: pointer;
}

.user-entry:hover {
	background-color: rgba(255, 255, 255, 0.05);
}

.selected-users .badge {
	font-size: 0.9rem;
	padding: 0.4rem 0.8rem;
}
</style>
