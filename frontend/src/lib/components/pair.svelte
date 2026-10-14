<script lang="ts">
    import Card from './card.svelte'
    const { pair }: {pair: CardPair} = $props()
</script>

<div class="pair-outside">
    <div class="bottom-card">
        <Card card={pair.attack_card} draggable={false} />
    </div>
    {#if pair.is_defended}
        <div class="top-card">
            <Card card={pair.defense_card} draggable={false} />
        </div>
    {/if}
</div>


<style>
    .pair-outside {
        position: relative;

        width: 100%;
        height: 100%;
    }

    .bottom-card {
        position: absolute;
    }

    .top-card {
        position: relative;
        top: 10%;
        left: 7%;
    }
</style>