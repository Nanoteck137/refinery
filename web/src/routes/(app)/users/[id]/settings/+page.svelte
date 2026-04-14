<script lang="ts">
  import { Button, Card } from "@nanoteck137/nano-ui";
  import ChangeDisplayName from "./ChangeDisplayName.svelte";
  import ApiToken from "./ApiToken.svelte";
  import NewApiTokenModal from "./NewApiTokenModal.svelte";
  import { Plus } from "lucide-svelte";
  import QuickCodeModal from "./QuickCodeModal.svelte";

  let { data } = $props();

  let openNewApiTokenModal = $state(false);
  let openQuickCodeModal = $state(false);
</script>

<Card.Root>
  <ChangeDisplayName />

  <div class="flex flex-col items-center gap-4 border-b p-6">
    <h2 class="text-bold text-center text-xl">
      Quick Code
      <Button
        variant="ghost"
        size="icon"
        onclick={() => {
          openQuickCodeModal = true;
        }}
      >
        <Plus />
      </Button>
    </h2>
  </div>

  <div class="flex flex-col items-center gap-4 border-b p-6">
    <h2 class="text-bold text-center text-xl">
      API Tokens
      <Button
        variant="ghost"
        size="icon"
        onclick={() => {
          openNewApiTokenModal = true;
        }}
      >
        <Plus />
      </Button>
    </h2>

    <div class="flex w-full flex-col sm:max-w-[460px]">
      {#each data.tokens as token}
        <ApiToken {token} />
      {/each}
    </div>
  </div>
</Card.Root>

<NewApiTokenModal bind:open={openNewApiTokenModal} />
<QuickCodeModal bind:open={openQuickCodeModal} />
