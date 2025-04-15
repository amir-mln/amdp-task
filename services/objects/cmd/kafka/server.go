// this package would be responsible for polling the outbox table
// and publishing the related events, commands, ... . It will also listen
// to incoming messages and call the appropriate handler.
package kafka
