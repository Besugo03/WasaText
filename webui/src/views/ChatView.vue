<script>
import DiscordMessage from "./DiscordMessage.vue";
export default {
	components: {
		DiscordMessage,
	},
	mounted() {
		const chatID = this.$route.params.id;
		console.log(`Chat ID: ${chatID}`);
		this.chatID = chatID;
		this.messages = [];
		this.refreshChat();
	},
	data() {
		return {
			chatID: this.$route.params.id,
			conversationInfo: null,
			messagesIDs: null,
			messages: [],
		};
	},
	watch: {
		$route(to, from) {
			this.chatID = to.params.id;
			console.log(`Chat ID: ${this.chatID}`);
			this.refreshChat();
		},
		messages(newVal) {
			console.log("Messages updated", newVal.length);
			this.scrollToElement();
		},
	},
	methods: {
		scrollToElement() {
			const el = this.$refs.scrollToMe;

			if (el) {
				// Use el.scrollIntoView() to instantly scroll to the element
				el.scrollIntoView({ behavior: "smooth" });
			}
		},
		async getMessagesDetails(messageID) {
			// contact the api to get the message data
			// if the request fails, set the message data to a default value.
			// if the request succeeds, set the message data to the returned value.
			let response = await this.$axios.get(
				`/conversations/message/info/${messageID}`,
				{
					headers: {
						Authorization: `Bearer ${this.$token.value}`,
					},
				}
			);
			return response.data;
		},
		async refreshChat() {
			// contact the api to get the chat data
			// if the request fails, set the chat data to a default value.
			// if the request succeeds, set the chat data to the returned value.
			let response = await this.$axios.get(
				`/conversations/info/${this.chatID}`,
				{
					headers: {
						Authorization: `Bearer ${this.$token.value}`,
					},
				}
			);

			this.conversationInfo = response;
			// console.log(this.conversationInfo);
			this.messagesIDs = this.conversationInfo.data.messages;
			if (this.messagesIDs === null) {
				this.messages = [];
				return;
			}

			let retrievedMessages = [];
			// for each messageID, call the getMessageDetails to get the message data
			// add each one to retreivedMessages, and once done, set this.messages to retrievedMessages
			for (let i = 0; i < this.messagesIDs.length; i++) {
				let messageData = await this.getMessagesDetails(
					this.messagesIDs[i]
				);
				retrievedMessages.push(messageData);
			}
			this.messages = retrievedMessages;
		},
	},
};
</script>

<template>
	<div class="chat-wrapper d-flex flex-column">
		{{ conversationInfo }}
		<DiscordMessage
			username="John Doe"
			user-picture="https://gratisography.com/wp-content/uploads/2024/11/gratisography-augmented-reality-800x525.jpg"
			timestamp="9:42 AM"
			:is-forwarded="false"
		>
			Hello, this is a forwarded message.
			<img
				src="https://gratisography.com/wp-content/uploads/2024/11/gratisography-augmented-reality-800x525.jpg"
			/>
		</DiscordMessage>
		<div
			class="flex-grow-1 overflow-auto d-flex flex-column justify-content-end"
			style="min-height: 0"
		>
			<div
				v-for="message in messages"
				:key="message.id"
				class="d-flex justify-content-start"
			>
				<div class="bg-secondary bg-opacity-50 p-2 rounded m-2">
					{{ message.content }}
				</div>
			</div>
		</div>
		<div ref="scrollToMe"></div>
	</div>
</template>

<style></style>
