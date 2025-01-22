<template>
	<div
		class="message-container p-3 hover-bg"
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
					<span class="username fw-bold me-2">{{ username }}</span>
					<span class="timestamp text-muted small">{{
						formattedTimestamp
					}}</span>
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
				<div class="message-body">
					<slot></slot>
				</div>
			</div>
		</div>
	</div>
</template>

<script>
export default {
	props: {
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
	},
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
	},
};
</script>

<style>
.message-content {
	max-width: calc(100% - 56px); /* 40px avatar + 16px margin */
}
.hover-bg {
	transition: background-color 0.2s ease;
	border-radius: 8px;
}

.message-body {
	color: #2e3338;
	word-break: break-word;
	overflow-wrap: break-word;
	max-width: 100%;
}

.message-body img {
	max-width: 100%;
	height: auto;
	border-radius: 4px;
	margin-top: 4px;
	display: block;
	cursor: pointer;
	max-height: 400px; /* Adjust this value as needed */
	object-fit: contain;
}

.message-body img:hover {
	box-shadow: 0 0 8px rgba(0, 0, 0, 0.1);
}

.message-body {
	overflow: hidden;
}

.hover-bg:hover {
	background-color: rgba(4, 4, 5, 0.07);
}

.forwarded-indicator {
	font-size: 0.75rem;
	line-height: 1;
}

.username {
	color: #060607;
}

.message-body {
	color: #2e3338;
	word-break: break-word;
}

.forwarded-message {
	/* opacity: 0.8; */
	position: relative;
}

.forwarded-message::before {
	content: "";
	position: absolute;
	left: 0;
	top: 0;
	bottom: 0;
	width: 2px;
	background-color: #00aff4;
	border-radius: 2px;
}
</style>
