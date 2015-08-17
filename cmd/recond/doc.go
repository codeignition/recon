/*

recond is the daemon that runs on the machines.
It collects the system metrics such as percentage of CPU used, memory consumption, filesystem metrics, network metrics, etc.
It sends the update every 5 seconds to the marksman server which is the metrics aggregator server.

*/
package main
