<script>
import { RouterLink, RouterView } from "vue-router";
import { ref, provide } from "vue";
import UserEntry from "./components/UserEntry.vue";
import ConvertToGroupModal from "./components/ConvertToGroupModal.vue";
export default {
	components: { UserEntry, ConvertToGroupModal },

	setup() {
		const refreshTrigger = ref(0);
		provide("refreshTrigger", refreshTrigger);
		return { refreshTrigger };
	},

	// in the data, we will save our conversations.
	data() {
		return {
			// token is our authentication token
			errormsg: null,
			loading: false,
			conversationIDs: null,
			// this is an array of objects. Each object has the following properties:
			// their key is the conversation ID
			// their contents are the results of querying the database for the conversation details.
			conversationDetails: null,
			userEntries: [],
			token: null,
			pollingInterval: 5000,
			pollingTimer: null,

			// Add participants modal
			showConvertToGroupModal: false,
			newParticipants: [],
			currentParticipants: [],
			isGroupChat: false,
			isConversation: false,
		};
	},

	watch: {
		"$route.params.id": {
			immediate: true,
			async handler(newId) {
				if (newId) {
					const details = await this.getConversationDetails(newId);
					console.log("details: " + details);
					this.isGroupChat = details?.isGroup || false;
					this.currentParticipants = details?.participants || [];
				}
			},
		},
	},

	methods: {
		async handleParticipantsUpdate() {
			// Refresh the current conversation details
			const details = await this.getConversationDetails(
				this.$route.params.id
			);
			this.currentParticipants = details?.participants || [];
			this.isGroupChat = details?.isGroup || false;

			// Force refresh the header section
			this.refreshTrigger++;
		},
		async removeParticipant(username) {
			try {
				await this.$axios.put(
					"/conversations/kickuser",
					{
						conversationId: this.$route.params.id,
						username: username,
					},
					{
						headers: {
							Authorization: `Bearer ${this.$token.value}`,
						},
					}
				);
				await this.handleParticipantsUpdate(); // Refresh after removal
			} catch (error) {
				console.error("Failed to remove participant:", error);
			}
		},

		async leaveGroup() {
			await this.removeParticipant(this.$username.value);
			this.$router.push("/");
		},
		navigateToSettings() {
			const isGroup = this.isGroupChat;
			const route = isGroup
				? `/group/settings/${this.$route.params.id}`
				: `/profile/settings`;
			this.$router.push(route);
		},

		async convertToGroup(participants) {
			try {
				// Placeholder API call - replace with actual endpoint
				await this.$axios.put(
					`/conversations/convert-to-group/${this.$route.params.id}`,
					{ participants },
					{
						headers: {
							Authorization: `Bearer ${this.$token.value}`,
						},
					}
				);

				this.refreshTrigger++;
				this.showConvertToGroupModal = false;
			} catch (error) {
				console.error("Conversion failed:", error);
			}
		},
		startPolling() {
			console.log("Starting polling");
			this.pollingTimer = setInterval(async () => {
				if (!this.loading) {
					// Prevent overlapping requests
					await this.refreshConversations();
					console.log("Polling for new messages");
				}
			}, this.pollingInterval);
		},
		navigateToChat(conversationID) {
			this.$router.push(`/chat/${conversationID}`);
		},

		async getToken() {
			console.log("stored token : ", localStorage.getItem("token"));
			console.log("stored username : ", localStorage.getItem("username"));
			if (localStorage.getItem("token")) {
				console.log("Token already exists");
				this.$token.value = localStorage.getItem("token");
				this.$username.value = localStorage.getItem("username");
				return this.$token.value;
			}
			// if (localStorage.getItem("token")) {
			// 	console.log("Token already exists");
			// 	this.$token.value = localStorage.getItem("token");
			// 	this.$username.value = localStorage.getItem("username");
			// 	return this.$token.value;
			// }
			this.loading = true;
			this.errormsg = null;
			if (this.$username.value === "") {
				console.log("No username set");
				this.loading = false;
				return;
			}
			if (this.$token.value) {
				// console.log("Token already exists: " + this.$token.value);
				return;
			}
			try {
				let response = await this.$axios.post(
					"/session",
					{
						// some params
						username: this.$username.value,
					},
					{
						headers: {
							"Content-Type": "application/json",
						},
					}
				);
				this.token = response.data.identifier;
				this.$token.value = this.token;
				localStorage.setItem("token", this.token);
				localStorage.setItem("username", this.$username.value);
			} catch (e) {
				this.errormsg = e.toString();
			}
			this.loading = false;
			return this.token;
		},

		async getConversationIDs() {
			this.loading = true;
			this.errormsg = null;
			try {
				let response = await this.$axios.get(`/conversations`, {
					headers: {
						Authorization: `Bearer ${this.$token.value}`,
					},
				});
				this.conversationIDs = response.data.conversations;
				// it's passed as an array, so convert it from the Proxy { <target>: (1) […], <handler>: {…} }
				// to a plain array
				this.conversationIDs = this.conversationIDs;
			} catch (e) {
				this.errormsg = e.toString();
			}
			this.loading = false;
			return this.conversationIDs;
		},

		async getConversationDetails(conversationID) {
			this.loading = true;
			this.errormsg = null;
			try {
				let response = await this.$axios.get(
					`/conversations/info/${conversationID}`,
					{
						headers: {
							Authorization: `Bearer ${this.$token.value}`,
						},
					}
				);
				this.conversationDetails = response.data;
			} catch (e) {
				this.errormsg = e.toString();
			}
			this.loading = false;
			return this.conversationDetails;
		},

		// this gets called for each conversation, and gives back a freshly-baked user-entry object, just how the component wants it.
		async initializeUserEntry(conversationDetails) {
			// our user-entry expect 3 props: imageUrl, username, lastText
			// The conditions are as follows:
			// we get is_group from the conversationDetails object
			// for the image :
			// if the conversation is a group, we use the group image
			// if the conversation is not a group, we should have either 1 or 2 users.
			// if there is only 1 user, we use a default image (a conversation with yourself)
			// if there are 2 users, we use the image of the other user.

			// console.log(
			// 	"Conversation details: " + JSON.stringify(conversationDetails)
			// );

			let conversationImage;
			if (conversationDetails.isGroup) {
				conversationImage = conversationDetails.picture;
			} else {
				if (conversationDetails.participants.length == 1) {
					conversationImage =
						"https://i.ytimg.com/vi/ow-i4-Uj9NA/maxresdefault.jpg";
				} else {
					// this is the case in which we have 2 users.
					// we need to find the user that is not the current user.
					let chatParticipants = conversationDetails.participants;
					// remove yourself from the list of participants
					let otherUser = chatParticipants.filter(
						(user) => user != this.$username.value
					)[0];

					// probably the case where i got an "undefined" user was just an artifact
					// created from the backend and not a bad frontend request.
					// TODO note that i got that edgecase by removing people from conversatoions.
					if (otherUser == undefined || otherUser == "") {
						otherUser = this.$username.value;
					}

					let response = await this.$axios.get(
						`/users/photo/${otherUser}`,
						{
							headers: {
								Authorization: `Bearer ${this.$token.value}`,
							},
						}
					);

					conversationImage = response.data.photo;
				}
			}
			if (conversationImage == "") {
				conversationImage =
					"https://external-content.duckduckgo.com/iu/?u=https%3A%2F%2Fstatic.vecteezy.com%2Fsystem%2Fresources%2Fpreviews%2F004%2F141%2F669%2Foriginal%2Fno-photo-or-blank-image-icon-loading-images-or-missing-image-mark-image-not-available-or-image-coming-soon-sign-simple-nature-silhouette-in-frame-isolated-illustration-vector.jpg&f=1&nofb=1&ipt=93ff810ba334bbdf63f2253d5d5cda42742059f5549d3f6c058d10ad55a79084&ipo=images";
			}

			// for the username:
			// if the conversation is a group, we use the group name
			// if the conversation is not a group, we should have either 1 or 2 users.
			// if there is only 1 user, we are stil in the "chat with yourself" scenario, and we use a default name.
			// if there are 2 users, we use the name of the other user.

			let conversationUsername;
			if (conversationDetails.isGroup) {
				conversationUsername = conversationDetails.name;
			} else {
				if (conversationDetails.participants.length == 1) {
					conversationUsername = "Chat with yourself. Crazy!";
				} else {
					// this is the case in which we have 2 users.
					// we need to find the user that is not the current user.
					let chatParticipants = conversationDetails.participants;
					// remove yourself from the list of participants
					let otherUser = chatParticipants.filter(
						(user) => user != this.$username.value
					)[0];

					conversationUsername = otherUser;
				}
			}
			if (conversationUsername == "") {
				conversationUsername = "No name conversation";
			}

			// for the last text:
			// it's really the same. We already get the last text and its details all ready from the conversationDetails object.
			// only thing is, if the conversation is a group, we should display who sent the last message, like "Diuk : MESSAGE"
			// so we swap the last message with this modified version that has the name of the sender before it
			// P.S. the name is whatever is before the first space from the user's name.

			let conversationLastText;
			if (conversationDetails.isGroup) {
				let lastMessageContent =
					conversationDetails.lastMessage.content;
				let sender =
					conversationDetails.lastMessage.sender.split(" ")[0];
				conversationLastText = `${sender} : ${lastMessageContent}`;
			} else {
				conversationLastText = conversationDetails.lastMessage.content;
			}

			// get the time of the last message
			let lastMessageTime = conversationDetails.lastMessage.timestamp;

			// converastionID which the component needs to redirect to the correct chat
			let conversationID = conversationDetails.conversationId;

			return {
				imageUrl: conversationImage,
				username: conversationUsername,
				lastText: conversationLastText,
				timestamp: lastMessageTime,
				conversationID: conversationID,
			};
		},

		async refreshConversations() {
			await this.getToken();
			if (!this.$token.value) {
				console.log("No token set");
				return;
			}
			await this.getConversationIDs();
			// console.log(JSON.stringify(this.conversationIDs));

			// for each conversation ID, get the conversation details. saving it all in a temporary object.
			let retrievedEntries = [];
			let retrievedConversationDetails = {};

			if (
				this.conversationIDs == null ||
				this.conversationIDs.length === 0
			) {
				console.log("No conversations yet");
				return;
			}
			for (let conversationID of this.conversationIDs) {
				let details = await this.getConversationDetails(conversationID);
				retrievedConversationDetails[conversationID] = details;
			}

			this.conversationDetails = retrievedConversationDetails;
			// for each conversation, initialize a user-entry object.
			for (let conversationID in this.conversationDetails) {
				let conversationDetails =
					this.conversationDetails[conversationID];
				let userEntry = await this.initializeUserEntry(
					conversationDetails
				);
				retrievedEntries.push(userEntry);
			}

			// sort the user entries by timestamp, by most recent. (the format is a string yyyy-mm-dd hh:mm:ss)
			retrievedEntries.sort(
				(a, b) => new Date(b.timestamp) - new Date(a.timestamp)
			);

			// console.log(
			// 	"Retrieved entries: " + JSON.stringify(retrievedEntries)
			// );
			// console.log("User entries: " + JSON.stringify(this.userEntries));
			// if the user entries are different from the retrieved entries, update the user entries.

			if (
				JSON.stringify(this.userEntries) !==
				JSON.stringify(retrievedEntries)
			) {
				console.log("New messages!");
				this.userEntries = [];
				this.userEntries = retrievedEntries;
				this.refreshTrigger++;
				console.log("refreshTrigger value : " + this.refreshTrigger);
			} else {
				console.log("No new messages...");
			}
			return this.userEntries;
		},

		async handleLogin(username) {
			this.$username.value = username;
			this.$token.value = "";
			localStorage.setItem("username", this.$username.value);
			localStorage.setItem("token", "");
			this.getToken();
			// clear the previous user entries
			this.userEntries = [];
			this.$router.push("/welcomescreen");
			console.log("Navigating to welcome screen");
			await this.refreshConversations().then(() => {});
		},
	},

	// mounted is a lifecycle hook that runs after the component has been added to the DOM
	// we will use it to get the token, but we need to make the function async to use await.
	async mounted() {
		// await this.refreshConversations();
		await this.getToken();
		this.startPolling();
	},

	beforeUnmount() {
		if (this.pollingTimer) {
			clearInterval(this.pollingTimer);
			this.pollingTimer = null;
		}
	},
};
</script>

