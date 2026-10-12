<script lang="ts">
	import { globalState } from "$lib/state/state.svelte";
	import { shID } from "$lib/util";
	import Card from './card.svelte';
	import Hand from "./hand.svelte";
	import Opponent from "./opponent.svelte";

    
    let game = $derived(globalState.lobby?.game_state)
    let card_pairs = $derived(globalState.lobby?.game_state?.turn.table_state)
    let pos = $derived(globalState.lobby?.position)
    let max_players = $derived(globalState.lobby?.max_players)
    // FIXME: Allow for updating of the state, while keeping the typesafety (|| [] seems to break the reactivity)
    let players = $derived(globalState.lobby?.game_state?.players || [])

    let game_started = $derived(globalState.lobby?.is_started)


    const playerPosOffset = (pl: GamePlayer, off: number) => pl.position == ((pos || 0) + off) % (max_players || 1)

    let left = $derived(players.find((pl) => playerPosOffset(pl, 1)))
    let front  = $derived(players.find((pl) => playerPosOffset(pl, 2)))
    let right = $derived(players.find((pl) => playerPosOffset(pl, 3)))

    function dragover(e) {
        e.preventDefault(); 
        e.dataTransfer.dropEffect = 'move';
    }
    
    function drop(e, pair: number) {
        e.preventDefault(); 
        console.log(card_pairs)
        card_pairs[pair] = {
            'attack_card': {
                suit: e.dataTransfer.getData('suit'),
                rank: e.dataTransfer.getData('rank')
            }
        }

        console.log(e.dataTransfer.getData('suit'), e.dataTransfer.getData('rank'))
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

    {#if game_started}
        {#each Array(6).keys() as pair (pair)}
            <div class="card-pair" ondragover={dragover} ondrop={(e)=>drop(e, pair)} role="group">
                {pair}
                <Card card={card_pairs[pair]['attack_card']} draggable={false} />
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

    {#if game != undefined}
        <div class="deck">
            <Card card={game.trump} draggable={false} />
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