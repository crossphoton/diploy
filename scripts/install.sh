echo Fetching binary
wget https://github.com/crossphoton/diploy/releases/download/v0.1.1/diploy
echo Done
echo Giving permissions
chmod +x diploy

echo Running setup
./diploy server setup