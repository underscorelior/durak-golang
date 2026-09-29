<script lang="ts">
	import { leaveLobby, startGame } from '$lib/event/actions';
	import { globalState } from '$lib/state/state.svelte';
	import { shID } from '$lib/util';
	import Card from './card.svelte';
	import Table from './table.svelte';

	const lobby = $derived(globalState.lobby);
	let cards = $derived(globalState.lobby?.game_state?.hand);

</script>


{#if lobby == null}
	<h1>You aren't in a lobby</h1>
{:else}
	<h1>{lobby.lobby_code} ({lobby.players.length}/{lobby.max_players})</h1>
	<p style="margin-bottom:5px">
		Host: {shID(lobby.players.find((p) => {
			return p.user_id == lobby.host_id;
		})?.user_id)}
	</p>
	<button onclick={() => leaveLobby()} style="margin-bottom:10px;">Leave Lobby</button>
	<button onclick={() => startGame()} hidden={lobby.host_id != globalState.user.user_id || globalState.lobby?.is_started} disabled={lobby.players.length <= 1}>Start Game</button>

	<Table />

	{#if lobby.is_started}
		<div style="display: flex; flex-direction: row; gap: 5 5;">
			{#each cards as card, i (i)}
				<Card {card} />
			{/each}
		</div>
	{/if}
{/if}
