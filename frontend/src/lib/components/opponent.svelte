<script lang="ts">
	import { globalState } from "$lib/state/state.svelte";
	import { shID } from "$lib/util";

    const { player } : {player?: GamePlayer} = $props()

    let defender = $derived(globalState.lobby?.game_state?.turn.defender_id)
    let attackers = $derived(globalState.lobby?.game_state?.turn.attacker_ids)
    let initial_attacker = $derived(globalState.lobby?.game_state?.turn.initial_attacker_id)

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

    {#if player.user_id == defender}
        Defender
    {/if}

    {#if attackers?.includes(player.user_id)}
        Attacker
    {/if}

    {#if player.user_id == initial_attacker}
        Initial
    {/if}
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