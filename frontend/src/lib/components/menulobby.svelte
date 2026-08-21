<script lang="ts">
	const { lobby }: { lobby: MenuLobby } = $props();
</script>

<p>
	{lobby.lobby_id.slice(0, 8)} ({lobby.player_count}/{lobby.max_players})
</p>
