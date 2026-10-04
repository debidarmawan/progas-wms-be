# Test Requirements — progas-wms-be

Standard Go `testing` package only (no testify/assert libraries are in `go.mod`).
Tests are table-driven with `t.Run` subtests.

## 1. Commands
- All tests: `go test ./...`
- Race detector: `go test -race ./...`
- Coverage: `go test -cover ./...` (targeted: `go test -cover ./helper/...`)
- Every change must keep `go build ./... && go vet ./...` green.

## 2. Placement & naming
- Test files live beside the code as `<file>_test.go` in the SAME package (white-box), e.g. `helper/finance_test.go` (`package helper`).
- Test funcs: `TestXxx`; subtests: short descriptive names (`"overpayment still counts as paid"`).
- Use the standard library only; do not add assertion frameworks.

## 3. What MUST be covered
Business logic and pure functions are the priority — they are deterministic and DB-free:
- `helper/` domain functions: weight/pricing (`CylinderFilledWeightKg`, `SumCylinderWeight`), status guards (`CanIssueOnDO`, `CanDeliverOnExchange`, `CanReceiveOnExchange`, `ValidateDOCylinders`, `ValidateExchangeOutCylinders`, `ValidateExchangeInCylinders`, `ValidateBarcodeListsUnique`), quota (`WouldExceedQuota`), invoice status (`DetermineInvoiceStatus`), aging (`DaysAtCustomer`, `IsOverdueAtCustomer`), pagination (`NormalizePagination`, `BuildPaginationMeta`), search (`NormalizeSearch`, `HasSearch`, `SearchPattern`).
- Usecase decision logic that has been extracted into pure helpers/validators.
- Enum behaviour: `IsValid()` accepts every valid value and rejects unknown ones.
- Validation tags behave (`helper.ValidateStruct` over a sample DTO).

## 4. Required tests per change type
- New business rule → add a table-driven test with happy path + boundary + failure cases.
- Bug fix → add a regression test that reproduces the bug first.
- New status/state → test every allowed AND disallowed transition.
- New helper function → must ship with a test.
- New endpoint → at minimum cover the underlying usecase/helper logic at unit level.

## 5. Table-driven template
```go
func TestSomething(t *testing.T) {
	cases := []struct {
		name string
		in   int
		want int
	}{
		{"zero", 0, 0},
		{"positive", 3, 3},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := Something(c.in); got != c.want {
				t.Errorf("Something(%v) = %v, want %v", c.in, got, c.want)
			}
		})
	}
}
```
This mirrors `helper/finance_test.go`.

## 6. Testability guidance
- Keep DB/GORM out of unit tests. If logic needs the DB, extract the pure decision into `helper/` and test that.
- Repository methods depend on `*gorm.DB`; prefer testing the orchestration/guards rather than GORM itself.
- Usecases depend on repository INTERFACES, so hand-written fakes (no mock framework) are acceptable for orchestration tests.
- Cover error paths: invalid input, not-found, over-quota, status mismatch, empty slices.

## 7. Expectations & guardrails
- Do not regress existing coverage; add tests when you touch a covered file.
- Do not assert on formatting of internals, timestamps, or UUID values.
- Keep tests deterministic (no bare `time.Now()` without a seam; no network).
- CI currently has no test gate — if you add one, use `go vet ./... && go test ./...`.

## 8. Definition of done (testing)
A task is not complete until: existing tests pass, new behaviour has tests, and `go test ./...` is green locally.
