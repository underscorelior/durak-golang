<script lang="ts">
	import { makeAction } from '$lib/event/actions';
	import { globalState } from '$lib/state/state.svelte';
	import {
		Action,
		canAttack,
		canDefend,
		canPass,
		canTake,
		cardDefends,
		isAttacking,
		isDefending
	} from '$lib/game';
	import Card from './card.svelte';
	import Pair from './pair.svelte';
	import Hand from './hand.svelte';
	import { isCardEmpty } from '$lib/util';

	let game_state = $derived(globalState.lobby?.game_state);

	let card_pairs = $derived(game_state?.turn.table_state);

	let pos = $derived(globalState.lobby?.position);
	let max_players = $derived(globalState.lobby?.max_players);

	let players = $derived(globalState.lobby?.players ?? []);

	let game_started = $derived(globalState.lobby?.is_started);

	const playerPosOffset = (pl: Player, off: number) =>
		pl.position == ((pos || 0) + off) % (max_players || 1);

	let player_id = $derived(globalState.user.user_id ?? '');
	let left = $derived(players.find((pl) => playerPosOffset(pl, 1)));
	let top = $derived(players.find((pl) => playerPosOffset(pl, 2)));
	let right = $derived(players.find((pl) => playerPosOffset(pl, 3)));

	let buttonDisabled = $derived(canPass() || canTake());
	let buttonShown = $derived(isAttacking() || isDefending());
	let buttonText = $derived(buttonShown ? (isAttacking() ? 'Pass' : 'Take') : '');

	function dragover(e: DragEvent) {
		e.preventDefault();
		e.dataTransfer.dropEffect = 'move';
	}

	function drop(e: DragEvent, pair: number) {
		e.preventDefault();
		const card = JSON.parse(e.dataTransfer.getData('card'));

		if (!card_pairs) {
			return;
		}

		const atk_card = card_pairs[pair].attack_card;

		console.log('test', isCardEmpty(atk_card));

		if (isCardEmpty(atk_card) && canAttack()) {
			makeAction(Action.Attack, card, pair);
		}

		if (!isCardEmpty(atk_card) && canDefend() && cardDefends(card_pairs[pair], card)) {
			makeAction(Action.Defense, card, pair);
		}
	}

	function buttonAction() {
		if (isAttacking() && canPass()) {
			makeAction(Action.Pass);
		}

		if (isDefending() && canTake()) {
			makeAction(Action.Take);
		}
	}
</script>

<div class="table">
	<div class="player-box" id="player-top" style="margin-inline: auto">
		{#if top}
			<Hand opponent user_id={top.user_id} orientation={2} />
		{/if}
	</div>
	<div class="player-box" id="player-left">
		{#if left}
			<Hand opponent user_id={left.user_id} orientation={1} />
		{/if}
	</div>
	<div class="tabletop">
		{#if game_started && card_pairs}
			{#each card_pairs as pair, i (i)}
				<div class="card-pair" ondragover={dragover} ondrop={(e) => drop(e, i)} role="group">
					<Pair {pair} />
				</div>
			{/each}
		{/if}
	</div>
	<div class="player-box" id="player-right">
		{#if right}
			<Hand opponent user_id={right.user_id} orientation={3} />
		{/if}
	</div>
	<div class="player-box" id="player-bottom" style="margin-inline: auto">
		{#if player_id}
			<Hand opponent={false} user_id={player_id} orientation={0} />
		{/if}
		{#if game_started}
			<button
				onclick={() => buttonAction()}
				disabled={!buttonDisabled}
				style={buttonShown ? '' : 'visibility: hidden;'}
				class="action-button"
			>
				{buttonText}
			</button>
		{/if}
	</div>

	{#if game_state}
		<div class="deck">
			<Card card={game_state.trump} draggable={false} />
		</div>
	{/if}
</div>

<style>
	.table {
		display: grid;
		grid-template-columns: 0.75fr repeat(3, 1fr) 0.75fr;
		grid-template-rows: 0.75fr repeat(3, 1fr) 0.75fr;
		width: 70%;

		gap: 1px 1px;
		margin-inline: auto;
	}

	.tabletop {
		grid-column: span 3;
		grid-row: span 3;
		background-color: aliceblue;
		border: 1px solid gray;

		display: grid;
		grid-template-columns: repeat(3, 1fr);
		grid-template-rows: repeat(2, 1fr);
	}

	.card-pair {
		/* border: 1px solid black; */
		display: flex;
		align-items: center;
		justify-content: center;
	}

	.deck {
		grid-column-start: 5;
		grid-row-start: 5;
		background-color: black;
	}

	.action-button {
		position: absolute;
		bottom: 0;
		right: 0;
	}

	.player-box {
		position: relative;

		display: flex;
		align-items: center;
		justify-content: center;

		min-width: 0;
		min-height: 0;

		border: 1px solid black;
		box-sizing: border-box;
	}

	#player-top {
		grid-column: span 3 / span 3;
		grid-column-start: 2;
	}

	#player-left {
		grid-row: span 3 / span 3;
		grid-row-start: 2;
	}

	#player-right {
		grid-row: span 3 / span 3;
		grid-row-start: 2;
		grid-column-start: 5;
	}

	#player-bottom {
		grid-column: span 3 / span 3;
		grid-column-start: 2;
		grid-row-start: 5;
	}
</style>
