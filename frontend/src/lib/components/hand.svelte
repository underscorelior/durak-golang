<script lang="ts">
	import { globalState } from "$lib/state/state.svelte";
	import Card from './card.svelte';

    
    let hand = $derived(globalState.lobby?.game_state?.hand);
    let middle = $derived((hand || []).length/2)

    let attacking = $derived(globalState.lobby?.game_state?.turn.initial_attacker_id == globalState.user.user_id)
</script>

<div class="hand" style="--bg: {attacking? 'red':'aliceblue'}">
    {#each hand as card, index (index)}
        <Card {card} {index} {middle} />
    {/each}
</div>

<style>
    .hand {
        display: flex;
        flex-direction: row;
        gap: 5px 5px;
        background-color: var(--bg);
    }
</style>
