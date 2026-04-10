package main

import "fmt"

// =============================================================================
// Assignment 7 – OS: Banker's Algorithm + Deadlock Detection (Go)
// =============================================================================

// ─────────────────────────────────────────────────────────────────────────────
// BANKER'S ALGORITHM – Deadlock Avoidance
// ─────────────────────────────────────────────────────────────────────────────

type BankersResult struct {
	Safe     bool
	Sequence []int
	Steps    []string
	Need     [][]int
}

func bankersAlgorithm(nProc, nRes int, available []int, maxM, alloc [][]int) BankersResult {
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
	var steps []string

	steps = append(steps, fmt.Sprintf("  Available (Work): %v", work))
	steps = append(steps, "  Need Matrix:")
	for i := 0; i < nProc; i++ {
		steps = append(steps, fmt.Sprintf("    P%d: %v", i, need[i]))
	}

	count := 0
	for count < nProc {
		found := false
		for i := 0; i < nProc; i++ {
			if finish[i] {
				continue
			}
			canRun := true
			for j := 0; j < nRes; j++ {
				if need[i][j] > work[j] {
					canRun = false
					break
				}
			}
			if canRun {
				for j := 0; j < nRes; j++ {
					work[j] += alloc[i][j]
				}
				finish[i] = true
				seq = append(seq, i)
				steps = append(steps, fmt.Sprintf("  → P%d can execute. Work now: %v", i, work))
				count++
				found = true
			}
		}
		if !found {
			break
		}
	}

	safe := true
	for _, f := range finish {
		if !f {
			safe = false
			break
		}
	}
	return BankersResult{Safe: safe, Sequence: seq, Steps: steps, Need: need}
}

// ─────────────────────────────────────────────────────────────────────────────
// DEADLOCK DETECTION ALGORITHM
// ─────────────────────────────────────────────────────────────────────────────

type DetectionResult struct {
	Deadlocked []int
	Steps      []string
}

func deadlockDetection(nProc, nRes int, available []int, request, alloc [][]int) DetectionResult {
	work := make([]int, nRes)
	copy(work, available)

	finish := make([]bool, nProc)
	for i := 0; i < nProc; i++ {
		hasAlloc := false
		for j := 0; j < nRes; j++ {
			if alloc[i][j] > 0 {
				hasAlloc = true
				break
			}
		}
		finish[i] = !hasAlloc
	}

	var steps []string
	var initFinished []int
	for i, f := range finish {
		if f {
			initFinished = append(initFinished, i)
		}
	}
	steps = append(steps, fmt.Sprintf("  Work (Available): %v", work))
	steps = append(steps, fmt.Sprintf("  Initially-finished (zero-alloc): %v", initFinished))

	changed := true
	for changed {
		changed = false
		for i := 0; i < nProc; i++ {
			if finish[i] {
				continue
			}
			canFinish := true
			for j := 0; j < nRes; j++ {
				if request[i][j] > work[j] {
					canFinish = false
					break
				}
			}
			if canFinish {
				for j := 0; j < nRes; j++ {
					work[j] += alloc[i][j]
				}
				finish[i] = true
				steps = append(steps, fmt.Sprintf("  → P%d can complete. Work now: %v", i, work))
				changed = true
			}
		}
	}

	var deadlocked []int
	for i := 0; i < nProc; i++ {
		if !finish[i] {
			deadlocked = append(deadlocked, i)
		}
	}
	return DetectionResult{Deadlocked: deadlocked, Steps: steps}
}

// ─────────────────────────────────────────────────────────────────────────────
// HELPERS
// ─────────────────────────────────────────────────────────────────────────────

func printMatrix(label string, m [][]int, nProc, nRes int) {
	fmt.Printf("  %s:\n", label)
	hdr := "    "
	for j := 0; j < nRes; j++ {
		hdr += fmt.Sprintf("R%d ", j)
	}
	fmt.Println(hdr)
	for i, row := range m {
		line := fmt.Sprintf("  P%d  ", i)
		for _, v := range row {
			line += fmt.Sprintf("%2d ", v)
		}
		fmt.Println(line)
	}
}

