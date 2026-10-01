<script lang="ts">
	import { shID } from "$lib/util";

    const { player } : {player?: GamePlayer} = $props()


</script>

{#if player != undefined}
    <div>
        {shID(player.user_id)} ({player.position})
    </div>

    <div style="position: relative; margin: 20px 40px;">
        {#each Array(player.hand_size) as card, i (i)}
            <div class="card-placeholder" style="--index: {i}; --middle: {5/2}">{card}</div>
        {/each}
    </div>
{/if}


<style>
    .card-placeholder {
        width: 25px;
        height: 35px;
        background-color: aliceblue;
        border: 1px solid black;
        border-radius: 5px;
        position: absolute;
        transform-origin: center 500%;
        transform: translate(-50%, -50%) rotate(calc((var(--middle) - var(--index)) * -2.5deg));
    }
</style>