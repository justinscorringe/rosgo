package main

// Disable code generation, since we're trying to void doing this:   //go:generate gengo msg std_msgs/String

// IMPORT REQUIRED PACKAGES.

// TODO - Why is the syntax for import different to everywhere else?

import (
	"fmt"
	"github.com/edwinhayes/rosgo/ros"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// DEFINE PUBLIC STRUCTURES.

// DEFINE PRIVATE STRUCTURES.

// DEFINE PUBLIC GLOBALS.

// DEFINE PRIVATE GLOBALS.

var g_node ros.Node

var subscribers map[string]ros.Subscriber

var pub ros.Publisher

// DEFINE PUBLIC STATIC FUNCTIONS.

// DEFINE PRIVATE STATIC FUNCTIONS.

func callback(msg *ros.DynamicMessage, event ros.MessageEvent) {
	g_node.Logger().Info("Received: ", event.ConnectionHeader["topic"], " : ", msg.Type().Name(), " : ", msg)

	// Try not to loopback into oblivion...
	if event.ConnectionHeader["topic"] == "/more_chatter" {
		return
	}

	// If the message was just a string, republish it.
	if msg.Type().Name() == "std_msgs/String" {
		g_node.Logger().Info("Republishing...")
		out_msg := ros.Message(msg)
		pub.Publish(out_msg)
	}
}

func poll_for_topics(node ros.Node, quit <-chan bool) {
	// Create a ticker to tell us to do stuff periodically.
	ticker := time.NewTicker(1 * time.Second)

	// Loop forever, or until rx on chan quit.
	node.Logger().Info("Starting goroutine to poll for topics...")
	for {
		select {
		case <-quit:
			node.Logger().Info("Stopping polling for topics...")
			return
		case <-ticker.C:
			// Fetch list of available topics (i.e. those with publishers) from the master.
			topic_list := node.GetPublishedTopics("")

			// Try to iterate over each of the topics in the list.
			for _, v := range topic_list {
				topic := v.([]interface{})
				topic_name := topic[0].(string)
				topic_type := topic[1].(string)

				// Check if we have a subscriber for this topic already.
				if _, ok := subscribers[topic_name]; !ok {
					// Apparently not, so we try to subscribe.
					node.Logger().Info("Attempting to subscribe to topic: ", topic_name)

					// Create a generic message, which tries to look up the important checksum via gengo.
					var m *ros.DynamicMessageType
					var err error
					if m, err = ros.NewDynamicMessageType(topic_type); err != nil {
						node.Logger().Info("Couldn't set message type: ", topic_type, " : Error: ", err)
						continue
					}

					// Then subscribe to the topic, and if we're successful, keep a note so we don't try to subscribe again.
					s := node.NewSubscriber(topic_name, m, callback)
					if s != nil {
						subscribers[topic_name] = s
					}
				}
			}
		}
	}
	// Not all done, since defer?
}

func main() {
	// Create our node.
	node, err := ros.NewNode("listener", os.Args)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	defer node.Shutdown()
	g_node = node

	// Configure node logging.
	node.Logger().SetSeverity(ros.LogLevelInfo)

	// We'll keep a list of ROS subscribers, so we can identify topics which we still need to subscribe to.
	subscribers = make(map[string]ros.Subscriber)

	// Spawn a routine to look for new topics which get published.
	quit_poll_for_topics := make(chan bool)
	go poll_for_topics(node, quit_poll_for_topics)
	defer close(quit_poll_for_topics)

	// Setup a signal handler to catch the keyboard interrupt.
	quit_mainloop := make(chan os.Signal, 2)
	signal.Notify(quit_mainloop, os.Interrupt, syscall.SIGTERM)

	// Create a publisher for rebroadcasting the messages we recieve.
	var out_msg *ros.DynamicMessageType
	if out_msg, err = ros.NewDynamicMessageType("std_msgs/String"); err != nil {
		node.Logger().Error("Oh noes!")
		return
	}
	pub = node.NewPublisher("more_chatter", out_msg)

	// Wait forever.
	node.Logger().Info("Spinning...")
	defer node.Logger().Info("Shutting down...")
	for {
		select {
		case <-quit_mainloop:
			node.Logger().Debug("Received SIGTERM.")
			return
		default:
			node.SpinOnce()
		}
	}

	// Not all done, since defer?
}

// ALL DONE.
