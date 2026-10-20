<script lang="ts">
    import { createMutation } from "@tanstack/svelte-query";
    import { api } from "@/lib/api/api";
    import { auth } from "@/lib/auth/auth.svelte";

    const logoutMutation = createMutation(() => ({
        mutationFn: api.auth.logout,
        onSuccess() {
            auth.setUnauthenticated();
        },
    }));

    async function logout() {
        logoutMutation.mutate();
    }
</script>

<button
    disabled={logoutMutation.isPending}
    class="border rounded p-1 border-gray-200 shadow hover:shadow-none transition-shadow"
    on:click={() => logout()}
    type="button">Logout</button
>
