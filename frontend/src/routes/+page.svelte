<script lang="ts">
	import LobbyPreview from '$lib/components/lobbypreview.svelte';
	import { createLobby } from '$lib/event/actions';
	import { connectionState, globalState } from '$lib/state/state.svelte';

	let lobbies = $derived(globalState.menu.lobbies);
</script>


<div class="lobby-list">
	{#if lobbies === null}
		<p>No Lobbies Found</p>
	{:else}
		{#each lobbies as lobby (lobby.lobby_code)}
			<LobbyPreview {lobby} />
		{/each}
	{/if}
</div>

<button onclick={() => createLobby()} disabled={!connectionState.connected}>
	Create a new lobby
</button>

{JSON.stringify(connectionState)}


<style>
	.lobby-list {
		max-height: 65vh;
		height: 65vh;
		width: 90%;
		margin-inline: auto;
		overflow-x: hidden;
		overflow-y: scroll;

		border: 1px solid black;
	}
</style>