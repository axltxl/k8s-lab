# Variables
# ------------------------------------
@vm_box = "bento/ubuntu-24.04"
@vm_box_version = "202502.21.0"
@vm_cpus = 2
@vm_memory_worker_nodes = "2048"
@vm_memory_control_plane = "2048"
@k8s_num_worker_nodes = 1

# Constants
# ------------------------------------
K8S_API_SERVER_IP = "192.168.0.11" # Control plane host IP
K8S_CNI_NETWORK_CIDR = "172.16.0.0/16" # Pod network CIDR

Vagrant.configure("2") do |config|

    config.vm.box = @vm_box
    config.vm.box_version = @vm_box_version

    config.vm.provider "virtualbox" do |vb|
        vb.cpus = @vm_cpus
        vb.gui = true
    end

    # General provisioning
    config.vm.provision "shell", inline: <<-SHELL
        sudo apt-get update -y

        # Disable swap partition
        # source: https://kubernetes.io/docs/setup/production-environment/tools/kubeadm/install-kubeadm/#swap-configuration
        sudo swapoff -a
        sudo sed -i 's@/swap@#/swap@' /etc/fstab

        # Disable swap on boot
        # sudo sed -i '/^[^#]*\/swap.img/s/^/#/' /etc/fstab


        # Install container runtime (containerd)
        # Add Docker's official GPG key:
        if [ -z "$(which containerd)" ]; then
            sudo apt-get install ca-certificates curl
            sudo install -m 0755 -d /etc/apt/keyrings
            sudo curl -fsSL https://download.docker.com/linux/ubuntu/gpg -o /etc/apt/keyrings/docker.asc
            sudo chmod a+r /etc/apt/keyrings/docker.asc

            # Add the repository to Apt sources:
            echo \
            "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.asc] https://download.docker.com/linux/ubuntu \
            $(. /etc/os-release && echo "${UBUNTU_CODENAME:-$VERSION_CODENAME}") stable" | \
            sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
            sudo apt-get update
            sudo apt-get install containerd.io

            # Configure containerd
            # ------------------------------------
            sudo mkdir -p /etc/containerd
            sudo cp /vagrant/files/remote/containerd/config.toml /etc/containerd/config.toml

            # disable a bug in ubuntu 22.04
            # which prevents you from deleting pods
            # -------------------------------------
            sudo systemctl stop apparmor.service
            sudo systemctl disable apparmor.service
        fi

        # Enable IPv4 forwarding
        # ------------------------------------
        # sysctl params required by setup, params persist across reboots
        if [ ! -f /etc/sysctl.d/k8s.conf ]; then
            # sysctl params required by setup, params persist across reboots
            cat << EOF | sudo tee /etc/sysctl.d/k8s.conf
net.ipv4.ip_forward = 1
EOF

            # Apply sysctl params without reboot
            sudo sysctl --system
        fi

        # Install kubelet and kubectl
        # ------------------------------------
        if [ -z "$(which kubectl)" ]; then
            sudo apt-get install -y apt-transport-https ca-certificates curl gpg

            # Download the public signing key for the Kubernetes package repositories.
            curl -fsSL https://pkgs.k8s.io/core:/stable:/v1.32/deb/Release.key | sudo gpg --dearmor -o /etc/apt/keyrings/kubernetes-apt-keyring.gpg

            # Add the Kubernetes package repository to the system's sources list.
            echo 'deb [signed-by=/etc/apt/keyrings/kubernetes-apt-keyring.gpg] https://pkgs.k8s.io/core:/stable:/v1.32/deb/ /' | sudo tee /etc/apt/sources.list.d/kubernetes.list

            # Update the package index and install kubeadm, kubelet, and kubectl.
            sudo apt-get update
            sudo apt-get install -y kubelet kubeadm kubectl
            sudo apt-mark hold kubelet kubeadm kubectl
        fi

        # Restart services regardless of the state
        # (in this way, we make sure services are running with the latest configuration)
        sudo systemctl restart containerd
    SHELL

    # Control plane
    # ----
    config.vm.define "cp" do |cp|
        cp.vm.box = @vm_box
        cp.vm.box_version = @vm_box_version

        # Network configuration
        cp.vm.network "private_network", ip: K8S_API_SERVER_IP

        # Resource configuration
        cp.vm.provider "virtualbox" do |vb|
            vb.memory = @vm_memory_control_plane
        end

        cp.vm.provision "shell", inline: <<-SHELL
            # Set hostname
            sudo hostnamectl set-hostname control-plane

            # Bootstrapping the Kubernetes control plane
            # ------------------------------------
            if [ ! -d /etc/kubernetes/pki ]; then
                # Initialize the Kubernetes cluster
                # (TLS PKI is generated automatically at /etc/kubernetes/pki)
                # ------------------------------------
                # see: files/remote/kubeadm/kubeadm-config.yaml for details
                # on how the cluster is being configured
                # (e.g. pod network CIDR, certificate SANs, etc.)
                sudo kubeadm init --config /vagrant/files/remote/k8s/kubeadm-config.yaml

            fi

            # Configure kubectl for the vagrant user
            # ------------------------------------
            if [ ! -f /home/vagrant/.kube/config ]; then
                # Create the .kube directory
                mkdir -p /home/vagrant/.kube
                # Copy the admin.conf file to the .kube directory
                sudo cp -i /etc/kubernetes/admin.conf /home/vagrant/.kube/config
                # Set the ownership of the .kube directory and its contents to the vagrant user
                sudo chown -R vagrant:vagrant /home/vagrant/.kube
            fi

            # Set up kubelet flags
            cat << EOF | sudo tee /etc/default/kubelet > /dev/null
