---
tags: ["token","临时凭证"]
---
# get-session-token
---
## Example

```bash
function mfa_dev() {
    local PROFILE=dev
    local MFA_PROFILE_NAME=dev-mfa
    local DEV_AWS_ID=012345678912
    if [ -z $1 ];then
        echo "Please input the code: "
        read code
    else
        code=$1
    fi

    # MFA获取数据
    local SESSION_JSON=$(aws --profile ${PROFILE} sts get-session-token --serial-number arn:aws:iam::$DEV_AWS_ID:mfa/karl.huang --token-code $code --output json)
    if [ $? -ne 0 ]; then
        echo "get-session-token failed, nothing to do"
        return
    fi

    # 提取关键数据
    local MFA_ACCESS_KEY=$(echo $SESSION_JSON | jq -r '.Credentials.AccessKeyId')
    local MFA_SECRET_ACCESS_KEY=$(echo $SESSION_JSON | jq -r '.Credentials.SecretAccessKey')
    local MFA_SESSION_TOKEN=$(echo $SESSION_JSON | jq -r '.Credentials.SessionToken')
    local MFA_EXPIRATION=$(echo $SESSION_JSON | jq -r '.Credentials.Expiration')

    # 配置新的profile
    aws --profile $MFA_PROFILE_NAME configure set aws_access_key_id $MFA_ACCESS_KEY
    aws --profile $MFA_PROFILE_NAME configure set aws_secret_access_key $MFA_SECRET_ACCESS_KEY
    aws --profile $MFA_PROFILE_NAME configure set aws_session_token $MFA_SESSION_TOKEN
    aws --profile $MFA_PROFILE_NAME configure set region ap-northeast-1
    aws --profile $MFA_PROFILE_NAME configure set output json
    echo "New credentials have been set successfully. (profile: $MFA_PROFILE_NAME, expiration: $MFA_EXPIRATION)"
}
```

https://awscli.amazonaws.com/v2/documentation/api/latest/reference/sts/get-session-token.html
