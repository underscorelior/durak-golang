<script lang="ts">
	import { isCardEmpty, rankToString, suitToColor, suitToSymbol } from '$lib/util';

	const {
		card,
		index = 0,
		distance = 0,
		draggable = true,
		inHand = false
	}: {
		card: Card;
		index?: number;
		distance?: number;
		draggable?: boolean;
		inHand?: boolean;
	} = $props();

	const rank = $derived(rankToString[card.rank]);
	const suitSymbol = $derived(suitToSymbol[card.suit]);
	const suitColor = $derived(suitToColor[card.suit]);

	function dragstart(evt: DragEvent) {
		evt.dataTransfer?.setData('card', JSON.stringify(card));
	}
</script>

{#if !isCardEmpty(card)}
	<div
		role="button"
		tabindex="0"
		ondragstart={(evt) => dragstart(evt)}
		class="card {inHand ? 'in-hand' : ''}"
		style="--suit-color: {suitColor}; --index: {index}; --distance: {distance}"
		draggable={draggable ? 'true' : 'false'}
	>
		<p class="rank-number left">
			{rank}
		</p>

		<div class="rank-art">
			{suitSymbol}
		</div>

		<p class="rank-number right">{rank}</p>
	</div>
{:else}
	<div
		class="card {inHand ? 'in-hand' : 'empty'}"
		style="--index: {index}; --distance: {distance}"
	></div>
{/if}

<style>
	.card {
		display: grid;
		grid-template-columns: repeat(5, 1fr);
		background-color: antiquewhite;
		aspect-ratio: 2.5 / 3.5;
		width: 50px;
		border-radius: 10px;
		border: 1px solid black;
		padding: 3px;
		color: var(--suit-color);
	}

	.card.in-hand {
		transition: all;
		transition-duration: 100ms;
		transition-timing-function: cubic-bezier();

		position: relative;
		transform-origin: bottom center;
		transform: rotate(calc(var(--distance) * 1.25deg)) translateY(calc(abs(var(--distance)) * 2px));
		left: calc(abs(var(--index)) * -10px);
	}

	.card:hover {
		scale: 1.05;
	}

	.card.empty {
		background-color: rgba(0, 0, 0, 0.2);
	}

	.card.empty:hover {
		scale: 1;
	}

	.rank-art {
		display: flex;
		grid-column: span 3;
		/* background-color: aliceblue; */
		text-align: center;
		justify-content: center;
		align-items: center;
		height: 100%;
		margin: 0px;
	}

	.rank-number {
		height: 100%;
		margin: 0px;
		display: flex;
		text-align: center;
	}

	.left {
	}

	.right {
		align-items: flex-end;
		text-align: center;
	}
</style>
