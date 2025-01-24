[file name]: ProfileSettings.vue [file content]
<template>
	<div
		class="p-4 d-flex flex-column text-white"
		style="min-height: calc(100vh - 3.5rem)"
	>
		<h2 class="fw-semibold mb-4">
			{{ isGroup ? "Group Settings" : "Profile Settings" }}
		</h2>

		<!-- Current Profile Display -->
		<div class="d-flex flex-column align-items-center mb-5">
			<img
				:src="currentPhoto || defaultImage"
				class="rounded-circle border border-5 border-secondary mb-3"
				style="width: 10rem; height: 10rem; object-fit: cover"
			/>
			<h4 class="text-center">
				{{ currentName }}
			</h4>
		</div>

		<!-- Settings Form -->
		<div class="row justify-content-center">
			<div class="col-md-8 col-lg-6">
				<!-- Image Upload -->
				<div class="mb-4">
					<label class="form-label"
						>Change
						{{ isGroup ? "Group Image" : "Profile Picture" }}</label
					>
					<div class="d-flex gap-3 align-items-center">
						<input
							type="file"
							ref="fileInput"
							accept="image/*"
							class="form-control bg-dark text-white"
							@change="handleImageUpload"
						/>
						<button
							class="btn btn-outline-danger"
							@click="clearImage"
							:disabled="!newPhoto"
						>
							Reset
						</button>
					</div>
				</div>

				<!-- Name Change -->
				<div class="mb-4">
					<label class="form-label"
						>Change
						{{ isGroup ? "Group Name" : "Display Name" }}</label
					>
					<input
						type="text"
						class="form-control bg-dark text-white"
						v-model="editedName"
						:placeholder="currentName"
					/>
				</div>
				<!-- Action Buttons -->
				<div class="d-flex gap-2 justify-content-end">
					<button class="btn btn-secondary" @click="resetFields">
						Cancel
					</button>
					<button
						class="btn btn-primary"
						@click="saveChanges"
						:disabled="!hasChanges || isLoading"
					>
						<span
							v-if="isLoading"
							class="spinner-border spinner-border-sm"
						></span>
						{{ isLoading ? "Saving..." : "Save Changes" }}
					</button>
				</div>

				<!-- Error Message -->
				<div v-if="errorMsg" class="alert alert-danger mt-3">
					{{ errorMsg }}
				</div>
			</div>
		</div>
	</div>
</template>

<script>
export default {
	data() {
		return {
			imgbbApiKey: "21faba3731fc41dba67d6f0211b0cbca", // Add imgBB API key
			isGroup: false,
			currentName: "",
			currentPhoto: "",
			editedName: "",
			newPhoto: null,
			isLoading: false,
			errorMsg: "",
			defaultImage:
				"https://cdn.pixabay.com/photo/2015/10/05/22/37/blank-profile-picture-973460_1280.png",
		};
	},
	computed: {
		hasChanges() {
			return (
				this.editedName !== this.currentName || this.newPhoto !== null
			);
		},
	},
	async mounted() {
		this.isGroup = this.$route.path.includes("/group/");
		await this.loadCurrentSettings();
	},
	methods: {
		async loadCurrentSettings() {
			try {
				if (this.isGroup) {
					const response = await this.$axios.get(
						`/conversations/info/${this.$route.params.id}`,
						{
							headers: {
								Authorization: `Bearer ${this.$token.value}`,
							},
						}
					);
					this.currentName = response.data.name;
					this.currentPhoto = response.data.picture;
				} else {
					console.log("not group");
					console.log(this.$username.value);
					const response = await this.$axios.get(
						"/users/photo/" + this.$username.value,
						{
							headers: {
								Authorization: `Bearer ${this.$token.value}`,
							},
						}
					);
					// this.currentName = response.data.username;
					this.currentName = this.$username.value;
					this.currentPhoto = response.data.photo;
				}
				this.editedName = this.currentName;
			} catch (error) {
				this.errorMsg = "Failed to load current settings";
				console.error(error);
			}
		},

		handleImageUpload(event) {
			const file = event.target.files[0];
			if (!file) return;

			this.newPhoto = file;
			const reader = new FileReader();
			reader.onload = (e) => {
				this.currentPhoto = e.target.result;
			};
			reader.readAsDataURL(file);
		},

		clearImage() {
			this.newPhoto = null;
			this.currentPhoto = this.isGroup
				? this.currentPhoto // Keep group image if exists
				: this.defaultImage; // Reset to default for profile
			this.$refs.fileInput.value = "";
		},

		async saveChanges() {
			this.isLoading = true;
			this.errorMsg = "";

			try {
				let imageUrl = this.currentPhoto;

				// Upload new photo to imgBB if changed
				if (this.newPhoto) {
					const imgbbFormData = new FormData();
					imgbbFormData.append("image", this.newPhoto);

					// Upload to imgBB
					const imgbbResponse = await fetch(
						`https://api.imgbb.com/1/upload?key=${this.imgbbApiKey}`,
						{
							method: "POST",
							body: imgbbFormData,
						}
					);

					const imgbbData = await imgbbResponse.json();
					if (!imgbbData.success) {
						throw new Error("Image upload failed");
					}

					imageUrl = imgbbData.data.url;

					// Update photo URL in backend
					const endpoint = this.isGroup
						? `/conversations/photo`
						: `/user/photo`;

					await this.$axios.put(
						endpoint,
						{
							photo: imageUrl,
							conversationId: this.$route.params.id,
						},
						{
							headers: {
								Authorization: `Bearer ${this.$token.value}`,
								"Content-Type": "application/json",
							},
						}
					);
				}

				// Update name if changed
				if (this.editedName !== this.currentName) {
					if (this.isGroup) {
						await this.$axios.put(
							`/conversations/name`,
							{
								name: this.editedName,
								conversationId: this.$route.params.id,
							},
							{
								headers: {
									Authorization: `Bearer ${this.$token.value}`,
									"Content-Type": "application/json",
								},
							}
						);
					} else {
						await this.$axios.put(
							"/user/username",
							{ username: this.editedName },
							{
								headers: {
									Authorization: `Bearer ${this.$token.value}`,
								},
							}
						);
					}

					// Update currentName if successful
					this.currentName = this.editedName;
					this.$username.value = this.editedName;
				}

				// Refresh data
				await this.loadCurrentSettings();
				this.newPhoto = null;
			} catch (error) {
				this.errorMsg =
					error.response?.data ||
					error.message ||
					"Failed to save changes";
				console.error(error);
			} finally {
				this.isLoading = false;
			}

			this.$router.push(
				this.isGroup ? `/chat/${this.$route.params.id}` : "/profile"
			);
		},

		resetFields() {
			this.editedName = this.currentName;
			this.clearImage();
		},
	},
};
</script>

<style scoped>
.form-control {
	border-color: #495057;
}

.form-control:focus {
	background-color: #212529;
	border-color: #6c757d;
	box-shadow: none;
}
</style>
