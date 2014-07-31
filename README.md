chapelco-weather-goajs
======================

This is a Golang and AngularJS project. The first version of this website/project can be found at my other corresponding repository EntilZha/chapelco-weather-ruby. The first version is written in Ruby/Sinatra/jQuery.

I created this version as a way to learn Go and AngularJS, so it is my first project in both.

The website can be found at chapelco.heroku.com. It pulls data from a weather station at the mid-mountain of Chapelco Ski Resort in Argentina. The weather station is at 1700M at a place called Puesto Fijo. The data is saved to google drive as a public file where my application can query it periodically for updated data using a dbf parser (which I contributed a patch to in order to query from urls in addition to files).

In creating this project I have done my best to follow best practices, but there is still a lot to learn!
