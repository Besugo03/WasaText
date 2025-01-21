<template></template>

<script>
export default {
	mounted() {
		const chatID = this.$route.params.id;
		console.log(`Chat ID: ${chatID}`);
	},
};
</script>

<style></style>
