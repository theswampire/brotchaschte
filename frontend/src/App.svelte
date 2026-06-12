<script lang="ts">
    import "./app.css";
    import { Router } from "sv-router";
    import "sv-router/generated";
    import { QueryClientProvider, QueryClient } from "@tanstack/svelte-query";
    import AuthProvider from "./lib/AuthProvider.svelte";
    import ApiClientProvider from "./lib/api/ApiClientProvider.svelte";

    const queryClient = new QueryClient();
</script>

<QueryClientProvider client={queryClient}>
    <ApiClientProvider>
        <AuthProvider>
            <Router />
        </AuthProvider>
    </ApiClientProvider>
</QueryClientProvider>
