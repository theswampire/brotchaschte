<script lang="ts">
  import { api } from "@/lib/api/api";
  import {
    createMutation,
    createQuery,
    useQueryClient,
  } from "@tanstack/svelte-query";
  import { p } from "sv-router/generated";

  const queryClient = useQueryClient();
  const query = createQuery(() => ({
    queryKey: ["users"],
    queryFn: api.users.getUsers,
  }));

  const createUserMutation = createMutation(() => ({
    mutationFn: api.users.createUser,
    onSuccess() {
      queryClient.invalidateQueries({ queryKey: ["users"] });
    },
  }));

  function onCreateHandler(ev: SubmitEvent) {
    ev.preventDefault();
    const data = new FormData(ev.currentTarget as HTMLFormElement);
    createUserMutation.mutate(data);
  }
</script>

<section class="absolute inset-0 bg-white p-4">
  <h1 class="font-bold text-xl">Benutzer erstellen</h1>
  <form onsubmit={onCreateHandler}>
    <label for="email" class="block">Nutzername</label>
    <input
      id="email"
      type="text"
      name="email"
      class="border rounded py-1 px-2 border-gray-200"
    />
    <label for="password" class="block">Passwort</label>
    <input
      id="password"
      type="password"
      name="password"
      class="border rounded py-1 px-2 border-gray-200"
    />
    <button
      type="submit"
      disabled={createUserMutation.isPending}
      class="block mt-2 mb-5 border rounded p-1 border-gray-200 shadow hover:shadow-none transition-shadow disabled:bg-gray-200 disabled:text-gray-500 disabled:shadow-none"
      >Erstellen</button
    >
  </form>
  <h1 class="font-bold text-xl">Alle Benutzer</h1>
  {#if query.isPending}
    Lädt...
  {:else if query.isSuccess}
    <ul>
      {#each query.data.data.users as user}
        <li>
          <a href={p("/users/:userId", { params: { userId: user.email } })}
            >{user.email}</a
          >
        </li>
      {/each}
    </ul>
  {/if}
</section>
