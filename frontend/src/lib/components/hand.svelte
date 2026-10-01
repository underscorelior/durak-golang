<script lang="ts">
	import { globalState } from "$lib/state/state.svelte";
	import Card from './card.svelte';

    
    let hand = $derived(globalState.lobby?.game_state?.hand);
    let middle = $derived((hand || []).length/2)
</script>

<div class="hand">
    {#each hand as card, index (index)}
        <Card {card} {index} {middle} />
    {/each}
</div>

<style>
    .hand {
        display: flex;
        flex-direction: row;
        gap: 5px 5px;
    }
</style>
