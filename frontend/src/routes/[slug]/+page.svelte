<script lang="ts">
	import Lobby from "$lib/components/lobby.svelte";
	import { joinLobby } from "$lib/event/actions";
	import { globalState } from "$lib/state/state.svelte";
    import { page } from '$app/state';
    import { goto } from '$app/navigation';
    import { resolve } from '$app/paths';
	import { untrack } from "svelte";

    const lobby_code = $derived(page.params.slug);

    // TODO: Find out when this rerenders and if connected is false, what to do then (i.e. lag?)
    $effect(() => {
        if (globalState.menu.connection.connected)
        untrack(() => {
            if (globalState.lobby === null) {
                if (!lobby_code) {
                    goto(resolve('/'));
                    return
                }
                joinLobby(lobby_code)
            }
        })
    })
    
</script>

<Lobby/>