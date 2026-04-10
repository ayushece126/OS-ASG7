"""
Assignment 7 - OS: Banker's Algorithm (Deadlock Avoidance) + Deadlock Detection
Author: Implemented in Python
"""

# ─────────────────────────────────────────────────────────────────────────────
# BANKER'S ALGORITHM – Deadlock Avoidance
# Uses: Need = Max - Allocation; Safety Sequence check via work/finish
# ─────────────────────────────────────────────────────────────────────────────

def bankers_algorithm(num_processes, num_resources, available, max_matrix, allocation):
    """
    Banker's Algorithm for Deadlock Avoidance.
    Returns (is_safe, safety_sequence, steps)
    """
    need = [[max_matrix[i][j] - allocation[i][j]
             for j in range(num_resources)]
            for i in range(num_processes)]

    work    = available[:]
    finish  = [False] * num_processes
    safe_seq = []
    steps    = []

    steps.append(f"  Available (Work): {work}")
    steps.append(f"  Need Matrix:")
    for i in range(num_processes):
        steps.append(f"    P{i}: {need[i]}")

    count = 0
    while count < num_processes:
        found = False
        for i in range(num_processes):
            if finish[i]:
                continue
            # Check if need[i] <= work
            if all(need[i][j] <= work[j] for j in range(num_resources)):
                # Simulate execution: release resources
                work = [work[j] + allocation[i][j] for j in range(num_resources)]
                finish[i] = True
                safe_seq.append(i)
                steps.append(f"  → P{i} can execute. Work now: {work}")
                count += 1
                found = True
        if not found:
            break   # deadlock / unsafe

    is_safe = all(finish)
    return is_safe, safe_seq, steps, need


# ─────────────────────────────────────────────────────────────────────────────
# DEADLOCK DETECTION ALGORITHM
# Uses: Request matrix (current requests), Allocation, Available
# ─────────────────────────────────────────────────────────────────────────────

def deadlock_detection(num_processes, num_resources, available, request, allocation):
    """
    Deadlock Detection Algorithm.
    Returns (deadlocked_processes, steps)
    """
    work   = available[:]
    finish = [False if any(allocation[i][j] > 0 for j in range(num_resources))
              else True
              for i in range(num_processes)]
    steps  = [f"  Work (Available): {work}",
              f"  Initially-finished (zero-alloc): "
              f"{[i for i in range(num_processes) if finish[i]]}"]

    changed = True
    while changed:
        changed = False
        for i in range(num_processes):
            if finish[i]:
                continue
            if all(request[i][j] <= work[j] for j in range(num_resources)):
                work = [work[j] + allocation[i][j] for j in range(num_resources)]
                finish[i] = True
                steps.append(f"  → P{i} can complete. Work now: {work}")
                changed = True

    deadlocked = [i for i in range(num_processes) if not finish[i]]
    return deadlocked, steps


# ─────────────────────────────────────────────────────────────────────────────
# PRETTY PRINTER
# ─────────────────────────────────────────────────────────────────────────────

def print_matrix(label, matrix, n_proc, n_res):
    print(f"  {label}:")
    header = "    " + " ".join(f"R{j}" for j in range(n_res))
    print(header)
    for i, row in enumerate(matrix):
        print(f"  P{i}  " + " ".join(f"{v:2}" for v in row))


def run_test_case(tc_id, title, n_proc, n_res, available, max_m, alloc, request_m):
    sep = "=" * 70
    print(f"\n{sep}")
    print(f"  TEST CASE {tc_id}: {title}")
    print(sep)

    # ── 1. Banker's Algorithm ──────────────────────────────────────────────
    print("\n[BANKER'S ALGORITHM – Deadlock Avoidance]")
    print(f"  Processes: {n_proc}  |  Resources: {n_res}")
    print(f"  Available: {available}")
    print_matrix("Max", max_m, n_proc, n_res)
    print_matrix("Allocation", alloc, n_proc, n_res)

    safe, seq, b_steps, need = bankers_algorithm(n_proc, n_res, available, max_m, alloc)
    print_matrix("Need (Max - Alloc)", need, n_proc, n_res)
    print("\n  Execution Trace:")
    for s in b_steps:
        print(s)

    if safe:
        print(f"\n  ✅ SAFE STATE → Safety Sequence: {' → '.join('P'+str(p) for p in seq)}")
    else:
        print(f"\n  ❌ UNSAFE STATE — System may deadlock! (Partial seq: {seq})")

    # ── 2. Deadlock Detection ──────────────────────────────────────────────
    print("\n[DEADLOCK DETECTION ALGORITHM]")
    print_matrix("Request (current)", request_m, n_proc, n_res)

    deadlocked, d_steps = deadlock_detection(n_proc, n_res, available, request_m, alloc)
    print("\n  Execution Trace:")
    for s in d_steps:
        print(s)

    if deadlocked:
        print(f"\n  🔴 DEADLOCK DETECTED! Deadlocked Processes: "
              f"{['P'+str(p) for p in deadlocked]}")
    else:
        print(f"\n  🟢 NO DEADLOCK — All processes can complete.")

    print()
    return safe, deadlocked


# ─────────────────────────────────────────────────────────────────────────────
# TEST CASES
# ─────────────────────────────────────────────────────────────────────────────

