# OS Assignment 7 — Banker's Algorithm & Deadlock Detection

**Subject:** Operating Systems | **Topic:** Deadlock Avoidance & Detection

---

## 1. Theory

### 1.1 Banker's Algorithm (Deadlock Avoidance)
Proposed by Dijkstra, the Banker's Algorithm proactively prevents deadlock by only granting a resource request if the resulting state is **safe**.

| Structure | Description |
|-----------|-------------|
| `Available[j]` | Free instances of resource Rj |
| `Max[i][j]` | Maximum demand of Pi for Rj |
| `Allocation[i][j]` | Resources currently allocated to Pi |
| `Need[i][j]` | Remaining need = Max[i][j] - Allocation[i][j] |

```
Safety Check:
  Work = Available; Finish[all] = false
  Find i: Finish[i]=false AND Need[i] <= Work
    Work += Allocation[i]; Finish[i] = true
  If all Finish[i]=true -> SAFE; else -> UNSAFE
```

### 1.2 Deadlock Detection Algorithm
Detects **existing** deadlock using the current `Request` matrix (actual pending requests), not worst-case `Need`.

```
Detection:
  Work = Available
  Finish[i] = (Allocation[i] == 0)
  Find i: Finish[i]=false AND Request[i] <= Work
    Work += Allocation[i]; Finish[i] = true
  Processes with Finish[i]=false -> DEADLOCKED
```

---

## 2. Python Source Code

```python
def bankers_algorithm(num_processes, num_resources, available, max_matrix, allocation):
    need = [[max_matrix[i][j] - allocation[i][j]
             for j in range(num_resources)]
            for i in range(num_processes)]

    work = available[:]
    finish = [False] * num_processes
    safe_seq = []

    count = 0
    while count < num_processes:
        found = False
        for i in range(num_processes):
            if finish[i]:
                continue
            if all(need[i][j] <= work[j] for j in range(num_resources)):
                work = [work[j] + allocation[i][j] for j in range(num_resources)]
                finish[i] = True
                safe_seq.append(i)
                count += 1
                found = True
        if not found:
            break

    return all(finish), safe_seq, need


def deadlock_detection(num_processes, num_resources, available, request, allocation):
    work = available[:]
    finish = [False if any(allocation[i][j] > 0 for j in range(num_resources))
              else True
              for i in range(num_processes)]

    changed = True
    while changed:
        changed = False
        for i in range(num_processes):
            if finish[i]:
                continue
            if all(request[i][j] <= work[j] for j in range(num_resources)):
                work = [work[j] + allocation[i][j] for j in range(num_resources)]
                finish[i] = True
                changed = True

    return [i for i in range(num_processes) if not finish[i]]
```

**Run:** `python3 bankers.py`

---

## 3. Go Source Code

```go
package main

import "fmt"

func bankersAlgorithm(nProc, nRes int, available []int, maxM, alloc [][]int) (bool, []int, [][]int) {
    need := make([][]int, nProc)
    for i := 0; i < nProc; i++ {
        need[i] = make([]int, nRes)
        for j := 0; j < nRes; j++ {
            need[i][j] = maxM[i][j] - alloc[i][j]
        }
    }
    work := make([]int, nRes)
    copy(work, available)
    finish := make([]bool, nProc)
    var seq []int

    count := 0
    for count < nProc {
        found := false
        for i := 0; i < nProc; i++ {
            if finish[i] { continue }
            ok := true
            for j := 0; j < nRes; j++ {
                if need[i][j] > work[j] { ok = false; break }
            }
            if ok {
                for j := 0; j < nRes; j++ { work[j] += alloc[i][j] }
                finish[i] = true
                seq = append(seq, i)
                count++; found = true
            }
        }
        if !found { break }
    }
    safe := true
    for _, f := range finish { if !f { safe = false; break } }
    return safe, seq, need
}

func deadlockDetection(nProc, nRes int, available []int, request, alloc [][]int) []int {
    work := make([]int, nRes)
    copy(work, available)
    finish := make([]bool, nProc)
    for i := 0; i < nProc; i++ {
        hasAlloc := false
        for j := 0; j < nRes; j++ { if alloc[i][j] > 0 { hasAlloc = true; break } }
        finish[i] = !hasAlloc
    }
    changed := true
    for changed {
        changed = false
        for i := 0; i < nProc; i++ {
            if finish[i] { continue }
            ok := true
            for j := 0; j < nRes; j++ { if request[i][j] > work[j] { ok = false; break } }
            if ok {
                for j := 0; j < nRes; j++ { work[j] += alloc[i][j] }
                finish[i] = true; changed = true
            }
        }
    }
    var dead []int
    for i := 0; i < nProc; i++ { if !finish[i] { dead = append(dead, i) } }
    return dead
}
```

