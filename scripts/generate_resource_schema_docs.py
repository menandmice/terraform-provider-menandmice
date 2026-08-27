#!/usr/bin/env python3
"""
Extracts the Terraform schema (attributes, types, required/optional/computed,
defaults, descriptions) directly from the Go source of this provider's
resources, and writes one Markdown file per resource describing exactly what
a customer needs to know to write/import that resource.

This does not use the Terraform example .tf files or the tfplugindocs output
in docs/ - it parses the schema.Schema map literals in menandmice/*.go so the
generated reference always matches the actual provider code.

Usage:
    python scripts/generate_resource_schema_docs.py

Output:
    generated/resource-schemas/<resource_type>.md   (gitignored)
"""

import argparse
import os
import re
import sys

REPO_ROOT = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
SRC_DIR = os.path.join(REPO_ROOT, "menandmice")
OUT_DIR = os.path.join(REPO_ROOT, "generated", "resource-schemas")

TYPE_NAMES = {
    "TypeString": "String",
    "TypeBool": "Boolean",
    "TypeInt": "Number",
    "TypeFloat": "Float",
    "TypeList": "List",
    "TypeSet": "Set",
    "TypeMap": "Map",
}


def load_go_sources():
    """Return {filename: text} for every non-test .go file in menandmice/."""
    sources = {}
    for name in os.listdir(SRC_DIR):
        if name.endswith(".go") and not name.endswith("_test.go"):
            path = os.path.join(SRC_DIR, name)
            with open(path, "r", encoding="utf-8") as f:
                sources[name] = f.read()
    return sources


def find_brace_matches(text):
    """
    Scan Go source and return {index_of_'{': index_of_matching_'}'} for every
    *structural* brace, i.e. one that isn't inside a string, rune literal, or
    comment. Braces that appear inside comments/strings are ignored entirely,
    which is what lets us skip over commented-out schema attributes below.
    """
    matches = {}
    stack = []
    state = "code"
    i, n = 0, len(text)
    while i < n:
        c = text[i]
        if state == "code":
            if c == "/" and i + 1 < n and text[i + 1] == "/":
                state = "line_comment"
                i += 2
                continue
            if c == "/" and i + 1 < n and text[i + 1] == "*":
                state = "block_comment"
                i += 2
                continue
            if c == '"':
                state = "dstring"
            elif c == "`":
                state = "backtick"
            elif c == "'":
                state = "rune"
            elif c == "{":
                stack.append(i)
            elif c == "}":
                if stack:
                    matches[stack.pop()] = i
            i += 1
        elif state == "line_comment":
            if c == "\n":
                state = "code"
            i += 1
        elif state == "block_comment":
            if c == "*" and i + 1 < n and text[i + 1] == "/":
                state = "code"
                i += 2
                continue
            i += 1
        elif state == "dstring":
            if c == "\\":
                i += 2
                continue
            if c == '"':
                state = "code"
            i += 1
        elif state == "backtick":
            if c == "`":
                state = "code"
            i += 1
        elif state == "rune":
            if c == "\\":
                i += 2
                continue
            if c == "'":
                state = "code"
            i += 1
    return matches


def build_masked(text):
    """Return a same-length copy of text with the contents of every string,
    rune literal, and comment blanked out (kept as spaces/newlines). Used to
    locate schema field keywords (Default:, Required:, ...) while ignoring
    look-alike text that happens to sit inside a Description string, e.g.
    `Description: "... (Default: True)"`."""
    out = []
    state = "code"
    i, n = 0, len(text)
    while i < n:
        c = text[i]
        if state == "code":
            if c == "/" and i + 1 < n and text[i + 1] == "/":
                state = "line_comment"
                out.append("  ")
                i += 2
                continue
            if c == "/" and i + 1 < n and text[i + 1] == "*":
                state = "block_comment"
                out.append("  ")
                i += 2
                continue
            if c == '"':
                state = "dstring"
                out.append(" ")
            elif c == "`":
                state = "backtick"
                out.append(" ")
            elif c == "'":
                state = "rune"
                out.append(" ")
            else:
                out.append(c)
            i += 1
        elif state == "line_comment":
            if c == "\n":
                state = "code"
                out.append(c)
            else:
                out.append(" ")
            i += 1
        elif state == "block_comment":
            if c == "*" and i + 1 < n and text[i + 1] == "/":
                state = "code"
                out.append("  ")
                i += 2
                continue
            out.append(" " if c != "\n" else c)
            i += 1
        elif state == "dstring":
            if c == "\\":
                out.append("  ")
                i += 2
                continue
            if c == '"':
                state = "code"
                out.append(" ")
            else:
                out.append(" " if c != "\n" else c)
            i += 1
        elif state == "backtick":
            if c == "`":
                state = "code"
                out.append(" ")
            else:
                out.append(" " if c != "\n" else c)
            i += 1
        elif state == "rune":
            if c == "\\":
                out.append("  ")
                i += 2
                continue
            if c == "'":
                state = "code"
                out.append(" ")
            else:
                out.append(" " if c != "\n" else c)
            i += 1
    return "".join(out)


