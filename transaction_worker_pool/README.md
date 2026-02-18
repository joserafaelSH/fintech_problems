# Transaction Worker Pool Study

A simple app to consume from a RabbitMQ queue and send for the
go routine (Wokers) the payload to be processed.
The Goal of this project is uderstand how to use channels, 
infinity consumer and the usage of select statements with channels.

Receive message from RabbitMQ -> InputChann -> Process -> Result chann -> Save the result of transaction 
                                                      |-> Error chann -> Place to log or retry (not implemented)

To avoid mutex/locks/races, I choose a single point to write in the database.  