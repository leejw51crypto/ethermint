import "package:flutter/material.dart";

class NodePage extends StatefulWidget {
  NodeState createState() => NodeState();
}

class NodeState extends State<NodePage> {
  Widget build(BuildContext context) {
    return Scaffold(
        appBar: AppBar(title: Text("Node Setup")),
        body: Container(color: Colors.blue, child: Text("OK")));
  }
}
