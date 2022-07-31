# custom-tcp

run TCP server:
   go run .\main.go

execute TCP client to interact with the TCP server:
   go run .\main.go -connect localhost

send the STOP command to exit the TCP client and server:
   STOP
  
send the STRING command followed by a message to send it to the TCP server:
   STRING Hello :)

send the IMAGEGOB command followed by an image path to send it to the TCP server using GOB:
   IMAGEGOB "path/image.jpg"
   
