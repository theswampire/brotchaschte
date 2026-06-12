<script lang="ts">
    import { navigate } from "sv-router/generated";
    import {
        type AuthContext,
        auth,
        initAuthStuff,
        navigateToLoginWithCallback,
        setAuthContext,
    } from "./auth.svelte";
    import { createQuery } from "@tanstack/svelte-query";
    import { registerAuthCallbacks } from "./api/axios";
    import { getCookie } from "./cookie";

    const initial_authctx: AuthContext = {
        is_authenticated: false,
        user: null,
    };

    const authcookie = getCookie("authsession");
    if (authcookie) {
        console.log(`cookie: ${authcookie}`);
        if (authcookie.startsWith('"authenticated,')) {
            const email = atob(authcookie.substring(15, authcookie.length - 1));
            console.log(`\t- ${email}`);
            initial_authctx.is_authenticated = true;
            initial_authctx.user = email;
        }
    }
    let authctx = $state<AuthContext>(initial_authctx);

    setAuthContext(authctx);
    let { children } = $props();
    $effect(() => {
        console.log(
            `user is now: ${authctx.is_authenticated ? "authenticated" : "unauthenticated"}`,
        );
        if (!authctx.is_authenticated) {
            navigateToLoginWithCallback(location.pathname);
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

    initAuthStuff();
    registerAuthCallbacks({ logout: auth.setUnauthenticated });
</script>

{#if query.isPending}
    <p>Loading Auth...</p>
{/if}

{#if !query.isPending}
    {@render children()}
{/if}
