<script lang="ts">
  import { Button } from "@nanoteck137/nano-ui";

  let { data } = $props();
</script>

{#if !data.user}
  <div class="mt-16 flex flex-col items-center justify-center gap-8 p-4">
    <h1
      class="bg-gradient-to-tr from-logo-1 via-logo-2 to-logo-3 bg-clip-text text-5xl font-bold text-transparent"
    >
      Refinery
    </h1>
    <Button size="lg" href="/login">Login</Button>
  </div>
{:else}
  <div class="flex flex-col gap-8">
    <div class="flex flex-col items-center gap-4 border-b pb-4">
      <a
        class="bg-gradient-to-tr from-logo-1 via-logo-2 to-logo-3 bg-clip-text px-8 text-2xl font-medium text-transparent"
        href="/"
      >
        Refinery
      </a>

      <h1 class="text-center text-2xl font-bold">
        Welcome, {data.user.displayName}!
      </h1>
    </div>

    <!-- <div class="flex items-center gap-4">
      <a
        class="bg-gradient-to-tr from-logo-1 via-logo-2 to-logo-3 bg-clip-text text-2xl font-medium text-transparent"
        href="/"
      >
        Refinery
      </a>

      <span class="text-xl font-bold">|</span>

      <h1 class="text-2xl font-bold">
        Welcome, {data.user.displayName}!
      </h1>
    </div> -->
  </div>
{/if}
