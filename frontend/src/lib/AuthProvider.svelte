<script lang="ts">
    import { navigate } from "sv-router/generated";
    import { type AuthContext, setAuthContext } from "./auth.svelte";
    import { createQuery } from "@tanstack/svelte-query";

    let authctx = $state<AuthContext>({ is_authenticated: false, user: null });
    setAuthContext(authctx);
    let { children } = $props();
    $effect(() => {
        console.log(
            `user is now: ${authctx.is_authenticated ? "authenticated" : "unauthenticated"}`,
        );
        if (!authctx.is_authenticated) {
            navigate("/login");
        }
    });

    const query = createQuery(() => ({
        queryKey: ["authCheck"],
        async queryFn() {
            const response = await fetch("/api/v1/auth/check");
            authctx.is_authenticated = response.ok;
            return await response.text();
        },
    }));
</script>

{#if query.isPending}
    <p>Loading Auth...</p>
{/if}

{#if !query.isPending}
    {@render children()}
{/if}
