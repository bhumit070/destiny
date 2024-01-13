#!/bin/bash
shellFile=""
if [ -n "$($SHELL -c 'echo $ZSH_VERSION')" ]; then
   shellFile=".zshrc"
elif [ -n "$($SHELL -c 'echo $BASH_VERSION')" ]; then
   shellFile=".bashrc"
else
   echo "it is something else"
fi
platform=''
if [[ "$OSTYPE" == "linux-gnu"* ]]; then
   platform="linux-amd64"
elif [[ "$OSTYPE" == "darwin"* ]]; then
   if [[ $(sysctl -n machdep.cpu.brand_string) == *"Apple"* ]]; then
      platform="darwin-arm64"
   else
      platform="darwin-amd64"
   fi
else
   echo "something else"
fi
platform="destiny-$platform"
tagName=$(curl -s https://api.github.com/repos/bhumit070/destiny/releases/136918205 | awk -F'"tag_name":' '{print $2}' | awk -F'"' '{print $2}' | xargs)
downloadableUrl="https://github.com/bhumit070/destiny/releases/download/$tagName/$platform"
echo "downloading..."
installDir="$HOME/.destiny"
destinationPath="$installDir/destiny"
mkdir -p "$installDir"
curl -# -L "$downloadableUrl" -o "$destinationPath" && chmod +x "$destinationPath"

echo "PATH=$PATH:$installDir" >>"$HOME/$shellFile"
echo "installation completed."
