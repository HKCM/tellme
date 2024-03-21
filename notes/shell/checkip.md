---
tags: ["getip"]
---
# CheckIP
---
## Example

```bash
curl checkip.amazonaws.com

# EC2
curl http://169.254.169.254/latest/meta-data/local-ipv4 # Get private IPv4
curl http://169.254.169.254/latest/meta-data/public-ipv4 # Get public IPv4
```
