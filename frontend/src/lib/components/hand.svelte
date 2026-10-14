<script lang="ts">
	import { globalState } from '$lib/state/state.svelte';
	import { fanDistance, findGamePlayerByID, findLobbyPlayerByID, shID } from '$lib/util';
	import Card from './card.svelte';

	const {
		opponent,
		user_id,
		orientation
	}: { opponent: boolean; user_id: string; orientation: number } = $props();

	let game_state = $derived(globalState.lobby?.game_state);
	let hand = $derived(game_state?.hand);

	let game_started = $derived(globalState.lobby?.is_started);

	let hand_size = $derived(findGamePlayerByID(user_id)?.hand_size ?? 0);
</script>

<div style="--orientation: {orientation};" class="container">
	<div class="nameplate">
		{shID(user_id)} ({findLobbyPlayerByID(user_id)?.position})
		<br />
		{#if game_started}{hand_size}{/if}
	</div>
	{#if game_started}
		<div class="hand">
			{#if opponent}
				{#each Array(hand_size).keys() as index (index)}
					<Card
						card={{ suit: 0, rank: 0 }}
						{index}
						distance={fanDistance(index, hand_size)}
						inHand
					/>
				{/each}
			{:else}
				{#each hand as card, index (index)}
					<Card {card} {index} distance={fanDistance(index, hand_size)} inHand />
				{/each}
			{/if}
		</div>
	{/if}
</div>

<style>
	.container {
		rotate: calc(var(--orientation) * 90deg);
		margin-inline: auto;
		width: 100%;
		height: 100%;
		justify-content: center;
		align-items: center;
	}

	.nameplate {
		position: absolute;
		background-color: rgba(25, 25, 25, 0.5);
		color: aliceblue;
		padding: 4px;
		border-radius: 4px;
		rotate: calc(var(--orientation) * -90deg);
	}

	.hand {
		display: flex;
		flex-direction: row;
	}
</style>
