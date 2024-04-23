## Samsung MagicInfo Server Address Changer

This program allows the remote and bulk (or individual) change of the MagicInfo server address of Samsung displays.

### Usage:
The usage is very simple:

1. Create a connection data file based on the `dumb.env` example file.
2. Create a file with the MAC addresses of the displays. The MAC address format is quite flexible, examples can be found in the `mac_address_example.txt` file.
3. Run the binary file suitable for your system from the `bin` folder, remembering to add the path to the .env file and the list of MAC addresses as arguments.

    For example, if you have such a file and folder structure:
    ```txt
    |_project
      |_.env
      |_mac.txt
      |_bin
        |_magicinfo_change_url_mac
        |_magicinfo_change_url_linux
        |_magicinfo_change_url_win_x64.exe
    ```
    Then the program launch will look like this<br>
    `./bin/magicinfo_change_url_mac -p mac.txt`

4. If you need to change the address for only one MAC address, specify it with the -a argument <br>`./bin/magicinfo_change_url_mac -a 1234567890AB`

## License

This project licensed under the MIT License - see the LICENSE.md file for more details.

## Contact Information

If you have any questions or problems, please let us know via email at [rptl8sr@gmail.com](mailto:rptl8sr@gmail.com).
