<script lang="ts">
  import { api } from "@/lib/api/api";
  import { createMutation, useQueryClient } from "@tanstack/svelte-query";
  import { navigate, route } from "sv-router/generated";

  const email = route.getParams("/users/:userId").userId;

  const queryClient = useQueryClient();
  const editPasswordMutation = createMutation(() => ({
    mutationKey: ["password", email],
    mutationFn: api.users.updateUserPassword,
    onSuccess() {
      queryClient.invalidateQueries({ queryKey: ["users"] });
      navigate("/users", { viewTransition: true });
    },
  }));
  const deleteUserMutation = createMutation(() => ({
    mutationKey: ["delete", email],
    mutationFn: api.users.deleteUser,
    onSuccess() {
      queryClient.invalidateQueries({ queryKey: ["users"] });
      navigate("/users", { viewTransition: true });
    },
  }));

  let changed = $state(false);

  function editPasswordHandler(
    ev: SubmitEvent & {
      currentTarget: EventTarget & HTMLFormElement;
    },
  ) {
    ev.preventDefault();
    const data = new FormData(ev.currentTarget);
    editPasswordMutation.mutate(data);
  }
</script>

<section class="absolute inset-0 bg-white p-4">
  <h1 class="font-bold text-xl">Benutzer {email}</h1>

  <form onsubmit={editPasswordHandler}>
    <input name="email" type="hidden" value={email} />
    <label for="password">Passwort</label>
    <input
      id="password"
      type="password"
      name="password"
      class="border rounded py-1 px-2 border-gray-200 w-20"
      oninput={() => (changed = true)}
    />
    <button
      disabled={!changed ||
        editPasswordMutation.isPending ||
        deleteUserMutation.isPending}
      type="submit"
      class="border rounded p-1 border-gray-200 shadow hover:shadow-none transition-shadow disabled:bg-gray-200 disabled:text-gray-500 disabled:shadow-none mt-5"
      >Passwort setzen</button
    >
  </form>

  <button
    disabled={editPasswordMutation.isPending || deleteUserMutation.isPending}
    class="border rounded p-1 border-gray-200 shadow hover:shadow-none transition-shadow disabled:bg-gray-200 disabled:text-gray-500 disabled:shadow-none mt-5"
    onclick={() => deleteUserMutation.mutate(email)}
    type="button">Benutzer löschen</button
  >
</section>
