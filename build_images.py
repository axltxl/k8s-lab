import os
import sys
import subprocess
import yaml

from pathlib import Path

def get_yaml_root_key_value(yaml_path, key):
    if not os.path.exists(yaml_path):
        raise FileNotFoundError(f"File not found: {yaml_path}")

    with open(yaml_path, 'r') as f:
        try:
            content = yaml.safe_load(f)
        except yaml.YAMLError as e:
            raise ValueError(f"Failed to parse YAML: {e}")

    if key not in content:
        raise KeyError(f"Key '{key}' not found at root level of {yaml_path}")

    return content[key]

def main():
    # Root dir = 3 levels up from current script location
    script_dir = Path(__file__).resolve().parent
    project_root = script_dir.parents[2]
    os.chdir(project_root)

    # Read the registry host IP from YAML
    config_path = project_root / "config.yaml"
    docker_registry_host = get_yaml_root_key_value(config_path, "k8s_api_server_ip")

    # Path to the app
    app_dir = project_root / "apps" / "todo" / "api"
    os.chdir(app_dir)

    # Build and push image
    image = f"{docker_registry_host}:5000/k8s-lab/todod"

    print(f"📦 Building Docker image: {image}")
    subprocess.run(["docker", "build", ".", "-t", image], check=True)

    print(f"🚀 Pushing image to registry: {image}")
    subprocess.run(["docker", "push", image], check=True)

if __name__ == "__main__":
    try:
        import yaml
    except ImportError:
        print("Missing 'PyYAML'. Install it with: pip install pyyaml", file=sys.stderr)
        sys.exit(1)

    try:
        main()
    except Exception as e:
        print(f"❌ Error: {e}", file=sys.stderr)
        sys.exit(1)
