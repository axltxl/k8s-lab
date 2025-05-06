# Variables
@vm_box = "bento/ubuntu-24.04"
@vm_box_version = "202502.21.0"
@vm_cpus = 2
@vm_memory = "2048"

# Control plane
# ----
Vagrant.configure("2") do |config|
    config.vm.define "master" do |master|
        master.vm.box = @vm_box
        master.vm.box_version = @vm_box_version

        master.vm.provider "virtualbox" do |vb|
            vb.memory = @vm_memory
            vb.cpus = @vm_cpus
        end

        master.vm.provision "shell", inline: <<-EOF
        EOF
    end
end

# Worker nodes
# ----
Vagrant.configure("2") do |config|

    config.vm.box = @vm_box
    config.vm.box_version = @vm_box_version

    config.vm.provider "virtualbox" do |vb|
        vb.memory = @vm_memory
        vb.cpus = @vm_cpus
    end

    config.vm.provision "shell", inline: <<-EOF
    EOF

    config.vm.define "worker1" do |worker|
    end

    config.vm.define "worker2" do |worker|
    end
end
