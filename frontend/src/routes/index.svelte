<script lang="ts">
	import { getAuthContext } from "@/lib/auth.svelte";
	import { createMutation } from "@tanstack/svelte-query";

	console.log("Home");
	const authctx = getAuthContext();

	const logoutMutation = createMutation(() => ({
		async mutationFn() {
			const response = await fetch("/api/v1/auth/logout");
			if (!response.ok) {
				console.error(response.statusText);
				console.error(await response.text());
			}
			authctx.is_authenticated = false;
		},
	}));

	async function logout() {
		logoutMutation.mutate();
	}
</script>

<section>
	Home

	<button
		disabled={logoutMutation.isPending}
		class="border rounded p-1 border-gray-200 shadow hover:shadow-none transition-shadow"
		on:click={() => logout()}
		type="button">Logout</button
	>
</section>
