<script lang="ts">
	import favicon from '$lib/assets/favicon.svg';
	import { updateUser } from '$lib/event/actions';
	import { globalState } from '$lib/state/state.svelte';
	import connectWebsocket from '$lib/ws/connection';
	import { onMount } from 'svelte';

	let { children } = $props();

	let name = $state(globalState.user.name || '')

	onMount(() => {
		connectWebsocket();
	});
</script>

<svelte:head>
	<link rel="icon" href={favicon} />
</svelte:head>

{@render children()}

<footer style="margin-top: 50px;">
	<input placeholder="name" bind:value={name}/>
	<button onclick={()=>updateUser(name)}>Update</button>
</footer>
