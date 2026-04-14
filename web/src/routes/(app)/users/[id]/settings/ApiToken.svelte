<script lang="ts">
  import { invalidateAll } from "$app/navigation";
  import { getApiClient, handleApiError } from "$lib";
  import type { ApiToken } from "$lib/api/types";
  import ConfirmBox from "$lib/components/ConfirmBox.svelte";
  import { Button, Dialog } from "@nanoteck137/nano-ui";
  import { Eye, Trash } from "lucide-svelte";
  import toast from "svelte-5-french-toast";

  type Props = {
    token: ApiToken;
  };

  const { token }: Props = $props();
  const apiClient = getApiClient();

  let openTokenShow = $state(false);
  let openDeleteModal = $state(false);
</script>

<div class="flex items-center justify-between border-b py-2">
  <p>{token.name}</p>
  <div class="flex gap-2">
    <Button
      variant="outline"
      size="icon"
      onclick={() => {
        openTokenShow = true;
      }}
    >
      <Eye />
    </Button>

    <Button
      size="icon"
      variant="destructive"
      onclick={() => {
        openDeleteModal = true;
      }}
    >
      <Trash />
    </Button>
  </div>
</div>

<Dialog.Root bind:open={openTokenShow}>
  <Dialog.Content>
    <Dialog.Header>
      <Dialog.Title>Token Value</Dialog.Title>
    </Dialog.Header>

    <p>{token.id}</p>
  </Dialog.Content>
</Dialog.Root>

<ConfirmBox
  bind:open={openDeleteModal}
  title="Delete Token?"
  description="Are you sure you want to delete this token? This action cannot be undone."
  confirmText="Delete"
  onResult={async () => {
    const res = await apiClient.deleteApiToken(token.id);
    if (!res.success) {
      return handleApiError(res.error);
    }

    toast.success("Successfully deleted api token");
    invalidateAll();
  }}
/>
