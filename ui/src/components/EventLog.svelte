<script lang="ts">
  import { eventLog } from '$lib/store';

  let logContainer: HTMLDivElement;

  // Auto-scroll to bottom on new entries
  $: if ($eventLog && logContainer) {
    setTimeout(() => {
      logContainer.scrollTop = logContainer.scrollHeight;
    }, 0);
  }
</script>

<aside class="event-log">
  <div class="event-log-header">
    <span class="event-log-title">Live Log</span>
    <span class="log-count">{$eventLog.length}</span>
  </div>

  <div class="log-entries" bind:this={logContainer}>
    {#if $eventLog.length === 0}
      <div class="log-empty">No events yet</div>
    {:else}
      {#each $eventLog as entry, i (i)}
        <div class="log-entry {entry.type}">
          {entry.message}
        </div>
      {/each}
    {/if}
  </div>
</aside>

<style>
  .event-log {
    position: fixed;
    top: 64px;
    right: 0;
    bottom: 0;
    width: 264px;
    background: var(--bg-subtle);
    border-left: 1px solid var(--border);
    display: flex;
    flex-direction: column;
    z-index: 50;
    overflow: hidden;
  }

  .event-log-header {
    padding: 0.875rem 1rem 0.75rem;
    border-bottom: 1px solid var(--border);
    display: flex;
    align-items: center;
    justify-content: space-between;
    flex-shrink: 0;
  }

  .event-log-title {
    font-size: 0.7rem;
    font-weight: 700;
    letter-spacing: 0.1em;
    text-transform: uppercase;
    color: var(--text-muted);
  }

  .log-count {
    font-size: 0.7rem;
    font-weight: 700;
    color: var(--text-muted);
    background: var(--bg-elevated);
    border-radius: var(--r-full);
    min-width: 20px;
    height: 20px;
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 0 5px;
    font-variant-numeric: tabular-nums;
  }

  .log-entries {
    flex: 1;
    overflow-y: auto;
    padding: 0.5rem;
    display: flex;
    flex-direction: column;
    gap: 0.25rem;
  }

  .log-empty {
    padding: 1.5rem 0.75rem;
    text-align: center;
    font-size: 0.78rem;
    color: var(--text-muted);
  }

  .log-entry {
    padding: 0.4rem 0.625rem;
    border-radius: var(--r-sm);
    font-size: 0.78rem;
    line-height: 1.45;
    word-break: break-word;
    animation: entryIn 0.2s var(--ease);
  }

  @keyframes entryIn {
    from { opacity: 0; transform: translateX(6px); }
    to   { opacity: 1; transform: translateX(0); }
  }

  .log-entry.info {
    color: var(--text-secondary);
  }

  .log-entry.success {
    color: var(--success);
    background: var(--success-dim);
  }

  .log-entry.error {
    color: var(--danger);
    background: var(--danger-dim);
  }
</style>


