## MAVPEAK 

Simple monitoring interface for mavryk bakers.

### Setup

Right now there are two supported ways to run mavpeak:
- as module of mavbake
- as a standalone server with ami

#### mavbake

Since mavbake 0.13.0-alpha the mavpeak is natively supported module. You can setup it in 3 simple steps:
1. `mavbake setup --peak`
2. Adjust configuration as needed
    - NOTE: *You can autodetect basic setups with `mavpeak -root-dir <path to where are apps stored> -autodetect-configuration <config file name>`*
3. `mavbake start --peak`

Sample minimal configuration:
```hjson
{
    listen: 0.0.0.0:8733
    app_root: /mavbake
    modules: {
        mavbake: {
            bakers: [
                tz1S5WxdZR5f9NzsPXhr7L9L1vrEb5spZFur
            ]
        }
        mavpay: {
            payout_wallet: tz1X7U9XxVz6NDxL4DSZhijME61PW45bYUJE
        }
    }
}
```

NOTE: you can use only modules you prefer, for example if you don't want to use mavpay, you can remove it from the configuration.
NOTE 2: mavpay should be 0.17.0 or higher to provide proper experience.

#### Standalone

You can run mavpeak as a standalone server as a binary or on linux as a service with [ami-mavpeak](https://github.com/mavryk-network/ami-mavpeak).

##### As binary

1. Download the latest release from the [releases page](https://github.com/mavryk-network/mavpeak)
2. Add configuration file `config.hjson` to the same directory as the binary
- Sample standalone minimal configuration:
```hjson
{
    listen: 0.0.0.0:8733
    app_root: /mavbake
    modules: {
        mavbake: {
            applications: null
            bakers: [
                tz1S5WxdZR5f9NzsPXhr7L9L1vrEb5spZFur
            ]
        }
		# mavpay requires ami-mavpay package, it is used to run automatic and manual payouts
        mavpay: {
            applications: null
            payout_wallet: tz1X7U9XxVz6NDxL4DSZhijME61PW45bYUJE
        }
    }
}
```
3. Run the binary
4. Open the browser and navigate to `http://localhost:8733`
5. Enjoy

##### As an ami based service

Refer to the [ami-mavpeak readme](https://github.com/mavryk-network/ami-mavpeak) for the installation and usage instructions.

### Advanced Configuration

```hjson
{
	# Id to show in the header
    id: ""
	# Address to listen on
    listen: 127.0.0.1:8733
    app_root: /mavbake
    modules: {
        mavbake: {
			# uncomment bellow to disable mavbake package monitoring
            # applications: null
            bakers: [
				# list of bakers to monitor for balances and rights
                tz1P6WKJu2rcbxKiKRZHKQKmKrpC9TfW1AwM
            ]
        }
        mavpay: {
			# can be null to disable mavpay package monitoring
            applications: {
				# path to mavpay ami package, either absolute or relative to parent directory peak
                mavpay: mavpay
            }
            payout_wallet: tz1X7U9XxVz6NDxL4DSZhijME61PW45bYUJE
            payout_wallet_preferences: {
                balance_warning_threshold: 100
                balance_error_threshold: 50
            }
			# forces all operations to be dry run
            force_dry_run: true
        }
    }
	
	# List of reference nodes to connect to
	# The reference nodes are used to get the rights and blocks if the baker's node is not available
    nodes: {
        "Mavryk Dynamics": {
            address: https://rpc.mavryk.network/
            is_rights_provider: true
            is_block_provider: false
        }
        mvkt: {
            address: https://mainnet.rpc.mavryk.network/
            is_rights_provider: false
            is_block_provider: true
	        # reports error if node not available, use for baker's node
            is_essential: false
        }
    }
	# The mode mavpeak should operate in
	# auto - if bound to localhost, it will operate in private mode if not, it will operate in public mode
	# public - assumes public environment, only readonly operations are allowed
	# private - assumes private environment, all operations are allowed
    mode: auto
}
``` 