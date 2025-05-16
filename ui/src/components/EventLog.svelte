<script lang="ts">
    import { eventLog} from '$lib/store';

    let logContainer: HTMLDivElement;

    $: {
        // Auto-scroll to bottom when new log is added
        if (logContainer) {
            setTimeout(() => {
            logContainer.scrollTop = logContainer.scrollHeight;
        }, 0);
        }
    }
</script>

<div class="event-log">
  <h3>Game Log</h3>
  <div class="log-entries" bind:this={logContainer}>
    {#each $eventLog as log}
      <div class="log-entry {log.type}">
        {log.message}
      </div>
    {/each}
  </div>
</div>


