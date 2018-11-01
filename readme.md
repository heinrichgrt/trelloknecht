# Trelloknecht
A tool to print Trello cards on a label printer.

There was a python thingie around that basically did the same. That software is not maintained anymore.  So I wrote something replacing it. Since there is no direct printing in golang implemented so far, I call lp to print to a cups enabled printer. This works on linux and MacOS. I do not know if it works on more bizarre or exotic operating software as well. Feel free to port it to whatever you want. 

## Introducion
Mark the Trello card you want to be printed with a label e.g. "PRINTME". This software scans a list of boards, finds the cards with the label, prints them and replaces the label with a new label e.g. "PRINTED". 

## Requirements
- A Trello Board
- Cups is installed and working. 
- A label Printer, we are using Brother QL XXX, is configuered on cups. 
- A go runtime environment 
- A "computer" running this software. 


## Installation 

go get github.com/heinrichgrt/trelloknecht


## Setup
- In Trello add a printing user to your organisation or use an existing technical user.
- Create Access token for this user if not already in place. Google will tell you how to achive this goal. 
- Invite this user to all the Trello boards you want to print from.
- Choose card labels for the state: "To be printed" and "Printed" and create them on every board. E.g. "PRINTME" and "PRINTED". 
- Create a file with the access token
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
```
Create a .token file with your access key and token:
```
cat .token
trelloAppKey=asdfasdfasdfasfasddfasdfasdf
trelloToken=hjklhjklhjklhjklhjklhjklhjklhjklhjklhjklhjklhjklhjklhjklhjkl
```
edit config.cfg to point to your board

Start the software:
./trelloknecht -configfile config.cfg -tokenfile .token 

That should be it. I will add some more detailed information later.

have fun

