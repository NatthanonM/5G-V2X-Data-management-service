# 5G-V2X-Data-management-service

## Installation
### Step 1 - Install [Go](https://golang.org/) (version 1.15+)

### Step 2 - Install Make
**For Windows** 

Recommended use [Cygwin](https://www.cygwin.com/)

**For MacOS**

Need to install command line utilities

1. Open "Terminal" (it is located in Applications/Utilities)

2. In the terminal window, run the command ```xcode-select --install```

3. In the windows that pops up, click Install, and agree to the Terms of Service.

Once the installation is complete, the command line utilities should be set up property.

**For Linux (Debian, Ubuntu)**

1. ```sudo apt-get update```

2. ```sudo apt-get install build-essential```


### Step 3 - Install [Protobuf](https://github.com/protocolbuffers/protobuf/releases/) (version 3.12.0+)

3.1 - Set PATH Environment

3.2 - Go to this repository directory

3.2 - Install protoc-gen-go ```go get github.com/golang/protobuf/protoc-gen-go```

### Step 4 - Install dependencies

4.1 - Go to this repository directory

4.2 - run ```go mod vendor```

### Step 5 - Create generated code from Protobuf

5.1 - Go to this repository directory

5.2 - run ```make proto```

## Start project

run ```make start```
