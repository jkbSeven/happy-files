# happy-files
## client (draft)
To sign up use `hf signup`

To listen for a file transfer use `hf listen [-y]`, if -y then auto accept every file transfer  

To send a file use `hf send USERNAME FILE`, where USERNAME is a username of the receiver and FILE is a path to the file  

To see history of transfered files use `hf history [--received|--sent]`, by default prints received and sent files to stdout  

To enable/disable whitelist (only users from this list can communicate with you) use `hf whitelist -e|-d`; to add user(s) use `hf whitelist -a USERNAME...`; to remove user(s) use `hf whitelist -r USERNAME...`; to clear the whitelist use `hf whitelist -c`  

To enable/disable blacklist (users from this list can NOT communicate with you) use `hf blacklist -e|-d`; to add user(s) use `hf blacklist -a USERNAME...`; to remove user(s) use `hf blacklist -r USERNAME...`; to clear the blacklist use `hf blacklist -c`  

To make changes in config use `hf config [-hsd] [-u USERNAME] [-e EMAIL] [-k PRIVATE_KEY] [--server-ip IP] [--server-port PORT] [--server-public PUBLIC_KEY] [--download-path PATH]`,  
-u Specify your username  
-e Specify your email  
-k Specify your private key  
-h Store history of file transfers  
-s Don't store history of file transfers (secret)  
-d Restore config to deafult

If you make changes to the configuration while happy-files client is already running, you need to restart it - changes aren't automatically applied.  


## Roadmap
1. Basic peer2peer file transfer ✓
    * basic message structure with field size prefixes ✓
    * verifying message codes ✓
    * listening for transfer and updating server with client's currently opened port ✓
    * sending a stream of bytes from file ✓
3. Full encryption and identity verifaction
    * update structure of messages to include signatures, session keys etc.
    * implementation of chacha20
    * implementation of poly1305
    * RSA signatures and transfer handshake
4. Improve error handling
5. Server side database with users' data
6. Logging and collecting transfer history
7. UPnP (+ TCP hole punching?)
    * add entry to router's NA(P)T table with client's current listening port
8. Basic CLI
    * implement commands listed in client section above
    * implement commands for launching and configuring the server
9. TUI with BubbleTea
    * provide more user friendly experience than regular cli
