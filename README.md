# sccapi

This program queries suse customer center about information like registration codes, products, systems, repositories urls etc. by using scc rest api v4.

The current version of this programm supports all GET requests for /organization/... 

To find more about rest api of scc go here: [SUSE Connect API](https://scc.suse.com/connect/v4/documentation "SUSE Customer Center rest api")

## benefits

* Customer of SUSE could query some information from SUSE Customer Center quickly without the need to search for individual user login and password.
* Often Customer's admins just need to find registration code for certain product
* No need to use web browser exploring scc.suse.com
* Admins could create further reports and send data into internal places or via email as report.
* Information like subscribed products, systems can be queried.

## Usage:

Download the compiled binary and put it to your favorit local directory which is in environment PATH.

```
wget https://github.com/bjin01/sccapi/raw/main/scc -O /usr/local/sbin/scc
```
add execute permission to the downloaded binary:
```
chmod +x /usr/local/sbin/scc
```

Create a config file where you save your scc mirror credentials. You could find it in your scc.suse.com portal under "proxies".
The file has this entries in yaml format.
```
user_name: UC2000880 
password: 5xa9b5cu6
```

Now you can run it by giving the scc credentials stored in your config file and provide parameter -get and information you want.

```
# scc -config scccred.yaml -get subscriptions
```
