## What is this?
This project in inspired by the embedded OS development I did at a prior internship. This repository demonstrates the high level concepts of a project I worked on during my tenure, with a very stripped down data set (for obvious reasons). 

The intent is to generate meaningful analytics on errors to discover other potential affected hardware, driver and firmware.

***None of the data used here is real :)***

## The idea
You can't run every test on every type of hardware when developing embedded software. When a test fails on a hardware target, we can deduct other hardware/software at risk with a reasonable degree of accuracy based off of the target's configuration. This serves as a proof of concept and demonstrates how something like this could be implemented.  

## What we did
As a demo, we are taking advantage of the infrastructure which is pre-existing, namely an ELK stack that is ever growing. The intent is to show the type of cool reports we can be generating from this,

## Why not GraphQL?
To be honest, graph based data types are perfect for this type of work, why we used ES and aggregations is simply to take advantage of the pre-existing infrasturecture.

This demo must show the viability of this concept with existing dev power, remapping all our data from ES to somewhere else was not something our teams would have bought.
## How do I try this out?
Pull and ```run docker-compose up --build``` :)

This was used for a short little demo, it will auto import some test data, and we used POSTMAN to make post calls.


## Credits
I found an incredibly helpful jist here for json requests in GO: https://gist.github.com/Tinker-S/52ae0f981d7b86e0b34f
