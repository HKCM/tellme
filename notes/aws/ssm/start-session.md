---
tags: ["ssm连接"]
---
# start-session
---
## Example

```bash
# 通过SSM连接到机器,如果机器处于STOP状态 则启动并等待
# 核心 aws ssm start-session --target i-04247b9xxxx
function connect() {
    local instanceId=i-04247b9xxxx
    local profile=dev-mfa
    local state=$(aws ec2 describe-instances --profile ${profile} --instance-ids ${instanceId} --query "Reservations[*].Instances[*].State.Name" --output text)
    if [ ${state} = "stopped" ]; then
        echo "Instance ${instanceId} is Stopped, starting instance... Please wait..."
        aws ec2 start-instances --instance-ids ${instanceId} --region ap-northeast-1 --profile ${profile} --output text
        aws ec2 wait instance-status-ok --profile ${profile} --instance-ids ${instanceId}
    fi
    echo "Instance ${instanceId} is ready.."
    echo "ssm start session... to server (${instanceId})"
    aws --profile ${profile} ssm start-session --target ${instanceId}
}
```
