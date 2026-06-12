<script lang="ts">
	import { auth } from "@/lib/auth.svelte";
	import { createMutation } from "@tanstack/svelte-query";
	import { p } from "sv-router/generated";

	console.log("Home");

	const logoutMutation = createMutation(() => ({
		async mutationFn() {
			const response = await fetch("/api/v1/auth/logout");
			if (!response.ok) {
				console.error(response.statusText);
				console.error(await response.text());
			}
			auth.setUnauthenticated();
		},
	}));

	async function logout() {
		logoutMutation.mutate();
	}
</script>

<section>
	Home

	<a href={p("/kek")}>Kek</a>
	<a href={p("/")}>Home</a>
	<a href={p("/login")}>Login</a>

	<button
		disabled={logoutMutation.isPending}
		class="border rounded p-1 border-gray-200 shadow hover:shadow-none transition-shadow"
		on:click={() => logout()}
		type="button">Logout</button
	>
</section>
