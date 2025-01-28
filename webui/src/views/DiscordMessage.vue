<template>
	<div
		class="message-container p-3 hover-bg w-100"
		:class="{
			'forwarded-message': isForwarded,
			'reply-active': isReplyVisualActive,
		}"
	>
		<div class="d-flex">
			<!-- User Avatar -->
			<div class="flex-shrink-0 me-3">
				<img
					:src="userPicture"
					alt="User avatar"
					class="rounded-circle"
					width="40"
					height="40"
					:style="{ objectFit: 'cover' }"
				/>
			</div>

			<!-- Message Content -->
			<div class="flex-grow-1 message-content">
				<!-- Message Header -->
				<div class="d-flex align-items-baseline mb-1">
					<span class="username fw-bold me-2 text-white"
						>{{ username }}
					</span>
					<span class="timestamp text-muted small">{{
						formattedTimestamp
					}}</span>
					<!-- Message actions -->
					<div class="message-actions ms-auto">
						<!-- Reply to message Button -->
						<button
							class="btn-action text-muted small me-2"
							title="Reply to message"
							@click="
								$emit('reply', messageID);
								toggleReplyVisual();
							"
						>
							<svg
								viewBox="0 0 24 24"
								fill="none"
								width="16"
								height="16"
								stroke="currentColor"
							>
								<path
									d="M20 17V15.8C20 14.1198 20 13.2798 19.673 12.638C19.3854 12.0735 18.9265 11.6146 18.362 11.327C17.7202 11 16.8802 11 15.2 11H4M4 11L8 7M4 11L8 15"
									stroke-width="2.3"
									stroke-linecap="round"
									stroke-linejoin="round"
								></path>
							</svg>
						</button>
						<!-- Forward button -->
						<button
							class="btn-action text-muted small"
							@click="$emit('forward', messageID)"
							title="Forward"
						>
							<svg
								width="16"
								height="16"
								viewBox="0 0 24 24"
								fill="none"
								stroke="currentColor"
							>
								<path
									d="M13.9546 5.18341L18.9324 9.60815C19.863 10.4353 20.3283 10.8489 20.4998 11.3373C20.6503 11.7662 20.6503 12.2335 20.4998 12.6624C20.3283 13.1508 19.863 13.5644 18.9324 14.3916L13.9546 18.8163C13.5323 19.1917 13.3211 19.3794 13.1418 19.3861C12.986 19.3919 12.8364 19.3247 12.7372 19.2044C12.6231 19.0659 12.6231 18.7834 12.6231 18.2184V15.4284C10.195 15.4284 7.63044 16.2083 5.75782 17.5926C4.78293 18.3133 4.29546 18.6737 4.1098 18.6595C3.92883 18.6456 3.81398 18.575 3.72008 18.4196C3.62374 18.2603 3.70883 17.7624 3.879 16.7666C4.98397 10.3004 9.43394 8.57129 12.6231 8.57129V5.78134C12.6231 5.21632 12.6231 4.93381 12.7372 4.79531C12.8364 4.67498 12.986 4.6078 13.1418 4.61363C13.3211 4.62034 13.5323 4.80803 13.9546 5.18341Z"
								/>
							</svg>
						</button>

						<!-- Delete button -->
						<button
							class="btn-action text-muted small ms-2"
							@click="$emit('delete', messageID)"
							title="Delete"
						>
							<svg
								width="16"
								height="16"
								viewBox="0 0 24 24"
								fill="none"
								stroke="currentColor"
							>
								<path
									d="M3 6h18M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2M10 11v6M14 11v6"
								/>
							</svg>
						</button>

						<!-- Add Comment Button -->
						<button
							class="btn-action text-muted small ms-2"
							title="Add Comment"
							@click="toggleCommentInput"
						>
							<svg
								viewBox="0 0 24 24"
								width="16"
								height="16"
								fill="none"
								stroke="currentColor"
							>
								<path
									d="M7 9H17M7 13H17M21 20L17.6757 18.3378C17.4237 18.2118 17.2977 18.1488 17.1656 18.1044C17.0484 18.065 16.9277 18.0365 16.8052 18.0193C16.6672 18 16.5263 18 16.2446 18H6.2C5.07989 18 4.51984 18 4.09202 17.782C3.71569 17.5903 3.40973 17.2843 3.21799 16.908C3 16.4802 3 15.9201 3 14.8V7.2C3 6.07989 3 5.51984 3.21799 5.09202C3.40973 4.71569 3.71569 4.40973 4.09202 4.21799C4.51984 4 5.0799 4 6.2 4H17.8C18.9201 4 19.4802 4 19.908 4.21799C20.2843 4.40973 20.5903 4.71569 20.782 5.09202C21 5.51984 21 6.0799 21 7.2V20Z"
									stroke-width="1.344"
									stroke-linecap="round"
									stroke-linejoin="round"
								></path>
							</svg>
						</button>
					</div>
				</div>

				<!-- Forwarded Indicator -->
				<div
					v-if="isForwarded"
					class="forwarded-indicator text-muted small mb-1"
				>
					<img
						src="https://static-00.iconduck.com/assets.00/forward-icon-2048x1450-aygcxglo.png"
						alt="Forwarded"
						width="16"
					/>
					Forwarded
				</div>
				<div
					v-if="replyingTo"
					class="text-white d-flex p-2 rounded mb-1"
					style="background-color: #282c34"
				>
					<div class="pe-1">Replying to</div>
					<div class="pe-1">{{ repliedMessage.sender }} :</div>
					<div class="text-muted">"{{ repliedMessage.content }}"</div>
				</div>
				<!-- Message Body -->
				<div class="message-body text-white">
					<div v-if="isSingleImage" class="single-image">
						<img :src="content" />
					</div>
					<template v-else>
						<template v-for="(part, index) in processedContent">
							<span
								v-if="part.type === 'text'"
								:key="`text-${index}`"
							>
								{{ part.content }}
							</span>
							<img
								v-else
								:key="`img-${index}`"
								:src="part.content"
								class="inline-image"
							/>
						</template>
					</template>
				</div>

				<!-- Add comment input -->
				<div v-if="showCommentInput" class="comment-input mt-2">
					<!-- <div class="comment-input mt-2"> -->
					<input
						ref="commentInput"
						v-model="newComment"
						@keyup.enter="submitComment"
						class="form-control form-control-sm bg-dark text-white"
						placeholder="Type a comment..."
					/>
				</div>
				<!-- Display comments -->
				<div
					v-if="comments.length > 0"
					class="comments mt-2 d-flex flex-wrap gap-2"
				>
					<!-- show the comment with primary bg if the current user was part of the commenters for it -->
					<div
						v-for="(comment, index) in aggregatedComments"
						:key="index"
						:style="{
							backgroundColor: comment.commenters.includes(
								$token.value
							)
								? '#6f562c'
								: '#282c34',
						}"
						class="comment-badge badge d-flex align-items-center gap-2"
						@click="addDuplicateComment(comment.text)"
					>
						<span>{{ comment.text }}</span>
						<span class="count-badge bg-dark rounded-pill px-2">{{
							comment.count
						}}</span>
					</div>
				</div>
				<!-- Read / Received display -->
				<div>
					<div
						v-if="!isReadByAll && isReceivedByAll"
						class="text-muted small"
					>
						<svg
							viewBox="0 0 24 24"
							width="16"
							height="16"
							fill="none"
						>
							<path
								id="Vector"
								d="M4 12L8.94975 16.9497L19.5572 6.34326"
								stroke="#f48e00"
								stroke-width="2"
								stroke-linecap="round"
								stroke-linejoin="round"
							></path>
						</svg>
					</div>
					<div v-if="isReceivedByAll" class="text-muted">
						<svg
							viewBox="0 0 24 24"
							width="16"
							height="16"
							fill="none"
						>
							<path
								id="Vector"
								d="M7 12L11.9497 16.9497L22.5572 6.34326M2.0498 12.0503L6.99955 17M17.606 6.39355L12.3027 11.6969"
								stroke="#f48e00"
								stroke-width="2"
								stroke-linecap="round"
								stroke-linejoin="round"
							></path>
						</svg>
					</div>
				</div>
			</div>
		</div>
	</div>
