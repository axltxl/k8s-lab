import os
import sys
import subprocess
import yaml
import time

from pathlib import Path

# Root dir = 3 levels up from current script location
script_dir = Path(__file__).resolve().parent

# Set the KUBECONFIG environment variable
kubeconfig_path = script_dir / "files" / "local" / "k8s" / "users" / "k8s-lab-admin.kubeconfig"
os.environ["KUBECONFIG"] = str(kubeconfig_path)

def generate_tag():
    """
    Generates a unique tag based on the current Git commit hash and the current Unix timestamp.

    Returns:
      str: A tag string in the format "<git_commit_hash>-<timestamp>", where <git_commit_hash> is the
         first 7 characters of the current Git commit hash and <timestamp> is the current time in seconds since the epoch.

    Raises:
      RuntimeError: If the Git commit hash cannot be retrieved.
    """
    try:
      result = subprocess.run(["git", "rev-parse", "HEAD"], check=True, capture_output=True, text=True)
      git_commit_hash = result.stdout.strip()[:7]

      # Get current number of seconds since epoch
      # Format the tag as "git_commit_hash-time_now"
      tag = f"{git_commit_hash}-{int(time.time())}"
      return tag
    except subprocess.CalledProcessError as e:
      raise RuntimeError(f"Failed to get Git commit hash: {e}")

def get_config_from_yaml(yaml_path):
    """
    Loads and parses a YAML configuration file.

    Args:
      yaml_path (str): The path to the YAML file to be loaded.

    Returns:
      dict: The contents of the YAML file as a Python dictionary.

    Raises:
      FileNotFoundError: If the specified YAML file does not exist.
      ValueError: If the YAML file cannot be parsed.
    """
    if not os.path.exists(yaml_path):
        raise FileNotFoundError(f"File not found: {yaml_path}")

    with open(yaml_path, 'r') as f:
        try:
            content = yaml.safe_load(f)
        except yaml.YAMLError as e:
            raise ValueError(f"Failed to parse YAML: {e}")

    return content

# Load the configuration from the YAML file
config = get_config_from_yaml(script_dir / "config.yaml")

# def build_and_push_docker_image(project_root, image_name, image_tag=None):
def build_and_push_docker_image(app_name, app_dir, image_tag=None):
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
      - The registry host IP is retrieved from the "k8s_cplane_addr" key in the
        "config.yaml" file located in the project root.
      - The Docker image is tagged with the format:
        "<docker_registry_host>:5000/k8s-lab/todod".
      - Ensure Docker is installed and running, and the user has permissions to
        build and push images.
    """
    global script_dir
    global config

    # Read the registry host IP from YAML
    docker_registry_host = config.get("k8s_cplane_addr")

    # Path to the app
    os.chdir(app_dir)

    # Build and push image
    if image_tag is None:
        image_tag = 'latest'
    image = f"{docker_registry_host}:5000/k8s-lab/{app_name}:{image_tag}"

    print(f"📦 Building Docker image: {image}")
    subprocess.run(["docker", "build", ".", "-t", image], check=True)

    print(f"🚀 Pushing image to registry: {image}")
    subprocess.run(["docker", "push", image], check=True)


def k8s_deploy_manifests():
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

    global script_dir
    os.chdir(script_dir)

    # Path to the manifests
    manifests_dir = script_dir / "k8s" / "manifests"
    # If there are no YAML files, return
    if not any(manifests_dir.glob("*.yaml")):
        print("No Kubernetes manifests found to apply.")
        return

    # Apply the manifests
    subprocess.run(["kubectl", "apply", "-Rf", manifests_dir], check=True)


def k8s_install_helm_charts(image_tag):
    """
    Installs Helm charts for a project.

    This function sets up the environment, installs the Helm chart for the Todo application,
    and patches the deployment with the correct image.

    Args:
      image_tag (str): The tag for the Docker image to be used in the deployment.

    Steps:
      1. Changes the working directory to the project root.
      2. Sets the KUBECONFIG environment variable to point to the kubeconfig file.
      3. Reads the Docker registry host IP from the project's `config.yaml` file.
      4. Changes the working directory to the Helm charts directory.
      5. Installs the Helm chart for the Todo application using `helm install`.
      6. Reads and updates the `todod-deployment.yaml` file with the new container image.
      7. Applies the updated deployment manifest using `kubectl apply`.

    Raises:
      subprocess.CalledProcessError: If any `helm` or `kubectl` command fails.
      FileNotFoundError: If required files (e.g., kubeconfig, config.yaml, or deployment manifest) are missing.
      KeyError: If the expected keys are not found in the YAML files.
    """

    global script_dir
    global config
    os.chdir(script_dir)

    # Path to Helm charts
    helm_charts_dir = script_dir / "k8s" / "charts"

    # docker registry host IP from YAML
    docker_registry_host = config.get("k8s_cplane_addr")

    # Install Helm chart
    for directory in helm_charts_dir.iterdir():
        if directory.is_dir():
            release_name = directory.name
            chart = str(directory)
            print(f"📦 Installing Helm chart: {release_name}")
            subprocess.run([
                "helm", "upgrade", "--install", release_name, chart,
                "-n", release_name,  # Namespace for the release
                "--create-namespace",  # Create namespace if it doesn't exist
                "--set", f"deployment.image.repository={docker_registry_host}:5000/k8s-lab/{release_name}",
                "--set", f"deployment.image.tag={image_tag}"
            ], check=True)
            print(f"🚀 Helm chart {release_name} installed successfully.")

def main():

    # Build and push the Docker images
    print("🔨 Building and pushing Docker images...")

    # Get image tag
    image_tag = generate_tag()

    # Build todo application image
    build_and_push_docker_image("todod", f"{script_dir}/apps/todo/api", image_tag=image_tag)

    # Apply base Kubernetes manifests
    print("📦 Applying Kubernetes manifests...")
    k8s_deploy_manifests()

    # Install all Helm charts
    print("📦 Installing Helm charts...")
    k8s_install_helm_charts(image_tag=image_tag)

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
