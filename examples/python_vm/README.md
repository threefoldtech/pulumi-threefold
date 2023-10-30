To run this example, you'll need `pulumi` and `python` already installed. The recommended way to install the required Python packages is to create a virtual environment and install them inside:

```
python -m venv venv
source venv/bin/activate
pip install -r requirements.txt
```

Remember to source the virtual environment again anytime you want to use the deployment file with a new shell.

If you didn't already specify somewhere to store your Pulumi state, you can store it in a local folder like this:

```
mkdir state
pulumi login file://./state
```

Finally, bring the deployment up and follow the prompts to complete it. Make sure to set your mnemonic and network first:

```
export MNEMONIC="your words here ..."
export NETWORK="main"
pulumi up
```

Then to destroy the deployment:

```
pulumi down
```
