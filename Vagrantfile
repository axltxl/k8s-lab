# Variables
@vm_box = "bento/ubuntu-24.04"
@vm_box_version = "202502.21.0"
@vm_cpus = 2
@vm_memory = "2048"

Vagrant.configure("2") do |config|

    config.vm.box = @vm_box
    config.vm.box_version = @vm_box_version

    config.vm.provider "virtualbox" do |vb|
        vb.memory = @vm_memory
        vb.cpus = @vm_cpus
    end

    # General provisioning
    config.vm.provision "shell", inline: <<-EOF
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
            sudo apt-get install -y kubelet kubectl
            sudo apt-mark hold kubelet kubectl

            # Enable and start the kubelet service.
            sudo systemctl enable --now kubelet
        fi

    EOF

    # Control plane
    # ----
    config.vm.define "master" do |master|
        master.vm.box = @vm_box
        master.vm.box_version = @vm_box_version

        # Network configuration
        master.vm.network "private_network", ip: "192.168.0.11"

        master.vm.provider "virtualbox" do |vb|
            vb.memory = @vm_memory
            vb.cpus = @vm_cpus
        end

        master.vm.provision "shell", inline: <<-EOF
            # Install kubeadm
            # ------------------------------------
            if [ -z "$(which kubeadm)" ]; then
                sudo apt-get install -y kubeadm
                sudo apt-mark hold kubeadm
            fi
        EOF
    end

    # Worker nodes
    # ----
    config.vm.define "n1" do |node|
        # Network configuration
        node.vm.network "private_network", ip: "192.168.0.21"
    end

    config.vm.define "n2" do |node|
        # Network configuration
        node.vm.network "private_network", ip: "192.168.0.22"
    end
end
