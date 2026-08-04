#!/usr/bin/env python3
"""Deterministic, package-local checks for the canonical AGENTS adapter."""

import re
import shutil
import tempfile
from pathlib import Path

PACKAGE = Path(__file__).resolve().parent.parent
FIXTURES = Path(__file__).resolve().parent / "fixtures"
TEMPLATE = PACKAGE / "templates" / "constitution-template.md"
COMMAND = PACKAGE / "commands" / "speckit.constitution.md"
MANIFEST = PACKAGE / "preset.yml"
README = PACKAGE / "README.md"
LOCAL = Path(".specify/memory/constitution.local.md")
CONSTITUTION = Path(".specify/memory/constitution.md")


def check(condition, message):
    if not condition:
        raise AssertionError(message)


def source_state(project):
    agents = sorted(path for path in project.rglob("AGENTS.md") if path.is_file())
    local = project / LOCAL
    return agents, local if local.is_file() else None


def rule_lines(path):
    return [
        line.strip()[2:].strip()
        for line in path.read_text(encoding="utf-8").splitlines()
        if line.strip().startswith(("- ", "* "))
    ] + [
        line.strip()
        for line in path.read_text(encoding="utf-8").splitlines()
        if re.match(r"^[A-Za-z][A-Za-z ]+: .+", line.strip())
    ]


def resolve(project):
    agents, local = source_state(project)
    rules = []
    for path in agents:
        for rule in rule_lines(path):
            rules.append((path.relative_to(project).as_posix(), rule))
    if local:
        for rule in rule_lines(local):
            rules.append((LOCAL.as_posix(), rule))

    by_topic = {}
    effective = {}
    duplicates = []
    conflicts = []
    for source, rule in rules:
        topic, _, value = rule.partition(":")
        topic = topic.strip().lower()
        value = value.strip().lower()
        effective[topic] = (source, rule)
        previous = by_topic.get(topic)
        if previous and previous[1] == value:
            duplicates.append((previous[0], source, rule))
        elif previous and previous[1] != value:
            conflicts.append((previous[0], source, rule))
        else:
            by_topic[topic] = (source, value)

    missing = []
    if not agents:
        missing.append("AGENTS.md")
    if not local:
        missing.append(LOCAL.as_posix())
    return {
        "agents": [path.relative_to(project).as_posix() for path in agents],
        "local": local.relative_to(project).as_posix() if local else None,
        "rules": rules,
        "effective": effective,
        "duplicates": duplicates,
        "conflicts": conflicts,
        "missing": missing,
    }


def command_simulation(project, arguments):
    constitution = project / CONSTITUTION
    if not constitution.exists():
        constitution.parent.mkdir(parents=True, exist_ok=True)
        constitution.write_bytes(TEMPLATE.read_bytes())
    elif constitution.read_bytes() != TEMPLATE.read_bytes():
        return resolve(project), False

    result = resolve(project)
    explicit = bool(
        re.search(
            r"\b(spec kit|speckit)[ -](?:specific )?rule\b|"
            r"\b(?:add|change|update|remove)\b.*\bspec kit\b",
            arguments,
            re.IGNORECASE,
        )
    )
    wrote = False
    if explicit:
        local = project / LOCAL
        local.parent.mkdir(parents=True, exist_ok=True)
        local.write_text(
            "# Explicit Spec Kit-local request\n\n" + arguments + "\n",
            encoding="utf-8",
        )
        wrote = True
    return result, wrote


def plan_equivalent_consumption(result):
    """Model the downstream plan input without invoking a model or CLI."""
    return "\n".join(
        [
            "Effective governance sources:",
            *[source for source, _ in result["rules"]],
            "Effective governance rules:",
            *[rule for _, rule in result["effective"].values()],
            "Conflicts: " + str(len(result["conflicts"])),
        ]
    )


def test_manifest_and_text_contracts():
    manifest = MANIFEST.read_text(encoding="utf-8")
    check('schema_version: "1.0"' in manifest, "manifest schema is missing")
    check('id: "canonical-agents-adapter"' in manifest, "manifest id is missing")
    check('version: "0.1.0"' in manifest, "manifest version is missing")
    check(
        'file: "templates/constitution-template.md"' in manifest
        and 'replaces: "constitution-template"' in manifest,
        "constitution template mapping is missing",
    )
    check(
        'file: "commands/speckit.constitution.md"' in manifest
        and 'replaces: "speckit.constitution"' in manifest,
        "constitution command mapping is missing",
    )

    template = TEMPLATE.read_text(encoding="utf-8")
    command = COMMAND.read_text(encoding="utf-8")
    readme = README.read_text(encoding="utf-8")
    for fixture_rule in (
        "Follow repository formatting",
        "Run the repository checks",
        "Keep generated specifications reviewable",
        "Use pytest for project tests",
    ):
        check(fixture_rule not in template, "template copied a fixture project rule")
    for text, label in ((template, "template"), (command, "command")):
        check("AGENTS.md" in text, label + " does not name AGENTS governance")
        check("constitution.local.md" in text, label + " omits local governance")
        check("constitution.md" in text, label + " omits fixed infrastructure")
        check(
            "MUST NOT" in text or "never" in text.lower(),
            label + " lacks safety language",
        )
    check(
        re.match(r"\A---\ndescription: .*\n---\n", command),
        "frontmatter shape is invalid",
    )
    check(
        "copy" in command.lower() and "when it is missing" in command.lower(),
        "command does not initialize a missing fixed constitution",
    )
    check(
        "differs" in command.lower() and "do not overwrite" in command.lower(),
        "command does not report fixed-constitution drift",
    )
    check("specify preset add --dev" in readme, "README omits installation")
    check(
        "temporary" in readme and "model" in readme, "README omits harness limitation"
    )


