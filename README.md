# cinderlcient-grpc
 A grpc server to execute `cinder local-attach`

 The goal here is to run a container with cinderlcient in a sidecare container
 to say a CSI Node Server.  The CSI Node Server could then issue the attach
 call to the sidecar (running this GRPC server), and thus attach the volume
 to the K8s node and pass back the path to the CSI Node Server.

 It's a cheap and dirty way to implement a whole bunch of Cinder connectors in
 one shot and use stand-alone Cinder pretty easy in K8s.


