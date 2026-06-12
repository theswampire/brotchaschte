<script lang="ts">
	import { createMutation } from "@tanstack/svelte-query";
	import type { Path } from "sv-router";
	import { navigate, type Routes, route } from "sv-router/generated";
	import { auth } from "@/lib/auth.svelte";

	let callback_url = route.search.callback;

	if (!callback_url) {
		callback_url = "/";
	}
	console.log(callback_url);

	const callback = new URL(`thismessage:${callback_url}`);

	$effect(() => {
		if (auth.isAuthenticated()) {
			navigate(callback.pathname as Path<Routes>, {
				search: callback.search,
				hash: callback.hash,
			});
		}
	});

	const submitMutation = createMutation(() => ({
		async mutationFn(formData: FormData) {
			const response = await fetch("/api/v1/auth/login", {
				method: "POST",
				body: formData,
			});
			if (!response.ok) {
				throw new Error("Anmeldung fehlgeschlagen");
			}
			auth.setAuthenticated("not implemented");
			navigate(callback.pathname as Path<Routes>, {
				search: callback.search,
				hash: callback.hash,
			});
		},
		async onSuccess(data, variables, onMutateResult, context) {
			console.log("yess");
		},
		onError(error, variables, onMutateResult, context) {
			console.log("grrr");
		},
	}));

	async function handleSubmit(event: SubmitEvent) {
		event.preventDefault();
		const form = event.currentTarget as HTMLFormElement;
		const formData = new FormData(form);

		submitMutation.mutate(formData);
	}
</script>

<section>
	Login
	<div>
		<form
			class="border-gray-300 border rounded-sm p-3 w-fit mx-auto bg-gray-50 mt-10 shadow"
			onsubmit={handleSubmit}
		>
			<label
				for="email"
				class="text-sm font-semibold text-violet-950 bg-white"
				>E-Mail:</label
			>
			<input
				id="email"
				name="email"
				type="text"
				required
				class="border-2 border-blue-300 rounded px-2 block"
			/>
			<label for="password" class="text-sm font-semibold text-violet-950"
				>Password:</label
			>
			<input
				id="password"
				name="password"
				type="password"
				required
				class="border-2 border-blue-300 rounded px-2 block"
			/>
			<button
				type="submit"
				disabled={submitMutation.isPending}
				class="border bg-blue-700 text-white border-blue-200 text-sm p-1 mt-2 font-semibold w-full rounded shadow disabled:bg-blue-700/60"
			>
				Login
			</button>
		</form>
		{#if submitMutation.isError}
			<p>{submitMutation.error.message}</p>
		{/if}
	</div>
</section>
