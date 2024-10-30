# ndoujin-cli

CLI tool to download doujins from nhentai.net website.

[Download](https://github.com/devyuji/ndoujin-cli/releases/latest)

## How to use

- Download the binary from [github](https://github.com/devyuji/ndoujin-cli/releases/latest) for your operating system.
- Open a terminal and go to the place where you downloaded the binary.
- Run the binary as

  ```bash
  <./ndoujin-cli-*> download <code / url> -p <folder location>
  ```

  If you don't specify the -p flag then it will download it in the same folder.

  You can also create a `config.json` file and set a download location.

  ```json
  {
    "path": "<download location>"
  }
  ```

  If you use `-p` flag and also have a `config.json` it will ignore the -p flag.

- To download in bulk Create a `code.txt` file and enter code one line at a time.

  ```bash
    534101
    533999
  ```

## If the website is using Cloudflare protection or any other then follow this steps

- Create a `config.json` file and add cookies, user-agent parameter to it

  ```json
  {
    "cookies": {
      "<cookies name>": "<cookie value>"
    },
    "user-agent": "<YOUR USER AGENT>"
  }
  ```

  You can get the user-agent from your browser just search `"what is my user agent"` and paste it on config.json file.

  **_Just make sure the cookies and user agent are from the same browser only._**
