export const ServeConfig = [
    {
        "Name": "sftp",
        "Description": "SFTP server",
        "Prefix": "sftp",
        "Options": [
            {
                "Name": "user",
                "Help": "User name for authentication",
                "Provider": "",
                "Default": "",
                "Value": null,
                "ShortOpt": "",
                "Hide": 0,
                "Required": false,
                "IsPassword": false,
                "NoPrefix": false,
                "Advanced": false
            },
            {
                "Name": "pass",
                "Help": "Password for authentication",
                "Provider": "",
                "Default": "",
                "Value": null,
                "ShortOpt": "",
                "Hide": 0,
                "Required": false,
                "IsPassword": true,
                "NoPrefix": false,
                "Advanced": false
            },
            {
                "Name": "authorized-keys",
                "Help": "Authorized keys file (default \"~/.ssh/authorized_keys\")",
                "Provider": "",
                "Default": "",
                "Value": null,
                "ShortOpt": "",
                "Hide": 0,
                "Required": false,
                "IsPassword": false,
                "NoPrefix": false,
                "Advanced": true
            },
            {
                "Name": "gid",
                "Help": "Override the gid field set by the filesystem (not supported on Windows) (default 1000)",
                "Provider": "",
                "Default": "",
                "Value": null,
                "ShortOpt": "",
                "Hide": 0,
                "Required": false,
                "IsPassword": false,
                "NoPrefix": false,
                "Advanced": true
            },
            {
                "Name": "uid",
                "Help": "Override the uid field set by the filesystem (not supported on Windows) (default 1000)",
                "Provider": "",
                "Default": "",
                "Value": null,
                "ShortOpt": "",
                "Hide": 0,
                "Required": false,
                "IsPassword": false,
                "NoPrefix": false,
                "Advanced": true
            },
            {
                "Name": "dir-perms",
                "Help": "Directory permissions (default 777)",
                "Provider": "",
                "Default": "",
                "Value": null,
                "ShortOpt": "",
                "Hide": 0,
                "Required": false,
                "IsPassword": false,
                "NoPrefix": false,
                "Advanced": true
            },
            {
                "Name": "file-perms",
                "Help": "File permissions (default 666)",
                "Provider": "",
                "Default": "",
                "Value": null,
                "ShortOpt": "",
                "Hide": 0,
                "Required": false,
                "IsPassword": false,
                "NoPrefix": false,
                "Advanced": true
            }
        ]
    },
    {
        "Name": "nfs",
        "Description": "NFS with NFSv3 server",
        "Prefix": "nfs",
        "Options": [
            {
                "Name": "gid",
                "Help": "Override the gid field set by the filesystem (not supported on Windows) (default 1000)",
                "Provider": "",
                "Default": "",
                "Value": null,
                "ShortOpt": "",
                "Hide": 0,
                "Required": false,
                "IsPassword": false,
                "NoPrefix": false,
                "Advanced": true
            },
            {
                "Name": "uid",
                "Help": "Override the uid field set by the filesystem (not supported on Windows) (default 1000)",
                "Provider": "",
                "Default": "",
                "Value": null,
                "ShortOpt": "",
                "Hide": 0,
                "Required": false,
                "IsPassword": false,
                "NoPrefix": false,
                "Advanced": true
            },
            {
                "Name": "dir-perms",
                "Help": "Directory permissions (default 777)",
                "Provider": "",
                "Default": "",
                "Value": null,
                "ShortOpt": "",
                "Hide": 0,
                "Required": false,
                "IsPassword": false,
                "NoPrefix": false,
                "Advanced": true
            },
            {
                "Name": "file-perms",
                "Help": "File permissions (default 666)",
                "Provider": "",
                "Default": "",
                "Value": null,
                "ShortOpt": "",
                "Hide": 0,
                "Required": false,
                "IsPassword": false,
                "NoPrefix": false,
                "Advanced": true
            }
        ]
    },
    {
        "Name": "s3",
        "Description": "Amazon S3 Compliant Storage Providers including AWS, Minio, Ceph, Digital Ocean Spaces, and more",
        "Prefix": "s3",
        "Options": [
            {
                "Name": "auth-key",
                "Help": "Set key pair for v4 authorization: access_key_id,secret_access_key",
                "Provider": "",
                "Default": "",
                "Value": null,
                "ShortOpt": "",
                "Hide": 0,
                "Required": true,
                "IsPassword": false,
                "NoPrefix": false,
                "Advanced": false,
                "Repeatable": true
            },
            {
                "Name": "gid",
                "Help": "Override the gid field set by the filesystem (not supported on Windows) (default 1000)",
                "Provider": "",
                "Default": "",
                "Value": null,
                "ShortOpt": "",
                "Hide": 0,
                "Required": false,
                "IsPassword": false,
                "NoPrefix": false,
                "Advanced": true
            },
            {
                "Name": "uid",
                "Help": "Override the uid field set by the filesystem (not supported on Windows) (default 1000)",
                "Provider": "",
                "Default": "",
                "Value": null,
                "ShortOpt": "",
                "Hide": 0,
                "Required": false,
                "IsPassword": false,
                "NoPrefix": false,
                "Advanced": true
            },
            {
                "Name": "dir-perms",
                "Help": "Directory permissions (default 777)",
                "Provider": "",
                "Default": "",
                "Value": null,
                "ShortOpt": "",
                "Hide": 0,
                "Required": false,
                "IsPassword": false,
                "NoPrefix": false,
                "Advanced": true
            },
            {
                "Name": "file-perms",
                "Help": "File permissions (default 666)",
                "Provider": "",
                "Default": "",
                "Value": null,
                "ShortOpt": "",
                "Hide": 0,
                "Required": false,
                "IsPassword": false,
                "NoPrefix": false,
                "Advanced": true
            }
        ]
    },
    {
        "Name": "webdav",
        "Description": "Webdav server",
        "Prefix": "webdav",
        "Options": [
            {
                "Name": "user",
                "Help": "User name for authentication",
                "Provider": "",
                "Default": "",
                "Value": null,
                "ShortOpt": "",
                "Hide": 0,
                "Required": false,
                "IsPassword": false,
                "NoPrefix": false,
                "Advanced": false
            },
            {
                "Name": "pass",
                "Help": "Password for authentication",
                "Provider": "",
                "Default": "",
                "Value": null,
                "ShortOpt": "",
                "Hide": 0,
                "Required": false,
                "IsPassword": true,
                "NoPrefix": false,
                "Advanced": false
            },
            {
                "Name": "uid",
                "Help": "Override the uid field set by the filesystem (not supported on Windows) (default 1000)",
                "Provider": "",
                "Default": "",
                "Value": null,
                "ShortOpt": "",
                "Hide": 0,
                "Required": false,
                "IsPassword": false,
                "NoPrefix": false,
                "Advanced": true
            },
            {
                "Name": "dir-perms",
                "Help": "Directory permissions (default 777)",
                "Provider": "",
                "Default": "",
                "Value": null,
                "ShortOpt": "",
                "Hide": 0,
                "Required": false,
                "IsPassword": false,
                "NoPrefix": false,
                "Advanced": true
            },
            {
                "Name": "file-perms",
                "Help": "File permissions (default 666)",
                "Provider": "",
                "Default": "",
                "Value": null,
                "ShortOpt": "",
                "Hide": 0,
                "Required": false,
                "IsPassword": false,
                "NoPrefix": false,
                "Advanced": true
            }
        ]
    },
];