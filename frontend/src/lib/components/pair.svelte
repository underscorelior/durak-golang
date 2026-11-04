<script lang="ts">
	import { EMPTY_CARD, isCardEmpty } from '$lib/util';
	import Card from './card.svelte';
	const { pair }: { pair: CardPair } = $props();
</script>

<div class="pair-outside">
	<div class="bottom-card">
		<Card card={pair.attack_card} draggable={false} />
	</div>
	{#if !isCardEmpty(pair.defense_card ?? EMPTY_CARD)}
		<div class="top-card">
			<Card card={pair.defense_card ?? EMPTY_CARD} draggable={false} />
		</div>
	{/if}
</div>

<style>
	.pair-outside {
		position: relative;
	}

	.bottom-card {
		position: relative;
	}

	.top-card {
		position: absolute;
		top: 20%;
		left: 7%;
	}
</style>
