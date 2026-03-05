# AWS-DR-Audit

![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)
![Open Source](https://img.shields.io/badge/Open%20Source-%E2%9D%A4-brightgreen)
![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go)
![AWS SDK](https://img.shields.io/badge/AWS_SDK-v2-FF9900?logo=amazonaws)

> Free, open-source CLI tool to scan your AWS account for Disaster Recovery gaps before a regional failure finds them for you.

## Why This Tool Exists

Cloud regions fail. emerging-market AWS regions carry concentrated infrastructure risk. Organizations that haven't implemented cross-region DR are one outage away from catastrophic data loss.

`aws-dr-audit` is 100% free and open-source. It tells you exactly where your DR posture is broken, before AWS breaks it for you.

## Checks

| Service | Check | What It Validates |
|---------|-------|-------------------|
| S3 | Public Access Block | All 4 block settings are enabled |
| S3 | Versioning | Object versioning is active |
| S3 | Cross-Region Replication | CRR rule copies objects to fallback region |
| RDS | Multi-AZ | Standby replica in separate AZ |
| RDS | Cross-Region Snapshots | Backups copied to fallback region |

## Install
```bash
git clone https://github.com/Mohamed-M-Meth/aws-dr-audit.git
cd aws-dr-audit
go build -o aws-dr-audit .
```

## Usage
```bash
# Full audit
./aws-dr-audit audit --region me-central-1 --fallback eu-central-1

# S3 only
./aws-dr-audit s3 --region me-central-1 --fallback eu-central-1

# RDS only
./aws-dr-audit rds --region me-central-1 --fallback eu-central-1

# JSON output
./aws-dr-audit audit --output json
```

## Permissions Required
```json
{
  "Effect": "Allow",
  "Action": [
    "s3:ListAllMyBuckets",
    "s3:GetBucketVersioning",
    "s3:GetReplicationConfiguration",
    "s3:GetBucketPublicAccessBlock",
    "rds:DescribeDBInstances",
    "rds:DescribeDBSnapshots"
  ],
  "Resource": "*"
}
```

## Contributing

Contributions are welcome! This project is open-source and free for everyone.
```bash
# Fork the repo, then:
git checkout -b feature/your-feature
git commit -m "feat: add your feature"
git push origin feature/your-feature
```

## Tech Stack

- Language: Go 1.21+
- AWS SDK: aws-sdk-go-v2
- CLI Framework: spf13/cobra
- Auth: Environment variables / AWS profiles / IAM roles

## Author

Built by Mohamed M. Meth | Powered by AI 🤖

- GitHub: [@Mohamed-M-Meth](https://github.com/Mohamed-M-Meth)

This tool is a personal learning project built while exploring AWS automation. As I am currently focused on my RHCSA and Ansible certifications, I welcome feedback and contributions from the community to improve its logic


## License

MIT License — Free to use, modify, and distribute.