func seqStr(seq []int) string {
	s := ""
	for idx, p := range seq {
		if idx > 0 {
			s += " → "
		}
		s += fmt.Sprintf("P%d", p)
	}
	return s
}

func runTestCase(tcID int, title string, nProc, nRes int,
	available []int, maxM, alloc, requestM [][]int) (bool, []int) {

	sep := "======================================================================"
	fmt.Printf("\n%s\n  TEST CASE %d: %s\n%s\n", sep, tcID, title, sep)

	// ── 1. Banker's Algorithm ──────────────────────────────────────────────
	fmt.Printf("\n[BANKER'S ALGORITHM – Deadlock Avoidance]\n")
	fmt.Printf("  Processes: %d  |  Resources: %d\n", nProc, nRes)
	fmt.Printf("  Available: %v\n", available)
	printMatrix("Max", maxM, nProc, nRes)
	printMatrix("Allocation", alloc, nProc, nRes)

	br := bankersAlgorithm(nProc, nRes, available, maxM, alloc)
	printMatrix("Need (Max - Alloc)", br.Need, nProc, nRes)
	fmt.Println("\n  Execution Trace:")
	for _, s := range br.Steps {
		fmt.Println(s)
	}
	if br.Safe {
		fmt.Printf("\n  ✅ SAFE STATE → Safety Sequence: %s\n", seqStr(br.Sequence))
	} else {
		fmt.Printf("\n  ❌ UNSAFE STATE — System may deadlock! (Partial seq: %v)\n", br.Sequence)
	}

	// ── 2. Deadlock Detection ──────────────────────────────────────────────
	fmt.Printf("\n[DEADLOCK DETECTION ALGORITHM]\n")
	printMatrix("Request (current)", requestM, nProc, nRes)

	dr := deadlockDetection(nProc, nRes, available, requestM, alloc)
	fmt.Println("\n  Execution Trace:")
	for _, s := range dr.Steps {
		fmt.Println(s)
	}
	if len(dr.Deadlocked) > 0 {
		pstr := ""
		for _, p := range dr.Deadlocked {
			pstr += fmt.Sprintf("P%d ", p)
		}
		fmt.Printf("\n  🔴 DEADLOCK DETECTED! Deadlocked Processes: [%s]\n", pstr)
	} else {
		fmt.Printf("\n  🟢 NO DEADLOCK — All processes can complete.\n")
	}
	fmt.Println()
	return br.Safe, dr.Deadlocked
}

// ─────────────────────────────────────────────────────────────────────────────
// MAIN – Test Cases
// ─────────────────────────────────────────────────────────────────────────────