ENTRY_KEY_RE = re.compile(r'"([A-Za-z0-9_]+)"\s*:\s*(?:&schema\.Schema)?\s*\{')


def iter_top_level_entries(text, matches, body_start, body_end):
    """Yield (name, entry_start, entry_end) for each top-level attribute in a
    `map[string]*schema.Schema{ ... }` body, skipping anything commented out."""
    cursor = body_start
    while cursor < body_end:
        m = ENTRY_KEY_RE.search(text, cursor, body_end)
        if not m:
            break
        brace_pos = m.end() - 1
        if brace_pos not in matches:
            # brace lived inside a comment/string - not a real entry, keep scanning
            cursor = m.end()
            continue
        close = matches[brace_pos]
        yield m.group(1), brace_pos + 1, close
        cursor = close + 1


def locate_field(masked, field_name):
    """Find where a field keyword lives in code (not inside a string/comment)
    and return its start offset, or None."""
    m = re.search(r"(?<![A-Za-z0-9_])" + field_name + r"\s*:", masked)
    return m.start() if m else None


def extract_string_literal(original, masked, field_name):
    pos = locate_field(masked, field_name)
    if pos is None:
        return None
    rest = original[pos:]
    m = re.match(field_name + r'\s*:\s*"((?:[^"\\]|\\.)*)"', rest)
    if m:
        return m.group(1).replace('\\"', '"')
    m = re.match(field_name + r"\s*:\s*`([^`]*)`", rest)
    if m:
        return " ".join(m.group(1).split())
    return None


def extract_bool(original, masked, field_name):
    pos = locate_field(masked, field_name)
    if pos is None:
        return False
    return re.match(field_name + r"\s*:\s*true\b", original[pos:]) is not None


def extract_default(original, masked):
    pos = locate_field(masked, "Default")
    if pos is None:
        return None
    m = re.match(r"Default\s*:\s*([^,\n]+)", original[pos:])
    if not m:
        return None
    return m.group(1).strip().rstrip(",")


def extract_int(original, masked, field_name):
    pos = locate_field(masked, field_name)
    if pos is None:
        return None
    m = re.match(field_name + r"\s*:\s*(\d+)", original[pos:])
    return m.group(1) if m else None


def extract_string_list(original, masked, field_name):
    pos = locate_field(masked, field_name)
    if pos is None:
        return []
    m = re.match(field_name + r"\s*:\s*\[\]string\{([^}]*)\}", original[pos:])
    if not m:
        return []
    return re.findall(r'"([^"]*)"', m.group(1))


