<template>
	<!-- Keep the same template structure -->
	<div
		class="message-send-container sticky-bottom bg-dark p-3 rounded-top w-100"
	>
		<div class="row g-2 d-flex flex-row">
			<div class="w-auto h-auto flex-col text-white" v-if="messageStart">
				<div>
					<svg
						data-v-6128a7d3=""
						viewBox="0 0 24 24"
						fill="none"
						width="16"
						height="16"
						stroke="#ffffff"
					>
						<path
							data-v-6128a7d3=""
							d="M20 17V15.8C20 14.1198 20 13.2798 19.673 12.638C19.3854 12.0735 18.9265 11.6146 18.362 11.327C17.7202 11 16.8802 11 15.2 11H4M4 11L8 7M4 11L8 15"
							stroke-width="2.3"
							stroke-linecap="round"
							stroke-linejoin="round"
						></path>
					</svg>
				</div>
				{{ replyingTo }}
			</div>
			<div class="col">
				<div class="d-flex align-items-center gap-2">
					<textarea
						ref="messageInput"
						v-model="messageText"
						class="form-control form-control-lg form-control-dark bg-dark text-white"
						placeholder="Type a message..."
						rows="1"
						@keydown.enter.exact.prevent="sendMessage"
						style="resize: vertical; max-height: 150px"
					></textarea>
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
							<span v-if="!isUploading">
								<svg
									viewBox="-2.4 -2.4 28.80 28.80"
									width="24"
									height="24"
								>
									<path
										opacity="0.5"
										fill-rule="evenodd"
										clip-rule="evenodd"
										d="M3 14.25C3.41421 14.25 3.75 14.5858 3.75 15C3.75 16.4354 3.75159 17.4365 3.85315 18.1919C3.9518 18.9257 4.13225 19.3142 4.40901 19.591C4.68577 19.8678 5.07435 20.0482 5.80812 20.1469C6.56347 20.2484 7.56459 20.25 9 20.25H15C16.4354 20.25 17.4365 20.2484 18.1919 20.1469C18.9257 20.0482 19.3142 19.8678 19.591 19.591C19.8678 19.3142 20.0482 18.9257 20.1469 18.1919C20.2484 17.4365 20.25 16.4354 20.25 15C20.25 14.5858 20.5858 14.25 21 14.25C21.4142 14.25 21.75 14.5858 21.75 15V15.0549C21.75 16.4225 21.75 17.5248 21.6335 18.3918C21.5125 19.2919 21.2536 20.0497 20.6517 20.6516C20.0497 21.2536 19.2919 21.5125 18.3918 21.6335C17.5248 21.75 16.4225 21.75 15.0549 21.75H8.94513C7.57754 21.75 6.47522 21.75 5.60825 21.6335C4.70814 21.5125 3.95027 21.2536 3.34835 20.6517C2.74643 20.0497 2.48754 19.2919 2.36652 18.3918C2.24996 17.5248 2.24998 16.4225 2.25 15.0549C2.25 15.0366 2.25 15.0183 2.25 15C2.25 14.5858 2.58579 14.25 3 14.25Z"
										fill="#ffffff"
									></path>
									<path
										fill-rule="evenodd"
										clip-rule="evenodd"
										d="M12 2.25C12.2106 2.25 12.4114 2.33852 12.5535 2.49392L16.5535 6.86892C16.833 7.17462 16.8118 7.64902 16.5061 7.92852C16.2004 8.20802 15.726 8.18678 15.4465 7.88108L12.75 4.9318V16C12.75 16.4142 12.4142 16.75 12 16.75C11.5858 16.75 11.25 16.4142 11.25 16V4.9318L8.55353 7.88108C8.27403 8.18678 7.79963 8.20802 7.49393 7.92852C7.18823 7.64902 7.16698 7.17462 7.44648 6.86892L11.4465 2.49392C11.5886 2.33852 11.7894 2.25 12 2.25Z"
										fill="#ffffff"
									></path>
								</svg>
							</span>
							<span
								v-else
								class="spinner-border spinner-border-sm"
							></span>
						</button>
					</div>

					<button
						class="btn btn-primary btn-lg"
						:disabled="!messageText"
						@click="sendMessage"
					>
						<!-- Send -->
						<svg
							viewBox="-2.4 -1 28.80 28.80"
							width="24"
							fill="none"
						>
							<path
								d="M11.5003 12H5.41872M5.24634 12.7972L4.24158 15.7986C3.69128 17.4424 3.41613 18.2643 3.61359 18.7704C3.78506 19.21 4.15335 19.5432 4.6078 19.6701C5.13111 19.8161 5.92151 19.4604 7.50231 18.7491L17.6367 14.1886C19.1797 13.4942 19.9512 13.1471 20.1896 12.6648C20.3968 12.2458 20.3968 11.7541 20.1896 11.3351C19.9512 10.8529 19.1797 10.5057 17.6367 9.81135L7.48483 5.24303C5.90879 4.53382 5.12078 4.17921 4.59799 4.32468C4.14397 4.45101 3.77572 4.78336 3.60365 5.22209C3.40551 5.72728 3.67772 6.54741 4.22215 8.18767L5.24829 11.2793C5.34179 11.561 5.38855 11.7019 5.407 11.8459C5.42338 11.9738 5.42321 12.1032 5.40651 12.231C5.38768 12.375 5.34057 12.5157 5.24634 12.7972Z"
								stroke="#ffffff"
								stroke-width="2"
								stroke-linecap="round"
								stroke-linejoin="round"
							></path>
						</svg>
					</button>
				</div>
			</div>
		</div>
	</div>
</template>

<script>
export default {
	props: {
		// string to start the message with
		messageStart: {
			type: String,
			default: "",
			required: false,
		},
	},
	name: "MessageSend",
	// mounted() {
	// 	this.$root.$on("reply", this.handleReply);
	// },
	// beforeUnmount() {
	// 	this.$root.$off("reply", this.handleReply);
	// },
	data() {
		return {
			messageText: "",
			isUploading: false,
			imgbbApiKey: "IMGUBB_API_KEY_PLACEHOLDER", // Replace with your ImgBB API key
			messageStartString: "",
			replyingTo: null,
		};
	},
	watch: {
		messageStart(newVal) {
			this.messageStartString = newVal;
			// set the replyingto to the messageStartString
			this.replyingTo = newVal.split(" ")[2];
		},
	},
	methods: {
		handleReply(message) {
			console.log("Replying to message:", message);
			this.replyingTo = message;
			this.messageText = `> ${message.username} said: ${message.content}\n`;
			this.$refs.messageInput.focus();
		},
		sendMessage() {
			if (this.messageText.trim()) {
				this.$emit("send", this.messageStartString + this.messageText);
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
	height: 80px;
	bottom: 0;
	z-index: 1000;
}

textarea.form-control {
	min-height: 45px;
}
</style>
