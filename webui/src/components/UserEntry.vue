<template>
	<!-- make the user entry take up the whole space of its container -->
	<div class="w-100 nav-item">
		<button
			class="nav-link w-100 d-flex p-2 rounded border-0"
			@mouseover="isFocused = true"
			@mouseleave="isFocused = false"
			:class="{ 'bg-dark': !isFocused, 'bg-secondary': isFocused }"
		>
			<img
				:src="imageUrl"
				class="rounded-circle"
				alt="avatar-icon"
				width="50"
				height="50"
			/>
			<!-- username  and last text -->
			<div class="ms-3 d-flex flex-column justify-content-evenly">
				<h5 class="mb-0 fw-semibold text-light align-self-start">
					{{ username }}
				</h5>
				<p
					class="mb-0 align-self-start text-truncate"
					style="max-width: 200px"
				>
					<!-- last text -->
					<span class="text-light">{{ lastText }}</span>
				</p>
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
	},
	data() {
		return {
			isFocused: false,
		};
	},
};
</script>
