# happy-files
## client (draft)
To sign up use `hf signup`
To listen for a file transfer use `hf listen [-y]`, if -y then auto accept every file transfer  
To send a file use `hf send USERNAME FILE`, where USERNAME is a username of the receiver and FILE is a path to the file  
To see history of transfered files use `hf history [--received|--sent]`, by default prints received and sent files to stdout  
To enable/disable whitelist (only users from this list can communicate with you) use `hf whitelist -e|-d`; to add user(s) use `hf whitelist -a USERNAME...`; to remove user(s) use `hf whitelist -r USERNAME...`; to clear the whitelist use `hf whitelist -c`  
To enable/disable blacklist (only users from this list can communicate with you) use `hf blacklist -e|-d`; to add user(s) use `hf blacklist -a USERNAME...`; to remove user(s) use `hf blacklist -r USERNAME...`; to clear the blacklist use `hf blacklist -c`  
To make changes in config use `hf config [-hs] [-u USERNAME] [-k PRIVATE_KEY] [--server-ip IP] [--server-port PORT] [--server-public PUBLIC_KEY] [--download-path PATH] [default]`,  
-u Specify your username  
-e Specify your email  
-k Specify your private key  
-h Store history of file transfers  
-s Don't store history of file transfers (secret)  
To restore your config to default use `hf config default` (...paste default config below...) 
If you make changes to the configuration while happy-files client is already running, you need to restart it - changes aren't automatically applied.  
