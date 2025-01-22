<template>
	<!-- make the user entry take up the whole space of its container -->
	<div class="nav-item w-100">
		<button
			class="nav-link d-flex p-2 rounded border-0 w-100"
			@mouseover="isFocused = true"
			@mouseleave="isFocused = false"
			:class="{ 'bg-dark': !isFocused, 'bg-secondary': isFocused }"
			@click="$emit('open-conversation', this.conversationID)"
		>
			<!-- class="rounded-circle border border-3 border-primary" bouna idea per una notifica -->
			<img
				:src="imageUrl"
				class="rounded-circle"
				alt="avatar-icon"
				width="50"
				height="50"
				style="object-fit: cover"
			/>
			<!-- username  and last text -->
			<div
				class="ms-3 d-flex flex-column justify-content-evenly text-light overflow-auto w-100"
			>
				<div class="d-flex flex-row justify-content-between">
					<h5
						class="mb-0 fw-semibold align-self-start text-light text-start text-truncate"
					>
						{{ username }}
					</h5>
					<!-- timestamp -->
					<div class="text-secondary">
						{{ convertTimestamp(timestamp) }}
					</div>
				</div>
				<div class="d-flex w-100">
					<p
						class="mb-0 align-self-start text-light text-start text-truncate"
					>
						<!-- last text -->
						{{ lastText }}
					</p>
				</div>
			</div>
		</button>
	</div>
</template>

<script>
export default {
	props: {
		imageUrl: {
			type: String,
			required: true,
		},
		username: {
			type: String,
			required: true,
		},
		lastText: {
			type: String,
			required: true,
		},
		timestamp: {
			type: String,
			required: true,
		},
		conversationID: {
			type: String,
			required: true,
		},
	},
	data() {
		return {
			isFocused: false,
		};
	},

	methods: {
		// Emit the user entry event
		convertTimestamp(timestamp) {
			const date = new Date(timestamp);
			// if the date is NaN, return a default string.
			// (this is the case where you don't have any messages in the convo yet.)
			if (isNaN(date)) {
				return "start chatting!";
			}
			// if the date is today, return the time
			if (
				date.getDate() === new Date().getDate() &&
				date.getMonth() === new Date().getMonth() &&
				date.getFullYear() === new Date().getFullYear()
			) {
				return `${date.getHours()}:${date.getMinutes()}`;
			}

			// if the date is not today, return the date
			return `${date.getDate()}/${date.getMonth()}/${date.getFullYear()}`;
		},
	},
};
</script>