**Run:** `go run bankers.go`

---

## 4. Test Case Results

All 7 test cases executed in both Python and Go — **identical results** in every case.

### TC-1: Classic Textbook Example (SAFE)
- **Setup:** 5 Processes, 3 Resources | Available: [3, 3, 2]
- **Banker's:** SAFE — Sequence: `P1 -> P3 -> P4 -> P0 -> P2`
- **Detection:** No Deadlock

### TC-2: Unsafe State — Zero Available
- **Setup:** 4 Processes, 3 Resources | Available: [0, 0, 0]
- **Banker's:** UNSAFE — No process can execute (empty sequence)
- **Detection:** DEADLOCK — All 4 processes deadlocked (P0, P1, P2, P3)

### TC-3: Small 3P/2R — Safe
- **Setup:** 3 Processes, 2 Resources | Available: [2, 1]
- **Banker's:** SAFE — Sequence: `P1 -> P2 -> P0`
- **Detection:** No Deadlock

### TC-4: Zero-Need Deadlock Breaker
- **Setup:** 5 Processes, 2 Resources | Available: [0, 0]
- P2 has Need=[0,0] so executes first, releases resources, unblocks others
- **Banker's:** SAFE — Sequence: `P2 -> P4 -> P0 -> P1 -> P3`
- **Detection:** No Deadlock

### TC-5: Critical Divergence — Banker UNSAFE but No Current Deadlock
- **Setup:** 6 Processes, 4 Resources | Available: [3, 1, 1, 2]
- P1's Need=[1,1,1,4] — cannot be satisfied in worst case (R3 deficit)
- P1's actual current Request=[1,0,0,0] — small and satisfiable
- **Banker's:** UNSAFE (P1 cannot get max resources)
- **Detection:** NO DEADLOCK (all current requests are satisfiable)

### TC-6: Complete Deadlock — All Processes Stuck
- **Setup:** 3 Processes, 2 Resources | Available: [0, 0]
- Circular wait with no slack — P0 waits on P1/P2, P1 waits on P0/P2
- **Banker's:** UNSAFE — empty sequence
- **Detection:** DEADLOCK — P0, P1, P2 all deadlocked

### TC-7: Edge Case — Single Process
- **Setup:** 1 Process, 1 Resource | Available: [5], Alloc: [0]
- **Banker's:** SAFE — Sequence: P0
- **Detection:** No Deadlock (zero allocation, immediately finished)

---

## 5. Critical Conclusions

### Conclusion 1: Banker's Algorithm is Conservative by Design
Banker's uses **worst-case future demand** (Need = Max - Allocation), not current requests. TC-5 proves this: P1's current request is small but its Need is large, so Banker's marks the state unsafe even though no deadlock currently exists. **The algorithm trades system concurrency for safety guarantees.**

### Conclusion 2: Detection Reacts; Banker's Prevents
- **Banker's** is **proactive** — denies requests that *might* lead to deadlock
- **Detection** is **reactive** — identifies deadlock *after* it has formed

Same underlying graph traversal logic, completely different input semantics (Need vs Request).

### Conclusion 3: A Zero-Need Process is a Deadlock Breaker (TC-4)
If any process Pi satisfies `Need[i] = [0,0,...,0]`, it can always execute regardless of available resources. It releases its allocations, potentially unblocking other processes. **In TC-4, P2 with Need=[0,0] saved the entire system from deadlock despite Available=[0,0].**

