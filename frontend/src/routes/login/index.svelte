<script lang="ts">
	import { createMutation } from "@tanstack/svelte-query";
	import axios from "axios";
	import type { Path } from "sv-router";
	import { navigate, type Routes, route } from "sv-router/generated";
	import type { LoginRes } from "@/lib/api/auth";
	import { apiV1client } from "@/lib/api/axios";
	import { auth, getAuthContext } from "@/lib/auth.svelte";
	import { api } from "@/lib/api/api";

	let callback_url = route.search.callback;

	if (!callback_url) {
		callback_url = "/";
	}
	console.log(callback_url);

	const callback = new URL(`thismessage:${callback_url}`);
	const authctx = getAuthContext();

	$effect(() => {
		if (authctx.is_authenticated) {
			navigate(callback.pathname as Path<Routes>, {
				search: callback.search,
				hash: callback.hash,
			});
		}
	});

	const submitMutation = createMutation(() => ({
		mutationFn: api.auth.login,
		onSuccess(data) {
			auth.setAuthenticated(data.data.email);
			navigate(callback.pathname as Path<Routes>, {
				search: callback.search,
				hash: callback.hash,
			});
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
		{#if submitMutation.isError && axios.isAxiosError(submitMutation.error)}
			<p>{submitMutation.error.response?.data.error}</p>
		{/if}
		{#if submitMutation.isError && !axios.isAxiosError(submitMutation.error)}
			<p>{submitMutation.error}</p>
		{/if}
	</div>
</section>