func main() {
	fmt.Println("╔══════════════════════════════════════════════════════════════════╗")
	fmt.Println("║    OS Assignment 7 – Banker's Algo + Deadlock Detection (Go)     ║")
	fmt.Println("╚══════════════════════════════════════════════════════════════════╝")

	type summary struct {
		name string
		safe bool
		dead []int
	}
	var results []summary

	// TC-1: Classic textbook (safe)
	safe, dead := runTestCase(
		1, "Classic Textbook Example (SAFE)", 5, 3,
		[]int{3, 3, 2},
		[][]int{{7, 5, 3}, {3, 2, 2}, {9, 0, 2}, {2, 2, 2}, {4, 3, 3}},
		[][]int{{0, 1, 0}, {2, 0, 0}, {3, 0, 2}, {2, 1, 1}, {0, 0, 2}},
		[][]int{{0, 0, 0}, {2, 0, 2}, {0, 0, 0}, {1, 0, 0}, {0, 0, 2}},
	)
	results = append(results, summary{"TC-1: Classic (Safe)", safe, dead})

	// TC-2: Unsafe state
	safe, dead = runTestCase(
		2, "Unsafe State (Banker rejects)", 4, 3,
		[]int{0, 0, 0},
		[][]int{{2, 1, 2}, {1, 2, 1}, {2, 2, 1}, {1, 1, 2}},
		[][]int{{1, 0, 1}, {0, 1, 0}, {1, 1, 0}, {0, 0, 1}},
		[][]int{{1, 1, 1}, {1, 1, 1}, {1, 1, 1}, {1, 1, 1}},
	)
	results = append(results, summary{"TC-2: Unsafe (Deadlock)", safe, dead})

	// TC-3: Small system safe
	safe, dead = runTestCase(
		3, "Small System – 3 Processes, 2 Resources (SAFE)", 3, 2,
		[]int{2, 1},
		[][]int{{3, 2}, {2, 2}, {3, 3}},
		[][]int{{1, 0}, {1, 1}, {1, 1}},
		[][]int{{0, 0}, {1, 0}, {1, 1}},
	)
	results = append(results, summary{"TC-3: Small 3P/2R (Safe)", safe, dead})

	// TC-4: Partial deadlock
	safe, dead = runTestCase(
		4, "Partial Deadlock – 5 Processes, only some stuck", 5, 2,
		[]int{0, 0},
		[][]int{{2, 1}, {1, 2}, {1, 1}, {2, 2}, {1, 1}},
		[][]int{{1, 0}, {0, 1}, {1, 1}, {1, 0}, {0, 1}},
		[][]int{{1, 1}, {1, 1}, {0, 0}, {1, 2}, {1, 0}},
	)
	results = append(results, summary{"TC-4: Partial Deadlock", safe, dead})

	// TC-5: Large system 6P/4R
	safe, dead = runTestCase(
		5, "Large System – 6 Processes, 4 Resources (SAFE)", 6, 4,
		[]int{3, 1, 1, 2},
		[][]int{
			{3, 3, 2, 2}, {1, 2, 3, 4}, {1, 3, 5, 0},
			{2, 0, 1, 1}, {4, 2, 0, 0}, {3, 1, 0, 1},
		},
		[][]int{
			{1, 0, 0, 0}, {0, 1, 2, 0}, {1, 3, 5, 0},
			{0, 0, 1, 1}, {1, 0, 0, 0}, {0, 0, 0, 0},
		},
		[][]int{
			{0, 0, 0, 0}, {1, 0, 0, 0}, {0, 0, 0, 0},
			{2, 0, 0, 0}, {3, 2, 0, 0}, {3, 1, 0, 1},
		},
	)
	results = append(results, summary{"TC-5: Large 6P/4R (Safe)", safe, dead})

	// TC-6: Complete deadlock
	safe, dead = runTestCase(
		6, "Complete Deadlock – All Processes Stuck", 3, 2,
		[]int{0, 0},
		[][]int{{2, 2}, {2, 2}, {2, 2}},
		[][]int{{1, 0}, {0, 1}, {1, 1}},
		[][]int{{1, 2}, {2, 1}, {0, 1}},
	)
	results = append(results, summary{"TC-6: Full Deadlock", safe, dead})

	// TC-7: Edge single process
	safe, dead = runTestCase(
		7, "Edge Case – 1 Process, 1 Resource (trivially safe)", 1, 1,
		[]int{5},
		[][]int{{5}},
		[][]int{{0}},
		[][]int{{0}},
	)
	results = append(results, summary{"TC-7: Edge 1P/1R", safe, dead})

	// SUMMARY
	fmt.Println("\n======================================================================")
	fmt.Println("  SUMMARY OF ALL TEST CASES")
	fmt.Println("======================================================================")
	fmt.Printf("  %-35s %-10s %s\n", "Test Case", "Safe?", "Deadlocked Procs")
	fmt.Println("  " + "------------------------------------------------------------")
	for _, r := range results {
		dl := "None"
		if len(r.dead) > 0 {
			dl = ""
			for _, p := range r.dead {
				dl += fmt.Sprintf("P%d ", p)
			}
		}
		safeStr := "YES"
		if !r.safe {
			safeStr = "NO"
		}
		fmt.Printf("  %-35s %-10s %s\n", r.name, safeStr, dl)
	}
	fmt.Println()
}
