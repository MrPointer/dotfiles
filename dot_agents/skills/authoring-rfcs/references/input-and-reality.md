# Input And Reality

Use this reference when gathering RFC inputs, deciding whether to load architecture-pattern guidance, verifying current reality, or normalizing settled decisions.

## Inputs To Read

Read available design context before drafting:

- user-approved direction and constraints from the current conversation
- active anchor, if one exists
- existing RFCs, design documents, ADRs, architecture docs, domain docs, and process docs relevant to the affected area
- source files that define current architecture, contracts, schemas, configuration, storage formats, or runtime flow

If the work is greenfield or the codebase has no relevant implementation yet, state that explicitly in the RFC and separate assumptions from verified facts.

## Architecture Pattern Companion

Load `applying-architecture-patterns` before drafting when the RFC involves any of these:

- new backend system architecture
- refactoring monoliths or tightly coupled backend boundaries
- Clean Architecture, Hexagonal Architecture, ports/adapters, DDD, aggregates, repositories, domain events, bounded contexts, or microservice decomposition
- testability, mockability, or framework-independence as a design goal

Use the companion skill as design vocabulary, tradeoff guidance, and a pitfall checklist. Do not copy its content into the RFC. Do not apply a pattern unless it is justified by the current system and goals.

## Current-State Model

Build a concise current-state model before writing proposed design:

- existing components and their responsibilities
- current data/control flow through the affected path
- existing contracts, public interfaces, configuration surfaces, storage formats, or external integrations
- constraints created by shipped behavior, persisted data, compatibility requirements, operational practices, or team conventions
- current pain points that the design intentionally addresses

Every important `currently` claim in the RFC should be backed by a source reference, documentation reference, or explicit user statement.

## Normalize Decisions

Before drafting, reduce brainstorming or discussion output into settled design inputs:

- chosen approach and why it was chosen
- goals and non-goals
- hard constraints and compatibility requirements
- rejected alternatives with concrete reasons
- open questions that remain genuinely unresolved

If an unresolved question changes the architecture, stop and ask the user one focused question. Non-blocking questions can remain in the RFC as open questions.
