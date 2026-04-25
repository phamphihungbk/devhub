import sys
from pathlib import Path

sys.path.insert(0, str(Path(__file__).resolve().parents[2]))

from scaffolders.template_runner import run_template_plugin


if __name__ == "__main__":
    run_template_plugin(Path(__file__).resolve().parent)
