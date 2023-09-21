current_dir = $(shell pwd)

run:
	rm -rf ${current_dir}/state 
	mkdir ${current_dir}/state
	pulumi login --cloud-url file://${current_dir}/state
	pulumi stack init test
	pulumi up --yes

destroy:
	pulumi destroy --yes
	pulumi stack rm --yes
	pulumi logout