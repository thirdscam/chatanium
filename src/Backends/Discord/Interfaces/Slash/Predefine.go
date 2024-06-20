package Slash

// pre-defined commands: can use slash interface without this definition
// using module system. but if you need to add runtime built-in command, maybe help this feature.
// also predefined commands have not affected by ACL, have top priority.
// (eg. if command name is same from module-defined and pre-defined command, runtime use pre-defined commands)
var PREDEFINED_COMMANDS = Commands{
	// Schema: Handler,
}
