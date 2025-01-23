<script>
export default {
	data: function () {
		return {
			errormsg: null,
			loading: false,
			some_data: null,
			username: "",
		};
	},
	methods: {
		async refresh() {
			this.loading = true;
			this.errormsg = null;
			this.loading = false;
		},
		login() {
			this.$username = this.username; // Set the global variable to the input value
			this.$token.value = ""; // Clear the token
			this.$emit("login", this.username); // Emit the login event
			console.log(`Logged in as: ${this.$username}`);
		},
	},
	mounted() {
		this.refresh();
	},
};
</script>

<template>
	<div style="min-height: calc(100vh - 3.5rem)">
		<div
			class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom text-white"
		>
			<h1 class="h2">Welcome! Please log in</h1>
		</div>
		{{ some_data }}
		<ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
		<!-- textbox for entering the username. -->
		<div class="input-group mb-3">
			<input
				type="text"
				class="form-control"
				placeholder="Username"
				aria-label="Username"
				v-model="username"
				@keyup.enter="login"
			/>
			<!-- on click, it sets this.$username to the contents of the textbox  -->
			<button
				class="btn btn-outline-secondary btn-primary text-light"
				type="button"
				id="button-addon2"
				@click="login"
			>
				Log in
			</button>
		</div>
	</div>
</template>

<style></style>