### Conclusion 4: Both Algorithms Are O(n^2 * m) — Polynomial
- Outer loop: at most n iterations
- Inner loop: scan n processes, check m resources

**Time Complexity:** O(n^2 * m) — efficient for real OS workloads.

### Conclusion 5: Unsafe Does NOT Equal Deadlocked (TC-5)
- **Unsafe:** "Cannot guarantee all processes finish if each requests its maximum"
- **Deadlocked:** "Right now, processes are stuck in circular wait"

```
Deadlocked states ⊂ Unsafe states ⊂ All states
```

A system can be unsafe but not deadlocked. A deadlocked system is always unsafe.

### Conclusion 6: Zero Available Resources Guarantees Deadlock Risk (TC-2, TC-6)
When the OS distributes all resource instances among blocked processes with no slack, deadlock is inevitable. **Operating systems should always keep a reserve of each resource — never allocate 100% of any resource class.**

### Conclusion 7: Safety Sequence Is Not Unique
Multiple valid safety sequences may exist for the same state. Both implementations find the first valid sequence via greedy linear scan. A production scheduler could optimize within the set of valid sequences for throughput, priority, or fairness.

### Conclusion 8: Python vs Go — Identical Correctness, Different Ergonomics

| Aspect | Python | Go |
|--------|--------|----|
| Code length | ~180 lines | ~230 lines |
| Type safety | Dynamic (runtime errors) | Static (compile-time safety) |
| Matrix init | Concise list comprehension | Verbose nested loops |
| Performance | Slower (interpreted) | ~5-10x faster (compiled) |
| Use case | Prototyping/education | OS-level/production systems |

For production OS components, **Go's static typing and compiled performance** are critical advantages.

### Conclusion 9: Detection's `finish[]` Initialization is Subtle
Processes with **zero allocation** are pre-marked `finish[i]=true` in Detection. They have no resources to release, so treating them as done prevents false deadlock reports. TC-5's P5 (zero allocation) is immediately finished — correctly so.

Banker's does not need this nuance because it considers all processes as active participants in the worst case.

### Conclusion 10: Combine Both Algorithms in Production
- Use **Banker's Algorithm** for request-time admission control (real-time/interactive systems)
- Use **Detection** periodically, especially when CPU utilization drops unexpectedly
- Finding a deadlock requires extra recovery steps: preemption, rollback, or process termination

**Combined strategy:** Banker's prevents; Detection is the safety net.

---

## 6. Summary Table

| Test Case | Processes | Resources | Available | Banker Safe? | Deadlock? | Key Pattern |
|-----------|-----------|-----------|-----------|:------------:|:---------:|-------------|
| TC-1: Classic | 5 | 3 | [3,3,2] | YES | No | Standard textbook example |
| TC-2: Zero Available | 4 | 3 | [0,0,0] | NO | YES — P0,P1,P2,P3 | Zero available = deadlock |
| TC-3: Small 3P/2R | 3 | 2 | [2,1] | YES | No | Minimal multi-resource case |
| TC-4: Zero-Need Breaker | 5 | 2 | [0,0] | YES | No | P2 (need=0) breaks the chain |
| TC-5: Critical Divergence | 6 | 4 | [3,1,1,2] | **NO** | **No** | **Unsafe != Deadlocked** |
| TC-6: Full Deadlock | 3 | 2 | [0,0] | NO | YES — P0,P1,P2 | Circular wait, no slack |
| TC-7: Edge 1P/1R | 1 | 1 | [5] | YES | No | Trivial single-process case |

---

> **Key Finding:** TC-5 is the most important result — Banker's Algorithm marked the state UNSAFE while the Deadlock Detection Algorithm found NO DEADLOCK. This proves that deadlock avoidance is strictly more conservative than detection, and that *unsafe != deadlocked*. Both Python and Go produced byte-for-byte identical logical results across all 7 test cases.
