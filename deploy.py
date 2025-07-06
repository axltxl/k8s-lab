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

def build_and_push_docker_image(project_root):
    """
    Builds and pushes a Docker image for the specified project.

    This function changes the working directory to the project root, reads the
    Docker registry host IP from a YAML configuration file, navigates to the
    application directory, builds a Docker image, and pushes it to the specified
    Docker registry.

    Args:
      project_root (Path): The root directory of the project.

    Raises:
      FileNotFoundError: If the configuration file or application directory does not exist.
      subprocess.CalledProcessError: If the Docker build or push command fails.

    Notes:
      - The registry host IP is retrieved from the "k8s_api_server_ip" key in the
        "config.yaml" file located in the project root.
      - The Docker image is tagged with the format:
        "<docker_registry_host>:5000/k8s-lab/todod".
      - Ensure Docker is installed and running, and the user has permissions to
        build and push images.
    """
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


def k8s_deploy_manifests(project_root):
    """
    Deploys Kubernetes manifests for a project.

    This function sets up the environment, applies Kubernetes manifests, and patches
    a specific deployment with an updated container image.

    Args:
      project_root (Path): The root directory of the project.

    Steps:
      1. Changes the working directory to the project root.
      2. Sets the KUBECONFIG environment variable to point to the kubeconfig file.
      3. Reads the Docker registry host IP from the project's `config.yaml` file.
      4. Changes the working directory to the manifests directory.
      5. Applies all Kubernetes manifests in the manifests directory using `kubectl apply`.
      6. Reads and updates the `todod-deployment.yaml` file with the new container image.
      7. Applies the updated deployment manifest using `kubectl apply`.

    Raises:
      subprocess.CalledProcessError: If any `kubectl` command fails.
      FileNotFoundError: If required files (e.g., kubeconfig, config.yaml, or deployment manifest) are missing.
      KeyError: If the expected keys are not found in the YAML files.
    """

    os.chdir(project_root)

    # Set the KUBECONFIG environment variable
    kubeconfig_path = project_root / "files" / "local" / "k8s" / "users" / "k8s-lab-admin.kubeconfig"
    os.environ["KUBECONFIG"] = str(kubeconfig_path)

    # Read the registry host IP from YAML
    config_path = project_root / "config.yaml"
    docker_registry_host = get_yaml_root_key_value(config_path, "k8s_api_server_ip")

    # Path to the manifests
    manifests_dir = project_root / "k8s" / "manifests"

    # Apply the manifests
    subprocess.run(["kubectl", "apply", "-Rf", manifests_dir], check=True)

def main():
    # Root dir = 3 levels up from current script location
    script_dir = Path(__file__).resolve().parent

    # Build and push the Docker images
    print("🔨 Building and pushing Docker images...")
    build_and_push_docker_image(script_dir)

    # Apply base Kubernetes manifests
    print("📦 Applying Kubernetes manifests...")
    k8s_deploy_manifests(script_dir)

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
