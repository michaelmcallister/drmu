# DrMu

the *D*ynamic *R*oute53 *M*ikroTik *U*pdater 

## Description

This utility creates a webserver, that accepts requests to the URL "/drmu/update/:domainName/:ipAddress"

Where `:domainName` is the subdomain name (eg. host1.example.com), and `:ipAddress` is the IP address to point to. Subdomains that don't exist, will be created, and existing ones updated.

## Getting Started

### Installing

* git clone https://github.com/michaelmcallister/drmu
* cd drmu
* go build cmd/drmu.go

### Executing program

Ensure that your environment has AWS credentials setup correctly, see the following
resource for further information.

* https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-configure.html

Ensure you pass the `-zone` flag with the hosted zone ID you wish to change.

## Usage
```
$./drmu --help
Usage of ./drmu:
  -address string
        address to bind on (default "localhost:8000")
  -alsologtostderr
        log to standard error as well as files
  -log_backtrace_at value
        when logging hits line file:N, emit a stack trace
  -log_dir string
        If non-empty, write log files in this directory
  -logtostderr
        log to standard error instead of files
  -stderrthreshold value
        logs at or above this threshold go to stderr
  -ttl int
        TTL for each record created in Route53 (default 300)
  -v value
        log level for V logs
  -vmodule value
        comma-separated list of pattern=N settings for file-filtered logging
  -zone string
        Route53 Hosted Zone ID (default "HOSTED_ZONE_ID")
```

## Integrating with MikroTik

Add [the script]("mikrotik.txt) in Mikrotik via `/system script add`, and use the scheduler 
to trigger it, see the following resources:

* https://wiki.mikrotik.com/wiki/Manual:System/Scheduler
* https://wiki.mikrotik.com/wiki/Manual:Scripting

## License

This project is licensed under the MIT License - see the LICENSE file for details