<script lang="ts">
  import Spacer from "$lib/components/Spacer.svelte";
  import { Card, ScrollArea } from "@nanoteck137/nano-ui";

  const { data, children } = $props();
</script>

<Card.Root>
  <div class="relative w-full overflow-hidden rounded-xl shadow-lg">
    <div
      class="h-full w-full bg-gradient-to-tr from-logo-1 via-logo-2 to-logo-3 p-6"
    >
      <h1 class="text-3xl font-bold">
        {data.userData.displayName}
      </h1>
    </div>
  </div>

  <div class="flex w-full justify-center gap-4">
    <ScrollArea class="px-2 pb-4 sm:pb-0" orientation="horizontal">
      <div class="flex gap-4 pt-4">
        <a
          class="border-b-2 px-2 py-3 text-sm font-medium hover:brightness-75"
          href="/users/{data.userData.id}"
        >
          Overview
        </a>

        {#if data.userData.id === data.user?.id}
          <a
            class="border-b-2 px-2 py-3 text-sm font-medium hover:brightness-75"
            href="/users/{data.userData.id}/settings"
          >
            Settings
          </a>
        {/if}
      </div>
    </ScrollArea>
  </div>
</Card.Root>

<Spacer size="md" />

{@render children()}
