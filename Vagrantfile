# Variables
@vm_box = "bento/ubuntu-24.04"
@vm_box_version = "202502.21.0"
@vm_cpus = 2
@vm_memory = "3184"
@k8s_num_worker_nodes = 1

Vagrant.configure("2") do |config|

    config.vm.box = @vm_box
    config.vm.box_version = @vm_box_version

    config.vm.provider "virtualbox" do |vb|
        vb.memory = @vm_memory
        vb.cpus = @vm_cpus
    end

    # General provisioning
    config.vm.provision "shell", inline: <<-SHELL
        sudo apt-get update -y

        # Disable swap
        # source: https://kubernetes.io/docs/setup/production-environment/tools/kubeadm/install-kubeadm/#swap-configuration
        sudo swapoff -a

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
        sudo systemctl restart kubelet
    SHELL

    # Control plane
    # ----
    config.vm.define "cp" do |cp|
        cp.vm.box = @vm_box
        cp.vm.box_version = @vm_box_version

        # Network configuration
        cp.vm.network "private_network", ip: "192.168.0.11"

        cp.vm.provider "virtualbox" do |vb|
            vb.memory = @vm_memory
            vb.cpus = @vm_cpus
        end

        cp.vm.provision "shell", inline: <<-SHELL
            # Set hostname
            sudo hostnamectl set-hostname control-plane

            if [ ! -d /etc/kubernetes/pki ]; then
                # Initialize the Kubernetes cluster
                # (TLS PKI is generated automatically at /etc/kubernetes/pki)
                # ------------------------------------
                sudo kubeadm init
            fi
        SHELL
    end

    # Worker nodes
    # ----
    # FIXME: include this when the time comes
    # (1..@k8s_num_worker_nodes).each do |i|
    #     config.vm.define "n#{i}" do |node|
    #         node.vm.provision "shell", inline: <<-SHELL
    #             # Set hostname
    #             sudo hostnamectl set-hostname n#{i}
    #         SHELL
    #         # Network configuration
    #         node.vm.network "private_network", ip: "192.168.0.2#{i}"
    #     end
    # end
end
