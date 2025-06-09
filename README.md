# xk6-ansible-config-loader

A k6 extension to get access to your ansible variables in tests. Built for [k6](https://github.com/grafana/k6) using [xk6](https://github.com/grafana/xk6).

## Build

To build a `k6` binary with this extension, first ensure you have the prerequisites:

- [Go toolchain](https://go101.org/article/go-toolchain.html)
- Git

Then:

1. Download `xk6`:

  ```bash
  go install github.com/grafana/xk6/cmd/xk6@latest
  ```

2. Build the binary:

  ```bash
  xk6 build --with github.com/grafana/xk6-ssh@latest
  ```

This will result in a `k6` binary in the current directory.

## Example

To use this extension you need to provide a config, see [extension-config.yaml](./test/testdata/extension-config.yaml).  
For each ansible inventory you need to provide a new config file. Also you are just able to provide one vault password per inventory.  
Encrypted variables will be returned decrypted!  

Provide the path to your config file via environment variable.  

```bash
./k6 -e CONFIG_PATH=./test/testdata/extension-config.yaml run ./test/ansible.test.js
```

In your test file, call the `getConfig` method with the path to your config file.  

```typescript
  const config = ansibleConfigLoader.getConfig(__ENV.CONFIG_PATH);
```

You will recieve an object that contains all your data in following structure:  

```json
"group_config":[
   {
      "group_name":"loadbalancer",
      "hosts":[
         {
            "host_name":"lb-aaa.example.com",
            "host_vars":{
               
            }
         },
         {
            "host_name":"lb-bbb.example.com",
            "host_vars":{
               
            }
         }
      ],
      "group_vars":{
         "lb_url":"example.com",
         "example":"loadbalancer"
      }
   } 
 ],
 "global_config": {
  "plain_vars_file": "testtest",
  "test_test": "no_tests",
  "top": "secret",
  "secret": "top"
 }
```

## Build for Development

1. Get dependencies
- `go mod tidy`

2. Build binary
- `xk6 build --with xk6-ansible-config-loader=.`

## Testing Locally

`./test/exec-tests.sh`
