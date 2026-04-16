import shutil
from pathlib import Path

from jinja2 import Environment, StrictUndefined

from .io import fail


_template_env = Environment(
    autoescape=False,
    keep_trailing_newline=True,
    undefined=StrictUndefined,
    variable_start_string="[[",
    variable_end_string="]]",
)


def scaffold_from_template(service_dir: Path, template_name: str, replacements: dict[str, str]) -> None:
    templates_root = Path(__file__).resolve().parents[3] / "templates"
    template_dir = templates_root / template_name
    scaffold_from_directory(service_dir, template_dir, replacements, templates_root / "charts" / "app")


def scaffold_from_directory(
    service_dir: Path,
    template_dir: Path,
    replacements: dict[str, str],
    common_chart_dir: Path | None = None,
) -> None:
    template_dir = template_dir.resolve()

    if not template_dir.is_dir():
        fail(f"template directory does not exist: {template_dir}")

    if service_dir.exists():
        shutil.rmtree(service_dir)
    service_dir.mkdir(parents=True, exist_ok=True)

    copy_template_tree(template_dir, service_dir, replacements)

    if common_chart_dir is not None and common_chart_dir.is_dir():
        chart_target_dir = service_dir / "charts" / "app"
        chart_target_dir.mkdir(parents=True, exist_ok=True)
        copy_template_tree(common_chart_dir, chart_target_dir, replacements)


def copy_template_tree(source_dir: Path, target_dir: Path, replacements: dict[str, str]) -> None:
    for path in source_dir.rglob("*"):
        relative_path = path.relative_to(source_dir)
        destination = target_dir / relative_path

        if path.is_dir():
            destination.mkdir(parents=True, exist_ok=True)
            continue

        destination.parent.mkdir(parents=True, exist_ok=True)
        copy_template_file(path, destination, replacements)


def copy_template_file(source: Path, destination: Path, replacements: dict[str, str]) -> None:
    try:
        raw = source.read_text(encoding="utf-8")
    except UnicodeDecodeError:
        shutil.copy2(source, destination)
        return

    destination.write_text(render_template(raw, replacements), encoding="utf-8")


def render_template(content: str, replacements: dict[str, str]) -> str:
    return _template_env.from_string(content).render(**replacements)
