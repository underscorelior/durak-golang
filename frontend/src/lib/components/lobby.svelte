<script lang="ts">
	import { leaveLobby, startGame } from '$lib/event/actions';
	import { globalState } from '$lib/state/state.svelte';

	const lobby = $derived(globalState.lobby);
</script>

{#if lobby == null}
	<h1>You aren't in a lobby</h1>
{:else}
	<h1>{lobby.lobby_code} ({lobby.players.length}/{lobby.max_players})</h1>
	<h2>
		Host: {lobby.players.find((p) => {
			return p.user_id == lobby.host_id;
		})?.user_id}
	</h2>
	<p>{JSON.stringify(lobby.players)}</p>
	<br />
	<button onclick={() => leaveLobby()}>Leave Lobby</button>
	<button onclick={() => startGame()} hidden={lobby.host_id != globalState.user.user_id} disabled={lobby.players.length <= 1}>Start Game</button>
{/if}
