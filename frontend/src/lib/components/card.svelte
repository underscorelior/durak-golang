<script lang="ts">
	import { rankToString, suitToColor, suitToSymbol } from "$lib/util";

    const { card, index = 0, middle = 0, draggable = true }: { card: Card, index?: number, middle?: number, draggable?: boolean } = $props()

    const rank = $derived(rankToString[card.rank])
    const suitSymbol = $derived(suitToSymbol[card.suit])
    const suitColor = $derived(suitToColor[card.suit])
</script>


<div class="card" style="--suit-color: {suitColor}; --index: {index}; --middle: {middle}" draggable={draggable ? "true" : 'false'}>
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
        
        transition: all;
        transition-duration: 100ms;
        transition-timing-function: cubic-bezier();

        /* position: absolute; */
        /* transform-origin: center 500%; */
        /* transform: translate(-50%, -50%) rotate(calc((var(--middle) - var(--index)) * -2.5deg)); */
    }

    .card:hover {
        scale: 1.05;
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