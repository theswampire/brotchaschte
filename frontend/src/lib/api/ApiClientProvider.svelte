<script lang="ts">
    import { initApiClient } from "./axios";

    let { children } = $props();

    initApiClient();
</script>

{@render children()}