if __name__ == "__main__":
    print("╔══════════════════════════════════════════════════════════════════╗")
    print("║   OS Assignment 7 – Banker's Algo + Deadlock Detection (Python)  ║")
    print("╚══════════════════════════════════════════════════════════════════╝")

    results = []

    # ── TC-1: Classic textbook (safe) ─────────────────────────────────────
    r = run_test_case(
        tc_id=1, title="Classic Textbook Example (SAFE)",
        n_proc=5, n_res=3,
        available=[3, 3, 2],
        max_m=[
            [7, 5, 3],
            [3, 2, 2],
            [9, 0, 2],
            [2, 2, 2],
            [4, 3, 3],
        ],
        alloc=[
            [0, 1, 0],
            [2, 0, 0],
            [3, 0, 2],
            [2, 1, 1],
            [0, 0, 2],
        ],
        request_m=[         # small requests → no deadlock
            [0, 0, 0],
            [2, 0, 2],
            [0, 0, 0],
            [1, 0, 0],
            [0, 0, 2],
        ]
    )
    results.append(("TC-1: Classic (Safe)", r))

    # ── TC-2: Unsafe state ────────────────────────────────────────────────
    r = run_test_case(
        tc_id=2, title="Unsafe State (Banker rejects)",
        n_proc=4, n_res=3,
        available=[0, 0, 0],   # no resources left
        max_m=[
            [2, 1, 2],
            [1, 2, 1],
            [2, 2, 1],
            [1, 1, 2],
        ],
        alloc=[
            [1, 0, 1],
            [0, 1, 0],
            [1, 1, 0],
            [0, 0, 1],
        ],
        request_m=[            # every process is waiting
            [1, 1, 1],
            [1, 1, 1],
            [1, 1, 1],
            [1, 1, 1],
        ]
    )
    results.append(("TC-2: Unsafe (Deadlock)", r))

    # ── TC-3: Small system 3P/2R safe ───────────────────────────────────
    r = run_test_case(
        tc_id=3, title="Small System – 3 Processes, 2 Resources (SAFE)",
        n_proc=3, n_res=2,
        available=[2, 1],
        max_m=[
            [3, 2],
            [2, 2],
            [3, 3],
        ],
        alloc=[
            [1, 0],
            [1, 1],
            [1, 1],
        ],
        request_m=[
            [0, 0],
            [1, 0],
            [1, 1],
        ]
    )
    results.append(("TC-3: Small 3P/2R (Safe)", r))

    # ── TC-4: Partial deadlock (some processes deadlocked) ───────────────
    r = run_test_case(
        tc_id=4, title="Partial Deadlock – 5 Processes, only 2 stuck",
        n_proc=5, n_res=2,
        available=[0, 0],
        max_m=[
            [2, 1],
            [1, 2],
            [1, 1],
            [2, 2],
            [1, 1],
        ],
        alloc=[
            [1, 0],
            [0, 1],
            [1, 1],   # P2 has all it needs
            [1, 0],
            [0, 1],
        ],
        request_m=[
            [1, 1],   # P0 waiting
            [1, 1],   # P1 waiting
            [0, 0],   # P2 done
            [1, 2],   # P3 waiting – never satisfied
            [1, 0],   # P4 waiting – never satisfied
        ]
    )
    results.append(("TC-4: Partial Deadlock", r))

    # ── TC-5: Large system 6P/4R safe ───────────────────────────────────
    r = run_test_case(
        tc_id=5, title="Large System – 6 Processes, 4 Resources (SAFE)",
        n_proc=6, n_res=4,
        available=[3, 1, 1, 2],
        max_m=[
            [3, 3, 2, 2],
            [1, 2, 3, 4],
            [1, 3, 5, 0],
            [2, 0, 1, 1],
            [4, 2, 0, 0],
            [3, 1, 0, 1],
        ],
        alloc=[
            [1, 0, 0, 0],
            [0, 1, 2, 0],
            [1, 3, 5, 0],
            [0, 0, 1, 1],
            [1, 0, 0, 0],
            [0, 0, 0, 0],
        ],
        request_m=[
            [0, 0, 0, 0],
            [1, 0, 0, 0],
            [0, 0, 0, 0],
            [2, 0, 0, 0],
            [3, 2, 0, 0],
            [3, 1, 0, 1],
        ]
    )
    results.append(("TC-5: Large 6P/4R (Safe)", r))

    # ── TC-6: Complete deadlock (all processes stuck) ────────────────────
    r = run_test_case(
        tc_id=6, title="Complete Deadlock – All Processes Stuck",
        n_proc=3, n_res=2,
        available=[0, 0],
        max_m=[
            [2, 2],
            [2, 2],
            [2, 2],
        ],
        alloc=[
            [1, 0],
            [0, 1],
            [1, 1],
        ],
        request_m=[
            [1, 2],  # waiting for R1, but R1 held by P1
            [2, 1],  # waiting for R0, but R0 held by P0
            [0, 1],  # waiting for R1, but R1 held by P1
        ]
    )
    results.append(("TC-6: Full Deadlock", r))

    # ── TC-7: Single process single resource ─────────────────────────────
    r = run_test_case(
        tc_id=7, title="Edge Case – 1 Process, 1 Resource (trivially safe)",
        n_proc=1, n_res=1,
        available=[5],
        max_m=[[5]],
        alloc=[[0]],
        request_m=[[0]]
    )
    results.append(("TC-7: Edge 1P/1R", r))

    # ── SUMMARY ───────────────────────────────────────────────────────────
    print("\n" + "=" * 70)
    print("  SUMMARY OF ALL TEST CASES")
    print("=" * 70)
    print(f"  {'Test Case':<35} {'Safe?':<10} {'Deadlocked Procs'}")
    print("  " + "-" * 60)
    for name, (safe, dead) in results:
        dl = ", ".join("P"+str(p) for p in dead) if dead else "None"
        print(f"  {name:<35} {'YES' if safe else 'NO':<10} {dl}")
    print()
