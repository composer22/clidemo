### [Dockerized] (http://www.docker.com) [clidemo](https://registry.hub.docker.com/u/composer22/clidemo/)


A docker image for clidemo. This is created as a single "static" executable using a lightweight image.

To make:

cd docker
./build.sh

Once it completes, you can run the server:
```
docker run --name <containername> -p <hosthttpport>:8080 -p <hostpprofport>:6060 -d composer22/clidemo
```
Example 8081 and 6061 are host ports e.g. OSX; the others are in the container
We also alter attrs by passing addition params.  Here below, the name of the server is changed from NoName to SanFrancisco:
```
docker run --name tester -p 8081:8080 -p 6061:6060 -d composer22/clidemo -N SanFrancisco
```
NOTE:  If you are using boot2docker on OSX you also need to map these ports in VirtualBox
if you want to access the server from Terminal or your tools. The easiest way is to launch the OSX app, navigate to the bootdocker-vm, open the settings. From there, select network/adaptor1/advanced/port forwarding. Add an entry such as this:

Example:

name: API Server
protocol: TCP
host IP: 127.0.0.1  // this is OSX
host port: 8081   // this is OSX
guest IP: nil
guest port: 8080

save, save

You can also do this in while the VM for boot2docker is running in Virtualbox.

Also make sure your /etc/hosts file in OSX has an entry to 127.0.0.1 via localhost if you use this localhost as a name.

#### Options

For additional unix tools in a small image use "FROM busybox" instead of "FROM scratch" in Dockerfile.final
