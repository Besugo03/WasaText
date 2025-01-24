<script>
export default {
	data() {
		return {
			profilePicture: null,
		};
	},
	methods: {
		async getProfilePicture() {
			// contact the api at /users/photo/$username.value
			// to get the user's profile pic.
			// if the request fails, set the profile pic to a default image.
			// if the request succeeds, set the profile pic to the returned image.

			// wait for the token to be set before making the request
			while (this.$token.value == "") {
				await new Promise((resolve) => setTimeout(resolve, 100));
			}
			try {
				let response = await this.$axios.get(
					`/users/photo/${this.$username.value}`,
					{
						headers: {
							Authorization: `Bearer ${this.$token.value}`,
						},
					}
				);
				this.profilePicture = response.data.photo;
			} catch (error) {
				this.profilePicture = "default.jpg";
			}
			console.log(this.profilePicture);
		},
	},
	async mounted() {
		await this.getProfilePicture();
	},
};
</script>

<template>
	<div
		class="p-4 d-flex flex-column align-items-center text-white"
		style="min-height: calc(100vh - 3.5rem)"
	>
		<h2 class="fw-semibold">Welcome, {{ $username.value }}!</h2>
		<img
			:src="profilePicture"
			alt="Profile Picture"
			class="rounded-circle border border-5 border-secondary"
			style="width: 10rem; height: 10rem; object-fit: cover"
		/>
		<h2>Begin chatting or change your settings from the navbar</h2>
	</div>
</template>

<style></style>
