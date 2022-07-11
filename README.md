# Fictitious Scripting Language (FSL) 

The goal of this program is to create a basic system that can process
our fictitious scripting language (FSL). FSL is written in JSON.

This system receives FSL as input. The FSL defines variables and functions.
Functionality is limited to create, delete, update, add, subtract, multiply,
divide, print, as well as function calls. Variables are always numeric.

## Requirements

The finished project must support receiving multiple FSL scripts.
Functions and variables must persist between FSL scripts.
Resolve conflicts by overwriting existing variables or functions.

The system will create a representation of the script processed.
The init function is immediately called after each script is processed.

The input is a JSON object of named variables and named functions.
Variables are defined as a key value pair.
References to variables are preceded by a hash mark (#).

A function is an array of command objects.
An attribute called “cmd” is required and will define which operation to perform.
All parameters passed to a function are referenced by a $.

Function calls are defined in the “cmd” attribute by preceding the function name
with a hash mark (#).

## Sample Script

```
{
  "var1":1,
  "var2":2,
  
  "init": [
    {"cmd" : "#setup" }
  ],
  
  "setup": [
    {"cmd":"update", "id": "var1", "value":3.5},
    {"cmd":"print", "value": "#var1"},
    {"cmd":"#sum", "id": "var1", "value1":"#var1", "value2":"#var2"},
    {"cmd":"print", "value": "#var1"},
    {"cmd":"create", "id": "var3", "value":5},
    {"cmd":"delete", "id": "var1"},
    {"cmd":"#printAll"}
  ],
  
  "sum": [
      {"cmd":"add", "id": "$id", "operand1":"$value1", "operand2":"$value2"}
  ],

  "printAll":
  [
    {"cmd":"print", "value": "#var1"},
    {"cmd":"print", "value": "#var2"},
    {"cmd":"print", "value": "#var3"}
  ]
}
```
### Sample Output
```
3.5
5.5
undefined
2
5
```
## Installation and usage
This program requires go 1.18+. The program takes script files paths as program arguments.
For example, the following command can be used to run the program with sample scripts
which are included in the repository:
```
go run cmd/fslengine.go ./sample/sample-script1.txt ./sample/sample-script2.txt
```
The output of the program is going to look like this:
```
3.5000
5.5000
undefined
2.0000
5.0000
undefined
2.0000
5.0000
```