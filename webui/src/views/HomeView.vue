<script>
export default {
	data: function () {
		return {
			errormsg: null,
			loading: false,
			some_data: null,
		};
	},
	methods: {
		async refresh() {
			this.loading = true;
			this.errormsg = null;
			// try {
			// 	let response = await this.$axios.post(
			// 		"/session",
			// 		{
			// 			// some params
			// 			username: "diuk",
			// 		},
			// 		{
			// 			headers: {
			// 				"Content-Type": "application/json",
			// 			},
			// 		}
			// 	);
			// 	let token = response.data.identifier;
			// 	console.log(token);

			// 	response = await this.$axios.get(`/conversations/info/${1}`, {
			// 		headers: {
			// 			Authorization: `Bearer ${token}`,
			// 		},
			// 	});

			// 	response = await this.$axios.get(`/conversations`, {
			// 		headers: {
			// 			Authorization: `Bearer ${token}`,
			// 		},
			// 	});

			// 	this.some_data = response.data;
			// } catch (e) {
			// 	this.errormsg = e.toString();
			// }
			this.loading = false;
		},
	},
	mounted() {
		this.refresh();
	},
};
</script>

<template>
	<div>
		<div
			class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom"
		>
			<div class="d-flex align-items-center">
				<img
					src="../assets/whatsapp2.png"
					class="me-4"
					alt="logo"
					style="height: 75px"
				/>
				<h1 class="h2">WASApp</h1>
			</div>
			<div class="btn-toolbar mb-2 mb-md-0">
				<div class="btn-group me-2">
					<button
						type="button"
						class="btn btn-sm btn-outline-secondary"
						@click="refresh"
					>
						Refresh
					</button>
					<button
						type="button"
						class="btn btn-sm btn-outline-secondary"
						@click="exportList"
					>
						Export
					</button>
				</div>
				<div class="btn-group me-2">
					<button
						type="button"
						class="btn btn-sm btn-outline-primary"
						@click="newItem"
					>
						New
					</button>
				</div>
			</div>
		</div>
		{{ some_data }}
		<ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
	</div>
</template>

<style></style>
