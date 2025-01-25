<script>
import DiscordMessage from "./DiscordMessage.vue";
import MessageSend from "../components/MessageSend.vue";
import ForwardMessageModal from "../components/ForwardMessageModal.vue";
import { inject } from "vue";
export default {
	inject: ["refreshTrigger"],
	setup() {
		// refreshtrigger
		// const refreshTrigger = inject("refreshTrigger");
	},
	components: {
		DiscordMessage,
		MessageSend,
		ForwardMessageModal,
	},
	mounted() {
		const chatID = this.$route.params.id;
		console.log(`Chat ID: ${chatID}`);
		this.chatID = chatID;
		this.messages = [];
		this.MessageEntries = [];
		this.refreshChat();
		// listen to the new-messages event from app.vue
	},
	beforeUnmount() {},
	data() {
		return {
			chatID: this.$route.params.id,
			conversationInfo: null,
			messagesIDs: null,
			messages: [],
			MessageEntries: [],
			showForwardModal: false,
			ForwardMessageID: null,
			ForwardConversationID: null,
		};
	},
	watch: {
		refreshTrigger() {
			console.log("ChatView refreshTriggered");
			this.refreshChat();
		},
		$route(to, from) {
			this.chatID = to.params.id;
			console.log(`Chat ID: ${this.chatID}`);
			this.refreshChat();
		},
		MessageEntries(newVal) {
			console.log("Messages updated", newVal.length);
			// this.scrollToElement();

			// setTimeout(() => {
			// 	this.scrollToElement();
			// }, 100);
		},
	},
	methods: {
		async handleMessageSend(messageContent) {
			console.log("Message : ", messageContent);
			// contact the api to send the message

			// if the request fails, show an error message.
			// if the request succeeds, refresh the chat.

			this.$axios
				.post(
					"/conversations/message",
					{
						conversationId: this.chatID,
						content: messageContent,
						isForwarded: false,
					},
					{
						headers: {
							Authorization: `Bearer ${this.$token.value}`,
						},
					}
				)
				.then(() => {
					this.refreshChat();
				})
				.catch((error) => {
					console.error(error);
					this.$emit("error", "Failed to send message");
				});
		},
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
		async prepareMessageEntries() {
			this.MessageEntries = [];
			// for each message in this.messages, we need to get the props to generate the DiscordMessage component.
			// username is in the "sender" field of the message
			// user-picture needs to be retrieved by calling the api with the "sender" field of the message
			// timestamp is in the "timestamp" field of the message, and needs to be formatted
			// is-forwarded is in the "isForwarded" field of the message
			// content is in the "content" field of the message
			// once we have all the props, we can add the message to the MessageEntries array
			for (let i = 0; i < this.messages.length; i++) {
				let message = this.messages[i];
				let sender = message.sender;
				let userPicture = await this.$axios.get(
					`/users/photo/${sender}`,
					{
						headers: {
							Authorization: `Bearer ${this.$token.value}`,
						},
					}
				);
				userPicture = userPicture.data.photo;
				let timestamp = message.timestamp;
				// format the timestamp. if the day is today, show the time only. if the day is not today, show the date + " at " + time.
				let date = new Date(timestamp);
				let today = new Date();
				if (
					date.getDate() === today.getDate() &&
					date.getMonth() === today.getMonth() &&
					date.getFullYear() === today.getFullYear()
				) {
					timestamp = date.toLocaleTimeString([], {
						hour: "2-digit",
						minute: "2-digit",
					});
				} else {
					timestamp = date.toLocaleDateString([], {
						year: "numeric",
						month: "short",
						day: "numeric",
					});
					timestamp += " at ";
					timestamp += date.toLocaleTimeString([], {
						hour: "2-digit",
						minute: "2-digit",
					});
				}

				let messageID = message.messageId;
				console.log(message);
				let isForwarded = message.isForwarded;
				let content = message.content;
				this.MessageEntries.push({
					username: sender,
					userPicture: userPicture,
					timestamp: timestamp,
					isForwarded: isForwarded,
					content: content,
					messageID: messageID,
				});
				await this.$nextTick();
				this.scrollToElement();
			}
			// reverse the array so the messages are displayed in the correct order
			// this.MessageEntries.reverse();
		},
		async refreshChat() {
			this.messagesIDs = null;
			this.messages = [];
			this.MessageEntries = [];
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
			await this.prepareMessageEntries();
		},
		forwardMessage(messageID) {
			console.log(`Forwarding message ${messageID}`);
			this.showForwardModal = true;
			this.ForwardMessageID = messageID;
			// emit the forward event with the messageID
			// this.$emit("forward", messageID);
		},
		deleteMessage(messageID) {
			console.log(`Deleting message ${messageID}`);
			// contact the api to delete the message
			// if the request fails, show an error message.
			// if the request succeeds, refresh the chat.
			this.$axios
				.delete(`/conversations/message/delete/${messageID}`, {
					headers: {
						Authorization: `Bearer ${this.$token.value}`,
					},
				})
				.then(() => {
					this.refreshChat();
				})
				.catch((error) => {
					console.error(error);
					this.$emit("error", "Failed to delete message");
				});
		},
		closeForwardModal() {
			console.log("Closing forward modal");
			// this.$emit("close-forward-modal");
			this.showForwardModal = false;
			this.ForwardMessageID = null;
			this.refreshChat();
		},
	},
};
</script>

<template>
	<div>
		<div
			class="chat-wrapper d-flex flex-column"
			style="min-height: calc(100vh - 8.4rem)"
		>
			<ForwardMessageModal
				:conversations="$root.userEntries"
				:messageId="ForwardMessageID"
				:show="showForwardModal"
				@close="closeForwardModal"
			/>
			<!-- {{ conversationInfo }} -->
			<!-- {{ messages }} -->
			<!-- class="flex-grow-1 overflow-auto d-flex flex-column justify-content-start" -->
			<div
				class="flex-grow-1 overflow-auto d-flex flex-column justify-content-start"
				style="min-height: 0"
			>
				<div
					v-for="messageEntry in MessageEntries"
					:key="messageEntry.username + messageEntry.timestamp"
					class="d-flex justify-content-star"
				>
					<DiscordMessage
						:username="messageEntry.username"
						:user-picture="messageEntry.userPicture"
						:timestamp="messageEntry.timestamp"
						:is-forwarded="messageEntry.isForwarded"
						:content="messageEntry.content"
						:messageID="messageEntry.messageID"
						@forward="forwardMessage"
						@delete="deleteMessage"
					></DiscordMessage>
				</div>
				<!-- The div below is position at the end of the conversation.
				I found out that if the scroll method was called to the conversation wrapper,
				it ended up scrolling to the TOP of the conversation, which whas like not scrolling at all.
				Adding this fixed the problem. -->
				<div ref="scrollToMe" class="bg-dark"></div>
			</div>
		</div>
		<MessageSend @send="handleMessageSend" />
	</div>
</template>

<style>
.chat-wrapper {
	position: relative;
}

/* Add this to ensure messages don't get hidden behind the sticky input */
.message-send-container {
	background: white;
	padding: 1rem;
	box-shadow: 0 -2px 10px rgba(0, 0, 0, 0.05);
}
</style>
