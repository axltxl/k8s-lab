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