<template>
	<div>
		<header
			class="navbar navbar-dark sticky-top bg-dark flex-md-nowrap p-0 shadow"
			style="height: 3.5rem"
		>
			<div
				class="d-flex flex-row align-items-center justify-content-between w-100"
			>
				<a class="navbar-brand me-0 px-3 fs-6" href="#/">
					WASApp
					<img
						src="https://external-content.duckduckgo.com/iu/?u=https%3A%2F%2Fcdn-icons-png.flaticon.com%2F512%2F8708%2F8708467.png&f=1&nofb=1&ipt=37149496808e91dee11682f8051a4a4c27dd377f3033d141398e7f5a86805d27&ipo=images"
						width="30"
						class=""
					/>
				</a>
				<!-- if the current path (link) ends with "chat/NUMBER", display the div -->
				<!-- In App.vue's template section, modify the header conditional -->
				<div
					v-if="$route.path.split('/')[1] === 'chat'"
					class="d-flex flex-row align-items-center position-relative"
					id="chatHeader"
				>
					<img
						:src="
							userEntries.find(
								(entry) =>
									entry.conversationID === $route.params.id
							).imageUrl
						"
						class="rounded-circle"
						alt="avatar-icon"
						width="30"
						height="30"
						style="object-fit: cover"
					/>
					<h5 class="text-light fw-semibold ms-3 mb-0">
						{{
							userEntries.find(
								(entry) =>
									entry.conversationID === $route.params.id
							).username
						}}
					</h5>

					<!-- Add participants hover and settings -->
					<div class="header-actions ms-2">
						<button
							class="btn btn-link text-light"
							@click="navigateToSettings"
						>
							<svg
								width="20"
								height="20"
								viewBox="0 0 24 24"
								fill="none"
							>
								<path
									d="M19.43 12.98c.04-.3.04-.6.04-.98s0-.68-.04-.98l2.11-1.65c.19-.15.24-.42.12-.64l-2-3.46c-.12-.22-.39-.3-.61-.22l-2.49 1c-.52-.4-1.08-.73-1.69-.98l-.38-2.65c-.03-.24-.24-.42-.49-.42h-4c-.25 0-.46.18-.49.42l-.38 2.65c-.61.25-1.17.59-1.69.98l-2.49-1c-.22-.08-.49 0-.61.22l-2 3.46c-.13.22-.07.49.12.64l2.11 1.65c-.04.3-.04.6-.04.98s0 .68.04.98l-2.11 1.65c-.19.15-.24.42-.12.64l2 3.46c.12.22.39.3.61.22l2.49-1c.52.4 1.08.73 1.69.98l.38 2.65c.03.24.24.42.49.42h4c.25 0 .46-.18.49-.42l.38-2.65c.61-.25 1.17-.59 1.69-.98l2.49 1c.22.08.49 0 .61-.22l2-3.46c.12-.22.07-.49-.12-.64l-2.11-1.65zm-7.43 2.52c-1.93 0-3.5-1.57-3.5-3.5s1.57-3.5 3.5-3.5 3.5 1.57 3.5 3.5-1.57 3.5-3.5 3.5z"
									fill="currentColor"
								/>
							</svg>
						</button>

						<!-- Participants hover card -->
						<div
							class="participants-popover bg-dark p-3 rounded shadow"
							style="width: 250px"
						>
							<div
								class="d-flex justify-content-between align-items-center mb-3"
							>
								<h6 class="text-white mb-0">Participants</h6>
								<button
									class="btn btn-sm btn-primary"
									@click.stop="showConvertToGroupModal = true"
								>
									Add Members
								</button>
							</div>

							<div class="mb-3">
								<div
									v-for="participant in currentParticipants"
									:key="participant"
									class="d-flex justify-content-between align-items-center mb-2"
								>
									<span class="text-light">{{
										participant
									}}</span>
									<button
										v-if="participant !== $username.value"
										class="btn btn-sm btn-danger"
										@click.stop="
											removeParticipant(participant);
											this.refreshTrigger++;
										"
									>
										Remove
									</button>
								</div>
							</div>

							<button
								class="btn btn-warning btn-sm w-100"
								@click.stop="leaveGroup"
							>
								Leave Group
							</button>
						</div>
						<!-- Convert to group button -->
						<ConvertToGroupModal
							v-if="showConvertToGroupModal"
							:current-participants="currentParticipants"
							@close="showConvertToGroupModal = false"
							@added-participants="handleParticipantsUpdate"
						/>
					</div>
				</div>
				<!-- empty div to center the current chat -->
				<div></div>
				<button
					class="navbar-toggler position-absolute d-md-none collapsed"
					type="button"
					data-bs-toggle="collapse"
					data-bs-target="#sidebarMenu"
					aria-controls="sidebarMenu"
					aria-expanded="false"
					aria-label="Toggle navigation"
				>
					<span class="navbar-toggler-icon"></span>
				</button>
			</div>
		</header>
		<div class="container-fluid">
			<div class="row">
				<nav
					id="sidebarMenu"
					class="d-md-block bg-dark sidebar collapse col-md-3 col-lg-2"
					style="border-top-right-radius: 1rem"
				>
					<div class="position-sticky sidebar-sticky">
						<h6
							class="sidebar-heading d-flex justify-content-between align-items-center px-3 mt-4 mb-1 text-muted text-uppercase"
						>
							<span>General</span>
						</h6>
						<ul class="nav flex-column">
							<li class="nav-item">
								<RouterLink to="/" class="nav-link text-light">
									<svg class="feather">
										<use
											href="/feather-sprite-v4.29.0.svg#key"
										/>
									</svg>
									Change user
								</RouterLink>
							</li>
							<li class="nav-item">
								<RouterLink
									to="/settings"
									class="nav-link text-light"
								>
									<svg
										viewBox="0 0 24 24"
										fill="none"
										width="18"
										xmlns="http://www.w3.org/2000/svg"
										class="feather mb-0"
									>
										<path
											d="M14.3365 12.3466L14.0765 11.9195C13.9082 12.022 13.8158 12.2137 13.8405 12.4092C13.8651 12.6046 14.0022 12.7674 14.1907 12.8249L14.3365 12.3466ZM9.6634 12.3466L9.80923 12.8249C9.99769 12.7674 10.1348 12.6046 10.1595 12.4092C10.1841 12.2137 10.0917 12.022 9.92339 11.9195L9.6634 12.3466ZM4.06161 19.002L3.56544 18.9402L4.06161 19.002ZM19.9383 19.002L20.4345 18.9402L19.9383 19.002ZM16 8.5C16 9.94799 15.2309 11.2168 14.0765 11.9195L14.5965 12.7737C16.0365 11.8971 17 10.3113 17 8.5H16ZM12 4.5C14.2091 4.5 16 6.29086 16 8.5H17C17 5.73858 14.7614 3.5 12 3.5V4.5ZM7.99996 8.5C7.99996 6.29086 9.79082 4.5 12 4.5V3.5C9.23854 3.5 6.99996 5.73858 6.99996 8.5H7.99996ZM9.92339 11.9195C8.76904 11.2168 7.99996 9.948 7.99996 8.5H6.99996C6.99996 10.3113 7.96342 11.8971 9.40342 12.7737L9.92339 11.9195ZM9.51758 11.8683C6.36083 12.8309 3.98356 15.5804 3.56544 18.9402L4.55778 19.0637C4.92638 16.1018 7.02381 13.6742 9.80923 12.8249L9.51758 11.8683ZM3.56544 18.9402C3.45493 19.8282 4.19055 20.5 4.99996 20.5V19.5C4.70481 19.5 4.53188 19.2719 4.55778 19.0637L3.56544 18.9402ZM4.99996 20.5H19V19.5H4.99996V20.5ZM19 20.5C19.8094 20.5 20.545 19.8282 20.4345 18.9402L19.4421 19.0637C19.468 19.2719 19.2951 19.5 19 19.5V20.5ZM20.4345 18.9402C20.0164 15.5804 17.6391 12.8309 14.4823 11.8683L14.1907 12.8249C16.9761 13.6742 19.0735 16.1018 19.4421 19.0637L20.4345 18.9402Z"
											fill="#000000"
										></path>
									</svg>
									Change your profile
								</RouterLink>
							</li>
							<li class="nav-item">
								<RouterLink
									to="/link2"
									class="nav-link text-light"
								>
									<svg
										viewBox="0 0 24 24"
										width="18"
										fill="none"
										class="feather mb-1"
									>
										<path
											d="M15.7955 15.8111L21 21M18 10.5C18 14.6421 14.6421 18 10.5 18C6.35786 18 3 14.6421 3 10.5C3 6.35786 6.35786 3 10.5 3C14.6421 3 18 6.35786 18 10.5Z"
											stroke-width="2"
											stroke-linecap="round"
											stroke-linejoin="round"
										></path>
									</svg>
									Find a user
								</RouterLink>
							</li>
							<li class="nav-item">
								<RouterLink
									to="/create-conversation"
									class="nav-link text-light"
								>
									<svg
										viewBox="0 0 24 24"
										width="18"
										fill="none"
										class="feather mb-1"
									>
										<path
											d="M12 4v16m8-8H4"
											stroke="currentColor"
											stroke-width="2"
											stroke-linecap="round"
										/>
									</svg>
									Create Group Chat
								</RouterLink>
							</li>
						</ul>

						<h6
							class="sidebar-heading d-flex justify-content-between align-items-center px-3 mt-4 mb-1 text-muted text-uppercase"
						>
							<span>Your chats</span>
						</h6>

						<!-- create an user entry for each one we prepared earlier -->
						<user-entry
							v-for="(conversation, index) in userEntries"
							:key="index"
							:imageUrl="conversation.imageUrl"
							:username="conversation.username"
							:lastText="conversation.lastText"
							:timestamp="conversation.timestamp"
							:conversationID="conversation.conversationID"
							@open-conversation="navigateToChat"
						></user-entry>
						<div
							v-if="
								(!userEntries || userEntries.length === 0) &&
								$username.value === ''
							"
						>
							<h5 class="text-light fw-bold p-3 bg-secondary">
								No chats yet. Please log in.
							</h5>
						</div>
						<div
							v-if="
								(!userEntries || userEntries.length === 0) &&
								$username.value !== ''
							"
						>
							<h5 class="text-light fw-semibold p-3 bg-secondary">
								start a chat, {{ $username.value }}!
							</h5>
						</div>

						<ul class="nav flex-column">
							<li class="nav-item">
								<!-- <RouterLink
									:to="'/some/' + 'variable_here' + '/path'"
									class="nav-link text-light"
								>
									<svg class="feather">
										<use
											href="/feather-sprite-v4.29.0.svg#file-text"
										/>
									</svg>
									Item 1
								</RouterLink> -->
							</li>
						</ul>
					</div>
				</nav>

				<!-- here goes the main screen -->
				<!-- se metto vh-100 poi mi appaiono due scrollbar su una 
				  chat piena ma una chat vuota e tutte le altre schermate funzionano... -->
				<!-- se non lo metto si vedono tutti bordi bianchi nelle chat quasi vuote o nelle altre schermate -->
				<main
					class="col-md-9 col-lg-10 ms-sm-auto h-100"
					style="background-color: #16191c"
				>
					<RouterView :key="$route.fullPath" @login="handleLogin" />
				</main>
			</div>
		</div>
	</div>
</template>

<style>
#chatHeader .header-actions {
	position: relative;
	margin-left: 0.5rem;
}

.participants-popover {
	position: absolute;
	top: 100%;
	left: 0;
	width: 200px;
	z-index: 1000;
	display: none;
}

#chatHeader:hover .participants-popover {
	display: block;
}

.convert-modal {
	position: fixed;
	top: 50%;
	left: 50%;
	transform: translate(-50%, -50%);
	background: #2d2d2d;
	padding: 20px;
	border-radius: 8px;
	z-index: 1001;
}
</style>
