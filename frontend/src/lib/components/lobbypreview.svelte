<script lang="ts">
	import { joinLobby } from '$lib/event/actions';
	import { globalState } from '$lib/state/state.svelte';

	const { lobby }: { lobby: LobbyPreview } = $props();
</script>

<h1>
	{lobby.host_name}
</h1>
<p>
	{lobby.lobby_code.slice(0, 8)} ({lobby.player_count}/{lobby.max_players}) | {new Date(
		lobby.created_at
		// ).toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: 'numeric' })}
	).toLocaleString('en-US', { hour: 'numeric', minute: 'numeric' })}
</p>
<button
	onclick={() => joinLobby(lobby.lobby_code)}
	disabled={globalState.lobby != null || lobby.player_count >= lobby.max_players}
>
	Join Lobby
</button>
