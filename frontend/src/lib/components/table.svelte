<script lang="ts">
	import { globalState } from "$lib/state/state.svelte";
	import { shID } from "$lib/util";
	import Card from './card.svelte';

    
    let trump = $derived(globalState.game?.trump)
</script>

<div class="table">
    <div class="player-box" id="player-top" style="margin-inline: auto">
        2,0
        -
    </div>
    <div class="player-box" id="player-left">0,2-</div>
    <div class="tabletop">1,1-4,4

    </div>
    <div class="player-box" id="player-right">5,2-</div>
    <div class="player-box" id="player-bottom" style="margin-inline: auto">
    2,5
        {shID(globalState.user.user_id)} ({globalState.lobby?.position})
    </div>

    {#if trump !== undefined}
        <div class="deck">
            <Card card={trump} draggable={false} />
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