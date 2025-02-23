<h1 align="center">Chatanium</h1>
<p align="center">Dynamic Chatbot Runtime</p>

![Design Architecture](https://github.com/user-attachments/assets/cef88e8c-6689-4a08-8231-ec80d43bbea2)

Chatanium is a dynamic chatbot runtime.

You can create and load your own bot modules through various online chat providers such as Discord.<br/>
In addition, you can add providers that are not directly related to chatting, such as running a web interface for bot management.

## Design Architecture

### Support for Dynamic Module Insertion

It supports shared objects compiled based on Go plugins.<br/>

Before running the bot runtime, users can create a `modules` folder at the top level of the repository folder and place modules there to run.
In the future, we plan to support remote module insertion like gRPC to support various module insertion interfaces.

### Customizable Interface

To maximize compatibility, we did not use a unified interface between providers and modules.

Therefore, you can use existing online chat provider libraries (e.g., [bwmarrin/discordgo](https://github.com/bwmarrin/discordgo) for Discord) as is.<br/>
Currently, we have created a Discord backend using [bwmarrin/discordgo](https://github.com/bwmarrin/discordgo).

Check out the modules below to see how modules are actually constructed.

## Modules

See how modules should actually be created.

### for Discord

* **[thirdscam/chatanium-musicbot](https://github.com/thirdscam/chatanium-musicbot)**<br/>
Allows you to play songs in voice channels.

* **[thirdscam/chatanium-snowflake](https://github.com/thirdscam/chatanium-snowflake)**<br/>
It tells you the creation time of a Snowflake ID generated in Discord.
