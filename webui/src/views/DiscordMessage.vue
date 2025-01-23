<template>
	<div
		class="message-container p-3 hover-bg w-100"
		:class="{ 'forwarded-message': isForwarded }"
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
					<div class="message-actions ms-auto">
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
								>{{ part.content }}</span
							>
							<img
								v-else
								:key="`img-${index}`"
								:src="part.content"
								class="inline-image"
							/>
						</template>
					</template>
				</div>
			</div>
		</div>
	</div>
</template>

<script>
export default {
	// components: { ForwardSVG },
	props: {
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
	emits: ["forward", "delete"],
	computed: {
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

			const textAfter = this.content.slice(lastIndex);
			if (textAfter) parts.push({ type: "text", content: textAfter });

			return parts;
		},
	},
};
</script>

<style>
.message-actions {
	opacity: 0;
	transition: opacity 0.2s ease;
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