</template>

<script>
export default {
	data() {
		return {
			showCommentInput: false,
			newComment: "",
			comments: [],
			isLoadingComments: false,
			repliedMessage: {},
			isReplyVisualActive: false,
		};
	},
	mounted() {
		this.loadComments();
		if (this.replyingTo !== null) {
			this.getReplyingMessageInfo(this.replyingTo).then((data) => {
				this.repliedMessage = data;
			});
		}
	},
	methods: {
		toggleReplyVisual() {},
		// gets the content and the user that sent the message we are replying to
		async getReplyingMessageInfo(messageID) {
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
		toggleCommentInput() {
			this.showCommentInput = !this.showCommentInput;
			if (!this.showCommentInput) this.newComment = "";
			// make it so that the input is focused when it's shown
			if (this.showCommentInput) {
				this.$nextTick(() => {
					this.$refs.commentInput.focus();
				});
			}
		},

		async loadComments() {
			try {
				this.isLoadingComments = true;
				const response = await this.$axios.get(
					`/conversations/message/getcomments/${this.messageID}`,
					{
						headers: {
							Authorization: `Bearer ${this.$token.value}`,
						},
					}
				);
				this.comments = response.data;
			} catch (error) {
				console.error("Error loading comments:", error);
			} finally {
				this.isLoadingComments = false;
			}
		},

		async submitComment() {
			if (!this.newComment.trim()) return;

			try {
				await this.$axios.post(
					`/conversations/message/comment/${this.messageID}`,
					{ content: this.newComment },
					{
						headers: {
							Authorization: `Bearer ${this.$token.value}`,
						},
					}
				);
				await this.loadComments();
				this.newComment = "";
				this.showCommentInput = false;
			} catch (error) {
				console.error("Error submitting comment:", error);
			}
		},

		async removeComment(text) {
			let response = await this.$axios.get(
				`/conversations/message/getcomments/${this.messageID}`,
				{
					headers: {
						Authorization: `Bearer ${this.$token.value}`,
					},
				}
			);

			// since we can only have one same-text comment per user, we can safely assume that there is only one comment with the same text and senderId
			const ownComment = response.data.filter(
				(comment) =>
					comment.content === text &&
					comment.senderId === this.$token.value
			)[0];

			try {
				await this.$axios.delete(
					`/conversations/message/comment/${ownComment.commentId}`,
					{
						headers: {
							Authorization: `Bearer ${this.$token.value}`,
						},
					}
				);
				await this.loadComments();
			} catch (error) {
				console.error("Error removing comment:", error);
			}
		},

		async addDuplicateComment(text) {
			// first, we need to check if we already have a comment with the same content sent from the same user
			// contact the api to get the comments info
			// if the request fails, set the comments info to a default value.
			// if the request succeeds, set the comments info to the returned value.
			let response = await this.$axios.get(
				`/conversations/message/getcomments/${this.messageID}`,
				{
					headers: {
						Authorization: `Bearer ${this.$token.value}`,
					},
				}
			);

			const comments = response.data;
			if (
				comments.some(
					(comment) =>
						comment.content === text &&
						comment.senderId === this.$token.value
				)
			) {
				console.log("Comment already exists from you");
				this.removeComment(text);
				return;
			}

			try {
				await this.$axios.post(
					`/conversations/message/comment/${this.messageID}`,
					{ content: text },
					{
						headers: {
							Authorization: `Bearer ${this.$token.value}`,
						},
					}
				);
				await this.loadComments();
			} catch (error) {
				console.error("Error adding duplicate comment:", error);
			}
		},
	},

	// components: { ForwardSVG },
	props: {
		isRead: {
			type: Boolean,
			required: true,
		},
		isReadByAll: {
			type: Boolean,
			required: true,
		},
		isReceivedByAll: {
			type: Boolean,
			required: true,
		},
		isReceived: {
			type: Boolean,
			required: true,
		},
		replyingTo: {
			type: String,
			required: false,
		},
		messageID: {
			type: String,
			required: true,
		},
		username: {
			type: String,
			required: true,
		},
		userPicture: {
			type: String,
			default:
				"https://cdn.pixabay.com/photo/2015/10/05/22/37/blank-profile-picture-973460_960_720.png",
		},
		timestamp: {
			type: [String, Date],
			required: true,
		},
		isForwarded: {
			type: Boolean,
			default: false,
		},
		content: {
			type: String,
			required: true,
		},
	},
	emits: ["forward", "delete", "reply"],
	computed: {
		aggregatedComments() {
			const counts = {};
			const commenters = {};

			// Iterate through the comments to aggregate counts and collect commenters
			this.comments.forEach((comment) => {
				const content = comment.content;
				const commenter = comment.senderId; // Assuming `senderId` is the property for the commenter's ID

				// Initialize the count and commenters array if the content is encountered for the first time
				if (!counts[content]) {
					counts[content] = 0;
					commenters[content] = [];
				}

				// Increment the count for this content
				counts[content] += 1;

				// Add the commenter to the list if they haven't already been added
				if (!commenters[content].includes(commenter)) {
					commenters[content].push(commenter);
				}
			});

			// Transform the aggregated data into the desired format
			return Object.entries(counts).map(([text, count]) => ({
				text,
				count,
				commenters: commenters[text], // Include the list of commenters
			}));
		},
		formattedTimestamp() {
			if (this.timestamp instanceof Date) {
				return this.timestamp.toLocaleTimeString([], {
					hour: "2-digit",
					minute: "2-digit",
				});
			}
			return this.timestamp;
		},
		isSingleImage() {
			const imageRegex =
				/^https?:\/\/\S+?\.(?:png|jpg|jpeg|gif|webp)(?:\?\S*)?$/i;
			return imageRegex.test(this.content.trim());
		},
		processedContent() {
			const imageRegex =
				/(https?:\/\/\S+?\.(?:png|jpg|jpeg|gif|webp)(?:\?\S*)?)/gi;
			const parts = [];
			let lastIndex = 0;
			let match;

			while ((match = imageRegex.exec(this.content)) !== null) {
				const textBefore = this.content.slice(lastIndex, match.index);
				if (textBefore)
					parts.push({ type: "text", content: textBefore });
				parts.push({ type: "image", content: match[1] });
				lastIndex = match.index + match[0].length;
			}

			// Add the remaining text after the last image
			const textAfter = this.content.slice(lastIndex);
			if (textAfter) parts.push({ type: "text", content: textAfter });

			return parts;
		},
	},
};
</script>

<style>
/* .reply-active {
	background-color: #ffc42f23;
} */

.message-actions {
	opacity: 0;
	transition: opacity 0.2s ease;
}

.message-container {
}

.hover-bg:hover .message-actions {
	opacity: 1;
}

.btn-action {
	background: none;
	border: none;
	padding: 2px;
	cursor: pointer;
	transition: color 0.2s ease;
}

.btn-action:hover {
	color: #ffa200 !important;
	fill: #ffa200 !important;
}
.message-content {
	max-width: calc(100% - 56px);
}
.hover-bg {
	transition: background-color 0.2s ease;
	border-radius: 8px;
}

.message-body {
	color: #2e3338;
	word-break: break-word;
	overflow-wrap: break-word;
}

.single-image img {
	max-width: 100%;
	height: auto;
	border-radius: 4px;
	margin-top: 4px;
	display: block;
	cursor: pointer;
	max-height: 400px;
	object-fit: contain;
}

.single-image img:hover {
	box-shadow: 0 0 8px rgba(0, 0, 0, 0.1);
}

.inline-image {
	height: 1.5rem;
	vertical-align: middle;
	margin: 0 0.1em;
	display: inline-block;
	border-radius: 2px;
	object-fit: contain;
	/* max-height: 1.5em; */
}

.hover-bg:hover {
	background-color: rgba(255, 255, 255, 0.09) !important;
}

.forwarded-indicator {
	font-size: 0.75rem;
	line-height: 1;
}

.username {
	color: #060607;
}

.forwarded-message {
	position: relative;
}

.forwarded-message::before {
	content: "";
	position: absolute;
	left: 0;
	top: 0;
	bottom: 0;
	width: 2px;
	background-color: #f48e00;
	border-radius: 2px;
}
</style>

<style scoped>
.comment-input {
	width: 200px;
}

.comment-badge {
	cursor: pointer;
	transition: transform 0.2s ease;
}

.comment-badge:hover {
	transform: scale(1.05);
}

.count-badge {
	font-size: 0.75rem;
}
</style>
