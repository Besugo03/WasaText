<template>
	<div
		class="convert-modal text-white p-4 rounded"
		style="width: 500px; background: #2d2d2d"
		v-if="show"
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
				class="user-entry d-flex align-items-center p-2 rounded mb-2"
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
		<div v-if="uninitiatedUserConversations.length > 0">
			<h5 class="mb-0 mt-2">Or forward to a new chat...</h5>

			<div class="user-list overflow-auto" style="max-height: 300px">
				<div
					v-for="user in uninitiatedUserConversations"
					:key="user"
					class="user-entry d-flex align-items-center p-2 rounded mb-2"
					@click="forwardToNewChat(user)"
				>
					<img
						src="https://via.placeholder.com/40"
						class="rounded-circle me-3"
						alt="Conversation"
						width="40"
						height="40"
						style="object-fit: cover"
					/>
					<div class="flex-grow-1">
						<h5 class="mb-0 text-light">{{ user }}</h5>
						<small class="text-muted">New chat</small>
					</div>
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
	mounted() {
		console.log("Conversations: ");
		console.log(this.conversations.map((conv) => conv.username));
		this.getUninitiatedUserConversations();
	},
	props: {
		conversations: {
			type: Array,
			required: true,
		},
		messageId: {
			type: String,
			required: true,
		},
		show: {
			type: Boolean,
			required: true,
		},
		newConversations: {
			type: Array,
			required: false,
		},
	},
	data() {
		return {
			searchQuery: "",
			errorMsg: "",
			successMsg: "",
			isProcessing: false,
			uninitiatedUserConversations: [],
		};
	},
	computed: {
		filteredConversations() {
			return this.conversations.filter((conv) =>
				conv.username
					.toLowerCase()
					.includes(this.searchQuery.toLowerCase())
			);
		},
	},
	methods: {
		async forwardToNewChat(username) {
			this.isProcessing = true;
			this.errorMsg = "";
			this.successMsg = "";

			try {
				let response = await this.$axios.post(
					"/conversations/create",
					{
						participants: [username],
					},
					{
						headers: {
							Authorization: `Bearer ${this.$token.value}`,
						},
					}
				);
				let conversationId = response.data.conversationId;
				await this.selectConversation(conversationId);
			} catch (error) {
				this.errorMsg =
					error.response?.data?.message ||
					"Failed to create new chat";
				console.error("Forward error:", error);
			} finally {
				this.isProcessing = false;
			}
		},
		// this function gets all the User conversations (1-to-1) that have not been initiated yet.
		async getUninitiatedUserConversations() {
			// get all the users
			let response = await this.$axios.get("/users", {
				headers: {
					Authorization: `Bearer ${this.$token.value}`,
				},
			});
			let users = [];
			users = response.data;
			// get all the conversations for the current user
			response = await this.$axios.get("/conversations", {
				headers: {
					Authorization: `Bearer ${this.$token.value}`,
				},
			});
			let conversations = response.data.conversations;

			// remove ourselves from the users
			users = users.filter((user) => user !== this.$username.value);
			// filter out the conversations that we have already initiated
			console.log("Users: ", users);
			console.log("Conversations: ", conversations);

			// filter out all the conversations that have only a user in users as participants
			for (let i = 0; i < conversations.length; i++) {
				let conversation = conversations[i];
				let participants;
				// query the db for the participants of the conversation
				response = await this.$axios.get(
					`/conversations/info/${conversation}`,
					{
						headers: {
							Authorization: `Bearer ${this.$token.value}`,
						},
					}
				);
				participants = response.data.participants;
				console.log(
					`Participants for conversation ${conversation}: `,
					participants
				);
				// if it's a 1-on-1 conversation
				if (participants.length === 2) {
					console.log(`Conversation ${conversation} is a 1-on-1`);
					// remove ourselves from the participants
					participants = participants.filter(
						(participant) => participant !== this.$username.value
					);
					// if the remaining participant is in users, the conversation has been initiated and thus we dont care
					console.log("Participants: ", participants);
					if (users.includes(participants[0])) {
						// remove participants[0] from users
						users = users.filter(
							(user) => user !== participants[0]
						);
						continue;
					}
				}
			}

			// users now stores all the users that we have not initiated a conversation with
			console.log("Users: ", users);
			this.uninitiatedUserConversations = users;
		},
		async selectConversation(targetConversationId) {
			this.isProcessing = true;
			this.errorMsg = "";
			this.successMsg = "";

			console.log("conversationId: ", targetConversationId);
			console.log("messageId: ", this.messageId);
			console.log("token: ", this.$token.value);

			try {
				await this.$axios.post(
					"/conversations/message/forward",
					{
						conversationId: targetConversationId,
						messageId: this.messageId.toString(),
					},
					{
						headers: {
							Authorization: `Bearer ${this.$token.value}`,
						},
					}
				);

				this.successMsg = "Message forwarded successfully!";
				setTimeout(() => this.$emit("close"), 100);
				this.$router.push(`/chat/${targetConversationId}`);
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
	background-color: rgba(159, 159, 159, 0.1);
	transition: background-color 0.2s ease;
	cursor: pointer;
}

.user-entry:hover {
	background-color: rgba(139, 139, 139, 0.376);
}
</style>
