<script setup>
import { RouterLink, RouterView } from "vue-router";
import UserEntry from "./components/UserEntry.vue";
</script>
<script>
export default {
	components: { UserEntry },

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
			pollingInterval: 3000,
			pollingTimer: null,
		};
	},

	methods: {
		navigateToChat(conversationID) {
			this.$router.push(`/chat/${conversationID}`);
		},

		async getToken() {
			this.loading = true;
			this.errormsg = null;
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
						Authorization: `Bearer ${this.token}`,
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
							Authorization: `Bearer ${this.token}`,
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

					// query the database for the user's image
					let response = await this.$axios.get(
						`/users/photo/${otherUser}`,
						{
							headers: {
								Authorization: `Bearer ${this.token}`,
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
			console.log("logged in. Token: " + this.token);
			await this.getConversationIDs();
			// console.log(JSON.stringify(this.conversationIDs));

			// for each conversation ID, get the conversation details. saving it all in a temporary object.
			let retrievedConversationDetails = {};
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
				this.userEntries.push(userEntry);
			}

			// sort the user entries by timestamp, by most recent. (the format is a string yyyy-mm-dd hh:mm:ss)
			this.userEntries.sort(
				(a, b) => new Date(b.timestamp) - new Date(a.timestamp)
			);
		},

		async handleLogin(username) {
			this.$username.value = username;
			// clear the previous user entries
			this.userEntries = [];
			await this.refreshConversations();
			this.$router.push("/welcomescreen");
		},
	},

	// mounted is a lifecycle hook that runs after the component has been added to the DOM
	// we will use it to get the token, but we need to make the function async to use await.
	async mounted() {
		await this.refreshConversations();
	},

	beforeUnmount() {
		clearInterval(this.pollingTimer);
	},
};
</script>

<template>
	<div>
		<header
			class="navbar navbar-dark sticky-top bg-dark flex-md-nowrap p-0 shadow"
		>
			<a class="navbar-brand me-0 px-3 fs-6" href="#/">
				WASApp
				<img
					src="https://external-content.duckduckgo.com/iu/?u=https%3A%2F%2Fcdn-icons-png.flaticon.com%2F512%2F8708%2F8708467.png&f=1&nofb=1&ipt=37149496808e91dee11682f8051a4a4c27dd377f3033d141398e7f5a86805d27&ipo=images"
					width="30"
					class=""
				/>
			</a>
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
		</header>

		<div class="container-fluid">
			<div class="row">
				<nav
					id="sidebarMenu"
					class="d-md-block bg-dark sidebar collapse col-md-3 col-lg-2"
					style="border-top-right-radius: 1rem"
				>
					<div class="position-sticky pt-3 sidebar-sticky">
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
									to="/link1"
									class="nav-link text-light"
								>
									<svg class="feather">
										<use
											href="/feather-sprite-v4.29.0.svg#layout"
										/>
									</svg>
									Menu item 1
								</RouterLink>
							</li>
							<li class="nav-item">
								<RouterLink
									to="/link2"
									class="nav-link text-light"
								>
									<svg class="feather">
										<use
											href="/feather-sprite-v4.29.0.svg#key"
										/>
									</svg>
									Menu item 2
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
				<main
					class="col-md-9 col-lg-10 ms-sm-auto h-100"
					style="background-color: #16191c"
				>
					<RouterView @login="handleLogin" />
				</main>
			</div>
		</div>
	</div>
</template>

<style></style>
