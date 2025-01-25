<template>
	<div
		class="convert-modal text-white p-4 rounded"
		style="width: 500px; background: #2d2d2d"
	>
		<div class="d-flex justify-content-between align-items-center mb-4">
			<h4 class="mb-0">Forward Message</h4>
			<button
				class="btn btn-sm btn-close btn-close-white"
				@click="$emit('close')"
			></button>
		</div>

		<div class="mb-3">
			<input
				type="text"
				class="form-control bg-dark text-white"
				placeholder="Search conversations..."
				v-model="searchQuery"
			/>
		</div>

		<div class="user-list overflow-auto" style="max-height: 300px">
			<div
				v-for="conversation in filteredConversations"
				:key="conversation.conversationID"
				class="user-entry d-flex align-items-center p-2 bg-dark rounded mb-2"
				@click="selectConversation(conversation.conversationID)"
			>
				<img
					:src="conversation.imageUrl"
					class="rounded-circle me-3"
					alt="Conversation"
					width="40"
					height="40"
					style="object-fit: cover"
				/>
				<div class="flex-grow-1">
					<h5 class="mb-0 text-light">{{ conversation.username }}</h5>
					<small class="text-muted">{{
						conversation.lastText
					}}</small>
				</div>
			</div>
		</div>

		<div v-if="errorMsg" class="alert alert-danger mt-3">
			{{ errorMsg }}
		</div>
		<div v-if="successMsg" class="alert alert-success mt-3">
			{{ successMsg }}
		</div>

		<div class="mt-4 d-flex justify-content-end gap-2">
			<button class="btn btn-secondary" @click="$emit('close')">
				Cancel
			</button>
		</div>
	</div>
</template>

<script>
export default {
	props: {
		conversations: {
			type: Array,
			required: true,
		},
		messageId: {
			type: String,
			required: true,
		},
	},
	data() {
		return {
			searchQuery: "",
			errorMsg: "",
			successMsg: "",
			isProcessing: false,
		};
	},
	computed: {
		filteredConversations() {
			return this.conversations.filter(
				(conv) =>
					conv.username
						.toLowerCase()
						.includes(this.searchQuery.toLowerCase()) &&
					conv.conversationID !== this.$route.params.id // Exclude current conversation
			);
		},
	},
	methods: {
		async selectConversation(targetConversationId) {
			this.isProcessing = true;
			this.errorMsg = "";
			this.successMsg = "";

			try {
				await this.$axios.post(
					"/conversations/message/forward",
					{
						conversationId: targetConversationId,
						messageId: this.messageId,
						isForwarded: true,
					},
					{
						headers: {
							Authorization: `Bearer ${this.$token.value}`,
						},
					}
				);

				this.successMsg = "Message forwarded successfully!";
				setTimeout(() => this.$emit("close"), 2000);
			} catch (error) {
				this.errorMsg =
					error.response?.data?.message ||
					"Failed to forward message";
				console.error("Forward error:", error);
			} finally {
				this.isProcessing = false;
			}
		},
	},
};
</script>

<style scoped>
.convert-modal {
	position: fixed;
	top: 50%;
	left: 50%;
	transform: translate(-50%, -50%);
	z-index: 1000;
}

.user-entry {
	transition: background-color 0.2s ease;
	cursor: pointer;
}

.user-entry:hover {
	background-color: rgba(255, 255, 255, 0.05);
}
</style>
