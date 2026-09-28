#!/usr/bin/env bash
# Generate shields.io endpoint JSON badges from badges/theme.json.
# Run: make badges
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
THEME="$ROOT/badges/theme.json"
OUT="$ROOT/badges"

if [[ ! -f "$THEME" ]]; then
  echo "missing $THEME" >&2
  exit 1
fi

mkdir -p "$OUT"

python3 - "$THEME" "$OUT" "$ROOT" <<'PY'
import json, os, re, subprocess, sys
from pathlib import Path

theme_path, out_dir, root = Path(sys.argv[1]), Path(sys.argv[2]), Path(sys.argv[3])
theme = json.loads(theme_path.read_text(encoding="utf-8"))

label_color = theme["labelColor"].lstrip("#")
color = theme["color"].lstrip("#")
style = theme.get("style", "flat-square")
repo = theme["repo"]
branch = theme.get("branch", "master")

version = os.environ.get("VERSION", "").strip()
if not version:
    try:
        version = subprocess.run(
            ["git", "describe", "--tags", "--abbrev=0"],
            cwd=root, capture_output=True, text=True, check=True,
        ).stdout.strip()
    except Exception:
        version = "untagged"

go_mod = (root / "go.mod").read_text(encoding="utf-8")
gm = re.search(r"(?m)^go\s+(\d+(?:\.\d+)*)", go_mod)
go_ver = gm.group(1) if gm else "1"
dep_count = len(re.findall(r"(?m)^\s*require\s+[^\s(]|^\t[^\s/]+/\S+", go_mod))

def write(name: str, payload: dict) -> None:
    payload.setdefault("schemaVersion", 1)
    payload.setdefault("labelColor", label_color)
    payload.setdefault("color", color)
    payload.setdefault("style", style)
    path = out_dir / name
    path.write_text(json.dumps(payload, indent=2) + "\n", encoding="utf-8")
    print(f"wrote {path.relative_to(root)}")

ver_msg = version
if re.fullmatch(r"\d+(\.\d+)*", version):
    ver_msg = f"v{version}"
write("version.json", {
    "label": "version",
    "message": ver_msg,
    "namedLogo": "github",
    "logoColor": color,
})

write("license.json", {
    "label": "license",
    "message": theme.get("license", "0BSD"),
})

write("go.json", {
    "label": "go",
    "message": go_ver,
    "namedLogo": "go",
    "logoColor": color,
})

write("deps.json", {
    "label": "dependencies",
    "message": str(dep_count),
})

write("accuracy.json", {
    "label": "accuracy",
    "message": theme["accuracy"],
})

write("differential.json", {
    "label": "differential",
    "message": theme["references"],
})

openssf_msg = "scorecard"
try:
    import urllib.request
    req = urllib.request.Request(
        f"https://api.scorecard.dev/projects/github.com/{repo}",
        headers={"Accept": "application/json", "User-Agent": "mgrs-go-badges"},
    )
    with urllib.request.urlopen(req, timeout=15) as resp:
        score_data = json.load(resp)
    score = score_data.get("score")
    if isinstance(score, (int, float)):
        openssf_msg = f"{score:.1f}/10"
except Exception:
    pass

write("openssf.json", {
    "label": "openssf",
    "message": openssf_msg,
    "namedLogo": "openssf",
    "logoColor": color,
})

raw = f"https://raw.githubusercontent.com/{repo}/badges"
enc = __import__("urllib.parse").parse.quote

def endpoint(file: str) -> str:
    return f"https://img.shields.io/endpoint?url={enc(f'{raw}/{file}', safe='')}"

ci = (
    f"https://img.shields.io/github/actions/workflow/status/{repo}/ci.yml"
    f"?branch={branch}&style={style}&label=ci&labelColor={label_color}&color={color}"
)

lines = [
    f'[![CI]({ci})](https://github.com/{repo}/actions/workflows/ci.yml)',
    f'[![OpenSSF]({endpoint("openssf.json")})](https://scorecard.dev/viewer/?uri=github.com/{repo})',
    f'[![version]({endpoint("version.json")})](https://github.com/{repo}/releases)',
    f'[![license]({endpoint("license.json")})](https://github.com/{repo}/blob/{branch}/LICENSE)',
    f'[![go]({endpoint("go.json")})](https://go.dev/dl/)',
    f'[![deps]({endpoint("deps.json")})](https://github.com/{repo}/blob/{branch}/go.mod)',
    f'[![accuracy]({endpoint("accuracy.json")})](https://github.com/{repo}#verification)',
    f'[![differential]({endpoint("differential.json")})](https://github.com/{repo}#verification)',
]

snippet = out_dir / "README.snippet.md"
snippet.write_text(" ".join(lines) + "\n", encoding="utf-8")
print(f"wrote {snippet.relative_to(root)}")
PY
