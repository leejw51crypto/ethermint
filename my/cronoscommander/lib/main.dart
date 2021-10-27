import 'package:flutter/material.dart';
import "dart:io";
import 'dart:convert';

void main() {
  runApp(MyApp());
}

class MyApp extends StatelessWidget {
  // This widget is the root of your application.
  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'Flutter Demo',
      theme: ThemeData(
        primarySwatch: Colors.blue,
      ),
      home: MyHomePage(title: 'CRONOS COMMANDER'),
    );
  }
}

class MyHomePage extends StatefulWidget {
  MyHomePage({Key? key, required this.title}) : super(key: key);

  final String title;

  @override
  _MyHomePageState createState() => _MyHomePageState();
}

class _MyHomePageState extends State<MyHomePage> {
  int _counter = 0;

  String _chainid = "ethermint-2";
  String _keyname = "mykey";
  int _node1Pid = 0;

  void _incrementCounter() {
    setState(() {
      _counter++;
    });
  }

  void changeNode1(int pid) {
    setState(() {
      _node1Pid = pid;
    });
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: Text(widget.title),
      ),
      body: Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: <Widget>[
            Container(
                padding: EdgeInsets.all(8),
                child: Row(children: [
                  Expanded(
                      child: Column(children: [
                    TextFormField(
                        initialValue: _chainid,
                        decoration: const InputDecoration(
                          filled: true,
                          hintText: 'Enter a chainid2...',
                          labelText: 'chainid',
                        )),
                    TextFormField(
                        initialValue: _keyname,
                        decoration: const InputDecoration(
                          filled: true,
                          hintText: 'Enter a keyname...',
                          labelText: 'Keyname',
                        ))
                  ])),
                ])),
            Container(
              padding: EdgeInsets.all(8),
              child: ElevatedButton(
                style: ElevatedButton.styleFrom(
                  primary: Colors.blue, // background
                  onPrimary: Colors.white, // foreground
                ),
                onPressed: () {
                  print('chainid=$_chainid   keyname=$_keyname');
                },
                child: Text('Check2'),
              ),
            ),
            Container(
              padding: EdgeInsets.all(8),
              child: ElevatedButton(
                style: ElevatedButton.styleFrom(
                  primary: Colors.blue, // background
                  onPrimary: Colors.white, // foreground
                ),
                onPressed: () {},
                child: Text('Activate Node0'),
              ),
            ),
            Container(
              padding: EdgeInsets.all(8),
              child: Row(children: [
                Container(
                  color: Colors.yellow,
                  padding: EdgeInsets.all(8),
                  child: Text('node1 PID= $_node1Pid?'),
                ),
                Container(
                  padding: EdgeInsets.all(8),
                  child: ElevatedButton(
                    style: ElevatedButton.styleFrom(
                      primary: Colors.blue, // background
                      onPrimary: Colors.white, // foreground
                    ),
                    onPressed: () async {
                      if (_node1Pid != 0) {
                        var killresult =
                            Process.killPid(_node1Pid, ProcessSignal.sighup);
                        print('kill the proess $_node1Pid  result $killresult');
                        changeNode1(0);
                        return;
                      }
                      var process = await Process.start('ethermintd', ['start'],
                          //var process = await Process.start('timeout', ['10000'],
                          runInShell: false,
                          mode: ProcessStartMode.detachedWithStdio);
                      stdout.addStream(process.stdout);
                      stderr.addStream(process.stderr);
                      //await Future.delayed(Duration(seconds: 5));
                      var pid = process.pid;
                      print('process $pid created2');
                      changeNode1(process.pid);

                      //Process.killPid(process.pid);
                    },
                    child: Text('Activate Node1'),
                  ),
                ),
              ]),
            ),
          ],
        ),
      ),
      // This trailing comma makes auto-formatting nicer for build methods.
    );
  }
}
