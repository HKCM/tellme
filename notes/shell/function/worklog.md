---
tags: ["工作记录"]
---
# worklog
---
## Example

```bash
# 用于记录，可以写入.function中
function worklog(){
  local YEAR=$(date +%Y)
  local WEEK=$(date +%U)
  local dir_path="$HOME/.wroklog/${YEAR}"
  local file_path="${dir_path}/week${WEEK}.md"
  if [ -z $1 ]; then
    [ -d "$dir_path" ] || mkdir -p "$dir_path"
    [ ! -e "$file_path" ] || touch "$file_path"
    vim $file_path
    return
  fi

  case $1 in
  show)
    cat ${file_path}
    ;;
  showall)
    tree -L 2 ${dir_path}
    ;;
  *)
    echo "$0         # 创建新的log"
    echo "$0 show    # 显示本周log"
    echo "$0 showall # 显示本年log"
    echo "$0 [show|showall|help] "
    ;;
  esac
}
```
