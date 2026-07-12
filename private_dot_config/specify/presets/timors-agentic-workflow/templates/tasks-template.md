---

description: "Task list template for feature implementation"
---

# Tasks: [FEATURE NAME]

**Input**: Design documents from `/specs/[###-feature-name]/`

**Prerequisites**: plan.md (required), spec.md (required for user stories), research.md, data-model.md, contracts/

**Tests**: Generate tests by default for testable behavior. When new tests are
not practical, state either `Existing coverage: <concrete evidence expected>`
or `Not applicable: <concrete reason>` in that phase.

**Organization**: Tasks are grouped by user story for independent delivery.
Shared foundation and cross-cutting work have dedicated phases only when they
cannot be owned by one story.

## Format: `[ID] [P?] [Story?] Description with exact path`

- Every task is `- [ ] TNNN [P?] [Story?] <action with exact path>`.
- IDs start at `T001` and increase sequentially without gaps.
- `[P]` is task-local: the task touches different files and does not depend on
  an incomplete task. It does not declare execution-group concurrency.
- `[Story]` is required for story tasks and omitted from setup, foundation, and
  cross-cutting phases.
- Every task description names its exact project-relative file path.

<!--
  __SPECKIT_COMMAND_TASKS__ replaces all bracketed prompts with feature-specific
  content. Keep task descriptions, checks, goals, independent tests, and
  checkpoints here. Execution groups, global dependencies, models, workers,
  capabilities, contracts, and data flows belong only in execution-plan.md.
-->

## Phase 1: Setup

**Purpose**: [Feature-specific initialization purpose]

**Test Expectation**: [Tests required | Existing coverage: concrete evidence expected | Not applicable: concrete reason]

- [ ] T001 [Action] in [exact/project-relative/path]

**Checkpoint**: [Observable setup state]

---

## Phase 2: Foundation

**Purpose**: [Shared prerequisite that cannot be owned by one user story]

**Test Expectation**: [Tests required | Existing coverage: concrete evidence expected | Not applicable: concrete reason]

- [ ] T002 [Action] in [exact/project-relative/path]

**Checkpoint**: [Observable foundation state]

---

## Phase 3: User Story 1 - [Title] (Priority: P1)

**Goal**: [Behavior or value delivered by this story]

**Independent Test**: [Observable steps and result proving this story alone]

**Test Expectation**: [Tests required | Existing coverage: concrete evidence expected | Not applicable: concrete reason]

### Tests for User Story 1

- [ ] T003 [P] [US1] [Test behavior] in [exact/project-relative/test-path]

### Implementation for User Story 1

- [ ] T004 [US1] [Implement behavior] in [exact/project-relative/source-path]

**Checkpoint**: [Observable, independently testable story result]

---

<!-- Repeat one phase per remaining user story in priority order. Continue TNNN
without gaps and use the matching [USN] label. -->

## Final Phase: Cross-Cutting Concerns

**Purpose**: [Work that necessarily spans completed stories]

**Test Expectation**: [Tests required | Existing coverage: concrete evidence expected | Not applicable: concrete reason]

- [ ] TNNN [Action] in [exact/project-relative/path]

**Checkpoint**: [Observable final state and regression result]

## Task Checks

- Every task has a sequential `TNNN`, an exact path, and the applicable
  story label.
- Every `[P]` marker is justified by that task's files and prerequisites.
- Every story has a goal, independent test, and behavioral checkpoint.
- Every phase records an allowed concrete test expectation.
