# Trelloknecht
A tool to print Trello cards on a label printer.

There was a python thingie around that basically did the same. That software is not maintained anymore.  So I wrote something replacing it. Since there is no direct printing in golang implemented so far, I call cups cli `lp` to print to a cups enabled printer. This works on linux and MacOS. I do not know if it works on more bizarre or exotic operating systems as well. Feel free to port it to whatever you want. 

## Basic Concept
This software connects to a trello board, gets all cards with a certain label. Creates a printjob for each matching card. Finally the tagging label will be removed. 

## Introducion
Mark the Trello card you want to be printed with a label e.g. "PRINTME". This software scans a list of boards, finds the cards with the label, prints them and replaces the label with a new label e.g. "PRINTED". 

## Requirements
- A Trello board
- A User with a valid trello API-Token. This user must be able to read/write to the trello board.
- The label marking a card to be printed, must exist. 
- Cups is installed and working. 
- A label Printer, we are using Brother QL XXX, is configuered on cups. 
- A system where go code can be compiled. 
- A "computer" running this software. 


## Installation 
```
go get github.com/heinrichgrt/trelloknecht
```

## Setup
- In Trello add a printing user to your organisation or use an existing technical user.
- Create Access token for this user if not already in place. Google will tell you how to achieve this goal. 
- Invite this user to all the Trello boards you want to print from.
- Choose and create a card label for the state: "To be printed" on every board you want to print from. 
- Create a file with the access token like:
``` 
cat .token
trelloAppKey=   yourappkeywithoutquotes
trelloToken=    yourtokenwithoutquotes
````
- Edit config.cfg to your needs. 
- Start the software. 
- Add the "PRINTME" to a card. 
- Wait until your label printer prints the label. 
- done. 

## Detailed Setup on a Raspberry Pi
We manage our printers with Rasperry Pis. This is cheap and reliable. They can act as build or service monitors while running the trello printer soft as well. The raspi is fine but not really a number cruncher. Therefore we build the software on a different platform. In my case it is a mac. 

``` 
go get github.com/heinrichgrt/trelloknecht
cd $GOPATH/src/github.com/heinrichgrt/trelloknecht
go get -d ./...
go build 

# Create a .token file with your access key and token:
cat .token
trelloAppKey=asdfasdfasdfasfasddfasdfasdf
trelloToken=hjklhjklhjklhjklhjklhjklhjklhjklhjklhjklhjklhjklhjklhjklhjkl

# edit config.cfg to point to your board
# Start the software:
./trelloknecht -configfile config.cfg -tokenfile .token 
```
That should be it. I will add some more detailed information later.

have fun

# Issues
- Textrendering is done by the pdf-lib. This will sometimes look funny, espacially if a headline is too long
- The exec/eval of the cups print command needs some improvement