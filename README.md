# Github Resource Pack Downloader
This is a simple script I have made to automatically download a specified file from the latest release of a Github Repository with the intention of being used for Minecraft Resource Packs.

## How to use
After the script has been installed to your computer's default terminal path (see below for how to do that), it can be run with the following command:

```bash
githubpackdownloader <repo> <filename> <destination path>
```
![Screenshot From 2025-04-09 21-49-07](https://github.com/user-attachments/assets/7af9953b-087f-4c38-afc2-e70a32be0264)

`<filename>` must match the file you want to download from the release assets.

#### Example in terminal

```bash
githubpackdownloader itslilscorp/MCParks-Resource-Pack-Updated mcparkspack-1.21.zip /home/jishy/.local/share/PrismLauncher/instances/1.21.1/minecraft/resourcepacks/
```
![Screenshot From 2025-04-09 21-52-30](https://github.com/user-attachments/assets/0ff7e743-9db7-4794-bbf7-064cd94280db)

#### Example in Prism Launcher
Using this script in Prism Launcher is very useful, as you can set it to run automatically as a pre-launch command for any instance(s) of your choosing.
![Screenshot From 2025-04-09 21-55-21](https://github.com/user-attachments/assets/4464ab5a-5253-48c2-b28e-09f2f4d2a292)

## How to install
This will be a basic tutorial on installing depending on your computer's operating system if needed.
*No Go Installation is required*

### Windows:
1. Create a path directory if it doesn't already exist
```bash
mkdir C:\bin
```
2. Move the downloaded from the releases tab into that folder.
3. Add the script to your PATH
   - Go to **Advanced system settings**
   - Click **Environment Variables**
   - Under **System Variables**, select **Path**, then **Edit**
   - Click **New** and add in `C:\bin`
   - Click **OK** for all the boxes.
4. Verify it is working by opening a command-prompt window and running
```bash
githubpackdownloader
```
It should respond with a message clarifying how to use the command. 

### MacOS and Linux:
1. Move the file downloaded from the releases tab into your PATH directory
```bash
sudo mv <file> /usr/local/bin/
```
2. Verify it is working by opening a command-prompt window and running
```bash
githubpackdownloader
```
