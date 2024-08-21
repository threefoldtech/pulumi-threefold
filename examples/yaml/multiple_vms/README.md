# multiple vms

This examples deploys 3 vms on GE farms and prints planetary and mycelium IPs

To run examples, make sure you have a mnemonic and a network set.

```bash

export MNEMONIC="mnemonic words"
export NETWORK="network" # dev, qa, test, main -> default is dev
```

- Go to the examples directory `cd yaml/multiple_vms`
- Run the example `make run`
- To cleanup the resources that you created `make destroy`


## example of the run

```
➜  multiple_vms git:(development_readme_edits) ✗ make run    
rm -rf /home/xmonader/wspace/pulumi-threefold/examples/yaml/multiple_vms/state 
mkdir /home/xmonader/wspace/pulumi-threefold/examples/yaml/multiple_vms/state
pulumi login --cloud-url file:///home/xmonader/wspace/pulumi-threefold/examples/yaml/multiple_vms/state
Logged in to xmonader-ThinkPad-E580 as xmonader (file:///home/xmonader/wspace/pulumi-threefold/examples/yaml/multiple_vms/state)
pulumi stack init test
Created stack 'test'
Enter your passphrase to protect config/secrets:  
Re-enter your passphrase to confirm:  
pulumi up --yes
Enter your passphrase to unlock config/secrets
    (set PULUMI_CONFIG_PASSPHRASE or PULUMI_CONFIG_PASSPHRASE_FILE to remember):  
Enter your passphrase to unlock config/secrets
Previewing update (test):
     Type                              Name                   Plan       Info
 +   pulumi:pulumi:Stack               pulumi-threefold-test  create     2 messages
 +   ├─ pulumi:providers:threefold     provider               create     
 +   ├─ threefold:provider:Scheduler   scheduler              create     
 +   ├─ threefold:provider:Network     network                create     
 +   └─ threefold:provider:Deployment  deployment             create     

Diagnostics:
  pulumi:pulumi:Stack (pulumi-threefold-test):
    5:12PM INF starting peer session=tf-1370865 twin=10466

    threefold grid provider setup

Outputs:
    mycelium_ip1      : output<string>
    mycelium_ip2      : output<string>
    mycelium_ip3      : output<string>
    node_deployment_id: output<string>
    planetary_ip1     : output<string>
    planetary_ip2     : output<string>
    planetary_ip3     : output<string>

Resources:
    + 5 to create

Updating (test):
     Type                              Name                   Status              Info
 +   pulumi:pulumi:Stack               pulumi-threefold-test  created (55s)       2 messages
 +   ├─ pulumi:providers:threefold     provider               created (0.00s)     
 +   ├─ threefold:provider:Scheduler   scheduler              created (18s)       
 +   ├─ threefold:provider:Network     network                created (5s)        
 +   └─ threefold:provider:Deployment  deployment             created (26s)       

Diagnostics:
  pulumi:pulumi:Stack (pulumi-threefold-test):
    5:12PM INF starting peer session=tf-1370948 twin=10466

    threefold grid provider setup

Outputs:
    mycelium_ip1      : "48b:f156:b548:4ba8:ff0f:c6e1:8a56:62d5"
    mycelium_ip2      : "48b:f156:b548:4ba8:ff0f:9474:96c9:6a4e"
    mycelium_ip3      : "48b:f156:b548:4ba8:ff0f:2d39:7cc9:b410"
    node_deployment_id: {
        310: 608143
    }
    planetary_ip1     : "300:c0f2:7fbb:8e29:cc2d:85a5:7b14:ea52"
    planetary_ip2     : "300:c0f2:7fbb:8e29:2e5b:3c5d:6fb4:27dc"
    planetary_ip3     : "300:c0f2:7fbb:8e29:736c:d326:6286:3f5b"

Resources:
    + 5 created

Duration: 56s

```