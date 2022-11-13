#! /bin/bash
INSTANCE_ID=i-01b524cbf5af100ca
ALLOCATION_ID=eipalloc-0419b4f0ed170c227
INSTANCE_STATE=$(aws ec2 describe-instances --instance-ids $INSTANCE_ID --output text --query 'Reservations[*].Instances[*].State.Name');

if test "$INSTANCE_STATE" = "running";
then
	echo $INSTANCE_ID "is "$INSTANCE_STATE
else
	echo $INSTANCE_ID "is "$INSTANCE_STATE
	aws ec2 start-instances --instance-ids $INSTANCE_ID
fi

aws ec2 associate-address --instance-id $INSTANCE_ID --allocation-id ALLOCATION_ID
