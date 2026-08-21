<script lang="ts">
	import favicon from '$lib/assets/favicon.svg';
	import connectWebsocket from '$lib/ws/connection';
	import { onMount } from 'svelte';

	let { children } = $props();

	onMount(() => {
		connectWebsocket();
	});
</script>

<svelte:head>
	<link rel="icon" href={favicon} />
</svelte:head>

{@render children()}
