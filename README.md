# envaws [![Build Status](http://img.shields.io/travis/gpestana/envaws.svg?style=flat-square)](http://travis-ci.org/gpestana/envaws)  ![Release Version](https://img.shields.io/badge/release-0.1-blue.svg)


Launch a subprocess with environment variables using data from AWS S3 or AWS
parameter store.

Envaws provides a convenient way to launch a subprocess with environment 
variables populated from AWS services S3 and SSM. The tool is inspired by 
[envconsul](https://github.com/hashicorp/envconsul).

Envaws supports [12-factor applications](https://12factor.net/) which get their 
configuration via the environment. Environment variables are dynamically 
populated from S3 or SSM, but the application is unaware; applications just read 
environment variables. This enables extreme flexibility and portability for 
applications across systems.

## Installation

### Pre-compiled

1) Download the binary

```
$ curl -so envaws https://raw.githubusercontent.com/gpestana/envaws/master/bin/envaws
```

2) Make binary executable

```
$ chmod 755 envaws
```

3) Move the binary into your `$PATH`.

```
$ mv envaws /usr/local/bin/envaws
$ chmod +x /usr/local/bin/envaws
```

## Quick Example

## Usage

## Command Line Interface (CLI)

## Debugging

## Contributing



