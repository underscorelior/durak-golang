<script lang="ts">
	import { globalState } from "$lib/state/state.svelte";
	import { shID } from "$lib/util";
	import Card from './card.svelte';
	import Hand from "./hand.svelte";
	import Opponent from "./opponent.svelte";

    
    let game = $derived(globalState.lobby?.game_state)
    let pos = $derived(globalState.lobby?.position)
    let max_players = $derived(globalState.lobby?.max_players)
    let players = $derived(globalState.lobby?.game_state?.players || [])


    const playerPosOffset = (pl: GamePlayer, off: number) => pl.position == ((pos || 0) + off) % (max_players || 1)

    let left = $derived(players.find((pl) => playerPosOffset(pl, 1)))
    let front  = $derived(players.find((pl) => playerPosOffset(pl, 2)))
    let right = $derived(players.find((pl) => playerPosOffset(pl, 3)))
    
</script>

<div class="table">
    <div class="player-box" id="player-top" style="margin-inline: auto">
        <Opponent player={front} />
    </div>
    <div class="player-box" id="player-left">
        <Opponent player={left} />
    </div>
    <div class="tabletop">1,1-4,4

    </div>
    <div class="player-box" id="player-right">
        <Opponent player={right} />
    </div>
    <div class="player-box" id="player-bottom" style="margin-inline: auto">
        {shID(globalState.user.user_id)} ({globalState.lobby?.position})
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