def parse_attr(full_text, masked_text, matches, name, start, end):
    body = full_text[start:end]
    masked_body = masked_text[start:end]

    attr = {
        "name": name,
        "type": None,
        "element_type": None,
        "children": None,
    }

    # An attribute's own fields (Required, Default, ConflictsWith, ...) can
    # share a keyword with fields nested inside its Elem sub-schema (e.g. a
    # nested "mask" attribute also has ConflictsWith). Strip that nested span
    # out before scanning for the parent's own fields so nothing bleeds up.
    own_body, own_masked = body, masked_body
    elem_match = re.search(r"Elem:\s*&schema\.(Resource|Schema)\s*\{", masked_body)
    if elem_match:
        elem_kind = elem_match.group(1)
        brace_pos = start + elem_match.end() - 1
        if brace_pos in matches:
            elem_close = matches[brace_pos]
            cut = elem_close + 1 - start
            own_body = body[: elem_match.start()] + body[cut:]
            own_masked = masked_body[: elem_match.start()] + masked_body[cut:]

            if elem_kind == "Schema":
                elem_body = full_text[brace_pos + 1: elem_close]
                elem_masked_body = masked_text[brace_pos + 1: elem_close]
                elem_type_pos = locate_field(elem_masked_body, "Type")
                if elem_type_pos is not None:
                    tm = re.match(r"Type\s*:\s*schema\.(\w+)", elem_body[elem_type_pos:])
                    if tm:
                        attr["element_type"] = tm.group(1)
            else:  # Resource -> nested attribute set
                nested_map_match = re.search(
                    r"Schema:\s*map\[string\]\*schema\.Schema\s*\{",
                    masked_text[brace_pos:elem_close],
                )
                if nested_map_match:
                    nmap_brace_pos = brace_pos + nested_map_match.end() - 1
                    if nmap_brace_pos in matches:
                        nmap_close = matches[nmap_brace_pos]
                        attr["children"] = [
                            parse_attr(full_text, masked_text, matches, cname, cstart, cend)
                            for cname, cstart, cend in iter_top_level_entries(
                                full_text, matches, nmap_brace_pos + 1, nmap_close
                            )
                        ]

    attr.update(
        {
            "description": extract_string_literal(own_body, own_masked, "Description") or "",
            "deprecated": extract_string_literal(own_body, own_masked, "Deprecated"),
            "required": extract_bool(own_body, own_masked, "Required"),
            "optional": extract_bool(own_body, own_masked, "Optional"),
            "computed": extract_bool(own_body, own_masked, "Computed"),
            "force_new": extract_bool(own_body, own_masked, "ForceNew"),
            "sensitive": extract_bool(own_body, own_masked, "Sensitive"),
            "default": extract_default(own_body, own_masked),
            "max_items": extract_int(own_body, own_masked, "MaxItems"),
            "min_items": extract_int(own_body, own_masked, "MinItems"),
            "conflicts_with": extract_string_list(own_body, own_masked, "ConflictsWith"),
            "required_with": extract_string_list(own_body, own_masked, "RequiredWith"),
            "exactly_one_of": extract_string_list(own_body, own_masked, "ExactlyOneOf"),
            "at_least_one_of": extract_string_list(own_body, own_masked, "AtLeastOneOf"),
        }
    )

    type_pos = locate_field(own_masked, "Type")
    if type_pos is not None:
        tm = re.match(r"Type\s*:\s*schema\.(\w+)", own_body[type_pos:])
        if tm:
            attr["type"] = tm.group(1)

    return attr


def parse_resource_func(sources, func_name):
    func_re = re.compile(r"func\s+" + re.escape(func_name) + r"\s*\(\)\s*\*schema\.Resource\s*\{")
    for filename, text in sources.items():
        m = func_re.search(text)
        if not m:
            continue
        matches = find_brace_matches(text)
        masked_text = build_masked(text)
        func_brace_pos = m.end() - 1
        if func_brace_pos not in matches:
            continue
        func_end = matches[func_brace_pos]
        func_masked_body = masked_text[func_brace_pos:func_end]

        schema_map_match = re.search(r"Schema:\s*map\[string\]\*schema\.Schema\s*\{", func_masked_body)
        if not schema_map_match:
            continue
        map_brace_pos = func_brace_pos + schema_map_match.end() - 1
        if map_brace_pos not in matches:
            continue
        map_close = matches[map_brace_pos]

        attrs = [
            parse_attr(text, masked_text, matches, name, astart, aend)
            for name, astart, aend in iter_top_level_entries(text, matches, map_brace_pos + 1, map_close)
        ]
        return filename, attrs
    return None, None


def render_type(attr):
    base = TYPE_NAMES.get(attr["type"], attr["type"] or "Unknown")
    if attr["type"] in ("TypeList", "TypeSet") and attr["element_type"]:
        base += f" of {TYPE_NAMES.get(attr['element_type'], attr['element_type'])}"
    elif attr["type"] == "TypeMap" and attr["element_type"]:
        base += f" of {TYPE_NAMES.get(attr['element_type'], attr['element_type'])}"
    return base


def render_notes(attr):
    notes = []
    if attr.get("deprecated"):
        notes.append(f"Deprecated: {attr['deprecated']}")
    if attr["force_new"]:
        notes.append("Changing this value forces recreation of the resource.")
    if attr["sensitive"]:
        notes.append("Value is treated as sensitive.")
    if attr["default"] is not None:
        notes.append(f"Default: `{attr['default']}`.")
    if attr["max_items"]:
        notes.append(f"Max items: {attr['max_items']}.")
    if attr["min_items"]:
        notes.append(f"Min items: {attr['min_items']}.")
    if attr["conflicts_with"]:
        notes.append(f"Conflicts with: {', '.join('`%s`' % x for x in attr['conflicts_with'])}.")
    if attr["required_with"]:
        notes.append(f"Required with: {', '.join('`%s`' % x for x in attr['required_with'])}.")
    if attr["exactly_one_of"]:
        notes.append(f"Exactly one of: {', '.join('`%s`' % x for x in attr['exactly_one_of'])}.")
    if attr["at_least_one_of"]:
        notes.append(f"At least one of: {', '.join('`%s`' % x for x in attr['at_least_one_of'])}.")
    return " ".join(notes)