def test_source_combinations_and_safety():
    expected = {
        "root-only": (1, False, True, 0),
        "nested-agents": (2, False, True, 0),
        "local-addition": (1, True, False, 0),
        "direct-conflict": (1, True, False, 1),
        "missing-agents": (0, True, True, 0),
        "missing-local": (1, False, True, 0),
        "both-missing": (0, False, True, 0),
    }
    template_before = TEMPLATE.read_bytes()
    for name, (agent_count, has_local, has_missing, conflict_count) in expected.items():
        with tempfile.TemporaryDirectory(prefix="canonical-agents-") as directory:
            project = Path(directory) / "project"
            shutil.copytree(FIXTURES / name, project)
            before = {path: path.read_bytes() for path in project.rglob("AGENTS.md")}
            result = resolve(project)
            check(len(result["agents"]) == agent_count, name + " AGENTS count mismatch")
            check(bool(result["local"]) == has_local, name + " local presence mismatch")
            check(
                bool(result["missing"]) == has_missing,
                name + " missing-source mismatch",
            )
            check(
                len(result["conflicts"]) == conflict_count,
                name + " conflict count mismatch",
            )
            if name == "direct-conflict":
                check(
                    any(
                        source == LOCAL.as_posix()
                        for source, _ in result["effective"].values()
                    ),
                    "local precedence did not win the Spec Kit conflict",
                )
            check(
                all(path.read_bytes() == contents for path, contents in before.items()),
                name + " changed AGENTS bytes during resolution",
            )
            consumed = plan_equivalent_consumption(result)
            check(
                "Effective governance sources:" in consumed,
                name + " plan input missing sources",
            )
            check(
                all(source in consumed for source in result["agents"]),
                name + " plan input omitted an AGENTS source",
            )
            if result["local"]:
                check(
                    result["local"] in consumed,
                    name + " plan input omitted local source",
                )
            if name == "both-missing":
                check(not result["rules"], "both-missing invented governance")
                check(
                    set(result["missing"]) == {"AGENTS.md", LOCAL.as_posix()},
                    "both-missing did not report both sources",
                )

            result, wrote = command_simulation(project, "")
            check(not wrote, name + " no-op created a local file")
            check(
                all(path.read_bytes() == contents for path, contents in before.items()),
                name + " no-op changed AGENTS bytes",
            )
            check(
                TEMPLATE.read_bytes() == template_before,
                "resolver changed during no-op",
            )

    with tempfile.TemporaryDirectory(prefix="canonical-agents-write-") as directory:
        project = Path(directory) / "project"
        shutil.copytree(FIXTURES / "root-only", project)
        agents_before = (project / "AGENTS.md").read_bytes()
        local = project / LOCAL
        check(not local.exists(), "write fixture unexpectedly had a local file")
        _, wrote = command_simulation(
            project, "Add a Spec Kit-specific rule: use the local review checklist."
        )
        check(
            wrote and local.exists(),
            "explicit Spec Kit request did not write local file",
        )
        check(
            (project / "AGENTS.md").read_bytes() == agents_before,
            "explicit write changed AGENTS",
        )
        check(
            TEMPLATE.read_bytes() == template_before, "explicit write changed resolver"
        )


def test_fixed_constitution_lifecycle():
    template_bytes = TEMPLATE.read_bytes()
    with tempfile.TemporaryDirectory(prefix="canonical-agents-constitution-") as directory:
        project = Path(directory) / "project"
        shutil.copytree(FIXTURES / "root-only", project)
        constitution = project / CONSTITUTION

        check(not constitution.exists(), "fixture unexpectedly has a constitution")
        command_simulation(project, "")
        check(
            constitution.read_bytes() == template_bytes,
            "missing constitution was not initialized exactly",
        )

        command_simulation(project, "")
        check(
            constitution.read_bytes() == template_bytes,
            "matching constitution changed on rerun",
        )

        drifted = b"# User-edited constitution\n"
        constitution.write_bytes(drifted)
        local = project / LOCAL
        _, wrote = command_simulation(
            project, "Add a Spec Kit-specific rule: require plan review."
        )
        check(
            constitution.read_bytes() == drifted,
            "drifted constitution was overwritten",
        )
        check(not wrote and not local.exists(), "drift did not block local governance writes")


def main():
    test_manifest_and_text_contracts()
    test_source_combinations_and_safety()
    test_fixed_constitution_lifecycle()
    print("canonical-agents-adapter: all deterministic checks passed")


if __name__ == "__main__":
    main()
