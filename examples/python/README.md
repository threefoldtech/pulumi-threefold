To run the Python examples, you'll need `pulumi` and `python3` already installed. The `pulumi_threefold` Python package is included with this repository.

The recommended way to install the required Python packages is to create a virtual environment and install them inside:

It's recommended to create a virtual environment to install the required Python packages. You can use a single virtual environment to run all of the examples:

```
python -m venv venv
source venv/bin/activate
pip install -r requirements.txt
```

Remember that you'll need to source the virtual environment again anytime you open a new shell and want to use Pulumi to manage ThreeFold deployments written in Python.

Next, if you didn't already, login to Pulumi. This example will store Pulumi state locally in your home directory (check `--help` for more info):

```
pulumi login --local
```

Finally, bring the deployment up and follow the prompts to complete it. Make sure to set your mnemonic and network first:

```
cd virtual_machine
export MNEMONIC="your words here ..."
export NETWORK="main"
pulumi up
```

Then to destroy the deployment:

```
pulumi down
```
