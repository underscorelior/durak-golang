<script lang="ts">
	import { rankToString, suitToColor, suitToSymbol } from "$lib/util";

    const { card, draggable = true }: { card: Card, draggable?: boolean } = $props()

    const rank = $derived(rankToString[card.rank])
    const suitSymbol = $derived(suitToSymbol[card.suit])
    const suitColor = $derived(suitToColor[card.suit])
</script>


<div class="card" style="--suit-color: {suitColor}" draggable={draggable ? "true" : 'false'}>
    <p class="rank-number left">
        {rank}
    </p>

    <div class="rank-art">
        {suitSymbol}
    </div>

    <p class="rank-number right">{rank}</p>
    
</div>

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
        margin: 10px;
        color: var(--suit-color);
        
    }

    .rank-art {
        display: flex;
        grid-column: span 3;
        background-color: aliceblue;
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