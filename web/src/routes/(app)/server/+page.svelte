<script lang="ts">
  import { PUBLIC_COMMIT, PUBLIC_VERSION } from "$env/static/public";
  import { getApiClient, handleApiError } from "$lib";
  import { Button } from "@nanoteck137/nano-ui";
  import { Play } from "lucide-svelte";
  import { onMount } from "svelte";
  import toast from "svelte-5-french-toast";
  import { z } from "zod";

  // const { data } = $props();
  const apiClient = getApiClient();

  const TaskSyncStateEventTask = z.object({
    name: z.string(),
    isRunning: z.boolean(),
  });
  type TaskSyncStateEventTaskTy = z.infer<typeof TaskSyncStateEventTask>;

  const TaskSyncStateEvent = z.object({
    tasks: z.array(TaskSyncStateEventTask),
  });

  let tasks = $state<TaskSyncStateEventTaskTy[]>([]);

  onMount(() => {
    const eventSource = new EventSource(apiClient.url.sseHandler());

    eventSource.addEventListener("connected", () => {
      console.log("Connected to SSE handler");
    });

    eventSource.addEventListener("task-sync-state", (e) => {
      const data = TaskSyncStateEvent.parse(JSON.parse(e.data));

      console.log("tasks state", data);
      tasks = data.tasks;

      // isSyncing = data.isRunning;
      // errors = data.errors;
      // numArtists = data.numArtists;
      // numAlbums = data.numAlbums;
      // numTracks = data.numTracks;

      // artistSyncTime = data.artistsSyncDurationMs;
      // albumSyncTime = data.albumsSyncDurationMs;
      // trackSyncTime = data.tracksSyncDurationMs;
      // totalSyncTime = data.totalSyncDurationMs;
    });

    return () => {
      eventSource.close();
    };
  });
</script>

<p>Server Page (W.I.P)</p>

<p>Version: {PUBLIC_VERSION}</p>
<p>Commit: {PUBLIC_COMMIT}</p>

{#each tasks as task}
  <div class="flex items-center gap-2">
    <p>{task.name} - Running: {task.isRunning}</p>
    {#if !task.isRunning}
      <Button
        variant="ghost"
        size="icon"
        onclick={async () => {
          const res = await apiClient.runTask(task.name);
          if (!res.success) {
            return handleApiError(res.error);
          }

          toast.success("Dispatched task");
        }}
      >
        <Play />
      </Button>
    {/if}
  </div>
{/each}
