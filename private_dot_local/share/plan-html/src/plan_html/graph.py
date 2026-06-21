"""Strict Mermaid-style dependency graph parsing for plan-html."""

import re
from dataclasses import dataclass

from plan_html.markdown import section_body


@dataclass(frozen=True)
class GraphNode:
    id: str
    label: str
    subplan: str


@dataclass(frozen=True)
class GraphEdge:
    source: str
    target: str


@dataclass(frozen=True)
class PlanGraph:
    direction: str
    nodes: list[GraphNode]
    edges: list[GraphEdge]
    warnings: list[str]


MERMAID_BLOCK_RE = re.compile(r"```mermaid\s*\n(.*?)\n```", re.S | re.I)
FLOWCHART_RE = re.compile(r"^flowchart\s+(LR|TD)\s*$", re.I)
NODE_RE = re.compile(r'^\s*(SP(?P<num>\d{2}))\s*\["(?P<label>[^"]+)"\]\s*$')
EDGE_RE = re.compile(r"^\s*(SP\d{2})\s*-->\s*(SP\d{2})\s*$")


def parse_dependency_graph(md: str, known_subplans: set[str] | None = None):
    body = section_body(md, "Dependency Graph")
    if not body:
        return None
    match = MERMAID_BLOCK_RE.search(body)
    if not match:
        return PlanGraph("LR", [], [], ["Dependency Graph section has no mermaid block."])

    known_subplans = known_subplans or set()
    direction = "LR"
    nodes: dict[str, GraphNode] = {}
    edges: list[GraphEdge] = []
    warnings: list[str] = []

    for raw_line in match.group(1).splitlines():
        line = raw_line.strip()
        if not line:
            continue

        flow = FLOWCHART_RE.match(line)
        if flow:
            direction = flow.group(1).upper()
            continue

        node = NODE_RE.match(line)
        if node:
            node_id = node.group(1)
            subplan = node.group("num")
            if known_subplans and subplan not in known_subplans:
                warnings.append(f"Graph node {node_id} has no matching sub-plan file.")
            nodes[node_id] = GraphNode(node_id, node.group("label"), subplan)
            continue

        edge = EDGE_RE.match(line)
        if edge:
            source, target = edge.groups()
            edges.append(GraphEdge(source, target))
            for node_id in (source, target):
                if node_id not in nodes:
                    subplan = node_id.removeprefix("SP")
                    if known_subplans and subplan not in known_subplans:
                        warnings.append(
                            f"Graph edge references {node_id}, which has no matching sub-plan file."
                        )
                    nodes.setdefault(node_id, GraphNode(node_id, node_id, subplan))
            continue

        warnings.append(f"Unsupported graph line: {line}")

    return PlanGraph(direction, list(nodes.values()), edges, warnings)


def graph_to_dict(graph: PlanGraph):
    return {
        "direction": graph.direction,
        "nodes": [node.__dict__ for node in graph.nodes],
        "edges": [edge.__dict__ for edge in graph.edges],
        "warnings": graph.warnings,
    }
