#!/usr/bin/env python3
"""Deterministic contract checks for optional human review projections."""

from pathlib import Path


PACKAGE = Path(__file__).resolve().parent.parent
MANIFEST = PACKAGE / "preset.yml"
COMPATIBILITY = PACKAGE / "references" / "protocol-compatibility.md"
REVIEW = PACKAGE / "references" / "human-review.md"
EXECUTION_PLAN_TEMPLATE = PACKAGE / "templates" / "execution-plan-template.md"
SPECIFY = PACKAGE / "commands" / "speckit.specify.md"
PLAN = PACKAGE / "commands" / "speckit.plan.md"
ANALYZE = PACKAGE / "commands" / "speckit.analyze.md"
README = PACKAGE / "README.md"


def check(condition, message):
    if not condition:
        raise AssertionError(message)


def text(path):
    return path.read_text(encoding="utf-8")


def test_manifest_inventory_and_command_mappings():
    manifest = text(MANIFEST)
    compatibility = text(COMPATIBILITY)

    check(REVIEW.is_file(), "shared human-review reference is missing")
    check(SPECIFY.is_file(), "wrapped specify command is missing")
    check(
        'name: "speckit.specify"' in manifest
        and 'file: "commands/speckit.specify.md"' in manifest
        and 'strategy: "wrap"' in manifest,
        "manifest does not wrap speckit.specify",
    )
    check("optional human review" in manifest, "preset metadata omits human review")
    check(
        "commands/speckit.specify.md" in compatibility
        and "references/human-review.md" in compatibility,
        "package inventory omits human-review surfaces",
    )
    check(
        "speckit.specify" in compatibility and "human-review.md" in compatibility,
        "compatibility mapping rules omit the specify review projection",
    )


def test_preset_version_boundary_keeps_protocol_stable():
    manifest = text(MANIFEST)
    execution_plan_template = text(EXECUTION_PLAN_TEMPLATE)

    check('version: "0.2.1"' in manifest, "preset version is not 0.2.1")
    check(
        'protocol_version: "0.1.0"' in manifest,
        "execution-plan protocol version is not 0.1.0",
    )
    check(
        "**Preset Protocol Version**: 0.1.0" in execution_plan_template,
        "execution-plan template protocol version is not 0.1.0",
    )


def test_execution_model_repeats_the_group_tier_and_rationale():
    execution_plan_template = text(EXECUTION_PLAN_TEMPLATE)
    expected_model = (
        "#### Execution Model\n\n"
        "**Model Tier**: Mid-tier\n"
        "**Rationale**: [Provider-neutral rationale for the semantic tier selected "
        "in Execution Groups.]"
    )

    check(
        execution_plan_template.count(expected_model) == 2,
        "execution-model sections must contain matching model-tier and rationale lines",
    )


def test_completion_projections_preserve_upstream_and_exclude_analyze():
    reference = text(REVIEW)
    normalized_reference = " ".join(reference.lower().split())
    specify = text(SPECIFY)
    plan = text(PLAN)
    analyze = text(ANALYZE)

    for command, name in ((specify, "specify"), (plan, "plan")):
        deferral = "## Completion Report And Handoff Deferral"
        review_heading = "## Optional " + (
            "Specification" if name == "specify" else "Planning"
        ) + " Human Review"
        check("## Package Integrity Gate" in command, name + " omits package gate")
        check("references/human-review.md" in command, name + " omits shared review")
        check("{CORE_TEMPLATE}" in command, name + " does not retain upstream composition")
        check(deferral in command, name + " omits final-report deferral")
        pre_core = command[: command.index("{CORE_TEMPLATE}")]
        normalized_pre_core = " ".join(pre_core.split())
        for phrase in (
            "Execute the complete upstream workflow",
            "artifact generation",
            "validation",
            "mandatory post-execution hooks",
            "Do not emit the upstream final Completion Report",
            "handoffs",
            "optional human-review step below has completed or is declined",
        ):
            check(
                " ".join(phrase.split()) in normalized_pre_core,
                name + " deferral omits: " + phrase,
            )
        check(
            command.index(deferral) < command.index("{CORE_TEMPLATE}") < command.index(review_heading),
            name + " does not defer reporting before core while retaining post-core review",
        )
        check(
            "after" in command.lower() and "validated" in command.lower(),
            name + " does not place review after validated artifact generation",
        )
        check("scripts:" not in command, name + " overrides inherited scripts")
        check("agent_scripts:" not in command, name + " overrides inherited agent scripts")
    check("references/human-review.md" not in analyze, "analyze advertises human review")

    for phrase in (
        "guided walkthrough",
        "independent review",
        "optional",
        "Pause after each review unit",
        "immediately update every",
        "affected artifact owned by the current phase",
        "no persisted preference",
        "no approval record",
        "no workflow gate",
        "Do not force a literal section-by-section or file-by-file walkthrough",
    ):
        check(
            " ".join(phrase.lower().split()) in normalized_reference,
            "review protocol omits: " + phrase,
        )


def test_phase_guidance_and_docs_are_present():
    specify = " ".join(text(SPECIFY).lower().split())
    plan = " ".join(text(PLAN).lower().split())
    readme = text(README).lower()

    for phrase in (
        "problem/outcome/scope",
        "scenarios/priorities",
        "acceptance behavior/edge cases",
        "measurable success criteria",
    ):
        check(" ".join(phrase.split()) in specify, "specify guidance omits: " + phrase)
    for phrase in (
        "technical approach/research decisions",
        "interfaces/integration/data/contracts",
        "testing/validation",
        "documentation impact/unresolved design choices",
    ):
        check(" ".join(phrase.split()) in plan, "plan guidance omits: " + phrase)
    check("human review" in readme, "README omits human-review behavior")


def test_current_phase_ownership_prevents_specification_reconfirmation():
    reference = " ".join(text(REVIEW).lower().split())
    plan = " ".join(text(PLAN).lower().split())

    for phrase in (
        "select only consequential decisions owned or introduced by the current phase",
        "authoritative baselines for the current phase",
        "must not become review units requiring renewed confirmation",
    ):
        check(
            " ".join(phrase.split()) in reference,
            "shared review ownership rule omits: " + phrase,
        )

    for phrase in (
        "treat the current `spec.md` requirements, scenarios, boundaries, and success criteria as authoritative inputs",
        "review only technical/design decisions introduced or derived during planning",
        "does not introduce a meaningful planning choice must not be presented merely to reconfirm the specification",
        "correct the plan rather than reopen the specification",
        "route it to the stock owning workflow; do not relitigate it inside planning review",
    ):
        check(
            " ".join(phrase.split()) in plan,
            "plan ownership guidance omits: " + phrase,
        )


def main():
    test_manifest_inventory_and_command_mappings()
    test_preset_version_boundary_keeps_protocol_stable()
    test_execution_model_repeats_the_group_tier_and_rationale()
    test_completion_projections_preserve_upstream_and_exclude_analyze()
    test_phase_guidance_and_docs_are_present()
    test_current_phase_ownership_prevents_specification_reconfirmation()
    print("timors-agentic-workflow human-review contract: all checks passed")


if __name__ == "__main__":
    main()
