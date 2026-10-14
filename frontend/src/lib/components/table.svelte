<script lang="ts">
	import { makeAction } from "$lib/event/actions";
	import { globalState } from "$lib/state/state.svelte";
	import { shID } from "$lib/util";
    import {Action, isAttacking, isDefending} from "$lib/game";
	import Card from './card.svelte';
	import Hand from "./hand.svelte";
	import Opponent from "./opponent.svelte";
	import Pair from "./pair.svelte";

    
    let game_state = $derived(globalState.lobby?.game_state)


    let card_pairs = $derived(game_state?.turn.table_state)
    
    let pos = $derived(globalState.lobby?.position)
    let max_players = $derived(globalState.lobby?.max_players)
    
    let players = $derived(game_state?.players ?? [])

    let game_started = $derived(globalState.lobby?.is_started)


    const playerPosOffset = (pl: GamePlayer, off: number) => pl.position == ((pos || 0) + off) % (max_players || 1)

    let left = $derived(players.find((pl) => playerPosOffset(pl, 1)))
    let front  = $derived(players.find((pl) => playerPosOffset(pl, 2)))
    let right = $derived(players.find((pl) => playerPosOffset(pl, 3)))

    function dragover(e: DragEvent) {
        e.preventDefault(); 
        e.dataTransfer.dropEffect = 'move';
    }
    
    function drop(e: DragEvent, pair: number) {
        e.preventDefault(); 
        const card = JSON.parse(e.dataTransfer.getData('card'))

        if (!card_pairs) {
            return
        }

        const atk_card = card_pairs[pair].attack_card

        if (atk_card.rank != 0 || atk_card.suit != 0 && isDefending()) {
            makeAction(Action.Defense, card, pair)
        }
        
        if (atk_card.rank == 0 && atk_card.suit == 0 && isAttacking()) {
            makeAction(Action.Attack, card, pair)
        }
    }
</script>

<div class="table">
    <div class="player-box" id="player-top" style="margin-inline: auto">
        <Opponent player={front} />
    </div>
    <div class="player-box" id="player-left">
        <Opponent player={left} />
    </div>
    <div class="tabletop">

    {#if game_started && card_pairs}
        {#each card_pairs as pair, i (i)}
            <div class="card-pair" ondragover={dragover} ondrop={(e)=>drop(e, i)} role="group">
                <Pair {pair} />
            </div>
        {/each}
    {/if}
    </div>
    <div class="player-box" id="player-right">
        <Opponent player={right} />
    </div>
    <div class="player-box" id="player-bottom" style="margin-inline: auto">
        {shID(globalState.user.user_id)} ({globalState.lobby?.position})
        <br />
        <br />
        <br />
        {#if globalState.lobby?.is_started}
            <div>
                <Hand/>
            </div>
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
        width: 100%;
        aspect-ratio: 1.5;
        gap: 1px 1px;
    }

    .tabletop {
        grid-column: span 3;
        grid-row: span 3;
        background-color: aliceblue;
        width: 100%;
        height: 100%;
        border: 1px solid gray;

        display: grid;
        grid-template-columns: repeat(3, 1fr);
        grid-template-rows: repeat(2, 1fr);
    }

    .card-pair {
        border: 1px solid black;
    }

    .deck {
        grid-column-start: 5;
        grid-row-start: 5;
        background-color: black;
    }

	.player-box {
		border: 1px solid black;
		text-align: center;
        width: 100%;
        height: 100%;
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