def category(attr):
    if attr["required"]:
        return "Required"
    if attr["optional"]:
        return "Optional"
    if attr["computed"]:
        return "Read-Only"
    return "Unknown"


def render_attr_lines(attrs, lines, depth=0):
    indent = "  " * depth
    for attr in sorted(attrs, key=lambda a: a["name"]):
        line = f"{indent}- `{attr['name']}` ({render_type(attr)}, {category(attr)})"
        if attr["description"]:
            line += f" - {attr['description']}"
        notes = render_notes(attr)
        if notes:
            line += f" {notes}"
        lines.append(line)
        if attr["children"]:
            render_attr_lines(attr["children"], lines, depth + 1)


def render_markdown(resource_type, func_name, source_file, attrs):
    lines = []
    lines.append(f"# `{resource_type}`")
    lines.append("")
    lines.append(
        f"Auto-generated from `menandmice/{source_file}` (`{func_name}`) by "
        "`scripts/generate_resource_schema_docs.py`. Do not edit by hand - "
        "re-run the script instead."
    )
    lines.append("")

    for label, pred in (
        ("Required", lambda a: a["required"]),
        ("Optional", lambda a: a["optional"]),
        ("Read-Only", lambda a: a["computed"]),
    ):
        group = [a for a in attrs if pred(a)]
        if not group:
            continue
        lines.append(f"## {label}")
        lines.append("")
        render_attr_lines(group, lines)
        lines.append("")

    return "\n".join(lines).rstrip() + "\n"


def get_resources_map(provider_text):
    """Extract {resource_type: func_name} from provider.go's ResourcesMap."""
    m = re.search(r"ResourcesMap:\s*map\[string\]\*schema\.Resource\s*\{", provider_text)
    if not m:
        raise RuntimeError("Could not find ResourcesMap in provider.go")
    matches = find_brace_matches(provider_text)
    brace_pos = m.end() - 1
    close = matches[brace_pos]
    body = provider_text[brace_pos + 1: close]
    return dict(re.findall(r'"(menandmice_\w+)"\s*:\s*(\w+)\s*\(\s*\)', body))


def parse_args(argv=None):
    parser = argparse.ArgumentParser(
        prog="generate_resource_schema_docs.py",
        description=(
            "Extracts the Terraform schema (attributes, types, "
            "required/optional/computed, defaults, descriptions) directly "
            "from the Go source of this provider's resources, and writes "
            "one Markdown file per resource describing exactly what a "
            "customer needs to know to write/import that resource.\n\n"
            "This parses the schema.Schema map literals in menandmice/*.go "
            "(driven by the ResourcesMap in menandmice/provider.go) rather "
            "than the example .tf files or the tfplugindocs output in "
            "docs/, so the generated reference always matches the actual "
            "provider code."
        ),
        epilog=(
            "Output:\n"
            f"  One <resource_type>.md file per resource is written to\n"
            f"  {os.path.relpath(OUT_DIR, REPO_ROOT)}{os.sep} (gitignored) by default.\n\n"
            "Example:\n"
            "  python scripts/generate_resource_schema_docs.py\n"
            "  python scripts/generate_resource_schema_docs.py --out-dir /tmp/docs"
        ),
        formatter_class=argparse.RawDescriptionHelpFormatter,
    )
    parser.add_argument(
        "--out-dir",
        default=OUT_DIR,
        help=(
            "Directory to write the generated <resource_type>.md files into "
            f"(default: {os.path.relpath(OUT_DIR, REPO_ROOT)}{os.sep}, relative to the repo root)"
        ),
    )
    return parser.parse_args(argv)


def main(argv=None):
    args = parse_args(argv)

    sources = load_go_sources()
    provider_text = sources.get("provider.go")
    if not provider_text:
        print("provider.go not found in menandmice/", file=sys.stderr)
        return 1

    resources = get_resources_map(provider_text)
    if not resources:
        print("No resources found in provider.go ResourcesMap", file=sys.stderr)
        return 1

    os.makedirs(args.out_dir, exist_ok=True)

    for resource_type, func_name in sorted(resources.items()):
        source_file, attrs = parse_resource_func(sources, func_name)
        if attrs is None:
            print(f"warning: could not parse schema for {resource_type} ({func_name})", file=sys.stderr)
            continue
        markdown = render_markdown(resource_type, func_name, source_file, attrs)
        out_path = os.path.join(args.out_dir, f"{resource_type}.md")
        with open(out_path, "w", encoding="utf-8") as f:
            f.write(markdown)
        print(f"wrote {out_path}")

    return 0


if __name__ == "__main__":
    sys.exit(main())
