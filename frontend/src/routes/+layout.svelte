<script lang="ts">
	import favicon from '$lib/assets/favicon.svg';
	import { updateUser } from '$lib/event/actions';
	import { globalState } from '$lib/state/state.svelte';
	import connectWebsocket from '$lib/ws/connection';
	import { onMount } from 'svelte';

	let { children } = $props();

	let name = $state(globalState.user.name || '')
	let text = $state("")

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
	<div>
		<textarea bind:value={text}></textarea>
		<button onclick={()=>{text = JSON.stringify(globalState || {'undef': undefined})}}>update</button>
	</div>
</footer>