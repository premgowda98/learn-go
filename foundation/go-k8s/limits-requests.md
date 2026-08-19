# Kubernetes CPU Limits vs Requests: Concise Notes

## Key Concepts
- **CPU Request**: Minimum CPU guaranteed to a pod. Scheduler uses this to place pods on nodes.
- **CPU Limit**: Maximum CPU a pod can use. If exceeded, the pod is throttled.

## Why CPU Limits Are Harmful
- CPU limits can cause unnecessary throttling, even when CPU is available.
- Throttling leads to degraded performance and higher latency.
- CPU is a renewable resource—using 100% now doesn't affect future availability.

## Three Analogies
1. **No limits, no requests**: Greedy pods can starve others (CPU starvation).
2. **With limits**: Even if resources are available, limits prevent usage, causing unnecessary failures.
3. **No limits, with requests**: Each pod gets its guaranteed share, but can use more if available—best outcome.

## Best Practices
- **Set CPU requests for all pods.**
- **Do not set CPU limits** unless you have a specific reason.
- Accurate requests prevent starvation without causing throttling.

## Memory Is Different
- Always set both memory requests and limits (set them equal).
- Memory is not compressible—if a pod exceeds its memory limit, it is killed.

## Visual Summary
<img src="https://cdn.prod.website-files.com/633e9bad8f71dfa75ae4c9db/672344d761a0dd3fa875e9be_6357fcc4b3a1634d362a408a_CPU%2520Limits.webp" alt="CPU Limits vs Requests" width="600" height="600">

---
For more details, see the full article: [Stop Using CPU Limits on Kubernetes](https://home.robusta.dev/blog/stop-using-cpu-limits)
