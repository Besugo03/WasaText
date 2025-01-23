<template>
	<!-- Keep the same template structure -->
	<div
		class="message-send-container sticky-bottom bg-dark p-3 rounded-top w-100"
	>
		<div class="row g-2">
			<div class="col">
				<div class="d-flex align-items-center gap-2">
					<textarea
						ref="messageInput"
						v-model="messageText"
						class="form-control form-control-lg"
						placeholder="Type a message..."
						rows="1"
						@keydown.enter.exact.prevent="sendMessage"
						style="resize: vertical; max-height: 150px"
					></textarea>

					<button
						class="btn btn-primary btn-lg"
						:disabled="!messageText"
						@click="sendMessage"
					>
						Send
					</button>
				</div>
			</div>
			<div class="col-auto">
				<input
					type="file"
					ref="fileInput"
					accept="image/*"
					class="d-none"
					@change="handleImageUpload"
				/>
				<button
					class="btn btn-outline-secondary btn-lg h-100"
					:disabled="isUploading"
					@click="$refs.fileInput.click()"
				>
					<span v-if="!isUploading">📷 Upload Image</span>
					<span
						v-else
						class="spinner-border spinner-border-sm"
					></span>
				</button>
			</div>
		</div>
	</div>
</template>

<script>
export default {
	name: "MessageSend",
	data() {
		return {
			messageText: "",
			isUploading: false,
			imgbbApiKey: "21faba3731fc41dba67d6f0211b0cbca", // Replace with your ImgBB API key
		};
	},
	methods: {
		sendMessage() {
			if (this.messageText.trim()) {
				this.$emit("send", this.messageText);
				this.messageText = "";
			}
		},

		async handleImageUpload(event) {
			const file = event.target.files[0];
			if (!file) return;

			if (!file.type.startsWith("image/")) {
				alert("Please upload an image file");
				return;
			}

			this.isUploading = true;

			try {
				const imageUrl = await this.uploadToImgBB(file);
				this.insertImageLink(imageUrl);
			} catch (error) {
				console.error("Image upload failed:", error);
				alert("Image upload failed: " + error.message);
			} finally {
				this.isUploading = false;
				this.$refs.fileInput.value = "";
			}
		},

		async uploadToImgBB(file) {
			const formData = new FormData();
			formData.append("image", file);

			const response = await fetch(
				`https://api.imgbb.com/1/upload?key=${this.imgbbApiKey}`,
				{
					method: "POST",
					body: formData,
				}
			);

			const data = await response.json();

			if (!response.ok || !data.success) {
				throw new Error(
					data.error?.message || "Failed to upload image"
				);
			}

			return data.data.image.url;
		},

		insertImageLink(url) {
			this.messageText += ` ${url}`;
			this.$refs.messageInput.focus();
		},
	},
};
</script>

<style scoped>
/* Keep existing styles */
.message-send-container {
	position: sticky;
	bottom: 0;
	z-index: 1000;
}

textarea.form-control {
	min-height: 45px;
}
</style>
