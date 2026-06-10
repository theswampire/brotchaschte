<script lang="ts">
    import "./app.css";
    import { Router } from "sv-router";
    import "sv-router/generated";
    import { QueryClientProvider, QueryClient } from "@tanstack/svelte-query";
    import AuthProvider from "./lib/AuthProvider.svelte";

    const queryClient = new QueryClient();
</script>

<QueryClientProvider client={queryClient}>
    <AuthProvider>
        <Router />
    </AuthProvider>
</QueryClientProvider>
