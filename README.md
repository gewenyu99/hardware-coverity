## What is this?
This project in inspired by the embedded OS development I did at a prior internship.

The intent is to generate meaningful analytics on errors to discover other potential affected hardware, driver and firmware.

***None of the data used here is real :)***

## The idea
You can't run every test on every type of hardware when developing embedded software. When a test fails on a hardware target, we can deduct other hardware/software at risk with a reasonable degree of accuracy based off of the target's configuration. This serves as a proof of concept and demonstrates how something like this could be implemented.  

## What we did
As a demo, we are taking advantage of the infrastructure which is pre-existing, namely an ELK stack that is ever growing. The intent is to show the type of cool reports we can be generating from this,

## Why not GraphQL?
To be honest, graph based data types are perfect for this type of work, why we used ES and aggregations is simply to take advantage of the pre-existing infrasturecture.

## How do I try this out?