KUBELET_EXTRA_ARGS="--node-ip=#{K8S_API_SERVER_IP}"
EOF

            # Restart kubelet service
            sudo systemctl restart kubelet


            # Explicitly set the KUBECONFIG environment variable
            # (otherwise kubectl will not work, even if the config file is in the right place)
            export KUBECONFIG=/home/vagrant/.kube/config

            if ! which cilium > /dev/null; then
                # Install the CNI plugin (Cilium)
                # Pod network configuration
                # source: https://docs.cilium.io/en/stable/gettingstarted/k8s-install-default/
                # ------------------------------------
                CILIUM_CLI_VERSION=$(curl -s https://raw.githubusercontent.com/cilium/cilium-cli/main/stable.txt)
                CLI_ARCH=amd64
                if [ "$(uname -m)" = "aarch64" ]; then CLI_ARCH=arm64; fi
                curl -L --fail --remote-name-all https://github.com/cilium/cilium-cli/releases/download/${CILIUM_CLI_VERSION}/cilium-linux-${CLI_ARCH}.tar.gz{,.sha256sum}
                sha256sum --check cilium-linux-${CLI_ARCH}.tar.gz.sha256sum
                sudo tar xzvfC cilium-linux-${CLI_ARCH}.tar.gz /usr/local/bin
                rm cilium-linux-${CLI_ARCH}.tar.gz{,.sha256sum}

                # Install Cilium CNI plugin
                # ------------------------------------
                cilium install \
                    --set ipam.operator.clusterPoolIPv4PodCIDRList='#{K8S_CNI_NETWORK_CIDR}' \
                    --set k8sServiceHost=#{K8S_API_SERVER_IP} \
                    --set k8sServicePort=6443
            fi

            # Generate the join command for worker nodes into a script
            # ------------------------------------
            # (this command will be used by worker nodes to join the cluster)
            k8s_node_join_script=/vagrant/scripts/remote/temp/k8s_node_join.sh
            kubeadm token create --print-join-command > $k8s_node_join_script
            chmod +x $k8s_node_join_script

            # Install kubeadduser: an utility script to add users to the cluster
            # ------------------------------------
            if [ ! -f /usr/local/bin/kubeadduser ]; then
                sudo cp /vagrant/scripts/remote/k8s/kubeadduser /usr/local/bin/kubeadduser
                sudo chmod +x /usr/local/bin/kubeadduser
            fi

            # Create a default admin user for the cluster
            # This will create a kubeconfig for the k8s-lab-admin user
            # whose credentials will be store in: /vagrant/files/local/k8s/users
            # and can be used to access the cluster
            # ------------------------------------
            kubeadduser #{K8S_API_SERVER_IP}
        SHELL
    end

    # Worker nodes
    # ----
    (1..@k8s_num_worker_nodes).each do |i|

        # Define the worker node IP address
        k8s_node_ip = "192.168.0.2#{i}"

        config.vm.define "n#{i}" do |node|
            node.vm.provision "shell", inline: <<-SHELL
                # Set hostname
                sudo hostnamectl set-hostname n#{i}

                # Set up kubelet flags
                cat << EOF | sudo tee /etc/default/kubelet > /dev/null
KUBELET_EXTRA_ARGS="--node-ip=#{k8s_node_ip}"
EOF


                # Restart kubelet service
                sudo systemctl restart kubelet

                # Join the Kubernetes cluster
                # (see: control-plane provisioning for details)
                # ------------------------------------
                sudo /vagrant/scripts/remote/temp/k8s_node_join.sh
            SHELL

            # Network configuration
            node.vm.network "private_network", ip: k8s_node_ip

            # Resource configuration
            node.vm.provider "virtualbox" do |vb|
                vb.memory = @vm_memory_worker_nodes
            end
        end
    end
end
